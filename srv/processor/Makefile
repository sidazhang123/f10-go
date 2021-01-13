
GOPATH:=$(shell go env GOPATH)


.PHONY: build
build:

	 GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o processor-srv main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t processor-srv:latest
