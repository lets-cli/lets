package cmd

import (
	"fmt"

	"github.com/lets-cli/lets/internal/upgrade"
	"github.com/lets-cli/lets/internal/upgrade/registry"
	"github.com/spf13/cobra"
)

type upgraderFactory func() (upgrade.Upgrader, error)

func initUpgradeCommand(version string) *cobra.Command {
	return initUpgradeCommandWith(func() (upgrade.Upgrader, error) {
		return upgrade.NewBinaryUpgrader(registry.NewGithubRegistry(), version)
	})
}

func initUpgradeCommandWith(createUpgrader upgraderFactory) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade lets to latest version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrader, err := createUpgrader()
			if err != nil {
				return fmt.Errorf("can not self-upgrade binary: %w", err)
			}

			if err := upgrader.Upgrade(cmd.Context()); err != nil {
				return fmt.Errorf("can not self-upgrade binary: %w", err)
			}

			return nil
		},
	}

	return upgradeCmd
}
