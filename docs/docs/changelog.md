---
id: changelog
title: Changelog
---

## [Unreleased](https://github.com/lets-cli/lets/releases/tag/v0.0.X)

## [0.0.57](https://github.com/lets-cli/lets/releases/tag/v0.0.57)

* `[Dependency]` update go to `1.24`
* `[Added]` support custom top-level keywords that start with `x-`
* `[Added]` check for invalid top-level keywords during config parsing
* `[Added]` support YAML aliases in `env` - env will be merged aliases mapping

## [0.0.56](https://github.com/lets-cli/lets/releases/tag/v0.0.56)

This tag is not released due to build issues

## [0.0.55](https://github.com/lets-cli/lets/releases/tag/v0.0.55)

* `[Added]` `lets self` command that is ment to be a new home for all lets own commands such as `completions` (soon) or `lsp`
* `[Added]` `lets self lsp` command that starts built-in `lsp` server with go to definition support and completions
* [`Development`] Since `lsp` implementation uses  `https://tree-sitter.github.io` (C library with go bindings) as a internal parser `lets` now build with `CGO_ENABLED=1`. If you are building on local machine, you do not have to specify `CGO_ENABLED` env variable. But you may have to install some build system dependencies in case compilatino fails.
* `[CI]` reworked release pipeline now supports go cross compilation
* `[Improvment]` split commands in help message into `Commands` and `Internal commands`
* `[Dependency]` update go to `1.23`
* `[Dependency]` update goreleaser to `1.63.x`
* `[Dependency]` update golangci-lint to `2.x` (also applied some lint fixes across codebase)

## [0.0.54](https://github.com/lets-cli/lets/releases/tag/v0.0.54)

