FROM golang:1.21.5-alpine3.18 as build
ARG lfs

RUN apk add --no-cache gcc musl-dev

COPY . /app

WORKDIR /app

RUN go build -ldflags="$lfs" -a xcc.go

FROM alpine:3.18

COPY --from=build /app/xcc /usr/local/bin/xcc

ENTRYPOINT xcc
