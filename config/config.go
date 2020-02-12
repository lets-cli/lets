package config

import (
	"fmt"
	"github.com/kindritskyiMax/lets/commands/command"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	// COMMANDS is a top-level directive. Includes all commands to run
	COMMANDS = "commands"
	SHELL    = "shell"
	ENV      = "env"
	MIXINS   = "mixins"
)

var validConfigFields = strings.Join([]string{COMMANDS, SHELL, ENV, MIXINS}, " ")
var validMixinConfigFields = strings.Join([]string{COMMANDS, ENV}, " ")

// Config is a struct for loaded config file
type Config struct {
	WorkDir  string
	FilePath string
	Commands map[string]command.Command
	Shell    string
	Env      map[string]string
	isMixin  bool // if true, we consider config as mixin and apply different parsing and validation
}

type ParseError struct {
	Path struct {
		Full  string
		Field string
	}
	Err error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse config: %s", e.Err)
}

func newConfigParseError(msg string, name string, field string) error {
	fields := []string{name, field}
	fullPath := strings.Join(fields, ".")
	return &ParseError{
		Path: struct {
			Full  string
			Field string
		}{
			Full:  fullPath,
			Field: field,
		},
		Err: fmt.Errorf("field %s: %s", fullPath, msg),
	}
}

func newConfig() *Config {
	return &Config{
		Commands: make(map[string]command.Command),
		Env:      make(map[string]string),
	}
}

func newMixinConfig() *Config {
	cfg := newConfig()
	cfg.isMixin = true
	return cfg
}

// Load a config from file
func Load(filename string, rootDir string) (*Config, error) {
	failedLoadErr := func(err error) error {
		return fmt.Errorf("failed to load config file %s: %s", filename, err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		return nil, failedLoadErr(err)
	}
	if rootDir != "" {
		workDir = rootDir
	}
	absPath, err := filepath.Abs(filepath.Join(workDir, filename))
	if err != nil {
		return nil, failedLoadErr(err)
	}

	config, err := loadConfig(absPath)
	if err != nil {
		return nil, failedLoadErr(err)
	}

	config.WorkDir = filepath.Dir(absPath)
	config.FilePath = absPath

	if err = Validate(config); err != nil {
		return nil, failedLoadErr(err)
	}
	return config, nil
}

func loadConfig(filename string) (*Config, error) {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := newConfig()
	err = yaml.Unmarshal(fileData, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func loadMixinConfig(filename string) (*Config, error) {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := newMixinConfig()
	err = yaml.Unmarshal(fileData, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// UnmarshalYAML unmarshals a config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rawKeyValue := make(map[string]interface{})

	if err := unmarshal(&rawKeyValue); err != nil {
		return err
	}
	if c.isMixin {
		return unmarshalMixinConfig(rawKeyValue, c)
	}
	return unmarshalConfig(rawKeyValue, c)
}

func unmarshalConfigGeneral(rawKeyValue map[string]interface{}, cfg *Config) error {
	if cmds, ok := rawKeyValue[COMMANDS]; ok {
		if err := cfg.loadCommands(cmds.(map[interface{}]interface{})); err != nil {
			return err
		}
	}
	if env, ok := rawKeyValue[ENV]; ok {
		env, ok := env.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("env must be a mapping")
		}
		err := parseAndValidateEnv(env, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalConfig(rawKeyValue map[string]interface{}, cfg *Config) error {
	if err := validateTopLevelFields(rawKeyValue, validConfigFields); err != nil {
		return err
	}

	if err := unmarshalConfigGeneral(rawKeyValue, cfg); err != nil {
		return err
	}

	if shell, ok := rawKeyValue[SHELL]; ok {
		shell, ok := shell.(string)
		if !ok {
			return fmt.Errorf("shell must be a string")
		}
		cfg.Shell = shell
	} else {
		return fmt.Errorf("'shell' field is required")
	}

	if mixins, ok := rawKeyValue[MIXINS]; ok {
		mixins, ok := mixins.([]interface{})
		if !ok {
			return fmt.Errorf("mixins must be a list of string")
		}
		err := readAndValidateMixins(mixins, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalMixinConfig(rawKeyValue map[string]interface{}, cfg *Config) error {
	if err := validateTopLevelFields(rawKeyValue, validMixinConfigFields); err != nil {
		return err
	}
	return unmarshalConfigGeneral(rawKeyValue, cfg)
}

func readAndValidateMixins(mixins []interface{}, cfg *Config) error {
	for _, filename := range mixins {
		if filename, ok := filename.(string); ok {
			mixinCfg, err := loadMixinConfig(filename)
			if err != nil {
				return fmt.Errorf("failed to load mixin config: %s", err)
			}
			if err := mergeConfigs(cfg, mixinCfg); err != nil {
				return fmt.Errorf("failed to merge mixin config %s with main config: %s", filename, err)
			}
		} else {
			return newConfigParseError(
				"must be a string",
				MIXINS,
				"list item",
			)
		}
	}
	return nil
}

// Merge main and mixin configs. If there is a conflict - return error as we do not override values
// TODO add test
func mergeConfigs(mainCfg *Config, mixinCfg *Config) error {
	for _, mixinCmd := range mixinCfg.Commands {
		if _, conflict := mainCfg.Commands[mixinCmd.Name]; conflict {
			return fmt.Errorf("command %s from mixin is already declared in main config's commands", mixinCmd.Name)
		}
		mainCfg.Commands[mixinCmd.Name] = mixinCmd
	}
	for mixinEnvKey, mixinEnvVal := range mixinCfg.Env {
		if _, conflict := mainCfg.Env[mixinEnvKey]; conflict {
			return fmt.Errorf("env %s from mixin is already declared in main config's env", mixinEnvKey)
		}
		mainCfg.Env[mixinEnvKey] = mixinEnvVal
	}
	return nil
}

func parseAndValidateEnv(env map[interface{}]interface{}, cfg *Config) error {
	for name, value := range env {
		nameKey := name.(string)
		if value, ok := value.(string); ok {
			cfg.Env[nameKey] = value
		} else {
			return newConfigParseError(
				"must be a string",
				ENV,
				nameKey,
			)
		}
	}
	return nil
}

func (c *Config) loadCommands(cmds map[interface{}]interface{}) error {
	for key, value := range cmds {
		keyStr := key.(string)
		newCmd := command.NewCommand(keyStr)

		err := command.ParseAndValidateCommand(&newCmd, value.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
		c.Commands[keyStr] = newCmd
	}
	return nil
}
