package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lets-cli/lets/internal/checksum"
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
	t.Run("default group", func(t *testing.T) {
		text := dedent.Dedent(`
		cmd: [echo, Hello]
		`)
		command := CommandFixture(t, text)
		exp := "Common"

		if command.GroupName != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, command.GroupName)
		}
	})

	t.Run("provided custom group", func(t *testing.T) {
		text := dedent.Dedent(`
		group: Group Name
		cmd: [echo, Hello]
		`)
		command := CommandFixture(t, text)
		exp := "Group Name"

		if command.GroupName != exp {
			t.Errorf("wrong output. \nexpect %s \ngot:  %s", exp, command.GroupName)
		}
	})
}

func TestParseCommandChecksum(t *testing.T) {
	t.Run("old list syntax", func(t *testing.T) {
		text := dedent.Dedent(`
		checksum:
		  - foo.txt
		persist_checksum: true
		cmd: echo ok
		`)
		command := CommandFixture(t, text)

		if !command.PersistChecksum {
			t.Fatal("expected persisted checksum")
		}

		got := command.ChecksumSources[checksum.DefaultChecksumKey]
		if len(got) != 1 || got[0] != "foo.txt" {
			t.Fatalf("unexpected checksum sources: %v", got)
		}
	})

	t.Run("new files syntax", func(t *testing.T) {
		text := dedent.Dedent(`
		checksum:
		  files:
		    source:
		      - foo.txt
		  persist: true
		cmd: echo ok
		`)
		command := CommandFixture(t, text)

		if !command.PersistChecksum {
			t.Fatal("expected persisted checksum")
		}

		got := command.ChecksumSources["source"]
		if len(got) != 1 || got[0] != "foo.txt" {
			t.Fatalf("unexpected checksum sources: %v", got)
		}
	})

	t.Run("new files list syntax", func(t *testing.T) {
		text := dedent.Dedent(`
		checksum:
		  files:
		    - foo.txt
		cmd: echo ok
		`)
		command := CommandFixture(t, text)

		got := command.ChecksumSources[checksum.DefaultChecksumKey]
		if len(got) != 1 || got[0] != "foo.txt" {
			t.Fatalf("unexpected checksum sources: %v", got)
		}
	})

	t.Run("new sh syntax", func(t *testing.T) {
		text := dedent.Dedent(`
		checksum:
		  sh: echo 1234
		  persist: true
		cmd: echo ok
		`)
		command := CommandFixture(t, text)

		if command.ChecksumSh != "echo 1234" {
			t.Fatalf("unexpected checksum sh: %s", command.ChecksumSh)
		}

		if !command.PersistChecksum {
			t.Fatal("expected persisted checksum")
		}
	})

	t.Run("invalid configurations", func(t *testing.T) {
		tests := []struct {
			name    string
			text    string
			wantErr string
		}{
			{
				name: "both checksum files and sh set",
				text: dedent.Dedent(`
				checksum:
				  files:
				    - foo.txt
				  sh: echo 1234
				cmd: echo ok
				`),
				wantErr: "checksum must use only one of 'files' or 'sh'",
			},
			{
				name: "persist checksum without checksum",
				text: dedent.Dedent(`
				persist_checksum: true
				cmd: echo ok
				`),
				wantErr: "'persist_checksum' must be used with 'checksum'",
			},
			{
				name: "checksum persist without files or sh",
				text: dedent.Dedent(`
				checksum:
				  persist: true
				cmd: echo ok
				`),
				wantErr: "'persist_checksum' or 'checksum.persist' must be used with 'checksum.files' or 'checksum.sh'",
			},
			{
				name: "conflicting persist settings",
				text: dedent.Dedent(`
				persist_checksum: true
				checksum:
				  files:
				    - foo.txt
				  persist: false
				cmd: echo ok
				`),
				wantErr: "'persist_checksum' conflicts with 'checksum.persist'",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := bytes.NewBufferString(tt.text)
				command := &Command{}

				err := yaml.NewDecoder(buf).Decode(&command)
				if err == nil {
					t.Fatal("expected command fixture decode error")
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error to contain %q, got %q", tt.wantErr, err.Error())
				}
			})
		}
	})
}
