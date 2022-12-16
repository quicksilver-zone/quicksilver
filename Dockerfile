FROM golang:1.19-alpine3.17 AS builder
RUN apk add --no-cache git musl-dev openssl-dev linux-headers ca-certificates build-base

WORKDIR /src/app/

COPY go.mod go.sum* ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Cosmwasm - download correct libwasmvm version
RUN WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm | cut -d ' ' -f 2) && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/$(WASMVM_VERSION)/libwasmvm_muslc.$(uname -m).a \
      -O /lib/libwasmvm_muslc.a

# Cosmwasm - verify checksum
RUN wget https://github.com/CosmWasm/wasmvm/releases/download/$(WASMVM_VERSION)/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm_muslc.a | grep $(cat /tmp/checksums.txt | grep $(uname -m) | cut -d ' ' -f 1)

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    VERSION=$(echo $(git describe --tags) | sed 's/^v//') && \
    COMMIT=$(git log -1 --format='%H') && \
    go build \
      -mod=readonly \
      -tags "netgo,ledger,muslc" \
      -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name="quicksilver" \
              -X github.com/cosmos/cosmos-sdk/version.AppName="quicksilverd" \
              -X github.com/cosmos/cosmos-sdk/version.Version=$VERSION \
              -X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT \
              -X github.com/cosmos/cosmos-sdk/version.BuildTags='netgo,ledger,muslc' \
              -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -static'" \
      -trimpath \
      -o /src/app/build/ \
      ./...

# Add to a distroless container
FROM alpine:3.17
COPY --from=builder /src/app/build/quicksilverd /usr/local/bin/quicksilverd
RUN adduser -S -h /quicksilver -D quicksilver -u 1000
USER quicksilver
CMD ["quicksilverd", "start"]
