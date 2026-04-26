package cmd

import (
	"github.com/lets-cli/lets/internal/settings"
	"github.com/lets-cli/lets/internal/util"
	"github.com/spf13/cobra"
)

// InitSelfCmd intializes root 'self' subcommand.
func InitSelfCmd(rootCmd *cobra.Command, version string, appSettings settings.Settings) {
	initSelfCmd(rootCmd, version, appSettings, util.OpenURL, util.OpenEditor)
}

func initSelfCmd(
	rootCmd *cobra.Command,
	version string,
	appSettings settings.Settings,
	openURL func(string) error,
	openEditor func(string) error,
) {
	selfCmd := &cobra.Command{
		Use:     "self",
		Hidden:  false,
		Short:   "Manage lets CLI itself",
		GroupID: "internal",
		Args:    validateCommandArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(selfCmd)

	selfCmd.AddCommand(initConfigCommand(openEditor))
	selfCmd.AddCommand(initDocCommand(openURL))
	selfCmd.AddCommand(initFixCommand())
	selfCmd.AddCommand(initLspCommand(version))
	selfCmd.AddCommand(initSkillsCommand())
	selfCmd.AddCommand(initUpgradeCommand(version, appSettings))
}
