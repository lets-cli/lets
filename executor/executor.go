package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lets-cli/lets/checksum"
	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/docopt"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/logging"
	"golang.org/x/sync/errgroup"
)

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
	cfg *config.Config
	out io.Writer
}

func NewExecutor(cfg *config.Config, out io.Writer) *Executor {
	return &Executor{
		cfg: cfg,
		out: out,
	}
}

type Context struct {
	ctx     context.Context
	command *config.Command
	logger  *logging.ExecLogger
}

func NewExecutorCtx(ctx context.Context, command *config.Command) *Context {
	return &Context{
		ctx:     ctx,
		command: command,
		logger:  logging.NewExecLogger().Child(command.Name),
	}
}

func ChildExecutorCtx(ctx *Context, command *config.Command) *Context {
	return &Context{
		command: command,
		logger:  ctx.logger.Child(command.Name),
	}
}

// Execute executes command and it depends recursively
// Command can be executed in parallel.
func (e *Executor) Execute(ctx *Context) error {
	if e.cfg.Init != "" {
		if err := e.runCmd(ctx, &config.Cmd{Script: e.cfg.Init}); err != nil {
			return err
		}
	}

	if ctx.command.Cmds.Parallel {
		return e.executeParallel(ctx)
	}

	return e.execute(ctx)
}

// Execute main command and wait for result.
// Must be used only when Command.Cmd is string or []string.
func (e *Executor) execute(ctx *Context) error {
	command := ctx.command

	if env.DebugLevel() > 1 {
		ctx.logger.Debug("command %s", command.Dump())
	}

	defer func() {
		if command.After != "" {
			e.executeAfterScript(ctx)
		}
	}()

	if err := e.initCmd(ctx); err != nil {
		return err
	}

	if err := e.executeDepends(ctx); err != nil {
		return err
	}

	for _, cmd := range command.Cmds.Commands {
		if err := e.runCmd(ctx, cmd); err != nil {
			return err
		}
	}

	// persist checksum only if exit code 0
	return e.persistChecksum(ctx)
}

