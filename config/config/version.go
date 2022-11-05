package config

import (
	"errors"

	"github.com/lets-cli/lets/util"
)

type Version string

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (v *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ver string
	if err := unmarshal(&ver); err != nil {
		return errors.New("version must be a valid semver string")
	}

	_, err := util.ParseVersion(ver)
	if err != nil {
		return err
	}

	*v = Version(ver)

	return nil
}
