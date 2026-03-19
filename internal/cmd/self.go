package cmd

import (
	"github.com/lets-cli/lets/internal/util"
	"github.com/spf13/cobra"
)

var openURL = util.OpenURL

// InitSelfCmd intializes root 'self' subcommand.
func InitSelfCmd(rootCmd *cobra.Command, version string) {
	selfCmd := &cobra.Command{
		Use:     "self",
		Hidden:  false,
		Short:   "Manage lets CLI itself",
		GroupID: "internal",
		Args:    validateCommandArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return PrintHelpMessage(cmd)
		},
	}

	rootCmd.AddCommand(selfCmd)

	selfCmd.AddCommand(initDocCommand(openURL))
	selfCmd.AddCommand(initLspCommand(version))
}
