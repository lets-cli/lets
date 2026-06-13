---
name: lets
description: >-
  Use lets, the YAML-based CLI task runner. Activate when a repository has
  lets.yaml or the user asks to run, inspect, add, or debug lets tasks,
  commands, dependencies, options, or mixins.
license: MIT
---

# lets Task Runner

Use this skill when working in a project that uses `lets`, usually indicated by a `lets.yaml` file.

## Workflow

1. Inspect `lets.yaml` before changing or running tasks. Also check all files declared in `mixins` - those are lets config files that extend/override main `lets.yaml` file. Mixins with `-` at the beginning of the file name are optional and file may not exist. Used for gitignored mixins
2. Use `lets --help` to list commands and global flags.
3. Use `lets help <command>` to inspect a command's description, options, and usage before running it.
4. Prefer project-defined `lets` commands for build, test, lint, format, docs, and release workflows instead of invoking underlying tools directly.

## Configuration Notes

`lets.yaml` can define top-level `shell`, `env`, `before`, `init`, `mixins`, and `commands` fields.

Command entries commonly use:

- `cmd`: shell command to run.
- `depends`: commands that must run first.
- `env`: command-specific environment.
- `options`: docopt-style CLI options exposed as `LETSOPT_*` and `LETSCLI_*` environment variables.
- `work_dir`: working directory for the command.
- `after`: commands to run after the main command.
- `checksum` and `persist_checksum`: skip work when inputs have not changed.
- `ref` and `args`: reuse another command with arguments.
- `shell`: command-specific shell settings.
