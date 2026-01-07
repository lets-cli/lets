package config

import (
	"bytes"
	"testing"

	"github.com/lithammer/dedent"
	"gopkg.in/yaml.v3"
)

func CommandFixture(t *testing.T, text string) *Command {
	buf := bytes.NewBufferString(text)
	c := &Command{}
	if err := yaml.NewDecoder(buf).Decode(&c); err != nil {
		t.Fatalf("command fixture decode error: %s", err)
	}

	return c
}

func TestParseCommand(t *testing.T) {
	t.Run("default group_name", func(t *testing.T) {
		text := dedent.Dedent(`
		cmd: [echo, Hello]
		`)
		command := CommandFixture(t, text)
		exp := ""

		if command.GroupName != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, command.GroupName)
		}
	})

	t.Run("provided custom group_name", func(t *testing.T) {
		text := dedent.Dedent(`
		group_name: Group Name
		cmd: [echo, Hello]
		`)
		command := CommandFixture(t, text)
		exp := "Group Name"

		if command.GroupName != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, command.GroupName)
		}
	})
}
