package config

import (
	"fmt"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/util"
)

// Validate loaded config.
func validate(config *config.Config, letsVersion string) error {
	if err := validateVersion(config, letsVersion); err != nil {
		return err
	}

	if err := validateDepends(config); err != nil {
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

func validateDepends(cfg *config.Config) error {
	for _, cmd := range cfg.Commands {
		cmd := cmd
		err := cmd.Depends.Range(func(key string, value config.Dep) error {
			dependency, exists := cfg.Commands[key]
			
			if !exists {
				return fmt.Errorf(
					"command '%s' depends on command '%s' which is not exist",
					cmd.Name, key,
				)
			}

			if dependency.Cmds.Parallel {
				return fmt.Errorf(
					"command '%s' depends on command '%s', but parallel cmd is not allowed in depends yet",
					cmd.Name, dependency.Name,
				)
			}

			return nil
		})

		if err != nil {
			return err
		}

		for _, other := range cfg.Commands {
			if cmd.Name == other.Name {
				continue
			}

			// if any two commands have each other command in deps, raise error.
			if cmd.Depends.Has(other.Name) && other.Depends.Has(cmd.Name) {
				return fmt.Errorf(
					"command '%s' have circular depends on command '%s'",
					cmd.Name, other.Name,
				)
			}
		}
	}

	return nil
}
