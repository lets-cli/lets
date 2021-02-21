package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/upgrade"
	"github.com/lets-cli/lets/upgrade/registry"
)

// CreateRootCommand is where all the stuff begins
func CreateRootCommand(out io.Writer, cfg *config.Config) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lets",
		Short: "A CLI command runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoot(cmd)
		},
		TraverseChildren: true,
		Version:          cfg.Version,
		SilenceErrors:    true,
		SilenceUsage:     true,
	}

	initRootCommand(rootCmd, out, cfg)

	return rootCmd
}

func initRootCommand(rootCmd *cobra.Command, out io.Writer, cfg *config.Config) {
	initSubCommands(rootCmd, cfg, out)
	initCompletionCmd(rootCmd)
	initVersionFlag(rootCmd)
	initEnvFlag(rootCmd)
	initOnlyAndExecFlags(rootCmd)
	initUpgradeFlag(rootCmd)
}

func initVersionFlag(rootCmd *cobra.Command) {
	rootCmd.Flags().BoolP("version", "v", false, "version for lets")
}

func initEnvFlag(rootCmd *cobra.Command) {
	rootCmd.Flags().StringToStringP("env", "E", nil, "set env variable for running command KEY=VALUE")
}

func initOnlyAndExecFlags(cmd *cobra.Command) {
	cmd.Flags().StringArray("only", []string{}, "run only specified command(s) described in cmd as map")
	cmd.Flags().StringArray("exclude", []string{}, "run all but excluded command(s) described in cmd as map")
}

func initUpgradeFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("upgrade", false, "upgrade lets to latest version")
}

func runRoot(cmd *cobra.Command) error {
	selfUpgrade, err := cmd.Flags().GetBool("upgrade")
	if err != nil {
		return err
	}

	if selfUpgrade {
		return upgrade.Upgrade(registry.NewGithubRegistry())
	}
	return cmd.Help()
}
