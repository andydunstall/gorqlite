.PHONY: all
all: build

.PHONY: build
build:
	go build ./...

.PHONY: deps
deps:
	go get -t -d ./...
	go install github.com/golang/mock/mockgen@v1.6.0

# Runs all unit tests.
.PHONY: test
test:
	go test ./...

# Runs all system tests (and unit tests).
.PHONY: system-test
system-test:
	DEBUG=true go test ./... -v -tags system

# Runs a docker container setup to run system tests.
.PHONY: env
env: image
	docker run -it --rm --volume=$(shell pwd):/app -p 6060:6060 --workdir=/app --name gorqlite gorqlite /bin/bash

# Creates the environment test image.
.PHONY: image
image:
	docker build . -t gorqlite

# Generates all mocks used for unit tests.
.PHONY: generate
generate:
	go generate ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run --disable-all -E errcheck,gosimple,govet,ineffassign,staticcheck,structcheck,typecheck,varcheck,asciicheck,bidichk,bodyclose,gocritic,godox,gosec,revive,stylecheck,unparam,wrapcheck --skip-dirs tests
	# Lint tests with less strict rules.
	golangci-lint run --disable-all -E errcheck,gosimple,govet,ineffassign,staticcheck,structcheck,typecheck,varcheck,asciicheck,bidichk,bodyclose,gocritic,godox,revive,stylecheck,unparam,wrapcheck tests/...
