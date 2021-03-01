SHELL:=/bin/bash
.DEFAULT_GOAL := help

##
### Main targets
##

.PHONY: build
build: ## build the software
	make -C controllers/daq-app-manager build

include .makefile/dockerized.mk
include .makefile/help.mk