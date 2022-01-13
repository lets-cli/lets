FROM golang:1.17.5-buster

ENV GOPROXY https://proxy.golang.org
WORKDIR /app

RUN apt-get update && apt-get install git gcc

RUN cd /tmp && \
    git clone git://github.com/bats-core/bats-core && \
    git clone git://github.com/bats-core/bats-support.git /bats/bats-support && \
    git clone git://github.com/bats-core/bats-assert.git /bats/bats-assert && \
    cd bats-core && \
    ./install.sh /usr && \
    echo Bats installed

RUN go install gotest.tools/gotestsum@latest

COPY go.mod .
COPY go.sum .

RUN go mod download
