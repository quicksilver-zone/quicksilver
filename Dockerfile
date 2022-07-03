FROM golang:1.18-alpine3.15 AS builder
RUN apk add --no-cache make git gcc musl-dev openssl-dev linux-headers 

WORKDIR /src/app/

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN make build

# Add to a distroless container
FROM alpine:3.15
COPY --from=builder /src/app/build/quicksilverd /usr/local/bin/quicksilverd
RUN adduser -S -h /quicksilver -D quicksilver -u 1000
USER quicksilver
CMD ["quicksilverd", "start"]
