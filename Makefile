SHELL=bash
GOFILES=$(wildcard *.go)
BIN_NAME=humioctl
BIN_PATH=bin/$(BIN_NAME)
CLI_COMMAND ?= ""

$(BIN_PATH): $(GOFILES)

all: build

$(BIN_PATH): FORCE
	@echo "--> Building Humio CLI"
	go build -o $(BIN_PATH) cmd/humioctl/*.go

build: $(BIN_PATH)

clean:
	@echo "--> Cleaning"
	go clean
	@rm -rf bin/

snapshot:
	@echo "--> Building snapshot"
	goreleaser build --rm-dist --snapshot
	@rm -rf bin/

run: $(BIN_PATH)
	$(BIN_PATH) $(CLI_COMMAND)

.PHONY: build clean snapshot run FORCE

FORCE:
