FROM golang:1.19-alpine3.17 as build

WORKDIR /src/app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build

FROM alpine:edge

RUN apk add --no-cache ca-certificates

COPY --from=build /src/app/interchain-queries /usr/local/bin/interchain-queries

RUN adduser -S -h /icq -D icq -u 1000

USER icq
