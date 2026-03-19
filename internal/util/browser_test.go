package util

import (
	"reflect"
	"testing"
)

func TestBrowserCommand(t *testing.T) {
	t.Run("darwin", func(t *testing.T) {
		name, args, err := browserCommand("darwin", "https://lets-cli.org")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if name != "open" {
			t.Fatalf("expected open, got %q", name)
		}

		expectedArgs := []string{"https://lets-cli.org"}
		if !reflect.DeepEqual(args, expectedArgs) {
			t.Fatalf("expected args %v, got %v", expectedArgs, args)
		}
	})

	t.Run("linux", func(t *testing.T) {
		name, args, err := browserCommand("linux", "https://lets-cli.org")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if name != "xdg-open" {
			t.Fatalf("expected xdg-open, got %q", name)
		}

		expectedArgs := []string{"https://lets-cli.org"}
		if !reflect.DeepEqual(args, expectedArgs) {
			t.Fatalf("expected args %v, got %v", expectedArgs, args)
		}
	})

	t.Run("unsupported", func(t *testing.T) {
		_, _, err := browserCommand("windows", "https://lets-cli.org")
		if err == nil {
			t.Fatal("expected unsupported platform error")
		}
	})
}
