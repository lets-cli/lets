# lets
Simple CLI task runner

Lets takes the best from Makefile and bash and presents you a simple yet powerful tool for defining and running cli tasks and commands.

Just describe your commands in `lets.yaml` and lets will do the rest.

> Very alpha software. Things may and will change/break

* [Install](#install)
* [Usage](#usage)
* [Config](#letsyaml)
* [Build](#build)
* [Completion](#completion)

## Install

**Shell script**:

This will install `lets` binary to `/usr/local/bin` directory. But you can change install location to any directory you want

```bash
sudo curl -sfL https://raw.githubusercontent.com/lets-cli/lets/master/install.sh | sudo sh -s -- -b /usr/local/bin
```

**Binary (Cross-platform)**:

Download the version you need for your platform from [Lets Releases](https://github.com/lets-cli/lets/releases). 

Once downloaded, the binary can be run from anywhere.

Ideally, you should install it somewhere in your PATH for easy use. `/usr/local/bin` is the most probable location.

**Arch Linux**:

You can get binary release from https://aur.archlinux.org/packages/lets-bin/

If you are using `yay` as AUR helper:

```bash
yay -S lets-bin
```

Also you can get bleeding edge version from https://aur.archlinux.org/packages/lets-git/

```bash
yay -S lets-git
```

## Usage

1. Create `lets.yaml` file in your project directory
2. Add commands to `lets.yaml` config. [Config reference](#letsyaml)

```yaml
commands:
    echo:
      description: Echo text
      cmd: |
        echo "Hello"
    
    run:
      description: Run some command
      options: |
        Usage: lets run [--debug] [--level=<level>]
        
        Options:
          --debug, -d     Run with debug
          --level=<level> Log level
      cmd: |
        env
```

3. Run commands

```bash
lets echo
# will print Hello
lets run --debug --level=info
# will print
# LETSOPT_DEBUG=true
# LETSOPT_LEVEL=info#
# LETSCLI_DEBUG=--debug
# LETSCLI_LEVEL=--level info

```

## lets.yaml

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


### Top-level directives:

#### `shell` 
`key: shell`

`type: string`

Specify shell to use when running commands

Example:

```yaml
shell: bash
```

#### `global env` 
`key: env`

`type: string`

Specify global env for all commands.

Example:

```yaml
shell: bash
env:
    MY_GLOBAL_ENV: "123"
```

#### `global eval_env` 
`key: env`

`type: string`

Specify global eval_env for all commands.

Example:

```yaml
shell: bash
eval_env:
    CURRENT_UID: echo "`id -u`:`id -g`"
```

#### `mixins` 
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

#### `commands`
`key: commands`

`type: mapping`

Mapping of all available commands

Example:

```yaml
commands:
    test:
        description: Test something
```

### Command directives:

##### `description`
`key: description`

`type: string`

Short description of command - shown in help message

Example:

```yaml
commands:
    test:
        description: Test something
```

##### `cmd`
`key: cmd`

`type: string or array of strings`

Actual command to run in shell.

Can be either a string (also a multiline string) or an array of strings.

The main difference is that when specifying an array, all args, specified in terminal will be appended to cmd array 

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

##### `depends`
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


##### `options`
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


##### `env`
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


##### `eval_env`
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


##### `checksum`
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

## Build

To build a binary:

```bash
go build -o lets *.go
```

To install in system

```bash
go build -o lets *.go && sudo mv ./lets /usr/local/bin/lets
```

### Completion

You can use Bash/Zsh/Oh-My-Zsh completion in you terminal

* **Bash**

    Add `source <(lets completion -s bash)` to your `~/.bashrc` or `~/.bash_profile`

* **Zsh/Oh-My-Zsh**

    There is a [repo](https://github.com/lets-cli/lets-zsh-plugin) with zsh plugin with instructions
