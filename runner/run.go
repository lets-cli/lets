package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/docopt/docopt-go"
	"golang.org/x/sync/errgroup"

	"github.com/lets-cli/lets/commands/command"
	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/logging"
)

const (
	NoticeColor = "\033[1;36m%s\033[0m"
)

const noParent = ""

type RunOptions struct {
	Config  *config.Config
	RawArgs []string
}

// RunCommand runs parent command
// TODO maybe we should store commands map in config as map[string]*Command (as pointers)
func RunCommand(ctx context.Context, cmdToRun command.Command, cfg *config.Config, out io.Writer) error {
	if cmdToRun.CmdMap != nil {
		return runCmdAsMap(ctx, &cmdToRun, cfg, out)
	}

	return runCmd(&cmdToRun, cfg, out, noParent)
}

// format docopts error and adds usage string to output
func formatOptsUsageError(err error, opts docopt.Opts, cmdName string, rawOptions string) error {
	if opts == nil && err.Error() == "" {
		err = fmt.Errorf("no such option")
	}

	errTpl := fmt.Sprintf("failed to parse docopt options for cmd %s: %s", cmdName, err)

	return fmt.Errorf("%s\n\n%s", errTpl, rawOptions)
}

// Init Command before run:
// - parse docopt
// - calculate checksum
func initCmd(
	cmdToRun *command.Command,
	cfg *config.Config,
	isChildCmd bool,
) error {
	// parse docopts - only for parent
	if !isChildCmd {
		opts, err := command.ParseDocopts(cmdToRun.Args, cmdToRun.RawOptions)
		if err != nil {
			return formatOptsUsageError(err, opts, cmdToRun.Name, cmdToRun.RawOptions)
		}

		cmdToRun.Options = command.OptsToLetsOpt(opts)
		cmdToRun.CliOptions = command.OptsToLetsCli(opts)
	}

	// calculate checksum if needed
	if err := cmdToRun.ChecksumCalculator(cfg.WorkDir); err != nil {
		return err
	}

	// if command declared as persist_checksum we must read current persisted checksums into memory
	if cmdToRun.PersistChecksum {
		if command.ChecksumForCmdPersisted(cfg.DotLetsDir, cmdToRun.Name) {
			err := cmdToRun.ReadChecksumsFromDisk(cfg.DotLetsDir, cmdToRun.Name, cmdToRun.ChecksumMap)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Prepare cmd to be run:
// - set in/out
// - set dir
// - prepare environment
//
// NOTE: We intentionally do not passing ctx to exec.Command because we want to wait for process end.
// Passing ctx will change behavior of program drastically - it will kill process if context will be canceled.
//
func prepareCmdForRun(
	cmdToRun *command.Command,
	cmdScript string,
	cfg *config.Config,
	out io.Writer,
	parentName string,
) *exec.Cmd {
	cmd := exec.Command(cfg.Shell, "-c", cmdScript) // #nosec G204
	// setup std out and err
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Stdin = os.Stdin

	// set working directory for command
	cmd.Dir = cfg.WorkDir

	// setup env for command
	cmd.Env = composeEnvs(
		os.Environ(),
		convertEnvMapToList(cfg.Env),
		convertEnvMapToList(cmdToRun.Env),
		convertEnvMapToList(cmdToRun.OverrideEnv),
		convertEnvMapToList(cmdToRun.Options),
		convertEnvMapToList(cmdToRun.CliOptions),
		convertChecksumToEnvForCmd(cmdToRun.Checksum),
		convertChecksumMapToEnvForCmd(cmdToRun.ChecksumMap),
	)

	if cmdToRun.PersistChecksum {
		cmd.Env = composeEnvs(
			cmd.Env,
			convertChangedChecksumMapToEnvForCmd(
				cmdToRun.Checksum,
				cmdToRun.ChecksumMap,
				cmdToRun.GetPersistedChecksums(),
			),
		)
	}

	isChildCmd := parentName != ""

	if !isChildCmd {
		logging.Log.Debugf(
			"Executing command\nname: %s\ncmd: %s\nenv:\n%s",
			fmt.Sprintf(NoticeColor, cmdToRun.Name),
			fmt.Sprintf(NoticeColor, cmdToRun.Cmd),
			cmd.Env,
		)
	} else {
		logging.Log.Debugf(
			"Executing child command\nparent name: %s\nname: %s\ncmd: %s\nenv:\n%s",
			fmt.Sprintf(NoticeColor, parentName),
			fmt.Sprintf(NoticeColor, cmdToRun.Name),
			fmt.Sprintf(NoticeColor, cmdToRun.Cmd),
			cmd.Env,
		)
	}

	return cmd
}

// Run all commands from Depends in sequential order
func runDepends(cmdToRun *command.Command, cfg *config.Config, out io.Writer) error {
	for _, dependCmdName := range cmdToRun.Depends {
		dependCmd := cfg.Commands[dependCmdName]

		err := runCmd(&dependCmd, cfg, out, cmdToRun.Name)
		if err != nil {
			// must return error to root
			return err
		}
	}

	return nil
}

// Persist new calculated checksum to disk.
// This function mus be called only after command finished(exited) with status 0
func persistChecksum(cmdToRun command.Command, cfg *config.Config) error {
	if cmdToRun.PersistChecksum {
		err := command.PersistCommandsChecksumToDisk(cfg.DotLetsDir, cmdToRun)
		if err != nil {
			return err
		}
	}

	return nil
}

// Run command and wait for result.
// Must be used only when Command.Cmd is string or []string
func runCmd(
	cmdToRun *command.Command,
	cfg *config.Config,
	out io.Writer,
	parentName string,
) error {
	if err := initCmd(cmdToRun, cfg, parentName != noParent); err != nil {
		return err
	}

	if err := runDepends(cmdToRun, cfg, out); err != nil {
		return err
	}

	if err := runCmdScript(cmdToRun, cmdToRun.Cmd, cfg, out, parentName); err != nil {
		return err
	}

	// persist checksum only if exit code 0
	if err := persistChecksum(*cmdToRun, cfg); err != nil {
		return err
	}

	return nil
}

func runCmdScript(
	cmdToRun *command.Command,
	cmdScript string,
	cfg *config.Config,
	out io.Writer,
	parentName string,
) error {
	isChildCmd := parentName != ""

	cmd := prepareCmdForRun(cmdToRun, cmdScript, cfg, out, parentName)

	runErr := cmd.Run()
	if runErr != nil {
		if isChildCmd {
			return fmt.Errorf("failed to run child command '%s' from 'depends': %s", cmdToRun.Name, runErr)
		}

		return fmt.Errorf("failed to run command '%s': %s", cmdToRun.Name, runErr)
	}

	return nil
}

func filterCmdMap(
	parentCmdName string,
	cmdMap map[string]string,
	only []string,
	exclude []string,
) (map[string]string, error) {
	hasOnly := len(only) > 0
	hasExclude := len(exclude) > 0

	if !hasOnly && !hasExclude {
		return cmdMap, nil
	}

	filteredCmdMap := make(map[string]string)

	if hasOnly {
		// put only commands which in `only` list
		for _, cmdName := range only {
			cmdScript, ok := cmdMap[cmdName]
			if !ok {
				return nil, fmt.Errorf("no such sub-command '%s' in command '%s' used in 'only' flag", cmdName, parentCmdName)
			}

			filteredCmdMap[cmdName] = cmdScript
		}
	}

	if hasExclude {
		filteredCmdMap = cmdMap
		// delete all commands which in `exclude` list
		for _, cmdName := range exclude {
			_, ok := cmdMap[cmdName]
			if !ok {
				return nil, fmt.Errorf("no such sub-command '%s' in command '%s' used in 'exclude' flag", cmdName, parentCmdName)
			}

			delete(filteredCmdMap, cmdName)
		}
	}

	return filteredCmdMap, nil
}

// Run all commands from Command.CmdMap in parallel and wait for results.
// Must be used only when Command.Cmd is map[string]string
func runCmdAsMap(ctx context.Context, cmdToRun *command.Command, cfg *config.Config, out io.Writer) error {
	if err := initCmd(cmdToRun, cfg, false); err != nil {
		return err
	}

	if err := runDepends(cmdToRun, cfg, out); err != nil {
		return err
	}

	g, _ := errgroup.WithContext(ctx)

	cmdMap, err := filterCmdMap(cmdToRun.Name, cmdToRun.CmdMap, cmdToRun.Only, cmdToRun.Exclude)
	if err != nil {
		return err
	}

	for _, cmdExecScript := range cmdMap {
		cmdExecScript := cmdExecScript
		// wait for cmd to end in a goroutine with error propagation
		g.Go(func() error {
			return runCmdScript(cmdToRun, cmdExecScript, cfg, out, noParent)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	// persist checksum only if exit code 0
	if err := persistChecksum(*cmdToRun, cfg); err != nil {
		return err
	}

	return nil
}
