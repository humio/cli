#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd $DIR

npm install 1>/dev/null

SKIP_HUMIO=$1
PORT=8081

stop_humio() {
    if [[ "$SKIP_HUMIO" != "--skip-humio" ]]; then
        echo "==> Stopping Humio"
        docker stop $CONTAINER
    fi
}

CONTAINER="humioctl-e2e"
trap stop_humio SIGINT EXIT


if [[ "$SKIP_HUMIO" == "--skip-humio" ]]; then
    echo "==> Skipping Humio Startup"
else
    echo "==> Starting Humio (waiting)"

    docker run -d --rm --name "$CONTAINER" --env HUMIO_JVM_ARGS=-Xss2M \
        -p $PORT:8080 --ulimit="nofile=8192:8192" \
        humio/humio:latest

    MAX_RETRY=60
    RETRY_COUNT=0
    while ! curl http://localhost:8081/api/v1/status 
    do
    RETRY_COUNT=RETRY_COUNT+1
    if [[ "$MAX_RETRY" == "$RETRY_COUNT" ]]; then
        echo "Failed to start Humio"
        exit 1
    fi
    sleep 1
    done
fi

echo "==> Running Tests"

npx bats $DIR/tests.bats -p
