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
