---
id: env
title: Environment
---

### Default environment variables

`lets` has builtin environ variables which user can override before lets execution. E.g `LETS_DEBUG=1 lets test`

* `LETS_DEBUG` - enable debug messages
* `LETS_CONFIG` - changes default `lets.yaml` file path (e.g. LETS_CONFIG=lets.my.yaml)
* `LETS_CONFIG_DIR` - changes path to dir where `lets.yaml` file placed
* `LETS_CHECK_UPDATE` - disables background update checks and notifications
* `NO_COLOR` - disables colored output. See https://no-color.org/

### Environment variables available at command runtime

* `LETS_COMMAND_NAME` - string name of launched command
* `LETS_COMMAND_ARGS` - positional arguments for launched command, e.g. for `lets run --debug --config=test.ini` it will contain `--debug --config=test.ini`
* `LETS_COMMAND_WORK_DIR` - absolute path to the directory the command runs in: the root dir, or the command's `work_dir` if it sets one.
* `LETS_CONFIG` - absolute path to lets config file. For a remote config this is the URL it was loaded from.
* `LETS_CONFIG_DIR` - absolute path to the directory holding the config file. Use it to target the project rather than the directory you ran `lets` from. For a remote config this is the local cache directory.
* `LETS_OS` - current operating system name from Go runtime, for example `linux`, `darwin`, `windows`
* `LETS_ARCH` - current architecture name from Go runtime, for example `amd64`, `arm64`, `386`
* `LETS_SHELL` - shell from config or command.
* `LETSOPT_<>` - options parsed from command `options` (docopt string). E.g `lets run --env=prod --reload` will be `LETSOPT_ENV=prod` and `LETSOPT_RELOAD=true`
* `LETSCLI_<>` - options which values is a options usage. E.g `lets run --env=prod --reload` will be `LETSCLI_ENV=--env=prod` and `LETSCLI_RELOAD=--reload`

### Override command env with -E flag

### `env_file` precedence

`env_file` loads dotenv-style files from config and command definitions.

Precedence order is:

* process env
* builtin lets vars
* global `env`
* global `env_file`
* command `env`
* command `env_file`
* command options, `-E` / `--env`, checksum vars

When the same key is present in both directives, `env_file` wins over `env` at the same scope.

`env_file` file names are expanded after `env` is resolved. This means:

* global `env_file` can depend on global `env`
* command `env_file` can depend on merged global env and command `env`
* `env.sh` still does not read values from `env_file`

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
