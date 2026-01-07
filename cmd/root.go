package cmd

import (
	"fmt"
	"strings"
	"sort"

	"github.com/spf13/cobra"
	"github.com/lets-cli/lets/set"
)

// newRootCmd represents the base command when called without any subcommands.
func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lets",
		Short: "A CLI task runner",
		Args:  cobra.ArbitraryArgs,
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
func CreateRootCommand(version string) *cobra.Command {
	rootCmd := newRootCmd(version)

	initRootFlags(rootCmd)

	return rootCmd
}

func initRootFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringToStringP("env", "E", nil, "set env variable for running command KEY=VALUE")
	rootCmd.Flags().StringArray("only", []string{}, "run only specified command(s) described in cmd as map")
	rootCmd.Flags().StringArray("exclude", []string{}, "run all but excluded command(s) described in cmd as map")
	rootCmd.Flags().Bool("upgrade", false, "upgrade lets to latest version")
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

func buildGroupCommandHelp(cmd *cobra.Command, group *cobra.Group) string {
	help := ""
	cmds := []*cobra.Command{}

	// select commands that belong to the specified group
	for _, c := range cmd.Commands() {
		if c.GroupID == group.ID && (c.IsAvailableCommand() || c.Name() == "help") {
			cmds = append(cmds, c)
		}
	}

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
	help += fmt.Sprintf("%s\n", group.Title)

	for _, subgroupName := range subGroupNameList {
		intend := ""
		if len(subGroupNameList) > 1 {
			help += fmt.Sprintf("\n  %s\n", subgroupName)
			intend = "  "
		}
		for _, c := range cmds {
			if subgroup, ok := c.Annotations["SubGroupName"]; ok && subgroup == subgroupName {
				help += fmt.Sprintf("%s  %-*s %s\n", intend, cmd.NamePadding(), c.Name(), c.Short)
			}
		}
	}

	for _, c := range cmds {
		if _, ok := c.Annotations["SubGroupName"]; !ok {
			help += fmt.Sprintf("  %-*s %s\n", cmd.NamePadding(), c.Name(), c.Short)
		}
	}

	help += "\n"

	return help
}


func PrintRootHelpMessage(cmd *cobra.Command) error {
	help := ""
	help = fmt.Sprintf("%s\n\n%s", cmd.Short, help)

	// General
	help += "Usage:\n"
	if cmd.Runnable() {
		help += fmt.Sprintf("  %s\n", cmd.UseLine())
	}
	if cmd.HasAvailableSubCommands() {
		help += fmt.Sprintf("  %s [command]\n", cmd.CommandPath())
	}
	help += "\n"

	// Commands
	for _, group := range cmd.Groups() {
		help += buildGroupCommandHelp(cmd, group)
	}

	// Flags
	if cmd.HasAvailableLocalFlags() {
		help += "Flags:\n"
		help += cmd.LocalFlags().FlagUsagesWrapped(120)
		help += "\n"
	}

	// Usage
	help += fmt.Sprintf(`Use "%s help [command]" for more information about a command.`, cmd.CommandPath())

	_, err := fmt.Fprint(cmd.OutOrStdout(), help)
	return err
}

func PrintVersionMessage(cmd *cobra.Command) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "lets version %s\n", cmd.Version)
	return err
}
