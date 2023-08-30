.PHONY: dev run test cover fmt pre-commit docker-build docker-run

.DEFAULT: help
help:
	@echo "make dev"
	@echo "	setup development environment"
	@echo "make run"
	@echo "	run app"
	@echo "make test"
	@echo "	run go test"
	@echo "make cover"
	@echo "	run go test with -cover"
	@echo "make tidy"
	@echo "	run go mod tidy"
	@echo "make fmt"
	@echo "	run gofumpt"
	@echo "make pre-commit"
	@echo "	run pre-commit hooks"
	@echo "make docker-build"
	@echo "	build docker image"
	@echo "make docker-run"
	@echo "	run docker image"

check-gofumpt:
ifeq (, $(shell which gofumpt))
	$(error "gofumpt not in $(PATH), gofumpt (https://pkg.go.dev/mvdan.cc/gofumpt) is required")
endif

check-pre-commit:
ifeq (, $(shell which pre-commit))
	$(error "pre-commit not in $(PATH), pre-commit (https://pre-commit.com) is required")
endif

dev: check-pre-commit
	pre-commit install

run:
	go build -o tinysyslog-bin ./cmd/tinysyslogd && ./tinysyslog-bin

build:
	go build -o tinysyslog-bin ./cmd/tinysyslogd

test:
	go test -v ./...

cover:
	go test -cover -v ./...

tidy:
	go mod tidy

fmt: check-gofumpt
	gofumpt -l -w .

pre-commit: check-pre-commit
	pre-commit

docker-build:
	docker build -t tinysyslog .

docker-run:
	docker run -p 5140:5140/udp --rm tinysyslog
