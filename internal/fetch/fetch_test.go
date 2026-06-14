package fetch_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
	return (*recordingTracker)(p)
}

type recordingTracker recordingProgress

func (t *recordingTracker) Add(n int64) {
	p := (*recordingProgress)(t)
	p.adds = append(p.adds, n)
}

func (t *recordingTracker) Done(err error) {
	p := (*recordingProgress)(t)
	p.dones = append(p.dones, err)
}

func TestDownload(t *testing.T) {
	t.Run("downloads valid yaml", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write([]byte("shell: bash\ncommands: {}"))
		}))
		defer srv.Close()

		data, err := fetch.Download(t.Context(), srv.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "shell: bash\ncommands: {}" {
			t.Fatalf("unexpected data: %s", data)
		}
	})

	t.Run("accepts text/plain content type", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("shell: bash"))
		}))
		defer srv.Close()

		_, err := fetch.Download(t.Context(), srv.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("errors on unsupported content type", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{}"))
		}))
		defer srv.Close()

		_, err := fetch.Download(t.Context(), srv.URL)
		if err == nil {
			t.Fatal("expected error for unsupported content type")
		}
		if !strings.Contains(err.Error(), "unsupported content type") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("errors on 404", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		_, err := fetch.Download(t.Context(), srv.URL)
		if err == nil {
			t.Fatal("expected error for 404")
		}
		if !strings.Contains(err.Error(), "no such file at") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("errors on non-2xx", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		_, err := fetch.Download(t.Context(), srv.URL)
		if err == nil {
			t.Fatal("expected error for 500")
		}
		if !strings.Contains(err.Error(), "network error") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("reports progress for known size", func(t *testing.T) {
		body := []byte("shell: bash")
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.Header.Get("Accept-Encoding"); got == "gzip" {
				t.Fatalf("fetch should not request gzip because it hides content length")
			}

			w.Header().Set("Content-Type", "application/yaml")
			w.Header().Set("Content-Length", strconv.Itoa(len(body)))
			_, _ = w.Write(body)
		}))
		defer srv.Close()

		progress := &recordingProgress{}
		_, err := fetch.Download(
			t.Context(),
			srv.URL,
			fetch.WithProgress(fetch.SourceRemoteConfig, progress),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(progress.starts) != 1 {
			t.Fatalf("expected one progress start, got %d", len(progress.starts))
		}
		if progress.starts[0].Kind != fetch.SourceRemoteConfig {
			t.Fatalf("unexpected progress kind: %q", progress.starts[0].Kind)
		}
		if progress.starts[0].URL != srv.URL {
			t.Fatalf("unexpected progress url: %q", progress.starts[0].URL)
		}
		if progress.starts[0].TotalBytes != int64(len(body)) {
			t.Fatalf("expected total %d, got %d", len(body), progress.starts[0].TotalBytes)
		}
		if sum(progress.adds) != int64(len(body)) {
			t.Fatalf("expected add total %d, got %d", len(body), sum(progress.adds))
		}
		if len(progress.dones) != 1 || progress.dones[0] != nil {
			t.Fatalf("expected successful done, got %#v", progress.dones)
		}
	})

	t.Run("reports unknown size", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.(http.Flusher).Flush()
			_, _ = w.Write([]byte("shell: bash"))
		}))
		defer srv.Close()

		progress := &recordingProgress{}
		_, err := fetch.Download(
			t.Context(),
			srv.URL,
			fetch.WithProgress(fetch.SourceRemoteMixin, progress),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(progress.starts) != 1 {
			t.Fatalf("expected one progress start, got %d", len(progress.starts))
		}
		if progress.starts[0].TotalBytes > 0 {
			t.Fatalf("expected unknown total, got %d", progress.starts[0].TotalBytes)
		}
	})

	t.Run("does not start progress before validation", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{}"))
		}))
		defer srv.Close()

		progress := &recordingProgress{}
		_, err := fetch.Download(
			t.Context(),
			srv.URL,
			fetch.WithProgress(fetch.SourceRemoteConfig, progress),
		)
		if err == nil {
			t.Fatal("expected error")
		}
		if len(progress.starts) != 0 {
			t.Fatalf("expected no progress starts, got %d", len(progress.starts))
		}
	})
}

func sum(values []int64) int64 {
	var total int64
	for _, value := range values {
		total += value
	}

	return total
}
