package cmd

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"

	"github.com/lets-cli/lets/config/config"
	"github.com/lets-cli/lets/docopt"
	"github.com/spf13/cobra"
)

const zshCompletionText = `#compdef _lets lets

LETS_EXECUTABLE=lets

function _lets {
    local state

	_arguments -C -s \
		"1: :->cmds" \
		'*::arg:->args'

	case $state in
		cmds)
			_lets_commands
			;;
		args)
			_lets_command_options "${words[1]}"
			;;
	esac
}

# Check if in folder with correct lets.yaml file
_check_lets_config() {
	${LETS_EXECUTABLE} 1>/dev/null 2>/dev/null
	echo $?
}

_lets_commands () {
	local cmds

	if [ $(_check_lets_config) -eq 0 ]; then
		IFS=$'\n' cmds=($(${LETS_EXECUTABLE} completion --commands --verbose))
	else
		cmds=()
	fi
	_describe -t commands 'Available commands' cmds
}

_lets_command_options () {
	local cmd=$1

	if [ $(_check_lets_config) -eq 0 ]; then
		IFS=$'\n'
		_arguments -s $(${LETS_EXECUTABLE} completion --options=${cmd} --verbose)
	fi
}
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
func genBashCompletion(out io.Writer) error {
	tmpl, err := template.New("Main").Parse(bashCompletionText)
	if err != nil {
		return fmt.Errorf("error creating zsh completion template: %w", err)
	}

	return tmpl.Execute(out, nil)
}

// generate zsh completion script.
// if verbose passed - generate completion with description.
func genZshCompletion(out io.Writer, verbose bool) error {
	tmpl, err := template.New("Main").Parse(zshCompletionText)
	if err != nil {
		return fmt.Errorf("error creating zsh completion template: %w", err)
	}

	data := struct {
		Verbose string
	}{Verbose: ""}

	if verbose {
		data.Verbose = "--verbose"
	}

	return tmpl.Execute(out, data)
}

// generate string of commands joined with \n.
func getCommandsList(rootCmd *cobra.Command, out io.Writer, verbose bool) error {
	buf := new(bytes.Buffer)

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

	_, err := buf.WriteTo(out)
	if err != nil {
		return fmt.Errorf("can not generate commands list: %w", err)
	}

	return nil
}

type option struct {
	name string
	desc string
}

// generate string of command options joined with \n.
func getCommandOptions(command *config.Command, out io.Writer, verbose bool) error {
	if command.Docopts == "" {
		return nil
	}

	rawOpts, err := docopt.ParseOptions(command.Docopts, command.Name)
	if err != nil {
		return fmt.Errorf("can not parse docopts: %w", err)
	}

	var options []option

	for _, opt := range rawOpts {
		if strings.HasPrefix(opt.Name, "--") {
			options = append(options, option{name: opt.Name, desc: opt.Description})
		}
	}

	sort.SliceStable(options, func(i, j int) bool {
		return options[i].name < options[j].name
	})

	buf := new(bytes.Buffer)

	for _, option := range options {
		if verbose {
			desc := fmt.Sprintf("No description for option %s", option.name)

			if option.desc != "" {
				desc = strings.TrimSpace(option.desc)
			}

			buf.WriteString(fmt.Sprintf("%[1]s[%s]\n", option.name, desc))
		} else {
			buf.WriteString(fmt.Sprintf("%s\n", option.name))
		}
	}

	_, err = buf.WriteTo(out)
	if err != nil {
		return fmt.Errorf("can not generate command options list: %w", err)
	}

	return nil
}

func initCompletionCmd(rootCmd *cobra.Command, cfg *config.Config) {
	completionCmd := &cobra.Command{
		Use:    "completion",
		Hidden: true,
		Short:  "Generates completion scripts for bash, zsh",
		RunE: func(cmd *cobra.Command, args []string) error {
			shellType, err := cmd.Flags().GetString("shell")
			if err != nil {
				return fmt.Errorf("can not get flag 'shell': %w", err)
			}

			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return fmt.Errorf("can not get flag 'verbose': %w", err)
			}

			list, err := cmd.Flags().GetBool("list")
			if err != nil {
				return fmt.Errorf("can not get flag 'list': %w", err)
			}

			commands, err := cmd.Flags().GetBool("commands")
			if err != nil {
				return fmt.Errorf("can not get flag 'commands': %w", err)
			}

			if list {
				commands = true
			}

			optionsForCmd, err := cmd.Flags().GetString("options")
			if err != nil {
				return fmt.Errorf("can not get flag 'options': %w", err)
			}

			if optionsForCmd != "" {
				if cfg == nil {
					return fmt.Errorf("can not read config")
				}

				command, exists := cfg.Commands[optionsForCmd]
				if !exists {
					return fmt.Errorf("command %s not declared in config", optionsForCmd)
				}

				return getCommandOptions(command, cmd.OutOrStdout(), verbose)
			}

			if commands {
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

	completionCmd.Flags().StringP("shell", "s", "", "The type of shell (bash or zsh)")
	completionCmd.Flags().Bool("list", false, "Show list of commands [deprecated, use --commands]")
	completionCmd.Flags().Bool("commands", false, "Show list of commands")
	completionCmd.Flags().String("options", "", "Show list of options for command")
	completionCmd.Flags().Bool("verbose", false, "Verbose list of commands or options (with description) (only for zsh)")

	rootCmd.AddCommand(completionCmd)
}
