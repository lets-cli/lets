FROM golangci/golangci-lint:v1.45-alpine

RUN mkdir -p /.cache && chmod -R 777 /.cache
