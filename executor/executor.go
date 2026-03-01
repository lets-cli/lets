package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/docopt"
	"github.com/lets-cli/lets/logging"
)

// ExecuteError wraps execution errors and provides exit code extraction.
type ExecuteError struct {
	err error
}

func (e *ExecuteError) Error() string {
	return e.err.Error()
}

// ExitCode returns the exit code from the underlying ExitError, or 1 as default.
func (e *ExecuteError) ExitCode() int {
	var exitErr *exec.ExitError
	if ok := errors.As(e.err, &exitErr); ok {
		return exitErr.ExitCode()
	}
	return 1
}

func (e *ExecuteError) Unwrap() error {
	return e.err
}

// Context holds the execution context for a command.
type Context struct {
	ctx     context.Context
	command *config.Command
	logger  *logging.ExecLogger
}

// NewExecutorCtx creates a new execution context.
func NewExecutorCtx(ctx context.Context, command *config.Command) *Context {
	return &Context{
		ctx:     ctx,
		command: command,
		logger:  logging.NewExecLogger().Child(command.Name),
	}
}

// ChildExecutorCtx creates a child context for dependency execution.
func ChildExecutorCtx(ctx *Context, command *config.Command) *Context {
	return &Context{
		ctx:     ctx.ctx,
		command: command,
		logger:  ctx.logger.Child(command.Name),
	}
}

// Executor is the main command executor.
// It delegates to PipelineExecutor for the actual execution.
type Executor struct {
	pipeline *PipelineExecutor
}

// NewExecutor creates a new executor.
func NewExecutor(cfg *config.Config, out io.Writer) *Executor {
	return &Executor{
		pipeline: NewPipelineExecutor(cfg, out),
	}
}

// Execute runs the command through the execution pipeline.
func (e *Executor) Execute(ctx *Context) error {
	return e.pipeline.Execute(ctx)
}

// formatOptsUsageError formats a docopt parsing error with usage information.
func formatOptsUsageError(err error, opts docopt.Opts, cmdName string, rawOptions string) error {
	if opts == nil && err.Error() == "" {
		err = errors.New("no such option")
	}

	errTpl := fmt.Sprintf("failed to parse docopt options for cmd %s: %s", cmdName, err)
	return fmt.Errorf("%s\n\n%s", errTpl, rawOptions)
}
