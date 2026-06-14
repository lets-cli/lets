package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lets-cli/lets/internal/util"
	"github.com/spf13/cobra"
)

func initConfigCommand(openEditor func(string) error) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage lets user config",
		Long: strings.TrimSpace(`Manage the per-user lets settings file.

The user config is stored at ~/.config/lets/config.yaml and applies to all
projects on the machine.`),
		Args: validateCommandArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	configCmd.AddCommand(initConfigPathCommand())
	configCmd.AddCommand(initConfigEditCommand(openEditor))

	return configCmd
}

func initConfigPathCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print lets user config path",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := util.LetsUserFile("config.yaml")
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), path)

			return err
		},
	}
}

func initConfigEditCommand(openEditor func(string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Open lets user config in EDITOR",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := util.LetsUserFile("config.yaml")
			if err != nil {
				return err
			}

			if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
				return fmt.Errorf("creating config directory: %w", err)
			}

			if err := openEditor(path); err != nil {
				return fmt.Errorf("can not open config: %w", err)
			}

			return nil
		},
	}
}
