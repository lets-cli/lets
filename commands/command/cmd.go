package command

import (
	"fmt"
	"os"
	"strings"
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

func parseAndValidateCmd(cmd interface{}, newCmd *Command) error { //nolint:cyclop
	switch cmd := cmd.(type) {
	case string:
		// TODO pass args to command as is if option accepts_arguments: true
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
				return newParseCommandError(
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
	case map[interface{}]interface{}:
		cmdMap := make(map[string]string, len(cmd))

		for cmdName, cmdScript := range cmd {
			cmdName, cmdNameOk := cmdName.(string)
			if !cmdNameOk {
				return newParseCommandError(
					"cmd name must be string",
					newCmd.Name,
					CMD,
					cmdName,
				)
			}

			cmdScript, cmdScriptOK := cmdScript.(string)
			if !cmdScriptOK {
				return newParseCommandError(
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
		return newParseCommandError(
			"must be either string or list of string",
			newCmd.Name,
			CMD,
			"",
		)
	}

	return nil
}
