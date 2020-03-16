SHELL=bash
GOFILES=$(wildcard *.go)
BIN_NAME=humioctl
BIN_PATH=bin/$(BIN_NAME)
CLI_COMMAND ?= ""

$(BIN_PATH): $(GOFILES)

all: build

$(BIN_PATH): FORCE
	@echo "--> Building Humio CLI"
	go build -o $(BIN_PATH) main.go

build: $(BIN_PATH)

get:
	@echo "--> Fetching dependencies"
	go get

test:
	@echo "--> Testing"
	go test

clean:
	@echo "--> Cleaning"
	go clean
	@rm -rf bin/

dist: clean
	@echo "--> Tagging Git & Releasing"
	./scripts/dist.sh

run: $(BIN_PATH)
	$(BIN_PATH) $(CLI_COMMAND)

.PHONY: build get clean dist run FORCE

FORCE:
