package config

import (
	"fmt"
	"os"

	"github.com/lets-cli/lets/config/config"
	"gopkg.in/yaml.v3"
)

// TODO 1. find a better place for function
// TODO 2. find a better place to call this func and ensure oredr of call is fine
func postprocessRefArgs(cfg *config.Config) {
	for _, cmd := range cfg.Commands {
		if cmd.Ref == nil {
			continue
		}

		for idx, arg := range cmd.Ref.Args {
			// we have to expand env here on our own, since this args not came from users tty, and not expanded before lets
			cmd.Ref.Args[idx] = os.Expand(arg, func(key string) string {
				return cfg.Env.Mapping[key].Value
			})
		}
	}
}


func Load(configName string, configDir string, version string ) (*config.Config, error) {
	configPath, err := FindConfig(configName, configDir)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(configPath.AbsPath)
	if err != nil {
		return nil, err
	}

	c := config.NewConfig(
		configPath.WorkDir,
		configPath.AbsPath,
		configPath.DotLetsDir,
	)
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, fmt.Errorf("lets: failed to parse %s: %w", configPath.Filename, err)
	}

	if err = validate(c, version); err != nil {
		return nil, err
	}

	postprocessRefArgs(c)

	return c, nil
}
