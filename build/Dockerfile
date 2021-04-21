# Multi-stage docker build
# Build stage
FROM golang:1.14 AS builder

LABEL maintainer="LitmusChaos"

ARG TARGETPLATFORM

ADD . /chaos-operator
WORKDIR /chaos-operator

RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) && \
    export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2)

RUN go env

RUN CGO_ENABLED=0 go build -o /output/chaos-operator -v ./cmd/manager/

# Packaging stage
FROM registry.access.redhat.com/ubi8/ubi-minimal:8.3

LABEL maintainer="LitmusChaos"

COPY --from=builder /output/chaos-operator /

RUN addgroup -S litmus && adduser -S -G litmus 1001

USER 1001

CMD ["./chaos-operator"]

