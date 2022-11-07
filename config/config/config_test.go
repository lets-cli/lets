package config

import (
	"bytes"
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
		cfg.InitArgs(args)

		exp := `echo Hello 'World' '--foo' '--bar='{"age": 20}''`
		if script := cfg.Commands["hello"].Cmds.Commands[0].Script; script != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, script)
		}
	})

}
