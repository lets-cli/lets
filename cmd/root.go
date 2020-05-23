package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/logging"
	"github.com/lets-cli/lets/workdir"
)

// CreateRootCommand is where all the stuff begins
func CreateRootCommand(out io.Writer, version string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lets",
		Short: "A CLI command runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoot(cmd)
		},
		TraverseChildren: true,
		Version: version,
	}

	initRootCommand(rootCmd, out, version)

	return rootCmd
}

func initRootCommand(rootCmd *cobra.Command, out io.Writer, version string) {
	var conf *config.Config

	configPath, findCfgErr := config.FindConfig()
	if findCfgErr != nil {
		initErrCheck(rootCmd, findCfgErr)
	} else {
		cfg, cfgErr := config.Load(configPath, version)
		if cfgErr != nil {
			initErrCheck(rootCmd, cfgErr)
		}
		conf = cfg
	}

	if conf != nil {
		// create .lets only when there is valid config in work dir
		if createDirErr := workdir.CreateDotLetsDir(configPath.WorkDir); createDirErr != nil {
			initErrCheck(rootCmd, createDirErr)
		}

		initSubCommands(rootCmd, conf, out)
	}

	initCompletionCmd(rootCmd)
	initVersionFlag(rootCmd)
	initEnvFlag(rootCmd)
}

// InitErrCheck check if error occurred before root cmd execution.
// Main reason to do it in PreRun allows us to run root cmd as usual,
//	parse help flags if any provided or check if its help command.
//
// For example if config load failed with error (no lets.yaml in current dir) - print error and exit.
func initErrCheck(rootCmd *cobra.Command, err error) {
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if err != nil {
			logging.Log.Fatal(err)
		}
	}
}

func initVersionFlag(rootCmd *cobra.Command) {
	rootCmd.Flags().BoolP("version", "v", false, "version for lets")
}

func initEnvFlag(rootCmd *cobra.Command) {
	rootCmd.Flags().StringToStringP("env", "E", nil, "set env variable for running command KEY=VALUE")
}

func runRoot(cmd *cobra.Command) error {
	return cmd.Help()
}
