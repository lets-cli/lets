package cmd

import (
	"github.com/kindritskyiMax/lets/commands"
	"github.com/kindritskyiMax/lets/commands/command"
	"github.com/kindritskyiMax/lets/config"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdVersion returns a cobra command for fetching versions
func newCmdGeneric(cmdToRun command.Command, conf *config.Config, out io.Writer) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   cmdToRun.Name,
		Short: cmdToRun.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.RunCommand(cmdToRun, conf, out)
		},
		DisableFlagParsing: true, // we use docopt to parse flags
		SilenceUsage:       true,
	}
	return cobraCmd
}

// initialize all commands dynamically from config
func initSubCommands(rootCmd *cobra.Command, conf *config.Config, out io.Writer) {
	for _, cmdToRun := range conf.Commands {
		rootCmd.AddCommand(newCmdGeneric(cmdToRun, conf, out))
	}
}
