APP_NAME = tn

.PHONY: all build run run-race test test-race lint format check install-hooks demo

all: build

build:
	go build -o bin/tn .

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

demo:
	bin/tn capture "This is a simple note."
	bin/tn capture "Note with #tag1 #tag2 @home @office"
	bin/tn capture "A task :todo"
	bin/tn capture "A task with override date :done ^2025-06-20"
	bin/tn capture "A useful link https://example.com"
