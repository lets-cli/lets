package cmd

import (
	"bytes"
	"fmt"
	"github.com/lets-cli/lets/commands"
	"github.com/lets-cli/lets/commands/command"
	"github.com/lets-cli/lets/config"
	"github.com/spf13/cobra"
	"io"
)

// newCmdGeneric creates new cobra root sub command from Command
func newCmdGeneric(cmdToRun command.Command, conf *config.Config, out io.Writer) *cobra.Command {
	subCmd := &cobra.Command{
		Use:   cmdToRun.Name,
		Short: cmdToRun.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.RunCommand(cmdToRun, conf, out)
		},
		DisableFlagParsing: true, // we use docopt to parse flags
		SilenceUsage:       true,
	}
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
func initSubCommands(rootCmd *cobra.Command, conf *config.Config, out io.Writer) {
	for _, cmdToRun := range conf.Commands {
		rootCmd.AddCommand(newCmdGeneric(cmdToRun, conf, out))
	}
}
