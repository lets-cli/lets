package util

import (
	"path/filepath"
	"testing"
)

func TestLetsUserDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	dir, err := LetsUserDir()
	if err != nil {
		t.Fatalf("LetsUserDir() error = %v", err)
	}

	expected := filepath.Join(tmpDir, ".config", "lets")
	if dir != expected {
		t.Fatalf("expected %q, got %q", expected, dir)
	}
}

func TestLetsUserFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	path, err := LetsUserFile("state.yaml")
	if err != nil {
		t.Fatalf("LetsUserFile() error = %v", err)
	}

	expected := filepath.Join(tmpDir, ".config", "lets", "state.yaml")
	if path != expected {
		t.Fatalf("expected %q, got %q", expected, path)
	}
}
