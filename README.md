<img align="right" src="docs/images/cli-logo.png" style="width: 200px" />

# Humio CLI

[![Build Status](https://github.com/humio/cli/workflows/CI/badge.svg)](https://github.com/humio/cli/actions?query=workflow%3ACI)
[![Release Status](https://github.com/humio/cli/workflows/goreleaser/badge.svg)](https://github.com/humio/cli/actions?query=workflow%3Agoreleaser)
[![Go Report Card](https://goreportcard.com/badge/github.com/humio/cli)](https://goreportcard.com/report/github.com/humio/cli)

A CLI for managing and sending data to Humio.

_This repository also contains Humio GoLang API client you can
use to build your own tools._

## Installation

### MacOS

```bash
$ brew tap humio/humio
$ brew install humioctl
```

### Linux (via Snapcraft)

```bash
$ sudo snap install humioctl
```

### Nix

```bash
$ nix-env -i humioctl
```
### Any OS (via Go)

```bash
$ go get github.com/humio/cli/cmd/humioctl
```

### Manual installation

Download the latest release archive from the [releases](https://github.com/humio/cli/releases) page, unpack and enjoy!

## Usage

To get started just write:

```bash
$ humioctl
```

and you will be asked to connect to your Humio cluster.
To list all commands use:

```bash
$ humioctl help
```

Will print all the different options and commands currently supported.

## Sending Logs

```bash
$ tail -f /var/log/system.log | humio ingest

# or

$ humioctl ingest --tail=/var/log/system.log
```

You can have Humio's UI open and tail the newly imported data using the `-open`
flag.

## @label and @session

When streaming data into Humio using the `ingest` command all events
will be annotated with `@label` and `@session` attributes.  `@label`
gives you a way to tag your streams to easily find them again e.g.:

```bash
$ humioctl ingest -label=work-related -tail=/var/log/mysql.log
```

```java
@label = "work-related" | groupby(loglevel)
```

`@session` is a unique id that is generated for each execution of the `humio`
binary. This allows you to find results for this session and nothing else.

## Profiles and Environment Variables

To make it easier to switch between different Humio clusters, you can
configure a profile for each cluster. The configuration file, containing the
API token and server address for all profiles will be default be saved in
`$HOME/.humio/config.yaml`.

Adding a profile and making it the new default can be done using:

```bash
$ humioctl profiles add my-profile
$ humioctl profiles set-default my-profile

It is also possible to use environment variables, and these will take
precendence over the default profile.

```bash
# Your account API token. You can find your token in Humio's UI under
# 'Your Account' in the account menu.
HUMIO_TOKEN=<token>

# The address of the Humio server. E.g. https://cloud.humio.com/,
# or http://localhost:8080/
HUMIO_ADDRESS=<url>

# If access to the Humio server requires trusting a specific Certificate Authority,
# for validating the certificate, you can specify CA certificate in PEM format.
# You can either point to a file with the certificate or provide it directly.
HUMIO_CA_CERTIFICATE=<ca-certificate>

# If access to the Humio server uses an untrusted certificate and you
# are unable to provide a CA certificate, you can disable TLS certificate verification.
# NB: This should only ever be used on test/sandbox clusters where you are in full
# control of the involved systems and underlying network.
# Do not use this for prodution use-cases.
HUMIO_INSECURE=<bool>
```
