---
id: env
title: Lets environment
---

### Default environment variables

`lets` has builtin environ variables.

* `LETS_DEBUG` - enable debug messages
* `LETS_CONFIG` - changes default `lets.yaml` file path (e.g. LETS_CONFIG=lets.my.yaml)
* `LETS_CONFIG_DIR` - changes path to dir where `lets.yaml` file placed
* `NO_COLOR` - disables colored output. See https://no-color.org/
* `LETS_COMMAND_NAME` - string name of launched command
* `LETS_COMMAND_ARGS` - positional arguments for launched command, e.g. for `lets run --debug --config=test.ini` it will contain `--debug --config=test.ini`

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
