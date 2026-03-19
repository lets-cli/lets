package cli

import (
	"testing"

	cmdpkg "github.com/lets-cli/lets/internal/cmd"
	"github.com/spf13/cobra"
)

func TestAllowsMissingConfig(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		command := &cobra.Command{Use: "help"}
		if !allowsMissingConfig(command) {
			t.Fatal("expected help to allow missing config")
		}
	})

	t.Run("completion", func(t *testing.T) {
		root := cmdpkg.CreateRootCommand("v0.0.0-test", "")
		cmdpkg.InitCompletionCmd(root, nil)

		command, _, err := root.Find([]string{"completion"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !allowsMissingConfig(command) {
			t.Fatal("expected completion to allow missing config")
		}
	})

	t.Run("self subcommand", func(t *testing.T) {
		root := cmdpkg.CreateRootCommand("v0.0.0-test", "")
		cmdpkg.InitSelfCmd(root, "v0.0.0-test")

		command, _, err := root.Find([]string{"self", "doc"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !allowsMissingConfig(command) {
			t.Fatal("expected lets self doc to allow missing config")
		}
	})

	t.Run("top level doc does not match self", func(t *testing.T) {
		root := cmdpkg.CreateRootCommand("v0.0.0-test", "")
		root.AddCommand(&cobra.Command{Use: "doc"})

		command, _, err := root.Find([]string{"doc"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if allowsMissingConfig(command) {
			t.Fatal("expected top-level doc to require config")
		}
	})
}
