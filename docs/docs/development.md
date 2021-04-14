---
id: development
title: Development
---

## Build

We are suggesting to use `lets-dev` name for development binary, so u could
have stable `lets` version untouched.

To build a binary:

```bash
go build -o lets-dev *.go
```

To install in system

```bash
go build -o lets-dev *.go && sudo mv ./lets-dev /usr/local/bin/lets-dev
```

Or if you already have `lets` installed in your system:

```bash
lets build-and-install [--path=<path>]
```
`path` - your custom executable $PATH, defaults to `/usr/local/bin`

After install - check version of lets - `lets-dev --version` - it should be development

It will install `lets-dev` to /usr/local/bin/lets-dev, or whereever u`ve specified in path, and set version to development with current tag and timestamp

## Test

To run all tests:

```shell script
lets test
```

To run unit tests:

```shell script
lets test-unit
```

To get coverage:

```shell script
lets coverage
```

To test `lets` output we using `bats` - bash automated testing:

```shell script
lets test-bats

# or run one test

lets test-bats global_env.bats
```

## Release

To release a new version:

```bash
lets release 0.0.1 -m "implement some new feature"
```

This will create an annotated tag with 0.0.1 and run `git push --tags`


## Versioning

`lets` releases must be backward compatible. That means every new `lets` release must work with old configs.

For situations like e.g. new functionality, there is a `version` in `lets.yaml` which specifies **minimum required** `lets` version.

If `lets` version installed on the user machine is less than the one specified in config it will show and error and ask the user to upgrade `lets` version.
