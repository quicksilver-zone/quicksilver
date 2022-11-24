version ?= latest

git_commit = $(shell git rev-parse --short HEAD)
build_version = $(shell cat VERSION)
linker_flags = "-s -X main.GitCommit=${git_commit} -X main.Version=${build_version} -linkmode external -extldflags -static"

build-local:
	go build -ldflags=${linker_flags} -a xcc.go

build-docker:
	go build -ldflags=${linker_flags} -a xcc.go
	docker build -f Dockerfile.local . -t quicksilverzone/xcclookup:$(version)

docker-release: build-local
	docker push quicksilverzone/xcclookup:$(version)
