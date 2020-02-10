FROM golang:1.13.7

WORKDIR /app

RUN apt-get update && apt-get install git gcc
RUN go get -u honnef.co/go/tools/cmd/staticcheck
