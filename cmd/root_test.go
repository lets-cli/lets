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
	t.Run("should init sub commands", func(t *testing.T) {
		var args []string
		rootCmd, _ := newTestRootCmd(args)

		comp, _, _ := rootCmd.Find([]string{"completion"})
		if comp.Name() != "completion" {
			t.Errorf("no '%s' subcommand in the root command", "completion")
		}
		if len(rootCmd.Commands()) != 1 {
			t.Errorf("root cmd has different number of subcommands than expected. Exp: %d, Got: %d", 1, len(rootCmd.Commands()))
		}

		// TODO add all other subcommands
	})
}
