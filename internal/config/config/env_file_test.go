package config

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"gopkg.in/yaml.v3"
)

func loadConfigFixture(t *testing.T, workDir string, text string) *Config {
	t.Helper()

	cfg := NewConfig(workDir, filepath.Join(workDir, "lets.yaml"), filepath.Join(workDir, ".lets"))
	buf := bytes.NewBufferString(text)
	if err := yaml.NewDecoder(buf).Decode(cfg); err != nil {
		t.Fatalf("config fixture decode error: %s", err)
	}

	return cfg
}

func writeFixtureFile(t *testing.T, dir string, name string, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write fixture %s: %s", path, err)
	}
}

func TestParseEnvFiles(t *testing.T) {
	t.Run("parses scalar env_file form", func(t *testing.T) {
		text := "env_file: .env\n"

		var raw struct {
			EnvFiles *EnvFiles `yaml:"env_file"`
		}
		if err := yaml.NewDecoder(bytes.NewBufferString(text)).Decode(&raw); err != nil {
			t.Fatalf("unexpected decode error: %s", err)
		}

		if len(raw.EnvFiles.Items) != 1 {
			t.Fatalf("expected 1 env file, got %d", len(raw.EnvFiles.Items))
		}

		if got := raw.EnvFiles.Items[0]; got.Name != ".env" || !got.Required {
			t.Fatalf("unexpected env file: %#v", got)
		}
	})

	t.Run("parses map env_file form", func(t *testing.T) {
		text := dedent.Dedent(`
		env_file:
		  name: .env.prod
		  required: false
		`)

		var raw struct {
			EnvFiles *EnvFiles `yaml:"env_file"`
		}
		if err := yaml.NewDecoder(bytes.NewBufferString(text)).Decode(&raw); err != nil {
			t.Fatalf("unexpected decode error: %s", err)
		}

		if len(raw.EnvFiles.Items) != 1 {
			t.Fatalf("expected 1 env file, got %d", len(raw.EnvFiles.Items))
		}

		if got := raw.EnvFiles.Items[0]; got.Name != ".env.prod" || got.Required {
			t.Fatalf("unexpected env file: %#v", got)
		}
	})

	t.Run("parses mixed env_file forms", func(t *testing.T) {
		text := dedent.Dedent(`
		env_file:
		  - .env
		  - -.env.local
		  - name: .env.prod
		    required: false
		`)

		var raw struct {
			EnvFiles *EnvFiles `yaml:"env_file"`
		}
		if err := yaml.NewDecoder(bytes.NewBufferString(text)).Decode(&raw); err != nil {
			t.Fatalf("unexpected decode error: %s", err)
		}

		if len(raw.EnvFiles.Items) != 3 {
			t.Fatalf("expected 3 env files, got %d", len(raw.EnvFiles.Items))
		}

		if got := raw.EnvFiles.Items[0]; got.Name != ".env" || !got.Required {
			t.Fatalf("unexpected first env file: %#v", got)
		}

		if got := raw.EnvFiles.Items[1]; got.Name != ".env.local" || got.Required {
			t.Fatalf("unexpected second env file: %#v", got)
		}

		if got := raw.EnvFiles.Items[2]; got.Name != ".env.prod" || got.Required {
			t.Fatalf("unexpected third env file: %#v", got)
		}
	})

	t.Run("rejects env_file item without name", func(t *testing.T) {
		text := dedent.Dedent(`
		env_file:
		  - required: false
		`)

		var raw struct {
			EnvFiles *EnvFiles `yaml:"env_file"`
		}
		err := yaml.NewDecoder(bytes.NewBufferString(text)).Decode(&raw)
		if err == nil {
			t.Fatal("expected decode error")
		}

		if !strings.Contains(err.Error(), "env_file name can not be empty") {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("rejects map form with optional dash prefix", func(t *testing.T) {
		text := dedent.Dedent(`
		env_file:
		  name: -.env.local
		`)

		var raw struct {
			EnvFiles *EnvFiles `yaml:"env_file"`
		}
		err := yaml.NewDecoder(bytes.NewBufferString(text)).Decode(&raw)
		if err == nil {
			t.Fatal("expected decode error")
		}

		if !strings.Contains(err.Error(), "use required: false instead") {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("rejects invalid top-level env_file kind", func(t *testing.T) {
		envFiles := &EnvFiles{}
		err := envFiles.UnmarshalYAML(&yaml.Node{Kind: yaml.AliasNode})
		if err == nil {
			t.Fatal("expected unmarshal error")
		}

		if !strings.Contains(err.Error(), "env_file must be a string, map, or sequence") {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}

func TestEnvFilesLoad(t *testing.T) {
	workDir := t.TempDir()
	writeFixtureFile(t, workDir, ".env.first", "VALUE=first\n")
	writeFixtureFile(t, workDir, ".env.second", "VALUE=second\nSECOND=two\n")
	writeFixtureFile(t, workDir, ".env.invalid", "NOT VALID\n")

	cfg := Config{WorkDir: workDir}

	t.Run("later files override earlier files", func(t *testing.T) {
		envFiles := &EnvFiles{
			Items: []EnvFile{
				{Name: ".env.first", Required: true},
				{Name: ".env.second", Required: true},
			},
		}

		got, err := envFiles.Load(cfg, nil)
		if err != nil {
			t.Fatalf("unexpected load error: %s", err)
		}

		if got["VALUE"] != "second" || got["SECOND"] != "two" {
			t.Fatalf("unexpected env map: %#v", got)
		}
	})

	t.Run("skips optional missing files", func(t *testing.T) {
		envFiles := &EnvFiles{
			Items: []EnvFile{
				{Name: ".env.first", Required: true},
				{Name: ".env.missing", Required: false},
			},
		}

		got, err := envFiles.Load(cfg, nil)
		if err != nil {
			t.Fatalf("unexpected load error: %s", err)
		}

		if got["VALUE"] != "first" {
			t.Fatalf("expected VALUE from existing file, got %#v", got)
		}
	})

	t.Run("fails on missing required file", func(t *testing.T) {
		envFiles := &EnvFiles{
			Items: []EnvFile{{Name: ".env.missing", Required: true}},
		}

		_, err := envFiles.Load(cfg, nil)
		if err == nil {
			t.Fatal("expected load error")
		}

		if !strings.Contains(err.Error(), filepath.Join(workDir, ".env.missing")) {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("fails on invalid env file format with filename", func(t *testing.T) {
		envFiles := &EnvFiles{
			Items: []EnvFile{{Name: ".env.invalid", Required: true}},
		}

		_, err := envFiles.Load(cfg, nil)
		if err == nil {
			t.Fatal("expected load error")
		}

		if !strings.Contains(err.Error(), filepath.Join(workDir, ".env.invalid")) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}

func TestConfigSetupEnvWithEnvFile(t *testing.T) {
	workDir := t.TempDir()
	writeFixtureFile(t, workDir, ".env.global", "FROM_FILE=global-file\nOVERRIDE=from-file\n")
	writeFixtureFile(t, workDir, ".env."+runtime.GOOS, "OS_FILE=os-file\n")

	cfg := loadConfigFixture(t, workDir, dedent.Dedent(`
	shell: bash
	env:
	  TARGET: global
	  OVERRIDE: from-env
	env_file:
	  - .env.${TARGET}
	  - .env.${LETS_OS}
	commands:
	  echo:
	    cmd: echo ok
	`))

	if err := cfg.SetupEnv(); err != nil {
		t.Fatalf("unexpected setup error: %s", err)
	}

	got := cfg.GetEnv()
	if got["FROM_FILE"] != "global-file" {
		t.Fatalf("expected FROM_FILE from env_file, got %#v", got)
	}

	if got["OVERRIDE"] != "from-file" {
		t.Fatalf("expected env_file to override env, got %#v", got)
	}

	if got["OS_FILE"] != "os-file" {
		t.Fatalf("expected OS_FILE from LETS_OS interpolation, got %#v", got)
	}
}

func TestCommandGetEnvWithEnvFile(t *testing.T) {
	workDir := t.TempDir()
	writeFixtureFile(t, workDir, ".env.global", "GLOBAL_FROM_FILE=global-file\n")
	writeFixtureFile(t, workDir, ".env.command.dev", "COMMAND_FROM_FILE=command-file\nCOMMAND_OVERRIDE=from-file\n")

	cfg := loadConfigFixture(t, workDir, dedent.Dedent(`
	shell: bash
	env:
	  TARGET: global
	  COMMAND_TARGET: command
	env_file: .env.${TARGET}
	commands:
	  echo:
	    env:
	      SUFFIX: dev
	      COMMAND_OVERRIDE: from-env
	    env_file: .env.${COMMAND_TARGET}.${SUFFIX}
	    cmd: echo ok
	`))

	if err := cfg.SetupEnv(); err != nil {
		t.Fatalf("unexpected setup error: %s", err)
	}

	cmd := cfg.Commands["echo"]
	got, err := cmd.GetEnv(*cfg, cfg.CommandBuiltinEnv(cmd, cfg.Shell, cfg.WorkDir))
	if err != nil {
		t.Fatalf("unexpected command env error: %s", err)
	}

	if got["COMMAND_FROM_FILE"] != "command-file" {
		t.Fatalf("expected command env_file value, got %#v", got)
	}

	if got["COMMAND_OVERRIDE"] != "from-file" {
		t.Fatalf("expected command env_file to override env, got %#v", got)
	}
}

func TestCommandGetEnvDoesNotReuseBuiltinEnvCache(t *testing.T) {
	workDir := t.TempDir()
	writeFixtureFile(t, workDir, ".env.one", "VALUE=from-one\n")
	writeFixtureFile(t, workDir, ".env.two", "VALUE=from-two\n")

	cfg := loadConfigFixture(t, workDir, dedent.Dedent(`
	shell: bash
	commands:
	  echo:
	    env_file: .env.${LETS_COMMAND_ARGS}
	    cmd: echo ok
	`))

	if err := cfg.SetupEnv(); err != nil {
		t.Fatalf("unexpected setup error: %s", err)
	}

	cmd := cfg.Commands["echo"]

	cmd.Args = []string{"one"}
	gotOne, err := cmd.GetEnv(*cfg, cfg.CommandBuiltinEnv(cmd, cfg.Shell, cfg.WorkDir))
	if err != nil {
		t.Fatalf("unexpected command env error: %s", err)
	}

	cmd.Args = []string{"two"}
	gotTwo, err := cmd.GetEnv(*cfg, cfg.CommandBuiltinEnv(cmd, cfg.Shell, cfg.WorkDir))
	if err != nil {
		t.Fatalf("unexpected command env error: %s", err)
	}

	if gotOne["VALUE"] != "from-one" {
		t.Fatalf("expected first env file to be used, got %#v", gotOne)
	}

	if gotTwo["VALUE"] != "from-two" {
		t.Fatalf("expected second env file to be used, got %#v", gotTwo)
	}
}
