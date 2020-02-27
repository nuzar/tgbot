SHELL = /bin/bash
BIN_NAME := "tgbot"
PKG := "github.com/nuzar/tgbot"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
DOCKER_HUB := "lyf0/tgbot:dev"

.PHONY: all build lint test clean run

all: build lint test

run: build
	@./$(BIN_NAME)

lint: build ## Lint the files
	@echo "check lint"
	@golangci-lint run

lintfix: ## Lint files and auto fix
	@echo "lint and autofix"
	@golangci-lint run --fix

test: ## Run unittests
	@echo "check test"
	@go test ${PKG_LIST}

build: ## Build the binary file
	@echo "make build"
	go build

clean: ## Remove previous build
	@echo "make clean"
	@rm -f $(BIN_NAME)

docker:
	@echo "build docker image"
	@docker build -t $() .
