package command

func parseAndValidateEnv(env interface{}, newCmd *Command) error {
	for name, value := range env.(map[interface{}]interface{}) {
		nameKey := name.(string)

		if value, ok := value.(string); ok {
			newCmd.Env[nameKey] = value
		} else {
			return newParseCommandError(
				"must be a string",
				newCmd.Name,
				ENV,
				nameKey,
			)
		}
	}

	return nil
}
