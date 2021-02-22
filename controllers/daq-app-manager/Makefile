SHELL:=/bin/bash
.DEFAULT_GOAL := help

#
# main variables
#

PYTHON ?= venv/bin/python3
PIP ?= venv/bin/pip3

##
### main targets
##

include .makefile/python.lint.mk
include .makefile/python.venv.mk

.PHONY: lint
lint: venv python.lint ## lint this code

.PHONY: build
build: venv ## build this code

.PHONY: test
test: venv ## run tests
	$(PYTHON) -m pytest --cov-report term --cov-report html:coverage --cov-config setup.cfg --cov=src/ test/

.PHONY: clean
clean: ## clean up
	rm -rf venv __pycache__

include .makefile/help.mk