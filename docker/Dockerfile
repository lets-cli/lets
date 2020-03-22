FROM golang:1.13.7

ENV GOPROXY https://proxy.golang.org
WORKDIR /app

RUN apt-get update && apt-get install git gcc

COPY go.mod .
COPY go.sum .

RUN cd /tmp && \
    git clone https://github.com/sstephenson/bats.git && \
    cd bats && \
    ./install.sh /usr

RUN go get ./...
