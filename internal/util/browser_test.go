package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestBrowserCommand(t *testing.T) {
	t.Run("darwin", func(t *testing.T) {
		cmd, err := browserCommand("darwin", "https://lets-cli.org")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cmd.Args[0] != "open" {
			t.Fatalf("expected open, got %q", cmd.Args[0])
		}

		expectedArgs := []string{"open", "https://lets-cli.org"}
		if !reflect.DeepEqual(cmd.Args, expectedArgs) {
			t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args)
		}
	})

	t.Run("linux", func(t *testing.T) {
		cmd, err := browserCommand("linux", "https://lets-cli.org")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cmd.Args[0] != "xdg-open" {
			t.Fatalf("expected xdg-open, got %q", cmd.Args[0])
		}

		expectedArgs := []string{"xdg-open", "https://lets-cli.org"}
		if !reflect.DeepEqual(cmd.Args, expectedArgs) {
			t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args)
		}
	})

	t.Run("unsupported", func(t *testing.T) {
		_, err := browserCommand("windows", "https://lets-cli.org")
		if err == nil {
			t.Fatal("expected unsupported platform error")
		}
		if !strings.Contains(err.Error(), "windows") {
			t.Fatalf("expected error to mention platform %q, got %q", "windows", err.Error())
		}
	})
}
