package commands

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/kindritskyiMax/lets/commands/command"
	"github.com/kindritskyiMax/lets/config"
	"github.com/kindritskyiMax/lets/logging"
	"io"
	"os"
	"os/exec"
)

const (
	NoticeColor = "\033[1;36m%s\033[0m"
)

type RunOptions struct {
	Config  *config.Config
	RawArgs []string
}

// RunCommand runs parent command
func RunCommand(cmdToRun command.Command, cfg *config.Config, out io.Writer) error {
	return runCmd(cmdToRun, cfg, out, "")
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

// format docopts error and adds usage string to output
func formatOptsUsageError(err error, opts docopt.Opts, cmdToRun command.Command) error {
	if opts == nil && err.Error() == "" {
		err = fmt.Errorf("no such option")
	}
	errTpl := fmt.Sprintf("failed to parse docopt options for cmd %s: %s", cmdToRun.Name, err)

	return fmt.Errorf("%s\n\n%s", errTpl, cmdToRun.RawOptions)
}

func runCmd(cmdToRun command.Command, cfg *config.Config, out io.Writer, parentName string) error {
	cmd := exec.Command(cfg.Shell, "-c", cmdToRun.Cmd)
	// setup std out and err
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Stdin = os.Stdin

	// parse docopts
	opts, err := command.ParseDocopts(cmdToRun)
	if err != nil {
		return formatOptsUsageError(err, opts, cmdToRun)
	}
	cmdToRun.Options = command.OptsToLetsOpt(opts)
	cmdToRun.CliOptions = command.OptsToLetsCli(opts)

	// setup env for command
	env := convertEnvForCmd(cmdToRun.Env)
	optsEnv := convertOptsToEnvForCmd(cmdToRun.Options)
	cliOptsEnv := convertOptsToEnvForCmd(cmdToRun.CliOptions)
	checksumEnv := convertChecksumToEnvForCmd(cmdToRun.Checksum)
	cmd.Env = composeEnvs(os.Environ(), env, optsEnv, cliOptsEnv, checksumEnv)
	if parentName == "" {
		logging.Log.Debugf(
			"Executing command\nname: %s\ncmd: %s\nenv:\n%s",
			fmt.Sprintf(NoticeColor, cmdToRun.Name),
			fmt.Sprintf(NoticeColor, cmdToRun.Cmd),
			cmd.Env,
		)
	} else {
		logging.Log.Debugf(
			"Executing child command\nparent name: %s\nname: %s\ncmd: %s\nenv:\n%s",
			fmt.Sprintf(NoticeColor, parentName),
			fmt.Sprintf(NoticeColor, cmdToRun.Name),
			fmt.Sprintf(NoticeColor, cmdToRun.Cmd),
			cmd.Env,
		)
	}

	// run depends commands
	for _, dependCmdName := range cmdToRun.Depends {
		dependCmd := cfg.Commands[dependCmdName]
		err := runCmd(dependCmd, cfg, out, cmdToRun.Name)
		if err != nil {
			// must return error to root
			return err
		}
	}

	return cmd.Run()
}
