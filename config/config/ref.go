package config

import (
	"errors"
	"os"

	"github.com/kballard/go-shellquote"
)

type Ref struct {
	Name string
	Args []string
}

type RefArgs []string

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

func (r *Ref) Clone() *Ref {
	if r == nil {
		return nil
	}

	return &Ref{
		Name: r.Name,
		Args: cloneArray(r.Args),
	}
}

func ExpandRefArgs(cfg *Config) {
	for _, cmd := range cfg.Commands {
		if cmd.Ref == nil {
			continue
		}

		for idx, arg := range cmd.Ref.Args {
			// we have to expand env here on our own, since this args not came from users tty, and not expanded before lets
			cmd.Ref.Args[idx] = os.Expand(arg, func(key string) string {
				return cfg.Env.Mapping[key].Value
			})
		}
	}
}
