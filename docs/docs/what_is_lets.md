---
id: what_is_lets
title: What is lets ?
sidebar_label: What is lets ?
---

Lets is a task runner.

You can think of it as a tool with a config where you can write tasks.

The task is usually your set of commands which you can type ten times a day, for example, you want to run tests in your project:

- pull the latest master
- spinup a database
- run migrations
- run tests (maybe run only one test file)

Or some initial setup script for your application:

- docker build -t myapp -f Dockerfile.dev .
- docker-compose up myapp postgres

This all can be represented in the task.

So is there are any of such tools out there ? 

Well, sure there are some.

Many developers know such a tool called Make.

So why not Make ?

Make is more like a build tool and was not intended to use as a task runner (but usually used because of the lack of alternatives).

Make has some sort of things which are bad/hard/no convinient for developers which use task runners on a daily basis.

Lets is a brand new task runner with a task-centric philosophy and written specifically to meet developers needs.

Lets features:

- yaml-based config - human-readable, recognizable and convenient format for such configs (also used by kubernetes, ansible, and many others)

- has support for global env
- has support for global computed env (known as `eval_env`)
- has support for per-command env 
- has support for per-command computed env (known as `eval_env`)
- has `checksum` support - a feature which helps to track file changes
- has checksum persistence
- written in Go - which means it is easy to read, write and test as well as contributing to project


