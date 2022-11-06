package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/lets-cli/lets/checksum"
	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/docopt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const GenericCmdNameTpl = "LETS_COMMAND_NAME"

type LogFn = func(string, ...interface{})

func cmdLogger(cmd *config.Command, parents ...*config.Command) LogFn {
	return func(format string, a ...interface{}) {
		cmdName := cmd.Name
		for _, p := range parents {
			if p != nil {
				cmdName = fmt.Sprintf("%s => %s", p.Name, cmdName)
			}
		}

		cmdName = color.GreenString("[%s]", cmdName)
		lets := color.BlueString("lets:")
		msg := color.BlueString(fmt.Sprintf(format, a...))
		msg = fmt.Sprintf("%s %s %s", lets, cmdName, msg)
		log.Debugf(msg)
	}
}

type ExecuteError struct {
	err error
}

func (e *ExecuteError) Error() string {
	return e.err.Error()
}

// ExitCode will return exit code from underlying ExitError or returns default error code.
func (e *ExecuteError) ExitCode() int {
	var exitErr *exec.ExitError
	if ok := errors.As(e.err, &exitErr); ok {
		return exitErr.ExitCode()
	}

	return 1 // default error code
}

type Executor struct {
	cmd       *config.Command
	parentCmd *config.Command // child command if parentCmd is not nil
	cfg       *config.Config
	out       io.Writer
	// debug logger with predefined cmd name
	cmdLog LogFn
}

func NewExecutor(cmd *config.Command, cfg *config.Config, out io.Writer) *Executor {
	return &Executor{
		cmd:    cmd,
		cfg:    cfg,
		out:    out,
		cmdLog: cmdLogger(cmd),
	}
}

func NewChildExecutor(cmd *config.Command, parentExecutor *Executor) *Executor {
	return &Executor{
		cmd:       cmd,
		parentCmd: parentExecutor.cmd,
		cfg:       parentExecutor.cfg,
		out:       parentExecutor.out,
		cmdLog:    cmdLogger(cmd, parentExecutor.cmd, parentExecutor.parentCmd),
	}
}

// Execute executes command and it depends recursively
// Command can be executed in parallel.
func (e *Executor) Execute(ctx context.Context) error {
	if e.cmd.Cmds.Parallel {
		return e.executeParallel(ctx)
	}

	return e.execute(ctx)
}

// Execute main command and wait for result.
// Must be used only when Command.Cmd is string or []string.
func (e *Executor) execute(ctx context.Context) error {
	e.cmdLog("command:\n  %s", e.cmd.Pretty())

	defer func() {
		if e.cmd.After != "" {
			e.executeAfterScript()
		}
	}()

	if err := e.initCmd(); err != nil {
		return err
	}

	if err := e.executeDepends(ctx); err != nil {
		return err
	}

	if cmd := e.cmd.Cmds.SingleCommand(); cmd != nil {
		if err := e.executeCmdScript(cmd.Script); err != nil {
			return err
		}
	}

	// persist checksum only if exit code 0
	return e.persistChecksum()
}

// Executes 'after' script after main 'cmd' script
// It allowed to fail and will print error
// Do not return error directly to root because we consider only 'cmd' exit code.
// Even if 'after' script failed we return exit code from 'cmd'.
// This behavior may change in the future if needed.
func (e *Executor) executeAfterScript() {
	cmd, err := e.newOsCommand(e.cmd.After)
	if err != nil {
		log.Printf("failed to run `after` script for command '%s': %s", e.cmd.Name, err)
		return
	}

	e.cmdLog("executing 'after':\n  cmd: %s\n  env: %s", e.cmd.After, fmtEnv(cmd.Env))

	if ExecuteError := cmd.Run(); ExecuteError != nil {
		log.Printf("failed to run `after` script for command '%s': %s", e.cmd.Name, ExecuteError)
	}
}

