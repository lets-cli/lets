package registry

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/lets-cli/lets/internal/fetch"
)

type recordingProgress struct {
	starts []fetch.ProgressInfo
	adds   []int64
	dones  []error
}

func (p *recordingProgress) Start(info fetch.ProgressInfo) fetch.ProgressTracker {
	p.starts = append(p.starts, info)
	return p
}

func (p *recordingProgress) Add(n int64) {
	p.adds = append(p.adds, n)
}

func (p *recordingProgress) Done(err error) {
	p.dones = append(p.dones, err)
}

func TestGithubRegistryGetLatestReleaseInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/releases/latest" {
			t.Fatalf("unexpected path %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
			t.Fatalf("unexpected accept header %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name":"v0.0.59","published_at":"2026-03-17T10:00:00Z"}`))
	}))
	defer server.Close()

	reg := NewGithubRegistry()
	reg.apiURI = server.URL

	release, err := reg.GetLatestReleaseInfo(context.Background())
	if err != nil {
		t.Fatalf("GetLatestReleaseInfo() error = %v", err)
	}
	if release.TagName != "v0.0.59" {
		t.Fatalf("expected tag v0.0.59, got %q", release.TagName)
	}
	if release.PublishedAt.IsZero() {
		t.Fatal("expected publishedAt to be parsed")
	}
}

func TestGithubRegistryGetLatestRelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name":"v0.0.59","published_at":"2026-03-17T10:00:00Z"}`))
	}))
	defer server.Close()

	reg := NewGithubRegistry()
	reg.apiURI = server.URL

	version, err := reg.GetLatestRelease(context.Background())
	if err != nil {
		t.Fatalf("GetLatestRelease() error = %v", err)
	}
	if version != "v0.0.59" {
		t.Fatalf("expected version v0.0.59, got %q", version)
	}
}

func TestGithubRegistryGetLatestPrerelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/releases" {
			t.Fatalf("unexpected path %q", got)
		}
		if got := r.URL.Query().Get("per_page"); got != "100" {
			t.Fatalf("unexpected per_page query %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
			t.Fatalf("unexpected accept header %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"tag_name":"v0.0.62","published_at":"2026-03-17T10:00:00Z","prerelease":false},
			{"tag_name":"v0.0.63-rc1","published_at":"2026-03-18T10:00:00Z","prerelease":true},
			{"tag_name":"v0.0.64-rc1","published_at":"2026-03-19T10:00:00Z","prerelease":true,"draft":true},
			{"tag_name":"v0.0.63-rc2","published_at":"2026-03-20T10:00:00Z","prerelease":true}
		]`))
	}))
	defer server.Close()

	reg := NewGithubRegistry()
	reg.apiURI = server.URL

	version, err := reg.GetLatestPrerelease(context.Background())
	if err != nil {
		t.Fatalf("GetLatestPrerelease() error = %v", err)
	}
	if version != "v0.0.63-rc2" {
		t.Fatalf("expected version v0.0.63-rc2, got %q", version)
	}
}

func TestGithubRegistryGetLatestPrereleaseNoMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"tag_name":"v0.0.62","published_at":"2026-03-17T10:00:00Z","prerelease":false},
			{"tag_name":"v0.0.63-rc1","published_at":"2026-03-18T10:00:00Z","prerelease":true,"draft":true}
		]`))
	}))
	defer server.Close()

	reg := NewGithubRegistry()
	reg.apiURI = server.URL

	_, err := reg.GetLatestPrerelease(context.Background())
	if err == nil {
		t.Fatal("expected no prerelease error")
	}
}

func TestGithubRegistryDownloadReleaseBinaryReportsProgress(t *testing.T) {
	archive := releaseArchive(t, "updated binary")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/releases/download/v0.0.2/lets_Test_x86_64.tar.gz" {
			t.Fatalf("unexpected path %q", got)
		}

		w.Header().Set("Content-Type", "application/gzip")
		w.Header().Set("Content-Length", strconv.Itoa(len(archive)))
		_, _ = w.Write(archive)
	}))
	defer server.Close()

	reg := NewGithubRegistry()
	reg.repoURI = server.URL

	dstPath := filepath.Join(t.TempDir(), "lets")
	progress := &recordingProgress{}
	if err := reg.DownloadReleaseBinary(context.Background(), "lets_Test_x86_64", "v0.0.2", dstPath, progress); err != nil {
		t.Fatalf("DownloadReleaseBinary() error = %v", err)
	}

	data, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read extracted binary: %v", err)
	}
	if string(data) != "updated binary" {
		t.Fatalf("expected extracted binary, got %q", data)
	}

	if len(progress.starts) != 1 {
		t.Fatalf("expected one progress start, got %d", len(progress.starts))
	}
	if progress.starts[0].Kind != fetch.SourceSelfUpdate {
		t.Fatalf("expected self update progress, got %q", progress.starts[0].Kind)
	}
	if progress.starts[0].URL != server.URL+"/releases/download/v0.0.2/lets_Test_x86_64.tar.gz" {
		t.Fatalf("unexpected progress URL %q", progress.starts[0].URL)
	}
	if progress.starts[0].TotalBytes != int64(len(archive)) {
		t.Fatalf("expected total bytes %d, got %d", len(archive), progress.starts[0].TotalBytes)
	}
	if len(progress.adds) == 0 {
		t.Fatal("expected progress byte updates")
	}
	if len(progress.dones) != 1 || progress.dones[0] != nil {
		t.Fatalf("expected successful progress done, got %#v", progress.dones)
	}
}

func releaseArchive(t *testing.T, content string) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	data := []byte(content)
	if err := tw.WriteHeader(&tar.Header{Name: "lets", Mode: 0o755, Size: int64(len(data))}); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write(data); err != nil {
		t.Fatalf("failed to write tar content: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("failed to close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("failed to close gzip: %v", err)
	}

	return buf.Bytes()
}
