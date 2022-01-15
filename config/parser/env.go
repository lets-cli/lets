package parser

import (
	"fmt"

	"github.com/lets-cli/lets/config/config"
)

func parseEnv(env interface{}, newCmd *config.Command) error {
	for name, value := range env.(map[string]interface{}) {
		newCmd.Env[name] = fmt.Sprintf("%v", value)
	}

	return nil
}

func parseEnvForConfig(env map[string]interface{}, cfg *config.Config) error {
	for name, value := range env {
		if value, ok := value.(string); ok {
			cfg.Env[name] = value
		} else {
			return newConfigParseError(
				"must be a string",
				ENV,
				name,
			)
		}
	}

	return nil
}
