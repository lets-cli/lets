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
	ENV    = "env"
)

var validFields = strings.Join([]string{COMMANDS, SHELL, ENV}, " ")

// Config is a struct for loaded config file
type Config struct {
	WorkDir  string
	FilePath string
	Commands map[string]command.Command
	Shell    string
	Env map[string]string
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

func newConfig() *Config {
	return &Config{
		Commands: make(map[string]command.Command),
		Env: make(map[string]string),
	}
}

// UnmarshalYAML unmarshals a config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rawKeyValue := make(map[string]interface{})

	if err := unmarshal(&rawKeyValue); err != nil {
		return err
	}

	if err := validateTopLevelFields(rawKeyValue, validFields); err != nil {
		return err
	}
	if cmds, ok := rawKeyValue[COMMANDS]; ok {
		if err := c.loadCommands(cmds.(map[interface{}]interface{})); err != nil {
			return err
		}
	}

	if shell, ok := rawKeyValue[SHELL]; ok {
		shell, ok := shell.(string)
		if !ok {
			return fmt.Errorf("shell must be a string")
		}
		c.Shell = shell
	} else {
		return fmt.Errorf("shell must be specified in config")
	}

	if env, ok := rawKeyValue[ENV]; ok {
		env, ok := env.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("env must be a mapping")
		}
		err := parseAndValidateEnv(env, c)
		if err != nil {
			return err
		}
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
