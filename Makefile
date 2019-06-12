SHELL=bash
GOFILES=$(wildcard *.go)
BIN_NAME=humioctl
BIN_PATH=bin/$(BIN_NAME)
CLI_COMMAND ?= ""

$(BIN_PATH): $(GOFILES)

all: build

$(BIN_PATH): $(GOFILES)
	@echo "--> Building Humio CLI"
	go build -o $(BIN_PATH) main.go

build: $(BIN_PATH)

get:
	@echo "--> Fetching dependencies"
	go get -v

test: get
	@echo "--> Testing"
	go test -v ./...

clean-integration:
	docker-compose down -v

test-integration: clean-integration
	docker-compose up --abort-on-container-exit --exit-code-from cli

clean:
	@echo "--> Cleaning"
	go clean
	@rm -rf bin/

dist: clean
	@echo "--> Tagging Git & Releasing"
	./scripts/dist.sh

run: $(BIN_PATH)
	$(BIN_PATH) $(CLI_COMMAND)

.PHONY: build get clean dist run