// format docopts error and adds usage string to output.
func formatOptsUsageError(err error, opts docopt.Opts, cmdName string, rawOptions string) error {
	if opts == nil && err.Error() == "" {
		// TODO how to get wrong option name
		err = fmt.Errorf("no such option")
	}

	errTpl := fmt.Sprintf("failed to parse docopt options for cmd %s: %s", cmdName, err)

	return fmt.Errorf("%s\n\n%s", errTpl, rawOptions)
}

// Init Command before execution:
// - parse docopt
// - calculate checksum.
func (e *Executor) initCmd() error {
	if !e.cmd.SkipDocopts {
		e.cmdLog("parse docopt")
		opts, err := docopt.Parse(e.cmd.Args, e.cmd.Docopts)
		if err != nil {
			// TODO if accept_args, just continue with what we got
			//  but this may  require changes in go-docopt
			return formatOptsUsageError(err, opts, e.cmd.Name, e.cmd.Docopts)
		}

		e.cmdLog("docopt parsed: %#v", opts)

		e.cmd.Options = docopt.OptsToLetsOpt(opts)
		e.cmd.CliOptions = docopt.OptsToLetsCli(opts)
	}

	// calculate checksum if needed
	if err := e.cmd.ChecksumCalculator(e.cfg.WorkDir); err != nil {
		return fmt.Errorf("failed to calculate checksum for command '%s': %w", e.cmd.Name, err)
	}

	// if command declared as persist_checksum we must read current persisted checksums into memory
	if e.cmd.PersistChecksum {
		if checksum.IsChecksumForCmdPersisted(e.cfg.ChecksumsDir, e.cmd.Name) {
			err := e.cmd.ReadChecksumsFromDisk(e.cfg.ChecksumsDir, e.cmd.Name, e.cmd.ChecksumMap)
			if err != nil {
				return fmt.Errorf("failed to read persisted checksum for command '%s': %w", e.cmd.Name, err)
			}
		}
	}

	return nil
}

func joinBeforeAndScript(before string, script string) string {
	if before == "" {
		return script
	}

	before = strings.TrimSpace(before)

	return strings.Join([]string{before, script}, "\n")
}

// Setup env for cmd.
func (e *Executor) setupEnv(cmd *exec.Cmd, shell string) error {
	defaultEnv := map[string]string{
		GenericCmdNameTpl:   e.cmd.Name,
		"LETS_COMMAND_ARGS": strings.Join(e.cmd.CommandArgs(), " "),
		"SHELL":             shell,
	}

	checksumEnvMap := getChecksumEnvMap(e.cmd.ChecksumMap)

	var changedChecksumEnvMap map[string]string
	if e.cmd.PersistChecksum {
		changedChecksumEnvMap = getChangedChecksumEnvMap(
			e.cmd.ChecksumMap,
			e.cmd.GetPersistedChecksums(),
		)
	}

	cmdEnv, err := e.cmd.GetEnv(*e.cfg)
	if err != nil {
		return err
	}

	envMaps := []map[string]string{
		defaultEnv,
		e.cfg.Env.Dump(),
		cmdEnv,
		e.cmd.Options,
		e.cmd.CliOptions,
		checksumEnvMap,
		changedChecksumEnvMap,
	}

	envList := os.Environ()
	for _, envMap := range envMaps {
		envList = append(envList, convertEnvMapToList(envMap)...)
	}

	cmd.Env = envList

	return nil
}

