package config

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lets-cli/lets/internal/fetch"
)

type recordingProgress struct {
	starts []fetch.ProgressInfo
}

func (p *recordingProgress) Start(info fetch.ProgressInfo) fetch.ProgressTracker {
	p.starts = append(p.starts, info)
	return noopTracker{}
}

type noopTracker struct{}

func (noopTracker) Add(int64) {}

func (noopTracker) Done(error) {}

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

func TestLoadRemote(t *testing.T) {
	validConfig := "shell: bash\ncommands:\n  hello:\n    cmd: echo hello\n"
	ctx := context.Background()

	t.Run("downloads and caches config", func(t *testing.T) {
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte(validConfig))
		}))
		defer srv.Close()

		t.Setenv("HOME", t.TempDir())
		t.Chdir(t.TempDir())

		cfg, err := LoadRemote(ctx, srv.URL, false, "0.0.0-test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := cfg.Commands["hello"]; !ok {
			t.Fatal("expected hello command")
		}
		if cfg.RemoteSource != srv.URL {
			t.Fatalf("expected RemoteSource=%q, got %q", srv.URL, cfg.RemoteSource)
		}
		if requests != 1 {
			t.Fatalf("expected 1 HTTP request, got %d", requests)
		}
	})

	t.Run("uses cache on second call without --no-cache", func(t *testing.T) {
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte(validConfig))
		}))
		defer srv.Close()

		t.Setenv("HOME", t.TempDir())
		t.Chdir(t.TempDir())

		if _, err := LoadRemote(ctx, srv.URL, false, "0.0.0-test"); err != nil {
			t.Fatalf("first call error: %v", err)
		}
		progress := &recordingProgress{}
		if _, err := LoadRemote(ctx, srv.URL, false, "0.0.0-test", WithProgress(progress)); err != nil {
			t.Fatalf("second call error: %v", err)
		}
		if requests != 1 {
			t.Fatalf("expected 1 HTTP request, got %d", requests)
		}
		if len(progress.starts) != 0 {
			t.Fatalf("expected no progress starts for cache hit, got %d", len(progress.starts))
		}
	})

	t.Run("reports progress for remote config download", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte(validConfig))
		}))
		defer srv.Close()

		t.Setenv("HOME", t.TempDir())
		t.Chdir(t.TempDir())

		progress := &recordingProgress{}
		if _, err := LoadRemote(ctx, srv.URL, false, "0.0.0-test", WithProgress(progress)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(progress.starts) != 1 {
			t.Fatalf("expected 1 progress start, got %d", len(progress.starts))
		}
		if progress.starts[0].Kind != fetch.SourceRemoteConfig {
			t.Fatalf("expected remote config progress, got %q", progress.starts[0].Kind)
		}
	})

	t.Run("re-downloads with --no-cache", func(t *testing.T) {
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte(validConfig))
		}))
		defer srv.Close()

		t.Setenv("HOME", t.TempDir())
		t.Chdir(t.TempDir())

		if _, err := LoadRemote(ctx, srv.URL, false, "0.0.0-test"); err != nil {
			t.Fatalf("prime cache error: %v", err)
		}
		if _, err := LoadRemote(ctx, srv.URL, true, "0.0.0-test"); err != nil {
			t.Fatalf("no-cache call error: %v", err)
		}
		if requests != 2 {
			t.Fatalf("expected 2 HTTP requests, got %d", requests)
		}
	})

	t.Run("falls back to cache when --no-cache download fails", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte(validConfig))
		}))

		t.Setenv("HOME", t.TempDir())
		t.Chdir(t.TempDir())
		url := srv.URL

		if _, err := LoadRemote(ctx, url, false, "0.0.0-test"); err != nil {
			t.Fatalf("prime cache error: %v", err)
		}

		srv.Close()

		cfg, err := LoadRemote(ctx, url, true, "0.0.0-test")
		if err != nil {
			t.Fatalf("expected fallback to cache, got error: %v", err)
		}
		if _, ok := cfg.Commands["hello"]; !ok {
			t.Fatal("expected hello command from cache")
		}
	})

	t.Run("errors when download fails with no cache", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		srv.Close()

		t.Setenv("HOME", t.TempDir())
		t.Chdir(t.TempDir())

		_, err := LoadRemote(ctx, srv.URL, false, "0.0.0-test")
		if err == nil {
			t.Fatal("expected error when no cache and download fails")
		}
	})

	t.Run("reports progress for remote mixin cache miss only", func(t *testing.T) {
		mixinConfig := "commands:\n  mixed:\n    cmd: echo mixed\n"
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte(mixinConfig))
		}))
		defer srv.Close()

		tempDir := t.TempDir()
		mainConfig := "shell: bash\nmixins:\n  - url: " + srv.URL + "\ncommands:\n  ok:\n    cmd: echo ok\n"
		if err := os.WriteFile(filepath.Join(tempDir, "lets.yaml"), []byte(mainConfig), 0o644); err != nil {
			t.Fatalf("write config: %v", err)
		}

		progress := &recordingProgress{}
		if _, err := LoadWithContext(ctx, "", tempDir, "0.0.0-test", WithProgress(progress)); err != nil {
			t.Fatalf("first load error: %v", err)
		}
		if len(progress.starts) != 1 {
			t.Fatalf("expected 1 progress start, got %d", len(progress.starts))
		}
		if progress.starts[0].Kind != fetch.SourceRemoteMixin {
			t.Fatalf("expected remote mixin progress, got %q", progress.starts[0].Kind)
		}

		cacheHitProgress := &recordingProgress{}
		if _, err := LoadWithContext(ctx, "", tempDir, "0.0.0-test", WithProgress(cacheHitProgress)); err != nil {
			t.Fatalf("second load error: %v", err)
		}
		if len(cacheHitProgress.starts) != 0 {
			t.Fatalf("expected no progress starts for remote mixin cache hit, got %d", len(cacheHitProgress.starts))
		}
		if requests != 1 {
			t.Fatalf("expected 1 HTTP request, got %d", requests)
		}
	})

	t.Run("re-downloads remote mixins with --no-cache", func(t *testing.T) {
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte("commands:\n  mixed:\n    cmd: echo mixed\n"))
		}))
		defer srv.Close()

		tempDir := t.TempDir()
		mainConfig := "shell: bash\nmixins:\n  - url: " + srv.URL + "\ncommands:\n  ok:\n    cmd: echo ok\n"
		if err := os.WriteFile(filepath.Join(tempDir, "lets.yaml"), []byte(mainConfig), 0o644); err != nil {
			t.Fatalf("write config: %v", err)
		}

		if _, err := LoadWithContext(ctx, "", tempDir, "0.0.0-test"); err != nil {
			t.Fatalf("prime cache error: %v", err)
		}
		progress := &recordingProgress{}
		if _, err := LoadWithContext(ctx, "", tempDir, "0.0.0-test", WithNoCache(), WithProgress(progress)); err != nil {
			t.Fatalf("no-cache load error: %v", err)
		}
		if requests != 2 {
			t.Fatalf("expected 2 HTTP requests, got %d", requests)
		}
		if len(progress.starts) != 1 {
			t.Fatalf("expected 1 progress start, got %d", len(progress.starts))
		}
		if progress.starts[0].Kind != fetch.SourceRemoteMixin {
			t.Fatalf("expected remote mixin progress, got %q", progress.starts[0].Kind)
		}
	})

	t.Run("does not replace remote mixin cache until downloaded mixin is valid", func(t *testing.T) {
		mixinConfig := "commands:\n  mixed:\n    cmd: echo mixed\n"
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte(mixinConfig))
		}))
		defer srv.Close()

		tempDir := t.TempDir()
		mainConfig := "shell: bash\nmixins:\n  - url: " + srv.URL + "\ncommands:\n  ok:\n    cmd: echo ok\n"
		if err := os.WriteFile(filepath.Join(tempDir, "lets.yaml"), []byte(mainConfig), 0o644); err != nil {
			t.Fatalf("write config: %v", err)
		}

		if _, err := LoadWithContext(ctx, "", tempDir, "0.0.0-test"); err != nil {
			t.Fatalf("prime cache error: %v", err)
		}

		mixinConfig = "commands:\n  broken:\n    xxx\n    cmd: echo broken\n"
		if _, err := LoadWithContext(ctx, "", tempDir, "0.0.0-test", WithNoCache()); err == nil {
			t.Fatal("expected no-cache load to fail for invalid downloaded mixin")
		}

		cfg, err := LoadWithContext(ctx, "", tempDir, "0.0.0-test")
		if err != nil {
			t.Fatalf("cached load error: %v", err)
		}
		if _, ok := cfg.Commands["mixed"]; !ok {
			t.Fatal("expected valid cached mixin command")
		}
	})
}

