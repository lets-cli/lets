package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/upgrade"
	"github.com/lets-cli/lets/upgrade/registry"
	"github.com/lets-cli/lets/workdir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// newRootCmd represents the base command when called without any subcommands.
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

// CreateRootCommandWithConfig used to run root command with all subcommands.
func CreateRootCommandWithConfig(out io.Writer, cfg *config.Config, version string) *cobra.Command {
	rootCmd := newRootCmd(version)

	initRootCommand(rootCmd, cfg)
	initSubCommands(rootCmd, cfg, out)

	return rootCmd
}

// CreateRootCommand used to run only root command without config.
func CreateRootCommand(version string) *cobra.Command {
	rootCmd := newRootCmd(version)

	initRootCommand(rootCmd, nil)

	return rootCmd
}

// ConfigErrorCheck will print error only if no args passed
// Main reason to do it in PreRun allows us to run root cmd as usual,
//	parse help flags if any provided or check if its help command.
//
// For example if config load failed with error (no lets.yaml in current dir) - print error and exit.
func ConfigErrorCheck(rootCmd *cobra.Command, err error) {
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() > 0 {
			return
		}

		log.Error(err)
		os.Exit(1)
	}
}

func initRootCommand(rootCmd *cobra.Command, cfg *config.Config) {
	initCompletionCmd(rootCmd, cfg)
	rootCmd.Flags().StringToStringP("env", "E", nil, "set env variable for running command KEY=VALUE")
	rootCmd.Flags().StringArray("only", []string{}, "run only specified command(s) described in cmd as map")
	rootCmd.Flags().StringArray("exclude", []string{}, "run all but excluded command(s) described in cmd as map")
	rootCmd.Flags().Bool("upgrade", false, "upgrade lets to latest version")
	rootCmd.Flags().Bool("init", false, "create a new lets.yaml in the current folder")
	rootCmd.Flags().Bool("no-depends", false, "skip 'depends' for running command")
}


func printHelpMessage(cmd *cobra.Command) error {
	help := cmd.UsageString()
	help = strings.Replace(help, "lets [command] --help", "lets help [command]", 1)
	_, err := fmt.Fprint(cmd.OutOrStdout(), help)
	return err
}

func runRoot(cmd *cobra.Command, version string) error {
	selfUpgrade, err := cmd.Flags().GetBool("upgrade")
	if err != nil {
		return fmt.Errorf("can not get flag 'upgrade': %w", err)
	}

	if selfUpgrade {
		upgrader, err := upgrade.NewBinaryUpgrader(registry.NewGithubRegistry(cmd.Context()), version)
		if err != nil {
			return fmt.Errorf("can not upgrade lets: %w", err)
		}

		return upgrader.Upgrade()
	}

	init, err := cmd.Flags().GetBool("init")
	if err != nil {
		return fmt.Errorf("can not get flag 'init': %w", err)
	}

	if init {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		if err := workdir.InitLetsFile(wd, version); err != nil {
			log.Fatal(err)
		}

		return nil
	}

	return printHelpMessage(cmd)
}
