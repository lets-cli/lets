package cmd

import (
	"bytes"
	"context"
	"github.com/spf13/cobra"
	"testing"
)

func newTestRootCmd(args []string) (rootCmd *cobra.Command, out *bytes.Buffer) {
	bufOut := new(bytes.Buffer)

	rootCommand := CreateRootCommand(context.Background(), bufOut, "test-version")
	rootCommand.SetOut(bufOut)
	rootCommand.SetErr(bufOut)
	rootCommand.SetArgs(args)

	return rootCommand, bufOut
}

func TestRootCmd(t *testing.T) {
	// TODO create test file in mem or to /tmp and read it from there
	t.Run("should init sub commands", func(t *testing.T) {
		var args []string
		rootCmd, _ := newTestRootCmd(args)

		expectedTotal := 20

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
