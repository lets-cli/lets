---
id: changelog
title: Changelog
---

## [Unreleased]

## [0.0.40](https://github.com/lets-cli/lets/releases/tag/v0.0.40)

* [Added] Allow override command arguments and env when using command in `depends`

   See example [in config docs](/docs/config#override-arguments-in-depends-command)

* [Added] Validate if commands declared in `depends` actually exist.
* [Refactoring] Refactored `runner` package, implemented `Runner` struct.
* [Added] Support `NO_COLOR` env variable to disable colored output. See https://no-color.org/
* [Added] `LETS_COMMAND_ARGS` - will contain command's positional args. [See config](/docs/env#default-environment-variables).
  
  Also, special bash env variables such as `"$@"` and `"$1"` etc. now available inside `cmd` script and work as expected. 
* [Added] `work_dir` directive for command. See [config](/docs/config#work_dir)
* [Added] `shell` directive for command. See [config](/docs/config#shell-1)
* [Added] `--init` flag. Run `lets --init` to create new `lets.yaml` with example command
* [Refactoring] updated `bats` test framework and adjusted all bats tests 
* [Added] `ref` directive to `command`. Allows to declare existing command with predefined args [See config](/docs/config#ref).
* [Added] `sh` and `checksum` execution modes for global level `env` and command level `env` [See config](/docs/config#env).
  `eval_env` is deprecated now, since `env` with `sh` execution mode does exactly the same


## [0.0.33](https://github.com/lets-cli/lets/releases/tag/v0.0.33)

* [Added] Allow templating in command `options` directive [docs](/docs/advanced_usage#command-templates)


## [0.0.32](https://github.com/lets-cli/lets/releases/tag/v0.0.32)

* [Fixed] Publish lets to homebrew


## [0.0.30](https://github.com/lets-cli/lets/releases/tag/v0.0.30)

* [Added] Build `lets` for `arm64 (M1)` arch
* [Deleted] Drop `386` arch builds
* [Added] Publish `lets` to homebrew
* [Added] `--upgrade` flag to make self-upgrades


## 0.0.29

* [Added] `after` directive to command.
  It allows to run some script after main `cmd`
  ```yaml
  commands:
    run:
      cmd: docker-compose up redis
      after: docker-compose stop redis
  ```

* [Added] `before` global directive to config.
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

* [Added] ignored minixs
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

* [Fix] Added environment variable value coercion.

  ```yaml
  commands:
    run:
      env:
        VERBOSE: 1
      cmd: docker-compose up
  ```

  Before 0.0.28 release this config vas invalid because `1` was not coerced to string `"1"`. Now it works as expected.

## 0.0.27

* Added `-E` (`--env`) command-line flag. It allows to set(override) environment variables for a running command.
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
