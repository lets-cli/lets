package config

import (
	"fmt"
	"path/filepath"

	"github.com/lets-cli/lets/env"
)

func FindConfig() (ConfigPath, error) {
	configFilename, workDir := env.GetConfigPathFromEnv()

	if configFilename == "" {
		configFilename = GetDefaultConfigPath()
	}

	failedFindErr := func(err error) error {
		return fmt.Errorf("failed to find config file %s: %s", configFilename, err)
	}

	// work dir is where to start looking for lets.yaml
	workDir, err := getWorkDir(configFilename, workDir)
	if err != nil {
		return ConfigPath{}, err
	}

	configAbsPath := ""

	// if user specified full path to config file
	if filepath.IsAbs(configFilename) {
		configAbsPath = configFilename
	} else {
		// try to find abs config path up in parent dir tree
		configAbsPath, err = getFullConfigPath(configFilename, workDir)
		if err != nil {
			return ConfigPath{}, failedFindErr(err)
		}
	}

	// just to be sure that work dir is correct
	workDir = filepath.Dir(configAbsPath)

	cp := ConfigPath{
		AbsPath: configAbsPath,
		WorkDir: workDir,
		Filename: configFilename,
	}

	return cp, nil
}
