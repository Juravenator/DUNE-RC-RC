SHELL:=/bin/bash
.DEFAULT_GOAL := help

#
## main targets
#

.PHONY: build
build:
	make -C controllers/daq-app-manager build

include .makefile/dockerized.mk
include .makefile/help.mk