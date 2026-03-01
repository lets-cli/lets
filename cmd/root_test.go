package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/lets-cli/lets/config/config"
	"github.com/spf13/cobra"
)

func newTestRootCmd(args []string) (rootCmd *cobra.Command) {
	root := CreateRootCommand("v0.0.0-test")
	root.SetArgs(args)
	InitCompletionCmd(root, nil)

	return root
}

func newTestRootCmdWithConfig(args []string) (rootCmd *cobra.Command, out *bytes.Buffer) {
	bufOut := new(bytes.Buffer)

	cfg := &config.Config{
		Commands: make(map[string]*config.Command),
	}
	cfg.Commands["foo"] = &config.Command{Name: "foo"}
	cfg.Commands["bar"] = &config.Command{Name: "bar"}

	root := CreateRootCommand("v0.0.0-test")
	root.SetArgs(args)
	root.SetOut(bufOut)
	root.SetErr(bufOut)

	InitCompletionCmd(root, cfg)
	InitSubCommands(root, cfg, true, out)

	return root, bufOut
}

func TestRootCmd(t *testing.T) {
	t.Run("should init completion subcommand", func(t *testing.T) {
		var args []string
		rootCmd := newTestRootCmd(args)

		expectedTotal := 1 //  completion

		comp, _, _ := rootCmd.Find([]string{"completion"})
		if comp.Name() != "completion" {
			t.Errorf("no '%s' subcommand in the root command", "completion")
		}
		totalCommands := len(rootCmd.Commands())
		if totalCommands != expectedTotal {
			t.Errorf(
				"root cmd has different number of subcommands than expected. Exp: %d, Got: %d",
				expectedTotal,
				totalCommands,
			)
		}
	})
}

func TestRootCmdWithConfig(t *testing.T) {
	t.Run("should init sub commands", func(t *testing.T) {
		var args []string
		rootCmd, _ := newTestRootCmdWithConfig(args)

		expectedTotal := 3 // foo, bar, completion

		comp, _, _ := rootCmd.Find([]string{"completion"})
		if comp.Name() != "completion" {
			t.Errorf("no '%s' subcommand in the root command", "completion")
		}
		totalCommands := len(rootCmd.Commands())
		if totalCommands != expectedTotal {
			t.Errorf(
				"root cmd has different number of subcommands than expected. Exp: %d, Got: %d",
				expectedTotal,
				totalCommands,
			)
		}
	})

	t.Run("should return exit code 2 for unknown command", func(t *testing.T) {
		rootCmd, _ := newTestRootCmdWithConfig([]string{"fo"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected unknown command error")
		}

		var exitCoder interface{ ExitCode() int }
		if !errors.As(err, &exitCoder) {
			t.Fatal("expected error with exit code")
		}

		if exitCode := exitCoder.ExitCode(); exitCode != 2 {
			t.Fatalf("expected exit code 2, got %d", exitCode)
		}

		if !strings.Contains(err.Error(), `unknown command "fo"`) {
			t.Fatalf("expected unknown command error, got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "Did you mean this?") {
			t.Fatalf("expected suggestions in error, got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "\tfoo\n") {
			t.Fatalf("expected foo suggestion, got %q", err.Error())
		}
	})

	t.Run("should return exit code 2 for unknown command with no suggestions", func(t *testing.T) {
		rootCmd, _ := newTestRootCmdWithConfig([]string{"zzzznotacommand"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected unknown command error")
		}

		var exitCoder interface{ ExitCode() int }
		if !errors.As(err, &exitCoder) {
			t.Fatal("expected error with exit code")
		}

		if exitCode := exitCoder.ExitCode(); exitCode != 2 {
			t.Fatalf("expected exit code 2, got %d", exitCode)
		}

		if !strings.Contains(err.Error(), `unknown command "zzzznotacommand"`) {
			t.Fatalf("expected unknown command error, got %q", err.Error())
		}

		if strings.Contains(err.Error(), "Did you mean this?") {
			t.Fatalf("expected no suggestions, got %q", err.Error())
		}
	})
}

func TestSelfCmd(t *testing.T) {
	t.Run("should return exit code 2 for unknown self subcommand", func(t *testing.T) {
		bufOut := new(bytes.Buffer)

		rootCmd := CreateRootCommand("v0.0.0-test")
		rootCmd.SetArgs([]string{"self", "ls"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufOut)
		InitSelfCmd(rootCmd, "v0.0.0-test")

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected unknown command error")
		}

		var exitCoder interface{ ExitCode() int }
		if !errors.As(err, &exitCoder) {
			t.Fatal("expected error with exit code")
		}

		if exitCode := exitCoder.ExitCode(); exitCode != 2 {
			t.Fatalf("expected exit code 2, got %d", exitCode)
		}

		if !strings.Contains(err.Error(), `unknown command "ls" for "lets self"`) {
			t.Fatalf("expected unknown self subcommand error, got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "Did you mean this?") {
			t.Fatalf("expected suggestions in error, got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "\tlsp\n") {
			t.Fatalf("expected lsp suggestion, got %q", err.Error())
		}
	})

	t.Run("should return exit code 2 for unknown self subcommand with no suggestions", func(t *testing.T) {
		bufOut := new(bytes.Buffer)

		rootCmd := CreateRootCommand("v0.0.0-test")
		rootCmd.SetArgs([]string{"self", "zzzznotacommand"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufOut)
		InitSelfCmd(rootCmd, "v0.0.0-test")

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected unknown command error")
		}

		var exitCoder interface{ ExitCode() int }
		if !errors.As(err, &exitCoder) {
			t.Fatal("expected error with exit code")
		}

		if exitCode := exitCoder.ExitCode(); exitCode != 2 {
			t.Fatalf("expected exit code 2, got %d", exitCode)
		}

		if !strings.Contains(err.Error(), `unknown command "zzzznotacommand" for "lets self"`) {
			t.Fatalf("expected unknown self subcommand error, got %q", err.Error())
		}

		if strings.Contains(err.Error(), "Did you mean this?") {
			t.Fatalf("expected no suggestions, got %q", err.Error())
		}
	})
}
