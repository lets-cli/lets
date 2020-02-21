package cmd

import (
	"fmt"
	"github.com/kindritskyiMax/lets/config"
	"github.com/kindritskyiMax/lets/logging"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
)

func isDebug() bool {
	debug, err := strconv.ParseBool(os.Getenv("LETS_DEBUG"))
	if err != nil {
		return false
	}
	return debug
}

// CreateRootCommand is where all the stuff begins
func CreateRootCommand(out io.Writer, version string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lets",
		Short: "A CLI command runner",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHelp(cmd)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initLogging(isDebug())
			return nil
		},
		Version: version,
	}
	// workaround to hide help sub command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "__lets-no-help",
		Hidden: true,
	})

	configPath, workDir := config.GetConfigPathFromEnv()

	if configPath == "" {
		configPath = config.GetDefaultConfigPath()
	}

	conf, err := config.Load(configPath, workDir)

	if err != nil {
		InitConfigErrCheck(rootCmd, err)
	} else {
		initSubCommands(rootCmd, conf, out)
	}

	return rootCmd
}

// InitConfigErrCheck check if config load failed with error, if so, print error and exit
// Doing it in PreRun allows us run root cmd as usual, parse help flags
// and only if no command were run and config load has failed - we print error
func InitConfigErrCheck(rootCmd *cobra.Command, cfgErr error) {
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cfgErr != nil {
			fmt.Printf("Error: %s\n", cfgErr)
			os.Exit(1)
		}
	}
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
