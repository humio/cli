<img align="right" src="docs/images/cli-logo.png" style="width: 200px" />

# Humio CLI

A CLI for managing and sending data to Humio.

Visit [humio/community](https://github.com/humio/community) to find parsers
that you can install with the CLI. We plan on adding dashboards to the community
repo as well.

_This repository also contains Humio GoLang API client you can
use to build your own tools._

## Usage

```bash
$ humio help
```

Will print all the different options and commands currently supported.

## Setup

You need to some environment variables

```
export HUMIO_API_TOKEN=<token>
export HUMIO_ADDR=<humio-url> // E.g. https://cloud.humio.com/
```

## Sending Logs

```bash
$ tail -f /var/log/system.log | humio ingest

# or

$ humio ingest -tail=/var/log/system.log
```

You can have Humio's UI open and tail the newly imported data using the `-open`
flag.

## Installing a Community Parser

Find a parser at [humio/community](https://github.com/humio/community) and
install it using the CLI.

For instance if you wanted to install an AccessLog parser you could use.

```bash
humio parsers install accesslog
```

This would install the parser at: `humio/comminity/parsers/accesslog/default.yaml`.
Since log formats can vary slightly you can install one of the other variations:

```bash
humio parsers install accesslog/utc
```

Which would install the `humio/comminity/parsers/accesslog/utc.yaml` parser.


## @label and @session

When streaming data into Humio using the `ingest` command all events
will be annotated with `@label` and `@session` attributes.  `@label`
gives you a way to tag your streams to easily find them again e.g.:

```bash
$ humio ingest -label=work-related -tail=/var/log/mysql.log
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

```
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
