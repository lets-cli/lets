---
id: advanced_usage
title: Advanced usage
---


In advanced usage we will start with a clean project and then we will add more commands to show how you can improve a developer experience in your project.

Assume you have a `node.js` project with a `run` command in `lets.yaml` from  [Basic usage](basic_usage.md)

```yaml
shell: bash

commands:
  run:
    description: Run nodejs server
    cmd: npm run server
```

### Env

You can add global or per-command `env`:

```yaml
shell: bash

env:
  DEBUG: "0"

commands:
  run:
    description: Run nodejs server
    env:
      NODE_ENV: development
    cmd: npm run server
```

### Eval env

Also if the value of the environment variable must be evaluated, you can add global or per-command `eval_env`:

```yaml
shell: bash

env:
  DEBUG: "0"

eval_env:
  CURRENT_UID: echo "`id -u`:`id -g`"
  CURRENT_USER_NAME: echo "`id -un`"

commands:
  run:
    description: Run nodejs server
    env:
      NODE_ENV: development
    cmd: npm run server
```

### Depends

You already can start your application, and like any other project your's also have dependencies. Dependencies can be added or deleted to project 

and developers have to know that there is some new dependency and it is needed to run `npm install` again.

You can do this - just add a new command and make it as a `run` command dependency, so each time you call `lets run` - dependant command will execute first.


```yaml
shell: bash

commands:
  build-deps:
    description: Install project dependencies
    cmd: npm install

  run:
    description: Run nodejs server
    depends:
      - build-deps
    cmd: npm run server
```

### Checksum

Now, each time you call `lets run` - `build-deps` will be executed first and this will guarantee that your dependencies are always up to date.

But we have one downside - run `npm install` may take some time and we do not want to wait.

`checksum` to the rescue.

Checksums allow you to know when some of the files have changed and made a decision based on that.

When you add `checksum` directive to a command - `lets` will calculate checksum from all of the files listed in `checksum` and put `LETS_CHECKSUM` env variable to command env.

`LETS_CHECKSUM` will have a checksum value.

We then can store this checksum somewhere in the file and check that stored checksum with a checksum from env.

Fortunately, `lets` have an option for that - `persist_checksum`.

If `persist_cheksum` used with `checksum` `lets` will store new checksum to `.lets` dir and each time you run a command `lets` will check if stored checksum changed from the one from env.

While using `persist_checksum`, `lets` will add new env variable to command env - `LETS_CHECKUM_CHANGED`.

You can learn more about checksum in [Checksum section](config.md#checksum)

```yaml
shell: bash

commands:
  build-deps:
    description: Install project dependencies
    checksum:
      - package.json
    persist_checksum: true
    cmd: |
      if [[ ${LETS_CHECKSUM_CHANGED} == true ]]; then
        npm install
      fi;

  run:
    description: Run nodejs server
    depends:
      - build-deps
    cmd: npm run server
    
```

So now `npm` install will be executed only on `package.json` change.

### Cmd as array

Now you have decided to add some frontend to your project. You decided to add a command to build js with a webpack.

`lets.yaml`

```yaml
shell: bash

commands:
  build-deps:
    description: Install project dependencies
    checksum:
      - package.json
    persist_checksum: true
    cmd: |
      if [[ ${LETS_CHECKSUM_CHANGED} == true ]]; then
        npm install
      fi;

  run:
    description: Run nodejs server
    depends:
      - build-deps
    cmd: npm run server

  
  js:
    description: Build project js
    cmd: npm run static
```

`package.json`

```json
{
    "scripts": {
        "static": "webpack"
    }
}
```

Now you want to run js with some options like `watch` or different config.

So lets update js command:

```yaml
js:
  description: Build project js
  cmd: 
    - npm 
    - run 
    - static
```

All we made is just rewrite `cmd` to be an array of strings. Now all positional arguments will be appended to cmd during `lets js` call.

**`lets js -- -w`** - this will pass `-w` option to webpack in `package.json`

### Options

Sooner or later you will come up with a convenient commands for your project.

`lets` options will help you with that.

Now you have a couple of environments in your project. And you want to be able to run a server with different environments.

Assume you have some configs:

- local.yaml
- stg.yaml
- prd.yaml

We can update `run` command using `options`:

```yaml
run:
  description: Run nodejs server
  depends:
    - build-deps
  options: |
    Usage: lets run [--stg] [--prd]  
  cmd: |
    CONFIG_PATH="local.yaml"
    if [[ -n ${LETSOPT_STG} ]]; then
        CONFIG_PATH="stg.yaml"
    elif [[ -n ${LETSOPT_PRD} ]]; then
        CONFIG_PATH="prd.yaml"
    fi
    npm run server -- config=$CONFIG_PATH
```

`options` is a string in a `docopt` format - http://docopt.org/.

`lets` knows how to parse docopt string and convert it in env variables.

In a few words, `lets` will capitalize on all options, replace `-` with `_` 
and append `LETSOPT_` prefix - so for `lets run --stg` we will get `LETSOPT_STG` env variable with no value as its a bool option.

Another variant of option usage:

```yaml
run:
  description: Run nodejs server
  depends:
    - build-deps
  options: |
    Usage: lets run [--config=<config>] 
  cmd: |
    npm run server -- config=${LETSOPT_CONFIG:-local.yaml}
```

In this example we also use options but unlike the previous example we using key-value options here.

So if we call `lets run --config stg.yaml` - `lets` will create `LETSOPT_CONFIG` env variable with value **stg.yaml**

One more example will show you another option `LETSCLI`.

`LETSCLI` is just a complementary env variable `lets` will create for each `LETSOPT`. 

So how does it works?

If we describe option `Usage: lets run --stg` we will actually get two env variables to one option:

- `LETSOPT_STG` with no value
- `LETSCLI_STG` with value `--stg`. It just basically stores CLI argument as is.

You can learn more about options in [Options section](config.md#options)

### Examples

There are a lot of variants how you can use `lets` in your project.

[Here](https://github.com/lets-cli/lets/tree/master/examples) you will find more examples with:

- python
- nodejs
- docker