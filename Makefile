git_commit = $(shell git rev-parse --short HEAD)
build_version = v$(shell cat VERSION)
linker_flags = "-s -X main.GitCommit=${git_commit} -X main.Version=${build_version} -linkmode external -extldflags -static"

build:
	go build -ldflags=${linker_flags} -a xcc.go

build-docker-local: build
	docker build -f Dockerfile.local . -t quicksilverzone/xcclookup:$(build_version)

build-docker:
	docker build --build-arg lfs=$(linker_flags) . -t quicksilverzone/xcclookup:$(build_version)

xbuild-docker:
	docker buildx build --platform linux/amd64 --build-arg lfs=$(linker_flags) . -t quicksilverzone/xcclookup:$(build_version)

docker-release:
	docker push quicksilverzone/xcclookup:$(build_version)

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --out-format=tab

lint-fix:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix --out-format=tab --issues-exit-code=0

.PHONY: lint lint-fix

format:
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run mvdan.cc/gofumpt -w .
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run github.com/client9/misspell/cmd/misspell -w
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run golang.org/x/tools/cmd/goimports -w -local github.com/ingenuity-build/quicksilver
.PHONY: format

mdlint:
	@echo "--> Running markdown linter"
	@$(DOCKER) run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md"

mdlint-fix:
	@$(DOCKER)  run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md" --fix

###############################################################################
###                                Protobuf                                 ###
###############################################################################

