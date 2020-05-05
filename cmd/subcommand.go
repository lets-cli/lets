package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/lets-cli/lets/commands/command"
	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/runner"
)

// cut all elements before command name
func prepareArgs(cmd command.Command, originalArgs []string) []string {
	nameIdx := 0

	for idx, arg := range originalArgs {
		if arg == cmd.Name {
			nameIdx = idx
		}
	}

	return originalArgs[nameIdx:]
}

// newCmdGeneric creates new cobra root sub command from Command
func newCmdGeneric(cmdToRun command.Command, conf *config.Config, out io.Writer) *cobra.Command {
	subCmd := &cobra.Command{
		Use:   cmdToRun.Name,
		Short: cmdToRun.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			only, exclude, err := parseAndValidateOnlyAndExclude(cmd)
			if err != nil {
				return err
			}

			cmdToRun.Only = only
			cmdToRun.Exclude = exclude
			cmdToRun.Args = prepareArgs(cmdToRun, os.Args)

			return runner.RunCommand(cmd.Context(), cmdToRun, conf, out)
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
		buf.WriteString(cmdToRun.RawOptions)

		_, err := buf.WriteTo(c.OutOrStdout())
		if err != nil {
			c.Println(err)
		}
	})
	initOnlyAndExecFlags(subCmd)

	return subCmd
}

// initialize all commands dynamically from config
func initSubCommands(rootCmd *cobra.Command, conf *config.Config, out io.Writer) {
	for _, cmdToRun := range conf.Commands {
		rootCmd.AddCommand(newCmdGeneric(cmdToRun, conf, out))
	}
}

func initOnlyAndExecFlags(cmd *cobra.Command) {
	cmd.Flags().StringArray("only", []string{}, "run only specified command(s) described in cmd as map")
	cmd.Flags().StringArray("exclude", []string{}, "run all but excluded command(s) described in cmd as map")
}

func parseAndValidateOnlyAndExclude(cmd *cobra.Command) (only []string, exclude []string, err error) {
	onlyCmds, err := cmd.Flags().GetStringArray("only")
	if err != nil {
		return []string{}, []string{}, err
	}

	excludeCmds, err := cmd.Flags().GetStringArray("exclude")
	if err != nil {
		return []string{}, []string{}, err
	}

	if len(excludeCmds) > 0 && len(onlyCmds) > 0 {
		return []string{}, []string{}, fmt.Errorf(
			"you must use either 'only' or 'exclude' flag but not both at the same time")
	}

	return onlyCmds, excludeCmds, nil
}
