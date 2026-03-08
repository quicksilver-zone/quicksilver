FROM golang:1.24.13-alpine3.23 as build
ARG VERSION
ARG COMMIT

RUN apk add --no-cache gcc musl-dev make

COPY . /app

WORKDIR /app

RUN make build

FROM alpine:3.23

COPY --from=build /app/xcclookup /usr/local/bin/xcclookup

ENTRYPOINT ["xcclookup"]