// The root dir is where lets was invoked, never where the config file sits.
func TestRootDirIsInvocationDir(t *testing.T) {
	ctx := context.Background()

	t.Run("local config in a child dir does not move the root", func(t *testing.T) {
		root := t.TempDir()
		configDir := filepath.Join(root, "sub")
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		writeFile(t, filepath.Join(configDir, "lets.yaml"), "shell: bash\ncommands:\n  hi:\n    cmd: echo hi\n")

		t.Chdir(root)

		cfg, err := LoadWithContext(ctx, "sub/lets.yaml", "", "0.0.0-test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertSameDir(t, "RootDir", cfg.RootDir, root)
		assertSameDir(t, "ConfigDir", cfg.ConfigDir, configDir)
	})

	t.Run("remote config roots at the cwd, not the cache dir", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte("shell: bash\ncommands:\n  hi:\n    cmd: echo hi\n"))
		}))
		defer srv.Close()

		root := t.TempDir()
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Chdir(root)

		cfg, err := LoadRemote(ctx, srv.URL, false, "0.0.0-test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertSameDir(t, "RootDir", cfg.RootDir, root)
		assertSameDir(t, "ConfigDir", cfg.ConfigDir, filepath.Join(home, ".config", "lets", "remote-configs"))

		// what a command actually receives, not just the field behind it
		builtin := cfg.BuiltinEnv(cfg.Shell)
		assertSameDir(t, "LETS_CONFIG_DIR", builtin["LETS_CONFIG_DIR"], cfg.ConfigDir)

		if builtin["LETS_CONFIG"] != srv.URL {
			t.Fatalf("expected LETS_CONFIG=%q, got %q", srv.URL, builtin["LETS_CONFIG"])
		}
	})
}

// A remote config's ConfigDir is its cache dir, which only holds the downloaded
// yaml, so a local mixin path there can never resolve.
func TestRemoteConfigRejectsLocalMixin(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write([]byte("shell: bash\nmixins:\n  - local.yaml\ncommands:\n  hi:\n    cmd: echo hi\n"))
	}))
	defer srv.Close()

	root := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Chdir(root)

	// present next to the cwd, to prove it is not picked up from there either
	writeFile(t, filepath.Join(root, "local.yaml"), "commands:\n  local:\n    cmd: echo local\n")

	_, err := LoadRemote(context.Background(), srv.URL, false, "0.0.0-test")
	if err == nil {
		t.Fatal("expected local mixin in a remote config to fail")
	}
	if !strings.Contains(err.Error(), "can only mix in URLs") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// macOS resolves t.TempDir() under /var, a symlink to /private/var, while
// os.Getwd reports the resolved path.
func assertSameDir(t *testing.T, name string, got string, want string) {
	t.Helper()

	resolved, err := filepath.EvalSymlinks(want)
	if err != nil {
		t.Fatalf("resolve %s: %v", want, err)
	}

	if got != resolved && got != want {
		t.Fatalf("expected %s=%q, got %q", name, resolved, got)
	}
}
