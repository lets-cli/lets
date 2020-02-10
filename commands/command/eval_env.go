package command

import "os/exec"

func evalEnvVariable(rawCmd string) (string, error) {
	cmd := exec.Command("sh", "-c", rawCmd)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func parseAndValidateEvalEnv(evalEnv interface{}, newCmd *Command) error {
	for name, value := range evalEnv.(map[interface{}]interface{}) {
		nameKey := name.(string)
		if value, ok := value.(string); ok {
			if computedVal, err := evalEnvVariable(value); err != nil {
				return err
			} else {
				newCmd.Env[nameKey] = computedVal
			}
		} else {
			return newCommandError(
				"must be a string",
				newCmd.Name,
				EVAL_ENV,
				nameKey,
			)
		}
		if computedVal, err := evalEnvVariable(value.(string)); err != nil {
			// TODO we have to fail here and log error for user
		} else {
			newCmd.Env[name.(string)] = computedVal
		}
	}
	return nil
}
