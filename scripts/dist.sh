#!/bin/bash

set -e

read -p "Releasing. Are you sure? [y/N]" -n 1 -r
echo    # (optional) move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi

scripts/bump_version.sh
goreleaser --rm-dist
