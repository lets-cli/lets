package config

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// that's how shell does it.
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

func CmdFixture(t *testing.T, text string, args []string) Cmds {
	buf := bytes.NewBufferString(text)
	var cmd struct {
		Cmd Cmds
	}
	os.Args = args
	if err := yaml.NewDecoder(buf).Decode(&cmd); err != nil {
		t.Fatalf("cmd fixture decode error: %s", err)
	}

	return cmd.Cmd
}

func TestCommandFieldCmd(t *testing.T) {
	t.Run("as string", func(t *testing.T) {
		cmd := CmdFixture(t, "cmd: echo Hello", []string{})
		exp := "echo Hello"
		if cmd.Commands[0].Script != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, cmd.Commands[0].Script)
		}
	})

	t.Run("as list", func(t *testing.T) {
		args := []string{"/bin/lets", "hello", "World", "--foo", `--bar='{"age": 20}'`}
		cmd := CmdFixture(t, "cmd: [echo, Hello]", args)
		exp := `echo Hello 'World' '--foo' '--bar='{"age": 20}''`
		if cmd.Commands[0].Script != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, cmd.Commands[0].Script)
		}
	})

	t.Run("as map", func(t *testing.T) {
		text := "cmd: \n  foo: echo Foo\n  bar: echo Bar"
		cmd := CmdFixture(t, text, []string{})
		expFoo := "echo Foo"
		expBar := "echo Bar"
		if cmdLen := len(cmd.Commands); cmdLen != 2 {
			t.Errorf("expect %d commands\ngot: %d", 2, cmdLen)
		}

		for _, command := range cmd.Commands {
			switch command.Name {
			case "foo":
				if command.Script != expFoo {
					t.Errorf("wrong output. \nexpect %s \ngot:  %s", expFoo, command.Script)
				}
			case "bar":
				if command.Script != expBar {
					t.Errorf("wrong output. \nexpect %s \ngot:  %s", expBar, command.Script)
				}
			default:
				t.Fatalf("unexpected command %s", command.Name)
			}
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
