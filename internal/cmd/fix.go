package cmd

import (
	"os"

	"github.com/lets-cli/lets/internal/config/migrate"
	"github.com/spf13/cobra"
)

func initFixCommand() *cobra.Command {
	var dryRun bool

	fixCmd := &cobra.Command{
		Use:   "fix",
		Short: "Apply lets config migrations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			configName, err := cmd.Root().Flags().GetString("config")
			if err != nil {
				return err
			}

			if configName == "" {
				configName = os.Getenv("LETS_CONFIG")
			}

			_, err = migrate.Fix(configName, os.Getenv("LETS_CONFIG_DIR"), dryRun, cmd.OutOrStdout())

			return err
		},
	}

	fixCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print migrated config without writing files")

	return fixCmd
}
