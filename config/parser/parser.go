package parser

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

type RemoteMixin struct {
	URL     string
	Version string

	mixinsDir string
}

// Filename is name of mixin file (hash from url).
func (rm *RemoteMixin) Filename() string {
	hasher := sha256.New()
	hasher.Write([]byte(rm.URL))

	if rm.Version != "" {
		hasher.Write([]byte(rm.Version))
	}

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// Path is abs path to mixin file (.lets/mixins/<filename>).
func (rm *RemoteMixin) Path() string {
	return filepath.Join(rm.mixinsDir, rm.Filename())
}

func (rm *RemoteMixin) persist(data []byte) error {
	f, err := os.OpenFile(rm.Path(), os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		return fmt.Errorf("can not open file %s to persist mixin: %w", rm.Path(), err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("can not write mixin to file %s: %w", rm.Path(), err)
	}

	return nil
}

func (rm *RemoteMixin) exists() bool {
	return util.FileExists(rm.Path())
}

func (rm *RemoteMixin) tryRead() ([]byte, error) {
	if !rm.exists() {
		return nil, nil
	}
	data, err := os.ReadFile(rm.Path())
	if err != nil {
		return nil, fmt.Errorf("can not read mixin config file at %s: %w", rm.Path(), err)
	}

	return data, nil
}

func (rm *RemoteMixin) download() ([]byte, error) {
	// TODO: maybe create a client for this?
	ctx, cancel := context.WithTimeout(context.Background(), 60*5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		rm.URL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 15 * 60 * time.Second, // TODO: move to client struct
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no such file at: %s", rm.URL)
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("network error: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}

func readAndValidateMixins(mixins []interface{}, cfg *config.Config) error {
	if err := cfg.CreateMixinsDir(); err != nil {
		return err
	}

	for _, mixin := range mixins {
		if filename, ok := mixin.(string); ok { //nolint:nestif
			configAbsPath, err := path.GetFullConfigPath(normalizeMixinFilename(filename), cfg.WorkDir)
			if err != nil {
				if isIgnoredMixin(filename) && errors.Is(err, path.ErrFileNotExists) {
					continue
				} else {
					// complain non-existed mixin only if its filename does not start with dash `-`
					return fmt.Errorf("failed to read mixin config: %w", err)
				}
			}
			fileData, err := os.ReadFile(configAbsPath)
			if err != nil {
				return fmt.Errorf("can not read mixin config file: %w", err)
			}

			mixinCfg := config.NewMixinConfig(cfg, filename)
			if err := parseMixinConfig(fileData, mixinCfg); err != nil {
				return fmt.Errorf("failed to load mixin config '%s': %w", filename, err)
			}

			if err := mergeConfigs(cfg, mixinCfg); err != nil {
				return fmt.Errorf("failed to merge mixin config %s with main config: %w", filename, err)
			}
		} else if mixinMapping, ok := mixin.(map[string]interface{}); ok {
			rm := &RemoteMixin{mixinsDir: cfg.MixinsDir}
			if url, ok := mixinMapping["url"]; ok {
				// TODO check if url is valid
				rm.URL, _ = url.(string)
			}

			if version, ok := mixinMapping["version"]; ok {
				rm.Version, _ = version.(string)
			}

			data, err := rm.tryRead()
			if err != nil {
				return err
			}

			if data == nil {
				data, err = rm.download()
				if err != nil {
					return err
				}
			}

			// TODO: what if multiple mixins have same commands
			//  1 option - fail and suggest use to namespace all commands in remote mixin
			//  2 option - namespace it (this may require specifying namespace in mixin config or in main config mixin section)
			mixinCfg := config.NewMixinConfig(cfg, rm.Filename())
			if err := parseMixinConfig(data, mixinCfg); err != nil {
				return fmt.Errorf("failed to load remote mixin config '%s': %w", rm.URL, err)
			}

			if err := mergeConfigs(cfg, mixinCfg); err != nil {
				return fmt.Errorf("failed to merge remote mixin config %s with main config: %w", rm.URL, err)
			}

			if err := rm.persist(data); err != nil {
				return fmt.Errorf("failed to persist remote mixin config %s: %w", rm.URL, err)
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

func parseCommands(cmds map[string]interface{}, cfg *config.Config) ([]config.Command, error) {
	var commands []config.Command
	for rawName, rawValue := range cmds {
		rawCmd := map[string]interface{}{}

		switch rawValue := rawValue.(type) {
		case map[string]interface{}:
			rawCmd = rawValue
		case map[interface{}]interface{}:
			for key, value := range rawValue {
				key, ok := key.(string)
				if !ok {
					return []config.Command{}, newConfigParseError(
						"command directive must be a string",
						rawName,
						"",
					)
				}
				rawCmd[key] = value
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
