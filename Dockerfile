# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY ./ ./

RUN apk --update add postgresql-client
RUN apk add --no-cache bash

RUN go mod download
RUN go build -o /docker-order-service ./cmd/main.go

RUN chmod +x entrypoint.sh

CMD ./entrypoint.sh
