package parser

import (
	"fmt"

	"github.com/lets-cli/lets/config/config"
)

const (
	nameKey = "name"
	argsKey = "args"
	envKey  = "env"
)

var (
	depKeys    = []string{nameKey, argsKey, envKey}
	depKeysMap = map[string]bool{
		nameKey: true,
		argsKey: true,
		envKey:  true,
	}
)

func parseDependsAsMap(dep map[interface{}]interface{}, cmdName string, idx int) (*config.Dep, error) { //nolint:cyclop
	name := ""
	args := []string{}
	env := map[string]string{}

	for k, v := range dep {
		key, ok := k.(string)
		if !ok {
			return nil, parseError(
				"key of depend must be a string",
				cmdName,
				DEPENDS,
				"",
			)
		}

		if _, exists := depKeysMap[key]; !exists {
			return nil, parseError(
				fmt.Sprintf("key of depend must be one of %s", depKeys),
				cmdName,
				DEPENDS,
				"",
			)
		}

		switch key {
		case nameKey:
			value, ok := v.(string)
			if !ok {
				return nil, &ParseCommandError{
					Name: cmdName,
					Err: fmt.Errorf(
						"field '%s': %s",
						fmt.Sprintf("%s.[%d][name:%s]", DEPENDS, idx, name),
						"value of 'name' must be a string (an existing command)",
					),
				}
			}
			name = value
		case argsKey:
			switch value := v.(type) {
			case string:
				args = append(args, value)
			case []interface{}:
				for _, arg := range value {
					arg, ok := arg.(string)
					if !ok {
						return nil, &ParseCommandError{
							Name: cmdName,
							Err: fmt.Errorf(
								"field '%s': %s",
								fmt.Sprintf("%s.[%d][name:%s]", DEPENDS, idx, name),
								fmt.Sprintf("value of 'args' must be an array of string, got array element: %#v", arg)),
						}
					}
					args = append(args, arg)
				}
			default:
				return nil, &ParseCommandError{
					Name: cmdName,
					Err: fmt.Errorf(
						"field '%s': %s",
						fmt.Sprintf("%s.[%d][name:%s]", DEPENDS, idx, name),
						fmt.Sprintf("value of 'args' must be a string or an array of string, got: %#v", v)),
				}
			}
		case envKey:
			for envName, envValue := range v.(map[interface{}]interface{}) {
				nameKey := envName.(string)
				env[nameKey] = fmt.Sprintf("%v", envValue)
			}
		}
	}

	return &config.Dep{
		Name: name,
		// args always must start with a dependency name, otherwise docopt will fail
		Args: append([]string{name}, args...),
		Env:  env,
	}, nil
}

func parseDepends(rawDepends interface{}, newCmd *config.Command) error {
	depends, ok := rawDepends.([]interface{})

	if !ok {
		return parseError(
			"must be a list of string (commands) or a list of maps",
			newCmd.Name,
			DEPENDS,
			"",
		)
	}

	dependencies := make(map[string]config.Dep, len(depends))
	dependsNames := make([]string, 0, len(depends))

	for idx, value := range depends {
		switch v := value.(type) {
		case string:
			dep := &config.Dep{Name: v, Args: []string{}}
			dependencies[dep.Name] = *dep
			dependsNames = append(dependsNames, dep.Name)
		case map[interface{}]interface{}:
			dep, err := parseDependsAsMap(v, newCmd.Name, idx)
			if err != nil {
				return err
			}
			dependencies[dep.Name] = *dep
			dependsNames = append(dependsNames, dep.Name)
		default:
			return parseError(
				"value of depends list must be a string",
				newCmd.Name,
				DEPENDS,
				"",
			)
		}
	}

	newCmd.Depends = dependencies
	newCmd.DependsNames = dependsNames

	return nil
}
