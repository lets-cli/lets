package command

import "fmt"

func parseAndValidateEnv(env interface{}, newCmd *Command) error {
	for name, value := range env.(map[interface{}]interface{}) {
		nameKey := name.(string)
		newCmd.Env[nameKey] = fmt.Sprintf("%v", value)
	}

	return nil
}
