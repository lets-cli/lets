---
id: examples
title: Examples
---

## What are these examples ?

While there is no such difference which project and which language is used with `lets`, in general we belive it would give a better understanding on how to do things with `lets` right if examples will be suited for different languages and tools. 

## Recomendations for writing `lets.yaml`

- if you have many projects (lets say - microservices) - it would be great to have one way to run and operate them when developing
    - `run` command - the main purpouse of this command is to run all in once. If all projects has this command its easier to remember.
    - `test` command - each projects should have a tests and a way to run them, either one file or all tests at once
    - `init` command - some kind of project initialization - creates missing files, dirs for developer, checks permissions, login to docker registry, checks inotify limits for tools such as webpack and other file watchers.

- split `lets.yaml` when it becomes big.
    If `lets.yaml` became bit, it would be great to split it in smaller, more specific files using `` for example