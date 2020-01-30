package cmd

import (
	"github.com/kindritskyiMax/lets/config"
	"github.com/spf13/cobra"
	"io"
)

type letsOptions struct {
	config  string
	rawArgs []string
}

// CreateRootCommand is where all the stuff begins
func CreateRootCommand() *cobra.Command {
	var opts letsOptions
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lets",
		Short: "A CLI command runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHelp(cmd)
		},
	}

	initRootCmdFlags(rootCmd, &opts)

	return rootCmd
}

func Execute(cmd *cobra.Command, conf *config.Config, out io.Writer) error {
	initSubCommands(cmd, conf, out)
	return cmd.Execute()
}

func runHelp(cmd *cobra.Command) error {
	return cmd.Help()
}

func initRootCmdFlags(rootCmd *cobra.Command, opts *letsOptions) {
	//cobra.OnInitialize(initConfig)
}
