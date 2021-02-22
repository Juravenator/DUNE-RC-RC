SHELL:=/bin/bash
.DEFAULT_GOAL := build

#
# main variables
#

src_files := $(shell find -name '*.go')

#
# main targets
#

.PHONY: build
build: build/consul-filler
	touch build

.PHONY: test
test:
	go test -tags testing -v ./...

.PHONY: clean
clean:
	rm -rf build ${RPM_NAME}

#
# dependencies
#

build/%: cmd/%.go $(src_files)
	go build -o $@ ./$(patsubst build/%,cmd/%.go,$@)
