package command

import (
	"fmt"
	"strings"

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
}

type Command struct {
	Name string
	// script to run
	Cmd string
	// map of named scripts to run in parallel
	CmdMap          map[string]string
	Description     string
	Env             map[string]string
	RawOptions      string
	Options         map[string]string
	CliOptions      map[string]string
	Depends         []string
	Checksum        string
	ChecksumMap     map[string]string
	PersistChecksum bool

	// prepared args - started from command name
	Args []string

	// run only specified commands from cmd map
	Only []string
	// run all but excluded commands from cmd map
	Exclude []string

	// if command has declared checksum
	hasChecksum    bool
	checksumSource map[string][]string
	// store loaded persisted checksums here
	persistedChecksums map[string]string
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

// env is not proper arg
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

// NewCommand creates new command struct
func NewCommand(name string) Command {
	return Command{
		Name: name,
		Env:  make(map[string]string),
	}
}

func (cmd *Command) ChecksumCalculator(workDir string) error {
	if len(cmd.checksumSource) == 0 {
		return nil
	}

	return calculateChecksumFromSource(workDir, cmd)
}

func (cmd *Command) GetPersistedChecksums() map[string]string {
	return cmd.persistedChecksums
}

// ParseAndValidateCommand parses and validates unmarshaled yaml
func ParseAndValidateCommand(newCmd *Command, rawCommand map[interface{}]interface{}) error {
	if err := validateCommandFields(rawCommand, validFields); err != nil {
		return err
	}

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
