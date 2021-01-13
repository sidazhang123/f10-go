
GOPATH:=$(shell go env GOPATH)


.PHONY: build
build:

	 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o scheduler-srv main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t scheduler-srv:latest
