FROM golang:1.18-alpine3.15 as build

COPY . /app
COPY ./quicksilver /quicksilver

WORKDIR /app

RUN go build xcc.go

FROM alpine

COPY --from=build /app/xcclookup /usr/local/bin/xcc

ENTRYPOINT xcc
