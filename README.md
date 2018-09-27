<img align="right" src="docs/images/cli-logo.png" style="width: 200px" />

# Humio CLI

A CLI for searching, streaming logs to, and manageing Humio.

__THIS IS WORK IN PROGRESS!__

You are welcome to try it out and contribute.  Warning: this readme is
also WIP.

## Usage

```bash
$ humio help
```

Will print all the different options and commands currently supported.

## Humio Env File

`humio` will look in the following file: `~/.humio-cli.env` for the
url, token, and repo options. This will save you added the options to
every command you type.

Example:

```bash
$> cat ~/.humio-cli.env
HUMIO_URL=https://cloud.humio.com
HUMIO_API_TOKEN=<token>
HUMIO_REPO=myrepo
```

You don't need all 3 env. variables, so if you often work in one Humio
system, but in different repositories leave out the `HUMIO_REPO` env.


## @name and @session

When streaming data into Humio using the `ingest` command all events
will be annotated with `@name` and `@session` attributes.  `@name`
gives you a way to tag your streams atito easily find them again,
e.g.:

```bash
$ humio -n work-related /usr/local/share/mysql.conf
```

```java
@name = "work-related" | groupby(loglevel)
```

`@session` is a unique id that is generated for each execution of the `humio`
binary. This allows you to find results for this session and nothing else.

## Developer Setup

```bash
$ make get
$ make build
```

### Making a new release

```bash
$ make dist
```
