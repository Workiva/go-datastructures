FROM ubuntu:bionic as base

# install go
ENV GO_VERSION='1.15.5'
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin:/usr/local/go/bin
RUN curl -o /tmp/go${GO_VERSION}.linux-amd64.tar.gz https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local/ -xzf /tmp/go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -Rf /tmp/go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -r /usr/local/go/doc /usr/local/go/api && \
    cd ~/

ARG GIT_SSH_KEY
ARG KNOWN_HOSTS_CONTENT
WORKDIR /go/src/github.com/Workiva/go-datastructures/
ADD . /go/src/github.com/Workiva/go-datastructures/

RUN mkdir /root/.ssh && \
    echo "$KNOWN_HOSTS_CONTENT" > "/root/.ssh/known_hosts" && \
    chmod 700 /root/.ssh/ && \
    umask 0077 && echo "$GIT_SSH_KEY" >/root/.ssh/id_rsa && \
    eval "$(ssh-agent -s)" && ssh-add /root/.ssh/id_rsa

ARG GOPATH=/go/
ENV PATH $GOPATH/bin:$PATH
RUN git config --global url.git@github.com:.insteadOf https://github.com
RUN echo "Starting the script section" && \
    go mode vendor && \
    echo "script section completed"

ARG BUILD_ARTIFACTS_DEPENDENCIES=/go/src/github.com/Workiva/go-datastructures/go.mod

FROM scratch
