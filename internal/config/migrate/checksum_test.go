package migrate

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func applyChecksumMigration(t *testing.T, input string) string {
	t.Helper()

	root := &yaml.Node{}
	if err := yaml.Unmarshal([]byte(input), root); err != nil {
		t.Fatalf("decode yaml: %s", err)
	}

	changed, err := ChecksumMigration{}.Apply(root)
	if err != nil {
		t.Fatalf("apply migration: %s", err)
	}

	if !changed {
		t.Fatal("expected migration to change config")
	}

	encoded, err := encodeYAML(root)
	if err != nil {
		t.Fatalf("encode yaml: %s", err)
	}

	return string(encoded)
}

func TestChecksumMigrationListWithPersistChecksum(t *testing.T) {
	got := applyChecksumMigration(t, `
shell: bash
commands:
  build:
    persist_checksum: true
    checksum:
      - go.mod
    cmd: go build
`)

	for _, want := range []string{
		"checksum:\n      files:\n        - go.mod\n      persist: true",
		"cmd: go build",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected migrated config to contain %q:\n%s", want, got)
		}
	}

	if strings.Contains(got, "persist_checksum") {
		t.Fatalf("expected persist_checksum to be removed:\n%s", got)
	}
}

func TestChecksumMigrationMapWithPersistChecksum(t *testing.T) {
	got := applyChecksumMigration(t, `
shell: bash
commands:
  build:
    persist_checksum: true
    checksum:
      deps:
        - go.mod
    cmd: go build
`)

	want := "checksum:\n      files:\n        deps:\n          - go.mod\n      persist: true"
	if !strings.Contains(got, want) {
		t.Fatalf("expected migrated config to contain %q:\n%s", want, got)
	}
}

func TestChecksumMigrationMovesPersistIntoNewChecksum(t *testing.T) {
	got := applyChecksumMigration(t, `
shell: bash
commands:
  build:
    persist_checksum: true
    checksum:
      sh: git rev-parse HEAD
    cmd: go build
`)

	want := "checksum:\n      sh: git rev-parse HEAD\n      persist: true"
	if !strings.Contains(got, want) {
		t.Fatalf("expected migrated config to contain %q:\n%s", want, got)
	}
}

func TestChecksumMigrationKeepsBlankLineBetweenCommands(t *testing.T) {
	got := applyChecksumMigration(t, `
shell: bash
commands:
  build:
    persist_checksum: true
    checksum:
      - go.mod
    cmd: go build

  test:
    persist_checksum: true
    checksum:
      - go.sum
    cmd: go test ./...
`)

	want := "    cmd: go build\n\n  test:"
	if !strings.Contains(got, want) {
		t.Fatalf("expected migrated config to keep command spacing %q:\n%s", want, got)
	}
}
