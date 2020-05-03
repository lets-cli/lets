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

// Prepare cmd to be run:
// - set in/out
// - set dir
// - parse docopt
// - calculate checksum
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
) (*exec.Cmd, error) {
	cmd := exec.Command(cfg.Shell, "-c", cmdScript) // #nosec G204
	// setup std out and err
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Stdin = os.Stdin

	// set working directory for command
	cmd.Dir = cfg.WorkDir

	isChildCmd := parentName != ""

	// parse docopts - only for parent
	if !isChildCmd {
		opts, err := command.ParseDocopts(cmdToRun.RawOptions)
		if err != nil {
			return nil, formatOptsUsageError(err, opts, cmdToRun.Name, cmdToRun.RawOptions)
		}

		cmdToRun.Options = command.OptsToLetsOpt(opts)
		cmdToRun.CliOptions = command.OptsToLetsCli(opts)
	}

	// calculate checksum if needed
	if err := cmdToRun.ChecksumCalculator(); err != nil {
		return nil, err
	}

	// if command declared as persist_checksum we must read current persisted checksums into memory
	var persistedChecksums map[string]string

	if cmdToRun.PersistChecksum {
		if command.ChecksumForCmdPersisted(cmdToRun.Name) {
			checksums, err := command.ReadChecksumsFromDisk(cmdToRun.Name, cmdToRun.ChecksumMap)
			if err != nil {
				return nil, err
			}

			persistedChecksums = checksums
		}
	}

	// setup env for command
	cmd.Env = composeEnvs(
		os.Environ(),
		convertEnvMapToList(cfg.Env),
		convertEnvMapToList(cmdToRun.Env),
		convertEnvMapToList(cmdToRun.Options),
		convertEnvMapToList(cmdToRun.CliOptions),
		convertChecksumToEnvForCmd(cmdToRun.Checksum),
		convertChecksumMapToEnvForCmd(cmdToRun.ChecksumMap),
	)

	if cmdToRun.PersistChecksum {
		cmd.Env = composeEnvs(
			cmd.Env,
			convertChangedChecksumMapToEnvForCmd(cmdToRun.Checksum, cmdToRun.ChecksumMap, persistedChecksums),
		)
	}

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

	return cmd, nil
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
func persistChecksum(cmdToRun command.Command) error {
	if cmdToRun.PersistChecksum {
		err := command.PersistCommandsChecksumToDisk(cmdToRun)
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
	if err := runDepends(cmdToRun, cfg, out); err != nil {
		return err
	}

	if err := runCmdScript(cmdToRun, cmdToRun.Cmd, cfg, out, parentName); err != nil {
		return err
	}

	// persist checksum only if exit code 0
	if err := persistChecksum(*cmdToRun); err != nil {
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

	cmd, err := prepareCmdForRun(cmdToRun, cmdScript, cfg, out, parentName)
	if err != nil {
		return err
	}

	runErr := cmd.Run()
	if runErr != nil {
		if isChildCmd {
			return fmt.Errorf("failed to run child command '%s' from 'depends': %s", cmdToRun.Name, runErr)
		}

		return fmt.Errorf("failed to run command '%s': %s", cmdToRun.Name, runErr)
	}

	return nil
}

// Run all commands from Command.CmdMap in parallel and wait for results.
// Must be used only when Command.Cmd is map[string]string
func runCmdAsMap(ctx context.Context, cmdToRun *command.Command, cfg *config.Config, out io.Writer) error {
	if err := runDepends(cmdToRun, cfg, out); err != nil {
		return err
	}

	g, _ := errgroup.WithContext(ctx)

	for _, cmdExecScript := range cmdToRun.CmdMap {
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
	if err := persistChecksum(*cmdToRun); err != nil {
		return err
	}

	return nil
}
