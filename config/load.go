package config

import (
	"fmt"
	"os"

	"github.com/lets-cli/lets/config/config"
	"gopkg.in/yaml.v3"
)

func Load(configName string, configDir string, version string) (*config.Config, error) {
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

	if err := c.SetupEnv(); err != nil {
		return nil, err
	}

	return c, nil
}
