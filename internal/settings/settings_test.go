package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
)

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	oldValue, hadValue := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset %s: %v", key, err)
	}

	t.Cleanup(func() {
		if hadValue {
			_ = os.Setenv(key, oldValue)
			return
		}

		_ = os.Unsetenv(key)
	})
}

func TestLoadFile(t *testing.T) {
	t.Run("uses defaults when file is missing", func(t *testing.T) {
		unsetEnv(t, "NO_COLOR")
		unsetEnv(t, "LETS_CHECK_UPDATE")

		cfg, err := LoadFile(filepath.Join(t.TempDir(), "missing.yaml"))
		if err != nil {
			t.Fatalf("LoadFile() error = %v", err)
		}

		if cfg.NoColor {
			t.Fatal("expected no_color default to be false")
		}
		if cfg.Theme != "default" {
			t.Fatalf("expected theme default to be %q, got %q", "default", cfg.Theme)
		}
		if !cfg.UpgradeNotify {
			t.Fatal("expected upgrade_notify default to be true")
		}
	})

	t.Run("loads file values", func(t *testing.T) {
		unsetEnv(t, "NO_COLOR")
		unsetEnv(t, "LETS_CHECK_UPDATE")

		path := filepath.Join(t.TempDir(), "config.yaml")
		err := os.WriteFile(path, []byte("no_color: true\ntheme: synthwave\nupgrade_notify: false\n"), 0o644)
		if err != nil {
			t.Fatalf("failed to write settings file: %v", err)
		}

		cfg, err := LoadFile(path)
		if err != nil {
			t.Fatalf("LoadFile() error = %v", err)
		}

		if !cfg.NoColor {
			t.Fatal("expected no_color to be true")
		}
		if cfg.Theme != "synthwave" {
			t.Fatalf("expected theme to be %q, got %q", "synthwave", cfg.Theme)
		}
		if cfg.UpgradeNotify {
			t.Fatal("expected upgrade_notify to be false")
		}
	})

	t.Run("env overrides file values", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		err := os.WriteFile(path, []byte("no_color: false\ntheme: ansi\nupgrade_notify: true\n"), 0o644)
		if err != nil {
			t.Fatalf("failed to write settings file: %v", err)
		}

		t.Setenv("NO_COLOR", "")
		t.Setenv("LETS_CHECK_UPDATE", "1")

		cfg, err := LoadFile(path)
		if err != nil {
			t.Fatalf("LoadFile() error = %v", err)
		}

		if !cfg.NoColor {
			t.Fatal("expected NO_COLOR to override settings file")
		}
		if cfg.Theme != "ansi" {
			t.Fatalf("expected theme to remain %q, got %q", "ansi", cfg.Theme)
		}
		if cfg.UpgradeNotify {
			t.Fatal("expected LETS_CHECK_UPDATE to disable notifications")
		}
	})

	t.Run("rejects invalid theme", func(t *testing.T) {
		unsetEnv(t, "NO_COLOR")
		unsetEnv(t, "LETS_CHECK_UPDATE")

		path := filepath.Join(t.TempDir(), "config.yaml")
		err := os.WriteFile(path, []byte("theme: vaporwave\n"), 0o644)
		if err != nil {
			t.Fatalf("failed to write settings file: %v", err)
		}

		_, err = LoadFile(path)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("rejects unknown fields", func(t *testing.T) {
		unsetEnv(t, "NO_COLOR")
		unsetEnv(t, "LETS_CHECK_UPDATE")

		path := filepath.Join(t.TempDir(), "config.yaml")
		err := os.WriteFile(path, []byte("wat: true\n"), 0o644)
		if err != nil {
			t.Fatalf("failed to write settings file: %v", err)
		}

		_, err = LoadFile(path)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	unsetEnv(t, "NO_COLOR")
	unsetEnv(t, "LETS_CHECK_UPDATE")

	configPath := filepath.Join(tmpDir, ".config", "lets", "config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0o755)
	if err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	err = os.WriteFile(configPath, []byte("no_color: true\ntheme: ansi\n"), 0o644)
	if err != nil {
		t.Fatalf("failed to write settings file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.NoColor {
		t.Fatal("expected loaded no_color to be true")
	}
	if cfg.Theme != "ansi" {
		t.Fatalf("expected loaded theme to be %q, got %q", "ansi", cfg.Theme)
	}
}

func TestApply(t *testing.T) {
	previous := color.NoColor
	t.Cleanup(func() {
		color.NoColor = previous
	})

	color.NoColor = false

	Settings{NoColor: true}.Apply()

	if !color.NoColor {
		t.Fatal("expected Apply to disable colors")
	}
}
