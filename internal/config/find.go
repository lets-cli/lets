package config

import (
	"fmt"
	"path/filepath"

	"github.com/lets-cli/lets/internal/config/path"
	"github.com/lets-cli/lets/internal/util"
	"github.com/lets-cli/lets/internal/workdir"
	log "github.com/sirupsen/logrus"
)

const defaultConfigFile = "lets.yaml"

type PathInfo struct {
	Filename string
	AbsPath  string
	// ConfigDir is the directory holding the config file. Config assembly
	// (mixin paths) resolves against it.
	ConfigDir string
	// RootDir is where commands run by default. Everything a command reads or
	// runs resolves against it, unless the command sets work_dir.
	RootDir string
	// .lets abs path
	DotLetsDir string
}

// FindConfig will try to find best match for config file.
// Rules are:
// - if specified configName - try to load only that file
// - if specified configDir - try to look for a config only in that dir - don't do recursion
// - if not specified any of params above - try to find config recursively.
func FindConfig(configName string, configDirFlag string) (PathInfo, error) {
	configDirSpecifiedByUser := configDirFlag != ""

	if configName == "" {
		configName = defaultConfigFile
	}

	// searchDir is where to start looking for lets.yaml
	searchDir, err := getSearchDir(configName, configDirFlag)
	if err != nil {
		return PathInfo{}, err
	}

	log.Debugf("found %s config file in %s directory", configName, searchDir)

	configAbsPath := ""

	// if user specified full path to config file
	if filepath.IsAbs(configName) { //nolint:nestif
		configAbsPath = configName
	} else {
		if configDirSpecifiedByUser {
			configAbsPath, err = path.GetFullConfigPath(configName, searchDir)
			if err != nil {
				return PathInfo{}, err
			}
		} else {
			// try to find abs config path up in parent dir tree
			configAbsPath, err = path.GetFullConfigPathRecursive(configName, searchDir)
			if err != nil {
				return PathInfo{}, err
			}
		}
	}

	configDir := filepath.Dir(configAbsPath)

	dotLetsDir, err := workdir.GetDotLetsDir(configDir)
	if err != nil {
		return PathInfo{}, fmt.Errorf("can not get .lets absolute path: %w", err)
	}

	if err := util.SafeCreateDir(dotLetsDir); err != nil {
		return PathInfo{}, fmt.Errorf("can not create .lets dir: %w", err)
	}

	pathInfo := PathInfo{
		AbsPath:   configAbsPath,
		ConfigDir: configDir,
		// preserved as-is here; the root is decoupled from the config dir in a follow-up
		RootDir:    configDir,
		Filename:   configName,
		DotLetsDir: dotLetsDir,
	}

	return pathInfo, nil
}
