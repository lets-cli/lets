package config

import (
	"net/http"
	"net/http/httptest"
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

func TestLoadRemote(t *testing.T) {
	validConfig := "shell: bash\ncommands:\n  hello:\n    cmd: echo hello\n"

	t.Run("downloads and caches config", func(t *testing.T) {
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte(validConfig))
		}))
		defer srv.Close()

		t.Setenv("HOME", t.TempDir())

		cfg, err := LoadRemote(srv.URL, false, "0.0.0-test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := cfg.Commands["hello"]; !ok {
			t.Fatal("expected hello command")
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

		if _, err := LoadRemote(srv.URL, false, "0.0.0-test"); err != nil {
			t.Fatalf("first call error: %v", err)
		}
		if _, err := LoadRemote(srv.URL, false, "0.0.0-test"); err != nil {
			t.Fatalf("second call error: %v", err)
		}
		if requests != 1 {
			t.Fatalf("expected 1 HTTP request, got %d", requests)
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

		if _, err := LoadRemote(srv.URL, false, "0.0.0-test"); err != nil {
			t.Fatalf("prime cache error: %v", err)
		}
		if _, err := LoadRemote(srv.URL, true, "0.0.0-test"); err != nil {
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
		url := srv.URL

		if _, err := LoadRemote(url, false, "0.0.0-test"); err != nil {
			t.Fatalf("prime cache error: %v", err)
		}

		srv.Close()

		cfg, err := LoadRemote(url, true, "0.0.0-test")
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

		_, err := LoadRemote(srv.URL, false, "0.0.0-test")
		if err == nil {
			t.Fatal("expected error when no cache and download fails")
		}
	})
}
