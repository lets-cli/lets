# Domain Docs

This repo is a **Single-context repo**.

For this repo, the **Domain docs** are:

- the root **Context** at `CONTEXT.md`
- the root **ADR** directory at `docs/adr/`

`docs/agents/*` are agent docs. They describe workflow conventions for **Agent skills**, but they are not Domain docs.

## Before exploring, read these

- **`CONTEXT.md`** at the repo root
- **`docs/adr/`** at the repo root — read ADRs that touch the area you're about to work in

If these files don't exist, **proceed silently**. Don't flag their absence; don't suggest creating them upfront. The producer skill (`/grill-with-docs`) creates them lazily when terms or decisions actually get resolved.

## File structure

Single-context repo layout:

```
/
├── CONTEXT.md
├── docs/adr/
│   ├── 0001-some-decision.md
│   └── 0002-another-decision.md
└── src/
```

Do not expect a `CONTEXT-MAP.md` unless this repo grows into a **Multi-context repo**.

## Use the Context vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, or a test name), use the term as defined in the root **Context**. Don't drift to synonyms the glossary explicitly avoids.

If the concept you need isn't in the glossary yet, that's a signal — either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/grill-with-docs`).

## Flag ADR conflicts

If your output contradicts an existing **ADR**, surface it explicitly rather than silently overriding:

> _Contradicts ADR-0007 — but worth reopening because…_
