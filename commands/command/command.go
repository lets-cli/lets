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
		for _, value := range depends.([]interface{}) {
			// TODO validate if command is realy exists - in validate
			newCmd.Depends = append(newCmd.Depends, value.(string))
		}
	}

	if checksum, ok := rawCommand[CHECKSUM]; ok {
		if patterns, ok := checksum.([]interface{}); ok {
			var files []string
			for _, value := range patterns {
				// TODO validate if command is realy exists - in validate
				files = append(files, value.(string))
			}
			checksum, err := calculateChecksum(files)
			if err == nil {
				newCmd.Checksum = checksum
			} else {
				// TODO return error or caclulate checksum upper in the code
				fmt.Printf("error while checksum %s\n", err)
			}
		}
	}
	return nil
}
