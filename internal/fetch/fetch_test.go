package fetch_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lets-cli/lets/internal/fetch"
)

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
}
