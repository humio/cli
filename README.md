<img align="right" src="docs/images/cli-logo.png" style="width: 200px" />

# Humio CLI

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

## Installing a Community Parser

Find a parser at [humio/community](https://github.com/humio/community) and
install it using the CLI.

For instance if you wanted to install an AccessLog parser you could use.

```bash
humioctl parsers install accesslog
```

This would install the parser at: `humio/comminity/parsers/accesslog/default.yaml`.
Since log formats can vary slightly you can install one of the other variations:

```bash
humioctl parsers install accesslog/utc
```

Which would install the `humio/comminity/parsers/accesslog/utc.yaml` parser.


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

## Environment Variables

When use `humio login` it will write your API token and server address
into `.humioconfig` in your home dir.

The CLI will always have environment variables take precedence over
values form your `~/.humioconfig` file.

```bash
# Your account API token. You can find your token in Humio's UI under
# 'Your Account' in the account menu.
HUMIO_API_TOKEN=<token>

# A file containing the auth token to use for authorization
# useful in conjunction with Humio's generated root token file.
# If this is set it takes precedence over HUMIO_API_TOKEN.
HUMIO_TOKEN_FILE=<path>

# The address of the Humio server. E.g. https://cloud.humio.com/,
# or http://localhost:8080/
HUMIO_ADDR=<url>

# Disable color in terminal output
HUMIO_CLI_NO_COLOR=<bool>
```
