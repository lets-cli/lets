package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/upgrade"
	"github.com/lets-cli/lets/upgrade/registry"
	"github.com/lets-cli/lets/logging"
)

// newRootCmd represents the base command when called without any subcommands
func newRootCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "lets",
		Short: "A CLI command runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoot(cmd, version)
		},
		TraverseChildren: true,
		Version:          version,
		SilenceErrors:    true,
		SilenceUsage:     true,
	}
}

// CreateRootCommandWithConfig used to run root command with all subcommands
func CreateRootCommandWithConfig(out io.Writer, cfg *config.Config, version string) *cobra.Command {
	var rootCmd = newRootCmd(version)

	initRootCommand(rootCmd)
	initSubCommands(rootCmd, cfg, out)

	return rootCmd
}

// CreateRootCommand used to run only root command without config
func CreateRootCommand(version string) *cobra.Command {
	var rootCmd = newRootCmd(version)

	initRootCommand(rootCmd)

	return rootCmd
}

// ConfigErrorCheck will print error only if no args passed
func ConfigErrorCheck(rootCmd *cobra.Command, err error) {
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() > 0 {
			return
		}
		
		logging.Log.Error(err)
		os.Exit(1)
	}
}

func initRootCommand(rootCmd *cobra.Command) {
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

func runRoot(cmd *cobra.Command, version string) error {
	selfUpgrade, err := cmd.Flags().GetBool("upgrade")
	if err != nil {
		return err
	}

	if selfUpgrade {
		upgrader, err := upgrade.NewBinaryUpgrader(registry.NewGithubRegistry(), version)
		if err != nil {
			return err
		}

		return upgrader.Upgrade()
	}

	showVersion, err := cmd.Flags().GetBool("version")
	if err != nil {
		return err
	}

	if showVersion {
		logging.Log.Printf("lets version %s", version)
		return nil
	}

	return cmd.Help()
}
