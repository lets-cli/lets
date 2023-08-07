package executor

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/lets-cli/lets/internal/checksum"
	"github.com/lets-cli/lets/internal/config/config"
)

func TestInitCmdUsesCommandShellAndWorkDirForChecksum(t *testing.T) {
	tempDir := t.TempDir()
	projectDir := filepath.Join(tempDir, "project")

	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("can not create project dir: %s", err)
	}

	if err := os.WriteFile(filepath.Join(projectDir, "checksum.txt"), []byte("command-checksum"), 0o600); err != nil {
		t.Fatalf("can not write checksum file: %s", err)
	}

	cfg := &config.Config{
		Shell:   "sh",
		WorkDir: tempDir,
	}
	cmd := &config.Command{
		Name:        "checksum-cmd",
		Shell:       "bash",
		WorkDir:     projectDir,
		SkipDocopts: true,
		ChecksumCmd: "[[ -f checksum.txt ]] && cat checksum.txt",
	}

	executor := NewExecutor(cfg, nil)
	ctx := NewExecutorCtx(context.Background(), cmd)

	if err := executor.initCmd(ctx); err != nil {
		t.Fatalf("initCmd failed: %s", err)
	}

	if got := cmd.ChecksumMap[checksum.DefaultChecksumKey]; got != "command-checksum" {
		t.Fatalf("wrong checksum output. Expect: %s, got: %s", "command-checksum", got)
	}
}
