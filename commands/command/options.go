package command

func parseAndValidateOptions(options interface{}, newCmd *Command) error {
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
