package cmd

import (
	"github.com/spf13/cobra"
)

// InitSelfCmd intializes root 'self' subcommand.
func InitSelfCmd(rootCmd *cobra.Command, version string) error {
	selfCmd := &cobra.Command{
		Use:    "self",
		Hidden: true,
		Short:  "Entrypoint command for self management",
		RunE: func(cmd *cobra.Command, args []string) error {
			return PrintHelpMessage(cmd)
		},
	}

	rootCmd.AddCommand(selfCmd)

	selfCmd.AddCommand(initLspCommand(version))

	return nil
}
