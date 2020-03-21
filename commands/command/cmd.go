package command

import (
	"fmt"
	"os"
	"strings"
)

func stringPartition(s, sep string) (string, string, string) {
	sepPos := strings.Index(s, sep)
	if sepPos == -1 { // no separator found
		return s, "", ""
	}

	split := strings.SplitN(s, sep, 2)

	return split[0], sep, split[1]
}

// e.g if value is a json --key='{"value": 1}'
// it will wrap value in '' - --key=''{"value": 1}''
// and when escaped it become --key='{"value": 1}'
func escapeFlagValue(str string) string {
	if strings.Contains(str, "=") {
		key, sep, val := stringPartition(str, "=")
		str = strings.Join([]string{key, fmt.Sprintf("'%s'", val)}, sep)
	}

	return str
}

func parseAndValidateCmd(cmd interface{}, newCmd *Command) error {
	switch cmd := cmd.(type) {
	case string:
		newCmd.Cmd = cmd
	case []interface{}:
		cmdList := make([]string, len(cmd))

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
		// cut binary path and command name
		if len(os.Args) > 1 {
			cmdList = append(cmdList, os.Args[2:]...)
		} else if len(os.Args) == 1 {
			cmdList = append(cmdList, os.Args[1:]...)
		}

		var escapedCmdList []string
		for _, val := range cmdList {
			escapedCmdList = append(escapedCmdList, escapeFlagValue(val))
		}

		newCmd.Cmd = strings.TrimSpace(strings.Join(escapedCmdList, " "))
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
