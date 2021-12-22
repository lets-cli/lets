package parser

import (
	"fmt"
	"strings"

	"github.com/lets-cli/lets/config"
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
func newParseCommandError(msg string, name string, field string, meta string) error {
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

// parseAndValidateCommand parses and validates unmarshaled yaml.
func parseAndValidateCommand(newCmd *config.Command, rawCommand map[interface{}]interface{}) error { //nolint:cyclop
	if err := validateCommandFields(rawCommand, validFields); err != nil {
		return err
	}

	if cmd, ok := rawCommand[CMD]; ok {
		if err := parseAndValidateCmd(cmd, newCmd); err != nil {
			return err
		}
	}

	if after, ok := rawCommand[AFTER]; ok {
		if err := parseAndValidateAfter(after, newCmd); err != nil {
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

	if evalEnv, ok := rawCommand[EvalEnv]; ok {
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

	if persistChecksum, ok := rawCommand[PersistChecksum]; ok {
		if err := parseAndValidatePersistChecksum(persistChecksum, newCmd); err != nil {
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
