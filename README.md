# lets
Simple CLI task runner

> Very alpha software. Things may and will change/break


## Usage

1. Create `lets.yaml` file in your project directory
2. Add commands to `lets.yaml` config. [Config reference](#lets.yaml)

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
# LETSOPT_RUN=true
# LETSOPT_LEVEL=info

```

## lets.yaml

```yaml
shell: bash
commands:
    [name]:
      description: string
      options: docopt string
      checksum: array of files to calculate checksum (accessed via LETS_CHECKSUM env)
      env: array of evironment variables
      eval_env: array of evironment variables but each will we evaluated (run in shell)
      cmd: string or array of strings
```

## Install

**Shell script**:

This will install `lets` binary to `/usr/local/bin` directory. But you can change install location to any directory you want

```bash
sudo curl -sfL https://raw.githubusercontent.com/kindritskyiMax/lets/master/install.sh | sudo sh -s -- -b /usr/local/bin
```

**Binary (Cross-platform)**:

Download the version you need for your platform from [Lets Releases](https://github.com/kindritskyiMax/lets/releases). 

Once downloaded, the binary can be run from anywhere.

Ideally, you should install it somewhere in your PATH for easy use. `/usr/local/bin` is the most probable location.

## Build

To build a binary:

```bash
go build -o lets *.go
```

To install in system

```bash
go build -o lets *.go && sudo mv ./lets /usr/local/bin/lets
```

### TODO

Yet there is no binaries

- [x] yaml config
- [x] cobra
- [x] docopts as just strings
- [x] docopts as repeated flags, joined in string
- [x] pass opts as is if cmd is an array
- [x] file hashes (checksums)
- [ ] global checksums (check if some commands use checksum so we can skip its calculation)
- [ ] multiple checksums in one command (kv)
- [x] depends on other commands
- [ ] inherit configs
- [x] LETS_DEBUG env for debugging logs
- [ ] command to only calculate checksum
- [x] capture env from shell
- [ ] env computing
  - [ ] global env
  - [x] command env
- [ ] dogfood on ci
- [x] add version flag to lets
- [ ] add verbose flag to lets
