SHELL=bash
GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
BIN_NAME=humio

build:
	@echo "--> Building Humio CLI"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(BIN_NAME) $(GOFILES)

get:
	@echo "--> Fetching dependencies"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get .

clean:
	@echo "--> Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	@rm -rf bin/

dist: clean
	@echo "--> Tagging Git & Releasing"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) ./scripts/dist.sh

.PHONY: build get clean dist
