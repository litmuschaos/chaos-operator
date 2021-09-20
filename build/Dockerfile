# Multi-stage docker build
# Build stage
FROM golang:alpine AS builder

LABEL maintainer="LitmusChaos"

ARG TARGETPLATFORM=linux/amd64

ADD . /chaos-operator
WORKDIR /chaos-operator

RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) && \
    export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2)

RUN go env

RUN CGO_ENABLED=0 go build -o /output/chaos-operator -v ./cmd/manager/

# Packaging stage
FROM uditgaurav/alpine:test

LABEL maintainer="LitmusChaos"

ENV OPERATOR=/usr/local/bin/chaos-operator
COPY --from=builder /output/chaos-operator ${OPERATOR}

ENTRYPOINT ["/usr/local/bin/chaos-operator"]
