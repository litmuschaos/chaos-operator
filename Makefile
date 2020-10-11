# Makefile for building Chaos Operator
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

IS_DOCKER_INSTALLED = $(shell which docker >> /dev/null 2>&1; echo $$?)

OPERATOR_DOCKER_REPO ?= litmuschaos

ADMISSION_CONTROLLERS_DOCKER_IMAGE?=admission-controllers

# Specify the name for the binaries
WEBHOOK=admission-controllers

# Specify the date o build
BUILD_DATE = $(shell date +'%Y%m%d%H%M%S')

# list only our namespaced directories
PACKAGES = $(shell go list ./... | grep -v '/vendor/')

# docker info

OPERATOR_DOCKER_IMAGE ?= chaos-operator

IMAGE_TAG ?= latest

.PHONY: all
all: deps gotasks test dockerops admission-controllers-image

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
	@go build -o ${GOPATH}/src/github.com/litmuschaos/chaos-operator/build/_output/bin/chaos-operator -gcflags all=-trimpath=${GOPATH} -asmflags all=-trimpath=${GOPATH} github.com/litmuschaos/chaos-operator/cmd/manager

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
	sudo docker build -f build/Dockerfile -t $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) --build-arg BUILD_DATE=${BUILD_DATE} .
	REPONAME=$(DOCKER_REPO) IMGNAME=$(DOCKER_IMAGE) IMGTAG=$(DOCKER_TAG) ./buildscripts/push

unused-package-check:
	@echo "------------------"
	@echo "--> Check unused packages for the chaos-operator"
	@echo "------------------"
	@tidy=$$(go mod tidy); \
	if [ -n "$${tidy}" ]; then \
		echo "go mod tidy checking failed!"; echo "$${tidy}"; echo; \
	fi

.PHONY: admission-controllers-image
admission-controllers-image:
	@echo "----------------------------"
	@echo -n "--> admission-controllers image: "
	@echo "${OPERATOR_DOCKER_REPO}/${ADMISSION_CONTROLLERS_DOCKER_IMAGE}:${IMAGE_TAG}"
	@echo "----------------------------"
	@PNAME=${WEBHOOK} CTLNAME=${WEBHOOK} sh -c "'$(PWD)/buildscripts/build.sh'"
	@cp bin/${WEBHOOK}/${WEBHOOK} buildscripts/admission-controllers/
	@cd buildscripts/${WEBHOOK} && sudo docker build -t ${OPERATOR_DOCKER_REPO}/${ADMISSION_CONTROLLERS_DOCKER_IMAGE}:${IMAGE_TAG} --build-arg BUILD_DATE=${BUILD_DATE} .
	@rm buildscripts/${WEBHOOK}/${WEBHOOK}
