package command

import (
	"fmt"
	"os"
	"strings"
)

var (
	CMD         = "cmd"
	DESCRIPTION = "description"
	ENV         = "env"
	EVAL_ENV    = "eval_env"
	OPTIONS     = "options"
	DEPENDS     = "depends"
	CHECKSUM    = "checksum"
)

type Command struct {
	Name        string
	Cmd         string
	Description string
	Env         map[string]string
	RawOptions  string
	Options     map[string]string
	Depends     []string
	Checksum    string
}

// TODO interface{} must be replaced
func NewCommand(name string, rawCommand map[interface{}]interface{}) Command {
	newCmd := Command{
		Name: name,
		Env:  make(map[string]string),
	}

	if cmd, ok := rawCommand[CMD]; ok {
		// TODO not safe, need validation
		//  decide, validate here or top-level validate and return all errors at once
		switch cmd := cmd.(type) {
		case string:
			newCmd.Cmd = cmd
		case []interface{}:
			cmdList := make([]string, len(cmd))
			for _, v := range cmd {
				cmdList = append(cmdList, v.(string))
			}
			cmdList = append(cmdList, os.Args[1:]...)
			newCmd.Cmd = strings.Join(cmdList, " ")
		default:
			fmt.Println("default, must raise an error")
		}
		// TODO here we need to validate if cmd is an array
	}

	if desc, ok := rawCommand[DESCRIPTION]; ok {
		newCmd.Description = desc.(string)
	}

	if env, ok := rawCommand[ENV]; ok {
		// TODO dirty hacks
		for name, value := range env.(map[interface{}]interface{}) {
			newCmd.Env[name.(string)] = value.(string)
		}
	}

	if evalEnv, ok := rawCommand[EVAL_ENV]; ok {
		for name, value := range evalEnv.(map[interface{}]interface{}) {
			if computedVal, err := evalEnvVariable(value.(string)); err != nil {
				// TODO we have to fail here and log error for user
			} else {
				newCmd.Env[name.(string)] = computedVal
			}
		}
	}

	if options, ok := rawCommand[OPTIONS]; ok {
		newCmd.RawOptions = options.(string)
	}
	if depends, ok := rawCommand[DEPENDS]; ok {
		for _, value := range depends.([]interface{}) {
			// TODO validate if command is realy exists - in validate
			newCmd.Depends = append(newCmd.Depends, value.(string))
		}
	}

	if checksum, ok := rawCommand[CHECKSUM]; ok {
		if patterns, ok := checksum.([]interface{}); ok {
			var files []string
			for _, value := range patterns {
				// TODO validate if command is realy exists - in validate
				files = append(files, value.(string))
			}
			checksum, err := calculateChecksum(files)
			if err == nil {
				newCmd.Checksum = checksum
			} else {
				// TODO return error or caclulate checksum upper in the code
				fmt.Printf("error while checksum %s\n", err)
			}
		}
	}
	return newCmd
}
