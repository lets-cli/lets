---
id: architecture
title: Architecture
---

![Architecture diagram](/img/lets-architecture-diagram.png)

## Parser

At the start of lets application, parser tries to find `lets.yaml` file starting from current directory up to the `/`.

When config file is found, parser tries to read/parse and validate yaml config.

#### Mixins
Lets has feature called [mixins](config.md#mixins). When parser meets `mixins` directive,
it basically repeats all read/parse logic on minix files.

Since mixin config files have some limitations, although they are parsed the same way, validation is a bit different.

### How parsing works ?

`config.go:Config` struct implements `UnmarshalYAML` function, so when `yaml.Unmarshal` called with `Config` instance passed in,
custom unmarshalling code is executed.

Its common to make some normalization of commands and its data during parsing phase so the rest of the code
does not have to do any kind of normalization on its own.

### Validation

There are two validation phases.

First validation phase happens during unmarshalling and checks if:
  - directives names valid
  - directives types valid (array, map, string, number, etc.)
  - references to command in `depends` directive points to existing commands 

Second phase happens after we ensured that config is syntactically and semantically correct.

Int the second phase we are checking:
  - config version
  - circular dependencies in commands
  

## Cobra CLI Framework

We are using `Cobra` CLI framework and delegating to it most of the work related to parsing
command line arguments, help messages etc.


### Binding our config with Cobra

Now we have to bind our config to `Cobra`.

Cobra has a concept of `cobra.Command`. It is a representation of command in CLI application, for example:

```bash
git commit
git pull
```

`git` is a CLI applications and
`commit` and `pull` are commands.

In a traditional `lets` application commands will be what is declared in `lets.yaml` commands section.

To achieve this we are creating so-called `root` command and `subcommands` from config.

#### Root command

Root command is responsible for:
  - `lets` own command line flags such as `--version`, `--upgrade`, `--help` and so on.
  - `lets` commands autocompletion in terminal

#### Subcommands

Subcommand is created from our `Config.Commands` (see `initSubCommands` function).

In subcommand's `RunE` callback we are parsing/validation/normalizing command line arguments for this subcommand
and then finally executing command with `Runner`.

Since we are using `docopt` as an argument parser for subcommands, we don't let `Cobra` parse and interpret args,
and instead we are passing raw arguments as is to `Runner`.

## Runner

`Runner` is responsible for:

- parsing and preparing args using `docopt`
- calculating and storing command's checksums
- executing other commands from `depends` section
- preparing environment 
- running command in OS using `exec.Command`
