<img align="right" src="docs/images/cli-logo.png" style="width: 200px" />

# Humio CLI

A CLI for searching and streaming logs to Humio.

__THIS IS WORK IN PROGRESS__

You are welcome to try it out and contribute.
Warning: this readme is also WIP.

## Usage

```bash
$ humio -t <HUMIO_API_TOKEN> /var/log/system.log
```

```bash
$ log stream | humio
```

You can also store your `HUMIO_API_TOKEN` as an environment variable.
Here is an example for setting it in `Bash`.

```bash
$ echo "export HUMIO_API_TOKEN=<API_TOKEN>" >> ~/.bashrc
```

By default all logs are sent to your _Scratch_ (Beta Feature) dataspace. You can use the
`-d` flag to specify another dataspace:

```bash
$ echo "Hello Humio" | humio -d another-dataspace [...]
```

## @name and @session

All events are will be annotated with `@name` and `@session` attributes.
`@name` gives you a way to tag your streams atito easily find them again, e.g.:

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
