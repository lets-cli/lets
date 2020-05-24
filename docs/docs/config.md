---
id: config
title: Config reference
---


Config schema

* [shell](#shell)
* [mixins](#mixins)
* [env](#global-env)
* [eval_env](#global-eval_env)
* [commands](#commands)
    * [description](#description)
    * [cmd](#cmd)
    * [depends](#depends)
    * [options](#options)
    * [env](#env)
    * [eval_env](#eval_env)
    * [checksum](#checksum)
    * [persist_checksum](#persist_checksum)


## Top-level directives:

### Version 

`key: version`

`type: semver string`

Specify **minimum required** `lets` version to run this config.

Example:

```yaml
version: '0.0.20'
```


### Shell 

`key: shell`

`type: string`

`required: true`

Specify shell to use when running commands

Example:

```yaml
shell: bash
```

### Global env 

`key: env`

`type: map string => string`

Specify global env for all commands.

Example:

```yaml
shell: bash
env:
  MY_GLOBAL_ENV: "123"
```

### Global eval_env 

`key: env`

`type: map string => string`

Specify global eval_env for all commands.

Example:

```yaml
shell: bash
eval_env:
  CURRENT_UID: echo "`id -u`:`id -g`"
```

### Mixins

`key: mixins`

`type: list of string`

Allows to split `lets.yaml` into mixins (mixin config files).

To make `lets.yaml` small and readable its convenient to split main config into many smaller ones and include them

Example:

```yaml
# in lets.yaml
...
shell: bash
mixins:
  - test.yaml

commands:
  echo:
    cmd: echo Hi
    
# in test.yaml
...
commands:
  test:
    cmd: echo Testing...
```

And `lets test` works fine.

### Commands

`key: commands`

`type: mapping`

`required: true`

Mapping of all available commands

Example:

```yaml
commands:
  test:
    description: Test something
```

## Command directives:

### `description`

`key: description`

`type: string`

Short description of command - shown in help message

Example:

```yaml
commands:
  test:
    description: Test something
```

### `cmd`

`key: cmd`

```
type: 
  - string
  - array of strings
  - map of string => string
```

Actual command to run in shell.

Can be either:
 - a string (also a multiline string)
 - an array of strings - it will allow to append all arguments passed to command as is (see bellow)
 - a map of string => string - this will allow run commands in parallel

Example single string:

```yaml
commands:
  test:
    description: Test something
    cmd: go test ./... -v
```


Example multiline string:

```yaml
commands:
  test:
    description: Test something
    cmd: |
      echo "Running go tests..."
      go test ./... -v
```

Example array of strings:

```yaml
commands:
  test:
    description: Test something
    cmd: 
      - go
      - test
      - ./...
```


When run with cmd as array of strings:

```bash
lets test -v
```

the `-v` will be appended, so the resulting command to run will be `go test ./... -v`

Example of map of string => string

```yaml
commands:
  run:
    description: Test something
    cmd: 
      app: npm run app
      nginx: docker-compose up nginx
      redis: docker-compsoe up redis
```

There are two flags `--only` and `--exclude` you can use with cmd map.

There must be used before command name:

```bash
lets --only app run
```

### `depends`

`key: depends`

`type: array of string`

Specify what commands to run before the actual command. May be useful, when have one shared command.
For example, lets say you have command `build`, which builds docker image.

Example:

```yaml
commands:
  build:
    description: Build docker image
    cmd: docker build -t myimg . -f Dockerfile

  test:
    description: Test something
    depends: [build]
    cmd: go test ./... -v

  fmt:
    description: Format the code
    depends: [build]
    cmd: go fmt
```

`build` command will be executed each time you run `lets test` or `lets fmt`


### `options`

`key: options`

`type: string (multiline string)`

One of the most cool things about `lets` than it has built in docopt parsing.
All you need is to write a valid docopt for a command and lets will parse and inject all values for you.

More info [http://docopt.org](http://docopt.org)

When parsed, `lets` will provide two kind of env variabled:

- `LETSOPT_<VAR>`
- `LETSCLI_<VAR>`

How does it work?

Lets see an example:

```yaml
commands:
  echo-env:
    description: Echo env vars
    options:
      Usage: lets [--log-level=<level>] [--debug] <args>...
      Options:
        <args>...       List of required positional args
        --log-level,-l      Log level
        --debug,-d      Run with debug
    cmd: |
      echo ${LETSOPT_ARGS}
      app ${LETSCLI_DEBUG}
```

So here we have:

`args` - is a list of required positional args

`--log-level` - is a key-value flag, must be provided with some value

`--debug` - is a bool flag, if specified, means true, if no specified means false

In the env of `cmd` command there will be available two types of env variables:

`lets echo-env --log-level=info --debug one two three`

Parsed and formatted key values

```bash
echo LETSOPT_ARGS=${LETSOPT_ARGS} # LETSOPT_ARGS=one two three
echo LETSOPT_LOG_LEVEL=${LETSOPT_LOG_LEVEL} # LETSOPT_LOG_LEVEL=info
echo LETSOPT_DEBUG=${LETSOPT_DEBUG} # LETSOPT_DEBUG=true
```

Raw flags (useful if for example you wan to pass `--debug` as is)

```bash
echo LETSCLI_ARGS=${LETSCLI_ARGS} # LETSCLI_ARGS=one two three
echo LETSCLI_LOG_LEVEL=${LETSCLI_LOG_LEVEL} # LETSCLI_LOG_LEVEL=--log-level info
echo LETSCLI_DEBUG=${LETSCLI_DEBUG} # LETSCLI_DEBUG=--debug
```


### `env`

`key: env`

`type: mapping string => string`

Env is as simple as it sounds. Define additional env for a commmand: 

Example:

```yaml
commands:
  test:
    description: Test something
    env:
      GO111MODULE: "on"
      GOPROXY: https://goproxy.io
    cmd: go build -o lets *.go
```


### `eval_env`

`key: eval_env`

`type: mapping string => string`

Same as env but allows you to dynamically compute env:

Example:

```yaml
commands:
  test:
    description: Test something
    eval_env:
      CURRENT_UID: echo "`id -u`:`id -g`"
      CURRENT_USER_NAME: echo "`id -un`"
    cmd: go build -o lets *.go
```

Value will be executed in shell and result will be saved in env.


### `checksum`

`key: checksum`

`type: array of string | mapping string => array of string`

Checksum used for computing file hashed. It is useful when you depend on some files content changes.

In `checksum` you can specify:

- a list of file names 
- a list of file regext patterns (parsed via go `path/filepath.Glob`)

or

- a mapping where key is name of env variable and value is:
    - a list of file names 
    - a list of file regext patterns (parsed via go `path/filepath.Glob`)

Each time a command runs, `lets` will calculate the checksum of all files specified in `checksum`.

Result then can be accessed via `LETS_CHECKSUM` env variable.

If checksum is a mapping, e.g:

```yaml
commands:
  build:
    checksum:
      deps:
        - package.json
      doc:
        - Readme.md
```

Resulting env will be:

* `LETS_CHECKSUM_DEPS` - checksum of deps files
* `LETS_CHECKSUM_DOC` - checksum of doc files
* `LETS_CHECKSUM` - checksum of all checksums (deps and doc)

Checksum is calculated with `sha1`.

If you specify patterns, `lets` will try to find all matches and will calculate checksum of that files.

Example:

```yaml
shell: bash
commands:
  app-build:
    checksum: 
      - requirements-*.txt
    cmd: |
      docker pull myrepo/app:${LETS_CHECKSUM}
      docker run --rm myrepo/app${LETS_CHECKSUM} python -m app       
```


### `persist_checksum`

`key: persist_checksum`

`type: bool`

This feature is useful when you want to know that something has changed between two executions of a command.

`persist_checksum` can be used only if `checksum` declared for command.

If set to `true`, each run all calculated checksums will be stored to disk.

After each subsequent run `lets` will check if new checksum and stored checksum are different.

Result of that check will be exposed via `LETS_CHECKSUM_CHANGED` and `LETS_CHECKSUM_[checksum-name]_CHANGED` env variables. 

**IMPORTANT**: New checksum will override old checksum only if cmd has exit code **0** 

`LETS_CHECKSUM_CHANGED` will be true after the very first execution, because when you first run command, there is no checksum yet, so we calculating new checksum - that means that checksum has changed.

Example:

```yaml
commands:
  build:
    persist_checksum: true
    checksum:
      deps:
        - package.json
      doc:
        - Readme.md
```

Resulting env will be:

* `LETS_CHECKSUM_DEPS` - checksum of deps files
* `LETS_CHECKSUM_DOC` - checksum of doc files
* `LETS_CHECKSUM` - checksum of all checksums (deps and doc)

* `LETS_CHECKSUM_DEPS_CHANGED` - is checksum of deps files changed
* `LETS_CHECKSUM_DOC_CHANGED` - is checksum of doc files changed
* `LETS_CHECKSUM_CHANGED` - is checksum of all checksums (deps and doc) changed