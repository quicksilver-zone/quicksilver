FROM golang:1.20-alpine3.18 AS builder
RUN apk add --no-cache git musl-dev openssl-dev linux-headers ca-certificates build-base

WORKDIR /src/app/

COPY go.mod go.sum* ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Cosmwasm - download correct libwasmvm version
RUN WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm | cut -d ' ' -f 2) && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm_muslc.$(uname -m).a \
    -O /lib/libwasmvm_muslc.a && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm_muslc.a | grep $(cat /tmp/checksums.txt | grep $(uname -m) | cut -d ' ' -f 1)

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    LINK_STATICALLY=true make build

# Add to a distroless container
FROM alpine:3.18
COPY --from=builder /src/app/build/quicksilverd /usr/local/bin/quicksilverd
RUN adduser -S -h /quicksilver -D quicksilver -u 1000
USER quicksilver

CMD ["quicksilverd", "start"]
