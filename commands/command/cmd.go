package command

import (
	"os"
	"strings"
)

func parseAndValidateCmd(cmd interface{}, newCmd *Command) error {
	switch cmd := cmd.(type) {
	case string:
		newCmd.Cmd = cmd
	case []interface{}:
		cmdList := make([]string, len(cmd))
		for _, v := range cmd {
			cmdList = append(cmdList, v.(string))
		}
		// cut binary path and command name
		cmdList = append(cmdList, os.Args[2:]...)
		newCmd.Cmd = strings.TrimSpace(strings.Join(cmdList, " "))
	default:
		return newCommandError(
			"must be either string or list of string",
			newCmd.Name,
			CMD,
			"",
		)
	}
	return nil
}
