# Dependency Failure Tree

**Date:** 2026-03-13
**Status:** Design approved — ready for implementation

## Problem

When `lets` runs a command with a `depends` chain and a dependency fails, the current error output only names the immediate failing command:

```
failed to run command 'lint': exit status 1
```

There is no indication of how deep in the dependency chain the failure occurred or which parent commands triggered it.

## Goal

Print an indented tree of the full dependency chain whenever a command fails, so the user immediately knows where in the chain execution stopped. This applies to all failures — with or without a `depends` chain.

```
  deploy
    build
      lint  <-- failed here

ERRO[0000] failed to run command 'lint': exit status 1
```

## Design

### Data Model

A new `DependencyError` type in `internal/executor/dependency_error.go`:

```go
type DependencyError struct {
    Chain []string // outermost-first: ["deploy", "build", "lint"]
    Err   error    // the original ExecuteError
}

func (e *DependencyError) Error() string { return e.Err.Error() }

func (e *DependencyError) ExitCode() int {
    var exitErr *ExecuteError
    if errors.As(e.Err, &exitErr) {
        return exitErr.ExitCode()
    }
    return 1
}
```

A `prependToChain(name string, err error) error` helper:
- If `err` is already a `*DependencyError`: prepend `name` to `Chain` and return the same error.
- Otherwise: return `&DependencyError{Chain: []string{name}, Err: err}`.

### Chain Construction

Chain building happens in `execute()` and `executeParallel()`. **All error return paths** in both functions call `prependToChain` before returning. This includes errors from:
- `initCmd`
- `executeDepends`
- `runCmd` (each iteration)
- `persistChecksum`

**Call stack trace for `deploy → build → lint` (lint fails):**

1. `Execute(deploy)` → `execute(deploy)` → `executeDepends` calls `Execute(build)` →
2. `Execute(build)` → `execute(build)` → `executeDepends` calls `Execute(lint)` →
3. `Execute(lint)` → `execute(lint)` → `executeDepends` (empty, returns nil) → `runCmd(lint)` → returns `ExecuteError`
4. `execute(lint)` calls `prependToChain("lint", ExecuteError)` → creates `DependencyError{Chain: ["lint"], Err: ExecuteError}`
5. Bubbles up to `execute(build)` via `executeDepends`; calls `prependToChain("build", DependencyError)` → `DependencyError{Chain: ["build", "lint"], Err: ExecuteError}`
6. Bubbles up to `execute(deploy)` via `executeDepends`; calls `prependToChain("deploy", DependencyError)` → `DependencyError{Chain: ["deploy", "build", "lint"], Err: ExecuteError}`

**Single-command failure (no depends):** `runCmd` fails in `execute`, which calls `prependToChain("lint", ExecuteError)` → `DependencyError{Chain: ["lint"], Err: ExecuteError}`. Renders as a single-node tree.

**Global `init` script failure:** The `Execute()` method runs `cfg.Init` before dispatching to `execute`/`executeParallel`. If this fails, it returns the raw error — no tree is shown. This is intentional: the global init is not tied to any command name.

### Rendering

```go
// PrintDependencyTree writes an indented tree of the dependency chain to w.
// The failing node (last in chain) is annotated in red.
// fatih/color v1.16.0 automatically respects NO_COLOR.
func PrintDependencyTree(e *DependencyError, w io.Writer)
```

- 2 spaces of indentation per depth level (level 0 = 2 spaces, level 1 = 4 spaces, …).
- The last node gets `  <-- failed here` in red via `fatih/color`. `fatih/color` v1.16.0 automatically respects `NO_COLOR`.

Example — 3-level chain:

```
  deploy
    build
      lint  <-- failed here
```

Example — single node:

```
  lint  <-- failed here
```

### Error Handler in main.go

Tree is printed **before** `log.Error` so context appears above the error line:

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

`DependencyError` implements `ExitCode()` so `getExitCode` propagates the correct exit code from the innermost failing process. All execution errors — from both `execute` and `executeParallel` — bubble up through this single handler.

Complete example of final stderr output for the 3-level chain:

```
  deploy
    build
      lint  <-- failed here
ERRO[0000] failed to run command 'lint': exit status 1
```

### Files Changed

| File | Change |
|------|--------|
| `internal/executor/dependency_error.go` | New: `DependencyError` type, `prependToChain` helper, `PrintDependencyTree` |
| `internal/executor/executor.go` | Call `prependToChain` on all error return paths in `execute()` and `executeParallel()` |
| `cmd/lets/main.go` | Detect `*DependencyError` and call `PrintDependencyTree` before `log.Error` |

### Out of Scope

- Global `init` script failures (no command name to display).
- Parallel `cmd` arrays within a single command (not `depends`).
- No changes to debug logging or `ExecLogger`.
- No new CLI flags.

## Testing

### Unit Tests

In `internal/executor/dependency_error_test.go`:

- `prependToChain` with non-`DependencyError` input → creates single-node `DependencyError`.
- `prependToChain` with existing `DependencyError` → prepends correctly.
- `DependencyError.ExitCode()` propagates exit code from inner `ExecuteError`.
- `DependencyError.ExitCode()` returns 1 when inner error has no exit code.
- `PrintDependencyTree` output matches expected indentation for 1, 2, 3-node chains.

### Bats Integration Test

`tests/dependency_failure_tree.bats` with fixture `tests/dependency_failure_tree/lets.yaml`:

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

Tests run with `NO_COLOR=1` to avoid ANSI escape codes in assertions:

```bash
@test "dependency failure tree shows full chain" {
    cd "$BATS_TEST_DIRNAME/dependency_failure_tree"
    run env NO_COLOR=1 lets deploy
    assert_failure
    assert_line --index 0 "  deploy"
    assert_line --index 1 "    build"
    assert_line --index 2 --partial "      lint"
    assert_line --index 2 --partial "failed here"
}
```
