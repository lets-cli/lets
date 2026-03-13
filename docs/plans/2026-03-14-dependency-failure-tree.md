# Dependency Failure Tree Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** When a `lets` command (or any of its `depends`) fails, print an indented tree to stderr showing the full dependency chain with the failing node highlighted in red.

**Architecture:** A new `DependencyError` type carries the command chain as `[]string`. The `execute()` and `executeParallel()` functions in the executor wrap every error return with `prependToChain`, which builds the chain bottom-up as the error bubbles up. `main.go` detects this type and renders the tree before printing the error message.

**Tech Stack:** Go stdlib `errors`, `fmt`, `io`, `strings`; `github.com/fatih/color` v1.16.0 (already in `go.mod`) for red highlight; `testing` package for unit tests; bats-core + bats-assert for integration tests.

**Spec:** `docs/specs/2026-03-13-dependency-failure-tree-design.md`

---

## Chunk 1: DependencyError type, helpers, and unit tests

### Task 1: Create `dependency_error.go` with TDD

**Files:**
- Create: `internal/executor/dependency_error.go`
- Create: `internal/executor/dependency_error_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/executor/dependency_error_test.go`:

```go
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
	t.Run("wraps ExecuteError exit code", func(t *testing.T) {
		// Use a plain error in ExecuteError (not exec.ExitError) — ExitCode() returns 1
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
		if !strings.Contains(out, "  lint") {
			t.Errorf("expected '  lint' in output, got: %q", out)
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
```

- [ ] **Step 2: Run tests to confirm they fail**

```bash
go test ./internal/executor/ -run "TestPrependToChain|TestDependencyError|TestPrintDependencyTree" -v
```

Expected: compile error — `prependToChain`, `DependencyError`, `PrintDependencyTree` not defined.

- [ ] **Step 3: Implement `dependency_error.go`**

Create `internal/executor/dependency_error.go`:

```go
package executor

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

// DependencyError carries the full dependency chain when a command fails.
// Chain is outermost-first: e.g. ["deploy", "build", "lint"].
type DependencyError struct {
	Chain []string
	Err   error
}

func (e *DependencyError) Error() string { return e.Err.Error() }

// ExitCode propagates the exit code from the innermost ExecuteError, or returns 1.
func (e *DependencyError) ExitCode() int {
	var exitErr *ExecuteError
	if errors.As(e.Err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

// prependToChain prepends name to the chain in err if err is already a *DependencyError,
// otherwise wraps err in a new single-element DependencyError.
func prependToChain(name string, err error) error {
	var depErr *DependencyError
	if errors.As(err, &depErr) {
		depErr.Chain = append([]string{name}, depErr.Chain...)
		return depErr
	}

	return &DependencyError{Chain: []string{name}, Err: err}
}

// PrintDependencyTree writes an indented tree of the dependency chain to w.
// The failing node (last in chain) is annotated in red.
// Respects NO_COLOR automatically via fatih/color.
func PrintDependencyTree(e *DependencyError, w io.Writer) {
	red := color.New(color.FgRed).SprintFunc()

	for i, name := range e.Chain {
		indent := strings.Repeat("  ", i+1)
		if i == len(e.Chain)-1 {
			fmt.Fprintf(w, "%s%s  %s\n", indent, name, red("<-- failed here"))
		} else {
			fmt.Fprintf(w, "%s%s\n", indent, name)
		}
	}
}
```

- [ ] **Step 4: Run tests to confirm they pass**

```bash
go test ./internal/executor/ -run "TestPrependToChain|TestDependencyError|TestPrintDependencyTree" -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/executor/dependency_error.go internal/executor/dependency_error_test.go
git commit -m "Add DependencyError type with chain tracking and tree rendering"
```

---

## Chunk 2: Wire executor and main, add integration test

### Task 2: Update `executor.go` to wrap errors with `prependToChain`

**Files:**
- Modify: `internal/executor/executor.go`

- [ ] **Step 1: Update `execute()` — all error return paths**

In `internal/executor/executor.go`, replace `execute()` (lines 92–121):

```go
func (e *Executor) execute(ctx *Context) error {
	command := ctx.command

	if env.DebugLevel() > 1 {
		ctx.logger.Debug("command %s", command.Dump())
	}

	defer func() {
		if command.After != "" {
			e.executeAfterScript(ctx)
		}
	}()

	if err := e.initCmd(ctx); err != nil {
		return prependToChain(command.Name, err)
	}

	if err := e.executeDepends(ctx); err != nil {
		return prependToChain(command.Name, err)
	}

	for _, cmd := range command.Cmds.Commands {
		if err := e.runCmd(ctx, cmd); err != nil {
			return prependToChain(command.Name, err)
		}
	}

	// persist checksum only if exit code 0
	if err := e.persistChecksum(ctx); err != nil {
		return prependToChain(command.Name, err)
	}

	return nil
}
```

