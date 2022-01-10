package parser

import (
	"fmt"

	"github.com/lets-cli/lets/config/config"
)

func parseShell(rawShell interface{}, newCmd *config.Command) error {
	shell, ok := rawShell.(string)
	if !ok {
		return fmt.Errorf("shell must be a string")
	}

	newCmd.Shell = shell

	return nil
}
