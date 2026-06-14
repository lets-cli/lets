# Ubiquitous Language

## Lets configuration

| Term | Definition | Aliases to avoid |
| ---- | ---------- | ---------------- |
| **Settings** | Per-user lets behavior stored in `~/.config/lets/config.yaml`. | Config, `lets.yaml`, project config |
| **Project config** | Repository command and runtime configuration stored in `lets.yaml`. | Settings, user config |
| **Theme** | A named visual style for lets help and styled error output. | Color scheme, palette |
| **Default theme** | The standard lets theme. | Normal theme, builtin colors |
| **ANSI theme** | A theme limited to broadly supported ANSI terminal colors. | Plain theme, basic colors |
| **Synthwave theme** | A high-contrast neon lets theme. | Vaporwave, purple theme |

## Issue workflow

| Term | Definition | Aliases to avoid |
| ---- | ---------- | ---------------- |
| **Issue tracker** | The system of record where this repo's work items live. | Tickets, task board |
| **GitHub Issue** | A work item stored in `lets-cli/lets` on GitHub. | Local issue, scratch note |
| **Triage label** | A label string that marks an issue's triage state. | Tag, status |
| **Needs triage** | The triage state meaning a maintainer still needs to evaluate the issue. | Unreviewed, backlog |
| **Needs info** | The triage state meaning the issue is waiting on the reporter. | Blocked, pending |
| **Ready for agent** | The triage state meaning an AFK agent can execute the work without more human context. | AFK-ready, bot-ready |
| **Ready for human** | The triage state meaning the work requires human implementation. | Human-ready |
| **Won't fix** | The triage state meaning the work will not be actioned. | Rejected, closed-no-action |

## Documentation layout

| Term | Definition | Aliases to avoid |
| ---- | ---------- | ---------------- |
| **Domain docs** | The files that define project language and durable architectural decisions. | Repo docs, notes |
| **Context** | A `CONTEXT.md` file that defines the project's domain vocabulary for a scope. | Readme, design doc |
| **ADR** | An architecture decision record that captures a lasting technical decision. | Spec, random note |
| **Single-context repo** | A repo with one root `CONTEXT.md` and one root `docs/adr/`. | Monorepo, multi-context |
| **Multi-context repo** | A repo with a root `CONTEXT-MAP.md` and multiple per-context `CONTEXT.md` files. | Single-context |
| **Context map** | A root index file that points to per-context language docs. | Sitemap, overview |
| **Agent skill** | A reusable prompt-driven workflow that reads and writes repo-specific context. | Macro, script |

## People

| Term | Definition | Aliases to avoid |
| ---- | ---------- | ---------------- |
| **Maintainer** | A person who evaluates issues and applies triage labels. | Owner, reviewer |
| **Reporter** | A person who opens an issue or supplies follow-up information for it. | Submitter, user |
| **AFK agent** | An autonomous agent that can complete work without additional human context. | Bot, automation |
| **Human implementer** | A person needed when work cannot be delegated to an AFK agent. | Developer, coder |

## Relationships

- **Settings** apply to lets across all repositories, while **Project config** applies to one repository.
- A **Theme** belongs to **Settings** and is one of **Default theme**, **ANSI theme**, or **Synthwave theme**.
- A **GitHub Issue** lives in the **Issue tracker** and may have zero or more **Triage labels**.
- An issue should be **Ready for agent** only when a **Maintainer** has made it fully specified for an **AFK agent**.
- A **Single-context repo** has exactly one root **Context** and one shared root **ADR** directory.
- A **Multi-context repo** uses a **Context map** to point skills to the relevant **Context** files.

## Example dialogue

> **Dev:** "Should the new `theme` option live in **Settings** or the **Project config**?"
>
> **Domain expert:** "Put it in **Settings** — it changes lets itself, while the **Project config** in `lets.yaml` describes repo commands."
>
> **Dev:** "For this repo, does the **Issue tracker** mean **GitHub Issues** with **Triage labels**?"
>
> **Domain expert:** "Yes. A **Maintainer** evaluates each **GitHub Issue** and applies labels like **Needs triage** or **Ready for agent**."
>
> **Dev:** "And the repo is **Single-context**, so skills read one root **Context** and the root **ADR** directory?"
>
> **Domain expert:** "Exactly — no **Context map** is needed unless the repo grows into a **Multi-context repo**."

## Flagged ambiguities

- "settings" and "config" were used close together — use **Settings** for `~/.config/lets/config.yaml` and **Project config** for `lets.yaml`.
- "theme" and "color scheme" refer to the same user-facing concept here — use **Theme** in docs and CLI-facing language, and reserve "color scheme" for internal implementation details.
- "issue tracker" was used both abstractly and concretely — use **Issue tracker** for the system and **GitHub Issue** for an individual work item in this repo.
- "domain docs" was used for both the source-of-truth documents and the agent-facing instructions file — use **Domain docs** for `CONTEXT.md` and `docs/adr/`, and say "agent docs" when referring to `docs/agents/*`.
