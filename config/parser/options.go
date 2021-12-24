package parser

import "github.com/lets-cli/lets/config/config"

func parseOptions(options interface{}, newCmd *config.Command) error {
	if value, ok := options.(string); ok {
		newCmd.RawOptions = value
	} else {
		return parseError(
			"must be a string",
			newCmd.Name,
			OPTIONS,
			"",
		)
	}

	return nil
}
