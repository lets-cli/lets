---
id: what_is_lets
title: What is lets ?
sidebar_label: What is lets ?
slug: /
---

### Introduction

`Lets` is a task runner.

You can think of it as a tool with a config where you can write tasks.

The task is usually your set of cli commands which you want to group together and gave it a name.

For example, if you want to run tests in your project you may need to run next commands:


```bash
# spinup a database for tests
docker-compose up postgres
# apply database migrations
docker-compose run --rm sql alembic upgrdade head
# run some tets
docker-compose run --rm app pytest -x "test_login"
```

This all can be represented in one task - for example `lets test`

```yaml
command:
  test:
    description: Run integration tests
    cmd: |
      docker-compose up postgres
      docker-compose run --rm sql alembic upgrdade head
      docker-compose run --rm app pytest -x "test_login"
```

And execute - `lets test`. Now everyone in you team knows how to run tests.

### Why yet another task runner ?

So is there are any of such tools out there ?

Well, sure there are some.

Many developers know such a tool called `make`.

So why not `make` ?

`make` is more like a build tool and was not intended to be used as a task runner (but usually used because of the lack of alternatives or because it is install on basicaly every developer's machine).

`make` has some sort of things which are bad/hard/no convinient for developers which use task runners on a daily basis.

Lets is a brand new task runner with a task-centric philosophy and created specifically to meet developers needs.

### Features

- `yaml config` - human-readable, recognizable and convenient format for such configs (also used by kubernetes, ansible, and many others)
- `arguments parsing` - using http://docopt.org
- `global and per/command env`
- `global and per/command dynamic env` - can be computed at runtime
- `checksum` - a feature which helps to track file changes
- `written in Go` - which means it is easy to read, write and test as well as contributing to project

To see all features, [check out config documentation](config.md)