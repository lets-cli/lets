package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lets-cli/lets/internal/cmd"
	"github.com/lets-cli/lets/internal/config"
	"github.com/lets-cli/lets/internal/env"
	"github.com/lets-cli/lets/internal/executor"
	"github.com/lets-cli/lets/internal/logging"
	"github.com/lets-cli/lets/internal/set"
	"github.com/lets-cli/lets/internal/settings"
	"github.com/lets-cli/lets/internal/upgrade"
	"github.com/lets-cli/lets/internal/upgrade/registry"
	"github.com/lets-cli/lets/internal/workdir"
	"github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const updateCheckTimeout = 3 * time.Second

type updateCheckResult struct {
	notifier *upgrade.UpdateNotifier
	notice   *upgrade.UpdateNotice
}

func Main(version string, buildDate string) int {
	ctx := getContext()

	configDir := os.Getenv("LETS_CONFIG_DIR")

	appSettings, err := settings.Load()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "lets: settings error: %s\n", err)
		return 1
	}

	appSettings.Apply()

	logging.InitLogging(os.Stdout, os.Stderr)

	rootCmd := cmd.CreateRootCommand(version, buildDate)
	rootCmd.InitDefaultHelpFlag()
	rootCmd.InitDefaultVersionFlag()
	reinitCompletionCmd := cmd.InitCompletionCmd(rootCmd, nil)
	cmd.InitSelfCmd(rootCmd, version)
	rootCmd.InitDefaultHelpCmd()

	command, args, err := rootCmd.Traverse(os.Args[1:])
	if err != nil {
		log.Errorf("traverse commands error: %s", err)
		return getExitCode(err, 1)
	}

	rootFlags, err := parseRootFlags(args)
	if err != nil {
		log.Errorf("parse flags error: %s", err)
		return 1
	}

	if rootFlags.version {
		if err := cmd.PrintVersionMessage(rootCmd); err != nil {
			log.Errorf("print version error: %s", err)
			return 1
		}

		return 0
	}

	debugLevel := env.SetDebugLevel(rootFlags.debug)

	if debugLevel > 0 {
		log.SetLevel(log.DebugLevel)
	}

	if rootFlags.config == "" {
		rootFlags.config = os.Getenv("LETS_CONFIG")
	}

	cfg, err := config.Load(rootFlags.config, configDir, version)
	if err != nil {
		if failOnConfigError(rootCmd, command, rootFlags) {
			log.Errorf("config error: %s", err)
			return 1
		}
	}

	if cfg != nil {
		reinitCompletionCmd(cfg)
		cmd.InitSubCommands(rootCmd, cfg, rootFlags.all, os.Stdout)
	}

	if rootFlags.init {
		wd, err := os.Getwd()
		if err == nil {
			err = workdir.InitLetsFile(wd, version)
		}

		if err != nil {
			log.Errorf("can not create lets.yaml: %s", err)
			return 1
		}

		return 0
	}

	showUsage := rootFlags.help || (command.Name() == "help" && len(args) == 0) || (len(os.Args) == 1)

	if showUsage {
		if err := cmd.PrintRootHelpMessage(rootCmd); err != nil {
			log.Errorf("print help error: %s", err)
			return 1
		}

		return 0
	}

	updateCh, cancelUpdateCheck := maybeStartUpdateCheck(ctx, version, command, appSettings)
	defer cancelUpdateCheck()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if depErr, ok := errors.AsType[*executor.DependencyError](err); ok {
			log.Errorf("%s", depErr.TreeMessage())
			log.Errorf("%s", depErr.FailureMessage())

			return getExitCode(err, 1)
		}

		log.Errorf("%s", err.Error())

		return getExitCode(err, 1)
	}

	printUpdateNotice(updateCh)

	return 0
}

// getContext returns context and kicks of a goroutine
// which waits for SIGINT, SIGTERM and cancels global context.
//
// Note that since we setting stdin to command we run, that command
// will receive SIGINT, SIGTERM at the same time as we here,
// so command's process can begin finishing earlier than cancel will say it to.
func getContext() context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-ch
		log.Printf("signal received: %s", sig)
		cancel()
	}()

	return ctx
}

func getExitCode(err error, defaultCode int) int {
	type errorWithExitCode interface {
		error
		ExitCode() int
	}

	if errWithExitCode, ok := errors.AsType[errorWithExitCode](err); ok {
		return errWithExitCode.ExitCode()
	}

	return defaultCode
}

