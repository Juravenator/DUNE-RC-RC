SHELL:=/bin/bash
.DEFAULT_GOAL := help

include .makefile/dockerized.mk
include .makefile/help.mk