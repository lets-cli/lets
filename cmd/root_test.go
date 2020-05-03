package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"

	"github.com/lets-cli/lets/test"
)

func newTestRootCmd(args []string) (rootCmd *cobra.Command, out *bytes.Buffer) {
	bufOut := new(bytes.Buffer)

	rootCommand := CreateRootCommand(bufOut, "test-version")
	rootCommand.SetOut(bufOut)
	rootCommand.SetErr(bufOut)
	rootCommand.SetArgs(args)

	return rootCommand, bufOut
}

func TestRootCmd(t *testing.T) {
	t.Run("should init sub commands", func(t *testing.T) {
		configRaw := &test.SerializableTestConfig{
			Shell: "bash",
			Commands: map[string]map[string]string{
				"foo": {
					"cmd": "echo foo",
				},
				"bar": {
					"cmd": "echo bar",
				},
			},
		}
		cleanupConfig := test.NewTestConfig(configRaw)
		defer cleanupConfig()

		var args []string
		rootCmd, _ := newTestRootCmd(args)

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
