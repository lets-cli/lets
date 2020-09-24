---
id: changelog
title: Changelog
---

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
