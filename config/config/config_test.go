package config

import (
	"bytes"
	"maps"
	"testing"

	"github.com/lithammer/dedent"
	"gopkg.in/yaml.v3"
)

func ConfigFixture(t *testing.T, text string) *Config {
	buf := bytes.NewBufferString(text)
	c := NewConfig(".", ".", ".")
	if err := yaml.NewDecoder(buf).Decode(&c); err != nil {
		t.Fatalf("config fixture decode error: %s", err)
	}

	return c
}

func TestParseConfig(t *testing.T) {
	t.Run("append args to cmd as list", func(t *testing.T) {
		args := []string{"World", "--foo", `--bar='{"age": 20}'`}
		text := dedent.Dedent(`
		shell: bash
		commands:
		  hello:
		    cmd: [echo, Hello]
		`)
		cfg := ConfigFixture(t, text)
		cmd := cfg.Commands["hello"]
		cmd.Cmds.AppendArgs(args)

		exp := `echo Hello 'World' '--foo' '--bar='{"age": 20}''`
		if script := cmd.Cmds.Commands[0].Script; script != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, script)
		}
	})

	t.Run("parse env with alias", func(t *testing.T) {
		text := dedent.Dedent(`
		shell: bash

		x-default-env: &default-env
		  HELLO: WORLD

		env:
		  <<: *default-env
		  FOO: BAR

		commands:
		  hello:
		    cmd: [echo, Hello]
		`)
		cfg := ConfigFixture(t, text)

		env := cfg.Env.Dump()
		expected := map[string]string{
			"FOO":   "BAR",
			"HELLO": "WORLD",
		}
		if !maps.Equal(env, expected) {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", expected, env)
		}
	})

	t.Run("invalid alias name - does not start with x-", func(t *testing.T) {
		text := dedent.Dedent(`
		shell: bash

		default-env: &default-env
		  HELLO: WORLD

		env:
		  <<: *default-env
		  FOO: BAR

		commands:
		  hello:
		    cmd: [echo, Hello]
		`)

		buf := bytes.NewBufferString(text)
		c := NewConfig(".", ".", ".")
		err := yaml.NewDecoder(buf).Decode(&c)
		if err.Error() != "keyword 'default-env' not supported" {
			t.Errorf("config must not allow custom keywords")
		}
	})
}
