---
id: development
title: Development
---

## Build

To build a binary:

```bash
go build -o lets *.go
```

To install in system

```bash
go build -o lets *.go && sudo mv ./lets /usr/local/bin/lets
```

Or if you already have `lets` installed in your system:

```bash
lets build-and-install
```

After install - check version of lets - `lets --version` - it should be development

It will install `lets` to /usr/local/bin/lets and set version to development with current tag and timestamp

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
