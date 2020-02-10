package command

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	CMD         = "cmd"
	DESCRIPTION = "description"
	ENV         = "env"
	EVAL_ENV    = "eval_env"
	OPTIONS     = "options"
	DEPENDS     = "depends"
	CHECKSUM    = "checksum"
)

type Command struct {
	Name        string
	Cmd         string
	Description string
	Env         map[string]string
	RawOptions  string
	Options     map[string]string
	Depends     []string
	Checksum    string
}

type CommandError struct {
	Path struct {
		Full  string
		Field string
	}
	Err error
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("failed to parse command: %s", e.Err)
}

// env is not proper arg
func newCommandError(msg string, name string, field string, env string) error {
	fields := []string{name, field}
	if env != "" {
		fields = append(fields, env)
	}
	fullPath := strings.Join(fields, ".")
	return &CommandError{
		Path: struct {
			Full  string
			Field string
		}{
			Full:  fullPath,
			Field: field,
		},
		Err: errors.New(fmt.Sprintf("field %s: %s", fullPath, msg)),
	}
}

// TODO interface{} must be replaced
func NewCommand(name string) Command {
	newCmd := Command{
		Name: name,
		Env:  make(map[string]string),
	}
	return newCmd
}

// TODO rename DeserializeCommand or make it command method
// add readAndValidate func, with cmd, section name and return validated and error

//TODO maybe add validate package and implement funcs such as read string, read list, or so
// or create file per section and create read and validate functions with some piblic api like Parse or ReadAndValidate
func DeserializeCommand(newCmd *Command, rawCommand map[interface{}]interface{}) error {
	if cmd, ok := rawCommand[CMD]; ok {
		// TODO decide, validate here or top-level validate and return all errors at once
		switch cmd := cmd.(type) {
		case string:
			newCmd.Cmd = cmd
		case []interface{}:
			cmdList := make([]string, len(cmd))
			for _, v := range cmd {
				cmdList = append(cmdList, v.(string))
			}
			cmdList = append(cmdList, os.Args[1:]...)
			newCmd.Cmd = strings.Join(cmdList, " ")
		default:
			return newCommandError(
				"must be either string or list of string",
				newCmd.Name,
				CMD,
				"",
			)
		}
	}

	if desc, ok := rawCommand[DESCRIPTION]; ok {
		if value, ok := desc.(string); ok {
			newCmd.Description = value
		} else {
			return newCommandError(
				"must be a string",
				newCmd.Name,
				DESCRIPTION,
				"",
			)
		}
	}

	if env, ok := rawCommand[ENV]; ok {
		// TODO dirty hacks
		for name, value := range env.(map[interface{}]interface{}) {
			nameKey := name.(string)
			if value, ok := value.(string); ok {
				newCmd.Env[nameKey] = value
			} else {
				return newCommandError(
					"must be a string",
					newCmd.Name,
					ENV,
					nameKey,
				)
			}
		}
	}

	if evalEnv, ok := rawCommand[EVAL_ENV]; ok {
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
	}

	if options, ok := rawCommand[OPTIONS]; ok {
		if value, ok := options.(string); ok {
			newCmd.RawOptions = value
		} else {
			return newCommandError(
				"must be a string",
				newCmd.Name,
				OPTIONS,
				"",
			)
		}
	}
	// TODO continue validation
	if depends, ok := rawCommand[DEPENDS]; ok {
		if depends, ok := depends.([]interface{}); ok {
			for _, value := range depends {
				if value, ok := value.(string); ok {
					// TODO validate if command is really exists - in validate
					newCmd.Depends = append(newCmd.Depends, value)
				} else {
					return newCommandError(
						"value of depends list must be a string",
						newCmd.Name,
						DEPENDS,
						"",
					)
				}
			}
		} else {
			return newCommandError(
				"must be a list of string (commands)",
				newCmd.Name,
				DEPENDS,
				"",
			)
		}
	}

	if checksum, ok := rawCommand[CHECKSUM]; ok {
		patterns, ok := checksum.([]interface{})
		if !ok {
			return newCommandError(
				"must be a list of string (files of glob patterns)",
				newCmd.Name,
				CHECKSUM,
				"",
			)
		}

		var files []string
		for _, value := range patterns {
			if value, ok := value.(string); ok {
				files = append(files, value)
			} else {
				return newCommandError(
					"value of checksum list must be a string",
					newCmd.Name,
					CHECKSUM,
					"",
				)
			}
		}
		checksum, err := calculateChecksum(files)
		if err == nil {
			newCmd.Checksum = checksum
		} else {
			return errors.New(fmt.Sprintf("failed to calculate checksum: %s", err))
		}
	}
	return nil
}
