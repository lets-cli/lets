package commands

import (
	"fmt"
	"github.com/kindritskyiMax/lets/commands/command"
	"github.com/kindritskyiMax/lets/config"
	"github.com/kindritskyiMax/lets/logging"
	"io"
	"os/exec"
)

type RunOptions struct {
	Config  *config.Config
	RawArgs []string
}

// RunCommand runs parent command
func RunCommand(cmdToRun command.Command, cfg *config.Config, out io.Writer) error {
	return runCmd(cmdToRun, cfg, out, false)
}

func convertEnvForCmd(envMap map[string]string) []string {
	envList := make([]string, len(envMap))
	for name, value := range envMap {
		envList = append(envList, fmt.Sprintf("%s=%s", name, value))
	}
	return envList
}

func convertOptsToEnvForCmd(opts map[string]string) []string {
	envList := make([]string, len(opts))
	for name, value := range opts {
		envList = append(envList, fmt.Sprintf("%s=%s", name, value))
	}
	return envList
}

func convertChecksumToEnvForCmd(checksum string) []string {
	return []string{fmt.Sprintf("LETS_CHECKSUM=%s", checksum)}
}

func composeEnvs(envs ...[]string) []string {
	var composed []string
	for _, env := range envs {
		composed = append(composed, env...)
	}
	return composed
}

func runCmd(cmdToRun command.Command, cfg *config.Config, out io.Writer, isChild bool) error {
	// TODO get user's current shell
	cmd := exec.Command("sh", "-c", cmdToRun.Cmd)
	// setup std out and err
	cmd.Stdout = out
	cmd.Stderr = out

	// parse docopts
	opts, err := command.ParseDocopts(cmdToRun)
	if err != nil {
		return fmt.Errorf("failed to parse docopt options for cmd %s: %s", cmdToRun.Name, err)
	}
	cmdToRun.Options = opts

	// setup env for command
	env := convertEnvForCmd(cmdToRun.Env)
	optsEnv := convertOptsToEnvForCmd(cmdToRun.Options)
	checksumEnv := convertChecksumToEnvForCmd(cmdToRun.Checksum)
	cmd.Env = composeEnvs(env, optsEnv, checksumEnv)
	if !isChild {
		logging.Log.Debugf("Executing command %s with env:\n%s", cmdToRun.Name, cmd.Env)
	} else {
		logging.Log.Debugf("Executing depend command %s with env:\n%s", cmdToRun.Name, cmd.Env)
	}

	// run depends commands
	for _, dependCmdName := range cmdToRun.Depends {
		dependCmd := cfg.Commands[dependCmdName]
		err := runCmd(dependCmd, cfg, out, true)
		if err != nil {
			// must return error to root
			return err
		}
	}

	return cmd.Run()
}
