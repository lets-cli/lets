package config

import (
	"fmt"

	"github.com/lets-cli/lets/commands/command"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/util"
)

const (
	NoticeColor = "\033[1;36m%s\033[0m"
)

func withColor(msg string) string {
	if env.IsNotColorOutput() {
		return msg
	}

	return fmt.Sprintf(NoticeColor, msg)
}

// Validate loaded config.
func Validate(config *Config, letsVersion string) error {
	if err := validateCircularDepends(config); err != nil {
		return err
	}

	return validateVersion(config, letsVersion)
}

func validateVersion(cfg *Config, letsVersion string) error {
	// no version specified on config
	if cfg.Version == "" {
		return nil
	}

	cfgVersionParsed, err := util.ParseVersion(cfg.Version)
	if err != nil {
		return fmt.Errorf("failed to parse config version: %w", err)
	}

	letsVersionParsed, err := util.ParseVersion(letsVersion)
	if err != nil {
		return fmt.Errorf("failed to parse lets version: %w", err)
	}

	// in dev (where version is 0.0.0-dev) this predicate will be always false
	if letsVersionParsed.LessThan(*cfgVersionParsed) {
		return fmt.Errorf(
			"config version '%s' is not compatible with 'lets' version '%s'. "+
				"Please upgrade 'lets' to '%s' "+
				"using 'lets --upgrade' command or following documentation at https://lets-cli.org/docs/installation'",
			cfgVersionParsed,
			letsVersionParsed,
			cfgVersionParsed,
		)
	}

	return nil
}

// if any two commands have each other command in deps, raise error
// TODO not optimal function but works for now.
func validateCircularDepends(cfg *Config) error {
	for _, cmdA := range cfg.Commands {
		for _, cmdB := range cfg.Commands {
			if cmdA.Name == cmdB.Name {
				continue
			}

			if yes := depsIntersect(cmdA, cmdB); yes {
				return fmt.Errorf(
					"command '%s' have circular depends on command '%s'",
					withColor(cmdA.Name),
					withColor(cmdB.Name),
				)
			}
		}
	}

	return nil
}

func depsIntersect(cmdA command.Command, cmdB command.Command) bool {
	return util.IsStringInList(cmdA.Name, cmdB.Depends) && util.IsStringInList(cmdB.Name, cmdA.Depends)
}

func validateTopLevelFields(rawKeyValue map[string]interface{}, validFields []string) error {
	for key := range rawKeyValue {
		if !util.IsStringInList(key, validFields) {
			return fmt.Errorf("unknown top-level field '%s'", key)
		}
	}

	return nil
}
