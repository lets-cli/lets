package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docopt/docopt-go"

	"github.com/lets-cli/lets/commands/command"
	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/logging"
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

func makeEnvEntry(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}

func convertEnvMapToList(envMap map[string]string) []string {
	var envList []string
	for name, value := range envMap {
		envList = append(envList, makeEnvEntry(name, value))
	}

	return envList
}

func convertChecksumToEnvForCmd(checksum string) []string {
	return []string{makeEnvEntry("LETS_CHECKSUM", checksum)}
}

func convertChecksumMapToEnvForCmd(checksumMap map[string]string) []string {
	normalizeKey := func(origKey string) string {
		key := strings.ReplaceAll(origKey, "-", "_")
		key = strings.ToUpper(key)

		return key
	}

	var envList []string

	for name, value := range checksumMap {
		if name != "" {
			envList = append(envList, makeEnvEntry(fmt.Sprintf("LETS_CHECKSUM_%s", normalizeKey(name)), value))
		}
	}

	return envList
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
	cmd := exec.Command(cfg.Shell, "-c", cmdToRun.Cmd) // #nosec G204
	// setup std out and err
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Stdin = os.Stdin

	// set working directory for command
	cmd.Dir = cfg.WorkDir

	isChildCmd := parentName != ""

	// parse docopts - only for parent
	if !isChildCmd {
		opts, err := command.ParseDocopts(cmdToRun)
		if err != nil {
			return formatOptsUsageError(err, opts, cmdToRun)
		}

		cmdToRun.Options = command.OptsToLetsOpt(opts)
		cmdToRun.CliOptions = command.OptsToLetsCli(opts)
	}

	//calculate checksum if needed
	if err := cmdToRun.ChecksumCalculator(); err != nil {
		return err
	}

	// setup env for command
	cmd.Env = composeEnvs(
		os.Environ(),
		convertEnvMapToList(cfg.Env),
		convertEnvMapToList(cmdToRun.Env),
		convertEnvMapToList(cmdToRun.Options),
		convertEnvMapToList(cmdToRun.CliOptions),
		convertChecksumToEnvForCmd(cmdToRun.Checksum),
		convertChecksumMapToEnvForCmd(cmdToRun.ChecksumMap),
	)

	if !isChildCmd {
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

	runErr := cmd.Run()
	if runErr != nil {
		return runErr
	}

	// persist checksum only if exit code 0
	if cmdToRun.PersistChecksum {
		err := command.PersistCommandsChecksumToDisk(cmdToRun)
		if err != nil {
			return err
		}
	}

	return nil
}
