package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/logging"
	"github.com/lets-cli/lets/workdir"
)

// CreateRootCommand is where all the stuff begins
func CreateRootCommand(ctx context.Context, out io.Writer, version string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lets",
		Short: "A CLI command runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoot(cmd)
		},
		Version: version,
	}

	initRootCommand(ctx, rootCmd, out, version)

	return rootCmd
}

func initRootCommand(ctx context.Context, rootCmd *cobra.Command, out io.Writer, version string) {
	configPath, workDir := env.GetConfigPathFromEnv()

	if configPath == "" {
		configPath = config.GetDefaultConfigPath()
	}

	conf, cfgErr := config.Load(configPath, workDir, version)
	if cfgErr != nil {
		initErrCheck(rootCmd, cfgErr)
	} else {
		// create .lets only when there is valid config in work dir
		if createDirErr := workdir.CreateDotLetsDir(); createDirErr != nil {
			initErrCheck(rootCmd, createDirErr)
		}

		initSubCommands(ctx, rootCmd, conf, out)
	}

	initCompletionCmd(rootCmd)
	initVersionFlag(rootCmd)
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

func runRoot(cmd *cobra.Command) error {
	return cmd.Help()
}
