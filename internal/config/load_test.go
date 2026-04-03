package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("just load config", func(t *testing.T) {
		_, err := Load("", "", "0.0.0-test")
		if err != nil {
			t.Errorf("can not load test config: %s", err)
		}
	})

	t.Run("returns error for malformed local mixin", func(t *testing.T) {
		tempDir := t.TempDir()

		mainConfig := "shell: bash\nmixins: [mixin.yaml]\ncommands:\n  ok:\n    cmd: echo ok\n"
		if err := os.WriteFile(filepath.Join(tempDir, "lets.yaml"), []byte(mainConfig), 0o644); err != nil {
			t.Fatalf("write main config: %v", err)
		}

		mixinConfig := "commands:\n  test1:\n    xxx\n    cmd: echo Test\n"
		if err := os.WriteFile(filepath.Join(tempDir, "mixin.yaml"), []byte(mixinConfig), 0o644); err != nil {
			t.Fatalf("write mixin config: %v", err)
		}

		_, err := Load("", tempDir, "0.0.0-test")
		if err == nil {
			t.Fatal("expected malformed mixin error")
		}

		if !strings.Contains(err.Error(), "failed to read mixin config 'mixin.yaml'") {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(err.Error(), "can not parse mixin config mixin.yaml") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
