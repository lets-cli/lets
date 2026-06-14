package executor

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/lets-cli/lets/internal/config/config"
	"github.com/lets-cli/lets/internal/env"
	log "github.com/sirupsen/logrus"
)

// ScriptRunner executes a shell script in the context of a command.
type ScriptRunner func(command *config.Command, script string) error

// NewShellRunner returns a ScriptRunner that spawns a real OS process.
func NewShellRunner(cfg *config.Config, out io.Writer) ScriptRunner {
	r := &shellRunner{cfg: cfg, out: out}
	return r.run
}

type shellRunner struct {
	cfg *config.Config
	out io.Writer
}

func (r *shellRunner) run(command *config.Command, cmdScript string) error {
	script := joinBeforeAndScript(r.cfg.Before, cmdScript)

	shell := r.cfg.Shell
	if command.Shell != "" {
		shell = command.Shell
	}

	args := []string{"-c", script}
	if len(command.Args) > 0 {
		args = append(args, "--", strings.Join(command.Args, " "))
	}

	// shell and script come from developer-authored lets.yaml, not user runtime input.
	// User CLI args land in args as positional parameters ($1, $2, …) after "--", not as code.
	osCmd := exec.Command(shell, args...) //nolint:gosec
	osCmd.Stdout = r.out
	osCmd.Stderr = r.out
	osCmd.Stdin = os.Stdin

	osCmd.Dir = r.cfg.WorkDir
	if command.WorkDir != "" {
		osCmd.Dir = command.WorkDir
	}

	if err := r.setupEnv(osCmd, command, shell); err != nil {
		return err
	}

	if env.DebugLevel() > 1 {
		log.Debugf("executing:\n  script: %s\n  env: %s", cmdScript, fmtEnv(osCmd.Env))
	}

	return osCmd.Run()
}

func (r *shellRunner) setupEnv(osCmd *exec.Cmd, command *config.Command, shell string) error {
	defaultEnv := r.cfg.CommandBuiltinEnv(command, shell, osCmd.Dir)

	checksumEnvMap := getChecksumEnvMap(command.ChecksumMap)

	var changedChecksumEnvMap map[string]string
	if command.PersistChecksum {
		changedChecksumEnvMap = getChangedChecksumEnvMap(
			command.ChecksumMap,
			command.GetPersistedChecksums(),
		)
	}

	cmdEnv, err := command.GetEnv(*r.cfg, defaultEnv)
	if err != nil {
		return err
	}

	envMaps := []map[string]string{
		defaultEnv,
		r.cfg.GetEnv(),
		cmdEnv,
		command.Options,
		command.CliOptions,
		checksumEnvMap,
		changedChecksumEnvMap,
	}

	envList := os.Environ()
	for _, envMap := range envMaps {
		envList = append(envList, convertEnvMapToList(envMap)...)
	}

	osCmd.Env = envList

	return nil
}

func joinBeforeAndScript(before string, script string) string {
	if before == "" {
		return script
	}

	before = strings.TrimSpace(before)

	return strings.Join([]string{before, script}, "\n")
}
