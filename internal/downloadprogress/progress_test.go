package downloadprogress

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/lets-cli/lets/internal/fetch"
)

func TestObserver(t *testing.T) {
	t.Run("renders unchanged downloading label for known size", func(t *testing.T) {
		var out bytes.Buffer
		observer := New(&out, WithWidth(120), WithNoColor(true), WithThrottle(0), WithFinalPause(0))

		tracker := observer.Start(fetch.ProgressInfo{
			Kind:       fetch.SourceRemoteConfig,
			URL:        "https://user:pass@example.com/path/lets.yaml?token=secret#fragment",
			TotalBytes: 4,
		})
		tracker.Add(4)
		tracker.Done(nil)

		got := out.String()
		if !strings.Contains(got, "Downloading lets.yaml") ||
			!strings.Contains(got, "\n\r") ||
			!strings.Contains(got, "100% 4 B/4 B") {
			t.Fatalf("expected completed line, got %q", got)
		}
		if strings.Contains(got, "Downloaded") {
			t.Fatalf("did not expect label to change to Downloaded, got %q", got)
		}
		if !strings.Contains(got, "#") {
			t.Fatalf("expected completed line to include progress bar, got %q", got)
		}
		if strings.Contains(got, "secret") || strings.Contains(got, "user:pass") || strings.Contains(got, "fragment") {
			t.Fatalf("expected URL secrets to be redacted, got %q", got)
		}
	})

	t.Run("renders unknown size without percent", func(t *testing.T) {
		var out bytes.Buffer
		observer := New(&out, WithWidth(120), WithNoColor(true), WithThrottle(0))

		tracker := observer.Start(fetch.ProgressInfo{
			Kind:       fetch.SourceRemoteMixin,
			URL:        "https://example.com/mixin.yaml",
			TotalBytes: -1,
		})
		tracker.Add(1536)
		tracker.Done(nil)

		got := out.String()
		if !strings.Contains(got, "Downloading mixin.yaml") ||
			!strings.Contains(got, "\n\r") ||
			!strings.Contains(got, "1.5 KiB") {
			t.Fatalf("expected unknown-size completed line, got %q", got)
		}
		if strings.Contains(got, "Downloaded") {
			t.Fatalf("did not expect label to change to Downloaded, got %q", got)
		}
		if strings.Contains(got, "%") {
			t.Fatalf("did not expect percent for unknown size, got %q", got)
		}
	})

	t.Run("throttles intermediate renders", func(t *testing.T) {
		var out bytes.Buffer
		now := time.Unix(0, 0)
		observer := New(&out, WithWidth(120), WithNoColor(true), WithNow(func() time.Time { return now }))

		tracker := observer.Start(fetch.ProgressInfo{
			Kind:       fetch.SourceRemoteConfig,
			URL:        "https://example.com/lets.yaml",
			TotalBytes: -1,
		})
		startOutput := out.String()
		tracker.Add(1)
		if out.String() != startOutput {
			t.Fatalf("expected add before throttle to skip render")
		}

		now = now.Add(100 * time.Millisecond)
		tracker.Add(1)
		if out.String() == startOutput {
			t.Fatalf("expected add at throttle boundary to render")
		}
	})
}

func TestFormatBytes(t *testing.T) {
	tests := map[int64]string{
		42:          "42 B",
		1024:        "1.0 KiB",
		1536:        "1.5 KiB",
		1024 * 1024: "1.0 MiB",
	}

	for input, want := range tests {
		if got := formatBytes(input); got != want {
			t.Fatalf("formatBytes(%d) = %q, want %q", input, got, want)
		}
	}
}

func TestDownloadLabel(t *testing.T) {
	got := downloadLabel("https://user:pass@example.com/path/lets.yaml?token=secret#fragment")
	if got != "lets.yaml" {
		t.Fatalf("expected filename label, got %q", got)
	}
}

func TestProgressModel(t *testing.T) {
	observer := New(&bytes.Buffer{}, WithWidth(120), WithNoColor(true), WithFinalPause(0))
	model := newProgressModel(observer, "lets.yaml", 1107, make(chan struct{}))

	_, cmd := model.Update(progressMsg{read: 512, total: 1107})
	if cmd == nil {
		t.Fatal("expected progress update to return animation command")
	}

	updated, _ := model.Update(progressDoneMsg{read: 1107, total: 1107})
	finalModel := updated.(progressModel)
	got := finalModel.View().Content
	if strings.Contains(got, "Downloaded") {
		t.Fatalf("did not expect label to change to Downloaded, got %q", got)
	}
	if !strings.Contains(got, "100% 1.1 KiB/1.1 KiB") {
		t.Fatalf("expected final progress status, got %q", got)
	}
}
