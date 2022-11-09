package config

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lets-cli/lets/config/path"
	"github.com/lets-cli/lets/util"
	"github.com/lets-cli/lets/workdir"
	log "github.com/sirupsen/logrus"
)

const defaultConfigFile = "lets.yaml"

type PathInfo struct {
	Filename string
	AbsPath  string
	WorkDir  string
	// .lets abs path
	DotLetsDir string
}

// FindConfig will try to find best match for config file.
// Rules are:
// - if specified configName - try to load only that file
// - if specified configDir - try to look for a config only in that dir - don't do recursion
// - if not specified any of params above - try to find config recursively.
func FindConfig(configName string, configDir string) (PathInfo, error) {
	configDirSpecifiedByUser := configDir != ""

	if configName == "" {
		configName = defaultConfigFile
	}

	// work dir is where to start looking for lets.yaml
	workDir, err := getWorkDir(configName, configDir)
	if err != nil {
		return PathInfo{}, err
	}

	log.Debugf(color.BlueString("lets: found %s config file in %s directory", configName, workDir))

	configAbsPath := ""

	// if user specified full path to config file
	if filepath.IsAbs(configName) { //nolint:nestif
		configAbsPath = configName
	} else {
		if configDirSpecifiedByUser {
			configAbsPath, err = path.GetFullConfigPath(configName, workDir)
			if err != nil {
				return PathInfo{}, err
			}
		} else {
			// try to find abs config path up in parent dir tree
			configAbsPath, err = path.GetFullConfigPathRecursive(configName, workDir)
			if err != nil {
				return PathInfo{}, err
			}
		}
	}

	// just to be sure that work dir is correct
	workDir = filepath.Dir(configAbsPath)

	dotLetsDir, err := workdir.GetDotLetsDir(workDir)
	if err != nil {
		return PathInfo{}, fmt.Errorf("can not get .lets absolute path: %w", err)
	}

	if err := util.SafeCreateDir(dotLetsDir); err != nil {
		return PathInfo{}, fmt.Errorf("can not create .lets dir: %w", err)
	}

	pathInfo := PathInfo{
		AbsPath:    configAbsPath,
		WorkDir:    workDir,
		Filename:   configName,
		DotLetsDir: dotLetsDir,
	}

	return pathInfo, nil
}
