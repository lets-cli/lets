package cmd

import (
	"fmt"
	"io"

	"github.com/lets-cli/lets/internal/fetch"
	"github.com/lets-cli/lets/internal/progressbar"
	"github.com/lets-cli/lets/internal/settings"
	"github.com/lets-cli/lets/internal/upgrade"
	"github.com/lets-cli/lets/internal/upgrade/registry"
	"github.com/lets-cli/lets/internal/util"
	"github.com/spf13/cobra"
)

type upgraderFactory func(cmd *cobra.Command) (upgrade.Upgrader, error)

func initUpgradeCommand(version string, appSettings settings.Settings) *cobra.Command {
	return initUpgradeCommandWith(func(cmd *cobra.Command) (upgrade.Upgrader, error) {
		progress := upgradeProgress(cmd.ErrOrStderr(), appSettings)

		return upgrade.NewBinaryUpgrader(registry.NewGithubRegistry(), version, upgrade.WithProgress(progress))
	})
}

// upgradeProgress builds a progress observer for the upgrade download, or a
// no-op observer when the stream is not an interactive terminal.
func upgradeProgress(stderr io.Writer, appSettings settings.Settings) fetch.ProgressObserver { //nolint:ireturn // Returns NopObserver or *progressbar.Observer.
	if !util.IsTerminalWriter(stderr) {
		return fetch.NopObserver{}
	}

	return progressbar.New(
		stderr,
		progressbar.WithNoColor(appSettings.NoColor),
		progressbar.WithTheme(appSettings.Theme),
	)
}

func initUpgradeCommandWith(createUpgrader upgraderFactory) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade lets to latest version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrader, err := createUpgrader(cmd)
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
