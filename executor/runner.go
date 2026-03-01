package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/logging"
)

// CommandRunner is responsible for executing shell commands.
// It handles the low-level details of creating and running OS commands.
type CommandRunner struct {
	cfg    *config.Config
	out    io.Writer
	logger *logging.ExecLogger
}

// NewCommandRunner creates a new command runner.
func NewCommandRunner(cfg *config.Config, out io.Writer, logger *logging.ExecLogger) *CommandRunner {
	return &CommandRunner{
		cfg:    cfg,
		out:    out,
		logger: logger,
	}
}

// RunScript executes a single script within the context of a command.
func (r *CommandRunner) RunScript(command *config.Command, script string, cmdEnv map[string]string) error {
	osCmd, err := r.createOsCommand(command, script, cmdEnv)
	if err != nil {
		return err
	}

	r.logExecution(script, osCmd.Env)

	if err := osCmd.Run(); err != nil {
		return &ExecuteError{err: fmt.Errorf("failed to run command '%s': %w", command.Name, err)}
	}

	return nil
}

// RunScriptIgnoreError executes a script and logs errors instead of returning them.
// Used for 'after' scripts where we don't want to fail the main command.
func (r *CommandRunner) RunScriptIgnoreError(command *config.Command, script string, cmdEnv map[string]string) {
	osCmd, err := r.createOsCommand(command, script, cmdEnv)
	if err != nil {
		r.logger.Info("failed to create 'after' script command: %s", err)
		return
	}

	r.logger.Debug("executing 'after':\n  cmd: %s\n  env: %s", script, FormatEnvForDebug(osCmd.Env))

	if err := osCmd.Run(); err != nil {
		r.logger.Info("failed to run 'after' script: %s", err)
	}
}

// createOsCommand creates an exec.Cmd configured for execution.
func (r *CommandRunner) createOsCommand(
	command *config.Command,
	script string,
	cmdEnv map[string]string,
) (*exec.Cmd, error) {
	// Join before script if present
	fullScript := joinBeforeAndScript(r.cfg.Before, script)

	// Determine shell
	shell := r.cfg.Shell
	if command.Shell != "" {
		shell = command.Shell
	}

	// Build command arguments
	args := []string{"-c", fullScript}
	if len(command.Args) > 0 {
		args = append(args, "--", strings.Join(command.Args, " "))
	}

	osCmd := exec.Command(shell, args...)

	// Setup I/O
	osCmd.Stdout = r.out
	osCmd.Stderr = r.out
	osCmd.Stdin = os.Stdin

	// Set working directory
	osCmd.Dir = r.cfg.WorkDir
	if command.WorkDir != "" {
		osCmd.Dir = command.WorkDir
	}

	// Build environment
	envBuilder := BuildCommandEnv(command, r.cfg, cmdEnv)
	osCmd.Env = envBuilder.Build(GetBaseEnv())

	return osCmd, nil
}

func (r *CommandRunner) logExecution(script string, cmdEnv []string) {
	if env.DebugLevel() == 1 {
		r.logger.Debug("executing script:\n%s", script)
	} else if env.DebugLevel() > 1 {
		r.logger.Debug("executing:\nscript: %s\nenv: %s\n", script, FormatEnvForDebug(cmdEnv))
	}
}

// joinBeforeAndScript concatenates the before script with the main script.
func joinBeforeAndScript(before string, script string) string {
	if before == "" {
		return script
	}
	before = strings.TrimSpace(before)
	return before + "\n" + script
}
