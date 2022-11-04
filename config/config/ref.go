package config

import (
	"errors"

	"github.com/kballard/go-shellquote"
)

type Ref struct {
	Name string
	Args []string
}

type RefArgs []string;


// UnmarshalYAML implements yaml.Unmarshaler interface.
func (a *RefArgs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if *a == nil {
		*a = make(RefArgs, 0)
	}

	var arg string
	if err := unmarshal(&arg); err == nil {
		args, err := shellquote.Split(arg)
		if err != nil {
			return errors.New("can not parse args into list")
		}

		*a = append(*a, args...)

		return nil
	}

	var args []string
	if err := unmarshal(&args); err != nil {
		return err
	}

	*a = append(*a, args...)

	return nil
}