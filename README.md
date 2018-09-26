[![Build Status](https://travis-ci.org/hanks/awsudo-go.svg?branch=master)](https://travis-ci.org/hanks/awsudo-go)

# AWSUDO-GO

A tool to automate AWS API access using a SAML compliant identity provider. This project is inspired by the ruby version [https://github.com/electronicarts/awsudo](https://github.com/electronicarts/awsudo), and rewrite it in golang, and do some enhancements.

## Prerequisite

* okta account
* awscli, 1.16.17+
* AWS Role setup
* Docker, for development
* Linux/macOS, not test in Windows now

## Enhancements

1. Rewrite with golang, use cross-compile to create one single binary with no dependencies
2. Merge `.awsudo` and `.awsudo_roles` configs to single one config by `TOML`, and add `awsudo configure` command to setup the config, just like `aws configure` style
3. Merge `awsudo agent` and `awsudo` command, just use one command to do all tasks
4. Add `AWS Role Session Duration` and `Awsudo Agent Expiration` support, more secure
5. Add command `awsudo shutdown` to close agent server gracefully
6. Add output log for monitor and debug
7. Add docker support for development

## Downgrades

1. Only support okta now

## Installation

```bash
make install
```

and uninstall by:

```bash
make uninstall
```

## Usage

```bash
awsudo configure
awsudo prod-env aws s3 ls
```

It will call assumeRole API to fetch the credentials, and set them as environment variables, then
to run aws command.

### Demo

![demo.gif](./docs/images/demo.gif)

## Development

* `make test`, run unit test, coverage test, static analytics
* `make run`, just to run help command to as a start point
* `make build`, cross compile binaries, and put into `dist/bin` directory
* `make debug`, use `dlv` to do the `gdb-style` debug
* `make dev`, build docker image used in dev

## Contribution

**Waiting for your pull request**

## Licence

MIT Licence
