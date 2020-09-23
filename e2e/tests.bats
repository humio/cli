#!/usr/bin/env bats

humioctl="$BATS_TEST_DIRNAME/../bin/humioctl --address=http://localhost:8081/"

## Health Commands

@test "health" {
  $humioctl health
}

## Cluster Commands

@test "cluster nodes list" {
  run $humioctl cluster nodes list
  grep "localhost:8080" <<< "$output"
}

@test "cluster nodes show" {
  run $humioctl cluster nodes show 1
  [[ "$status" -eq 0 ]]
  grep "UUID" <<< "$output"
}

## View Commands

@test "views list" {
  run $humioctl views list
  [[ "$status" -eq 0 ]]
  grep "humio-audit" <<< "$output"
  grep "sandbox" <<< "$output"
}

@test "views show" {
  $humioctl views show humio-audit
}

@test "views show <unknown>" {
  run $humioctl views show foo
  [[ "$status" -eq 1 ]]
}

## User Commands

@test "users show" {
  run $humioctl users list
  [[ "$status" -eq 0 ]]
  grep "developer" <<< "$output"
}

@test "users add" {
  $humioctl users add --email "foo@acme.org" --root true --name "Anders" anders
}

@test "users add <already exists>" {
  run $humioctl users add --email "bar@acme.org" --root true --name "Peter" peter
  [[ "$status" -eq 0 ]]

  run $humioctl users add --email "bar@acme.org" --root true --name "Peter" peter
  [[ "$status" -eq 1 ]]
}