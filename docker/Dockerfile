FROM golang:1.13.7

ENV GOPROXY https://proxy.golang.org
WORKDIR /app

RUN apt-get update && apt-get install git gcc

# install dev tools, linters, etc.
RUN go get -u honnef.co/go/tools/cmd/staticcheck

COPY go.mod .
COPY go.sum .

RUN go get ./...
