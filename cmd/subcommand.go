package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/runner"
	"github.com/spf13/cobra"
)

// cut all elements before command name.
func prepareArgs(cmd config.Command, originalArgs []string) []string {
	nameIdx := 0

	for idx, arg := range originalArgs {
		if arg == cmd.Name {
			nameIdx = idx
		}
	}

	return originalArgs[nameIdx:]
}

func replaceGenericCmdPlaceholder(commandName string, cmd config.Command) string {
	genericCmdTplPlaceholder := fmt.Sprintf("${%s}", runner.GenericCmdNameTpl)
	// replace only one placeholder in options
	return strings.Replace(cmd.Docopts, genericCmdTplPlaceholder, commandName, 1)
}

// newCmdGeneric creates new cobra root sub command from Command.
func newCmdGeneric(cmdToRun config.Command, conf *config.Config, out io.Writer) *cobra.Command {
	subCmd := &cobra.Command{
		Use:   cmdToRun.Name,
		Short: cmdToRun.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			only, exclude, err := parseOnlyAndExclude(cmd)
			if err != nil {
				return err
			}

			cmdToRun.Only = only
			cmdToRun.Exclude = exclude
			cmdToRun.Args = prepareArgs(cmdToRun, os.Args)
			cmdToRun.CommandArgs = cmdToRun.Args[1:]

			envs, err := parseEnvFlag(cmd)
			if err != nil {
				return err
			}

			cmdToRun.Docopts = replaceGenericCmdPlaceholder(cmdToRun.Name, cmdToRun)

			cmdToRun.OverrideEnv = envs

			return runner.NewRunner(&cmdToRun, conf, out).Execute(cmd.Context())
		},
		// we use docopt to parse flags on our own, so any flag is valid flag here
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: false,
		SilenceUsage:       true,
	}

	// try print docopt as help for command
	subCmd.SetHelpFunc(func(c *cobra.Command, strings []string) {
		buf := new(bytes.Buffer)
		if cmdToRun.Description != "" {
			buf.WriteString(fmt.Sprintf("%s\n\n", cmdToRun.Description))
		}
		buf.WriteString(cmdToRun.Docopts)

		_, err := buf.WriteTo(c.OutOrStdout())
		if err != nil {
			c.Println(err)
		}
	})

	return subCmd
}

// initialize all commands dynamically from config.
func initSubCommands(rootCmd *cobra.Command, conf *config.Config, out io.Writer) {
	for _, cmdToRun := range conf.Commands {
		rootCmd.AddCommand(newCmdGeneric(cmdToRun, conf, out))
	}
}

func parseOnlyAndExclude(cmd *cobra.Command) (only []string, exclude []string, err error) {
	onlyCmds, err := cmd.Parent().Flags().GetStringArray("only")
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("can not get flag 'only': %w", err)
	}

	excludeCmds, err := cmd.Parent().Flags().GetStringArray("exclude")
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("can not get flag 'exclude': %w", err)
	}

	if len(excludeCmds) > 0 && len(onlyCmds) > 0 {
		return []string{}, []string{}, fmt.Errorf(
			"you must use either 'only' or 'exclude' flag but not both at the same time")
	}

	return onlyCmds, excludeCmds, nil
}

func parseEnvFlag(cmd *cobra.Command) (map[string]string, error) {
	// TraversChildren enabled for parent so we will have parent flags here
	envs, err := cmd.Parent().Flags().GetStringToString("env")
	if err != nil {
		return map[string]string{}, fmt.Errorf("can not get flag 'env': %w", err)
	}

	return envs, nil
}
