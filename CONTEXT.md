# Context

This is a single-context repo.

## Purpose

`lets` is a YAML-based CLI task runner. A repository defines project workflow in `lets.yaml`, and `lets` turns that declaration into runnable **Project commands**, built-in **Self commands**, help output, shell completion, and editor tooling. In marketing copy "task runner" is fine; in precise product language prefer **Project command** to "task".

## Language choices

- Use **Settings** for `~/.config/lets/config.yaml`.
- Use **Project config** for `lets.yaml`.
- Use **Project command** for a user-defined runnable unit.
- Use **Self command** for `lets self ...` functionality.
- Use **Theme** for user-facing output styling; reserve "color scheme" for implementation details.

## Configuration model

| Term | Definition | Aliases to avoid |
| ---- | ---------- | ---------------- |
| **Project config** | The repository-level declaration in `lets.yaml` that defines commands and shared execution rules. | Settings, user config |
| **Remote config** | A main Project config loaded from a URL instead of from a local `lets.yaml` file. | Remote root config, URL config |
| **Settings** | Per-user lets behavior stored in `~/.config/lets/config.yaml`. | Project config |
| **Mixin** | An additional config file merged additively into a Project config. | Override, patch |
| **Remote mixin** | A Mixin loaded from a URL and merged into the main Project config. | Remote config, URL mixin |
| **Global env** | Top-level environment entries shared by all Project commands. | Process env |
| **Env file** | A dotenv-style file loaded into command execution at global or command scope. | Settings file, config file |
| **Theme** | A named style for lets help and styled error output. | Color scheme, palette |

## Command model

| Term | Definition | Aliases to avoid |
| ---- | ---------- | ---------------- |
| **Project command** | A named runnable unit declared under `commands` in Project config. | Task, job |
| **Self command** | A built-in `lets self ...` command used to manage lets itself. | Project command |
| **Option** | A docopt-declared CLI input for a Project command. | Cobra flag, raw env var |
| **Dependency chain** | The ordered graph of Project commands reached through `depends`. | Pipeline, hook |
| **Reference command** | A Project command that reuses another command via `ref` and optional `args`. | Alias |
| **Command group** | A help-only label used to organize Project commands in rendered help. | Execution stage |
| **Environment** | The fully resolved environment passed to Project command execution. | Process env |
| **Checksum** | A SHA1 digest derived from configured files and exposed as environment variables. | Cache key |
| **Persisted checksum** | Stored checksum state used to detect change between invocations. | Build cache |

## Execution model

| Term | Definition | Aliases to avoid |
| ---- | ---------- | ---------------- |
| **Init script** | A top-level script run once per lets invocation before the first Project command executes. | Before script |
| **Before script** | A top-level script prepended to each Project command invocation, including dependencies. | Init script |
| **After script** | A command-scoped script run after a Project command execution attempt. | Cleanup hook |
| **Work dir** | The directory where a Project command runs after config and command resolution. | Repo root |
| **Download progress indicator** | A user-visible status shown while lets retrieves a Remote config or Remote mixin. | Progress bar |
| **Help surface** | The rendered CLI help for root and Project commands. | Docs page |
| **LSP surface** | The editor-facing language-server features exposed by `lets self lsp`. | CLI help |

## Relationships

- One **Project config** defines zero or more **Project commands**.
- A **Project command** inherits shell and shared execution rules from **Project config** unless it overrides them.
- A **Dependency chain** contains only **Project commands**; **Self commands** are outside that graph.
- A **Reference command** points to exactly one target **Project command**.
- An **Option** parsed for a **Project command** is exposed as `LETSOPT_*` and `LETSCLI_*` variables in the command **Environment**.
- A **Checksum** belongs to one **Project command** and may expose one aggregate digest plus named digests.
- **Settings** affect lets itself, while **Project config** affects repository workflow.
- **Mixins** add commands, env, before scripts, and env files; conflicting names are errors rather than implicit overrides.
- A **Command group** changes help organization only; it does not change execution behavior.

## Execution flow

1. Discover **Project config** from `--config`, `LETS_CONFIG`, or upward search for `lets.yaml`.
2. Load **Mixins** and merge them into one effective **Project config**.
3. Validate config structure, supported keywords, version requirements, references, and dependency graph.
4. Build CLI surfaces from **Project commands** plus built-in **Self commands**.
5. Parse **Option** values with docopt and expose them to the command **Environment**.
6. Resolve **Environment** from built-in lets vars, config env, env files, parsed options, explicit overrides, and checksum vars.
7. Run the **Init script** once per invocation.
8. Execute the **Dependency chain** before the requested **Project command**.
9. Run Project command scripts in the selected shell and **Work dir**.
10. Run the **After script** after the execution attempt; its own failure does not replace the main command result.
11. Persist **Persisted checksum** state only after successful execution.

## Precedence and boundaries

- **Settings** precedence is: environment variables > settings file > built-in defaults.
- **Environment** precedence at command runtime is: process env < built-in lets vars < global `env` < global `env_file` < command `env` < command `env_file` < parsed options / explicit env overrides / checksum vars.
- `NO_COLOR` disables color even if **Settings** choose a **Theme**.
- A **Theme** affects lets output itself, not the child processes launched by **Project commands**.
- **Project config** must declare `shell`; **Mixins** do not.
- `persist_checksum` is valid only when `checksum` is declared.
- Unknown top-level keywords fail config loading unless they use the `x-` prefix for custom extensions.
- **Agent skills** are installed and managed with `lets self skills`; they are not configured in **Project config**.

## Non-goals

- `lets` orchestrates shell commands; it does not replace the shell itself.
- `lets` is not a package manager.
- **Settings** are not a place to share repository workflow.
- **Self commands** are not part of a repository's **Project command** namespace.
- **Mixins** are additive composition, not implicit override layers.

## Related docs

- `UBIQUITOUS_LANGUAGE.md` is the working glossary and ambiguity log.
- Promote stable terminology from `UBIQUITOUS_LANGUAGE.md` into this file when it becomes durable project language.
