package config

import (
	"errors"
	"fmt"
	"maps"

	"github.com/lets-cli/lets/internal/checksum"
	"gopkg.in/yaml.v3"
)

// ChecksumFiles type for file based checksums.
type ChecksumFiles map[string][]string

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (c *ChecksumFiles) UnmarshalYAML(unmarshal func(any) error) error {
	if *c == nil {
		*c = make(ChecksumFiles)
	}

	var patterns []string
	if err := unmarshal(&patterns); err == nil {
		(*c)[checksum.DefaultChecksumKey] = patterns

		return nil
	}

	var patternsMap map[string][]string
	if err := unmarshal(&patternsMap); err != nil {
		return err
	}

	maps.Copy((*c), patternsMap)

	return nil
}

// Checksum type for env checksum uses.
type Checksum = ChecksumFiles

type CommandChecksum struct {
	Files   ChecksumFiles
	Sh      string
	Persist *bool
}

func (c *CommandChecksum) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		var files ChecksumFiles
		if err := node.Decode(&files); err == nil {
			c.Files = files

			return nil
		}

		return errors.New("checksum must be a list, map, or object")
	}

	var files ChecksumFiles
	if !isCommandChecksumObject(node) {
		if err := node.Decode(&files); err != nil {
			return err
		}

		c.Files = files

		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		value := node.Content[i+1]

		switch key {
		case "files":
			if err := value.Decode(&c.Files); err != nil {
				return err
			}
		case "sh":
			if err := value.Decode(&c.Sh); err != nil {
				return err
			}
		case "persist":
			var persist bool
			if err := value.Decode(&persist); err != nil {
				return err
			}

			c.Persist = &persist
		default:
			return fmt.Errorf("checksum.%s is not supported", key)
		}
	}

	if len(c.Files) > 0 && c.Sh != "" {
		return errors.New("checksum must use only one of 'files' or 'sh'")
	}

	return nil
}

func isCommandChecksumObject(node *yaml.Node) bool {
	for i := 0; i < len(node.Content); i += 2 {
		switch node.Content[i].Value {
		case "files", "sh", "persist":
			return true
		}
	}

	return false
}
