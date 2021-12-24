FROM golang:1.17.5-buster

ENV GOPROXY https://proxy.golang.org
WORKDIR /app

RUN apt-get update && apt-get install git gcc

RUN cd /tmp && \
    git clone https://github.com/sstephenson/bats.git && \
    cd bats && \
    ./install.sh /usr

RUN go install gotest.tools/gotestsum@latest

COPY go.mod .
COPY go.sum .

RUN go mod download
