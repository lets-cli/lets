FROM golang:1.16.0-buster

ENV GOPROXY https://proxy.golang.org
WORKDIR /app

RUN apt-get update && apt-get install git gcc

RUN cd /tmp && \
    git clone https://github.com/sstephenson/bats.git && \
    cd bats && \
    ./install.sh /usr

COPY go.mod .
COPY go.sum .

RUN go mod download
