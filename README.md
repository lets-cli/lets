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
```

3. Run commands

```bash
lets echo
# will print Hello
```

## lets.yaml

```yaml
commands:
    [name]:
      description: string
      options: docopt string
      checksum: array of files to calculate checksum (accessed via LETS_CHECKSUM env)
      env: array of evironment variables
      eval_env: array of evironment variables but each will we evaluated (run in shell)
      cmd: string or array of strings
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

### TODO

Yet there is no binaries

- [x] yaml config
- [x] cobra
- [x] docopts as just strings
- [x] docopts as repeated flags, joined in string
- [x] pass opts as is if cmd is an array
- [x] file hashes (checksums)
- [ ] global checksums
- [ ] multiple checksums in one command (kv)
- [x] depends on other commands
- [ ] inherit configs
- [x] LETS_DEBUG env for debugging logs
- [ ] command to only calculate checksum
- [ ] capture env from shell
- [ ] env computing
  - [ ] global env
  - [x] command env
- [ ] dogfood on ci
- [ ] add version flag to lets
- [ ] add verbose flag to lets
