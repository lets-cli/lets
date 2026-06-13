package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	skillpkg "github.com/lets-cli/lets/internal/skills"
)

func TestSelfSkillsCmd(t *testing.T) {
	t.Run("should show bundled lets skill", func(t *testing.T) {
		bufOut := new(bytes.Buffer)
		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "skills", "show"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(new(bytes.Buffer))
		InitSelfCmd(rootCmd, "v0.0.0-test")

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(bufOut.String(), "name: lets") {
			t.Fatalf("expected lets skill, got %q", bufOut.String())
		}
	})

	t.Run("should prompt and install local skill", func(t *testing.T) {
		repoDir := chdirTempGitRepo(t)
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "skills", "install"})
		rootCmd.SetIn(strings.NewReader("local\n"))
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(bufErr)
		InitSelfCmd(rootCmd, "v0.0.0-test")

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		skillPath := filepath.Join(repoDir, skillpkg.SkillsRelDir, skillpkg.LetsName, skillpkg.SkillFile)
		assertSkillFile(t, skillPath)
		if !strings.Contains(bufErr.String(), "Install lets skill:") {
			t.Fatalf("expected install prompt, got %q", bufErr.String())
		}
		if !strings.Contains(bufErr.String(), filepath.Join(repoDir, skillpkg.SkillsRelDir, skillpkg.LetsName)) {
			t.Fatalf("expected local install path in prompt, got %q", bufErr.String())
		}
		if !strings.Contains(bufOut.String(), "Installed ") {
			t.Fatalf("expected install output, got %q", bufOut.String())
		}
	})

	t.Run("should prompt and install global skill", func(t *testing.T) {
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)
		bufOut := new(bytes.Buffer)
		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "skills", "install"})
		rootCmd.SetIn(strings.NewReader("global\n"))
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(new(bytes.Buffer))
		InitSelfCmd(rootCmd, "v0.0.0-test")

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		skillPath := filepath.Join(homeDir, skillpkg.SkillsRelDir, skillpkg.LetsName, skillpkg.SkillFile)
		assertSkillFile(t, skillPath)
	})

	t.Run("should accept numbered picker choice", func(t *testing.T) {
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)
		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "skills", "install"})
		rootCmd.SetIn(strings.NewReader("2\n"))
		rootCmd.SetOut(new(bytes.Buffer))
		rootCmd.SetErr(new(bytes.Buffer))
		InitSelfCmd(rootCmd, "v0.0.0-test")

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		skillPath := filepath.Join(homeDir, skillpkg.SkillsRelDir, skillpkg.LetsName, skillpkg.SkillFile)
		assertSkillFile(t, skillPath)
	})

	t.Run("should not overwrite existing skill without force", func(t *testing.T) {
		repoDir := chdirTempGitRepo(t)
		skillPath := filepath.Join(repoDir, skillpkg.SkillsRelDir, skillpkg.LetsName, skillpkg.SkillFile)
		if err := os.MkdirAll(filepath.Dir(skillPath), 0o755); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(skillPath, []byte("custom"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		bufOut := new(bytes.Buffer)
		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "skills", "install", "--local"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(new(bytes.Buffer))
		InitSelfCmd(rootCmd, "v0.0.0-test")

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(skillPath)
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}
		if string(data) != "custom" {
			t.Fatalf("expected existing skill to remain unchanged, got %q", string(data))
		}
		if !strings.Contains(bufOut.String(), "already exists") {
			t.Fatalf("expected already exists output, got %q", bufOut.String())
		}
	})

	t.Run("should update installed skill", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		repoDir := chdirTempGitRepo(t)
		skillPath := filepath.Join(repoDir, skillpkg.SkillsRelDir, skillpkg.LetsName, skillpkg.SkillFile)
		if err := os.MkdirAll(filepath.Dir(skillPath), 0o755); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(skillPath, []byte("old"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		bufOut := new(bytes.Buffer)
		rootCmd := CreateRootCommand("v0.0.0-test", "")
		rootCmd.SetArgs([]string{"self", "skills", "update"})
		rootCmd.SetOut(bufOut)
		rootCmd.SetErr(new(bytes.Buffer))
		InitSelfCmd(rootCmd, "v0.0.0-test")

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertSkillFile(t, skillPath)
		if !strings.Contains(bufOut.String(), "Updated ") {
			t.Fatalf("expected update output, got %q", bufOut.String())
		}
	})
}

func chdirTempGitRepo(t *testing.T) string {
	t.Helper()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Fatalf("Chdir(%s) error = %v", oldWd, err)
		}
	})

	repoDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(repoDir, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Chdir(%s) error = %v", repoDir, err)
	}

	return repoDir
}

func assertSkillFile(t *testing.T, skillPath string) {
	t.Helper()

	data, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", skillPath, err)
	}
	if !bytes.Equal(data, skillpkg.LetsSkill()) {
		t.Fatalf("unexpected skill content at %s", skillPath)
	}
}
