package parser

import (
	"fmt"
	"os"

	"github.com/kballard/go-shellquote"
	"github.com/lets-cli/lets/config/config"
)

func parseArgs(rawArgs interface{}, newCmd *config.Command) error {
	switch args := rawArgs.(type) {
	case string:
		argsList, err := shellquote.Split(args)
		if err != nil {
			return parseError(
				"can not parse into args list",
				newCmd.Name,
				ARGS,
				err.Error(),
			)
		}

		newCmd.RefArgs = argsList
	case []string:
		newCmd.RefArgs = args
	case []interface{}:
		for _, arg := range args {
			newCmd.RefArgs = append(newCmd.RefArgs, fmt.Sprintf("%s", arg))
		}
	default:
		return parseError(
			"must be a string or a list of string",
			newCmd.Name,
			ARGS,
			"",
		)
	}

	return nil
}

func postprocessRefArgs(cfg *config.Config) {
	for _, cmd := range cfg.Commands {
		for idx, arg := range cmd.RefArgs {
			// we have to expand env here on our own, since this args not came from users tty, and not expanded before lets
			cmd.RefArgs[idx] = os.Expand(arg, func(key string) string {
				return cfg.Env[key]
			})
		}
	}
}
