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
	// TODO refactor to have list of cmds as basic strture - it is easier to unify them
	Name string
	// script to run
	Cmd string
	// script to run after cmd finished (cleanup, etc)
	After string
	// map of named scripts to run in parallel
	CmdMap map[string]string
	// if specified, overrides global shell for this particular command
	Shell string
	// if specified, overrides global workdir (where lets.yaml is located) for this particular command
	WorkDir     string
	Description string
	// env from command
	Env *Envs
	// store docopts from options directive
	Docopts     string
	SkipDocopts bool // default false
	Options     map[string]string
	CliOptions  map[string]string
	Depends     *Deps
	// store depends commands in order declared in config
	ChecksumMap     map[string]string
	PersistChecksum bool

	// args with command name
	// e.g. from 'lets run --debug' we will get [run, --debug]
	Args []string

	// run only specified commands from cmd map
	Only []string
	// run all but excluded commands from cmd map
	Exclude []string

	// if command has declared checksum
	// TODO drop HasChecksum
	HasChecksum     bool
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
		Cmd Cmds
		Description          string
		Shell string
		Env           *Envs
		EvalEnv           *Envs `yaml:"eval_env"`
		Options string
		Depends *Deps
		WorkDir string `yaml:"work_dir"`
		After string
		Ref string
		Checksum *Checksum
		PersistChecksum bool `yaml:"persist_checksum"`
	}

	if err := unmarshal(&cmd); err != nil {
		return err
	}

	if len(cmd.Cmd.commands) == 1 {
		c.Cmd = cmd.Cmd.commands[0].Script
	} else if cmd.Cmd.parallel {
		c.CmdMap = make(map[string]string, len(cmd.Cmd.commands))
		for _, cmd := range cmd.Cmd.commands {
			c.CmdMap[cmd.Name] = cmd.Script
		}
	}
	c.Description = cmd.Description
	c.Env = cmd.Env
	// support deprecated eval_env
	if !cmd.EvalEnv.Empty() {
		log.Debug("eval_env is deprecated, consider using 'env' with 'sh' executor")
	}
	cmd.EvalEnv.Range(func(name string, value Env) error {
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

	//TODO lol, why do we need this field ?
	c.HasChecksum = len(c.ChecksumSources) > 0
	c.PersistChecksum = cmd.PersistChecksum
	if !c.HasChecksum && c.PersistChecksum {
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
func (cmd Command) CommandArgs() []string {
	if len(cmd.Args) == 0 {
		return []string{}
	}

	return cmd.Args[1:]
}

func (c *Command) GetEnv(cfg Config) (map[string]string, error) {
	if err := c.Env.Execute(cfg); err != nil {
		// TODO: move execution to somevere else. probably make execution lazy and cached
		return nil, err
	}

	return c.Env.Dump(), nil
}

// NewCommand creates new command struct.
func NewCommand(name string) Command {
	return Command{
		Name:        name,
		SkipDocopts: false,
	}
}

func (c *Command) WithArgs(args []string) *Command {
	// TODO: use c.Clone() here
	newCmd := c
	newCmd.Args = args

	return newCmd
}

func (c *Command) FromRef(ref *Ref) *Command {
	// TODO: use c.Clone() here
	newCmd := c

	if len(newCmd.Args) == 0 {
		newCmd.Args = append([]string{c.Name}, ref.Args...)
	} else {
		newCmd.Args = append(newCmd.Args, ref.Args...)
	}

	return newCmd
}

func (c *Command) WithEnv(env *Envs) *Command {
	// TODO: use c.Clone() here
	newCmd := c
	newCmd.Env.Merge(env)

	return newCmd
}

func (cmd Command) Pretty() string {
	pretty, err := json.MarshalIndent(cmd, "", "  ")
	if err != nil {
		return ""
	}

	return string(pretty)
}

func (cmd *Command) Help() string {
	buf := new(bytes.Buffer)
	if cmd.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", cmd.Description))
	}

	if cmd.Docopts != "" {
		buf.WriteString(cmd.Docopts)
	}

	if buf.Len() == 0 {
		buf.WriteString(fmt.Sprintf("No help message for '%s'", cmd.Name))
	}

	return strings.TrimSuffix(buf.String(), "\n")
}

func (cmd *Command) ChecksumCalculator(workDir string) error {
	if len(cmd.ChecksumSources) == 0 {
		return nil
	}

	checksumMap, err := checksum.CalculateChecksumFromSources(workDir, cmd.ChecksumSources)
	if err != nil {
		return err
	}

	cmd.ChecksumMap = checksumMap

	return nil
}

func (cmd *Command) GetPersistedChecksums() map[string]string {
	return cmd.persistedChecksums
}

// ReadChecksumsFromDisk reads all checksums for cmd into map.
func (cmd *Command) ReadChecksumsFromDisk(checksumsDir string, cmdName string, checksumMap map[string]string) error {
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

	cmd.persistedChecksums = checksums

	return nil
}
