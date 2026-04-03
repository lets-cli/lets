package cmd

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/lets-cli/lets/internal/set"
	"github.com/spf13/cobra"
)

type unknownCommandError struct {
	message string
}

func (e *unknownCommandError) Error() string {
	return e.message
}

func (e *unknownCommandError) ExitCode() int {
	return 2
}

func buildUnknownCommandMessage(cmd *cobra.Command, arg string) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "unknown command %q for %q", arg, cmd.CommandPath())

	if cmd.DisableSuggestions {
		return builder.String()
	}

	if cmd.SuggestionsMinimumDistance <= 0 {
		cmd.SuggestionsMinimumDistance = 2
	}

	suggestions := cmd.SuggestionsFor(arg)
	if len(suggestions) == 0 {
		return builder.String()
	}

	builder.WriteString("\n\nDid you mean this?\n")

	for _, suggestion := range suggestions {
		builder.WriteByte('\t')
		builder.WriteString(suggestion)
		builder.WriteByte('\n')
	}

	return builder.String()
}

func validateCommandArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}

	return &unknownCommandError{
		message: buildUnknownCommandMessage(cmd, args[0]),
	}
}

// newRootCmd represents the base command when called without any subcommands.
func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lets",
		Short: "A CLI task runner",
		Args:  validateCommandArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return PrintHelpMessage(cmd)
		},
		TraverseChildren:   true,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Version:            version,
		// handle errors manually
		SilenceErrors: true,
		// print help message manyally
		SilenceUsage: true,
	}
	cmd.AddGroup(&cobra.Group{ID: "main", Title: "Commands:"}, &cobra.Group{ID: "internal", Title: "Internal commands:"})
	cmd.SetHelpCommandGroupID("internal")

	return cmd
}

// CreateRootCommand used to run only root command without config.
func CreateRootCommand(version string, buildDate string) *cobra.Command {
	rootCmd := newRootCmd(version)
	rootCmd.Annotations = map[string]string{"buildDate": buildDate}

	initRootFlags(rootCmd)

	return rootCmd
}

func initRootFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringToStringP("env", "E", nil, "set env variable for running command KEY=VALUE")
	rootCmd.Flags().StringArray("only", []string{}, "run only specified command(s) described in cmd as map")
	rootCmd.Flags().StringArray("exclude", []string{}, "run all but excluded command(s) described in cmd as map")
	rootCmd.Flags().Bool("init", false, "create a new lets.yaml in the current folder")
	rootCmd.Flags().Bool("no-depends", false, "skip 'depends' for running command")
	rootCmd.Flags().CountP("debug", "d", "show debug logs (or use LETS_DEBUG=1). If used multiple times, shows more verbose logs") //nolint:lll
	rootCmd.Flags().StringP("config", "c", "", "config file (default is lets.yaml)")
	rootCmd.Flags().Bool("all", false, "show all commands (including the ones with _)")
}

func PrintHelpMessage(cmd *cobra.Command) error {
	help := cmd.UsageString()
	help = fmt.Sprintf("%s\n\n%s", cmd.Short, help)
	help = strings.Replace(help, "lets [command] --help", "lets help [command]", 1)
	_, err := fmt.Fprint(cmd.OutOrStdout(), help)

	return err
}

func maxCommandNameLen(cmd *cobra.Command) int {
	commands := cmd.Commands()
	if len(commands) == 0 {
		return 0
	}

	maxCmd := slices.MaxFunc(commands, func(a, b *cobra.Command) int {
		return cmp.Compare(len(a.Name()), len(b.Name()))
	})

	return len(maxCmd.Name())
}

func rpad(s string, padding int) string {
	formattedString := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(formattedString, s)
}

func hasSubgroup(cmd *cobra.Command) bool {
	subgroups := make(map[string]struct{})

	for _, c := range cmd.Commands() {
		if subgroup, ok := c.Annotations["SubGroupName"]; ok && subgroup != "" {
			subgroups[subgroup] = struct{}{}
			if len(subgroups) > 1 {
				return true
			}
		}
	}

	return false
}

func writeGroupCommandHelpLine(builder *strings.Builder, prefix string, name string, padding int, suffix string, short string) {
	builder.WriteString("  ")
	builder.WriteString(prefix)
	builder.WriteString(rpad(name, padding))
	builder.WriteString(suffix)
	builder.WriteString("  ")
	builder.WriteString(short)
	builder.WriteByte('\n')
}

func buildGroupCommandHelp(cmd *cobra.Command, group *cobra.Group) string {
	cmds := []*cobra.Command{}

	// select commands that belong to the specified group
	for _, c := range cmd.Commands() {
		if c.GroupID == group.ID && (c.IsAvailableCommand() || c.Name() == "help") {
			cmds = append(cmds, c)
		}
	}

	padding := maxCommandNameLen(cmd)

	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Name() < cmds[j].Name()
	})

	// Create a list of subgroups
	subGroupNameSet := set.NewSet[string]()

	for _, c := range cmds {
		if subgroup, ok := c.Annotations["SubGroupName"]; ok && subgroup != "" {
			subGroupNameSet.Add(subgroup)
		}
	}

	subGroupNameList := subGroupNameSet.ToList()
	sort.Strings(subGroupNameList)

	// generate output
	var builder strings.Builder
	builder.WriteString(group.Title)
	builder.WriteByte('\n')

	intend := ""
	if hasSubgroup(cmd) {
		intend = "  "
	}

	for _, subgroupName := range subGroupNameList {
		if len(subGroupNameList) > 1 {
			builder.WriteString("\n  ")
			builder.WriteString(subgroupName)
			builder.WriteByte('\n')
		}

		for _, c := range cmds {
			if subgroup, ok := c.Annotations["SubGroupName"]; ok && subgroup == subgroupName {
				writeGroupCommandHelpLine(&builder, intend, c.Name(), padding, "", c.Short)
			}
		}
	}

	for _, c := range cmds {
		if _, ok := c.Annotations["SubGroupName"]; !ok {
			writeGroupCommandHelpLine(&builder, "", c.Name(), padding, intend, c.Short)
		}
	}

	builder.WriteByte('\n')

	return builder.String()
}

func PrintRootHelpMessage(cmd *cobra.Command) error {
	var builder strings.Builder
	builder.WriteString(cmd.Short)
	builder.WriteString("\n\n")

	// General
	builder.WriteString("Usage:\n")

	if cmd.Runnable() {
		fmt.Fprintf(&builder, "  %s\n", cmd.UseLine())
	}

	if cmd.HasAvailableSubCommands() {
		fmt.Fprintf(&builder, "  %s [command]\n", cmd.CommandPath())
	}

	builder.WriteByte('\n')

	// Commands
	for _, group := range cmd.Groups() {
		builder.WriteString(buildGroupCommandHelp(cmd, group))
	}

	// Flags
	if cmd.HasAvailableLocalFlags() {
		builder.WriteString("Flags:\n")
		builder.WriteString(cmd.LocalFlags().FlagUsagesWrapped(120))
		builder.WriteByte('\n')
	}

	// Usage
	fmt.Fprintf(&builder, `Use "%s help [command]" for more information about a command.`, cmd.CommandPath())

	_, err := fmt.Fprint(cmd.OutOrStdout(), builder.String())

	return err
}

func PrintVersionMessage(cmd *cobra.Command) error {
	msg := "lets version " + cmd.Version
	if buildDate := cmd.Annotations["buildDate"]; buildDate != "" {
		msg += fmt.Sprintf(" (%s)", buildDate)
	}

	_, err := fmt.Fprintln(cmd.OutOrStdout(), msg)

	return err
}