* `[Fixed]` `lets --init` now properly creates `lets.yaml`. Issue [#263](https://github.com/lets-cli/lets/issues/263)
* `[Dependency]` update go to `1.22`
* `[Fixed]` ensure `init` script does not get called twice. Issue [#256](https://github.com/lets-cli/lets/issues/256)
* `[Fixed]` do not fail if `sh` in env is empty. Issue [#235](https://github.com/lets-cli/lets/issues/235)
* `[Fixed]` support `arm64` in `lets --upgrade`. Issue [#254](https://github.com/lets-cli/lets/issues/254)

## [0.0.53](https://github.com/lets-cli/lets/releases/tag/v0.0.53)

* `[Fixed]` change `SHELL` env to `LETS_SHELL` because setting system variable `SHELL` to just `bash` without full path to binary cases and error in some cases.

## [0.0.52](https://github.com/lets-cli/lets/releases/tag/v0.0.52)

* `[Dependency]` update and pin goreleaser

## [0.0.51](https://github.com/lets-cli/lets/releases/tag/v0.0.51)

* `[Fixed]` `-E` flag now overrides env not only when env variable already declared in command's `env`, but also it overrides env variables that declared in global `env`. Also it sets env variable even if no env variable was declared either in global or command's `env`.

## [0.0.50](https://github.com/lets-cli/lets/releases/tag/v0.0.50)

* `[Dependency]` upgrade cobra to 1.6.0
* `[Dependency]` upgrade logrus to 1.9.0
* `[Fixed]` Removed builtin `--help` flag for subcommands. Now using `--help` will pass this flag to underlying `cmd` script.
* `[Added]` Add `--debug` (`-d`) debug flag. It works same as `LETS_DEBUG=1` env variable. It can be specified as `-dd` (or `LETS_DEBUG=2`). Lets then prints more verbose logs. Also debug logs are nicer now.
* `[Added]` Add `--config` `-c` flag. It works same as `LETS_CONFIG=<path to lets file>` env variable.
* `[Added]` `LETS_CONFIG` env variable now present at command runtime, and contains lets config filename. Default is `lets.yaml`.
* `[Added]` `LETS_CONFIG_DIR` env variable now present at command runtime, and contains absolute path to dir where lets config found.
* `[Added]` `LETS_COMMAND_WORKDIR` env variable now present at command runtime, and contains absolute path to dir where `command.work_dir` points.
* `[Added]` Add `init` directive to config. It is a script that will be executed only once before any other commands. It differs from `before` in a way that `before` is a script that is prepended to each command's script and thus will be execured every time a command executes. See [config reference for init](/docs/config#init)
* `[Refactoring]` Config parsing is reimplemented using `UnmarhallYAML`. This ends up in reduced size and complexity of parsing code.
* `[Refactoring]` `Command` now is clonable and this opened a possibility to reimplement `ref`, `depends` as map and `--no-depends` - now we clone a command and modify a brand new struct instead of mutating the same command (which was not safe).
* `[Refactoring]` `Command.Cmd` script was replaced with `Cmds` struct which represents a list of `Cmd`. This allowed generalizing so-called cmd-as-map into a list of commands that will be executed in parallel (see `Executor.executeParallel`).
* `[Refactoring]` Error reporting has changed in some places and if one is depending on particular error messages it will probably reak.
* `[Refactoring]` Simplified `Executor` by extracting commands filtering by `--only` and `--exclude` flags into `subcommand.go`.
* `[Added]` Command short syntax. See [config reference for short syntax](/docs/config#short-syntax). Example:

  Before:
  ```yaml
  commands:
    hello:
      cmd: echo Hello
  ```
  After:
  ```yaml
  commands:
    hello: echo Hello
  ```
* `[Added]` If command starts with `_` it is considered *hidden* and will not be printed in help message unless `--all` root flag is provided.
* `[Changed]` Improved rendering of `description` directive. When used in `lets help` usage message the text will be stripped to 120 chars and only text up to first `\n` will be taken. Thus help message will be rendered cosnsistently in width. It is still possible to have a lot of text in `description` and all this text will be used in `lets help [command]` help message. See [config reference for how help messages works](/docs/config#description)

## [0.0.49](https://github.com/lets-cli/lets/releases/tag/v0.0.49)

* `[Added]` remote mixins `experimental` support. See [config](/docs/config#remote-mixins-experimental) for more details.

## [0.0.48](https://github.com/lets-cli/lets/releases/tag/v0.0.48)

* `[Added]` `--no-depends` global option. Lets will skip `depends` for running command

  ```shell
  lets --no-depends run
  ```
## [0.0.47](https://github.com/lets-cli/lets/releases/tag/v0.0.47)

* `[Added]` completion for command options
* `[Dependency]` use fork of docopt.go with extended options parser
## [0.0.45](https://github.com/lets-cli/lets/releases/tag/v0.0.45)

* `[Fixed]` **`Breaking change`** Fix duplicate files for checksum.
  This will change checksum output if the same file has been read multiple times.
* `[Fixed]` Fix parsing for ref args when declared as string.
* `[Added]` ref `args` can be a list of string

## [0.0.44](https://github.com/lets-cli/lets/releases/tag/v0.0.44)

* `[Fixed]` Run ref declared in `depends` directive.

## [0.0.43](https://github.com/lets-cli/lets/releases/tag/v0.0.43)

* `[Noop]` Same as 0.0.42, deployed by accident.

## [0.0.42](https://github.com/lets-cli/lets/releases/tag/v0.0.42)

* `[Fixed]` Fixed publish to `aur` repository.

## [0.0.41](https://github.com/lets-cli/lets/releases/tag/v0.0.41)

* `[Fixed]` Tried to fixe publish to `aur` repository.

## [0.0.40](https://github.com/lets-cli/lets/releases/tag/v0.0.40)

* `[Added]` Allow override command arguments and env when using command in `depends`

   See example [in config docs](/docs/config#override-arguments-in-depends-command)

* `[Added]` Validate if commands declared in `depends` actually exist.
* `[Refactoring]` Refactored `executor` package, implemented `Executor` struct.
* `[Added]` Support `NO_COLOR` env variable to disable colored output. See https://no-color.org/
* `[Added]` `LETS_COMMAND_ARGS` - will contain command's positional args. [See config](/docs/env#default-environment-variables).

  Also, special bash env variables such as `"$@"` and `"$1"` etc. now available inside `cmd` script and work as expected.
* `[Added]` `work_dir` directive for command. See [config](/docs/config#work_dir)
* `[Added]` `shell` directive for command. See [config](/docs/config#shell-1)
* `[Added]` `--init` flag. Run `lets --init` to create new `lets.yaml` with example command
* `[Refactoring]` updated `bats` test framework and adjusted all bats tests
* `[Added]` `ref` directive to `command`. Allows to declare existing command with predefined args [See config](/docs/config#ref).
* `[Added]` `sh` and `checksum` execution modes for global level `env` and command level `env` [See config](/docs/config#env).
  `eval_env` is deprecated now, since `env` with `sh` execution mode does exactly the same


## [0.0.33](https://github.com/lets-cli/lets/releases/tag/v0.0.33)

* `[Added]` Allow templating in command `options` directive [docs](/docs/advanced_usage#command-templates)


## [0.0.32](https://github.com/lets-cli/lets/releases/tag/v0.0.32)

* `[Fixed]` Publish lets to homebrew


## [0.0.30](https://github.com/lets-cli/lets/releases/tag/v0.0.30)

* `[Added]` Build `lets` for `arm64 (M1)` arch
* `[Deleted]` Drop `386` arch builds
* `[Added]` Publish `lets` to homebrew
* `[Added]` `--upgrade` flag to make self-upgrades


## 0.0.29

* `[Added]` `after` directive to command.
  It allows to run some script after main `cmd`
  ```yaml
  commands:
    run:
      cmd: docker-compose up redis
      after: docker-compose stop redis
  ```

* `[Added]` `before` global directive to config.
  It allows to run some script before each main `cmd`
  ```yaml
  before: |
    function @docker-compose() {
      docker-compose --log-level ERROR $@
    }

  commands:
    run:
      cmd: @docker-compose up redis
  ```

* `[Added]` ignored minixs
  It allows to include mixin only if it exists - otherwise lets will ignore it.
  Useful for git-ignored files.

  Just add `-` prefix to mixin filename

  ```yaml
  mixins:
    - -my.yaml

  commands:
    run:
      cmd: docker-compose up redis
  ```


## 0.0.28

* `[Fixed]` Added environment variable value coercion.

  ```yaml
  commands:
    run:
      env:
        VERBOSE: 1
      cmd: docker-compose up
  ```

  Before 0.0.28 release this config vas invalid because `1` was not coerced to string `"1"`. Now it works as expected.

## 0.0.27

* `[Added]` `-E` (`--env`) command-line flag. It allows to set(override) environment variables for a running command.
  Example:

  ```bash
  # lets.yaml
  ...
  commands:
    greet:
      env:
        NAME: Morty
      cmd: echo "Hello ${NAME}"
  ...

  lets -E NAME=Rick greet
  ```

* Changed behavior of `persist_checksum` at first run. Now, if there was no checksum and we just calculated a new checksum, that means checksum has changed, hence `LETS_CHECKSUM_CHANGED` will be `true`.