// Prepare cmd to be executed:
// - set in/out
// - set dir
// - prepare environment
//
// NOTE: We intentionally do not passing ctx to exec.Command because we want to wait for process end.
// Passing ctx will change behavior of program drastically - it will kill process if context will be canceled.
func (e *Executor) newOsCommand(cmdScript string) (*exec.Cmd, error) {
	script := joinBeforeAndScript(e.cfg.Before, cmdScript)
	shell := e.cfg.Shell
	if e.cmd.Shell != "" {
		shell = e.cmd.Shell
	}

	args := []string{"-c", script}
	if len(e.cmd.CommandArgs()) > 0 {
		// for "--" see https://linux.die.net/man/1/bash
		args = append(args, "--", strings.Join(e.cmd.CommandArgs(), " "))
	}

	cmd := exec.Command(
		shell,
		args...,
	)

	// setup std out and err
	cmd.Stdout = e.out
	cmd.Stderr = e.out
	cmd.Stdin = os.Stdin

	// set working directory for command
	cmd.Dir = e.cfg.WorkDir
	if e.cmd.WorkDir != "" {
		cmd.Dir = e.cmd.WorkDir
	}

	if err := e.setupEnv(cmd, shell); err != nil {
		return nil, err
	}

	return cmd, nil
}

// Run all commands from Depends in sequential order.
func (e *Executor) executeDepends(ctx context.Context) error {
	return e.cmd.Depends.Range(func(depName string, dep config.Dep) error {
		e.cmdLog("running dependency '%s'", depName)
		dependCmd := e.cfg.Commands[depName]
		if dependCmd.Cmds.Parallel {
			// TODO: this must be ensured at the validation time, not at the runtime
			// forbid to run parallel command in depends
			return &ExecuteError{
				err: fmt.Errorf(
					"failed to run child command '%s' from 'depends': cmd as map is not allowed in depends yet",
					e.cmd.Name,
				),
			}
		}

		// by default, if depends command in simple format, skip docopts
		dependCmd.SkipDocopts = true
		if len(dep.Args) != 0 {
			dependCmd = dependCmd.WithArgs(dep.Args)
			dependCmd.SkipDocopts = false
		}

		if !dep.Env.Empty() {
			dependCmd = dependCmd.WithEnv(dep.Env)
		}

		// TODO: move working with ref to parsing
		if dependCmd.Ref != nil {
			dependCmd = e.cfg.Commands[dependCmd.Ref.Name].FromRef(dependCmd.Ref)
		}

		err := NewChildExecutor(dependCmd, e).Execute(ctx)
		if err != nil {
			// must return error to root
			return err
		}

		return nil
	})
}

// Persist new calculated checksum to disk.
// This function mus be called only after command finished(exited) with status 0.
func (e *Executor) persistChecksum() error {
	if e.cmd.PersistChecksum {
		e.cmdLog("persisting checksum")

		if err := e.cfg.CreateChecksumsDir(); err != nil {
			return err
		}

		err := checksum.PersistCommandsChecksumToDisk(
			e.cfg.ChecksumsDir,
			e.cmd.ChecksumMap,
			e.cmd.Name,
		)
		if err != nil {
			return fmt.Errorf("can not persist checksum to disk: %w", err)
		}
	}

	return nil
}

func (e *Executor) executeCmdScript(script string) error {
	cmd, err := e.newOsCommand(script)
	if err != nil {
		return err
	}

	e.cmdLog("executing:\n  cmd: %s\n  env: %s\n", script, fmtEnv(cmd.Env))

	if err := cmd.Run(); err != nil {
		return &ExecuteError{err: fmt.Errorf("failed to run command '%s': %w", e.cmd.Name, err)}
	}

	return nil
}

// Execute all commands from Cmds in parallel and wait for results.
func (e *Executor) executeParallel(ctx context.Context) (err error) {
	defer func() {
		if e.cmd.After != "" {
			e.executeAfterScript()
		}
	}()

	if err = e.initCmd(); err != nil {
		return err
	}

	if err = e.executeDepends(ctx); err != nil {
		return err
	}

	group, _ := errgroup.WithContext(ctx)

	for _, cmd := range e.cmd.Cmds.Commands {
		cmd := cmd
		// wait for cmd to end in a goroutine with error propagation
		group.Go(func() error {
			return e.executeCmdScript(cmd.Script)
		})
	}

	if err = group.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	// persist checksum only if exit code 0
	if err = e.persistChecksum(); err != nil {
		return fmt.Errorf("persist checksum error in command '%s': %w", e.cmd.Name, err)
	}

	return err
}
