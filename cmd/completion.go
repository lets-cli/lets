package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"strings"
	"text/template"
)

var tmpl = `To enable completion in your shell, run:
  eval "$(lets completion -s <shell>)"
You can add that to your '~/.bash_profile' to enable completion whenever you
start a new shell.
`

const zshCompletionText = `
#compdef lets

_list () {
	local cmds

	# Check if in folder with correct lets.yaml file
	lets 1>/dev/null 2>/dev/null
	if [ $? -eq 0 ]; then
		IFS=$'\n' cmds=($(lets completion --list {{.Short}}))
	else
		cmds=()
	fi
	_describe -t commands 'Available commands' cmds
}

_arguments -C -s "1: :{_list}" '*::arg:->args' --
`

// generate zsh completion script.
// if short passed - generate completion without description
func genZshCompletion(w io.Writer, short bool) error {
	tmpl, err := template.New("Main").Parse(zshCompletionText)
	if err != nil {
		return fmt.Errorf("error creating zsh completion template: %v", err)
	}
	data := struct {
		Short string
	}{Short: ""}
	if short {
		data.Short = "--short"
	}
	return tmpl.Execute(w, data)
}

// generate string of commands joined with \n
func getCommandsList(rootCmd *cobra.Command, w io.Writer, short bool) error {
	var buf = new(bytes.Buffer)
	for _, cmd := range rootCmd.Commands() {
		if !cmd.Hidden && cmd.Name() != "help" {
			if short {
				buf.WriteString(fmt.Sprintf("%s\n", cmd.Name()))
			} else {
				descr := fmt.Sprintf("No description for command %s", cmd.Name())
				if cmd.Short != "" {
					descr = cmd.Short
					descr = strings.ReplaceAll(descr, ":", " ")
					descr = strings.TrimSpace(descr)
				}
				buf.WriteString(fmt.Sprintf("%s:%s\n", cmd.Name(), descr))
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
		Short:  "Generates completion scripts",
		Long:   tmpl,
		RunE: func(cmd *cobra.Command, args []string) error {
			shellType, err := cmd.Flags().GetString("shell")
			if err != nil {
				return err
			}
			short, err := cmd.Flags().GetBool("short")
			if err != nil {
				return err
			}
			list, err := cmd.Flags().GetBool("list")
			if err != nil {
				return err
			}

			if list {
				return getCommandsList(rootCmd, cmd.OutOrStdout(), short)
			}

			switch shellType {
			case "bash":
				return rootCmd.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return genZshCompletion(cmd.OutOrStdout(), short)
			default:
				return fmt.Errorf("unsupported shell type %q", shellType)
			}
		},
	}
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringP("shell", "s", "bash", "The type of shell")
	completionCmd.Flags().Bool("list",  false, "Show list of commands")
	completionCmd.Flags().Bool("short", false, "Short completion without description (only for zsh)")
}
