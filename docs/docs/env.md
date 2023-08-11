---
id: env
title: Environment
---

### Default environment variables

`lets` has builtin environ variables which user can override before lets execution. E.g `LETS_DEBUG=1 lets test`

* `LETS_DEBUG` - enable debug messages
* `LETS_CONFIG` - changes default `lets.yaml` file path (e.g. LETS_CONFIG=lets.my.yaml)
* `LETS_CONFIG_DIR` - changes path to dir where `lets.yaml` file placed
* `NO_COLOR` - disables colored output. See https://no-color.org/

### Environment variables available at command runtime

* `LETS_COMMAND_NAME` - string name of launched command
* `LETS_COMMAND_ARGS` - positional arguments for launched command, e.g. for `lets run --debug --config=test.ini` it will contain `--debug --config=test.ini`
* `LETS_COMMAND_WORK_DIR` - absolute path to `work_dir` specified in command.
* `LETS_CONFIG` - absolute path to lets config file.
* `LETS_CONFIG_DIR` - absolute path to lets config file firectory.
* `LETS_SHELL` - shell from config or command.
* `LETSOPT_<>` - options parsed from command `options` (docopt string). E.g `lets run --env=prod --reload` will be `LETSOPT_ENV=prod` and `LETSOPT_RELOAD=true`
* `LETSCLI_<>` - options which values is a options usage. E.g `lets run --env=prod --reload` will be `LETSCLI_ENV=--env=prod` and `LETSCLI_RELOAD=--reload`

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
