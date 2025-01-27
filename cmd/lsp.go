package cmd

import (
	"github.com/lets-cli/lets/lsp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func initLspCommand(version string) *cobra.Command {
	lspCmd := &cobra.Command{
		Use:   "lsp",
		Short: "Language Server Protocol (LSP) server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := lsp.Run(cmd.Context(), version); err != nil {
				return errors.Wrap(err, "lsp error")
			}
			return nil
		},
	}

	return lspCmd
}
