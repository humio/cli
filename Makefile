SHELL=bash
GOFILES=$(wildcard *.go)
BIN_NAME=humio
BIN_PATH=bin/$(BIN_NAME)
CLI_COMMAND ?= ""

$(BIN_PATH): $(GOFILES)

build: $(BIN_PATH)
	@echo "--> Building Humio CLI"
	go build -o bin/$(BIN_NAME) $(GOFILES)

get:
	@echo "--> Fetching dependencies"
	govendor get

test:
	@echo "--> Testing"
	govendor test +local

clean:
	@echo "--> Cleaning"
	go clean
	@rm -rf bin/

dist: clean
	@echo "--> Tagging Git & Releasing"
	./scripts/dist.sh

run:
	$(MAKE) build
	$(BIN_PATH) $(CLI_COMMAND)

.PHONY: build get clean dist run
