package parser

import (
	"github.com/lets-cli/lets/config/config"
)

func parsePersistChecksum(persistChecksum interface{}, newCmd *config.Command) error {
	shouldPersist, ok := persistChecksum.(bool)

	if !ok {
		return parseError(
			"must be a bool",
			newCmd.Name,
			PersistChecksum,
			"",
		)
	}

	if !newCmd.HasChecksum {
		return parseError(
			"you must declare 'checksum' for command to use 'persist_checksum'",
			newCmd.Name,
			PersistChecksum,
			"",
		)
	}

	newCmd.PersistChecksum = shouldPersist

	return nil
}
