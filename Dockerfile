FROM golang:1.18.3-bullseye as builder

ENV GOPROXY https://proxy.golang.org
WORKDIR /app

RUN apt-get update && apt-get install -y \
    git gcc \
    zsh  # for zsh completion tests

RUN cd /tmp && \
    git clone https://github.com/bats-core/bats-core && \
    git clone https://github.com/bats-core/bats-support.git /bats/bats-support && \
    git clone https://github.com/bats-core/bats-assert.git /bats/bats-assert && \
    cd bats-core && \
    ./install.sh /usr && \
    echo Bats installed

RUN go install gotest.tools/gotestsum@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM golangci/golangci-lint:v1.50-alpine as linter

RUN mkdir -p /.cache && chmod -R 777 /.cache
