package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/executor"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

const (
	shortLimit = 120
)

// cut all elements after command name.
// [/bin/lets foo -x] will be [-x].
func prepareArgs(cmdName string, osArgs []string) []string {
	nameIdx := 0

	for idx, arg := range osArgs {
		if arg == cmdName {
			nameIdx = idx + 1
		}
	}

	return osArgs[nameIdx:]
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

func validateOnlyAndExclude(
	command *config.Command,
	only []string,
	exclude []string,
) error {
	commandNames := make([]string, len(command.Cmds.Commands))
	for idx, cmd := range command.Cmds.Commands {
		commandNames[idx] = cmd.Name
	}

	for _, name := range only {
		if !slices.Contains(commandNames, name) {
			return fmt.Errorf("no such cmd '%s' in command '%s' used in 'only' flag", name, command.Name)
		}
	}

	for _, name := range exclude {
		if !slices.Contains(commandNames, name) {
			return fmt.Errorf("no such cmd '%s' in command '%s' used in 'exclude' flag", name, command.Name)
		}
	}

	return nil
}

// Filter cmmds based on --only and --exclude values.
// Only and Exclude can not be both true at the same time.
func filterCmds(
	cmds config.Cmds,
	only []string,
	exclude []string,
) []*config.Cmd {
	hasOnly := len(only) > 0
	hasExclude := len(exclude) > 0

	if !hasOnly && !hasExclude {
		return cmds.Commands
	}

	filteredCmds := make([]*config.Cmd, 0)

	if hasOnly {
		// put only commands which in `only` list
		for _, cmd := range cmds.Commands {
			if slices.Contains(only, cmd.Name) {
				filteredCmds = append(filteredCmds, cmd)
			}
		}
	} else if hasExclude {
		// delete all commands which in `exclude` list
		for _, cmd := range cmds.Commands {
			if !slices.Contains(exclude, cmd.Name) {
				filteredCmds = append(filteredCmds, cmd)
			}
		}
	}

	return filteredCmds
}

// Replace command name placeholder if present
// E.g. if command name is foo, lets ${LETS_COMMAND_NAME} will be lets foo.
func setDocoptNamePlaceholder(c *config.Command) {
	c.Docopts = strings.Replace(c.Docopts, "${LETS_COMMAND_NAME}", c.Name, 1)
	c.Docopts = strings.Replace(c.Docopts, "$LETS_COMMAND_NAME", c.Name, 1)
}

type cmdFlags struct {
	only      []string
	exclude   []string
	env       map[string]string
	noDepends bool
}

func parseFlags(cmd *cobra.Command) (*cmdFlags, error) {
	flags := &cmdFlags{}

	only, exclude, err := parseOnlyAndExclude(cmd)
	if err != nil {
		return nil, err
	}

	flags.only = only
	flags.exclude = exclude

	// env from -E flag
	env, err := parseEnvFlag(cmd)
	if err != nil {
		return nil, err
	}

	flags.env = env

	noDepends, err := parseNoDepends(cmd)
	if err != nil {
		return nil, err
	}

	flags.noDepends = noDepends

	return flags, nil
}

func isHidden(cmdName string, showAll bool) bool {
	if strings.HasPrefix(cmdName, "_") {
		return !showAll
	}

	return false
}

// newSubcommand creates new cobra root subcommand from config.Command.
func newSubcommand(command *config.Command, conf *config.Config, showAll bool, out io.Writer) *cobra.Command {
	subCmd := &cobra.Command{
		Use:    command.Name,
		Short:  short(command.Description),
		Hidden: isHidden(command.Name, showAll),
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Args = append(command.Args, prepareArgs(command.Name, os.Args)...)
			command.Cmds.AppendArgs(args)

			flags, err := parseFlags(cmd)
			if err != nil {
				return err
			}

			if err := validateOnlyAndExclude(command, flags.only, flags.exclude); err != nil {
				return err
			}

			command.Env.MergeMap(flags.env)

			setDocoptNamePlaceholder(command)

			command.Cmds.Commands = filterCmds(command.Cmds, flags.only, flags.exclude)

			if flags.noDepends {
				command = command.Clone()
				command.Depends = &config.Deps{}
			}

			ctx := executor.NewExecutorCtx(cmd.Context(), command)
			return executor.NewExecutor(conf, out).Execute(ctx)
		},
		// we use docopt to parse flags on our own, so any flag is valid flag here
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Args:               cobra.ArbitraryArgs,
		// disables builtin --help flag
		DisableFlagParsing: true,
		// print help message manyally
		SilenceUsage: true,
	}

	subCmd.SetHelpFunc(func(c *cobra.Command, strings []string) {
		if _, err := fmt.Fprint(c.OutOrStdout(), command.Help()); err != nil {
			c.Println(err)
		}
	})

	return subCmd
}

// initialize all commands dynamically from config.
func InitSubCommands(rootCmd *cobra.Command, conf *config.Config, showAll bool, out io.Writer) {
	for _, cmdToRun := range conf.Commands {
		rootCmd.AddCommand(newSubcommand(cmdToRun, conf, showAll, out))
	}
}

func parseOnlyAndExclude(cmd *cobra.Command) ([]string, []string, error) {
	only, err := cmd.Parent().Flags().GetStringArray("only")
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("can not get flag 'only': %w", err)
	}

	exclude, err := cmd.Parent().Flags().GetStringArray("exclude")
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("can not get flag 'exclude': %w", err)
	}

	if len(exclude) > 0 && len(only) > 0 {
		return []string{}, []string{}, fmt.Errorf(
			"can not use '--only' and '--exclude' at the same time")
	}

	return only, exclude, nil
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
