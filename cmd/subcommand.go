package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/runner"
	"github.com/spf13/cobra"
)

const (
	shortLimit = 120
)

// cut all elements before command name.
func prepareArgs(cmdName string, originalArgs []string) []string {
	nameIdx := 0

	for idx, arg := range originalArgs {
		if arg == cmdName {
			nameIdx = idx
		}
	}

	return originalArgs[nameIdx:]
}

func short(text string) string {
	if idx := strings.Index(text, "\n"); idx >= 0 {
		return text[:idx]
	}

	if len(text) > shortLimit {
		return text[:shortLimit]
	}

	return text
}

// newCmdGeneric creates new cobra root sub command from Command.
func newCmdGeneric(command *config.Command, conf *config.Config, out io.Writer) *cobra.Command {
	subCmd := &cobra.Command{
		Use:   command.Name,
		Short: short(command.Description),
		RunE: func(cmd *cobra.Command, args []string) error {
			only, exclude, err := parseOnlyAndExclude(cmd)
			if err != nil {
				return err
			}

			// env from -E flag
			envs, err := parseEnvFlag(cmd)
			if err != nil {
				return err
			}

			noDepends, err := parseNoDepends(cmd)
			if err != nil {
				return err
			}

			command.Only = only
			command.Exclude = exclude
			// TODO: validate if only and exclude contains only existing commands (and also that command has Parallel true)

			command.Args = prepareArgs(command.Name, os.Args)
			command.Env.MergeMap(envs)

			// replace only one placeholder in options
			command.Docopts = strings.Replace(
				command.Docopts,
				fmt.Sprintf("${%s}", runner.GenericCmdNameTpl), command.Name,
				1,
			)

			if command.Ref != nil {
				command = conf.Commands[command.Ref.Name].FromRef(command.Ref)
			}

			if noDepends {
				command = command.Clone()
				command.Depends = &config.Deps{}
			}

			return runner.NewRunner(command, conf, out).Execute(cmd.Context())
		},
		// we use docopt to parse flags on our own, so any flag is valid flag here
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Args:               cobra.ArbitraryArgs,
		// disables builtin --help flag
		DisableFlagParsing: true,
		SilenceUsage:       true,
	}

	subCmd.SetHelpFunc(func(c *cobra.Command, strings []string) {
		if _, err := fmt.Fprint(c.OutOrStdout(), command.Help()); err != nil {
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
			"can not use '--only' and '--exclude' at the same time")
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

func parseNoDepends(cmd *cobra.Command) (bool, error) {
	noDepends, err := cmd.Parent().Flags().GetBool("no-depends")
	if err != nil {
		return false, fmt.Errorf("can not get flag 'no-depends': %w", err)
	}

	return noDepends, nil
}
