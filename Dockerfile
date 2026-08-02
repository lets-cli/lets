FROM golang:1.26.5-bookworm@sha256:1ecb7edf62a0408027bd5729dfd6b1b8766e578e8df93995b225dfd0944eb651 AS builder

ENV GOPROXY=https://proxy.golang.org
ENV CGO_ENABLED=0

WORKDIR /app

RUN apt-get update && apt-get install -y \
    git \
    zsh  # for zsh completion tests

RUN cd /tmp && \
    git clone https://github.com/bats-core/bats-core && \
    git clone https://github.com/bats-core/bats-support.git /bats/bats-support && \
    git clone https://github.com/bats-core/bats-assert.git /bats/bats-assert && \
    git -C bats-core checkout eb7f42f8d608ac693d7a4b67474f6714ea68cfc5 && \
    git -C /bats/bats-support checkout 24a72e14349690bcbf7c151b9d2d1cdd32d36eb1 && \
    git -C /bats/bats-assert checkout f1e9280eaae8f86cbe278a687e6ba755bc802c1a && \
    cd bats-core && \
    ./install.sh /usr && \
    echo Bats installed

RUN go install gotest.tools/gotestsum@v1.13.0

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM golangci/golangci-lint:v2.11.3-alpine@sha256:b1c3de5862ad0a95b4e45a993b0f00415835d687e4f12c845c7493b86c13414e AS linter

RUN mkdir -p /.cache && chmod -R 777 /.cache
