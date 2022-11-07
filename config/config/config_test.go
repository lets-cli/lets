package config

import (
	"bytes"
	"os"
	"testing"

	"github.com/lithammer/dedent"
	"gopkg.in/yaml.v3"
)

func ConfigFixture(t *testing.T, text string, args []string) *Config {
	buf := bytes.NewBufferString(text)
	c := NewConfig(".", ".", ".")
	os.Args = args
	if err := yaml.NewDecoder(buf).Decode(&c); err != nil {
		t.Fatalf("config fixture decode error: %s", err)
	}

	return c
}

func TestParseConfig(t *testing.T) {
	t.Run("append args to cmd as list", func(t *testing.T) {
		args := []string{"/bin/lets", "hello", "World", "--foo", `--bar='{"age": 20}'`}
		text := dedent.Dedent(`
		shell: bash
		commands:
		  hello:
		    cmd: [echo, Hello]
		`)
		cfg := ConfigFixture(t, text, args)

		exp := `echo Hello 'World' '--foo' '--bar='{"age": 20}''`
		if script := cfg.Commands["hello"].Cmds.Commands[0].Script; script != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, script)
		}
	})

}
