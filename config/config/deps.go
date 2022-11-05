package config

import (
	"errors"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)


type Dep struct {
	Name string
	Args []string
	Env  *Envs
}

type Deps struct {
	Keys    []string
	Mapping map[string]Dep
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (d *Deps) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.SequenceNode {
		return errors.New("lets: 'depends' must be a sequence")
	}

	for i := 0; i < len(node.Content); i += 1 {
		node := node.Content[i]

		var dep Dep
		if err := node.Decode(&dep); err != nil {
			return err
		}
		d.Set(dep.Name, dep)
	}

	return nil
}

func (d *Deps) Clone() *Deps {
	if d == nil {
		return nil
	}

	mapping := make(map[string]Dep, len(d.Mapping))
	for k, v := range d.Mapping {
		mapping[k] = v.Clone()
	}

	return &Deps{
		Keys: cloneArray(d.Keys),
		Mapping: mapping,
	}
}

// Range allows you to loop into the Deps in its right order
func (d *Deps) Range(yield func(key string, value Dep) error) error {
	if d == nil {
		return nil
	}

	for _, k := range d.Keys {
		if err := yield(k, d.Mapping[k]); err != nil {
			return err
		}
	}
	return nil
}


// Set sets a value to a given key
func (d *Deps) Set(key string, value Dep) {
	if d.Mapping == nil {
		d.Mapping = make(map[string]Dep, 1)
	}
	if !slices.Contains(d.Keys, key) {
		d.Keys = append(d.Keys, key)
	}
	d.Mapping[key] = value
}

// Get get a value by a given key
func (d *Deps) Get(key string) *Dep {
	if d == nil || d.Mapping == nil {
		return nil
	}

	dep, ok := d.Mapping[key]
	if !ok {
		return nil
	}

	return &dep
}

// Has checks if a value exists by a given key
func (d *Deps) Has(key string) bool {
	if d == nil || d.Mapping == nil {
		return false
	}

	_, ok := d.Mapping[key]
	return ok
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (d *Dep) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var cmdName string
	if err := unmarshal(&cmdName); err == nil {
		d.Name = cmdName
		return nil
	}

	var cmd struct {
		Name string
		Env *Envs
	}

	if err := unmarshal(&cmd); err != nil {
		return err
	}

	d.Name = cmd.Name
	d.Env = cmd.Env

	var cmdArgsStr struct {
		Args string
	}

	if err := unmarshal(&cmdArgsStr); err == nil {
		d.Args = append([]string{cmd.Name}, cmdArgsStr.Args)
		return nil
	}

	var cmdArgs struct {
		Args []string
	}

	if err := unmarshal(&cmdArgs); err != nil {
		return err
	}

	// args always must start with a dependency name, otherwise docopt will fail
	d.Args = append([]string{cmd.Name}, cmdArgs.Args...)

	return nil
}

func (d Dep) Clone() Dep {
	return Dep{
		Name: d.Name,
		Args: cloneArray(d.Args),
		Env: d.Env.Clone(),
	}
}
