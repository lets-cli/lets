package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/lets-cli/lets/commands/command"
	"github.com/lets-cli/lets/util"
	"github.com/lets-cli/lets/workdir"
)

var (
	// COMMANDS is a top-level directive. Includes all commands to run
	COMMANDS = "commands"
	SHELL    = "shell"
	ENV      = "env"
	EvalEnv  = "eval_env"
	MIXINS   = "mixins"
	VERSION  = "version"
	BEFORE   = "before"
)

const defaultConfigPath = "lets.yaml"

var validConfigFields = []string{COMMANDS, SHELL, ENV, EvalEnv, MIXINS, VERSION, BEFORE}
var validMixinConfigFields = []string{COMMANDS, ENV, EvalEnv, BEFORE}

var errFileNotExists = errors.New("file not exists")

type PathInfo struct {
	Filename string
	AbsPath  string
	WorkDir  string
}

// Config is a struct for loaded config file
type Config struct {
	// absolute path to work dir - where config is placed
	WorkDir  string
	FilePath string
	Commands map[string]command.Command
	Shell    string
	// before is a script which will be included before every cmd
	Before  string
	Env     map[string]string
	Version string
	isMixin bool // if true, we consider config as mixin and apply different parsing and validation
	// absolute path to .lets
	DotLetsDir string
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
	sep := "."

	if field == "" {
		sep = ""
	}

	fullPath := strings.Join(fields, sep)

	return &ParseError{
		Path: struct {
			Full  string
			Field string
		}{
			Full:  fullPath,
			Field: field,
		},
		Err: fmt.Errorf("field '%s': %s", fullPath, msg),
	}
}

func newConfig(workDir string, configAbsPath string) *Config {
	return &Config{
		Commands: make(map[string]command.Command),
		Env:      make(map[string]string),
		WorkDir:  workDir,
		FilePath: configAbsPath,
	}
}

func newMixinConfig(workDir string, configAbsPath string) *Config {
	cfg := newConfig(workDir, configAbsPath)
	cfg.isMixin = true

	return cfg
}

func Read(version string) (*Config, error) {
	configPath, err := FindConfig()
	if err != nil {
		return nil, err
	}

	cfg, err := LoadFromFile(configPath, version)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetDefaultConfigPath() string {
	return defaultConfigPath
}

// find config file recursively
// filename is a file to find and work dir is where to start.
func getFullConfigPathRecursive(filename string, workDir string) (string, error) {
	fileAbsPath, err := filepath.Abs(filepath.Join(workDir, filename))
	if err != nil {
		return "", err
	}

	if util.FileExists(fileAbsPath) {
		return fileAbsPath, nil
	}

	// else we get parent and try again up until we reach roof of fs
	parentDir := filepath.Dir(workDir)
	if parentDir == "/" {
		return "", fmt.Errorf("can not find config")
	}

	return getFullConfigPathRecursive(filename, parentDir)
}

// find config file non-recursively
// filename is a file to find and work dir is where to start.
func getFullConfigPath(filename string, workDir string) (string, error) {
	fileAbsPath, err := filepath.Abs(filepath.Join(workDir, filename))
	if err != nil {
		return "", err
	}

	if !util.FileExists(fileAbsPath) {
		return "", fmt.Errorf("%w: %s", errFileNotExists, fileAbsPath)
	}

	return fileAbsPath, nil
}

// workDir is where lets.yaml found or rootDir points to
func getWorkDir(filename string, rootDir string) (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get workdir for config %s: %s", filename, err)
	}

	if rootDir != "" {
		workDir = rootDir
	}

	return workDir, nil
}

// Load a config from file
func LoadFromFile(pathInfo PathInfo, letsVersion string) (*Config, error) {
	failedLoadErr := func(err error) error {
		return fmt.Errorf("failed to load config file %s: %s", pathInfo.Filename, err)
	}

	config := newConfig(pathInfo.WorkDir, pathInfo.AbsPath)

	err := loadConfig(pathInfo.AbsPath, config)
	if err != nil {
		return nil, failedLoadErr(err)
	}

	if err = Validate(config, letsVersion); err != nil {
		return nil, failedLoadErr(err)
	}

	dotLetsDir, err := workdir.GetDotLetsDir(pathInfo.WorkDir)
	if err != nil {
		return nil, err
	}

	config.DotLetsDir = dotLetsDir

	return config, nil
}

