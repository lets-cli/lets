package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/lets-cli/lets/checksum"
)

type Command struct {
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
	Env map[string]string
	// env from -E flag
	OverrideEnv map[string]string
	// store docopts from options directive
	Docopts     string
	SkipDocopts bool // default false
	Options     map[string]string
	CliOptions  map[string]string
	Depends     map[string]Dep
	// store depends commands in order declared in config
	DependsNames    []string
	ChecksumMap     map[string]string
	PersistChecksum bool

	// args with command name
	// e.g. from 'lets run --debug' we will get [run, --debug]
	Args []string
	// args without command name
	// e.g. from 'lets run --debug' we will get [--debug]
	CommandArgs []string

	// run only specified commands from cmd map
	Only []string
	// run all but excluded commands from cmd map
	Exclude []string

	// if command has declared checksum
	HasChecksum     bool
	ChecksumSources map[string][]string
	// store loaded persisted checksums here
	persistedChecksums map[string]string

	// ref is basically a command name to use with predefined args, env
	Ref string
	// can be specified only with ref
	RefArgs []string
}

// NewCommand creates new command struct.
func NewCommand(name string) Command {
	return Command{
		Name:        name,
		Env:         make(map[string]string),
		SkipDocopts: false,
	}
}

func (cmd Command) WithArgs(args []string) Command {
	newCmd := cmd
	newCmd.Args = args

	return newCmd
}

func (cmd Command) FromRef(refCommand Command) Command {
	newCmd := cmd

	if len(newCmd.Args) == 0 {
		newCmd.Args = append([]string{cmd.Name}, refCommand.RefArgs...)
	} else {
		newCmd.Args = append(newCmd.Args, refCommand.RefArgs...)
	}

	newCmd.CommandArgs = newCmd.Args[1:]

	return newCmd
}

func (cmd Command) WithEnv(env map[string]string) Command {
	newCmd := cmd
	for key, val := range env {
		newCmd.Env[key] = val
	}

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

	return buf.String()
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
func (cmd *Command) ReadChecksumsFromDisk(dotLetsDir string, cmdName string, checksumMap map[string]string) error {
	checksums := make(map[string]string, len(checksumMap)+1)

	for checksumName := range checksumMap {
		filename := checksumName
		if checksumName == checksum.DefaultChecksumKey {
			filename = checksum.DefaultChecksumFileName
		}
		checksumResult, err := checksum.ReadChecksumFromDisk(dotLetsDir, cmdName, filename)
		if err != nil {
			return err
		}

		checksums[checksumName] = checksumResult
	}

	cmd.persistedChecksums = checksums

	return nil
}
