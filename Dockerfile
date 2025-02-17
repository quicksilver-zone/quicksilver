FROM golang:1.23-alpine3.20 as build
ARG lfs

RUN apk add --no-cache gcc musl-dev

COPY . /app

WORKDIR /app

RUN go build -ldflags="$lfs" -a xcc.go

FROM alpine:3.20

COPY --from=build /app/xcc /usr/local/bin/xcc

ENTRYPOINT xcc
