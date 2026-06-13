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

1. Inspect `lets.yaml` before changing or running tasks. Also check all files declared in `mixins`; those files extend or override the main config. Mixins prefixed with `-` are optional and may be gitignored.
2. Use `lets --help` to list commands and global flags.
3. Use `lets help <command>` to inspect a command's description, options, and usage before running it.
4. Prefer project-defined `lets` commands for build, test, lint, format, docs, and release workflows instead of invoking underlying tools directly.
5. After editing config, run the affected `lets` command or the repository's standard test/lint task.

## Config Basics

The main config is usually `lets.yaml`. Create one with `lets --init` if the project does not have it.

Minimal config:

```yaml
shell: bash

commands:
  test:
    description: Run tests
    cmd: go test ./...
```

Top-level fields:

- `version`: minimum required lets version.
- `shell`: default shell for commands; commonly `bash`.
- `env`: global environment available to every command.
- `env_file`: dotenv files loaded for every command.
- `before`: script prepended to every command and dependency invocation.
- `init`: script run once per lets invocation before the first command.
- `mixins`: additional lets config files or remote mixins.
- `commands`: map of task names to command definitions.

## Writing Commands

Use short syntax for simple commands:

```yaml
commands:
  fmt: go fmt ./...
```

Use long syntax when you need descriptions, dependencies, env, options, checksums, or cleanup:

```yaml
commands:
  test:
    description: Run unit tests
    depends: [generate]
    env:
      GOFLAGS: -count=1
    cmd: go test ./...
```

Command fields:

- `cmd`: shell command to run.
- `description`: shown in command help.
- `work_dir`: run from a path relative to the config directory.
- `shell`: override shell for one command.
- `after`: cleanup script that runs even if `cmd` fails.
- `depends`: commands that must run first.
- `env`: command-specific environment.
- `options`: docopt-style CLI options exposed as `LETSOPT_*` and `LETSCLI_*` environment variables.
- `env_file`: command-specific dotenv files.
- `checksum` and `persist_checksum`: calculate file checksums and detect whether inputs changed.
- `ref` and `args`: reuse another command with arguments.
- `group`: organize commands in help output.

`cmd` can be a string, multiline string, array, or experimental map:

```yaml
commands:
  build:
    cmd: |
      echo "Building"
      go build ./cmd/lets

  test:
    cmd:
      - go
      - test
      - ./...

  dev:
    cmd:
      api: go run ./cmd/api
      web: npm run dev
```

With array `cmd`, extra CLI args are appended. For example `lets test -run TestName` runs `go test ./... -run TestName`.

For map `cmd`, users can select entries with global flags before the command name, such as `lets --only api dev` or `lets --exclude web dev`.

## Options And Arguments

Use `options` for user-facing command arguments. It is a docopt usage block.

```yaml
commands:
  release:
    description: Create a release
    options: |
      Usage: lets release <version> [--dry-run] [--message=<message>]

      Options:
        <version>             Version to release
        --dry-run             Print actions without changing anything
        --message=<message>   Release message
    cmd: |
      if [[ -n "${LETSOPT_DRY_RUN}" ]]; then
        echo "Dry run for ${LETSOPT_VERSION}"
      fi
      echo "Message: ${LETSOPT_MESSAGE}"
```

Rules:

- Positional args become `LETSOPT_<NAME>`, for example `<version>` becomes `LETSOPT_VERSION`.
- Long flags become uppercase with `-` converted to `_`, for example `--dry-run` becomes `LETSOPT_DRY_RUN`.
- `LETSOPT_*` contains parsed values such as `true`, `staging`, or positional args.
- `LETSCLI_*` contains the raw CLI fragment, useful when forwarding flags to another command.
- Use kebab-case for CLI flags, not snake_case.

## Environment

Global `env` applies to all commands. Command `env` extends or overrides it.

