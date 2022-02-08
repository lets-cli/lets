package parser

import (
	"github.com/lets-cli/lets/checksum"
	"github.com/lets-cli/lets/config/config"
)

func parseChecksum(checksum interface{}, newCmd *config.Command) error {
	patternsList, okList := checksum.([]interface{})
	patternsMap, okMap := checksum.(map[string]interface{})

	switch {
	case okList:
		return parseChecksumList(patternsList, newCmd)
	case okMap:
		return parseChecksumMap(patternsMap, newCmd)
	default:
		return parseError(
			"must be a list of string (files of glob patterns) or a map of lists of string",
			newCmd.Name,
			CHECKSUM,
			"",
		)
	}
}

func parseChecksumList(patternsList []interface{}, newCmd *config.Command) error {
	checksumSource := make(map[string][]string)

	for _, value := range patternsList {
		if value, ok := value.(string); ok {
			checksumSource[checksum.DefaultChecksumKey] = append(checksumSource[checksum.DefaultChecksumKey], value)
		} else {
			return parseError(
				"value of checksum list must be a string",
				newCmd.Name,
				CHECKSUM,
				"",
			)
		}
	}

	newCmd.ChecksumSources = checksumSource
	newCmd.HasChecksum = true

	return nil
}

func parseChecksumMap(patternsMap map[string]interface{}, newCmd *config.Command) error {
	checksumSources := make(map[string][]string)

	for key, patterns := range patternsMap {
		patterns, ok := patterns.([]interface{})

		if !ok {
			return parseError(
				"value of checksum map must be a list",
				newCmd.Name,
				CHECKSUM,
				"",
			)
		}

		for _, value := range patterns {
			if value, ok := value.(string); ok {
				checksumSources[key] = append(checksumSources[key], value)
			} else {
				return parseError(
					"value of checksum list must be a string",
					newCmd.Name,
					CHECKSUM,
					"",
				)
			}
		}
	}

	newCmd.ChecksumSources = checksumSources
	newCmd.HasChecksum = true

	return nil
}