func loadConfig(filename string, cfg *Config) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileData, cfg)
	if err != nil {
		return err
	}

	return nil
}

func loadMixinConfig(filename string, rootCfg *Config) (*Config, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := newMixinConfig(rootCfg.WorkDir, filename)

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
		cmdsMap, ok := cmds.(map[interface{}]interface{})
		if !ok {
			return newConfigParseError(
				"must be a mapping",
				COMMANDS,
				"",
			)
		}

		if err := cfg.loadCommands(cmdsMap); err != nil {
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

	if evalEnv, ok := rawKeyValue[EvalEnv]; ok {
		evalEnv, ok := evalEnv.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("eval_env must be a mapping")
		}

		err := parseAndValidateEvalEnv(evalEnv, cfg)
		if err != nil {
			return err
		}
	}

	if before, ok := rawKeyValue[BEFORE]; ok {
		before, ok := before.(string)
		if !ok {
			return fmt.Errorf("before must be a string")
		}

		err := parseAndValidateBefore(before, cfg)
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

	if version, ok := rawKeyValue[VERSION]; ok {
		versionParseErr := fmt.Errorf("version must be a valid semver string")

		version, ok := version.(string)
		if !ok {
			return versionParseErr
		}

		_, err := util.ParseVersion(version)
		if err != nil {
			return versionParseErr
		}

		cfg.Version = version
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

// Trim `-` prefix.
// Using this prefix we allow to include non-existed mixins (git-ignored for example)
func normalizeMixinFilename(filename string) string {
	return strings.TrimPrefix(filename, "-")
}

// Ignored means that its okay of minix does not exists.
// It can be a git-ignored file for example
func isIgnoredMixin(filename string) bool {
	return strings.HasPrefix(filename, "-")
}

func readAndValidateMixins(mixins []interface{}, cfg *Config) error {
	for _, filename := range mixins {
		if filename, ok := filename.(string); ok {
			configAbsPath, err := getFullConfigPath(normalizeMixinFilename(filename), cfg.WorkDir)
			if err != nil {
				if isIgnoredMixin(filename) && errors.Is(err, errFileNotExists) {
					continue
				} else {
					// complain non-existed mixin only if its filename does not starts with dash `-`
					return fmt.Errorf("failed to read mixin config: %s", err)
				}
			}

			mixinCfg, err := loadMixinConfig(configAbsPath, cfg)
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

func joinBeforeScripts(beforeScripts ...string) string {
	buf := new(bytes.Buffer)

	for _, script := range beforeScripts {
		buf.WriteString(script)
		buf.WriteString("\n")
	}

	return buf.String()
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

	mainCfg.Before = joinBeforeScripts(
		mainCfg.Before,
		mixinCfg.Before,
	)

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

func parseAndValidateEvalEnv(evalEnv map[interface{}]interface{}, cfg *Config) error {
	for name, value := range evalEnv {
		nameKey := name.(string)

		if value, ok := value.(string); ok {
			computedVal, err := command.EvalEnvVariable(value)
			if err != nil {
				return err
			}

			cfg.Env[nameKey] = computedVal
		} else {
			return newConfigParseError(
				"must be a string",
				EvalEnv,
				nameKey,
			)
		}
	}

	return nil
}

func parseAndValidateBefore(before string, cfg *Config) error {
	cfg.Before = before

	return nil
}

func (c *Config) loadCommands(cmds map[interface{}]interface{}) error {
	for key, value := range cmds {
		keyStr, ok := key.(string)
		if !ok {
			return newConfigParseError(
				"command name must be a string",
				COMMANDS,
				"",
			)
		}

		newCmd := command.NewCommand(keyStr)
		err := command.ParseAndValidateCommand(&newCmd, value.(map[interface{}]interface{}))

		if err != nil {
			return err
		}

		c.Commands[keyStr] = newCmd
	}

	return nil
}
