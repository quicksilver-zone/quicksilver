FROM golang:1.17 AS builder
RUN apt update && apt install liblz4-dev libsnappy-dev libzstd-dev libbz2-dev zlib1g-dev git build-essential -y
RUN git clone https://github.com/facebook/rocksdb.git /opt/rocksdb --branch v6.27.3
WORKDIR /opt/rocksdb
RUN make shared_lib install

WORKDIR /src/app/
COPY go.mod go.sum* ./
RUN go mod download
COPY . .

RUN bash -c 'CGO_CFLAGS="-I/opt/rocksdb/include"; CGO_LDFLAGS="-lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd"; make build'

# Add to a distroless container
FROM debian:bullseye
RUN apt update && apt install liblz4-dev libsnappy-dev libzstd-dev libbz2-dev zlib1g-dev -y
COPY --from=builder /src/app/build/quicksilverd /usr/local/bin/quicksilverd
COPY --from=builder /opt/rocksdb/librocksdb.so.6.27 /usr/lib/
RUN adduser --system --home /quicksilver --disabled-password --disabled-login quicksilver -u 1000
USER quicksilver
RUN mkdir -p /quicksilver/.quicksilverd/data/snapshots/
CMD ["quicksilverd", "start"]
