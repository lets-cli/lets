package parser

import (
	"fmt"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/set"
	log "github.com/sirupsen/logrus"
)

var (
	CMD             = "cmd"
	DESCRIPTION     = "description"
	WORKDIR         = "work_dir"
	SHELL           = "shell"
	PLUGINS         = "plugins"
	ENV             = "env"
	EvalEnv         = "eval_env"
	OPTIONS         = "options"
	DEPENDS         = "depends"
	CHECKSUM        = "checksum"
	PersistChecksum = "persist_checksum"
	AFTER           = "after"
	REF             = "ref"
	ARGS            = "args"
)

var directives = set.NewSet[string](
	CMD,
	DESCRIPTION,
	WORKDIR,
	SHELL,
	PLUGINS,
	ENV,
	EvalEnv,
	OPTIONS,
	DEPENDS,
	CHECKSUM,
	PersistChecksum,
	AFTER,
	REF,
	ARGS,
)

// parseCommand parses and validates unmarshaled yaml.
func parseCommand(newCmd *config.Command, rawCommand map[string]interface{}, cfg *config.Config) error {
	if err := validateCommendDirectives(rawCommand); err != nil {
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

	if workdir, ok := rawCommand[WORKDIR]; ok {
		if err := parseWorkDir(workdir, newCmd); err != nil {
			return err
		}
	}

	if shell, ok := rawCommand[SHELL]; ok {
		if err := parseShell(shell, newCmd); err != nil {
			return err
		}
	}

	if rawPlugins, ok := rawCommand[PLUGINS]; ok {
		plugins, ok := rawPlugins.(map[string]interface{})
		if !ok {
			return fmt.Errorf("plugins must be a mapping")
		}
		if err := parsePlugins(plugins, newCmd); err != nil {
			return err
		}
	}

	rawEnv := make(map[string]interface{})

	if env, ok := rawCommand[ENV]; ok {
		env, ok := env.(map[string]interface{})
		if !ok {
			return fmt.Errorf("env must be a mapping")
		}
		for k, v := range env {
			rawEnv[k] = v
		}
	}

	if evalEnv, ok := rawCommand[EvalEnv]; ok {
		log.Debug("eval_env is deprecated, consider using 'env' with 'sh' executor")
		evalEnv, ok := evalEnv.(map[string]interface{})
		if !ok {
			return fmt.Errorf("env must be a mapping")
		}
		for k, v := range evalEnv {
			rawEnv[k] = map[string]interface{}{"sh": v}
		}
	}

	envEntries, err := parseEnvEntries(rawEnv, cfg)
	if err != nil {
		return err
	}

	for _, entry := range envEntries {
		value, err := entry.Value()
		if err != nil {
			return parseDirectiveError(
				"env",
				fmt.Sprintf("can not get value for '%s' env variable", entry.Name()),
			)
		}

		newCmd.Env[entry.Name()] = value
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

	ref, refOk := rawCommand[REF]
	if refOk {
		if err := parseRef(ref, newCmd); err != nil {
			return err
		}
	}

	if refOk {
		if args, ok := rawCommand[ARGS]; ok {
			if err := parseArgs(args, newCmd); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateCommendDirectives(rawKeyValue map[string]interface{}) error {
	for key := range rawKeyValue {
		if !directives.Contains(key) {
			return fmt.Errorf("unknown command field '%s'", key)
		}
	}

	_, argsExist := rawKeyValue[ARGS] // # ifshort
	_, refExist := rawKeyValue[REF]

	if argsExist && !refExist {
		return fmt.Errorf("'args' can only be used with 'ref'")
	}

	return nil
}
