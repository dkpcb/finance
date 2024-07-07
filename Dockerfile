FROM golang:1.22-alpine as server-build

WORKDIR  /go/src/finatext

COPY . .

RUN apk upgrade --update && \
    apk --no-cache add git