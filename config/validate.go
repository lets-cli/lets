package config

import (
	"fmt"

	"github.com/lets-cli/lets/config/config"
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
func validate(config *config.Config, letsVersion string) error {
	if err := validateVersion(config, letsVersion); err != nil {
		return err
	}

	if err := validateCommandInDependsExists(config); err != nil {
		return err
	}

	if err := validateDependsCycle(config); err != nil {
		return err
	}

	return nil
}

func validateVersion(cfg *config.Config, letsVersion string) error {
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

	isDev := letsVersionParsed.PreRelease == "dev"
	if letsVersionParsed.LessThan(*cfgVersionParsed) && !isDev {
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

func validateCommandInDependsExists(cfg *config.Config) error {
	for _, cmd := range cfg.Commands {
		err := cmd.Depends.Range(func(key string, value config.Dep) error {
			if _, exists := cfg.Commands[key]; !exists {
				return fmt.Errorf(
					"command '%s' depends on command '%s' which is not exist",
					withColor(cmd.Name),
					withColor(key),
				)
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// if any two commands have each other command in deps, raise error.
func validateDependsCycle(cfg *config.Config) error {
	for _, cmdA := range cfg.Commands {
		for _, cmdB := range cfg.Commands {
			if cmdA.Name == cmdB.Name {
				continue
			}

			if cmdA.Depends.Has(cmdB.Name) && cmdB.Depends.Has(cmdA.Name) {
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
