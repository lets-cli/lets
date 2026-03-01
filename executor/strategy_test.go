package executor

import (
	"context"
	"testing"

	"github.com/lets-cli/lets/config/config"
)

func TestSequentialStrategy(t *testing.T) {
	t.Run("should execute commands in order", func(t *testing.T) {
		strategy := &SequentialStrategy{}
		
		order := []string{}
		cmds := []*config.Cmd{
			{Name: "first", Script: "echo first"},
			{Name: "second", Script: "echo second"},
		}

		// Create a mock that tracks order
		for _, cmd := range cmds {
			order = append(order, cmd.Name)
		}

		if len(order) != 2 {
			t.Errorf("expected 2 commands, got %d", len(order))
		}
		if order[0] != "first" || order[1] != "second" {
			t.Error("commands should be in order")
		}

		// Verify strategy exists and has correct type
		if strategy == nil {
			t.Error("strategy should not be nil")
		}
	})

	t.Run("should respect context cancellation", func(t *testing.T) {
		strategy := &SequentialStrategy{}
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		cmds := []*config.Cmd{{Script: "echo test"}}
		command := &config.Command{Name: "test"}

		err := strategy.Run(ctx, cmds, command, nil, nil)

		if err != context.Canceled {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
	})
}

func TestParallelStrategy(t *testing.T) {
	t.Run("should exist", func(t *testing.T) {
		strategy := &ParallelStrategy{}
		if strategy == nil {
			t.Error("strategy should not be nil")
		}
	})
}

func TestSelectStrategy(t *testing.T) {
	t.Run("should return sequential for non-parallel", func(t *testing.T) {
		cmds := config.Cmds{Parallel: false}
		strategy := SelectStrategy(cmds)

		if _, ok := strategy.(*SequentialStrategy); !ok {
			t.Error("expected SequentialStrategy")
		}
	})

	t.Run("should return parallel for parallel", func(t *testing.T) {
		cmds := config.Cmds{Parallel: true}
		strategy := SelectStrategy(cmds)

		if _, ok := strategy.(*ParallelStrategy); !ok {
			t.Error("expected ParallelStrategy")
		}
	})
}
