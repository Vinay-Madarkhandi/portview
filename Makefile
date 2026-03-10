# Makefile for PortView
# Build, test, and release targets.

BINARY   := portview
CMD      := ./cmd/portview
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -s -w -X main.version=$(VERSION)

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: all build clean test vet fmt lint release help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'

## build: Build for current platform
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

## test: Run all tests
test:
	go test ./... -v

## vet: Run go vet
vet:
	go vet ./...

## fmt: Format code
fmt:
	gofmt -s -w .

## lint: Run staticcheck (install: go install honnef.co/go/tools/cmd/staticcheck@latest)
lint:
	staticcheck ./...

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/

## release: Cross-compile for all platforms (output in dist/)
release: clean
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output="dist/$(BINARY)-$${os}-$${arch}"; \
		echo "Building $${output}..."; \
		GOOS=$${os} GOARCH=$${arch} go build -ldflags "$(LDFLAGS)" -o $${output} $(CMD); \
	done
	@echo ""
	@echo "Release binaries:"
	@ls -lh dist/

all: build

