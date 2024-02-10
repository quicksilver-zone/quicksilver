FROM golang:1.22.0-alpine3.19 as build

WORKDIR /src/app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=build /src/app/icq-relayer /usr/local/bin/icq-relayer

RUN adduser -S -h /icq -D icq -u 1000

USER icq
