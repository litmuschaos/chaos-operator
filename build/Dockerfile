FROM golang:1.14.2

ARG TARGETARCH

ENV OPERATOR=/usr/local/bin/chaos-operator \
    USER_UID=1001 \
    USER_NAME=chaos-operator

# install operator binary
COPY build/_output/bin/chaos-operator${TARGETARCH} ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
