# Multi-stage docker build
# Build stage
FROM golang:1.16 AS builder

LABEL maintainer="LitmusChaos"

ARG TARGETPLATFORM

ADD . /chaos-operator
WORKDIR /chaos-operator

RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) && \
    export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2)

RUN go env

RUN CGO_ENABLED=0 go build -o /output/chaos-operator -v ./cmd/manager/

FROM golang:alpine as cert
RUN apk --no-cache add ca-certificates

# Packaging stage
FROM scratch

LABEL maintainer="LitmusChaos"

COPY --from=cert /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV OPERATOR=/usr/local/bin/chaos-operator \
    USER_UID=1001 \
    USER_NAME=chaos-operator

COPY --from=builder /output/chaos-operator ${OPERATOR}

ENTRYPOINT ["/usr/local/bin/chaos-operator"]

USER ${USER_UID
