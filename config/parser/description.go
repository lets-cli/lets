package parser

import "github.com/lets-cli/lets/config/config"

func parseDescription(desc interface{}, newCmd *config.Command) error {
	if value, ok := desc.(string); ok {
		newCmd.Description = value
	} else {
		return parseError(
			"must be a string",
			newCmd.Name,
			DESCRIPTION,
			"",
		)
	}

	return nil
}
