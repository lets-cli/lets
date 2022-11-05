package runner

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

type RunError struct {
	err error
}

func (e *RunError) Error() string {
	return e.err.Error()
}

// ExitCode will return exit code from underlying ExitError or returns default error code.
func (e *RunError) ExitCode() int {
	var exitErr *exec.ExitError
	if ok := errors.As(e.err, &exitErr); ok {
		return exitErr.ExitCode()
	}

	return 1 // default error code
}

type Runner struct {
	cmd       *config.Command
	parentCmd *config.Command // child command if parentCmd is not nil
	cfg       *config.Config
	out       io.Writer
	// debug logger with predefined cmd name
	cmdLog LogFn
}

func NewRunner(cmd *config.Command, cfg *config.Config, out io.Writer) *Runner {
	return &Runner{
		cmd:    cmd,
		cfg:    cfg,
		out:    out,
		cmdLog: cmdLogger(cmd),
	}
}

func NewChildRunner(cmd *config.Command, parentRunner *Runner) *Runner {
	return &Runner{
		cmd:       cmd,
		parentCmd: parentRunner.cmd,
		cfg:       parentRunner.cfg,
		out:       parentRunner.out,
		cmdLog:    cmdLogger(cmd, parentRunner.cmd, parentRunner.parentCmd),
	}
}

// Execute runs command.
func (r *Runner) Execute(ctx context.Context) error {
	if r.parentCmd != nil {
		return r.runChild(ctx)
	}

	if r.cmd.Cmds.Parallel {
		return r.runParallel(ctx)
	}

	return r.run(ctx)
}

// Run main command and wait for result.
// Must be used only when Command.Cmd is string or []string.
func (r *Runner) run(ctx context.Context) error {
	r.cmdLog("command:\n  %s", r.cmd.Pretty())

	defer func() {
		if r.cmd.After != "" {
			r.runAfterScript()
		}
	}()

	if err := r.initCmd(); err != nil {
		return err
	}

	if err := r.runDepends(ctx); err != nil {
		return err
	}

	if len(r.cmd.Cmds.Commands) > 0 {
		if err := r.runCmdScript(r.cmd.Cmds.Commands[0].Script); err != nil {
			return err
		}
	}

	// persist checksum only if exit code 0
	return r.persistChecksum()
}

// Run command and wait for result.
// Must be used only when Command.Cmd is string or []string.
func (r *Runner) runChild(ctx context.Context) error {
	r.cmdLog("command:\n  %s", r.cmd.Pretty())

	defer func() {
		if r.cmd.After != "" {
			r.runAfterScript()
		}
	}()

	// never skip docopt for main command
	if err := r.initCmd(); err != nil {
		return err
	}

	if err := r.runDepends(ctx); err != nil {
		return err
	}

	script := r.cmd.Cmds.Commands[0].Script
	cmd, err := r.newOsCommand(script)
	if err != nil {
		return err
	}

	r.cmdLog("executing:\n  cmd: %s\n  env: %s\n", script, fmtEnv(cmd.Env))

	if err := cmd.Run(); err != nil {
		return &RunError{err: fmt.Errorf("failed to run child command '%s' from 'depends': %w", r.cmd.Name, err)}
	}

	// persist checksum only if exit code 0
	return r.persistChecksum()
}

// Runs 'after' script after main 'cmd' script
// It allowed to fail and will print error
// Do not return error directly to root because we consider only 'cmd' exit code.
// Even if 'after' script failed we return exit code from 'cmd'.
// This behavior may change in the future if needed.
func (r *Runner) runAfterScript() {
	cmd, err := r.newOsCommand(r.cmd.After)
	if err != nil {
		log.Printf("failed to run `after` script for command '%s': %s", r.cmd.Name, err)
		return
	}

	r.cmdLog("executing 'after':\n  cmd: %s\n  env: %s", r.cmd.After, fmtEnv(cmd.Env))

	if RunError := cmd.Run(); RunError != nil {
		log.Printf("failed to run `after` script for command '%s': %s", r.cmd.Name, RunError)
	}
}

