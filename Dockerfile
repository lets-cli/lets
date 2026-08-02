FROM golang:1.26.5-bookworm AS builder

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
    cd bats-core && \
    ./install.sh /usr && \
    echo Bats installed

RUN go install gotest.tools/gotestsum@v1.13.0

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM golangci/golangci-lint:v2.12.2-alpine AS linter

RUN mkdir -p /.cache && chmod -R 777 /.cache
