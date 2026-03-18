package config

import (
	"maps"

	"github.com/lets-cli/lets/internal/checksum"
)

// Checksum type for all checksum uses (env, command.env, command,checksum).
type Checksum map[string][]string

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (c *Checksum) UnmarshalYAML(unmarshal func(any) error) error {
	if *c == nil {
		*c = make(Checksum)
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
