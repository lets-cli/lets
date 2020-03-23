package cmd

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/logging"
)

// CreateRootCommand is where all the stuff begins
func CreateRootCommand(_ context.Context, out io.Writer, version string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lets",
		Short: "A CLI command runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHelp(cmd)
		},
		Version: version,
	}

	configPath, workDir := env.GetConfigPathFromEnv()

	if configPath == "" {
		configPath = config.GetDefaultConfigPath()
	}

	conf, err := config.Load(configPath, workDir)

	if err != nil {
		InitConfigErrCheck(rootCmd, err)
	} else {
		initSubCommands(rootCmd, conf, out)
	}

	initCompletionCmd(rootCmd)

	return rootCmd
}

// InitConfigErrCheck check if config load failed with error, if so, print error and exit
// Doing it in PreRun allows us run root cmd as usual, parse help flags
// and only if no command were run and config load has failed - we print error
func InitConfigErrCheck(rootCmd *cobra.Command, cfgErr error) {
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cfgErr != nil {
			logging.Log.Errorf("error: %s\n", cfgErr)
			os.Exit(1)
		}
	}
}

func runHelp(cmd *cobra.Command) error {
	return cmd.Help()
}
