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
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4

LABEL maintainer="LitmusChaos"

ENV OPERATOR=/usr/local/bin/chaos-operator

COPY --from=builder /output/chaos-operator ${OPERATOR}
RUN chown 65534:0 ${OPERATOR} && chmod 755 ${OPERATOR}

USER 65534

ENTRYPOINT ["/usr/local/bin/chaos-operator"]
