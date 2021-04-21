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

ENV OPERATOR=/usr/local/bin/chaos-operator \
    USER_UID=1001 \
    USER_NAME=chaos-operator

COPY --from=builder /output/chaos-operator ${OPERATOR}
COPY build/bin /usr/local/bin

RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
