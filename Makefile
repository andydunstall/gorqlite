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

.PHONY: generate
generate:
	go generate ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run --disable-all -E deadcode,errcheck,gosimple,govet,ineffassign,staticcheck,structcheck,typecheck,unused,varcheck,asciicheck,bidichk,bodyclose,gocritic,godox,gosec,revive,stylecheck,unparam,wrapcheck --skip-dirs tests
	# Lint tests with less strict rules.
	golangci-lint run --disable-all -E deadcode,errcheck,gosimple,govet,ineffassign,staticcheck,structcheck,typecheck,unused,varcheck,asciicheck,bidichk,bodyclose,gocritic,godox,revive,stylecheck,unparam,wrapcheck tests/...
