package cmd

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
	"github.com/lets-cli/fang"
	"github.com/lets-cli/lets/internal/config"
	"github.com/lets-cli/lets/internal/executor"
	"github.com/lets-cli/lets/internal/theme"
	"github.com/spf13/cobra"
)

func newGoldenRoot(t *testing.T) *cobra.Command {
	t.Helper()
	t.Setenv("__FANG_TEST_WIDTH", "80")
	root := CreateRootCommand("0.0.0-dev", "")
	root.InitDefaultHelpFlag()
	root.InitDefaultVersionFlag()
	InitSelfCmd(root, "0.0.0-dev")
	InitCompletionCmd(root, nil)
	root.InitDefaultHelpCmd()
	return root
}

func newGoldenRootFromYAML(t *testing.T, fixture string) *cobra.Command {
	t.Helper()
	root := newGoldenRoot(t)
	absPath, err := filepath.Abs(filepath.Join("testdata", "fixtures", fixture))
	if err != nil {
		t.Fatalf("fixture path: %v", err)
	}
	conf, err := config.Load(absPath, "", "0.0.0-dev")
	if err != nil {
		t.Fatalf("load config %s: %v", fixture, err)
	}
	InitSubCommands(root, conf, false, nil)
	return root
}

func goldenExecute(t *testing.T, root *cobra.Command, args []string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs(args)
	err := fang.Execute(
		context.Background(), root,
		fang.WithVersion(root.Version),
		fang.WithColorSchemeFunc(theme.DefaultColorScheme),
		fang.WithErrorHandler(ErrorHandler),
		fang.WithHelpRenderer(HelpRenderer),
	)
	if err != nil {
		golden.RequireEqual(t, stderr.Bytes())
		return
	}
	golden.RequireEqual(t, stdout.Bytes())
}

// TestHelpGolden covers root-level help rendering for various command configurations.
// To update golden files: go test ./internal/cmd/ -run TestHelpGolden -update
func TestHelpGolden(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		root := newGoldenRootFromYAML(t, "basic.yaml")
		goldenExecute(t, root, []string{"--help"})
	})

	t.Run("long_command", func(t *testing.T) {
		root := newGoldenRootFromYAML(t, "long_command.yaml")
		goldenExecute(t, root, []string{"--help"})
	})

	t.Run("grouped_commands", func(t *testing.T) {
		root := newGoldenRootFromYAML(t, "grouped_commands.yaml")
		goldenExecute(t, root, []string{"--help"})
	})

	t.Run("grouped_commands_long", func(t *testing.T) {
		root := newGoldenRootFromYAML(t, "grouped_commands_long.yaml")
		goldenExecute(t, root, []string{"--help"})
	})
}

// TestCommandHelpGolden covers per-subcommand help rendering with docopt options.
// To update golden files: go test ./internal/cmd/ -run TestCommandHelpGolden -update
func TestCommandHelpGolden(t *testing.T) {
	t.Run("long_description", func(t *testing.T) {
		root := newGoldenRootFromYAML(t, "command_help.yaml")
		goldenExecute(t, root, []string{"help", "test"})
	})

	t.Run("short_only", func(t *testing.T) {
		root := newGoldenRootFromYAML(t, "command_help.yaml")
		goldenExecute(t, root, []string{"help", "test2"})
	})
}

// TestErrorGolden covers error rendering: unknown commands and dependency trees.
// Error tests use synthetic commands to inject DependencyError directly,
// keeping them focused on the renderer rather than shell execution.
// To update golden files: go test ./internal/cmd/ -run TestErrorGolden -update
func TestErrorGolden(t *testing.T) {
	t.Run("command_not_found", func(t *testing.T) {
		root := newGoldenRootFromYAML(t, "basic.yaml")
		goldenExecute(t, root, []string{"zzzznotacommand"})
	})

	t.Run("command_not_found_with_suggestion", func(t *testing.T) {
		// "self" is added by newGoldenRoot; "slef" triggers the suggestion
		root := newGoldenRoot(t)
		goldenExecute(t, root, []string{"slef"})
	})

	t.Run("dependency_single", func(t *testing.T) {
		root := newGoldenRoot(t)
		root.AddCommand(&cobra.Command{
			Use:     "lint",
			GroupID: "main",
			RunE: func(*cobra.Command, []string) error {
				return &executor.DependencyError{
					Chain: []string{"lint"},
					Err:   fmt.Errorf("exit status 1"),
				}
			},
		})
		goldenExecute(t, root, []string{"lint"})
	})

	t.Run("dependency_chain", func(t *testing.T) {
		root := newGoldenRoot(t)
		root.AddCommand(&cobra.Command{
			Use:     "deploy",
			GroupID: "main",
			RunE: func(*cobra.Command, []string) error {
				return &executor.DependencyError{
					Chain: []string{"deploy", "build", "lint"},
					Err:   fmt.Errorf("exit status 1"),
				}
			},
		})
		goldenExecute(t, root, []string{"deploy"})
	})
}
