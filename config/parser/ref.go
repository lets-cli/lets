package parser

import (
	"github.com/lets-cli/lets/config/config"
)

func parseRef(rawRef interface{}, newCmd *config.Command) error {
	ref, ok := rawRef.(string)
	if !ok {
		return parseError(
			"must be a string",
			newCmd.Name,
			REF,
			"",
		)
	}

	newCmd.Ref = ref

	return nil
}
