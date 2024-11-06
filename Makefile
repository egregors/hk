SHELL := /bin/bash

PROJECT_NAME := "hk"
PKG := "github.com/egregors/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

BIN := "t-hk-srv"

## Common tasks

.PHONY: run
run:  ## Run dev version with watch dog (required watchexec)
	@watchexec -r -e go -- go run cmd/dev/main.go

.PHONY: build
build:  ## Build server and put bin into ~/go/bin/
	@go build -o $(BIN) cmd/prod/main.go
	mv ./$(BIN) ~/go/bin/

.PHONY: lint
lint:  ## Lint the files
	@golangci-lint run

## Help

.PHONY: help
help:  ## Show help message
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
	printf "%s\n\n" "Usage: make [task]"; \
	printf "%-20s %s\n" "task" "help" ; \
	printf "%-20s %s\n" "------" "----" ; \
	for help_line in $${help_lines[@]}; do \
		IFS=$$':' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf '\033[36m'; \
		printf "%-20s %s" $$help_command ; \
		printf '\033[0m'; \
		printf "%s\n" $$help_info; \
	done