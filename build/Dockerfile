# Building chaos-operator binary
FROM golang:1.14 AS builder

LABEL maintainer="LitmusChaos"

ADD . /chaos-operator
WORKDIR /chaos-operator

RUN CGO_ENABLED=0 go build -o /build/_output/bin/chaos-operator cmd/manager/main.go

# install operator binary
FROM registry.access.redhat.com/ubi7-dev-preview/ubi-minimal:7.6
ENV OPERATOR=/usr/local/bin/chaos-operator \
    USER_UID=1001 \
    USER_NAME=chaos-operator

COPY --from=builder /build/_output/bin/chaos-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
