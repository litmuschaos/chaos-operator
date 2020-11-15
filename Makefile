# Makefile for building Chaos Operator
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

IS_DOCKER_INSTALLED = $(shell which docker >> /dev/null 2>&1; echo $$?)

# list only our namespaced directories
PACKAGES = $(shell go list ./... | grep -v '/vendor/')

# docker info
DOCKER_REPO ?= litmuschaos
DOCKER_IMAGE ?= chaos-operator
DOCKER_TAG ?= latest

.PHONY: all
all: deps format lint build test dockerops dockerops-amd64

.PHONY: help
help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake deps      -- sets up dependencies for image build"
	@echo "\tmake gotasks   -- builds the chaos operator binary"
	@echo "\tmake dockerops -- builds & pushes the chaos operator image"
	@echo ""

.PHONY: deps
deps: _build_check_docker godeps unused-package-check

.PHONY: godeps
godeps:
	@echo ""
	@echo "INFO:\tverifying dependencies for chaos operator build ..."
	@go get -u -v golang.org/x/lint/golint
	@go get -u -v golang.org/x/tools/cmd/goimports

.PHONY: _build_check_docker
_build_check_docker:
	@if [ $(IS_DOCKER_INSTALLED) -eq 1 ]; \
		then echo "" \
		&& echo "ERROR:\tdocker is not installed. Please install it before build." \
		&& echo "" \
		&& exit 1; \
		fi;

.PHONY: gotasks
gotasks: format lint build

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
	@./build/go-multiarch-build.sh github.com/litmuschaos/chaos-operator/cmd/manager

.PHONY: test
test:
	@echo "------------------"
	@echo "--> Run Go Test"
	@echo "------------------"
	@go test ./... -coverprofile=coverage.txt -covermode=atomic -v

.PHONY: dockerops
dockerops:
	@echo "------------------"
	@echo "--> Build & Push chaos-operator docker image"
	@echo "------------------"
	sudo docker buildx build --file build/Dockerfile --progress plane --platform linux/arm64,linux/amd64 --tag $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) .
	REPONAME=$(DOCKER_REPO) IMGNAME=$(DOCKER_IMAGE) IMGTAG=$(DOCKER_TAG) ./buildscripts/push

.PHONY: dockerops-amd64
dockerops-amd64:
	@echo "--------------------------------------------"
	@echo "--> Build chaos-operator amd-64 docker image"
	@echo "--------------------------------------------"
	sudo docker build --file build/Dockerfile --tag $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) . --build-arg TARGETARCH=amd64
	@echo "--------------------------------------------"
	@echo "--> Push chaos-operator amd-64 docker image"
	@echo "--------------------------------------------"	
	sudo docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

unused-package-check:
	@echo "------------------"
	@echo "--> Check unused packages for the chaos-operator"
	@echo "------------------"
	@tidy=$$(go mod tidy); \
	if [ -n "$${tidy}" ]; then \
		echo "go mod tidy checking failed!"; echo "$${tidy}"; echo; \
	fi

