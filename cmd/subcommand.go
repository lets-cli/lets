package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/lets-cli/lets/commands/command"
	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/runner"
)

// newCmdGeneric creates new cobra root sub command from Command
func newCmdGeneric(ctx context.Context, cmdToRun command.Command, conf *config.Config, out io.Writer) *cobra.Command {
	subCmd := &cobra.Command{
		Use:   cmdToRun.Name,
		Short: cmdToRun.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.RunCommand(ctx, cmdToRun, conf, out)
		},
		DisableFlagParsing: true, // we use docopt to parse flags
		SilenceUsage:       true,
	}

	// try print docopt as help for command
	subCmd.SetHelpFunc(func(c *cobra.Command, strings []string) {
		buf := new(bytes.Buffer)
		if cmdToRun.Description != "" {
			buf.WriteString(fmt.Sprintf("%s\n\n", cmdToRun.Description))
		}
		buf.WriteString(cmdToRun.RawOptions)

		_, err := buf.WriteTo(c.OutOrStdout())
		if err != nil {
			c.Println(err)
		}
	})

	return subCmd
}

// initialize all commands dynamically from config
func initSubCommands(ctx context.Context, rootCmd *cobra.Command, conf *config.Config, out io.Writer) {
	for _, cmdToRun := range conf.Commands {
		rootCmd.AddCommand(newCmdGeneric(ctx, cmdToRun, conf, out))
	}
}
