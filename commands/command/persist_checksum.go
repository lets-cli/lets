package command

func parseAndValidatePersistChecksum(persistChecksum interface{}, newCmd *Command) error {
	shouldPersist, ok := persistChecksum.(bool)

	if !ok {
		return newParseCommandError(
			"must be a bool",
			newCmd.Name,
			PersistChecksum,
			"",
		)
	}

	newCmd.persistChecksum = shouldPersist

	return nil
}
