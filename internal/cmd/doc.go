package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const letsDocsURL = "https://lets-cli.org/docs/config"

func initDocCommand(openURL func(string) error) *cobra.Command {
	docCmd := &cobra.Command{
		Use:     "doc",
		Aliases: []string{"docs"},
		Short:   "Open lets documentation in browser",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := openURL(letsDocsURL); err != nil {
				return fmt.Errorf("can not open documentation: %w", err)
			}

			return nil
		},
	}

	return docCmd
}
