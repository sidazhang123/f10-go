GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:

	 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o config-grpc main.go

.PHONY: docker
docker:
	docker build . -t config-grpc:latest
