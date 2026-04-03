package cli

import (
	"testing"

	cmdpkg "github.com/lets-cli/lets/internal/cmd"
	"github.com/lets-cli/lets/internal/settings"
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

func TestShouldCheckForUpdate(t *testing.T) {
	defaultSettings := settings.Default()

	t.Run("should allow normal interactive commands", func(t *testing.T) {
		t.Setenv("CI", "")

		if !shouldCheckForUpdate(&cobra.Command{Use: "run"}, true, defaultSettings) {
			t.Fatal("expected update check to be enabled")
		}
	})

	t.Run("should skip non interactive sessions", func(t *testing.T) {
		if shouldCheckForUpdate(&cobra.Command{Use: "run"}, false, defaultSettings) {
			t.Fatal("expected non-interactive session to skip update check")
		}
	})

	t.Run("should skip when CI is set", func(t *testing.T) {
		t.Setenv("CI", "1")
		if shouldCheckForUpdate(&cobra.Command{Use: "run"}, true, defaultSettings) {
			t.Fatal("expected CI to skip update check")
		}
	})

	t.Run("should skip when notifier disabled in settings", func(t *testing.T) {
		disabled := settings.Default()
		disabled.UpgradeNotify = false

		if shouldCheckForUpdate(&cobra.Command{Use: "run"}, true, disabled) {
			t.Fatal("expected disabled settings to skip update check")
		}
	})

	t.Run("should skip completion and help commands", func(t *testing.T) {
		for _, command := range []*cobra.Command{{Use: "completion"}, {Use: "help"}} {
			if shouldCheckForUpdate(command, true, defaultSettings) {
				t.Fatalf("expected %q to skip update check", command.Name())
			}
		}
	})

	t.Run("should skip self subcommands", func(t *testing.T) {
		root := cmdpkg.CreateRootCommand("v0.0.0-test", "")
		cmdpkg.InitSelfCmd(root, "v0.0.0-test")

		for _, args := range [][]string{{"self"}, {"self", "doc"}, {"self", "upgrade"}} {
			command, _, err := root.Find(args)
			if err != nil {
				t.Fatalf("unexpected error for %v: %v", args, err)
			}

			if shouldCheckForUpdate(command, true, defaultSettings) {
				t.Fatalf("expected %v to skip update check", args)
			}
		}
	})
}

func TestParseRootFlags(t *testing.T) {
	t.Run("should reject legacy upgrade flag", func(t *testing.T) {
		_, err := parseRootFlags([]string{"--upgrade"})
		if err == nil {
			t.Fatal("expected legacy upgrade flag error")
		}

		if err.Error() != "--upgrade has been replaced with 'lets self upgrade'" {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
