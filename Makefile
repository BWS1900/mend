SHELL := /bin/bash
BIN   := mend
PKG   := ./...
LDFLAGS := -s -w \
  -X github.com/will/mend/internal/version.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev) \
  -X github.com/will/mend/internal/version.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo none) \
  -X github.com/will/mend/internal/version.Date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: all build test test-race cover fmt vet lint clean run snapshot release install

all: test build

build:
	go build -ldflags '$(LDFLAGS)' -o $(BIN) .

test:
	go test $(PKG)

test-race:
	go test -race $(PKG)

cover:
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out -o coverage.html

fmt:
	gofmt -s -w .

vet:
	go vet $(PKG)

clean:
	rm -f $(BIN) coverage.out coverage.html
	rm -rf dist/

run: build
	./$(BIN) examples/basic.md | head -40

snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean

install: build
	install -m 0755 $(BIN) $(GOPATH)/bin/$(BIN)
