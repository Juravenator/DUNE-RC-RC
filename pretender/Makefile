SHELL:=/bin/bash
.DEFAULT_GOAL := build

src_files := $(shell find -type f -name '*.go')

#######
# hard to remember golang vendoring commands
#######
.PHONY: dependencies
dependencies: ${src_files}
	# Update all direct and indirect dependencies to latest minor or patch upgrades
	go get -u ./...
	# Prune any no-longer-needed dependencies from go.mod and add any dependencies needed for other combinations of OS, architecture, and build tags
	go mod tidy
	# create vendor directory for offline building
	# go build does NOT use the vendor directory by default, this directory a fail-safe for takedowns of dependencies
	# go mod vendor

.PHONY: lint
lint: ${src_files}
	gofmt -s -w .

build/pretender: dependencies
	go build -o build/pretender -tags=production ./cmd/pretender.go

.PHONY: build
build: build/pretender

#######
# clean
#######
.PHONY: clean
clean:
	git clean -dfX
