# Multi-stage docker build
# Build stage
FROM golang:alpine AS builder

LABEL maintainer="LitmusChaos"

ARG TARGETPLATFORM

ADD . /chaos-operator
WORKDIR /chaos-operator

RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) && \
    export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2)

RUN go env

RUN CGO_ENABLED=0 go build -buildvcs=false -o /output/chaos-operator -v ./main.go

# Packaging stage
# Image source: https://github.com/litmuschaos/test-tools/blob/master/custom/hardened-alpine/infra/Dockerfile
# The base image is non-root (have litmus user) with default litmus directory.
FROM litmuschaos/infra-alpine

LABEL maintainer="LitmusChaos"

ENV OPERATOR=/usr/local/bin/chaos-operator
COPY --from=builder /output/chaos-operator ${OPERATOR}

ENTRYPOINT ["/usr/local/bin/chaos-operator"]
