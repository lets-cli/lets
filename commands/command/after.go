package command

func parseAndValidateAfter(after interface{}, newCmd *Command) error {
	switch after := after.(type) {
	case string:
		newCmd.After = after
	default:
		return newParseCommandError(
			"must be a string",
			newCmd.Name,
			AFTER,
			"",
		)
	}

	return nil
}
