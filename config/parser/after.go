package parser

import "github.com/lets-cli/lets/config/config"

func parseAfter(after interface{}, newCmd *config.Command) error {
	switch after := after.(type) {
	case string:
		newCmd.After = after
	default:
		return parseError(
			"must be a string",
			newCmd.Name,
			AFTER,
			"",
		)
	}

	return nil
}
