FROM golang:1.24.5-alpine3.21 AS builder
RUN apk add --no-cache git musl-dev openssl-dev linux-headers ca-certificates build-base

WORKDIR /src/app/

COPY go.mod go.sum* ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    LINK_STATICALLY=true make build

# Add to a distroless container
FROM alpine:3.20
COPY --from=builder /src/app/build/quicksilverd /usr/local/bin/quicksilverd
RUN adduser -S -h /quicksilver -D quicksilver -u 1000
USER quicksilver

CMD ["quicksilverd", "start"]
