---
id: development
title: Development
---

## Build

We are suggesting to use `lets-dev` name for development binary, so you could
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

It will install `lets-dev` to /usr/local/bin/lets-dev, or wherever you`ve specified in path, and set version to development with current tag and timestamp

## Test

To run all tests:

```bash
lets test
```

To run unit tests:

```bash
lets test-unit
```

To get coverage:

```bash
lets coverage
```

To test `lets` output we are using `bats` - bash automated testing:

```bash
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

### Prerelease

If you are not ready to release a new version yet, it is possible to create a prerelease version.

Prerelease version is no visible to `install.sh` script and you can be sure that no one will get this version accidentiall.

Also you do not need to revoke published version if it has some critical bugs.

To create a prerelease version you need to append a `-rcN` suffix to next version, for example:

```bash
lets release 0.0.1-rc1 -m "pre: implement some new feature"
```

This will create a `0.0.1-rc1` tag and push it to github. Github will create a prerelease version and build all the binaries.

Once you are ready to release a new version, just create a normal release:

```bash
lets release 0.0.1 -m "implement some new feature"
```

## Versioning

`lets` releases must be backward compatible. That means every new `lets` release must work with old configs.

For situations like e.g. new functionality, there is a `version` in `lets.yaml` which specifies **minimum required** `lets` version.

If `lets` version installed on the user machine is less than the one specified in config it will show and error and ask the user to upgrade `lets` version.
