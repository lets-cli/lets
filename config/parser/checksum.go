package parser

import "github.com/lets-cli/lets/config/config"

func parseChecksum(checksum interface{}, newCmd *config.Command) error { //nolint:cyclop
	patternsList, okList := checksum.([]interface{})
	patternsMap, okMap := checksum.(map[interface{}]interface{})
	checksumSource := make(map[string][]string)

	switch {
	case okList:
		for _, value := range patternsList {
			if value, ok := value.(string); ok {
				checksumSource[""] = append(checksumSource[""], value)
			} else {
				return parseError(
					"value of checksum list must be a string",
					newCmd.Name,
					CHECKSUM,
					"",
				)
			}
		}
	case okMap:
		for key, patterns := range patternsMap {
			key, ok := key.(string)
			if !ok {
				return parseError(
					"key of checksum list must be a string",
					newCmd.Name,
					CHECKSUM,
					"",
				)
			}

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
					checksumSource[key] = append(checksumSource[key], value)
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
	default:
		return parseError(
			"must be a list of string (files of glob patterns) or a map of lists of string",
			newCmd.Name,
			CHECKSUM,
			"",
		)
	}

	newCmd.ChecksumSource = checksumSource
	newCmd.HasChecksum = true

	return nil
}
