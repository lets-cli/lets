package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lets-cli/lets/internal/util"
	"gopkg.in/yaml.v3"
)

type EnvFile struct {
	Name     string
	Required bool
}

type EnvFiles struct {
	Items  []EnvFile
	loaded map[string]string
	ready  bool
}

func (e *EnvFile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var filename string
	// try parse as scalar
	if err := unmarshal(&filename); err == nil {
		e.Name = normalizeEnvFilename(filename)
		e.Required = !isOptionalEnvFilename(filename)
		if e.Name == "" {
			return errors.New("env_file name can not be empty")
		}

		return nil
	}

	var raw struct {
		Name     string
		Required *bool
	}
	// try parse as map
	if err := unmarshal(&raw); err != nil {
		return err
	}

	if raw.Name == "" {
		return errors.New("env_file name can not be empty")
	}

	if isOptionalEnvFilename(raw.Name) {
		return errors.New("env_file map form does not support '-' prefix in name; use required: false instead")
	}

	e.Name = raw.Name
	e.Required = true
	if raw.Required != nil {
		e.Required = *raw.Required
	}

	return nil
}

func (e *EnvFiles) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode, yaml.MappingNode:
		var item EnvFile
		if err := node.Decode(&item); err != nil {
			return err
		}

		e.Items = []EnvFile{item}
		return nil
	case yaml.SequenceNode:
		items := make([]EnvFile, 0, len(node.Content))
		for _, itemNode := range node.Content {
			var item EnvFile
			if err := itemNode.Decode(&item); err != nil {
				return err
			}
			items = append(items, item)
		}

		e.Items = items
		return nil
	default:
		return errors.New("env_file must be a string, map, or sequence")
	}
}

func (e *EnvFiles) Clone() *EnvFiles {
	if e == nil {
		return nil
	}

	items := make([]EnvFile, len(e.Items))
	copy(items, e.Items)

	return &EnvFiles{
		Items: items,
	}
}

func (e *EnvFiles) Append(other *EnvFiles) {
	if other == nil || len(other.Items) == 0 {
		return
	}

	e.Items = append(e.Items, other.Items...)
}

func (e *EnvFiles) Load(cfg Config, envMap map[string]string) (map[string]string, error) {
	if e == nil {
		return map[string]string{}, nil
	}

	if e.ready {
		return cloneMap(e.loaded), nil
	}

	loaded := make(map[string]string)

	for _, item := range e.Items {
		filename := expandWithEnv(item.Name, envMap)
		if strings.TrimSpace(filename) == "" {
			return nil, fmt.Errorf("env_file %q resolved to empty path", item.Name)
		}

		if !filepath.IsAbs(filename) {
			filename = filepath.Join(cfg.WorkDir, filename)
		}

		if !util.FileExists(filename) {
			if item.Required {
				return nil, fmt.Errorf("env_file %q does not exist", filename)
			}

			continue
		}

		values, err := godotenv.Read(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to parse env_file %q: %w", filename, err)
		}

		for key, value := range values {
			loaded[key] = value
		}
	}

	e.loaded = loaded
	e.ready = true

	return cloneMap(loaded), nil
}

func normalizeEnvFilename(filename string) string {
	return strings.TrimPrefix(filename, "-")
}

func isOptionalEnvFilename(filename string) bool {
	return strings.HasPrefix(filename, "-")
}

func expandWithEnv(value string, envMap map[string]string) string {
	return os.Expand(value, func(key string) string {
		if envMap != nil {
			if value, exists := envMap[key]; exists {
				return value
			}
		}

		// If lets own env does not have declared env var, fallback to os env
		return os.Getenv(key)
	})
}
