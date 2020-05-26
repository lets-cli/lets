---
id: changelog
title: Changelog
---

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