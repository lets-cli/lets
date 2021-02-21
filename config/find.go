package config

import (
	"fmt"
	"path/filepath"

	"github.com/lets-cli/lets/env"
)

// FindConfig will try to find best match for config file.
// Rules are:
// - if specified LETS_CONFIG - try to load only that file
// - if specified LETS_CONFIG_DIR - try to look for a config only in that dir - don't do recursion
// - if not specified any of env vars above - try to find config recursively.
func FindConfig() (PathInfo, error) {
	configFilename, workDir := env.GetConfigPathFromEnv()
	configDirFromEnv := workDir != ""

	if configFilename == "" {
		configFilename = GetDefaultConfigPath()
	}

	failedFindErr := func(err error, filename string) error {
		return fmt.Errorf("failed to find config file %s: %w", filename, err)
	}

	// work dir is where to start looking for lets.yaml
	workDir, err := getWorkDir(configFilename, workDir)
	if err != nil {
		return PathInfo{}, err
	}

	configAbsPath := ""

	// if user specified full path to config file
	if filepath.IsAbs(configFilename) { //nolint:nestif
		configAbsPath = configFilename
	} else {
		if configDirFromEnv {
			configAbsPath, err = getFullConfigPath(configFilename, workDir)
			if err != nil {
				return PathInfo{}, failedFindErr(err, configFilename)
			}
		} else {
			// try to find abs config path up in parent dir tree
			configAbsPath, err = getFullConfigPathRecursive(configFilename, workDir)
			if err != nil {
				return PathInfo{}, failedFindErr(err, configFilename)
			}
		}
	}

	// just to be sure that work dir is correct
	workDir = filepath.Dir(configAbsPath)

	pathInfo := PathInfo{
		AbsPath:  configAbsPath,
		WorkDir:  workDir,
		Filename: configFilename,
	}

	return pathInfo, nil
}
