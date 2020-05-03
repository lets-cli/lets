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

	return runCmdWait(ctx, &cmdToRun, cfg, out, noParent)
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
func prepareCmdForRun(
	ctx context.Context,
	cmdToRun *command.Command,
	cmdScript string,
	cfg *config.Config,
	out io.Writer,
	parentName string,
) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, cfg.Shell, "-c", cmdScript) // #nosec G204
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
func runDepends(ctx context.Context, cmdToRun *command.Command, cfg *config.Config, out io.Writer) error {
	for _, dependCmdName := range cmdToRun.Depends {
		dependCmd := cfg.Commands[dependCmdName]

		err := runCmdWait(ctx, &dependCmd, cfg, out, cmdToRun.Name)
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

// Run command and wait for result
func runCmdWait(
	ctx context.Context,
	cmdToRun *command.Command,
	cfg *config.Config,
	out io.Writer,
	parentName string,
) error {
	cmd, err := prepareCmdForRun(ctx, cmdToRun, cmdToRun.Cmd, cfg, out, parentName)
	if err != nil {
		return err
	}

	if err := runDepends(ctx, cmdToRun, cfg, out); err != nil {
		return err
	}

	runErr := cmd.Run()
	if runErr != nil {
		return fmt.Errorf("failed to run cmd: %s", runErr)
	}

	// persist checksum only if exit code 0
	if err := persistChecksum(*cmdToRun); err != nil {
		return err
	}

	return nil
}

// Start cmd and return without waiting for result.
// Cmd must be waited by caller.
func runCmdNoWait(
	ctx context.Context,
	cmdToRun *command.Command,
	cmdScript string,
	cfg *config.Config,
	out io.Writer,
) (*exec.Cmd, error) {
	cmd, err := prepareCmdForRun(ctx, cmdToRun, cmdScript, cfg, out, noParent)
	if err != nil {
		return nil, err
	}

	if err := runDepends(ctx, cmdToRun, cfg, out); err != nil {
		return nil, err
	}

	startErr := cmd.Start()
	if startErr != nil {
		return nil, fmt.Errorf("failed to start cmd: %s", startErr)
	}

	return cmd, nil
}

func runCmdAsMap(ctx context.Context, cmdToRun *command.Command, cfg *config.Config, out io.Writer) error {
	g, ctx := errgroup.WithContext(ctx)

	// TODO how do we use cmdName ???
	for _, cmdExecScript := range cmdToRun.CmdMap {
		cmdStarted, err := runCmdNoWait(ctx, cmdToRun, cmdExecScript, cfg, out)

		if err != nil {
			return err
		}

		// wait for cmd to end in a goroutine with error propagation
		g.Go(cmdStarted.Wait)
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
