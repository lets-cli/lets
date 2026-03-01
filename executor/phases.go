package executor

import (
	"context"
	"fmt"

	"github.com/lets-cli/lets/checksum"
	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/docopt"
	"github.com/lets-cli/lets/logging"
)

// Phase represents a single step in the command execution pipeline.
// Phases are executed in order, and each phase can access and modify the execution plan.
type Phase interface {
	// Name returns the phase name for logging/debugging.
	Name() string
	// Execute runs the phase. Returns error to abort the pipeline.
	Execute(ctx context.Context, plan *ExecutionPlan) error
}

// ExecutionPlan holds all the data needed for command execution.
// It's built up by phases and passed down the pipeline.
type ExecutionPlan struct {
	// Config is the global lets configuration.
	Config *config.Config
	// Command is the command being executed.
	Command *config.Command
	// CommandEnv holds the computed command environment variables.
	CommandEnv map[string]string
	// Runner executes shell commands.
	Runner *CommandRunner
	// Logger for this execution.
	Logger *logging.ExecLogger
	// ParentPlan for nested (dependency) executions.
	ParentPlan *ExecutionPlan
	// ExecutorRef back-reference for dependency execution.
	ExecutorRef *PipelineExecutor
}

// DocoptPhase parses command-line arguments using docopt.
type DocoptPhase struct{}

func (p *DocoptPhase) Name() string { return "docopt" }

func (p *DocoptPhase) Execute(ctx context.Context, plan *ExecutionPlan) error {
	cmd := plan.Command

	if cmd.SkipDocopts {
		plan.Logger.Debug("skipping docopt parsing")
		return nil
	}

	plan.Logger.Debug("parse docopt: %s, args: %s", cmd.Docopts, cmd.Args)

	opts, err := docopt.Parse(cmd.Name, cmd.Args, cmd.Docopts)
	if err != nil {
		return formatOptsUsageError(err, opts, cmd.Name, cmd.Docopts)
	}

	plan.Logger.Debug("docopt parsed: %v", opts)

	cmd.Options = docopt.OptsToLetsOpt(opts)
	cmd.CliOptions = docopt.OptsToLetsCli(opts)

	return nil
}

// ChecksumPhase calculates checksums for the command.
type ChecksumPhase struct{}

func (p *ChecksumPhase) Name() string { return "checksum" }

func (p *ChecksumPhase) Execute(ctx context.Context, plan *ExecutionPlan) error {
	cmd := plan.Command
	cfg := plan.Config

	// Calculate checksum if needed
	if err := cmd.ChecksumCalculator(cfg.WorkDir); err != nil {
		return fmt.Errorf("failed to calculate checksum for command '%s': %w", cmd.Name, err)
	}

	// Read persisted checksums if persist_checksum is enabled
	if cmd.PersistChecksum {
		if checksum.IsChecksumForCmdPersisted(cfg.ChecksumsDir, cmd.Name) {
			err := cmd.ReadChecksumsFromDisk(cfg.ChecksumsDir, cmd.Name, cmd.ChecksumMap)
			if err != nil {
				return fmt.Errorf("failed to read persisted checksum for command '%s': %w", cmd.Name, err)
			}
		}
	}

	return nil
}

// EnvPhase computes the command environment variables.
type EnvPhase struct{}

func (p *EnvPhase) Name() string { return "env" }

func (p *EnvPhase) Execute(ctx context.Context, plan *ExecutionPlan) error {
	cmdEnv, err := plan.Command.GetEnv(*plan.Config)
	if err != nil {
		return fmt.Errorf("failed to get command environment: %w", err)
	}
	plan.CommandEnv = cmdEnv
	return nil
}

// DependencyPhase executes command dependencies.
type DependencyPhase struct{}

func (p *DependencyPhase) Name() string { return "dependencies" }

func (p *DependencyPhase) Execute(ctx context.Context, plan *ExecutionPlan) error {
	return plan.Command.Depends.Range(func(depName string, dep config.Dep) error {
		plan.Logger.Debug("running dependency '%s'", depName)

		cmd := plan.Config.Commands[depName]
		cmd = cmd.Clone()

		if dep.HasArgs() {
			cmd.Args = dep.Args
			plan.Logger.Debug("dependency args overridden: '%s'", cmd.Args)
		}

		if !dep.Env.Empty() {
			cmd.Env.Merge(dep.Env)
			plan.Logger.Debug("dependency env overridden: '%s'", cmd.Env.Dump())
		}

		// Create child plan and execute
		childPlan := plan.ExecutorRef.createChildPlan(plan, cmd)
		return plan.ExecutorRef.executePlan(ctx, childPlan)
	})
}

// ExecutionPhase runs the actual command scripts.
type ExecutionPhase struct{}

func (p *ExecutionPhase) Name() string { return "execution" }

func (p *ExecutionPhase) Execute(ctx context.Context, plan *ExecutionPlan) error {
	strategy := SelectStrategy(plan.Command.Cmds)
	return strategy.Run(
		ctx,
		plan.Command.Cmds.Commands,
		plan.Command,
		plan.CommandEnv,
		plan.Runner,
	)
}

// ChecksumPersistPhase persists checksums after successful execution.
type ChecksumPersistPhase struct{}

func (p *ChecksumPersistPhase) Name() string { return "checksum_persist" }

func (p *ChecksumPersistPhase) Execute(ctx context.Context, plan *ExecutionPlan) error {
	cmd := plan.Command
	cfg := plan.Config

	if !cmd.PersistChecksum {
		return nil
	}

	plan.Logger.Debug("persisting checksum")

	if err := cfg.CreateChecksumsDir(); err != nil {
		return err
	}

	err := checksum.PersistCommandsChecksumToDisk(
		cfg.ChecksumsDir,
		cmd.ChecksumMap,
		cmd.Name,
	)
	if err != nil {
		return fmt.Errorf("can not persist checksum to disk: %w", err)
	}

	return nil
}

// AfterScriptPhase runs the command's after script.
// This phase is special - it runs in a deferred manner and doesn't fail the pipeline.
type AfterScriptPhase struct{}

func (p *AfterScriptPhase) Name() string { return "after_script" }

func (p *AfterScriptPhase) Execute(ctx context.Context, plan *ExecutionPlan) error {
	cmd := plan.Command
	if cmd.After == "" {
		return nil
	}

	// After scripts don't fail the pipeline - they just log errors
	plan.Runner.RunScriptIgnoreError(cmd, cmd.After, plan.CommandEnv)
	return nil
}
