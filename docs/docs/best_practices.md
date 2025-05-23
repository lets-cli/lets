---
id: best_practices
title: Best practices
---

### Naming conventions

Prefer single word over plural.

It is better to leverage semantics of `lets` as an intention to do something. For example it is natural saying `lets test` or `lets build` something.

`bad`

```
lets runs
```

`good`

```
lets run
```

---

`bad`

```
lets tests
```

`good`

```
lets test
```

### Default commands

If you have many projects (lets say - microservices) - it would be great to have one way to run and operate them when developing

- `run` command - the main purpouse of this command is to run all in once. If all projects has this command its easier to remember.
- `test` command - each projects should have a tests and a way to run them, either one file or all tests at once
- `init` command - some kind of project initialization - creates missing files, dirs for developer, checks permissions, login to docker registry, checks inotify limits for tools such as webpack and other file watchers.

### Split `lets.yaml` when it becomes big. 

If `lets.yaml` became big, it may be great to split it in a smaller, more specific files using `mixins` directive.
For example:

- **`lets.yaml`**
- **`lets.test.yaml`**
- **`lets.build.yaml`**
- **`lets.frontend.yaml`**
- **`lets.i18n.yaml`**

In each of these files we then hold all specific tasks.

### Use checskums

Checksums can help you decrease amount of task executions. How ? Lets see.

Suppose we have `js` project and we obviously holding all dependencies in `package.json`.
Also we are using Docker for reproducible development environment.

Dockerfile

```bash
FROM alpine:3.8

WORKDIR /work

COPY package.json .

RUN npm install

CMD ["npm start"]
```

What if we want to rebuild docker image every time we changed dependency ?

lets.yaml

```
shell: bash

commands:
  run:
    depends: 
      - build
    cmd: docker-compose up application

  build:
    checksum:
      - package.json
    persist_checksum: true  
    cmd: |
      if [[ ${LETS_CHECKSUM_CHANGED} == true ]]; then 
        docker-compose build application
      else
        Image is up to date
      fi
```

As you can see, we execute `build` command each time we execute `run` command (`depends`).

`persist_checksum` will save calculated checksum to `.lets` directory and all subsequent calls of `build` command will
read checksum from disk, calculate new checksum, and compare them. If `package.json` will change - we will rebuild the image.


### Initialize project using `init`

You can use `init` keyword to write a script that will do some initialization on lets startup, like creating some dirs, configs or installing project dependencies.

By default, `init` runs each time the `lets` program is executed.

You can make `init` conditional, by simply creating a file and checking if it exists at the start of `init` script.

Example:

```
shell: bash

init: |
  if [[ ! -f .lets/init_done ]]; then
    echo "calling init script"
    touch .lets/init_done
  fi
```

In this example we are checking for `.lets/init_done` file existence. If it does not exist, we will call init script and create `init_done` file as a marker of successfull init script invocation.

We are using `.lets` dir here because this dir will be created by `lets` itself and is generally a good place to create such files, but you are free to create files with any name and in any directory you want.
