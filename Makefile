version ?= latest

build-local:
	go build -ldflags "-linkmode external -extldflags -static" -a xcc.go
	docker build -f Dockerfile.local . -t quicksilverzone/xcclookup:$(version)

docker-release: build-local
	docker push quicksilverzone/xcclookup:$(version)
