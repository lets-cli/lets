package config

import (
	"encoding/json"
)

type Command struct {
	Name string
	// script to run
	Cmd string
	// script to run after cmd finished (cleanup, etc)
	After string
	// map of named scripts to run in parallel
	CmdMap      map[string]string
	Description string
	// env from command
	Env map[string]string
	// env from -E flag
	OverrideEnv map[string]string
	// store docopts from options directive
	Docopts         string
	SkipDocopts     bool // default false
	Options         map[string]string
	CliOptions      map[string]string
	Depends         map[string]Dep
	Checksum        string
	ChecksumMap     map[string]string
	PersistChecksum bool

	// prepared args - started from command name
	Args []string

	// run only specified commands from cmd map
	Only []string
	// run all but excluded commands from cmd map
	Exclude []string

	// if command has declared checksum
	HasChecksum    bool
	ChecksumSource map[string][]string
	// store loaded persisted checksums here
	persistedChecksums map[string]string
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

func (cmd Command) Pretty() string {
	pretty, _ := json.MarshalIndent(cmd, "", "  ")

	return string(pretty)
}

func (cmd *Command) ChecksumCalculator(workDir string) error {
	if len(cmd.ChecksumSource) == 0 {
		return nil
	}

	return calculateChecksumFromSource(workDir, cmd)
}

func (cmd *Command) GetPersistedChecksums() map[string]string {
	return cmd.persistedChecksums
}

// ReadChecksumsFromDisk reads all checksums for cmd into map.
func (cmd *Command) ReadChecksumsFromDisk(dotLetsDir string, cmdName string, checksumMap map[string]string) error {
	checksums := make(map[string]string, len(checksumMap)+1)

	for checksumName := range checksumMap {
		checksum, err := readOneChecksum(dotLetsDir, cmdName, checksumName)
		if err != nil {
			return err
		}

		checksums[checksumName] = checksum
	}

	checksum, err := readOneChecksum(dotLetsDir, cmdName, DefaultChecksumName)
	if err != nil {
		return err
	}

	checksums[DefaultChecksumName] = checksum

	cmd.persistedChecksums = checksums

	return nil
}
