.PHONY: build dev test fmt

.DEFAULT: help
help:
	@echo "make build"
	@echo "       run go build"
	@echo "make dev"
	@echo "       setup development environment"
	@echo "make test"
	@echo "       run go test"
	@echo "make fmt"
	@echo "       run go fmt"

build:
	go build ./cmd/tinysyslogd

dev:
	@type pre-commit > /dev/null || (echo "ERROR: pre-commit (https://pre-commit.com/) is required."; exit 1)
	pre-commit install

test:
	go test ./...

fmt:
	go fmt ./...

