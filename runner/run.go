package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/lets-cli/lets/checksum"
	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/config/parser"
	"github.com/lets-cli/lets/env"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const color = "\033[1;36m%s\033[0m"

const GenericCmdNameTpl = "LETS_COMMAND_NAME"

func colored(msg string, color string) string {
	if env.IsNotColorOutput() {
		return msg
	}

	return fmt.Sprintf(color, msg)
}

func debugf(format string, a ...interface{}) {
	prefixed := fmt.Sprintf("lets: %s", fmt.Sprintf(format, a...))
	log.Debugf(colored(prefixed, color))
}

type RunErr struct {
	err error
}

func (e *RunErr) Error() string {
	return e.err.Error()
}

// ExitCode will return exit code from underlying ExitError or returns default error code.
func (e *RunErr) ExitCode() int {
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
}

func NewRunner(cmd *config.Command, cfg *config.Config, out io.Writer) *Runner {
	return &Runner{
		cmd: cmd,
		cfg: cfg,
		out: out,
	}
}

func NewChildRunner(cmd *config.Command, parentRunner *Runner) *Runner {
	return &Runner{
		cmd:       cmd,
		parentCmd: parentRunner.cmd,
		cfg:       parentRunner.cfg,
		out:       parentRunner.out,
	}
}

// Execute runs command.
func (r *Runner) Execute(ctx context.Context) error {
	if r.parentCmd != nil {
		return r.runChild(ctx)
	}

	if r.cmd.CmdMap != nil {
		return r.runCmdAsMap(ctx)
	}

	return r.run(ctx)
}

// Run main command and wait for result.
// Must be used only when Command.Cmd is string or []string.
func (r *Runner) run(ctx context.Context) error {
	debugf("running command '%s': %s", r.cmd.Name, r.cmd.Pretty())

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

	if err := r.runCmdScript(r.cmd.Cmd); err != nil {
		return err
	}

	// persist checksum only if exit code 0
	return r.persistChecksum()
}

