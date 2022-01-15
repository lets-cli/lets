package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/lets-cli/lets/config/config"
)

// A workaround function which helps to prevent breaking
// strings with special symbols (' ', '*', '$', '#'...)
// When you run a command with an argument containing one of these, you put it into quotation marks:
// lets alembic -n dev revision --autogenerate -m "revision message"
// which makes shell understand that "revision message" is a single argument, but not two args
// The problem is, lets constructs a script string
// and then passes it to an appropriate interpreter (sh -c $SCRIPT)
// so we need to wrap args with quotation marks to prevent breaking
// This also solves problem with json params: --key='{"value": 1}' => '--key={"value": 1}'.
func escapeArgs(args []string) []string {
	var escapedArgs []string

	for _, arg := range args {
		// wraps every argument with quotation marks to avoid ambiguity
		// TODO: maybe use some kind of blacklist symbols to wrap only necessary args
		escapedArg := fmt.Sprintf("'%s'", arg)
		escapedArgs = append(escapedArgs, escapedArg)
	}

	return escapedArgs
}

func parseCmd(cmd interface{}, newCmd *config.Command) error { //nolint:cyclop
	switch cmd := cmd.(type) {
	case string:
		newCmd.Cmd = cmd
	case []interface{}:
		// a list of arguments to be appended to commands in lets.yaml
		var proxyArgs []string
		// cut binary path and command name
		if len(os.Args) > 1 {
			proxyArgs = os.Args[2:]
		} else if len(os.Args) == 1 {
			proxyArgs = os.Args[1:]
		}

		cmdList := make([]string, 0, len(cmd)+len(proxyArgs))

		for _, v := range cmd {
			if v == nil {
				return parseError(
					"got nil in cmd list",
					newCmd.Name,
					CMD,
					"",
				)
			}

			cmdList = append(cmdList, fmt.Sprintf("%s", v))
		}

		fullCommandList := append(cmdList, escapeArgs(proxyArgs)...)
		newCmd.Cmd = strings.TrimSpace(strings.Join(fullCommandList, " "))
	case map[string]interface{}:
		cmdMap := make(map[string]string, len(cmd))

		for cmdName, cmdScript := range cmd {
			cmdScript, cmdScriptOK := cmdScript.(string)
			if !cmdScriptOK {
				return parseError(
					"cmd name must be string",
					newCmd.Name,
					CMD,
					cmdScript,
				)
			}

			cmdMap[cmdName] = cmdScript
		}

		newCmd.CmdMap = cmdMap
	default:
		return parseError(
			`
must be one of
  - string
  - list of string
  - map of string to string
`,
			newCmd.Name,
			CMD,
			"",
		)
	}

	return nil
}
