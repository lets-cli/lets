package parser

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/config/path"
	"github.com/lets-cli/lets/util"
	"gopkg.in/yaml.v2"
)

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

func parseConfigGeneral(rawKeyValue map[string]interface{}, cfg *config.Config) error { //nolint:cyclop
	if cmds, ok := rawKeyValue[config.COMMANDS]; ok {
		cmdsMap, ok := cmds.(map[interface{}]interface{})
		if !ok {
			return newConfigParseError(
				"must be a mapping",
				config.COMMANDS,
				"",
			)
		}

		commands, err := parseCommands(cmdsMap)
		if err != nil {
			return err
		}

		for _, c := range commands {
			cfg.Commands[c.Name] = c
		}
	}

	if env, ok := rawKeyValue[ENV]; ok {
		env, ok := env.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("env must be a mapping")
		}

		err := parseAndValidateEnvForConfig(env, cfg)
		if err != nil {
			return err
		}
	}

	if evalEnv, ok := rawKeyValue[EvalEnv]; ok {
		evalEnv, ok := evalEnv.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("eval_env must be a mapping")
		}

		err := parseAndValidateEvalEnvForConfig(evalEnv, cfg)
		if err != nil {
			return err
		}
	}

	if before, ok := rawKeyValue[config.BEFORE]; ok {
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

func validateTopLevelFields(rawKeyValue map[string]interface{}, validFields []string) error {
	for key := range rawKeyValue {
		if !util.IsStringInList(key, validFields) {
			return fmt.Errorf("unknown top-level field '%s'", key)
		}
	}

	return nil
}

func parseConfig(rawKeyValue map[string]interface{}, cfg *config.Config) error { //nolint:cyclop
	if err := validateTopLevelFields(rawKeyValue, config.ValidConfigFields); err != nil {
		return err
	}

	if err := parseConfigGeneral(rawKeyValue, cfg); err != nil {
		return err
	}

	if version, ok := rawKeyValue[config.VERSION]; ok {
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

	if shell, ok := rawKeyValue[config.SHELL]; ok {
		shell, ok := shell.(string)
		if !ok {
			return fmt.Errorf("shell must be a string")
		}

		cfg.Shell = shell
	} else {
		return fmt.Errorf("'shell' field is required")
	}

	if mixins, ok := rawKeyValue[config.MIXINS]; ok {
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

// Trim `-` prefix.
// Using this prefix we allow to include non-existed mixins (git-ignored for example).
func normalizeMixinFilename(filename string) string {
	return strings.TrimPrefix(filename, "-")
}

// Ignored means that it is okay if minix does not exist.
// It can be a git-ignored file for example.
func isIgnoredMixin(filename string) bool {
	return strings.HasPrefix(filename, "-")
}

func readAndValidateMixins(mixins []interface{}, cfg *config.Config) error {
	for _, filename := range mixins {
		if filename, ok := filename.(string); ok { //nolint:nestif
			configAbsPath, err := path.GetFullConfigPath(normalizeMixinFilename(filename), cfg.WorkDir)
			if err != nil {
				if isIgnoredMixin(filename) && errors.Is(err, path.ErrFileNotExists) {
					continue
				} else {
					// complain non-existed mixin only if its filename does not starts with dash `-`
					return fmt.Errorf("failed to read mixin config: %w", err)
				}
			}

			fileData, err := os.ReadFile(configAbsPath)
			if err != nil {
				return fmt.Errorf("can not read mixin config file: %w", err)
			}

			mixinCfg := config.NewMixinConfig(cfg.WorkDir, filename, cfg.DotLetsDir)

			if err := parseMixinConfig(fileData, mixinCfg); err != nil {
				return fmt.Errorf("failed to load mixin config: %w", err)
			}

			if err := mergeConfigs(cfg, mixinCfg); err != nil {
				return fmt.Errorf("failed to merge mixin config %s with main config: %w", filename, err)
			}
		} else {
			return newConfigParseError(
				"must be a string",
				config.MIXINS,
				"list item",
			)
		}
	}

	return nil
}

func parseMixinConfig(data []byte, mixinCfg *config.Config) error {
	rawKeyValue := make(map[string]interface{})

	if err := yaml.Unmarshal(data, &rawKeyValue); err != nil {
		return fmt.Errorf("can not decode mixin config file: %w", err)
	}

	if err := validateTopLevelFields(rawKeyValue, config.ValidMixinConfigFields); err != nil {
		return err
	}

	return parseConfigGeneral(rawKeyValue, mixinCfg)
}

// Merge main and mixin configs. If there is a conflict - return error as we do not override values
// TODO add test.
func mergeConfigs(mainCfg *config.Config, mixinCfg *config.Config) error {
	for _, mixinCmd := range mixinCfg.Commands {
		if _, conflict := mainCfg.Commands[mixinCmd.Name]; conflict {
			return fmt.Errorf("parser %s from mixin is already declared in main config's commands", mixinCmd.Name)
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

func parseAndValidateBefore(before string, cfg *config.Config) error {
	cfg.Before = before

	return nil
}

func parseCommands(cmds map[interface{}]interface{}) ([]config.Command, error) {
	var commands []config.Command
	for key, value := range cmds {
		keyStr, ok := key.(string)
		if !ok {
			return []config.Command{}, newConfigParseError(
				"parser name must be a string",
				config.COMMANDS,
				"",
			)
		}

		newCmd := config.NewCommand(keyStr)

		err := parseAndValidateCommand(&newCmd, value.(map[interface{}]interface{}))
		if err != nil {
			return []config.Command{}, err
		}

		commands = append(commands, newCmd)
	}

	return commands, nil
}

func joinBeforeScripts(beforeScripts ...string) string {
	buf := new(bytes.Buffer)

	for _, script := range beforeScripts {
		buf.WriteString(script)
		buf.WriteString("\n")
	}

	return buf.String()
}

// Parse file data into config.
func Parse(data []byte, cfg *config.Config) error {
	rawKeyValue := make(map[string]interface{})

	if err := yaml.Unmarshal(data, &rawKeyValue); err != nil {
		return err
	}

	return parseConfig(rawKeyValue, cfg)
}
