package cmd

import (
	"github.com/spf13/cobra"
)

// InitSelfCmd intializes root 'self' subcommand.
func InitSelfCmd(rootCmd *cobra.Command, version string) {
	selfCmd := &cobra.Command{
		Use:    "self",
		Hidden: false,
		Short:  "Manage lets CLI itself",
		RunE: func(cmd *cobra.Command, args []string) error {
			return PrintHelpMessage(cmd)
		},
	}

	rootCmd.AddCommand(selfCmd)

	selfCmd.AddCommand(initLspCommand(version))
}
