GO ?= go

.PHONY: build vet test lint coverage integration

build:
	$(GO) build ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test -count=1 ./...

lint:
	golangci-lint run

coverage:
	$(GO) test -count=1 -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

integration:
	docker compose -f docker-compose.integration.yml up -d --wait
	$(GO) test -count=1 -tags=integration ./...
	docker compose -f docker-compose.integration.yml down -v
