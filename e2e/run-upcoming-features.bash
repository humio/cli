#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd $DIR

npm install 1>/dev/null

echo "==> Testing upcoming features"
HUMIO_PORT=8080 npx bats $DIR/upcoming-feature-tests.bats -p
