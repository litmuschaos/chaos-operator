# Makefile for building Chaos Exporter
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

IS_DOCKER_INSTALLED = $(shell which docker >> /dev/null 2>&1; echo $$?)

# list only our namespaced directories
PACKAGES = $(shell go list ./... | grep -v '/vendor/')

# stable release version of operator-sdk 
OPERATOR_SDK_VERSION ?= "v0.8.0"

# docker image details 
DOCKER_REPO ?= "litmuschaos"
DOCKER_IMAGE ?= "chaos-operator"
DOCKER_TAG ?= "ci"

.PHONY: all
all: format lint build test push 

.PHONY: help
help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake all   -- [default] builds the chaos exporter container"
	@echo ""

.PHONY: godeps
godeps:
	@echo ""
	@echo "INFO:\tverifying dependencies for chaos exporter build ..."
	@go get -u -v golang.org/x/lint/golint
	@go get -u -v golang.org/x/tools/cmd/goimports
	@go get -u -v github.com/golang/dep/cmd/dep

.PHONY: _build_check_docker
_build_check_docker:
	@if [ $(IS_DOCKER_INSTALLED) -eq 1 ]; \
		then echo "" \
		&& echo "ERROR:\tdocker is not installed. Please install it before build." \
		&& echo "" \
		&& exit 1; \
		fi;

.PHONY: _install_operator_sdk
_install_operator_sdk:
	@echo "------------------"
	@echo "Installing Operator SDK Stable Release Binary"
	curl -OJL https://github.com/operator-framework/operator-sdk/releases/download/\
        ${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu
	chmod +x operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu && \
        sudo cp operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu /usr/local/bin/operator-sdk && \
        rm operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu
          

.PHONY: deps
deps: _build_check_docker _install_operator_sdk godeps

.PHONY: format
format:
	@echo "------------------"
	@echo "--> Running go fmt"
	@echo "------------------"
	@go fmt $(PACKAGES)

.PHONY: lint
lint:
	@echo "------------------"
	@echo "--> Running golint"
	@echo "------------------"
	@golint $(PACKAGES)
	@echo "------------------"
	@echo "--> Running go vet"
	@echo "------------------"
	@go vet $(PACKAGES)

.PHONY: build  
build:
	@echo "------------------"
	@echo "--> Build Chaos Operator"
	@echo "------------------"
	operator-sdk build $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: test
test:
	@echo "------------------"
	@echo "--> Run Go Test"
	@echo "------------------"
	@go test ./... -v 

.PHONY: push
dockerops: 
	@echo "------------------"
	@echo "--> Push chaos-operator docker image" 
	@echo "------------------"
	REPONAME=$(DOCKER_REPO) IMGNAME=$(DOCKER_IMAGE) IMGTAG=$(DOCKER_TAG) ./buildscripts/push
