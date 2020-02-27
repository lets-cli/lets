package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"text/template"
)

var tmpl = `To enable completion in your shell, run:
  eval "$(lets completion -s <shell>)"
You can add that to your '~/.bash_profile' to enable completion whenever you
start a new shell.
`

var (
	zshCompletionText = `
#compdef lets
function _lets {
  local -a commands

  _arguments -C \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=({{range .Commands}}{{if not .Hidden}}
	{{if .Short}}"{{.Name}}:{{.Short}}"{{else}}"{{.Name}}:No description for command {{.Name}}"{{end}}{{end}}{{end}}
    )
    _describe "command" commands
    ;;
  esac
}

_lets
`
)

var (
	zshCompletionSimpleText = `
#compdef lets
function _lets {
  local -a commands

  _arguments -C \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=({{range .Commands}}{{if not .Hidden}}
	"{{.Name}}{{end}}{{end}}
    )
    _describe "command" commands
    ;;
  esac
}

_lets
`
)

func genZshCompletion(rootCmd *cobra.Command, w io.Writer) error {
	tmpl, err := template.New("Main").Parse(zshCompletionText)
	if err != nil {
		return fmt.Errorf("error creating zsh completion template: %v", err)
	}
	return tmpl.Execute(w, rootCmd)
}

// same as genZshCompletion but without description
func genZshCompletionSimple(rootCmd *cobra.Command, w io.Writer) error {
	tmpl, err := template.New("Main").Parse(zshCompletionSimpleText)
	if err != nil {
		return fmt.Errorf("error creating zsh completion template: %v", err)
	}
	return tmpl.Execute(w, rootCmd)
}

func initCompletionCmd(rootCmd *cobra.Command) {
	var completionCmd = &cobra.Command{
		Use:    "completion",
		Hidden: true,
		Short:  "Generates completion scripts",
		Long:   tmpl,
		RunE: func(cmd *cobra.Command, args []string) error {
			shellType, err := cmd.Flags().GetString("shell")
			simple, err := cmd.Flags().GetBool("simple")
			if err != nil {
				return err
			}

			switch shellType {
			case "bash":
				return rootCmd.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				if simple {
					return genZshCompletionSimple(rootCmd, cmd.OutOrStdout())
				}
				return genZshCompletion(rootCmd, cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell type %q", shellType)
			}
		},
	}
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringP("shell", "s", "bash", "The type of shell")
	completionCmd.Flags().Bool("simple", false, "Simple completion or with description (only for zsh)")
}