type RunOptions struct {
	Config  *config.Config
	RawArgs []string
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

// Init Command before run:
// - parse docopt
// - calculate checksum.
func (r *Runner) initCmd() error {
	if !r.cmd.SkipDocopts {
		r.cmdLog("parse docopt")
		opts, err := docopt.Parse(r.cmd.Args, r.cmd.Docopts)
		if err != nil {
			// TODO if accept_args, just continue with what we got
			//  but this may  require changes in go-docopt
			return formatOptsUsageError(err, opts, r.cmd.Name, r.cmd.Docopts)
		}

		r.cmdLog("docopt parsed: %#v", opts)

		r.cmd.Options = docopt.OptsToLetsOpt(opts)
		r.cmd.CliOptions = docopt.OptsToLetsCli(opts)
	}

	// calculate checksum if needed
	if err := r.cmd.ChecksumCalculator(r.cfg.WorkDir); err != nil {
		return fmt.Errorf("failed to calculate checksum for command '%s': %w", r.cmd.Name, err)
	}

	// if command declared as persist_checksum we must read current persisted checksums into memory
	if r.cmd.PersistChecksum {
		if checksum.IsChecksumForCmdPersisted(r.cfg.ChecksumsDir, r.cmd.Name) {
			err := r.cmd.ReadChecksumsFromDisk(r.cfg.ChecksumsDir, r.cmd.Name, r.cmd.ChecksumMap)
			if err != nil {
				return fmt.Errorf("failed to read persisted checksum for command '%s': %w", r.cmd.Name, err)
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
func (r *Runner) setupEnv(cmd *exec.Cmd, shell string) error {
	defaultEnv := map[string]string{
		GenericCmdNameTpl:   r.cmd.Name,
		"LETS_COMMAND_ARGS": strings.Join(r.cmd.CommandArgs(), " "),
		"SHELL":             shell,
	}

	checksumEnvMap := getChecksumEnvMap(r.cmd.ChecksumMap)

	var changedChecksumEnvMap map[string]string
	if r.cmd.PersistChecksum {
		changedChecksumEnvMap = getChangedChecksumEnvMap(
			r.cmd.ChecksumMap,
			r.cmd.GetPersistedChecksums(),
		)
	}

	cmdEnv, err := r.cmd.GetEnv(*r.cfg)
	if err != nil {
		return err
	}

	envMaps := []map[string]string{
		defaultEnv,
		r.cfg.Env.Dump(),
		cmdEnv,
		r.cmd.Options,
		r.cmd.CliOptions,
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

// Prepare cmd to be run:
// - set in/out
// - set dir
// - prepare environment
//
// NOTE: We intentionally do not passing ctx to exec.Command because we want to wait for process end.
// Passing ctx will change behavior of program drastically - it will kill process if context will be canceled.
func (r *Runner) newOsCommand(cmdScript string) (*exec.Cmd, error) {
	script := joinBeforeAndScript(r.cfg.Before, cmdScript)
	shell := r.cfg.Shell
	if r.cmd.Shell != "" {
		shell = r.cmd.Shell
	}

	args := []string{"-c", script}
	if len(r.cmd.CommandArgs()) > 0 {
		// for "--" see https://linux.die.net/man/1/bash
		args = append(args, "--", strings.Join(r.cmd.CommandArgs(), " "))
	}

	cmd := exec.Command(
		shell,
		args...,
	)

	// setup std out and err
	cmd.Stdout = r.out
	cmd.Stderr = r.out
	cmd.Stdin = os.Stdin

	// set working directory for command
	cmd.Dir = r.cfg.WorkDir
	if r.cmd.WorkDir != "" {
		cmd.Dir = r.cmd.WorkDir
	}

	if err := r.setupEnv(cmd, shell); err != nil {
		return nil, err
	}

	return cmd, nil
}

// Run all commands from Depends in sequential order.
func (r *Runner) runDepends(ctx context.Context) error {
	return r.cmd.Depends.Range(func(depName string, dep config.Dep) error {
		r.cmdLog("running dependency '%s'", depName)
		dependCmd := r.cfg.Commands[depName]
		if dependCmd.Cmds.Parallel {
			// TODO: this must be ensured at the validation time, not at the runtime
			// forbid to run parallel command in depends
			return &RunError{
				err: fmt.Errorf(
					"failed to run child command '%s' from 'depends': cmd as map is not allowed in depends yet",
					r.cmd.Name,
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
			dependCmd = r.cfg.Commands[dependCmd.Ref.Name].FromRef(dependCmd.Ref)
		}

		err := NewChildRunner(dependCmd, r).Execute(ctx)
		if err != nil {
			// must return error to root
			return err
		}

		return nil
	})
}

// Persist new calculated checksum to disk.
// This function mus be called only after command finished(exited) with status 0.
func (r *Runner) persistChecksum() error {
	if r.cmd.PersistChecksum {
		r.cmdLog("persisting checksum")

		if err := r.cfg.CreateChecksumsDir(); err != nil {
			return err
		}

		err := checksum.PersistCommandsChecksumToDisk(
			r.cfg.ChecksumsDir,
			r.cmd.ChecksumMap,
			r.cmd.Name,
		)
		if err != nil {
			return fmt.Errorf("can not persist checksum to disk: %w", err)
		}
	}

	return nil
}

func fmtEnv(env []string) string {
	buf := ""

	for _, entry := range env {
		buf = fmt.Sprintf("%s\n  %s", buf, entry)
	}

	return buf
}

func (r *Runner) runCmdScript(script string) error {
	cmd, err := r.newOsCommand(script)
	if err != nil {
		return err
	}

	r.cmdLog("executing:\n  cmd: %s\n  env: %s\n", script, fmtEnv(cmd.Env))

	if err := cmd.Run(); err != nil {
		return &RunError{err: fmt.Errorf("failed to run command '%s': %w", r.cmd.Name, err)}
	}

	return nil
}

// Run all commands from Cmds in parallel and wait for results.
func (r *Runner) runParallel(ctx context.Context) (err error) {
	defer func() {
		if r.cmd.After != "" {
			r.runAfterScript()
		}
	}()

	if err = r.initCmd(); err != nil {
		return err
	}

	if err = r.runDepends(ctx); err != nil {
		return err
	}

	group, _ := errgroup.WithContext(ctx)

	for _, cmd := range r.cmd.Cmds.Commands {
		cmd := cmd
		// wait for cmd to end in a goroutine with error propagation
		group.Go(func() error {
			return r.runCmdScript(cmd.Script)
		})
	}

	if err = group.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	// persist checksum only if exit code 0
	if err = r.persistChecksum(); err != nil {
		return fmt.Errorf("persist checksum error in command '%s': %w", r.cmd.Name, err)
	}

	return err
}
