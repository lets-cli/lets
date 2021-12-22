package cmd

import (
	"bytes"
	"testing"

	"github.com/lets-cli/lets/config/config"
	"github.com/spf13/cobra"
)

func newTestRootCmd(args []string) (rootCmd *cobra.Command) {
	rootCommand := CreateRootCommand("v0.0.0-test")
	rootCommand.SetArgs(args)

	return rootCommand
}

func newTestRootCmdWithConfig(args []string) (rootCmd *cobra.Command, out *bytes.Buffer) {
	bufOut := new(bytes.Buffer)

	testCfg := &config.Config{
		Commands: make(map[string]config.Command),
	}
	testCfg.Commands["foo"] = config.Command{
		Name: "foo",
	}
	testCfg.Commands["bar"] = config.Command{
		Name: "bar",
	}

	rootCommand := CreateRootCommandWithConfig(bufOut, testCfg, "v0.0.0-test")
	rootCommand.SetOut(bufOut)
	rootCommand.SetErr(bufOut)
	rootCommand.SetArgs(args)

	return rootCommand, bufOut
}

func TestRootCmd(t *testing.T) {
	t.Run("should init sub commands", func(t *testing.T) {
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
