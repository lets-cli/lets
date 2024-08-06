package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lets-cli/lets/cmd"
	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/executor"
	"github.com/lets-cli/lets/logging"
	"github.com/lets-cli/lets/set"
	"github.com/lets-cli/lets/upgrade"
	"github.com/lets-cli/lets/upgrade/registry"
	"github.com/lets-cli/lets/workdir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version = "0.0.0-dev"

func main() {
	ctx := getContext()

	configDir := os.Getenv("LETS_CONFIG_DIR")

	logging.InitLogging(os.Stdout, os.Stderr)

	rootCmd := cmd.CreateRootCommand(version)
	rootCmd.InitDefaultHelpFlag()
	rootCmd.InitDefaultVersionFlag()
	reinitCompletionCmd := cmd.InitCompletionCmd(rootCmd, nil)
	rootCmd.InitDefaultHelpCmd()

	command, args, err := rootCmd.Traverse(os.Args[1:])
	if err != nil {
		log.Errorf("lets: traverse commands error: %s", err)
		os.Exit(1)
	}

	rootFlags, err := parseRootFlags(args)
	if err != nil {
		log.Errorf("lets: parse flags error: %s", err)
		os.Exit(1)
	}

	if rootFlags.version {
		if err := cmd.PrintVersionMessage(rootCmd); err != nil {
			log.Errorf("lets: print version error: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
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
			log.Errorf("lets: config error: %s", err)
			os.Exit(1)
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
			log.Errorf("lets: can not create lets.yaml: %s", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	if rootFlags.upgrade {
		upgrader, err := upgrade.NewBinaryUpgrader(registry.NewGithubRegistry(ctx), version)
		if err == nil {
			err = upgrader.Upgrade()
		}

		if err != nil {
			log.Errorf("lets: can not self-upgrade binary: %s", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	showUsage := rootFlags.help || (command.Name() == "help" && len(args) == 0)

	if showUsage {
		if err := cmd.PrintHelpMessage(rootCmd); err != nil {
			log.Errorf("lets: print help error: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Error(err.Error())

		exitCode := 1
		if e, ok := err.(*executor.ExecuteError); ok { //nolint:errorlint
			exitCode = e.ExitCode()
		}

		os.Exit(exitCode)
	}
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
		log.Printf("lets: signal received: %s", sig)
		cancel()
	}()

	return ctx
}

// do not fail on config error in it is help (-h, --help) or --init or completion command
func failOnConfigError(root *cobra.Command, current *cobra.Command, rootFlags *flags) bool {
	rootCommands := set.NewSet("completion", "help")
	return (root.Flags().NFlag() == 0 && !rootCommands.Contains(current.Name())) && !rootFlags.help && !rootFlags.init
}

type flags struct {
	config  string
	debug   int
	help    bool
	version bool
	all     bool
	init    bool
	upgrade bool
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
			if !isFlagVisited("upgrade") {
				f.upgrade = true
			}
		}

		idx += 1 //nolint:revive,golint
	}

	return f, nil
}