- [ ] **Step 2: Update `executeParallel()` — all error return paths**

Replace `executeParallel()` (lines 362–399):

```go
func (e *Executor) executeParallel(ctx *Context) error {
	command := ctx.command

	defer func() {
		if command.After != "" {
			e.executeAfterScript(ctx)
		}
	}()

	if err := e.initCmd(ctx); err != nil {
		return prependToChain(command.Name, err)
	}

	if err := e.executeDepends(ctx); err != nil {
		return prependToChain(command.Name, err)
	}

	group, _ := errgroup.WithContext(ctx.ctx)

	for _, cmd := range command.Cmds.Commands {
		cmd := cmd
		group.Go(func() error {
			return e.runCmd(ctx, cmd)
		})
	}

	if err := group.Wait(); err != nil {
		return prependToChain(command.Name, err) //nolint:wrapcheck
	}

	// persist checksum only if exit code 0
	if err := e.persistChecksum(ctx); err != nil {
		return prependToChain(command.Name, err)
	}

	return nil
}
```

- [ ] **Step 3: Run unit tests to confirm nothing is broken**

```bash
go test ./...
```

Expected: all tests PASS (no new tests here — we're verifying no regressions).

- [ ] **Step 4: Commit**

```bash
git add internal/executor/executor.go
git commit -m "Wrap executor errors with prependToChain for dependency tree tracking"
```

---

### Task 3: Update `main.go` to print the tree on failure

**Files:**
- Modify: `cmd/lets/main.go`

- [ ] **Step 1: Update the error handler in `main()`**

In `cmd/lets/main.go`, replace lines 120–123:

```go
// before:
if err := rootCmd.ExecuteContext(ctx); err != nil {
    log.Error(err.Error())
    os.Exit(getExitCode(err, 1))
}
```

With:

```go
if err := rootCmd.ExecuteContext(ctx); err != nil {
    var depErr *executor.DependencyError
    if errors.As(err, &depErr) {
        executor.PrintDependencyTree(depErr, os.Stderr)
    }
    log.Error(err.Error())
    os.Exit(getExitCode(err, 1))
}
```

Add `"errors"` and `"github.com/lets-cli/lets/internal/executor"` to the import block if not already present.

- [ ] **Step 2: Build to confirm it compiles**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 3: Smoke test manually**

Create a quick temp `lets.yaml` (or use an existing fixture) with a failing depends chain and run `lets`. Confirm the tree appears above the error line.

- [ ] **Step 4: Commit**

```bash
git add cmd/lets/main.go
git commit -m "Print dependency failure tree in main.go error handler"
```

---

### Task 4: Add bats integration test

**Files:**
- Create: `tests/dependency_failure_tree/lets.yaml`
- Create: `tests/dependency_failure_tree.bats`

- [ ] **Step 1: Create the fixture**

Create `tests/dependency_failure_tree/lets.yaml`:

```yaml
shell: bash
commands:
  deploy:
    depends: [build]
    cmd: echo done
  build:
    depends: [lint]
    cmd: echo done
  lint:
    cmd: exit 1
```

- [ ] **Step 2: Create the bats test**

Create `tests/dependency_failure_tree.bats`:

```bash
load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/dependency_failure_tree
}

@test "dependency_failure_tree: shows full 3-level chain on failure" {
    run env NO_COLOR=1 lets deploy
    assert_failure
    assert_line --index 0 "  deploy"
    assert_line --index 1 "    build"
    assert_line --index 2 --partial "      lint"
    assert_line --index 2 --partial "failed here"
}

@test "dependency_failure_tree: single node when no depends" {
    run env NO_COLOR=1 lets lint
    assert_failure
    assert_line --index 0 --partial "  lint"
    assert_line --index 0 --partial "failed here"
}
```

- [ ] **Step 3: Run the bats tests**

```bash
lets test-bats dependency_failure_tree
```

Expected: both tests PASS.

- [ ] **Step 4: Run the full test suite**

```bash
go test ./...
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add tests/dependency_failure_tree/ tests/dependency_failure_tree.bats
git commit -m "Add bats integration test for dependency failure tree"
```

---

### Task 5: Update changelog

**Files:**
- Modify: `docs/docs/changelog.md`

- [ ] **Step 1: Add entry to Unreleased section**

Open `docs/docs/changelog.md` and add under the `Unreleased` section:

```markdown
* `[Added]` When a command or its `depends` chain fails, print an indented tree to stderr showing the full chain with the failing command highlighted
```

- [ ] **Step 2: Commit**

```bash
git add docs/docs/changelog.md
git commit -m "Update changelog for dependency failure tree feature"
```
