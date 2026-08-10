# ADR-0004 — The Root dir is the invocation dir

**Date:** 2026-08-10
**Status:** Accepted

## Context

`lets` had never specified which directory a **Project command** runs in. The behavior
that shipped was an accident of implementation, and it changed silently in `0.0.63`.

Up to `0.0.62`, `Command.WorkDir` was assigned unconditionally from
`filepath.Abs(cmd.WorkDir)`. For the overwhelming majority of commands — those without a
`work_dir` — that argument was `""`, and `filepath.Abs("")` returns the process cwd. The
executor then preferred `Command.WorkDir` over `Config.WorkDir` whenever it was non-empty,
which it now always was. So commands ran in the invocation dir, and `Config.WorkDir` — the
config file's directory, computed since 2020 — was dead code for its entire life.

`0.0.63` added an `if cmd.WorkDir != ""` guard while fixing checksum handling for
`work_dir`. That was correct in isolation, but it un-shadowed `Config.WorkDir` and moved
every command's working directory to the config file's directory.

Neither version was internally consistent. `Config.WorkDir` was read by some directives
and not others, so a single command definition resolved paths against two directories:

| | `0.0.62` | `0.0.63` |
| --- | --- | --- |
| `cmd` cwd | invocation dir | config dir |
| `env.sh` cwd | invocation dir | invocation dir |
| `checksum` | config dir | config dir / `work_dir` |
| `env_file` | config dir | config dir |
| `work_dir` base | invocation dir | invocation dir |

On `0.0.62`, `checksum: [data.txt]` and `cmd: cat data.txt` in the same command could
refer to different files. On `0.0.63`, `cmd` and `env.sh` disagreed instead. The root
cause of both was one field, `Config.WorkDir`, carrying two meanings: "where the config
is" and "where commands run".

**Remote configs** (ADR-0003) already special-cased this by pinning their root to the
invocation dir, since a cached config has no meaningful local directory. That special case
was evidence that the invocation dir was the right general answer.

## Decision

Split the conflated field into two, and define resolution in terms of when it happens.

- **Root dir** — the directory `lets` was invoked from. Commands run here by default.
- **Config dir** — the directory holding the config file.
- **Work dir** — a command's actual directory: the Root dir, or its `work_dir` if set.

Two rules:

1. **Config assembly** resolves against the config file that declares it. Local `mixins`
   paths are the only thing in this category.
2. **Everything a command reads or runs** resolves against that command's Work dir. This
   covers `cmd`, `checksum` paths, `env_file` paths at both global and command scope, and
   `env.sh` scripts. A relative `work_dir` resolves against the Root dir.

`--config` / `-c`, `--config-dir` and `LETS_CONFIG_DIR` select which config is loaded and
never change the Root dir.

`.lets/` is created in the Root dir, so a persisted checksum stays paired with the files it
was computed from.

Remote configs stop being a special case: their root is the invocation dir under the
general rule. Because their Config dir is a cache directory holding only the downloaded
YAML, a remote config declaring a local `mixins` path is rejected with an explicit error.

## Consequences

- `lets foo` does what typing `foo` at the prompt would do, and a command is readable
  without knowing where its config lives.
- `lets -c ../other/lets.yaml build` borrows another project's commands and runs them on
  the invoking project's files, which is the only useful reading of `-c`.
- A filename appearing in two directives of one command now means one file.
- Breaking relative to `0.0.63`: commands run in the invocation dir again.
- Breaking relative to every previous version: `checksum` and `env_file` follow the Work
  dir, and `.lets/` follows the Root dir.
- A relative `work_dir` now depends on where the user stands, so `work_dir: docs` works
  from the project root and fails from a subdirectory. Commands that must always target the
  project use `cd "${LETS_CONFIG_DIR}/…"` inside `cmd`. `work_dir` does not expand
  environment variables.
- `LETS_CONFIG_DIR` for a remote config reports the cache directory rather than the cwd.
  `$PWD` is the way to reach the project.

## Alternatives considered

**Keep `0.0.63` — Root dir is the Config dir.** Every path in a config would mean the same
thing regardless of where `lets` ran, and `lets x` would be identical from any directory.
Rejected because it makes "act on where I am" inexpressible: nothing exposes the invocation
dir, so it would have required a new `LETS_INVOCATION_DIR`. The chosen model needs no new
API — `$LETS_CONFIG_DIR` already provides the inverse escape hatch and always has.

**Make everything cwd-relative, including mixins.** Fully uniform, and the honest version
of pre-`0.0.63` behavior. Rejected because `mixins: [./common.yaml]` would break whenever
`lets` ran from a subdirectory, making configs with mixins unloadable from anywhere but
their own directory.

**Leave `checksum` and `env_file` resolving against the Config dir.** Closest to a literal
revert of `0.0.63`. Rejected because it preserves the original defect — one command
definition resolving the same filename against two directories.

**Put `.lets/` next to the config rather than in the Root dir.** Avoids stray `.lets/`
directories when running from subdirectories. Rejected because a persisted checksum stored
next to the config but computed from an arbitrary directory would describe different files
on different runs, making change detection unreliable.
