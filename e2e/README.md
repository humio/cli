# E2E Tests for humioctl

Runs a test suite against a Humio docker container.

The bats-code bash testing framwork is used to write tests.
The test runner is installed through NPM (see Setup below).

## Requirements

- NPM
- Docker (Optional)

## Setup

npm install

## Running the Tests

Make sure to build the CLI before running the tests:

```
$ make build
```

in the root of the project

### Humio tests in docker

```shell
$ ./run.bash
```

### Humio tests in running locally

```
$ ./run.bash --skip-humio
```

This is useful when developing tests so you don't
have to wait for the docker containers starting and stopping.