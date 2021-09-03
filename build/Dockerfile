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

RUN CGO_ENABLED=0 go build -o /output/chaos-operator -v ./cmd/manager/

RUN apk --no-cache add ca-certificates

# add a user here because addgroup and adduser are not available in scratch
RUN addgroup -S chaos-operator \
    && adduser -S -u 1001 -g chaos-operator chaos-operator

# Packaging stage
FROM scratch

LABEL maintainer="LitmusChaos"

ENV OPERATOR=/usr/local/bin/chaos-operator

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /output/chaos-operator ${OPERATOR}
COPY --from=builder /etc/passwd /etc/passwd

ENTRYPOINT ["/usr/local/bin/chaos-operator"]

USER 1001
