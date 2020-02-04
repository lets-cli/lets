package config

import (
	"fmt"
	"github.com/kindritskyiMax/lets/commands/command"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// COMMANDS is a top-level directive. Includes all commands to run
	COMMANDS = "commands"
)

// Config is a struct for loaded config file
type Config struct {
	WorkDir  string
	FilePath string
	Commands map[string]command.Command
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
	}
}

// UnmarshalYAML unmarshals a config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rawKeyValue := make(map[string]map[string]map[string]interface{})

	if err := unmarshal(&rawKeyValue); err != nil {
		return err
	}

	if cmds, ok := rawKeyValue[COMMANDS]; ok {
		if err := c.loadCommands(cmds); err != nil {
			return err
		}
	}
	return nil
}

// TODO refactor
func (c *Config) loadCommands(cmds map[string]map[string]interface{}) error {
	for key, value := range cmds {
		c.Commands[key] = command.NewCommand(key, value)
	}
	return nil
}
