package executor

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestPrependToChain(t *testing.T) {
	t.Run("new error creates single-element chain", func(t *testing.T) {
		orig := fmt.Errorf("something failed")
		result := prependToChain("lint", orig)

		depErr, ok := result.(*DependencyError)
		if !ok {
			t.Fatalf("expected *DependencyError, got %T", result)
		}
		if len(depErr.Chain) != 1 || depErr.Chain[0] != "lint" {
			t.Errorf("expected chain [lint], got %v", depErr.Chain)
		}
		if depErr.Err != orig {
			t.Errorf("expected original error to be preserved")
		}
	})

	t.Run("existing DependencyError gets name prepended", func(t *testing.T) {
		base := &DependencyError{
			Chain: []string{"lint"},
			Err:   fmt.Errorf("orig"),
		}
		result := prependToChain("build", base)

		depErr, ok := result.(*DependencyError)
		if !ok {
			t.Fatalf("expected *DependencyError, got %T", result)
		}
		if len(depErr.Chain) != 2 || depErr.Chain[0] != "build" || depErr.Chain[1] != "lint" {
			t.Errorf("expected chain [build lint], got %v", depErr.Chain)
		}
	})

	t.Run("three deep chain accumulates correctly", func(t *testing.T) {
		err := fmt.Errorf("exit 1")
		err = prependToChain("lint", err)
		err = prependToChain("build", err)
		err = prependToChain("deploy", err)

		depErr, ok := err.(*DependencyError)
		if !ok {
			t.Fatalf("expected *DependencyError, got %T", err)
		}
		want := []string{"deploy", "build", "lint"}
		if len(depErr.Chain) != 3 {
			t.Fatalf("expected chain length 3, got %d: %v", len(depErr.Chain), depErr.Chain)
		}
		for i, name := range want {
			if depErr.Chain[i] != name {
				t.Errorf("chain[%d]: want %q, got %q", i, name, depErr.Chain[i])
			}
		}
	})
}

func TestDependencyErrorExitCode(t *testing.T) {
	t.Run("returns 1 when ExecuteError wraps non-ExitError", func(t *testing.T) {
		execErr := &ExecuteError{err: fmt.Errorf("failed")}
		depErr := &DependencyError{Chain: []string{"lint"}, Err: execErr}
		if depErr.ExitCode() != 1 {
			t.Errorf("expected exit code 1, got %d", depErr.ExitCode())
		}
	})

	t.Run("returns 1 for non-ExecuteError", func(t *testing.T) {
		depErr := &DependencyError{Chain: []string{"lint"}, Err: fmt.Errorf("plain error")}
		if depErr.ExitCode() != 1 {
			t.Errorf("expected default exit code 1, got %d", depErr.ExitCode())
		}
	})
}

func TestDependencyErrorError(t *testing.T) {
	inner := fmt.Errorf("inner error message")
	depErr := &DependencyError{Chain: []string{"lint"}, Err: inner}
	if depErr.Error() != "inner error message" {
		t.Errorf("expected Error() to delegate to Err, got %q", depErr.Error())
	}
}

func TestPrintDependencyTree(t *testing.T) {
	t.Run("single node", func(t *testing.T) {
		depErr := &DependencyError{Chain: []string{"lint"}, Err: fmt.Errorf("fail")}
		var buf bytes.Buffer
		PrintDependencyTree(depErr, &buf)
		out := buf.String()
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d: %v", len(lines), lines)
		}
		if !strings.HasPrefix(lines[0], "  lint") {
			t.Errorf("expected line to start with '  lint', got: %q", lines[0])
		}
		if !strings.Contains(out, "failed here") {
			t.Errorf("expected 'failed here' annotation on lint line, got: %q", out)
		}
	})

	t.Run("three nodes with correct indentation", func(t *testing.T) {
		depErr := &DependencyError{
			Chain: []string{"deploy", "build", "lint"},
			Err:   fmt.Errorf("fail"),
		}
		var buf bytes.Buffer
		PrintDependencyTree(depErr, &buf)
		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")

		if len(lines) != 3 {
			t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
		}
		// index 0 = 2 spaces, index 1 = 4 spaces, index 2 = 6 spaces (outermost first)
		checks := []struct {
			prefix    string
			name      string
			hasFailed bool
		}{
			{"  ", "deploy", false},
			{"    ", "build", false},
			{"      ", "lint", true},
		}
		for i, c := range checks {
			if !strings.HasPrefix(lines[i], c.prefix+c.name) {
				t.Errorf("line %d: want prefix %q + name %q, got %q", i, c.prefix, c.name, lines[i])
			}
			if c.hasFailed && !strings.Contains(lines[i], "failed here") {
				t.Errorf("line %d: expected 'failed here' annotation, got %q", i, lines[i])
			}
			if !c.hasFailed && strings.Contains(lines[i], "failed here") {
				t.Errorf("line %d: unexpected 'failed here' annotation on non-failing node, got %q", i, lines[i])
			}
		}
	})
}
