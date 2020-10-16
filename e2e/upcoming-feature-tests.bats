#!/usr/bin/env bats

HUMIO_PORT=${HUMIO_PORT:-8081}
humioctl="$BATS_TEST_DIRNAME/../bin/humioctl --address=http://localhost:$HUMIO_PORT/"

load './node_modules/bats-support/load'
load './node_modules/bats-assert/load'

## Package Commands

@test "packages" {
  $humioctl packages
}

@test "package install valid humio/rubyapp from directory" {
  run $humioctl packages install developer ./packages/valid
  assert_success
}

@test "package listInstalled should contain humio/rubyapp" {
  run $humioctl packages listInstalled developer
  assert_success
  assert_output -p "humio/rubyapp"
}

@test "package uninstall humio/rubyapp" {
  run $humioctl packages uninstall developer humio/rubyapp
  assert_success
}

@test "package listInstalled should be empty" {
  run $humioctl packages listInstalled developer
  assert_success
  assert_line --index 0 -p "Listing installed packages in view developer" 
  refute_output -p "humio/rubyapp"
}

@test "package install invalid humio/rubyapp from directory should fail" {
  run $humioctl packages install developer ./packages/invalid
  assert_failure
}

@test "package install valid humio/rubyapp from zip" {
  run $humioctl packages install developer ./packages/humio--rubyapp--0.1.0.zip
  assert_success
}

@test "package listInstalled should contain humio/rubyapp" {
  run $humioctl packages listInstalled developer
  assert_success
  assert_output -p "humio/rubyapp"
}