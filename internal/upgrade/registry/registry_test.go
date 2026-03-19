package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

	reg := NewGithubRegistry(context.Background())
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

	reg := NewGithubRegistry(context.Background())
	reg.apiURI = server.URL

	version, err := reg.GetLatestRelease()
	if err != nil {
		t.Fatalf("GetLatestRelease() error = %v", err)
	}
	if version != "v0.0.59" {
		t.Fatalf("expected version v0.0.59, got %q", version)
	}
}
