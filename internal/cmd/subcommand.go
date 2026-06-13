package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/lets-cli/lets/internal/config/config"
	"github.com/lets-cli/lets/internal/docopt"
	"github.com/lets-cli/lets/internal/executor"
	"github.com/spf13/cobra"
)

const (
	shortLimit             = 120
	annotationSubGroupName = "SubGroupName"
	annotationHelpOptions  = "lets.helpOptions"
	annotationHelpUsage    = "lets.helpUsage"
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
	if before, _, ok := strings.Cut(text, "\n"); ok {
		return before
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

func replaceDocoptNamePlaceholder(docopts string, cmdName string) string {
	docopts = strings.Replace(docopts, "${LETS_COMMAND_NAME}", cmdName, 1)
	docopts = strings.Replace(docopts, "$LETS_COMMAND_NAME", cmdName, 1)

	return docopts
}

// Replace command name placeholder if present
// E.g. if command name is foo, lets ${LETS_COMMAND_NAME} will be lets foo.
func setDocoptNamePlaceholder(c *config.Command) {
	c.Docopts = replaceDocoptNamePlaceholder(c.Docopts, c.Name)
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

func hasHelpArg(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}

	return false
}

func buildCommandUse(commandName string, usage string) string {
	lines := make([]string, 0)

	for line := range strings.SplitSeq(usage, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = strings.TrimPrefix(line, "lets ")
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return commandName
	}

	return strings.Join(lines, "\n")
}

func buildCommandAnnotations(command *config.Command, docopts string) map[string]string {
	annotations := map[string]string{
		annotationSubGroupName: command.GroupName,
	}

	docoptParts := docopt.ParseDocoptParts(docopts)
	if usage := buildCommandUse(command.Name, docoptParts.Usage); usage != "" {
		annotations[annotationHelpUsage] = usage
	}

	helpOptions := docopt.ParseHelpOptions(docopts, command.Name)
	if len(helpOptions) == 0 {
		return annotations
	}

	payload, err := json.Marshal(helpOptions)
	if err != nil {
		return annotations
	}

	annotations[annotationHelpOptions] = string(payload)

	return annotations
}

// newSubcommand creates new cobra root subcommand from config.Command.
func newSubcommand(command *config.Command, conf *config.Config, showAll bool, out io.Writer) *cobra.Command {
	docopts := replaceDocoptNamePlaceholder(command.Docopts, command.Name)
	docoptParts := docopt.ParseDocoptParts(docopts)

	subCmd := &cobra.Command{
		Use:     command.Name,
		Example: docoptParts.Example,
		Short:   short(command.Description),
		Long:    command.Description,
		GroupID: "main",
		Hidden:  isHidden(command.Name, showAll),
		RunE: func(cmd *cobra.Command, args []string) error {
			if hasHelpArg(args) {
				return cmd.Help()
			}

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

	subCmd.Annotations = buildCommandAnnotations(command, docopts)

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
		return []string{}, []string{}, errors.New("can not use '--only' and '--exclude' at the same time")
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
