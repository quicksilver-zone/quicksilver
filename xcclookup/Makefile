.PHONY: build

QSVERSION ?= $(shell grep github.com/quicksilver-zone/quicksilver go.mod | grep -Eo 'v[0-9]\.[0-9]\.[0-9].*')
COMMIT ?= $(shell git log -1 --format='%H')
VERSION ?= $(shell git describe --tags --match 'xcclookup/*' | head -n1 | sed 's/.*\///')

ldflags = -X github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types.VERSION=xcclookup/$(VERSION) \
          -X github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types.COMMIT=$(COMMIT) \
          -X github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types.QUICKSILVER_VERSION=$(QSVERSION)

build:
	go build -ldflags '$(ldflags)'

docker:
	docker buildx build --platform linux/amd64 -f Dockerfile . -t quicksilverzone/xcclookup:$(VERSION) --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT)

docker-push:
	docker push quicksilverzone/xcclookup:$(VERSION)


###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint run

lint-fix:
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --fix --issues-exit-code=0

.PHONY: lint lint-fix

format:
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run mvdan.cc/gofumpt -w .
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run github.com/client9/misspell/cmd/misspell -w
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run golang.org/x/tools/cmd/goimports -w -local github.com/quicksilver-zone/quicksilver/xcclookup
.PHONY: format

mdlint:
	@echo "--> Running markdown linter"
	@$(DOCKER) run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md"

mdlint-fix:
	@$(DOCKER)  run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md" --fix

###############################################################################
###                                Protobuf                                 ###
###############################################################################

