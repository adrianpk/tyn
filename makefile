APP_NAME = tyn

.PHONY: all build run run-race test test-race lint format check install-hooks

all: build

build:
	go build -o bin/tyn .

run: build
	./bin/$(APP_NAME)

run-race:
	go run -race main.go

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	golangci-lint run ./...

format:
	go fmt ./...

check: format lint

install-hooks:
	ln -sf scripts/hooks/pre-commit .git/hooks/pre-commit
	chmod +x scripts/hooks/pre-commit