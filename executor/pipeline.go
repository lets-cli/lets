package executor

import (
	"context"
	"io"
	"sync"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/env"
)

// PipelineExecutor orchestrates command execution through a series of phases.
// It replaces the monolithic Executor with a composable, phase-based approach.
type PipelineExecutor struct {
	cfg *config.Config
	out io.Writer

	// Phases to execute for each command
	phases []Phase

	// Init tracking
	initOnce   sync.Once
	initScript string
}

// NewPipelineExecutor creates a new pipeline-based executor.
func NewPipelineExecutor(cfg *config.Config, out io.Writer) *PipelineExecutor {
	return &PipelineExecutor{
		cfg:        cfg,
		out:        out,
		initScript: cfg.Init,
		phases: []Phase{
			&DocoptPhase{},
			&ChecksumPhase{},
			&EnvPhase{},
			&DependencyPhase{},
			&ExecutionPhase{},
			&ChecksumPersistPhase{},
		},
	}
}

// Execute runs a command through the pipeline.
func (e *PipelineExecutor) Execute(ctx *Context) error {
	// Run global init script once
	if err := e.runInitOnce(ctx); err != nil {
		return err
	}

	plan := e.createPlan(ctx)
	return e.executePlanWithAfter(ctx.ctx, plan)
}

// runInitOnce executes the global init script exactly once.
func (e *PipelineExecutor) runInitOnce(ctx *Context) error {
	if e.initScript == "" {
		return nil
	}

	var initErr error
	e.initOnce.Do(func() {
		runner := NewCommandRunner(e.cfg, e.out, ctx.logger)
		// Create a minimal command for init script execution
		initCmd := &config.Command{Name: "__init__"}
		initErr = runner.RunScript(initCmd, e.initScript, nil)
	})
	return initErr
}

// createPlan creates an execution plan for a command.
func (e *PipelineExecutor) createPlan(ctx *Context) *ExecutionPlan {
	return &ExecutionPlan{
		Config:      e.cfg,
		Command:     ctx.command,
		Logger:      ctx.logger,
		Runner:      NewCommandRunner(e.cfg, e.out, ctx.logger),
		ExecutorRef: e,
	}
}

// createChildPlan creates a plan for a dependency command.
func (e *PipelineExecutor) createChildPlan(parent *ExecutionPlan, cmd *config.Command) *ExecutionPlan {
	return &ExecutionPlan{
		Config:      e.cfg,
		Command:     cmd,
		Logger:      parent.Logger.Child(cmd.Name),
		Runner:      NewCommandRunner(e.cfg, e.out, parent.Logger.Child(cmd.Name)),
		ParentPlan:  parent,
		ExecutorRef: e,
	}
}

// executePlanWithAfter runs the pipeline and ensures after script runs.
func (e *PipelineExecutor) executePlanWithAfter(ctx context.Context, plan *ExecutionPlan) error {
	// Defer after script execution
	defer func() {
		afterPhase := &AfterScriptPhase{}
		_ = afterPhase.Execute(ctx, plan)
	}()

	return e.executePlan(ctx, plan)
}

// executePlan runs a plan through all phases.
func (e *PipelineExecutor) executePlan(ctx context.Context, plan *ExecutionPlan) error {
	if env.DebugLevel() > 1 {
		plan.Logger.Debug("command %s", plan.Command.Dump())
	}

	for _, phase := range e.phases {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := phase.Execute(ctx, plan); err != nil {
				return err
			}
		}
	}

	return nil
}

// Pipeline allows customizing the execution phases.
type Pipeline struct {
	phases []Phase
}

// NewPipeline creates a new customizable pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{
		phases: make([]Phase, 0),
	}
}

// Add appends a phase to the pipeline.
func (p *Pipeline) Add(phase Phase) *Pipeline {
	p.phases = append(p.phases, phase)
	return p
}

// Phases returns the configured phases.
func (p *Pipeline) Phases() []Phase {
	return p.phases
}

// DefaultPipeline returns the standard execution pipeline.
func DefaultPipeline() *Pipeline {
	return NewPipeline().
		Add(&DocoptPhase{}).
		Add(&ChecksumPhase{}).
		Add(&EnvPhase{}).
		Add(&DependencyPhase{}).
		Add(&ExecutionPhase{}).
		Add(&ChecksumPersistPhase{})
}

// WithCustomPhases creates an executor with custom phases.
func NewPipelineExecutorWithPhases(cfg *config.Config, out io.Writer, phases []Phase) *PipelineExecutor {
	return &PipelineExecutor{
		cfg:        cfg,
		out:        out,
		initScript: cfg.Init,
		phases:     phases,
	}
}

// Ensure PipelineExecutor implements the same interface pattern as old Executor
var _ interface {
	Execute(*Context) error
} = (*PipelineExecutor)(nil)
