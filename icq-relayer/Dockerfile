FROM golang:1.24.5-alpine3.22 as build

ARG VERSION
ARG COMMIT

WORKDIR /src/app

RUN apk add --no-cache gcc musl-dev make git

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN make build

FROM alpine:3.22

RUN apk add --no-cache ca-certificates jq curl

COPY --from=build /src/app/icq-relayer /usr/local/bin/icq-relayer

RUN adduser -S -h /icq -D icq -u 1000

USER icq
