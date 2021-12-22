package parser

import "github.com/lets-cli/lets/config"

func parseAndValidateDepends(depends interface{}, newCmd *config.Command) error {
	if depends, ok := depends.([]interface{}); ok {
		for _, value := range depends {
			if value, ok := value.(string); ok {
				// TODO validate if command is really exists - in validate
				newCmd.Depends = append(newCmd.Depends, value)
			} else {
				return newParseCommandError(
					"value of depends list must be a string",
					newCmd.Name,
					DEPENDS,
					"",
				)
			}
		}
	} else {
		return newParseCommandError(
			"must be a list of string (commands)",
			newCmd.Name,
			DEPENDS,
			"",
		)
	}

	return nil
}
