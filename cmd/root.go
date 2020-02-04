package cmd

import (
	"github.com/kindritskyiMax/lets/config"
	"github.com/kindritskyiMax/lets/logging"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// CreateRootCommand is where all the stuff begins
func CreateRootCommand(conf *config.Config, out io.Writer) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lets",
		Short: "A CLI command runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHelp(cmd)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initLogging(os.Getenv("LETS_DEBUG") == "true")
			return nil
		},
	}
	// workaround to hide help sub command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	initSubCommands(rootCmd, conf, out)

	return rootCmd
}

func Execute(cmd *cobra.Command) error {
	return cmd.Execute()
}

func runHelp(cmd *cobra.Command) error {
	return cmd.Help()
}

func initLogging(verbose bool) {
	logger := logging.Log

	logger.Level = log.WarnLevel

	if verbose {
		logger.Level = log.DebugLevel
	}
	logger.Out = os.Stderr

	formatter := &logging.Formatter{}
	log.SetFormatter(formatter)
	logger.Formatter = formatter
}
