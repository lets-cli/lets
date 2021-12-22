package parser

import (
	"github.com/lets-cli/lets/config"
)

func parseAndValidatePersistChecksum(persistChecksum interface{}, newCmd *config.Command) error {
	shouldPersist, ok := persistChecksum.(bool)

	if !ok {
		return newParseCommandError(
			"must be a bool",
			newCmd.Name,
			PersistChecksum,
			"",
		)
	}

	if !newCmd.HasChecksum {
		return newParseCommandError(
			"you must declare 'checksum' for command to use 'persist_checksum'",
			newCmd.Name,
			PersistChecksum,
			"",
		)
	}

	newCmd.PersistChecksum = shouldPersist

	return nil
}
