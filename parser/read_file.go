package parser

import (
	"github.com/lets-cli/lets/config"
)

// TODO not sure where reader must be ?
func ReadConfig(version string) (*config.Config, error) {
	configPath, err := config.FindConfig()
	if err != nil {
		return nil, err
	}

	cfg, err := LoadFromFile(configPath, version)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
