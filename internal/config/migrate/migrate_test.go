package migrate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFixFilePreservesExistingFileMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lets.yaml")
	originalMode := os.FileMode(0o600)

	err := os.WriteFile(path, []byte(`
shell: bash
commands:
  build:
    persist_checksum: true
    checksum:
      - go.mod
    cmd: go build
`), originalMode)
	if err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	changed, _, _, err := fixFile(path, false, []Migration{ChecksumMigration{}})
	if err != nil {
		t.Fatalf("fix file: %v", err)
	}
	if !changed {
		t.Fatal("expected file to change")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat migrated file: %v", err)
	}

	if got := info.Mode().Perm(); got != originalMode {
		t.Fatalf("expected mode %s, got %s", originalMode, got)
	}
}
