# lets

@README.md

## Tools

Use `lets` task runner for all build/test/lint operations instead of raw commands. Run `lets build` first if binary is missing.

```bash
lets build [bin]              # build CLI with version metadata
lets build-and-install        # build and install lets binary locally
lets test                     # full suite: unit + bats + completions
lets test-unit                # Go unit tests only
lets test-bats [test]         # Docker-based Bats integration tests
lets lint                     # golangci-lint via Docker
lets fmt                      # go fmt ./...
lets coverage [--html]        # coverage report
lets run-docs                 # local docs dev server (docs/)
lets publish-docs             # deploy docs site
```

`lets test-unit`, `lets test-bats`, and `lets lint` require Docker. Use `go test ./...` locally for quick iteration without Docker.

## Agent skills

### Issue tracker

The Issue tracker for this repo is GitHub; work items live as GitHub Issues in `lets-cli/lets`. See `docs/agents/issue-tracker.md`.

### Triage labels

Triage uses the default Triage label names: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, and `wontfix`. See `docs/agents/triage-labels.md`.

### Domain docs

This is a Single-context repo; skills should read the root Context and root ADRs. See `docs/agents/domain.md`.

## Agent Behavior

- **Proactive execution** — Don't ask "Can I proceed?" for implementation. DO ask before changing success criteria, test thresholds, or what "working" means.
- **Test early, test real** — Don't accumulate 10 changes then debug. After each logical step: does it work? With realistic input, not just edge case that triggered the work.
- **Pushback** — Propose alternatives before implementing suboptimal approaches. Ask about design choices.
- **Unify, don't duplicate** — Merge nearly-identical structs/functions rather than adding variants.
- **No over-engineering** — Minimum complexity for current task. No speculative abstractions.
- **Terseness** — Comments for surprising/hairy logic only. Be extremely concise in communication.

## Package Structure

- `cmd/lets/main.go` — CLI entry point, flag parsing, signal handling
- `internal/cmd/` — Cobra command setup (root, subcommands, completion, LSP, self-update)
- `internal/config/` — config file discovery, loading, validation; `internal/config/config/` defines Config/Command/Mixin structs and YAML unmarshaling; `internal/config/path/` contains config path helpers
- `internal/executor/` — command execution, dependency resolution, env setup, checksum verification
- `internal/env/` — debug level state (`LETS_DEBUG`, levels 0-2)
- `internal/logging/` — logrus-based logging with command chain formatting
- `internal/lsp/` — Language Server Protocol: definition lookup, completion for depends, tree-sitter YAML parsing; `lets lsp` runs stdio-based server for IDE integration
- `internal/checksum/` — SHA1 file checksumming with glob patterns
- `internal/docopt/` — docopt argument parsing, produces `LETSOPT_*` and `LETSCLI_*` env vars
- `internal/upgrade/` — binary self-update from GitHub releases; `internal/upgrade/registry/` contains release registry implementation
- `internal/util/` — file/dir/version helpers
- `internal/workdir/` — `--init` scaffolding
- `internal/set/` — generic Set data structure
- `internal/test/` — test utilities (temp files, args helpers)

## Project Rules

- Follow `gofmt` exactly; tabs for indentation, ~120 char lines
- Unit tests as `*_test.go` next to source; Bats tests in `tests/*.bats`
- Fixtures in matching `tests/<scenario>/` folder, use `lets.yaml` unless variant needed
- Bats tests use `run` + `assert_success`/`assert_line` pattern
- Run at least `go test ./...` before considering work complete; `lets test-bats` for CLI-path changes
- Run `lets lint` to verify code quality before commit/push/PR creation
- If you discover non-obvious knowledge needed to make something work or avoid a known issue, document it in code comments or docs so it is not lost
- Add concise code comments for non-obvious logic, invariants, and surprising decisions; do not comment self-explanatory code
- **Golden tests** — `internal/cmd/testdata/*` are snapshot of the rendered help and error output. If you change anything that affects help or error rendering (flags, styles, section titles, error messages), regenerate them with `go test ./internal/cmd/ -run -update` (or `lets test-unit --update-golden` in Docker), then commit the updated `.golden` files. If you add a new rendering behaviour (new section, new error type, new command layout), add a corresponding golden test in `internal/cmd/help_golden_test.go` with a fixture YAML in `internal/cmd/testdata/fixtures/` if needed, then run with `-update` to create the golden file.
- Commits: short imperative subjects (`Add ...`, `Fix ...`, `Use ...`), explain non-obvious context in body
- **Changelog workflow**: add entries to the `Unreleased` section in `docs/docs/changelog.md` with each commit/PR. At release time, rename `Unreleased` to the new tag version
- Do not commit `lets.my.yaml`, generated binaries, `.lets/`, `coverage.out`, or `node_modules`
- CLI flags: kebab-case only (`--dry-run` not `--dry_run`)
- No "Generated by <agent>" in commits