```yaml
env:
  TARGET: dev
  IMAGE:
    sh: echo "app-${TARGET}"

commands:
  build:
    env:
      TAG:
        sh: git rev-parse --short HEAD
    cmd: docker build -t "${IMAGE}:${TAG}" .
```

Environment entries are evaluated in declaration order, so later entries can reference earlier ones. Command env can also reference global env.

Use `env_file` to load dotenv files:

```yaml
env:
  TARGET: dev

env_file:
  - .env
  - -.env.local

commands:
  up:
    env_file:
      - .env.${TARGET}
      - name: .env.required
        required: true
    cmd: docker compose up
```

Rules:

- `-filename` means optional, equivalent to `required: false`.
- `env_file` paths are resolved relative to the config directory, not `work_dir`.
- Values from env files override values from `env`.
- File names are expanded after available env values are resolved.

## Dependencies And Cleanup

Use `depends` for prerequisite commands. Dependencies run before the command.

```yaml
commands:
  build:
    cmd: docker build -t app .

  test:
    depends: [build]
    cmd: go test ./...
```

Dependencies can pass args or env to the dependency command:

```yaml
commands:
  test:
    depends:
      - name: build
        args: [--verbose]
        env:
          TARGET: test
    cmd: go test ./...
```

Use `after` for cleanup that must happen even when the command fails:

```yaml
commands:
  redis:
    cmd: docker compose up redis
    after: docker compose stop redis
```

Use `init` for setup that should run once per lets invocation. Avoid heavy `before` scripts because `before` runs before each command and dependency.

## Checksums

Use `checksum` when command behavior depends on file contents. lets calculates SHA1 checksums and exposes them as env vars.

```yaml
commands:
  deps:
    checksum:
      deps:
        - package.json
        - package-lock.json
    persist_checksum: true
    cmd: |
      if [[ "${LETS_CHECKSUM_DEPS_CHANGED}" == "true" ]]; then
        npm install
      fi
```

Rules:

- A list checksum exposes `LETS_CHECKSUM`.
- A named checksum map exposes `LETS_CHECKSUM_<NAME>` plus combined `LETS_CHECKSUM`.
- `persist_checksum: true` exposes `LETS_CHECKSUM_CHANGED` and named `LETS_CHECKSUM_<NAME>_CHANGED`.
- Persisted checksums update only after a successful command exit.
- Glob patterns are supported in checksum file lists.

## Mixins

Use `mixins` to split large configs or share common commands.

```yaml
shell: bash

mixins:
  - lets.build.yaml
  - -lets.my.yaml

commands:
  test:
    cmd: go test ./...
```

Rules:

- Local mixins are paths to other lets config files.
- A mixin name prefixed with `-` is optional and may not exist; this is useful for gitignored personal config.
- Remote mixins use `url` and optional `version`; lets caches them under `.lets/mixins`.
- When editing a command, search mixins too because commands may be defined, extended, or overridden there.

## Reusing Commands

Use experimental `ref` and `args` to create aliases with predefined arguments.

```yaml
commands:
  hello:
    cmd: echo Hello $@

  hello-world:
    ref: hello
    args: World
```

`ref` is only compatible with `args`; do not combine it with normal command directives.

Use `group` to organize help output:

```yaml
commands:
  build:
    group: Build
    cmd: npm run build
```

## Common Pitfalls

- Do not commit `.lets/`, generated binaries, coverage files, dependency directories, or personal mixins such as `lets.my.yaml` unless the repository explicitly tracks them.
- Prefer one clear command over several near-duplicate commands. Use `options`, `args`, or `ref` when that keeps the config simpler.
- Keep `description` first-line concise; help output uses the first line.
- Use multiline `cmd: |` for scripts with conditionals or multiple shell statements.
- Quote env values that look like booleans or numbers when they must stay strings.
- Remember `work_dir` changes where `cmd` runs, but `env_file` remains relative to the config directory.
- Prefer project-defined `lets` tasks for validation, for example `lets test`, `lets lint`, or `lets fmt` when present.
