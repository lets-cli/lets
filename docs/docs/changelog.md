---
id: changelog
title: Changelog
---

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
