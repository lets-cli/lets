---
id: where_commands_run
title: Where commands run
---

## Short answer

**`lets` runs your commands in the directory you ran `lets` from.** Not in the directory
`lets.yaml` lives in.

```yaml
shell: bash
commands:
  where:
    cmd: pwd
```

```console
$ cd ~/myproject && lets where
/home/you/myproject

$ cd ~/myproject/src && lets where     # lets.yaml is still up in ~/myproject
/home/you/myproject/src
```

If a command should always act on the project instead of on your current directory,
send it there yourself with `$LETS_CONFIG_DIR`:

```yaml
commands:
  build:
    cmd: cd "${LETS_CONFIG_DIR}" && go build ./...
```

## The three directories

`lets` distinguishes three directories. Most confusion comes from expecting one of them
to be another.

| | What it is | How to reach it |
| --- | --- | --- |
| **Root dir** | The directory you ran `lets` from. Commands run here by default. | `$PWD`, `$(pwd)` |
| **Config dir** | The directory holding `lets.yaml`. | `$LETS_CONFIG_DIR` |
| **Work dir** | The directory a specific command runs in: the root dir, or [`work_dir`](config.md#work_dir) if the command sets one. | `$LETS_COMMAND_WORK_DIR` |

The root dir and the config dir are the same directory whenever you run `lets` from the
project root, which is the common case. They only differ when you run `lets` from a
subdirectory, or point `-c` at a config somewhere else.

## What resolves against what

There are exactly two rules.

> **Config assembly** resolves against the config file that declares it.
> **Everything a command reads or runs** resolves against that command's work dir.

| Directive | Resolves against |
| --- | --- |
| [`mixins`](config.md#mixins) (local paths) | the config file that declares it |
| [`cmd`](config.md#cmd) working directory | the command's work dir |
| [`checksum`](config.md#checksum) file paths and globs | the command's work dir |
| [`env_file`](config.md#env_file) paths, global and command | the command's work dir |
| `env.sh` scripts | the command's work dir |
| [`work_dir`](config.md#work_dir), when relative | the root dir |

`mixins` is the only exception, and it has to be: a mixin is an include. If mixin paths
moved with your shell, a config would only be loadable from its own directory.

## Every way of pointing lets at a config

The root dir never depends on how `lets` found the config:

| you run | root dir |
| --- | --- |
| `cd proj && lets x` | `proj` |
| `cd proj && lets -c lets.yaml x` | `proj` |
| `cd proj && lets -c sub/lets.yaml x` | `proj` |
| `cd proj && lets -c /abs/path/lets.yaml x` | `proj` |
| `cd proj/deep && lets x` — config found up the tree | `proj/deep` |
| `cd proj/deep && lets -c ../lets.yaml x` | `proj/deep` |
| `cd proj && lets -c https://example.com/lets.yaml x` | `proj` |
| `cd proj && LETS_CONFIG_DIR=sub lets x` | `proj` |

`--config` / `-c`, `--config-dir` and `LETS_CONFIG_DIR` choose **which config is
loaded**. None of them changes **where commands run**.

## Remote configs

A [remote config](config.md#remote-configs) behaves exactly like a local one: commands
run in the directory you invoked `lets` from, which is the only sensible root, since the
config itself lives in a cache directory on your machine that holds nothing but the
downloaded YAML.

Two consequences:

- `LETS_CONFIG_DIR` points at that cache directory, so it is not useful for reaching your
  project. Use `$PWD` instead.
- A remote config can only mix in other URLs. A local `mixins` path is an error, because
  there is no local directory to resolve it against.

## Where `.lets` goes

The `.lets/` directory — persisted checksums and cached mixins — is created in the **root
dir**, so running `lets` from a subdirectory creates it there.

This keeps a persisted checksum next to the files it was computed from. If `.lets` sat
next to the config while checksums were computed from wherever you happened to stand, the
same stored checksum would describe different files on different runs, and change
detection would flip back and forth.

## Recipes

| You want | Write |
| --- | --- |
| Act on the directory the user is in | nothing — that is the default |
| Always act on the project, wherever `lets` is run from | `cmd: cd "${LETS_CONFIG_DIR}" && …` |
| Always act on a fixed subdirectory of the project | `cmd: cd "${LETS_CONFIG_DIR}/docs" && …` |
| Act on a subdirectory relative to where the user is | `work_dir: docs` |
| Know where the user actually is | `$PWD` |

`work_dir` takes a plain path and does not expand environment variables, so
`work_dir: ${LETS_CONFIG_DIR}/docs` will not work. Use `cd` inside `cmd` for that.

## Why it works this way

**Commands run where you are, because a config describes commands — it does not relocate
them.** `lets foo` should do what typing `foo` at your prompt would do. That makes a
command predictable from its own text: a relative path in `cmd` means what it looks like
it means, and you do not have to know where the config file happens to live to read it.

It also keeps `-c` honest. `lets -c ../other-project/lets.yaml build` borrows another
project's commands and runs them on **your** data. If `-c` also moved the root dir,
loading someone else's config would silently retarget everything, and there would be no
way to express "use these commands here".

**One work dir per command, for everything the command touches.** Before `0.0.63`,
`checksum` and `env_file` resolved against the config dir while `cmd` ran somewhere else,
so this command could hash one `data.txt` and read a different one:

```yaml
commands:
  build:
    checksum: [data.txt]
    cmd: cat data.txt
```

A single filename appearing twice in one command definition has to mean one file. Making
every command-scoped directive follow the same directory is what buys that, and it is
also what makes `work_dir` mean something coherent: it moves the whole command, not just
its script.

**Mixins are the exception because they are resolved at a different time.** Mixins are
read while the config is being assembled, before any command exists and before there is a
command work dir to speak of. A mixin path is part of the config's own structure, like an
`import`, so it belongs to the file that wrote it.

## History

This behavior was accidental for a long time before it was specified.

- Up to and including `0.0.62`, commands ran in the invocation dir — but by accident, via
  a `filepath.Abs("")` that silently defaulted to the process cwd. `checksum` and
  `env_file` resolved against the config dir at the same time, so the two disagreed.
- `0.0.63` fixed an unrelated `work_dir` bug and removed that accident, which moved the
  root dir to the config dir. `env.sh` kept running in the invocation dir, so the
  directives disagreed in a new way.
- The current release restores the invocation dir as the root and defines the resolution
  rules explicitly, so every command-scoped directive agrees.

See [ADR-0004](https://github.com/lets-cli/lets/blob/master/docs/adr/0004-root-dir-is-invocation-dir.md)
for the decision record, and the [changelog](changelog.md) for the exact breaking changes.
