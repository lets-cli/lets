package parser

import (
	"fmt"
	"strings"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/util"
)

var (
	CMD             = "cmd"
	DESCRIPTION     = "description"
	ENV             = "env"
	EvalEnv         = "eval_env"
	OPTIONS         = "options"
	DEPENDS         = "depends"
	CHECKSUM        = "checksum"
	PersistChecksum = "persist_checksum"
	AFTER           = "after"
)

var validFields = []string{
	CMD,
	DESCRIPTION,
	ENV,
	EvalEnv,
	OPTIONS,
	DEPENDS,
	CHECKSUM,
	PersistChecksum,
	AFTER,
}

type ParseCommandError struct {
	Name string
	Path struct {
		Full  string
		Field string
	}
	Err error
}

func (e *ParseCommandError) Error() string {
	return fmt.Sprintf("failed to parse '%s' command: %s", e.Name, e.Err)
}

// env is not proper arg.
func parseError(msg string, name string, field string, meta string) error {
	fields := []string{field}
	if meta != "" {
		fields = append(fields, meta)
	}

	fullPath := strings.Join(fields, ".")

	return &ParseCommandError{
		Name: name,
		Path: struct {
			Full  string
			Field string
		}{
			Full:  fullPath,
			Field: field,
		},
		Err: fmt.Errorf("field %s: %s", fullPath, msg),
	}
}

// parseCommand parses and validates unmarshaled yaml.
func parseCommand(newCmd *config.Command, rawCommand map[interface{}]interface{}) error { //nolint:cyclop
	if err := validateCommandFields(rawCommand, validFields); err != nil {
		return err
	}

	if cmd, ok := rawCommand[CMD]; ok {
		if err := parseCmd(cmd, newCmd); err != nil {
			return err
		}
	}

	if after, ok := rawCommand[AFTER]; ok {
		if err := parseAfter(after, newCmd); err != nil {
			return err
		}
	}

	if desc, ok := rawCommand[DESCRIPTION]; ok {
		if err := parseDescription(desc, newCmd); err != nil {
			return err
		}
	}

	if env, ok := rawCommand[ENV]; ok {
		if err := parseEnv(env, newCmd); err != nil {
			return err
		}
	}

	if evalEnv, ok := rawCommand[EvalEnv]; ok {
		if err := parseEvalEnv(evalEnv, newCmd); err != nil {
			return err
		}
	}

	if options, ok := rawCommand[OPTIONS]; ok {
		if err := parseOptions(options, newCmd); err != nil {
			return err
		}
	}

	if depends, ok := rawCommand[DEPENDS]; ok {
		if err := parseDepends(depends, newCmd); err != nil {
			return err
		}
	}

	if checksum, ok := rawCommand[CHECKSUM]; ok {
		if err := parseChecksum(checksum, newCmd); err != nil {
			return err
		}
	}

	if persistChecksum, ok := rawCommand[PersistChecksum]; ok {
		if err := parsePersistChecksum(persistChecksum, newCmd); err != nil {
			return err
		}
	}

	return nil
}

func validateCommandFields(rawKeyValue map[interface{}]interface{}, validFields []string) error {
	for key := range rawKeyValue {
		if !util.IsStringInList(key.(string), validFields) {
			return fmt.Errorf("unknown command field '%s'", key)
		}
	}

	return nil
}