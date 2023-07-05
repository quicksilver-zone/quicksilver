FROM golang:1.20-alpine3.17 as build
ARG lfs

RUN apk add --no-cache gcc musl-dev

COPY . /app

WORKDIR /app

RUN go build -ldflags="$lfs" -a xcc.go

FROM alpine:3.17

COPY --from=build /app/xcclookup /usr/local/bin/xcc

ENTRYPOINT xcc
