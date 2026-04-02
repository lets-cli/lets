package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

	t.Run("propagates real exit code from exec.ExitError", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "exit 2")
		runErr := cmd.Run()
		if runErr == nil {
			t.Fatal("expected command to fail")
		}
		execErr := &ExecuteError{err: runErr}
		depErr := &DependencyError{Chain: []string{"lint"}, Err: execErr}
		if depErr.ExitCode() != 2 {
			t.Errorf("expected exit code 2, got %d", depErr.ExitCode())
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

func TestDependencyErrorFailureMessage(t *testing.T) {
	t.Run("returns root cause for execute errors", func(t *testing.T) {
		depErr := &DependencyError{
			Chain: []string{"lint"},
			Err:   &ExecuteError{err: fmt.Errorf("failed to run command 'lint': %w", fmt.Errorf("exit status 1"))},
		}

		if got := depErr.FailureMessage(); got != "exit status 1" {
			t.Fatalf("expected root cause message, got %q", got)
		}
	})

	t.Run("keeps non execute errors intact", func(t *testing.T) {
		depErr := &DependencyError{
			Chain: []string{"lint"},
			Err:   fmt.Errorf("failed to calculate checksum for command 'lint': missing file"),
		}

		if got := depErr.FailureMessage(); got != "failed to calculate checksum for command 'lint': missing file" {
			t.Fatalf("expected original message, got %q", got)
		}
	})

	t.Run("keeps path context for execute errors", func(t *testing.T) {
		depErr := &DependencyError{
			Chain: []string{"lint"},
			Err: &ExecuteError{
				err: fmt.Errorf(
					"failed to run command 'lint': %w",
					&os.PathError{Op: "chdir", Path: "/tmp/missing", Err: syscall.ENOENT},
				),
			},
		}

		if got := depErr.FailureMessage(); got != "chdir /tmp/missing: no such file or directory" {
			t.Fatalf("expected path-aware message, got %q", got)
		}
	})
}

func TestPrintDependencyTree(t *testing.T) {
	t.Run("single node", func(t *testing.T) {
		depErr := &DependencyError{Chain: []string{"lint"}, Err: fmt.Errorf("fail")}
		var buf bytes.Buffer
		PrintDependencyTree(depErr, &buf)
		out := buf.String()
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
		}
		want := []string{
			dependencyTreeHeader,
			dependencyTreeIndent + dependencyTreeJoint + "lint" + dependencyTreeIndent + "<-- failed here",
		}

		for i := range want {
			if lines[i] != want[i] {
				t.Errorf("line %d: want %q, got %q", i, want[i], lines[i])
			}
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

		if len(lines) != 4 {
			t.Fatalf("expected 4 lines, got %d: %v", len(lines), lines)
		}
		want := []string{
			dependencyTreeHeader,
			dependencyTreeIndent + dependencyTreeJoint + "deploy",
			strings.Repeat(dependencyTreeIndent, 2) + dependencyTreeJoint + "build",
			strings.Repeat(dependencyTreeIndent, 3) + dependencyTreeJoint + "lint" +
				dependencyTreeIndent + "<-- failed here",
		}

		for i := range want {
			if lines[i] != want[i] {
				t.Errorf("line %d: want %q, got %q", i, want[i], lines[i])
			}
		}
	})
}
