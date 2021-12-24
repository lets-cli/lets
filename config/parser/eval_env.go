package parser

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lets-cli/lets/config/config"
)

// eval env value and trim result string
// TODO pass env from cfg.env - it will allow to use static env in eval_env
// TODO maybe use cfg.Shell instead of sh.
func evalEnvVariable(rawCmd string) (string, error) {
	cmd := exec.Command("sh", "-c", rawCmd)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("can not get output from eval_env script: %s: %w", rawCmd, err)
	}

	res := string(out)
	// TODO get rid of TrimSpace
	return strings.TrimSpace(res), nil
}

func parseAndValidateEvalEnv(evalEnv interface{}, newCmd *config.Command) error {
	for name, value := range evalEnv.(map[interface{}]interface{}) {
		nameKey := name.(string)

		if value, ok := value.(string); ok {
			computedVal, err := evalEnvVariable(value)
			if err != nil {
				return newParseCommandError(
					fmt.Sprintf("failed to eval: %s", err),
					newCmd.Name,
					EvalEnv,
					nameKey,
				)
			}

			newCmd.Env[nameKey] = computedVal
		} else {
			return newParseCommandError(
				"must be a string",
				newCmd.Name,
				EvalEnv,
				nameKey,
			)
		}
	}

	return nil
}

func parseAndValidateEvalEnvForConfig(evalEnv map[interface{}]interface{}, cfg *config.Config) error {
	for name, value := range evalEnv {
		nameKey := name.(string)

		if value, ok := value.(string); ok {
			computedVal, err := evalEnvVariable(value)
			if err != nil {
				return fmt.Errorf("can not evaluate env variable: %w", err)
			}

			cfg.Env[nameKey] = computedVal
		} else {
			return newConfigParseError(
				"must be a string",
				EvalEnv,
				nameKey,
			)
		}
	}

	return nil
}
