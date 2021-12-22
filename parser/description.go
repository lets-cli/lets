package parser

import "github.com/lets-cli/lets/config"

func parseAndValidateDescription(desc interface{}, newCmd *config.Command) error {
	if value, ok := desc.(string); ok {
		newCmd.Description = value
	} else {
		return newParseCommandError(
			"must be a string",
			newCmd.Name,
			DESCRIPTION,
			"",
		)
	}

	return nil
}
