package parser

import "github.com/lets-cli/lets/config/config"

func parseAndValidateOptions(options interface{}, newCmd *config.Command) error {
	if value, ok := options.(string); ok {
		newCmd.RawOptions = value
	} else {
		return newParseCommandError(
			"must be a string",
			newCmd.Name,
			OPTIONS,
			"",
		)
	}

	return nil
}
