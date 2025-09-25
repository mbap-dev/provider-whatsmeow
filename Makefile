SHELL := /bin/sh

GO       ?= go
GOFLAGS  ?=
BIN      ?= bin/whatsmeow-adapter

.PHONY: all check test build tidy vet fmt docker-build

all: check build

check: tidy vet test

tidy:
	$(GO) mod tidy

vet:
	$(GO) vet ./...

test:
	$(GO) test $(GOFLAGS) ./...

build:
	$(GO) build $(GOFLAGS) -o $(BIN) ./cmd/whatsmeow-adapter

docker-build:
	docker build -t provider-whatsmeow:local .

# Convenience target to run tests/build without pkg-config (CGO still on)
test-nopkg:
	$(GO) test -tags nopkgconfig ./...

build-nopkg:
	$(GO) build -tags nopkgconfig -o $(BIN) ./cmd/whatsmeow-adapter

