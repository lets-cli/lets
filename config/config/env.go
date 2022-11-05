package config

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/lets-cli/lets/checksum"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type Envs struct {
	Keys    []string
	Mapping map[string]Env
}

type Env struct {
	Name     string
	Value    string
	Sh       string
	Checksum Checksum
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (e *Envs) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("lets: env is not a map")
	}

	// keys accessed under even indexes
	// values accessed under odd indexes
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		var env Env
		if err := valueNode.Decode(&env); err != nil {
			return err
		}
		env.Name = keyNode.Value
		e.Set(keyNode.Value, env)
	}

	return nil
}

func (e *Envs) Clone() *Envs {
	if e == nil {
		return nil
	}

	mapping := make(map[string]Env, len(e.Mapping))
	for k, v := range e.Mapping {
		mapping[k] = v
	}
	return &Envs{
		Keys:    cloneArray(e.Keys),
		Mapping: mapping,
	}
}

func (e *Envs) Empty() bool {
	if e == nil {
		return true
	}

	return len(e.Keys) == 0
}

// Has checks if a value exists by a given key.
func (e *Envs) Has(key string) bool {
	if e == nil || e.Mapping == nil {
		return false
	}

	_, ok := e.Mapping[key]
	return ok
}

func (e *Envs) Dump() map[string]string {
	if e == nil {
		return map[string]string{}
	}

	envs := make(map[string]string, len(e.Keys))
	for _, k := range e.Keys {
		envs[k] = e.Mapping[k].Value
	}

	return envs
}

// Range allows you to loop into the envs in its right order.
func (e *Envs) Range(yield func(key string, value Env) error) error {
	if e == nil {
		return nil
	}
	for _, k := range e.Keys {
		if err := yield(k, e.Mapping[k]); err != nil {
			return err
		}
	}
	return nil
}

// Merge merges the given Envs into the existing Envs.
func (e *Envs) Merge(other *Envs) {
	_ = other.Range(func(key string, value Env) error {
		e.Set(key, value)
		return nil
	})
}

// MergeMap merges the given map into the existing Envs.
func (e *Envs) MergeMap(other map[string]string) {
	for key, value := range other {
		e.Set(key, Env{Name: key, Value: value})
	}
}

// Set sets a value to a given key.
func (e *Envs) Set(key string, value Env) {
	if e.Mapping == nil {
		e.Mapping = make(map[string]Env, 1)
	}
	if !slices.Contains(e.Keys, key) {
		e.Keys = append(e.Keys, key)
	}
	e.Mapping[key] = value
}

// eval env value and trim result string.
func executeScript(shell string, script string) (string, error) {
	cmd := exec.Command(shell, "-c", script)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("can not get output from eval_env script: %s: %w", script, err)
	}

	res := string(out)
	// TODO get rid of TrimSpace
	return strings.TrimSpace(res), nil
}

// Execute executes env entries for sh scrips and calculate checksums
// It is lazy and caches data on first call.
func (e *Envs) Execute(cfg Config) error {
	if e == nil {
		return nil
	}

	for _, key := range e.Keys {
		env := e.Mapping[key]
		if env.Sh != "" {
			result, err := executeScript(cfg.Shell, env.Sh)
			if err != nil {
				return err
			}
			env.Value = result
			e.Mapping[key] = env
		} else if len(env.Checksum) > 0 {
			result, err := checksum.CalculateChecksum(cfg.WorkDir, env.Checksum[checksum.DefaultChecksumKey])
			if err != nil {
				return err
			}

			env.Value = result
			e.Mapping[key] = env
		}
	}

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (e *Env) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err == nil {
		e.Value = str
		return nil
	}

	var sh struct {
		Sh string
	}

	if err := unmarshal(&sh); err != nil {
		return err
	}

	if sh.Sh != "" {
		e.Sh = sh.Sh
		return nil
	}

	var checksum struct {
		Checksum *Checksum
	}

	if err := unmarshal(&checksum); err != nil {
		return err
	}

	// TODO: current lets implementation does not support checksum map for env
	// TODO: probably we should deprecate command.checksum in favor of
	// cmd.env.VAR.checksum: [file1, file2]
	if len(*checksum.Checksum) > 0 {
		e.Checksum = *checksum.Checksum
		return nil
	}

	return nil
}
