.PHONY: all
all: build

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./...

.PHONY: system-test
system-test:
	DEBUG=true go test ./... -v -tags system

.PHONY: env
env: image
	docker run -it --rm --volume=$(shell pwd):/app --workdir=/app --name gorqlite gorqlite /bin/bash

.PHONY: image
image:
	docker build . -t gorqlite

.PHONY: fmt
fmt:
	go fmt ./...