// Run command and wait for result.
// Must be used only when Command.Cmd is string or []string.
func (r *Runner) runChild(ctx context.Context) error {
	debugf("running child command '%s': %s", r.cmd.Name, r.cmd.Pretty())

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

	cmd := r.prepareOsCommandForRun(r.cmd.Cmd)

	debugf(
		"executing child os command for %s -> %s\ncmd: %s\nenv: %s\n",
		r.parentCmd.Name, r.cmd.Name, r.cmd.Cmd, fmtEnv(cmd.Env),
	)

	if err := cmd.Run(); err != nil {
		return &RunErr{err: fmt.Errorf("failed to run child command '%s' from 'depends': %w", r.cmd.Name, err)}
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
	cmd := r.prepareOsCommandForRun(r.cmd.After)

	debugf("executing after script:\ncommand: %s\nscript: %s\nenv: %s", r.cmd.Name, r.cmd.After, fmtEnv(cmd.Env))

	if runErr := cmd.Run(); runErr != nil {
		log.Printf("failed to run `after` script for command '%s': %s", r.cmd.Name, runErr)
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
		debugf("parse docopt for command '%s'", r.cmd.Name)
		opts, err := parser.ParseDocopts(r.cmd.Args, r.cmd.Docopts)
		if err != nil {
			// TODO if accept_args, just continue with what we got
			//  but this may  require changes in go-docopt
			return formatOptsUsageError(err, opts, r.cmd.Name, r.cmd.Docopts)
		}

		debugf("raw docopt for command '%s': %#v", r.cmd.Name, opts)

		r.cmd.Options = parser.OptsToLetsOpt(opts)
		r.cmd.CliOptions = parser.OptsToLetsCli(opts)
	}

	// calculate checksum if needed
	if err := r.cmd.ChecksumCalculator(r.cfg.WorkDir); err != nil {
		return fmt.Errorf("failed to calculate checksum for command '%s': %w", r.cmd.Name, err)
	}

	// if command declared as persist_checksum we must read current persisted checksums into memory
	if r.cmd.PersistChecksum {
		if checksum.IsChecksumForCmdPersisted(r.cfg.DotLetsDir, r.cmd.Name) {
			err := r.cmd.ReadChecksumsFromDisk(r.cfg.DotLetsDir, r.cmd.Name, r.cmd.ChecksumMap)
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

	return strings.Join([]string{before, script}, "\n")
}

// Setup env for cmd.
func (r *Runner) setupEnv(cmd *exec.Cmd, shell string) {
	defaultEnv := map[string]string{
		GenericCmdNameTpl:   r.cmd.Name,
		"LETS_COMMAND_ARGS": strings.Join(r.cmd.CommandArgs, " "),
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

	envMaps := []map[string]string{
		defaultEnv,
		r.cfg.Env,
		r.cmd.Env,
		r.cmd.OverrideEnv,
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
}

// Prepare cmd to be run:
// - set in/out
// - set dir
// - prepare environment
//
// NOTE: We intentionally do not passing ctx to exec.Command because we want to wait for process end.
// Passing ctx will change behavior of program drastically - it will kill process if context will be canceled.
//
func (r *Runner) prepareOsCommandForRun(cmdScript string) *exec.Cmd {
	script := joinBeforeAndScript(r.cfg.Before, cmdScript)
	shell := r.cfg.Shell
	if r.cmd.Shell != "" {
		shell = r.cmd.Shell
	}

	cmd := exec.Command(
		shell,
		"-c",
		script,
		"--", // see https://linux.die.net/man/1/bash
		strings.Join(r.cmd.CommandArgs, " "),
	) // #nosec G204

	// setup std out and err
	cmd.Stdout = r.out
	cmd.Stderr = r.out
	cmd.Stdin = os.Stdin

	// set working directory for command
	cmd.Dir = r.cfg.WorkDir
	if r.cmd.WorkDir != "" {
		cmd.Dir = r.cmd.WorkDir
	}

	r.setupEnv(cmd, shell)

	return cmd
}

// Run all commands from Depends in sequential order.
func (r *Runner) runDepends(ctx context.Context) error {
	for _, depName := range r.cmd.DependsNames {
		dep := r.cmd.Depends[depName]
		debugf("running dependency '%s' for command '%s'", depName, r.cmd.Name)

		dependCmd := r.cfg.Commands[depName]
		if dependCmd.CmdMap != nil {
			// forbid to run depends command as map
			return &RunErr{
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
		if len(dep.Env) != 0 {
			dependCmd = dependCmd.WithEnv(dep.Env)
		}

		if dependCmd.Ref != "" {
			dependCmd = r.cfg.Commands[dependCmd.Ref].FromRef(dependCmd)
		}

		err := NewChildRunner(&dependCmd, r).Execute(ctx)
		if err != nil {
			// must return error to root
			return err
		}
	}

	return nil
}

// Persist new calculated checksum to disk.
// This function mus be called only after command finished(exited) with status 0.
func (r *Runner) persistChecksum() error {
	if r.cmd.PersistChecksum {
		debugf("persisting checksum for command '%s'", r.cmd.Name)

		err := checksum.PersistCommandsChecksumToDisk(
			r.cfg.DotLetsDir,
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
		buf = fmt.Sprintf("%s\n%s", buf, entry)
	}

	return buf
}

func (r *Runner) runCmdScript(cmdScript string) error {
	cmd := r.prepareOsCommandForRun(cmdScript)

	debugf("executing os command for '%s'\ncmd: %s\nenv: %s\n", r.cmd.Name, r.cmd.Cmd, fmtEnv(cmd.Env))

	if err := cmd.Run(); err != nil {
		return &RunErr{err: fmt.Errorf("failed to run command '%s': %w", r.cmd.Name, err)}
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
// Must be used only when Command.Cmd is map[string]string.
func (r *Runner) runCmdAsMap(ctx context.Context) (err error) {
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

	g, _ := errgroup.WithContext(ctx)

	cmdMap, err := filterCmdMap(r.cmd.Name, r.cmd.CmdMap, r.cmd.Only, r.cmd.Exclude)
	if err != nil {
		return err
	}

	for _, cmdExecScript := range cmdMap {
		cmdExecScript := cmdExecScript
		// wait for cmd to end in a goroutine with error propagation
		g.Go(func() error {
			return r.runCmdScript(cmdExecScript)
		})
	}

	if err = g.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	// persist checksum only if exit code 0
	if err = r.persistChecksum(); err != nil {
		return fmt.Errorf("persist checksum error in command '%s': %w", r.cmd.Name, err)
	}

	return err
}
