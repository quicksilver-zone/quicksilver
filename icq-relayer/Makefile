.PHONY: build

QSVERSION := $(shell grep github.com/quicksilver-zone/quicksilver go.mod | grep -Eo 'v[0-9]\.[0-9]\.[0-9].*')
COMMIT := $(shell git log -1 --format='%H')
VERSION=$(shell git describe --tags | grep 'icq-relayer' | head -n1 | sed 's/.*\///')

ldflags = -X github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner.VERSION=icq-relayer/$(VERSION) \
          -X github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner.COMMIT=$(COMMIT) \
          -X github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner.QUICKSILVER_VERSION=$(QSVERSION)

build:
	go build -ldflags '$(ldflags)' -o ../build

docker:
	docker build -f Dockerfile . -t quicksilverzone/interchain-queries:$(VERSION)

dockerx:
	docker buildx build --platform linux/amd64 -f Dockerfile . -t quicksilverzone/interchain-queries:$(VERSION)

docker-push:
	docker push quicksilverzone/interchain-queries:$(VERSION)
