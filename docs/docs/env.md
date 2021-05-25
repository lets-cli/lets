---
id: env
title: Lets environment
---

### Default environment variables

`lets` has builtin environ variables.

* `LETS_DEBUG` - enable debug messages
* `LETS_CONFIG` - changes default `lets.yaml` file path (e.g. LETS_CONFIG=lets.my.yaml)
* `LETS_CONFIG_DIR` - changes path to dir where `lets.yaml` file placed
* `LETS_NO_COLOR_OUTPUT` - disables colored output
* `LETS_COMMAND_NAME` - string name of launched command

### Override command env with -E flag

You can override environment for command with `-E` flag:

```yaml
shell: bash

commands:
  say:
    env:
      NAME: Rick
    cmd: echo Hello ${NAME}
```

**`lets say`** - prints `Hello Rick`

**`lets -E NAME=Morty say`** - prints `Hello Morty`

Alternatively:

**`lets --env NAME=Morty say`** - prints `Hello Morty`
