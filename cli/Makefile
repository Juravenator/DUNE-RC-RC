SHELL:=/bin/bash
.DEFAULT_GOAL := help

src_files := $(shell find -type f -name '*.go')

##
### Main targets
##
.PHONY: build
build: build/run-control ## build run-control CLI

.PHONY: lint
lint: ${src_files} ## lint code
	gofmt -s -w .
	go vet ./...

.PHONY: clean
clean: ## clean
	git clean -dfX
	rm -rf build

##
### Helper functions
##
.PHONY: dependencies
dependencies: ${src_files} ## update dependencies
	# Update all direct and indirect dependencies to latest minor or patch upgrades
	go get -u ./...
	# Prune any no-longer-needed dependencies from go.mod and add any dependencies needed for other combinations of OS, architecture, and build tags
	go mod tidy
	# create vendor directory for offline building
	# go build does NOT use the vendor directory by default, this directory is a fail-safe for takedowns of dependencies
	# go mod vendor

build/%: cmd/%.go $(src_files)
	go build -o $@ ./$(patsubst build/%,cmd/%.go,$@)

include .makefile/docker.mk
include .makefile/help.mk