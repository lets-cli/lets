---
id: config
title: Config reference
---


Config schema

* [shell](#shell)
* [mixins](#mixins)
* [env](#global-env)
* [eval_env](#global-eval_env)
* [before](#before)
* [commands](#commands)
    * [description](#description)
    * [cmd](#cmd)
    * [work_dir](#work_dir)
    * [after](#after)
    * [depends](#depends)
    * [options](#options)
    * [env](#env)
    * [eval_env](#eval_env)
    * [checksum](#checksum)
    * [persist_checksum](#persist_checksum)
    * [ref](#ref)
    * [args](#args)


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

`type: map string => string or map with execution mode`

Specify global env for all commands.

Env can be declared as static value or with execution mode:

Example:

```yaml
shell: bash
env:
  MY_GLOBAL_ENV: "123"
  MY_GLOBAL_ENV_2: 
    sh: echo "`id`"
  MY_GLOBAL_ENV_3:
    checksum: [Readme.md, package.json]
```

### Global eval_env 

**`Deprecated`**

`key: eval_env`

`type: mapping string => string`

> Since `env` now has `sh` execution mode, `eval_env` is deprecated.

Specify global eval_env for all commands.

Example:

```yaml
shell: bash
eval_env:
  CURRENT_UID: echo "`id -u`:`id -g`"
```

### Global before 

`key: before`

`type: string`

Specify global before script for all commands.

Example:

Run `redis` with docker-compose using log level ERROR

```yaml
shell: bash

before:
  function @docker-compose() {
    docker-compose --log-level ERROR $@
  }

  export XXX=123

commands:
  redis: |
    echo $XXX
    @docker-compose up redis
```

### Mixins

`key: mixins`

`type: list of string`

Allows to split `lets.yaml` into mixins (mixin config files).

To make `lets.yaml` small and readable it is convenient to split main config into many smaller ones and include them

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

### Ignored mixins

It is possible to specify mixin file which could not exist. It is convenient when you have
git-ignored file where you write your own commands.

To make `lets` read this mixin just add `-` prefix to filename

For example:

```yaml
shell: bash
mixins:
  - -my.yaml
```

Now if `my.yaml` exists - it will be loaded as a mixin. If it is not exist - `lets` will skip it.

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
  - map of string => string (experimental)
```

Actual command to run in shell.

Can be either:
 - a string (also a multiline string)
 - an array of strings - it will allow to append all arguments passed to command as is (see bellow)
 - a map of string => string - this will allow run commands in parallel `(experimental)`

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

`cmd` can be a map `(it is experimental feature)`.

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

### `work_dir`

`key: work_dir`

`type: string`

Specify work directory to run in. Path must be relative to project root. Be default command's workdir is project root (where lets.yaml located).

Example:

```yaml
commands:
  run-docs:
    description: Run docusaurus documentation live
    work_dir: docs
    cmd: npm start
```

### `shell`

`key: shell`

`type: string`

Specify shell to run command in.

Any shell can be used, not only sh-compatible, for example `python`.

Example:

```yaml
shell: bash
commands:
  run-sh:
    shell: /bin/sh
    cmd: echo Hi
    
  run-py:
    shell: python
    cmd: print('hi')
```


### `after`

`key: after`

`type: string`

Specify script to run after the actual command. May be useful, when we want to cleanup some resources or stop some services

`after` script is guaranteed to execute if specified, event if `cmd` exit code is not `0`

Example:

```yaml
commands:
  redis:
    description: Run redis
    cmd: docker-compose up redis
    after: docker-compose stop redis

  run:
    description: Run app and services
    cmd: 
      app: node server.js
      redis: docker-compose up redis
    after: |
      echo Stopping app and redis
      docker-compose stop redis
```

### `depends`

`key: depends`

`type: array of string or array or object`

Specify what commands to run before the actual command. May be useful, when you have one shared command.
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

#### Override arguments in depends command

It is possible to override arguments or env for any commands declared in depends.

For example we want:

- `build` command to be executed with `--verbose` flag in test `depends`.
- `alarm` command to be executed with `Something is happening` arg and `LEVEL: INFO` env in test `depends`.

```yaml
commands:
  greet:
    cmd: echo Hi developer

  alarm:
    options: |
      Usage: lets alarm <msg>
    env:
      LEVEL: DEBUG
    cmd: echo Alarm ${LETSOPT_MSG}

  build:
    description: Build docker image
    options: |
      lets build [--verbose]
    cmd: |
      if [[ -n ${LETSOPT_VERBOSE} ]]; then
        echo Building docker image
      fi
      docker build -t myimg . -f Dockerfile

  test:
    description: Test something
    depends:
      - greet
      - name: alarm
        args: Something is happening
        env:
          LEVEL: INFO
      - name: build:
        args: [--verbose]
    cmd: go test ./... -v
```

Running `lets test` will output: 

```shell
# lets test
# Hi developer
# Something is happening
# Building docker image
# ... continue building docker image
```

### `options`

`key: options`

`type: string (multiline string)`

One of the most cool things about `lets` than it has built in docopt parsing.
All you need is to write a valid docopt for a command and lets will parse and inject all values for you.

More info [http://docopt.org](http://docopt.org)

When parsed, `lets` will provide two kind of env variables:

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

Raw flags (useful if for example you want to pass `--debug` as is)

```bash
echo LETSCLI_ARGS=${LETSCLI_ARGS} # LETSCLI_ARGS=one two three
echo LETSCLI_LOG_LEVEL=${LETSCLI_LOG_LEVEL} # LETSCLI_LOG_LEVEL=--log-level info
echo LETSCLI_DEBUG=${LETSCLI_DEBUG} # LETSCLI_DEBUG=--debug
```


### `env`

`key: env`

`type: mapping string => string or map with execution mode`

Env is as simple as it sounds. Define additional env for a command: 

Env can be declared as static value or with execution mode:

Example:

```yaml
commands:
  test:
    description: Test something
    env:
      GO111MODULE: "on"
      GOPROXY: https://goproxy.io
      MY_ENV_1:
        sh: echo "`id`"
      MY_ENV_2:
        checksum: [Readme.md, package.json]
    cmd: go build -o lets *.go
```


### `eval_env`

**`Deprecated`**

`key: eval_env`

`type: mapping string => string`

> Since `env` now has `sh` execution mode, `eval_env` is deprecated.

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

Checksum used for computing file hashes. It is useful when you depend on some files content changes.

In `checksum` you can specify:

- a list of file names 
- a list of file regexp patterns (parsed via go `path/filepath.Glob`)

or

- a mapping where key is name of env variable and value is:
    - a list of file names 
    - a list of file regexp patterns (parsed via go `path/filepath.Glob`)

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

`LETS_CHECKSUM_CHANGED` will be true after the very first execution, because when you first run command, there is no checksum yet, so we are calculating new checksum - that means that checksum has changed.

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

### `ref`

`key: ref`

`type: string`

**`Experimental feature`**

NOTE: `ref` is not compatible (yet) with any directives except `args`. Actually `ref` is a special syntax that turns command into reference to command. It may be changed in the future.

Allows to run command with predefined arguments. Before this was implemented, if you had some commmand and wanted same command but with some predefined args, you had to use so called `lets-in-lets` hack.

Before:

```yaml
commands:
  ls:
    cmd: [ls]

  ls-table:
    cmd: lets ls -l
```


Now:

```yaml
commands:
  ls:
    cmd: [ls]

  ls-table:
    ref: ls
    args: -l
```

### `args`

`key: args`

`type: string`

**`Experimental feature`**

`args` is used only with [ref](#ref) and allows to set additional positional args to referenced command. See [ref](#ref) example.

