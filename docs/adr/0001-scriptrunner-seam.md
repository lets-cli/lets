# ADR-0001 — Introduce ScriptRunner seam in executor

**Date:** 2026-06-14  
**Status:** Accepted

## Context

`internal/executor/` orchestrates the full lifecycle of a **Project command** — docopt parsing, **Dependency chain** resolution, **Environment** layering, **Init script**, **After script**, and **Persisted checksum** — but it had no testable boundary between that orchestration and OS process spawning. Every function that needed to verify execution semantics had to spawn a real shell, making the unit test suite dependent on Docker and slow.

The two execution paths (`execute` and `executeParallel`) also duplicated their five phases verbatim, so any change had to be applied twice.

## Decision

Introduce a `ScriptRunner` function type as the seam between orchestration and OS process spawning:

```go
type ScriptRunner func(command *config.Command, script string) error
```

`NewExecutor` accepts a `ScriptRunner` parameter. The production implementation, `NewShellRunner`, wraps the shell-exec logic (script preparation, stdio wiring, **Work dir** resolution, **Environment** layering). Tests inject a closure — no test struct or mock framework required.

`Executor` no longer holds an `io.Writer`; output wiring belongs entirely to the runner.

## Consequences

- **Positive:** `Execute()` is now unit-testable without a real shell — inject a recording closure and assert on invocations and errors.
- **Positive:** All direct `os/exec` calls are confined to `shellRunner`; the orchestrator is pure Go logic.
- **Positive:** Enables a future dry-run mode by swapping in a no-op `ScriptRunner`.
- **Neutral:** Callers pass `executor.NewShellRunner(conf, out)` instead of `out` directly — a one-line change at each call site.
- **Neutral:** Debug-level-2 env logging moved into the runner; the orchestrator logs the script only.
