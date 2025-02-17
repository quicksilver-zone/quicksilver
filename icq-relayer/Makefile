.PHONY: build

QSVERSION ?= $(shell grep github.com/quicksilver-zone/quicksilver go.mod | grep -Eo 'v[0-9]\.[0-9]\.[0-9].*')
COMMIT ?= $(shell git log -1 --format='%H')
VERSION ?= $(shell git describe --tags --match 'icq-relayer/*' | head -n1 | sed 's/.*\///')

ldflags = -X github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner.VERSION=icq-relayer/$(VERSION) \
          -X github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner.COMMIT=$(COMMIT) \
          -X github.com/quicksilver-zone/quicksilver/icq-relayer/pkg/runner.QUICKSILVER_VERSION=$(QSVERSION)

build:
	go build -ldflags '$(ldflags)'

docker:
	docker buildx build --platform linux/amd64 -f Dockerfile . -t quicksilverzone/interchain-queries:$(VERSION) --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT)

docker-push:
	docker push quicksilverzone/interchain-queries:$(VERSION)
