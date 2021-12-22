package parser

import (
	"fmt"

	"github.com/lets-cli/lets/config"
)

func parseAndValidateEnv(env interface{}, newCmd *config.Command) error {
	for name, value := range env.(map[interface{}]interface{}) {
		nameKey := name.(string)
		newCmd.Env[nameKey] = fmt.Sprintf("%v", value)
	}

	return nil
}

// TODO split parsers for command and config
func parseAndValidateEnvForConfig(env map[interface{}]interface{}, cfg *config.Config) error {
	for name, value := range env {
		nameKey := name.(string)

		if value, ok := value.(string); ok {
			cfg.Env[nameKey] = value
		} else {
			return newConfigParseError(
				"must be a string",
				ENV,
				nameKey,
			)
		}
	}

	return nil
}
