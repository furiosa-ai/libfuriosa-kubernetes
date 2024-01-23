SHELL := /bin/bash

# make assumption that hwloc is installed with brew command "brew install hwloc"
ifeq ($(shell uname -s),Darwin)
    CGO_CFLAGS := -I/opt/homebrew/opt/hwloc/include
    CGO_LDFLAGS := -L/opt/homebrew/opt/hwloc/lib
endif


.PHONY: all
all: build fmt lint vet test tidy vendor

.PHONY: build
build:
	CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) go build ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) golangci-lint run

.PHONY: vet
vet:
	CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) go vet -v ./...

.PHONY: test
test:
	CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) go test ./...

.PHONY: cover
cover:
	CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

.PHONY: tidy
tidy:
	go mod tidy

.PHONY:vendor
vendor:
	go mod vendor