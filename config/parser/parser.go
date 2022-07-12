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
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type ConfigParseError struct {
	Path struct {
		Full  string
		Field string
	}
	Err error
}

func (e *ConfigParseError) Error() string {
	return fmt.Sprintf("failed to parse config: %s", e.Err)
}

// TODO refactor this.
func newConfigParseError(msg string, name string, field string) error {
	fields := []string{name, field}
	sep := "."

	if field == "" {
		sep = ""
	}

	fullPath := strings.Join(fields, sep)

	return &ConfigParseError{
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

func parseConfigGeneral(rawKeyValue map[string]interface{}, cfg *config.Config) error {
	rawEnv := make(map[string]interface{})

	if env, ok := rawKeyValue[ENV]; ok {
		env, ok := env.(map[string]interface{})
		if !ok {
			return fmt.Errorf("env must be a mapping")
		}
		for k, v := range env {
			rawEnv[k] = v
		}
	}

	if evalEnv, ok := rawKeyValue[EvalEnv]; ok {
		log.Debug("eval_env is deprecated, consider using 'env' with 'sh' executor")
		evalEnv, ok := evalEnv.(map[string]interface{})
		if !ok {
			return fmt.Errorf("eval_env must be a mapping")
		}

		for k, v := range evalEnv {
			rawEnv[k] = map[string]interface{}{"sh": v}
		}
	}

	envEntries, err := parseEnvEntries(rawEnv, cfg)
	if err != nil {
		return err
	}

	for _, entry := range envEntries {
		value, err := entry.Value()
		if err != nil {
			return parseDirectiveError(
				"env",
				fmt.Sprintf("can not get value for '%s' env variable", entry.Name()),
			)
		}

		cfg.Env[entry.Name()] = value
	}

	if before, ok := rawKeyValue[config.BEFORE]; ok {
		before, ok := before.(string)
		if !ok {
			return fmt.Errorf("before must be a string")
		}

		err := parseBefore(before, cfg)
		if err != nil {
			return err
		}
	}

	if cmds, ok := rawKeyValue[config.COMMANDS]; ok {
		cmdsMap, ok := cmds.(map[string]interface{})
		if !ok {
			return newConfigParseError(
				"must be a mapping",
				config.COMMANDS,
				"",
			)
		}

		commands, err := parseCommands(cmdsMap, cfg)
		if err != nil {
			return err
		}

		for _, c := range commands {
			cfg.Commands[c.Name] = c
		}
	}

	return nil
}

func parseConfig(rawKeyValue map[string]interface{}, cfg *config.Config) error {
	for key := range rawKeyValue {
		if !config.ValidConfigDirectives.Contains(key) {
			return fmt.Errorf("unknown top-level field '%s'", key)
		}
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

	if rawPlugins, ok := rawKeyValue[config.PLUGINS]; ok {
		plugins, ok := rawPlugins.(map[string]interface{})
		if !ok {
			return fmt.Errorf("plugins must be a mapping")
		}
		if err := parseConfigPlugins(plugins, cfg); err != nil {
			return err
		}
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

	postprocessRefArgs(cfg)

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
				return fmt.Errorf("failed to load mixin config '%s': %w", filename, err)
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

	for key := range rawKeyValue {
		if !config.ValidMixinConfigDirectives.Contains(key) {
			return fmt.Errorf("unknown top-level field '%s'", key)
		}
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

func parseBefore(before string, cfg *config.Config) error {
	cfg.Before = before

	return nil
}

func parseConfigPlugins(rawPlugins map[string]interface{}, cfg *config.Config) error {
	plugins := make(map[string]config.ConfigPlugin)

	for key, value := range rawPlugins {
		pluginConfig, ok := value.(map[string]interface{})
		if !ok {
			// TODO maybe print plugin configuration schema
			return fmt.Errorf("plugin %s configuration must be a mapping", key)
		}

		plugin := config.ConfigPlugin{Name: key}

		for configKey, configVal := range pluginConfig {
			switch configVal := configVal.(type) {
			case string:
				switch configKey {
				case "version":
					plugin.Version = configVal
				case "url":
					plugin.Url = configVal
				case "bin":
					plugin.Bin = configVal
				case "repo":
					plugin.Repo = configVal
				}
			}
		}

		plugins[key] = plugin
	}

	cfg.Plugins = plugins

	return nil
}

func parseCommands(cmds map[string]interface{}, cfg *config.Config) ([]config.Command, error) {
	var commands []config.Command
	for rawName, rawValue := range cmds {
		rawCmd := map[string]interface{}{}

		switch rawValue := rawValue.(type) {
		case map[string]interface{}:
			rawCmd = rawValue
		case map[interface{}]interface{}:
			for k, v := range rawValue {
				k, ok := k.(string)
				if !ok {
					return []config.Command{}, newConfigParseError(
						"command directive must be a string",
						rawName,
						"",
					)
				}
				rawCmd[k] = v
			}
		default:
			return []config.Command{}, newConfigParseError(
				"command name must be a string",
				config.COMMANDS,
				"",
			)
		}

		newCmd := config.NewCommand(rawName)

		err := parseCommand(&newCmd, rawCmd, cfg)
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
		if script == "" {
			continue
		}
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
