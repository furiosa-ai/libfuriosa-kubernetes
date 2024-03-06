SHELL := /bin/bash

# make assumption that hwloc is installed with brew command "brew install hwloc"
ifeq ($(shell uname -s),Darwin)
    CGO_CFLAGS := -I/opt/homebrew/opt/hwloc/include
    CGO_LDFLAGS := -L/opt/homebrew/opt/hwloc/lib
endif

define install_deps_function
    @UNAME_S=$$(uname -s); \
    if [ "$$UNAME_S" = "Linux" ]; then \
        echo "Installing for Ubuntu/Debian familly"; \
        sudo apt-get install hwloc libhwloc-dev; \
    elif [ "$$UNAME_S" = "Darwin" ]; then \
        echo "macOS detected. Installing using Homebrew..."; \
        brew install hwloc; \
    else \
        echo "Unsupported Operating System"; \
        exit 1; \
    fi
endef

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

.PHONY: install-deps
install-deps:
	$(call install_deps_function)

.PHONY: example
example:
	$(call build_examples_function,./example)

.PHONY: build-base
build-base:
	docker build . -t ghcr.io/furiosa-ai/libfuriosa-kubernetes:devel --progress=plain --platform=linux/amd64

.PHONY: build-base-no-cache
build-base-no-cache:
	docker build . --no-cache -t ghcr.io/furiosa-ai/libfuriosa-kubernetes:devel --progress=plain --platform=linux/amd64
