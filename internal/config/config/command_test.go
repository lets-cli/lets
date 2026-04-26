package config

import (
	"bytes"
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
}
