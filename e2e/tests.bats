#!/usr/bin/env bats

HUMIO_PORT=${HUMIO_PORT:-8081}
humioctl="$BATS_TEST_DIRNAME/../bin/humioctl --address=http://localhost:$HUMIO_PORT/"

load './node_modules/bats-support/load'
load './node_modules/bats-assert/load'

## Health Commands

@test "health" {
  run $humioctl health
}

## Cluster Commands

@test "cluster nodes list" {
  run $humioctl cluster nodes list
  assert_output -p "localhost:8080"
}

@test "cluster nodes show" {
  run $humioctl cluster nodes show 1
  assert_success
  assert_output -p "UUID"
}

## View Commands

@test "views list" {
  run $humioctl views list
  assert_success
  assert_output -p "humio-audit"
  assert_output -p "sandbox"
}

@test "views show" {
  $humioctl views show humio-audit
}

@test "views show <unknown>" {
  run $humioctl views show foo
  assert_failure 1
}

## User Commands

@test "users show" {
  run $humioctl users show developer
  assert_success
  assert_output -p "developer"
}

@test "users add" {
  $humioctl users add --email "foo@acme.org" --root true --name "Anders" anders
}

@test "users add <already exists>" {
  run $humioctl users add --email "bar@acme.org" --root true --name "Peter" peter
  assert_success

  run $humioctl users add --email "bar@acme.org" --root true --name "Peter" peter
  assert_failure 1
}

@test "users list" {
  run $humioctl users add --email "quux@acme.org" --root false --name "Jens" jens
  assert_success

  run $humioctl users list
  assert_success
  assert_output -p "Jens"
}

@test "users update" {
  $humioctl users add --email "abc@acme.org" --name "Odin" odin
  $humioctl users update --name 'Othello' odin
  
  run $humioctl users show odin 

  assert_success
  assert_output -p 'Othello'
}

# Status Commands

@test "status" {
  run $humioctl status
  assert_success
  assert_output -p 'Version'
}

# License Commands

@test "license install" {
  skip # TODO: This fails in the test build susseeds with manual tests.
       # we keep getting 'error installing license: Invalid input' in the test.
       # Why is that?
  
  # Humio will accept an expired license, but limit usage.
  run $humioctl license install "$BATS_TEST_DIRNAME/licenses/expired.pem"
  assert_success
  assert_output "OK"
}

@test "license install <invalid license>" {
  run $humioctl license install "$BATS_TEST_DIRNAME/licenses/invalid.pem"
  assert_failure
}

@test "license show" {
  $humioctl license show
}

# Alerts

@test "alerts list" {
  $humioctl alerts list humio
}

# Actions

@test "actions list" {
  $humioctl actions list humio
}

# Ingest Token Commands

@test "ingest-tokens add" {
  run $humioctl ingest-tokens add humio "test123" --parser "kv"
  assert_success
  assert_output -p "test123"
  assert_output -p "kv"
}

@test "ingest-tokens update change parser" {
  $humioctl ingest-tokens add humio "updateParser" --parser "kv"
  run $humioctl ingest-tokens update humio "updateParser" --parser "json"
  assert_output -p "updateParser"
  assert_output -p "json"
}

@test "ingest-tokens update remove parser" {
  $humioctl ingest-tokens add humio "removeParser" --parser "kv"
  run $humioctl ingest-tokens update humio "removeParser"
  assert_output -p "removeParser"
}

@test "ingest-tokens list" {
  $humioctl ingest-tokens add humio "foo"
  run $humioctl ingest-tokens list humio
  assert_success
  assert_output -p "Assigned Parser"
}

# Parser Commands

@test "parser install" {
  run $humioctl parsers install --file parsers/accesslog2.yaml humio
  assert_success
}

@test "parser remove" {
  $humioctl parsers install --file parsers/accesslog2.yaml --name accesslog3 humio
  run $humioctl parsers remove humio accesslog3
  assert_success
  assert_output -p "Successfully removed parser \"accesslog3\" from repository \"humio\""
}
