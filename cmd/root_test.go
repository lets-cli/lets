package cmd

import (
	"bytes"
	"github.com/kindritskyiMax/lets/test"
	"github.com/spf13/cobra"
	"testing"
)

func newTestRootCmd(args []string) (rootCmd *cobra.Command, out *bytes.Buffer) {
	rootCommand := CreateRootCommand()
	bufOut := new(bytes.Buffer)
	rootCommand.SetOut(bufOut)
	rootCommand.SetErr(bufOut)
	rootCommand.SetArgs(args)
	return rootCommand, bufOut
}

func TestRootCmd(t *testing.T) {

	t.Run("run root cmd w/o args", func(t *testing.T) {
		var args []string
		rootCmd, bufOut := newTestRootCmd(args)
		conf := test.GetTestConfig()

		err := Execute(rootCmd, conf, bufOut)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		outStr := bufOut.String()
		if !test.CheckIsDefaultOutput(outStr) {
			t.Errorf("not default output. got: %s", outStr)
		}

	})
}
