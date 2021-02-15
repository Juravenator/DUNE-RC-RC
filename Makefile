SHELL:=/bin/bash
.DEFAULT_GOAL := help

#
# main variables
#

# PROJECT_NAME=sql-exporter
# VERSION ?= $(shell git describe --always)
# ARCH ?= amd64
# RPM_NAME ?= cactus-${PROJECT_NAME}-${VERSION}-${RELEASE}.${ARCH}.rpm

# INSTALL_ROOT ?= /opt/ral/inventory-manager
# RPM_ROOT_FOLDER ?= rpmroot
# RPM_INSTALL_ROOT = ${RPM_ROOT_FOLDER}${INSTALL_ROOT}

##
### main targets
##

# .PHONY: lint
# lint: ## lint all code
# 	make -C htdocs lint
# 	make -C api lint
# 	# make -C syncer lint

# .PHONY: build
# build: ## compile all code
# 	make -C htdocs build
# 	make -C api build
# 	# make -C syncer build

# .PHONY: test
# test: ## test all code
# 	make -C htdocs test
# 	make -C api test
# 	# make -C syncer test

# .PHONY: clean
# clean: ## clean entire project
# 	make -C htdocs clean
# 	make -C api clean
# 	# make -C syncer clean

# .PHONY: rpm
# rpm: ${RPM_NAME}
# ${RPM_NAME}: build
# 	rm -rf ${RPM_ROOT_FOLDER}
# 	mkdir -p ${RPM_INSTALL_ROOT}/htdocs
# 	cp -r api ${RPM_INSTALL_ROOT}
# 	cp -r htdocs/build ${RPM_INSTALL_ROOT}/htdocs

# 	cd ${RPM_ROOT_FOLDER} && fpm \
# 	-s dir \
# 	-t rpm \
# 	-n ral-inventory-manager \
# 	-v ${VERSION} \
# 	-a ${ARCH} \
# 	-m "<cactus@cern.ch>" \
# 	--vendor CERN \
# 	--description "SQL exporter for the L1 Online Software at P5" \
# 	--url "https://gitlab.cern.ch/cms-cactus/ops/monitoring/sql-exporter" \
# 	--provides cactus_$(subst -,_,${PROJECT_NAME}) \
# 	.=/ && mv *.rpm ..



include .makefile/dockerized.mk
include .makefile/help.mk