package config

import (
	"fmt"
	"path/filepath"

	"github.com/lets-cli/lets/set"
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

var (
	ValidConfigDirectives = set.NewSet(
		COMMANDS, SHELL, ENV, EvalEnv, MIXINS, VERSION, BEFORE,
	)
	ValidMixinConfigDirectives = set.NewSet(
		COMMANDS, ENV, EvalEnv, BEFORE,
	)
)

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
	// absolute path to .lets/checksums
	ChecksumsDir string
	// absolute path to .lets/mixins
	MixinsDir string
}

func NewConfig(workDir string, configAbsPath string, dotLetsDir string) *Config {
	return &Config{
		Commands:     make(map[string]Command),
		Env:          make(map[string]string),
		WorkDir:      workDir,
		FilePath:     configAbsPath,
		DotLetsDir:   dotLetsDir,
		ChecksumsDir: filepath.Join(dotLetsDir, "checksums"),
		MixinsDir:    filepath.Join(dotLetsDir, "mixins"),
	}
}

func NewMixinConfig(cfg *Config, configAbsPath string) *Config {
	mixin := NewConfig(cfg.WorkDir, configAbsPath, cfg.DotLetsDir)
	mixin.isMixin = true

	return mixin
}

func (c *Config) CreateChecksumsDir() error {
	if err := util.SafeCreateDir(c.ChecksumsDir); err != nil {
		return fmt.Errorf("can not create %s: %w", c.ChecksumsDir, err)
	}

	return nil
}

func (c *Config) CreateMixinsDir() error {
	if err := util.SafeCreateDir(c.MixinsDir); err != nil {
		return fmt.Errorf("can not create %s: %w", c.MixinsDir, err)
	}

	return nil
}
