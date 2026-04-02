package lsp

import (
	"testing"
	"time"

	"github.com/tliron/glsp"
	lsp "github.com/tliron/glsp/protocol_3_16"
)

func TestDocumentDebouncerCoalescesRepeatedSchedules(t *testing.T) {
	events := make(chan string, 2)
	debouncer := newDocumentDebouncer(20*time.Millisecond, func(uri string) {
		events <- uri
	})
	defer debouncer.Stop()

	debouncer.Schedule("file:///tmp/lets.yaml")
	debouncer.Schedule("file:///tmp/lets.yaml")

	select {
	case got := <-events:
		if got != "file:///tmp/lets.yaml" {
			t.Fatalf("refresh uri = %q, want %q", got, "file:///tmp/lets.yaml")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for debounced refresh")
	}

	select {
	case got := <-events:
		t.Fatalf("unexpected extra refresh for %q", got)
	case <-time.After(60 * time.Millisecond):
	}
}

func TestTextDocumentDidChangeUsesLatestDocumentAfterDebounce(t *testing.T) {
	server := &lspServer{
		storage: newStorage(),
		parser:  newParser(logger),
		index:   newIndex(logger),
		log:     logger,
	}
	server.refresh = newDocumentDebouncer(20*time.Millisecond, server.refreshDocument)
	defer server.refresh.Stop()

	params := &lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{URI: "file:///tmp/lets.yaml"},
		},
	}

	params.ContentChanges = []any{
		lsp.TextDocumentContentChangeEventWhole{
			Text: `commands:
  build:
    cmd: echo build`,
		},
	}

	if err := server.textDocumentDidChange(&glsp.Context{}, params); err != nil {
		t.Fatalf("first textDocumentDidChange() error = %v", err)
	}

	params.ContentChanges = []any{
		lsp.TextDocumentContentChangeEventWhole{
			Text: `commands:
  release:
    cmd: echo release`,
		},
	}

	if err := server.textDocumentDidChange(&glsp.Context{}, params); err != nil {
		t.Fatalf("second textDocumentDidChange() error = %v", err)
	}

	deadline := time.Now().Add(300 * time.Millisecond)
	for time.Now().Before(deadline) {
		if _, ok := server.index.findCommand("release"); ok {
			if _, ok := server.index.findCommand("build"); ok {
				t.Fatal("expected stale command to be removed after debounced refresh")
			}

			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("timed out waiting for debounced document refresh")
}
