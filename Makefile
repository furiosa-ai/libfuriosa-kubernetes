SHELL := /bin/bash
BASE_IMAGE := registry.corp.furiosa.ai/furiosa/furiosa-smi:latest

ifeq ($(shell uname -s),Darwin)
    CGO_CFLAGS := "-I/usr/local/include"
    CGO_LDFLAGS := "-L/usr/local/lib"
endif

ifeq ($(shell uname), Linux)
export LD_LIBRARY_PATH := $(LD_LIBRARY_PATH):/usr/local/lib
endif

define build_examples_function
    @for dir in $(1)/*; do \
        if [ -d "$$dir" ] && [ -f "$$dir/$$(basename $$dir).go" ]; then \
            CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) go build -o "$$(basename $$dir)" $$dir/$$(basename $$dir).go; \
            echo "Built $$dir"; \
        fi \
    done
endef

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

.PHONY: example
example:
	$(call build_examples_function,./example)

.PHONY: base
base:
	docker build . -t registry.corp.furiosa.ai/furiosa/libfuriosa-kubernetes:devel --progress=plain --platform=linux/amd64 --build-arg BASE_IMAGE=$(BASE_IMAGE)

.PHONY: base-no-cache
base-no-cache:
	docker build . --no-cache -t registry.corp.furiosa.ai/furiosa/libfuriosa-kubernetes:devel --progress=plain --platform=linux/amd64 --build-arg BASE_IMAGE=$(BASE_IMAGE)

.PHONY: base-rel
base-rel:
	docker build . -t registry.corp.furiosa.ai/furiosa/libfuriosa-kubernetes:latest --progress=plain --platform=linux/amd64 --build-arg BASE_IMAGE=$(BASE_IMAGE)

.PHONY: base-no-cache-rel
base-no-cache-rel:
	docker build . --no-cache -t registry.corp.furiosa.ai/furiosa/libfuriosa-kubernetes:latest --progress=plain --platform=linux/amd64 --build-arg BASE_IMAGE=$(BASE_IMAGE)

