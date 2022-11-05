package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lets-cli/lets/checksum"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	Name string
	// Represents a list of commands (scripts)
	Cmds Cmds
	// script to run after cmd finished (cleanup, etc)
	After string
	// overrides global shell for this particular command
	Shell string
	// overrides global workdir (where lets.yaml is located) for this particular command
	WorkDir     string
	Description string
	// env from command
	Env *Envs
	// store docopts from options directive
	Docopts         string
	SkipDocopts     bool // default false
	Options         map[string]string
	CliOptions      map[string]string
	Depends         *Deps
	ChecksumMap     map[string]string
	PersistChecksum bool

	// args with command name
	// e.g. from 'lets run --debug' we will get [run, --debug]
	Args []string

	ChecksumSources map[string][]string
	// store loaded persisted checksums here
	persistedChecksums map[string]string

	// ref is basically a command name to use with predefined args, env
	Ref *Ref
}

type Commands map[string]*Command

func (c *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// TODO: implement short cmd syntax

	var cmd struct {
		Cmd             Cmds
		Description     string
		Shell           string
		Env             *Envs
		EvalEnv         *Envs `yaml:"eval_env"`
		Options         string
		Depends         *Deps
		WorkDir         string `yaml:"work_dir"`
		After           string
		Ref             string
		Checksum        *Checksum
		PersistChecksum bool `yaml:"persist_checksum"`
	}

	if err := unmarshal(&cmd); err != nil {
		return err
	}

	c.Cmds = cmd.Cmd
	c.Description = cmd.Description
	c.Env = cmd.Env

	// support for deprecated eval_env
	if !cmd.EvalEnv.Empty() {
		log.Debug("eval_env is deprecated, consider using 'env' with 'sh' executor")
	}
	_ = cmd.EvalEnv.Range(func(name string, value Env) error {
		c.Env.Set(name, Env{Name: name, Sh: value.Value})

		return nil
	})

	c.Shell = cmd.Shell
	c.Docopts = cmd.Options
	c.Depends = cmd.Depends
	c.WorkDir = cmd.WorkDir
	c.After = cmd.After
	// TODO: checksum must be refactored, first name of var is misleading
	if cmd.Checksum != nil {
		c.ChecksumSources = *cmd.Checksum
	}

	c.PersistChecksum = cmd.PersistChecksum
	if len(c.ChecksumSources) == 0 && c.PersistChecksum {
		return errors.New("'persist_checksum' must be used with 'checksum'")
	}

	// TODO: validate if ref points to real command ?
	if cmd.Ref != "" {
		// only parsing Args when ref is set
		var refArgs struct {
			Args *RefArgs
		}
		if err := unmarshal(&refArgs); err != nil {
			return err
		}
		c.Ref = &Ref{Name: cmd.Ref, Args: *refArgs.Args}
	}

	return nil
}

// args without command name
// e.g. from 'lets run --debug' we will get [--debug].
func (c *Command) CommandArgs() []string {
	if len(c.Args) == 0 {
		return []string{}
	}

	return c.Args[1:]
}

func (c *Command) GetEnv(cfg Config) (map[string]string, error) {
	if err := c.Env.Execute(cfg); err != nil {
		// TODO: move execution to somevere else. probably make execution lazy and cached
		return nil, err
	}

	return c.Env.Dump(), nil
}

func (c *Command) WithArgs(args []string) *Command {
	newCmd := c.Clone()
	newCmd.Args = args

	return newCmd
}

func (c *Command) FromRef(ref *Ref) *Command {
	newCmd := c.Clone()

	if len(newCmd.Args) == 0 {
		newCmd.Args = append([]string{c.Name}, ref.Args...)
	} else {
		newCmd.Args = append(newCmd.Args, ref.Args...)
	}

	return newCmd
}

func (c *Command) WithEnv(env *Envs) *Command {
	newCmd := c.Clone()
	newCmd.Env.Merge(env)

	return newCmd
}

func (c *Command) Clone() *Command {
	cmd := &Command{
		Name:               c.Name,
		Cmds:               c.Cmds.Clone(),
		After:              c.After,
		Shell:              c.Shell,
		WorkDir:            c.WorkDir,
		Description:        c.Description,
		Env:                c.Env.Clone(),
		Docopts:            c.Docopts,
		SkipDocopts:        c.SkipDocopts,
		Options:            cloneMap(c.Options),
		CliOptions:         cloneMap(c.CliOptions),
		Depends:            c.Depends.Clone(),
		ChecksumMap:        cloneMap(c.ChecksumMap),
		PersistChecksum:    c.PersistChecksum,
		ChecksumSources:    cloneMapArray(c.ChecksumSources),
		persistedChecksums: cloneMap(c.persistedChecksums),
		Ref:                c.Ref.Clone(),
		Args:               cloneArray(c.Args),
	}

	return cmd
}

func (c *Command) Pretty() string {
	pretty, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return ""
	}

	result := string(pretty)
	result = strings.TrimLeft(result, "{")
	result = strings.TrimRight(result, "}")

	return strings.TrimSpace(result)
}

func (c *Command) Help() string {
	buf := new(bytes.Buffer)
	if c.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", c.Description))
	}

	if c.Docopts != "" {
		buf.WriteString(c.Docopts)
	}

	if buf.Len() == 0 {
		buf.WriteString(fmt.Sprintf("No help message for '%s'", c.Name))
	}

	return strings.TrimSuffix(buf.String(), "\n")
}

func (c *Command) ChecksumCalculator(workDir string) error {
	if len(c.ChecksumSources) == 0 {
		return nil
	}

	checksumMap, err := checksum.CalculateChecksumFromSources(workDir, c.ChecksumSources)
	if err != nil {
		return err
	}

	c.ChecksumMap = checksumMap

	return nil
}

func (c *Command) GetPersistedChecksums() map[string]string {
	return c.persistedChecksums
}

// ReadChecksumsFromDisk reads all checksums for cmd into map.
func (c *Command) ReadChecksumsFromDisk(checksumsDir string, cmdName string, checksumMap map[string]string) error {
	checksums := make(map[string]string, len(checksumMap)+1)

	for checksumName := range checksumMap {
		filename := checksumName
		if checksumName == checksum.DefaultChecksumKey {
			filename = checksum.DefaultChecksumFileName
		}
		checksumResult, err := checksum.ReadChecksumFromDisk(checksumsDir, cmdName, filename)
		if err != nil {
			return err
		}

		checksums[checksumName] = checksumResult
	}

	c.persistedChecksums = checksums

	return nil
}