// do not fail on config error if it is help (-h, --help), --init, completion, or lets self.
func failOnConfigError(root *cobra.Command, current *cobra.Command, rootFlags *flags) bool {
	return (root.Flags().NFlag() == 0 && !allowsMissingConfig(current)) && !rootFlags.help && !rootFlags.init
}

func allowsMissingConfig(current *cobra.Command) bool {
	if current == nil {
		return false
	}

	switch current.Name() {
	case "completion", "help":
		return true
	}

	return isSelfCommand(current)
}

func isSelfCommand(current *cobra.Command) bool {
	for cmd := current; cmd != nil; cmd = cmd.Parent() {
		parent := cmd.Parent()
		if cmd.Name() == "self" && parent != nil && parent.Name() == "lets" {
			return true
		}
	}

	return false
}

func maybeStartUpdateCheck(
	ctx context.Context,
	version string,
	command *cobra.Command,
	appSettings settings.Settings,
) (<-chan updateCheckResult, context.CancelFunc) {
	if !shouldCheckForUpdate(command, isInteractiveStderr(), appSettings) {
		return nil, func() {}
	}

	log.Debugf("start update check")

	notifier, err := upgrade.NewUpdateNotifier(registry.NewGithubRegistry())
	if err != nil {
		return nil, func() {}
	}

	ch := make(chan updateCheckResult, 1)
	checkCtx, cancel := context.WithTimeout(ctx, updateCheckTimeout)

	go func() {
		notice, err := notifier.Check(checkCtx, version)
		if err != nil {
			upgrade.LogUpdateCheckError(err)
		}

		log.Debugf("update check done")

		ch <- updateCheckResult{
			notifier: notifier,
			notice:   notice,
		}
	}()

	return ch, cancel
}

func printUpdateNotice(updateCh <-chan updateCheckResult) {
	if updateCh == nil {
		return
	}

	select {
	case result := <-updateCh:
		if result.notice == nil {
			return
		}

		if _, err := fmt.Fprintln(os.Stderr, result.notice.Message()); err != nil {
			return
		}

		if err := result.notifier.MarkNotified(result.notice); err != nil {
			upgrade.LogUpdateCheckError(err)
		}
	default:
	}
}

func shouldCheckForUpdate(command *cobra.Command, interactive bool, appSettings settings.Settings) bool {
	if !interactive || !appSettings.UpgradeNotify || os.Getenv("CI") != "" {
		return false
	}

	if command == nil {
		return true
	}

	switch command.Name() {
	case "completion", "help":
		return false
	}

	return !isSelfCommand(command)
}

func isInteractiveStderr() bool {
	fd := os.Stderr.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

type flags struct {
	config  string
	debug   int
	help    bool
	version bool
	all     bool
	init    bool
}

// We can not parse --config and --debug flags using cobra.Command.ParseFlags
//
//	until we read config and initialize all subcommands.
//	Otherwise root command will parse all flags gready.
//
// For example in 'lets --config lets.my.yaml mysubcommand --config=myconfig'
//
//	cobra will parse all --config flags, but take only latest
//
// --config=myconfig, and this is wrong.
func parseRootFlags(args []string) (*flags, error) {
	f := &flags{}
	// if first arg is not a flag, then it is subcommand
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return f, nil
	}

	visited := set.NewSet[string]()

	isFlagVisited := func(name string) bool {
		if visited.Contains(name) {
			return true
		}

		visited.Add(name)

		return false
	}

	idx := 0
	for idx < len(args) {
		arg := args[idx]
		if !strings.HasPrefix(arg, "-") {
			// stop if arg is not a flag, it is probably a subcommand
			break
		}

		name, value, found := strings.Cut(arg, "=")
		switch name {
		case "--config", "-c":
			if !isFlagVisited("config") {
				if found {
					if value == "" {
						return nil, errors.New("--config must be set to value")
					}

					f.config = value
				} else if len(args[idx:]) > 0 {
					f.config = args[idx+1]
					idx += 2

					continue
				}
			}
		case "--debug", "-d", "-dd":
			if !isFlagVisited("debug") {
				f.debug = 1
				if arg == "-dd" {
					f.debug = 2
				}
			}
		case "--help", "-h":
			if !isFlagVisited("help") {
				f.help = true
			}
		case "--version", "-v":
			if !isFlagVisited("version") {
				f.version = true
			}
		case "--all":
			if !isFlagVisited("all") {
				f.all = true
			}
		case "--init":
			if !isFlagVisited("init") {
				f.init = true
			}
		case "--upgrade":
			return nil, errors.New("--upgrade has been replaced with 'lets self upgrade'")
		}

		idx += 1 //nolint:revive,golint
	}

	return f, nil
}
