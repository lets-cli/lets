package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lets-cli/lets/util"
)

var (
	// COMMANDS is a top-level directive. Includes all commands to run.
	COMMANDS = "commands"
	SHELL    = "shell"
	ENV      = "env"
	EvalEnv  = "eval_env"
	MIXINS   = "mixins"
	VERSION  = "version"
	BEFORE   = "before"
)

const defaultConfigPath = "lets.yaml"

var (
	ValidConfigFields      = []string{COMMANDS, SHELL, ENV, EvalEnv, MIXINS, VERSION, BEFORE}
	ValidMixinConfigFields = []string{COMMANDS, ENV, EvalEnv, BEFORE}
)

var (
	ErrFileNotExists  = errors.New("file not exists")
	errConfigNotFound = errors.New("can not find config")
)

type PathInfo struct {
	Filename string
	AbsPath  string
	WorkDir  string
}

// Config is a struct for loaded config file.
type Config struct {
	// absolute path to work dir - where config is placed
	WorkDir  string
	FilePath string
	Commands map[string]Command
	Shell    string
	// before is a script which will be included before every cmd
	Before  string
	Env     map[string]string
	Version string
	isMixin bool // if true, we consider config as mixin and apply different parsing and validation
	// absolute path to .lets
	DotLetsDir string
}

func NewConfig(workDir string, configAbsPath string) *Config {
	return &Config{
		Commands: make(map[string]Command),
		Env:      make(map[string]string),
		WorkDir:  workDir,
		FilePath: configAbsPath,
	}
}

func NewMixinConfig(workDir string, configAbsPath string) *Config {
	cfg := NewConfig(workDir, configAbsPath)
	cfg.isMixin = true

	return cfg
}

func GetDefaultConfigPath() string {
	return defaultConfigPath
}

// find config file recursively
// filename is a file to find and work dir is where to start.
func getFullConfigPathRecursive(filename string, workDir string) (string, error) {
	fileAbsPath, err := filepath.Abs(filepath.Join(workDir, filename))
	if err != nil {
		return "", fmt.Errorf("can not get absolute workdir path: %w", err)
	}

	if util.FileExists(fileAbsPath) {
		return fileAbsPath, nil
	}

	// else we get parent and try again up until we reach roof of fs
	parentDir := filepath.Dir(workDir)
	if parentDir == "/" {
		return "", errConfigNotFound
	}

	return getFullConfigPathRecursive(filename, parentDir)
}

// find config file non-recursively
// filename is a file to find and work dir is where to start.
func GetFullConfigPath(filename string, workDir string) (string, error) {
	fileAbsPath, err := filepath.Abs(filepath.Join(workDir, filename))
	if err != nil {
		return "", fmt.Errorf("can not get absolute workdir path: %w", err)
	}

	if !util.FileExists(fileAbsPath) {
		return "", fmt.Errorf("%w: %s", ErrFileNotExists, fileAbsPath)
	}

	return fileAbsPath, nil
}

// workDir is where lets.yaml found or rootDir points to.
func getWorkDir(filename string, rootDir string) (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get workdir for config %s: %w", filename, err)
	}

	if rootDir != "" {
		workDir = rootDir
	}

	return workDir, nil
}
