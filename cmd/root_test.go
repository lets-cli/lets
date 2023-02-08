package cmd

import (
	"bytes"
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
}
