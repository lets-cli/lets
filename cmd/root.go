package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// newRootCmd represents the base command when called without any subcommands.
func newRootCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "lets",
		Short: "A CLI task runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return PrintHelpMessage(cmd)
		},
		TraverseChildren:   true,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Version:            version,
		// handle errors manually
		SilenceErrors: true,
		// print help message manyally
		SilenceUsage: true,
	}
}

// CreateRootCommand used to run only root command without config.
func CreateRootCommand(version string) *cobra.Command {
	rootCmd := newRootCmd(version)

	initRootFlags(rootCmd)

	return rootCmd
}

func initRootFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringToStringP("env", "E", nil, "set env variable for running command KEY=VALUE")
	rootCmd.Flags().StringArray("only", []string{}, "run only specified command(s) described in cmd as map")
	rootCmd.Flags().StringArray("exclude", []string{}, "run all but excluded command(s) described in cmd as map")
	rootCmd.Flags().Bool("upgrade", false, "upgrade lets to latest version")
	rootCmd.Flags().Bool("init", false, "create a new lets.yaml in the current folder")
	rootCmd.Flags().Bool("no-depends", false, "skip 'depends' for running command")
	rootCmd.Flags().CountP("debug", "d", "show debug logs (or use LETS_DEBUG=1). If used multiple times, shows more verbose logs") //nolint:lll
	rootCmd.Flags().StringP("config", "c", "", "config file (default is lets.yaml)")
	rootCmd.Flags().Bool("all", false, "show all commands (including the ones with _)")
	rootCmd.Flags().Bool("lsp", false, "run lsp server")
}

func PrintHelpMessage(cmd *cobra.Command) error {
	help := cmd.UsageString()
	help = fmt.Sprintf("%s\n\n%s", cmd.Short, help)
	help = strings.Replace(help, "lets [command] --help", "lets help [command]", 1)
	_, err := fmt.Fprint(cmd.OutOrStdout(), help)
	return err
}

func PrintVersionMessage(cmd *cobra.Command) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "lets version %s\n", cmd.Version)
	return err
}
