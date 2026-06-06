package cmd

import (
	"fmt"
	"strings"

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

// newRootCmd creates root cobra command that represents the base command
// when called without any subcommands.
func newRootCmd(version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lets",
		Short: "A CLI task runner",
		Args:  validateCommandArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Version:            buildVersion(version, buildDate),
		Annotations:        map[string]string{"buildDate": buildDate},
		TraverseChildren:   true,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		// handle errors manually
		SilenceErrors: true,
		// print help message manyally
		SilenceUsage: true,
	}

	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		var err error
		if c == c.Root() {
			err = PrintRootHelpMessage(c)
		} else {
			err = PrintHelpMessage(c)
		}

		if err != nil {
			c.Println(err)
		}
	})
	cmd.AddGroup(&cobra.Group{ID: "main", Title: "Commands:"}, &cobra.Group{ID: "internal", Title: "Internal commands:"})
	cmd.SetHelpCommandGroupID("internal")

	return cmd
}

func buildVersion(version string, buildDate string) string {
	if buildDate != "" {
		version += fmt.Sprintf(" (%s)", buildDate)
	}

	return version
}

// CreateRootCommand used to run only root command without config.
func CreateRootCommand(version string, buildDate string) *cobra.Command {
	rootCmd := newRootCmd(version, buildDate)

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
