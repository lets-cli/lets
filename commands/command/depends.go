package command

func parseAndValidateDepends(depends interface{}, newCmd *Command) error {
	if depends, ok := depends.([]interface{}); ok {
		for _, value := range depends {
			if value, ok := value.(string); ok {
				// TODO validate if command is really exists - in validate
				newCmd.Depends = append(newCmd.Depends, value)
			} else {
				return newCommandError(
					"value of depends list must be a string",
					newCmd.Name,
					DEPENDS,
					"",
				)
			}
		}
	} else {
		return newCommandError(
			"must be a list of string (commands)",
			newCmd.Name,
			DEPENDS,
			"",
		)
	}
	return nil
}
