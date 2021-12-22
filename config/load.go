package config

import (
	"fmt"
	"os"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/config/parser"
)


func Load(version string) (*config.Config, error) {
	configPath, err := FindConfig()
	if err != nil {
		return nil, err
	}

	cfg := config.NewConfig(
		configPath.WorkDir, 
		configPath.AbsPath,
		configPath.DotLetsDir,
	)

	fileData, err := os.ReadFile(configPath.AbsPath)
	if err != nil {
		return nil, fmt.Errorf("can not read config file: %w", err)
	}
	
	err = parser.Parse(fileData, cfg)
	if err != nil {
		return nil, err
	}

	
	if err = validate(cfg, version); err != nil {
		return nil, err
	}

	return cfg, nil
}