package command

import (
	"errors"
	"fmt"
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

// NewCommand creates new command struct
func NewCommand(name string) Command {
	newCmd := Command{
		Name: name,
		Env:  make(map[string]string),
	}
	return newCmd
}

// ParseAndValidateCommand parses and validates unmarshaled yaml
func ParseAndValidateCommand(newCmd *Command, rawCommand map[interface{}]interface{}) error {
	if cmd, ok := rawCommand[CMD]; ok {
		if err := parseAndValidateCmd(cmd, newCmd); err != nil {
			return err
		}
	}

	if desc, ok := rawCommand[DESCRIPTION]; ok {
		if err := parseAndValidateDescription(desc, newCmd); err != nil {
			return err
		}
	}

	if env, ok := rawCommand[ENV]; ok {
		if err := parseAndValidateEnv(env, newCmd); err != nil {
			return err
		}
	}

	if evalEnv, ok := rawCommand[EVAL_ENV]; ok {
		if err := parseAndValidateEvalEnv(evalEnv, newCmd); err != nil {
			return err
		}
	}

	if options, ok := rawCommand[OPTIONS]; ok {
		if err := parseAndValidateOptions(options, newCmd); err != nil {
			return err
		}
	}

	if depends, ok := rawCommand[DEPENDS]; ok {
		if err := parseAndValidateDepends(depends, newCmd); err != nil {
			return err
		}
	}

	if checksum, ok := rawCommand[CHECKSUM]; ok {
		if err := parseAndValidateChecksum(checksum, newCmd); err != nil {
			return err
		}
	}
	return nil
}
