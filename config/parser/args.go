package parser

import (
	"github.com/lets-cli/lets/config/config"
)

func parseArgs(rawArgs interface{}, newCmd *config.Command) error {
	args, ok := rawArgs.(string)
	if !ok {
		return parseError(
			"must be a string",
			newCmd.Name,
			ARGS,
			"",
		)
	}

	newCmd.RefArgs = args

	return nil
}
