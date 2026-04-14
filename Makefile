SHELL := /bin/sh

GO       ?= go
GOFLAGS  ?=
BIN      ?= bin/whatsmeow-adapter
COMPOSE  ?= docker compose
COMPOSE_FILE ?= docker-compose.yml

.PHONY: all check test build tidy vet fmt docker-build deploy deploy-up deploy-down deploy-logs deploy-ps deploy-restart deploy-smoke

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

deploy: deploy-up

deploy-up:
	$(COMPOSE) -f $(COMPOSE_FILE) up -d --build --remove-orphans

deploy-down:
	COMPOSE_PROFILES=multi-worker $(COMPOSE) -f $(COMPOSE_FILE) down --remove-orphans

deploy-logs:
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f --tail=200

deploy-ps:
	$(COMPOSE) -f $(COMPOSE_FILE) ps

deploy-restart:
	$(COMPOSE) -f $(COMPOSE_FILE) up -d --build --force-recreate --remove-orphans

deploy-smoke:
	@echo "Waiting provider to be ready on http://localhost:8080/healthz ..."
	@i=0; \
	until curl -fsS http://localhost:8080/healthz >/dev/null 2>&1; do \
	  i=$$((i+1)); \
	  if [ $$i -ge 30 ]; then \
	    echo "Provider did not become ready in time"; \
	    exit 1; \
	  fi; \
	  sleep 2; \
	done; \
	echo "Provider is healthy"; \
	curl -fsS http://localhost:8080/api/v1/worker; \
	echo

# Convenience target to run tests/build without pkg-config (CGO still on)
test-nopkg:
	$(GO) test -tags nopkgconfig ./...

build-nopkg:
	$(GO) build -tags nopkgconfig -o $(BIN) ./cmd/whatsmeow-adapter
