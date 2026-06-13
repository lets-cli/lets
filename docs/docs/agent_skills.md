---
id: agent_skills
title: Agent Skills
---

Agent Skills are portable instructions that help AI agents discover how to work with a tool or project.

`lets` ships one bundled skill named `lets`. It explains how agents should inspect `lets.yaml`, discover available commands, prefer project-defined tasks, and safely modify lets configuration.

The feature is experimental and might change or be removed in a future release.

## What gets installed

The `lets` skill is installed as a standard Agent Skills directory:

```text
.agents/skills/lets/SKILL.md
```

For global installs, the same directory is created under your home directory:

```text
~/.agents/skills/lets/SKILL.md
```

Any compatible agent can discover the skill from those locations.

## Show the bundled skill

Print the bundled skill to stdout:

```bash
lets self skills show
```

Use this to inspect exactly what will be installed.

## Install the skill

Run the install command:

```bash
lets self skills install
```

Without flags, `lets` prompts you to choose local or global scope and shows the exact install path for each option.

Install for the current project:

```bash
lets self skills install --local
```

This writes to `.agents/skills/lets/` at the current Git repository root.

Install for the current user:

```bash
lets self skills install --global
```

This writes to `~/.agents/skills/lets/`.

Install to a custom skills directory:

```bash
lets self skills install --path /path/to/.agents/skills
```

Overwrite an existing installed skill:

```bash
lets self skills install --force
```

The optional skill name is accepted for compatibility:

```bash
lets self skills install lets
```

`lets` currently ships only the `lets` skill.

## Update the skill

Update installed copies to the version bundled in the current `lets` binary:

```bash
lets self skills update
```

You can also pass the skill name:

```bash
lets self skills update lets
```

Update checks the known local and global locations:

- `.agents/skills/lets/` at the current Git repository root
- `~/.agents/skills/lets/`

Skills installed with `--path` are not discovered by `update`; reinstall with `--path --force` to refresh a custom location.

## Remove the skill

There is no dedicated remove command. Delete the installed skill directory.

Remove the local project skill:

```bash
rm -rf .agents/skills/lets
```

Remove the global user skill:

```bash
rm -rf ~/.agents/skills/lets
```

After removal, compatible agents will stop discovering the bundled `lets` skill from that location.
