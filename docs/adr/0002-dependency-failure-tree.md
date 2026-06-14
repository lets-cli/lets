# ADR-0002 — Render Project command failures as a dependency tree

**Date:** 2026-06-14
**Status:** Accepted

## Context

This decision was originally captured as a design spec on 2026-03-13 and later promoted to an ADR.

When `lets` ran a **Project command** with a `depends` chain and a dependency failed, the error output only named the innermost failing command:

```text
failed to run command 'lint': exit status 1
```

Users could not see which parent commands triggered that failure or how deep in the **Dependency chain** execution stopped. The same ambiguity existed for single-command failures because the error carried no structured command context.

The fix needed to work for both serial and parallel command execution paths without changing the child process exit code that `lets` returns.

## Decision

Carry command-chain context through executor errors with a `DependencyError` type in `internal/executor`.

`DependencyError` stores:

- `Chain []string`: the failing **Dependency chain**, outermost-first, for example `[]string{"deploy", "build", "lint"}`.
- `Err error`: the original execution error.

It delegates `Error()` and `Unwrap()` to the original error, and its `ExitCode()` preserves the exit code from the innermost `ExecuteError` when one exists.

Add a `prependToChain(name, err)` helper at command boundaries:

- If `err` is already a `DependencyError`, return a new `DependencyError` with `name` prepended and the original cause preserved.
- Otherwise, wrap `err` in a single-node `DependencyError`.

Both `execute()` and `executeParallel()` call `prependToChain` for command-scoped error paths:

- command initialization (`initCmd`)
- dependency execution (`executeDepends`)
- command script execution (`runCmd`)
- persisted-checksum writes (`persistChecksum`)

The root **Init script** is intentionally excluded because it is not tied to a **Project command** name.

Render the chain in the CLI error renderer as a themed `command tree:` section, with the failing leaf annotated:

```text
command tree:
    └─ deploy
      └─ build
        └─ lint  <-- failed here
```

For non-terminal output, keep the plain error text path.

## Consequences

- **Positive:** Users see the full command context for dependency failures and single-command failures.
- **Positive:** The original error remains unwrap-able, so existing error handling and exit-code propagation continue to work.
- **Positive:** Serial and parallel executor paths share the same command-boundary wrapping rule.
- **Neutral:** Command execution failures now surface as `DependencyError` values at higher layers.
- **Neutral:** Failures inside a parallel `cmd` array are attributed to the owning **Project command**, not to an individual shell fragment.
- **Negative:** The error renderer now depends on executor-specific error structure to show the command tree.

## Related implementation

- `internal/executor/dependency_error.go`
- `internal/executor/executor.go`
- `internal/cmd/error.go`
- `internal/executor/dependency_error_test.go`
- `internal/cmd/help_golden_test.go`
- `tests/dependency_failure_tree.bats`