// Executes 'after' script after main 'cmd' script
// It allowed to fail and will print error
// Do not return error directly to root because we consider only 'cmd' exit code.
// Even if 'after' script failed we return exit code from 'cmd'.
// This behavior may change in the future if needed.
func (e *Executor) executeAfterScript(ctx *Context) {
	command := ctx.command

	osCmd, err := e.newOsCommand(command, command.After)
	if err != nil {
		ctx.logger.Info("failed to run `after` script: %s", err)
		return
	}

	ctx.logger.Debug("executing 'after':\n  cmd: %s\n  env: %s", command.After, fmtEnv(osCmd.Env))

	if ExecuteError := osCmd.Run(); ExecuteError != nil {
		ctx.logger.Info("failed to run `after` script: %s", ExecuteError)
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
func (e *Executor) initCmd(ctx *Context) error {
	cmd := ctx.command

	if !cmd.SkipDocopts {
		ctx.logger.Debug("parse docopt: docopt: %s, args: %s", cmd.Docopts, cmd.Args)
		opts, err := docopt.Parse(cmd.Name, cmd.Args, cmd.Docopts)
		if err != nil {
			// TODO if accept_args, just continue with what we got
			//  but this may  require changes in go-docopt
			return formatOptsUsageError(err, opts, cmd.Name, cmd.Docopts)
		}

		ctx.logger.Debug("docopt parsed: %v", opts)

		cmd.Options = docopt.OptsToLetsOpt(opts)
		cmd.CliOptions = docopt.OptsToLetsCli(opts)
	}

	// calculate checksum if needed
	if err := cmd.ChecksumCalculator(e.cfg.WorkDir); err != nil {
		return fmt.Errorf("failed to calculate checksum for command '%s': %w", cmd.Name, err)
	}

	// if command declared as persist_checksum we must read current persisted checksums into memory
	if cmd.PersistChecksum {
		if checksum.IsChecksumForCmdPersisted(e.cfg.ChecksumsDir, cmd.Name) {
			err := cmd.ReadChecksumsFromDisk(e.cfg.ChecksumsDir, cmd.Name, cmd.ChecksumMap)
			if err != nil {
				return fmt.Errorf("failed to read persisted checksum for command '%s': %w", cmd.Name, err)
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
func (e *Executor) setupEnv(osCmd *exec.Cmd, command *config.Command, shell string) error {
	defaultEnv := map[string]string{
		"LETS_COMMAND_NAME":     command.Name,
		"LETS_COMMAND_ARGS":     strings.Join(command.Args, " "),
		"LETS_COMMAND_WORK_DIR": osCmd.Dir,
		"LETS_CONFIG":           filepath.Base(e.cfg.FilePath),
		"LETS_CONFIG_DIR":       filepath.Dir(e.cfg.FilePath),
		"SHELL":                 shell,
	}

	checksumEnvMap := getChecksumEnvMap(command.ChecksumMap)

	var changedChecksumEnvMap map[string]string
	if command.PersistChecksum {
		changedChecksumEnvMap = getChangedChecksumEnvMap(
			command.ChecksumMap,
			command.GetPersistedChecksums(),
		)
	}

	cmdEnv, err := command.GetEnv(*e.cfg)
	if err != nil {
		return err
	}

	envMaps := []map[string]string{
		defaultEnv,
		e.cfg.GetEnv(),
		cmdEnv,
		command.Options,
		command.CliOptions,
		checksumEnvMap,
		changedChecksumEnvMap,
	}

	envList := os.Environ()
	for _, envMap := range envMaps {
		envList = append(envList, convertEnvMapToList(envMap)...)
	}

	osCmd.Env = envList

	return nil
}

// Prepare cmd to be executed:
// - set in/out
// - set dir
// - prepare environment
//
// NOTE: We intentionally do not passing ctx to exec.Command because we want to wait for process end.
// Passing ctx will change behavior of program drastically - it will kill process if context will be canceled.
func (e *Executor) newOsCommand(command *config.Command, cmdScript string) (*exec.Cmd, error) {
	script := joinBeforeAndScript(e.cfg.Before, cmdScript)
	shell := e.cfg.Shell
	if command.Shell != "" {
		shell = command.Shell
	}

	args := []string{"-c", script}
	if len(command.Args) > 0 {
		// for "--" see https://linux.die.net/man/1/bash
		args = append(args, "--", strings.Join(command.Args, " "))
	}

	osCmd := exec.Command(
		shell,
		args...,
	)

	// setup std out and err
	osCmd.Stdout = e.out
	osCmd.Stderr = e.out
	osCmd.Stdin = os.Stdin

	// set working directory for command
	osCmd.Dir = e.cfg.WorkDir
	if command.WorkDir != "" {
		osCmd.Dir = command.WorkDir
	}

	if err := e.setupEnv(osCmd, command, shell); err != nil {
		return nil, err
	}

	return osCmd, nil
}

// Run all commands from Depends in sequential order.
func (e *Executor) executeDepends(ctx *Context) error {
	return ctx.command.Depends.Range(func(depName string, dep config.Dep) error {
		ctx.logger.Debug("running dependency '%s'", depName)
		cmd := e.cfg.Commands[depName]

		cmd = cmd.Clone()

		if dep.HasArgs() {
			cmd.Args = dep.Args
			ctx.logger.Debug("dependency args overridden: '%s'", cmd.Args)
		}

		if !dep.Env.Empty() {
			cmd.Env.Merge(dep.Env)
			ctx.logger.Debug("dependency env overridden: '%s'", cmd.Env.Dump())
		}

		return e.Execute(ChildExecutorCtx(ctx, cmd))
	})
}

// Persist new calculated checksum to disk.
// This function mus be called only after command finished(exited) with status 0.
func (e *Executor) persistChecksum(ctx *Context) error {
	cmd := ctx.command

	if cmd.PersistChecksum {
		ctx.logger.Debug("persisting checksum")

		if err := e.cfg.CreateChecksumsDir(); err != nil {
			return err
		}

		err := checksum.PersistCommandsChecksumToDisk(
			e.cfg.ChecksumsDir,
			cmd.ChecksumMap,
			cmd.Name,
		)
		if err != nil {
			return fmt.Errorf("can not persist checksum to disk: %w", err)
		}
	}

	return nil
}

func (e *Executor) runCmd(ctx *Context, cmd *config.Cmd) error {
	command := ctx.command
	osCmd, err := e.newOsCommand(command, cmd.Script)
	if err != nil {
		return err
	}

	if env.DebugLevel() == 1 {
		ctx.logger.Debug("executing script:\n%s", cmd.Script)
	} else if env.DebugLevel() > 1 {
		ctx.logger.Debug("executing:\nscript: %s\nenv: %s\n", cmd.Script, fmtEnv(osCmd.Env))
	}

	if err := osCmd.Run(); err != nil {
		return &ExecuteError{err: fmt.Errorf("failed to run command '%s': %w", command.Name, err)}
	}

	return nil
}

// Execute all commands from Cmds in parallel and wait for results.
func (e *Executor) executeParallel(ctx *Context) (err error) {
	command := ctx.command

	defer func() {
		if command.After != "" {
			e.executeAfterScript(ctx)
		}
	}()

	if err = e.initCmd(ctx); err != nil {
		return err
	}

	if err = e.executeDepends(ctx); err != nil {
		return err
	}

	group, _ := errgroup.WithContext(ctx.ctx)

	for _, cmd := range command.Cmds.Commands {
		cmd := cmd
		// wait for cmd to end in a goroutine with error propagation
		group.Go(func() error {
			return e.runCmd(ctx, cmd)
		})
	}

	if err = group.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	// persist checksum only if exit code 0
	if err = e.persistChecksum(ctx); err != nil {
		return fmt.Errorf("persist checksum error in command '%s': %w", command.Name, err)
	}

	return err
}
