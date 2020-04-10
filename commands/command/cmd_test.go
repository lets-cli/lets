package command

import (
	"os"
	"strings"
	"testing"
)

// that's how shell does it
func simulateProcessShellArgs(inputCmdList []string) []string {
	var cmdList []string

	for _, arg := range inputCmdList {
		isEnquoted := len(arg) >= 2 && (arg[0] == '\'' && arg[len(arg)-1] == '\'')
		if isEnquoted {
			quoteless := arg[1 : len(arg)-1]
			cmdList = append(cmdList, quoteless)
		} else {
			cmdList = append(cmdList, arg)
		}
	}

	return cmdList
}

func TestCommandFieldCmd(t *testing.T) {
	t.Run("so subcommand in os.Args", func(t *testing.T) {
		testCmd := NewCommand("test-cmd")
		cmdArgs := "echo Hello"
		// mock args
		os.Args = []string{"bin_to_run"}
		err := parseAndValidateCmd(cmdArgs, &testCmd)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if testCmd.Cmd != cmdArgs {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", cmdArgs, testCmd.Cmd)
		}
	})

	t.Run("as string", func(t *testing.T) {
		testCmd := NewCommand("test-cmd")
		cmdArgs := "echo Hello"

		err := parseAndValidateCmd(cmdArgs, &testCmd)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if testCmd.Cmd != cmdArgs {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", cmdArgs, testCmd.Cmd)
		}
	})

	t.Run("as list", func(t *testing.T) {
		testCmd := NewCommand("test-cmd")
		var cmdArgs []interface{}
		cmdArgs = append(cmdArgs, "echo", "Hello")

		appendArgs := []string{"one", "two", "--there", `--four='{"age": 20}'`}
		// mock args
		os.Args = append([]string{"bin_to_run", "test-cmd"}, appendArgs...)

		err := parseAndValidateCmd(cmdArgs, &testCmd)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		exp := strings.Join(append([]string{"echo", "Hello"}, "'one'", "'two'", "'--there'", `'--four='{"age": 20}''`), " ")
		if testCmd.Cmd != exp {
			t.Errorf("wrong output. \nexpect: %s \ngot:    %s", exp, testCmd.Cmd)
		}
	})
}

func TestEscapeArguments(t *testing.T) {
	t.Run("escape value if json", func(t *testing.T) {
		jsonArg := `--kwargs={"age": 20}`
		escaped := escapeArgs([]string{jsonArg})[0]
		exp := `'--kwargs={"age": 20}'`
		if escaped != exp {
			t.Errorf("wrong output. \nexpect: %s \ngot:    %s", exp, escaped)
		}
	})

	t.Run("escape string with whitespace", func(t *testing.T) {
		letsCmd := "lets commitCrime"
		appendArgs := "-m 'azaza lalka'"
		fullCommand := strings.Join([]string{letsCmd, appendArgs}, " ")

		cmdList := simulateProcessShellArgs(strings.Split(fullCommand, " "))

		args := cmdList[2:]
		escapedArgs := escapeArgs(args)
		resultArgs := strings.Join(simulateProcessShellArgs(escapedArgs), " ")

		if resultArgs != appendArgs {
			t.Errorf("wrong output. \nexpect: %s \ngot:    %s", appendArgs, resultArgs)
		}
	})
}
