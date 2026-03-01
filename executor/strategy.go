package executor

import (
	"context"

	"github.com/lets-cli/lets/config/config"
	"golang.org/x/sync/errgroup"
)

// ExecutionStrategy defines how multiple commands are executed.
type ExecutionStrategy interface {
	// Run executes commands using the provided runner.
	Run(ctx context.Context, commands []*config.Cmd, command *config.Command, cmdEnv map[string]string, runner *CommandRunner) error
}

// SequentialStrategy executes commands one after another.
type SequentialStrategy struct{}

// Run executes commands sequentially, stopping on first error.
func (s *SequentialStrategy) Run(
	ctx context.Context,
	commands []*config.Cmd,
	command *config.Command,
	cmdEnv map[string]string,
	runner *CommandRunner,
) error {
	for _, cmd := range commands {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := runner.RunScript(command, cmd.Script, cmdEnv); err != nil {
				return err
			}
		}
	}
	return nil
}

// ParallelStrategy executes commands concurrently.
type ParallelStrategy struct{}

// Run executes all commands in parallel, returning the first error.
func (p *ParallelStrategy) Run(
	ctx context.Context,
	commands []*config.Cmd,
	command *config.Command,
	cmdEnv map[string]string,
	runner *CommandRunner,
) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, cmd := range commands {
		cmd := cmd // capture for goroutine
		g.Go(func() error {
			return runner.RunScript(command, cmd.Script, cmdEnv)
		})
	}

	return g.Wait()
}

// SelectStrategy returns the appropriate execution strategy based on command configuration.
func SelectStrategy(cmds config.Cmds) ExecutionStrategy {
	if cmds.Parallel {
		return &ParallelStrategy{}
	}
	return &SequentialStrategy{}
}
