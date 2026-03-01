package executor

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	t.Run("should create pipeline with phases", func(t *testing.T) {
		pipeline := NewPipeline().
			Add(&DocoptPhase{}).
			Add(&ChecksumPhase{}).
			Add(&EnvPhase{})

		phases := pipeline.Phases()

		if len(phases) != 3 {
			t.Errorf("expected 3 phases, got %d", len(phases))
		}
	})

	t.Run("default pipeline should have all phases", func(t *testing.T) {
		pipeline := DefaultPipeline()
		phases := pipeline.Phases()

		expectedPhases := []string{
			"docopt",
			"checksum",
			"env",
			"dependencies",
			"execution",
			"checksum_persist",
		}

		if len(phases) != len(expectedPhases) {
			t.Errorf("expected %d phases, got %d", len(expectedPhases), len(phases))
		}

		for i, phase := range phases {
			if phase.Name() != expectedPhases[i] {
				t.Errorf("phase %d: expected '%s', got '%s'", i, expectedPhases[i], phase.Name())
			}
		}
	})
}

func TestPhaseNames(t *testing.T) {
	tests := []struct {
		phase Phase
		name  string
	}{
		{&DocoptPhase{}, "docopt"},
		{&ChecksumPhase{}, "checksum"},
		{&EnvPhase{}, "env"},
		{&DependencyPhase{}, "dependencies"},
		{&ExecutionPhase{}, "execution"},
		{&ChecksumPersistPhase{}, "checksum_persist"},
		{&AfterScriptPhase{}, "after_script"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.phase.Name() != tc.name {
				t.Errorf("expected '%s', got '%s'", tc.name, tc.phase.Name())
			}
		})
	}
}
