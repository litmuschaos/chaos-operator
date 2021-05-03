# Makefile for building Chaos Operator
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

IS_DOCKER_INSTALLED = $(shell which docker >> /dev/null 2>&1; echo $$?)

# docker info
DOCKER_REPO ?= litmuschaos
DOCKER_IMAGE ?= chaos-operator
DOCKER_TAG ?= latest

.PHONY: all
all: deps unused-package-check build-chaos-operator test

.PHONY: help
help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake deps      -- sets up dependencies for image build"
	@echo ""

.PHONY: deps
deps: _build_check_docker godeps

.PHONY: _build_check_docker
_build_check_docker:
	@if [ $(IS_DOCKER_INSTALLED) -eq 1 ]; \
		then echo "" \
		&& echo "ERROR:\tdocker is not installed. Please install it before build." \
		&& echo "" \
		&& exit 1; \
		fi;

.PHONY: godeps
godeps:
	@echo ""
	@echo "INFO:\tverifying dependencies for chaos operator build ..."
	@go get -u -v golang.org/x/lint/golint
	@go get -u -v golang.org/x/tools/cmd/goimports

.PHONY: test
test:
	@echo "------------------"
	@echo "--> Run Go Test"
	@echo "------------------"
	@go test ./... -coverprofile=coverage.txt -covermode=atomic -v

unused-package-check:
	@echo "------------------"
	@echo "--> Check unused packages for the chaos-operator"
	@echo "------------------"
	@tidy=$$(go mod tidy); \
	if [ -n "$${tidy}" ]; then \
		echo "go mod tidy checking failed!"; echo "$${tidy}"; echo; \
	fi

gofmt-check:
	@echo "------------------"
	@echo "--> Check unused packages for the chaos-operator"
	@echo "------------------"
	@gfmt=$$(gofmt -s -l . | wc -l); \
	if [ "$${gfmt}" -ne 0 ]; then \
		echo "The following files were found to be not go formatted:"; \
   		gofmt -s -l .; \
   		exit 1; \
  	fi

.PHONY: build-chaos-operator build-chaos-operator-amd64 push-chaos-operator

build-chaos-operator:
	@docker buildx build --file build/Dockerfile --progress plane  --no-cache --platform linux/arm64,linux/amd64 --tag $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) .

build-for-amd64:
	@docker build -f build/Dockerfile  --no-cache -t $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) .  --build-arg TARGETPLATFORM="linux/amd64"

push-chaos-operator:
	@docker buildx build --file build/Dockerfile --progress plane --no-cache --push --platform linux/arm64,linux/amd64 --tag $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) .
