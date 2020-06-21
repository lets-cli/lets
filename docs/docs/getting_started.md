---
id: getting_started
title: Getting started with Lets
sidebar_label: Getting started with Lets
---

If you already have `lets.yaml` then just go to that directory and run `lets` to see all available commands.

If you are starting from scratch and want to create a new `lets.yaml`, please, take a look at [Creating new config](#creating-new-config) section.

### Config lookup

`lets` will be looking for a config starting from where you call `lets` and up to the `/`.

For example:

```bash
cd /home/me
touch lets.yaml

mkdir ./project
cd ./project

lets # it will use lets.yaml at /home/me/lets.yaml

touch lets.yaml

lets # it will use lets.yaml right here (at /home/me/project/lets.yaml)
```

## Creating new config

1. Create `lets.yaml` file in your project directory
2. Add commands to `lets.yaml` config. [Config reference](config.md)

```yaml
shell: /bin/sh

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
lets echo # will print 
# Hello
```

```bash
lets run --debug --level=info # will print
# LETSOPT_DEBUG=true
# LETSOPT_LEVEL=info
# LETSCLI_DEBUG=--debug
# LETSCLI_LEVEL=--level info

```

## Lets directory

At first run `lets` will create `.lets` directory in the current directory.

`lets` uses `.lets` to store some specific data such as checksums, etc.

It's better to add `.lets` to your `.gitignore` end exclude it in your favorite ide.
