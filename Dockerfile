FROM golang:1.18-alpine AS build-go

ARG GIT_SSH_KEY
ARG KNOWN_HOSTS_CONTENT
WORKDIR /go/src/github.com/Workiva/go-datastructures/
ADD . /go/src/github.com/Workiva/go-datastructures/

ARG GOPATH=/go/
ENV PATH $GOPATH/bin:$PATH
RUN echo "Starting the script section" && \
    go mod vendor && \
    echo "script section completed"

ARG BUILD_ARTIFACTS_DEPENDENCIES=/go/src/github.com/Workiva/go-datastructures/go.mod

FROM scratch
