SHELL=bash
GOFILES=$(wildcard *.go)
BIN_NAME=humioctl
BIN_PATH=bin/$(BIN_NAME)
CLI_COMMAND ?= ""
SCHEMA_CLUSTER?=${HUMIO_ENDPOINT}
SCHEMA_CLUSTER_API_TOKEN?=${HUMIO_TOKEN}

$(BIN_PATH): $(GOFILES)

all: build

$(BIN_PATH): FORCE
	go build -o $(BIN_PATH) ./cmd/humioctl

build: $(BIN_PATH)

clean:
	go clean
	@rm -rf bin/

snapshot:
	goreleaser build --clean --snapshot
	@rm -rf bin/

run: $(BIN_PATH)
	$(BIN_PATH) $(CLI_COMMAND)

update-schema:
	go run github.com/suessflorian/gqlfetch/gqlfetch@607d6757018016bba0ba7fd1cb9fed6aefa853b5 --endpoint ${SCHEMA_CLUSTER}/graphql --header "Authorization=Bearer ${SCHEMA_CLUSTER_API_TOKEN}" > internal/api/humiographql/schema/_schema.graphql
	printf "# Fetched from version %s" $$(curl --silent --location '${SCHEMA_CLUSTER}/api/v1/status' | jq -r ".version") >> internal/api/humiographql/schema/_schema.graphql

e2e: $(BIN_PATH)
	./e2e/run.bash

e2e-upcoming: $(BIN_PATH)
	./e2e/run-upcoming-features.bash

.PHONY: build clean snapshot run e2e e2e-upcoming FORCE

FORCE:
