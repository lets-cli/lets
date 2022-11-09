package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
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
	// args from 'lets run --debug' will become [--debug]
	Args []string

	ChecksumSources map[string][]string
	// store loaded persisted checksums here
	persistedChecksums map[string]string

	// ref is basically a command name to use with predefined args, env
	ref *ref
}

type Commands map[string]*Command

func (c *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var short string
	if err := unmarshal(&short); err == nil {
		c.Cmds = Cmds{
			Commands: []*Cmd{{Script: short}},
		}
		c.SkipDocopts = true
		return nil
	}

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
	if c.Docopts == "" {
		c.SkipDocopts = true
	}
	c.Depends = cmd.Depends
	workDir, err := filepath.Abs(cmd.WorkDir)
	if err != nil {
		return err
	}
	c.WorkDir = workDir
	c.After = cmd.After
	// TODO: checksum must be refactored
	if cmd.Checksum != nil {
		c.ChecksumSources = *cmd.Checksum
	}

	c.PersistChecksum = cmd.PersistChecksum
	if len(c.ChecksumSources) == 0 && c.PersistChecksum {
		return errors.New("'persist_checksum' must be used with 'checksum'")
	}

	if cmd.Ref != "" {
		var refArgs struct {
			Args *refArgs
		}
		if err := unmarshal(&refArgs); err != nil {
			return err
		}
		c.ref = &ref{Name: cmd.Ref}
		if refArgs.Args != nil {
			c.ref.Args = *refArgs.Args
		}
	}

	return nil
}

func (c *Command) GetEnv(cfg Config) (map[string]string, error) {
	if err := c.Env.Execute(cfg); err != nil {
		return nil, err
	}

	return c.Env.Dump(), nil
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
		ChecksumSources:    cloneMapSlice(c.ChecksumSources),
		persistedChecksums: cloneMap(c.persistedChecksums),
		Args:               cloneSlice(c.Args),
	}

	return cmd
}

func (c *Command) Dump() string {
	pretty, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return ""
	}

	result := string(pretty)
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
