package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"strings"
	"text/template"
)

const zshCompletionText = `#compdef lets

_list () {
	local cmds

	# Check if in folder with correct lets.yaml file
	lets 1>/dev/null 2>/dev/null
	if [ $? -eq 0 ]; then
		IFS=$'\n' cmds=($(lets completion --list {{.Verbose}}))
	else
		cmds=()
	fi
	_describe -t commands 'Available commands' cmds
}

_arguments -C -s "1: :{_list}" '*::arg:->args' --
`

const bashCompletionText = `_lets_completion() {
    cur="${COMP_WORDS[COMP_CWORD]}"
    COMPREPLY=( $(lets completion --list "${COMP_WORDS[@]:1:$((COMP_CWORD-1))}" -- ${cur} 2>/dev/null) )
    if [[ ${COMPREPLY} == "" ]]; then
        COMPREPLY=( $(compgen -f -- ${cur}) )
    fi
    return 0
}

complete -o filenames -F _lets_completion lets
`

// generate bash completion script.
func genBashCompletion(w io.Writer) error {
	tmpl, err := template.New("Main").Parse(bashCompletionText)
	if err != nil {
		return fmt.Errorf("error creating zsh completion template: %v", err)
	}
	return tmpl.Execute(w, nil)
}

// generate zsh completion script.
// if verbose passed - generate completion with description
func genZshCompletion(w io.Writer, verbose bool) error {
	tmpl, err := template.New("Main").Parse(zshCompletionText)
	if err != nil {
		return fmt.Errorf("error creating zsh completion template: %v", err)
	}
	data := struct {
		Verbose string
	}{Verbose: ""}
	if verbose {
		data.Verbose = "--verbose"
	}
	return tmpl.Execute(w, data)
}

// generate string of commands joined with \n
func getCommandsList(rootCmd *cobra.Command, w io.Writer, verbose bool) error {
	var buf = new(bytes.Buffer)
	for _, cmd := range rootCmd.Commands() {
		if !cmd.Hidden && cmd.Name() != "help" {
			if verbose {
				descr := fmt.Sprintf("No description for command %s", cmd.Name())
				if cmd.Short != "" {
					descr = cmd.Short
					descr = strings.TrimSpace(descr)
				}
				buf.WriteString(fmt.Sprintf("%s:%s\n", cmd.Name(), descr))
			} else {
				buf.WriteString(fmt.Sprintf("%s\n", cmd.Name()))
			}
		}
	}
	_, err := buf.WriteTo(w)
	return err
}

func initCompletionCmd(rootCmd *cobra.Command) {
	var completionCmd = &cobra.Command{
		Use:    "completion",
		Hidden: true,
		Short:  "Generates completion scripts for bash, zsh",
		RunE: func(cmd *cobra.Command, args []string) error {
			shellType, err := cmd.Flags().GetString("shell")
			if err != nil {
				return err
			}
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}
			list, err := cmd.Flags().GetBool("list")
			if err != nil {
				return err
			}

			if list {
				return getCommandsList(rootCmd, cmd.OutOrStdout(), verbose)
			}

			if shellType == "" {
				return cmd.Help()
			}

			switch shellType {
			case "bash":
				return genBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return genZshCompletion(cmd.OutOrStdout(), verbose)
			default:
				return fmt.Errorf("unsupported shell type %q", shellType)
			}
		},
	}
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringP("shell", "s", "", "The type of shell (bash or zsh)")
	completionCmd.Flags().Bool("list", false, "Show list of commands")
	completionCmd.Flags().Bool("verbose", false, "Verbose list of commands (with description) (only for zsh)")
}
