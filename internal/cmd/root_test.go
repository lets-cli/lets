package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/lets-cli/lets/internal/config/config"
	"github.com/lets-cli/lets/internal/upgrade"
	"github.com/spf13/cobra"
)

type testErrorWithExitCode interface {
	error
	ExitCode() int
}

func requireErrorWithExitCode(t *testing.T, err error) testErrorWithExitCode {
	t.Helper()

	errWithExitCode, ok := errors.AsType[testErrorWithExitCode](err)
	if !ok {
		t.Fatal("expected error with exit code")
	}

	return errWithExitCode
}

func newTestRootCmd(args []string) (rootCmd *cobra.Command) {
	root := CreateRootCommand("v0.0.0-test", "")
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

	root := CreateRootCommand("v0.0.0-test", "")
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

		errWithExitCode := requireErrorWithExitCode(t, err)

		if exitCode := errWithExitCode.ExitCode(); exitCode != 2 {
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

		errWithExitCode := requireErrorWithExitCode(t, err)

		if exitCode := errWithExitCode.ExitCode(); exitCode != 2 {
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

func TestPrintVersionMessage(t *testing.T) {
	t.Run("should include build date in parentheses when non-empty", func(t *testing.T) {
		buf := new(bytes.Buffer)
		root := CreateRootCommand("v0.0.0-test", "2024-01-15T10:30:00Z")
		root.SetOut(buf)

		err := PrintVersionMessage(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "v0.0.0-test") {
			t.Errorf("expected version in output, got %q", out)
		}
		if !strings.Contains(out, "(2024-01-15T10:30:00Z)") {
			t.Errorf("expected build date in parentheses, got %q", out)
		}
	})

	t.Run("should omit parentheses when build date is empty", func(t *testing.T) {
		buf := new(bytes.Buffer)
		root := CreateRootCommand("v0.0.0-test", "")
		root.SetOut(buf)

		err := PrintVersionMessage(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		out := buf.String()
		if strings.Contains(out, "(") {
			t.Errorf("expected no parentheses when build date is empty, got %q", out)
		}
	})
}

func TestSelfCmd(t *testing.T) {
	t.Run("should open documentation in browser", func(t *testing.T) {
		bufOut := new(bytes.Buffer)
		called := false
		gotURL := ""

		openURL := func(url string) error {
			called = true
			gotURL = url

			return nil
		}

		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "doc"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufOut)
		initSelfCmd(rootCmd, "v0.0.0-test", openURL)

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !called {
			t.Fatal("expected documentation opener to be called")
		}

		if gotURL != letsDocsURL {
			t.Fatalf("expected docs url %q, got %q", letsDocsURL, gotURL)
		}
	})

	t.Run("should return opener error for documentation command", func(t *testing.T) {
		bufOut := new(bytes.Buffer)

		openURL := func(url string) error {
			return errors.New("open failed")
		}

		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "doc"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufOut)
		initSelfCmd(rootCmd, "v0.0.0-test", openURL)

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected documentation opener error")
		}

		if !strings.Contains(err.Error(), "can not open documentation") {
			t.Fatalf("expected documentation error, got %q", err.Error())
		}
	})

	t.Run("should return exit code 2 for unknown self subcommand", func(t *testing.T) {
		bufOut := new(bytes.Buffer)

		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "ls"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufOut)
		InitSelfCmd(rootCmd, "v0.0.0-test")

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected unknown command error")
		}

		errWithExitCode := requireErrorWithExitCode(t, err)

		if exitCode := errWithExitCode.ExitCode(); exitCode != 2 {
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

		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "zzzznotacommand"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufOut)
		InitSelfCmd(rootCmd, "v0.0.0-test")

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected unknown command error")
		}

		errWithExitCode := requireErrorWithExitCode(t, err)

		if exitCode := errWithExitCode.ExitCode(); exitCode != 2 {
			t.Fatalf("expected exit code 2, got %d", exitCode)
		}

		if !strings.Contains(err.Error(), `unknown command "zzzznotacommand" for "lets self"`) {
			t.Fatalf("expected unknown self subcommand error, got %q", err.Error())
		}

		if strings.Contains(err.Error(), "Did you mean this?") {
			t.Fatalf("expected no suggestions, got %q", err.Error())
		}
	})

	t.Run("should run self upgrade subcommand", func(t *testing.T) {
		bufOut := new(bytes.Buffer)
		called := false

		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "upgrade"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufOut)
		selfCmd := &cobra.Command{
			Use:   "self",
			Short: "Manage lets CLI itself",
		}
		rootCmd.AddCommand(selfCmd)

		selfCmd.AddCommand(initUpgradeCommandWith(func() (upgrade.Upgrader, error) {
			return mockUpgraderFunc(func(ctx context.Context) error {
				called = true

				return nil
			}), nil
		}))

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !called {
			t.Fatal("expected upgrader to be called")
		}
	})

	t.Run("should return upgrader error for self upgrade command", func(t *testing.T) {
		bufOut := new(bytes.Buffer)

		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "upgrade"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufOut)
		selfCmd := &cobra.Command{
			Use:   "self",
			Short: "Manage lets CLI itself",
		}
		rootCmd.AddCommand(selfCmd)

		selfCmd.AddCommand(initUpgradeCommandWith(func() (upgrade.Upgrader, error) {
			return mockUpgraderFunc(func(ctx context.Context) error {
				return errors.New("upgrade failed")
			}), nil
		}))

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected upgrader error")
		}

		if !strings.Contains(err.Error(), "can not self-upgrade binary") {
			t.Fatalf("expected self-upgrade error, got %q", err.Error())
		}
	})
}

type mockUpgraderFunc func(context.Context) error

func (f mockUpgraderFunc) Upgrade(ctx context.Context) error {
	return f(ctx)
}
