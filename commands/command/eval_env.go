package command

import (
	"os/exec"
	"strings"
)

// eval env value and trim result string
// TODO pass env from cfg.env - it will allow to use static env in eval_env
// TODO maybe use cfg.Shell instead of sh
func EvalEnvVariable(rawCmd string) (string, error) {
	cmd := exec.Command("sh", "-c", rawCmd)
	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	res := string(out)
	// TODO get rid of TrimSpace
	return strings.TrimSpace(res), nil
}

func parseAndValidateEvalEnv(evalEnv interface{}, newCmd *Command) error {
	for name, value := range evalEnv.(map[interface{}]interface{}) {
		nameKey := name.(string)

		if value, ok := value.(string); ok {
			computedVal, err := EvalEnvVariable(value)
			if err != nil {
				return err
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
