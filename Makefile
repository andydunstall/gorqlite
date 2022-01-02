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

.PHONY: cover
cover:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out

# Runs all system tests (and unit tests).
.PHONY: system-test
system-test:
	DEBUG=true go test ./... -v -tags system

.PHONY: system-cover
system-cover:
	go test -coverprofile=cover.out ./... -tags system
	go tool cover -html=cover.out


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
	golangci-lint run --disable-all -E errcheck,gosimple,govet,ineffassign,staticcheck,structcheck,typecheck,varcheck,asciicheck,bidichk,bodyclose,gocritic,godox,gosec,revive,stylecheck,unparam,wrapcheck --skip-dirs tests --skip-files _test.go
	# Lint tests with less strict rules.
	golangci-lint run --disable-all -E errcheck,gosimple,govet,ineffassign,staticcheck,structcheck,typecheck,varcheck,asciicheck,bidichk,gocritic,godox,revive,stylecheck,unparam,wrapcheck tests/...
