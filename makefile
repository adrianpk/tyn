.PHONY: all build run run-race test test-race lint format check install-hooks list demo test-notification view-logs clear-logs

APP_NAME = tn

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

list:
	bin/tn list

demo:
	bin/tn capture "This is a simple note."
	bin/tn capture "Note with #tag1 #tag2 @home @office"
	bin/tn capture "A task :todo"
	bin/tn capture "A task with override date :done ^2025-06-20"
	bin/tn capture "A useful link https://example.com"

test-notification:
	@echo "Creating task with due date in 2 minutes..."
	@future_date=$$(date -d "+1 minutes" +"%Y-%m-%d-%H-%M-%S"); \
	bin/tn capture "Test notification task :todo ^$$future_date This task should generate a notification"
	@echo "The notification should appear soon."
	@echo "Press Ctrl+C when finished testing."

view-logs:
	@echo "Showing daemon logs (press Ctrl+C to exit):"
	@if [ -f ~/.tyn/daemon.log ]; then \
		tail -f ~/.tyn/daemon.log; \
	else \
		echo "No logs found at ~/.tyn/daemon.log"; \
	fi

clear-logs:
	@echo "Clearing daemon logs..."
	@if [ -f ~/.tyn/daemon.log ]; then \
		> ~/.tyn/daemon.log && echo "Logs cleared successfully."; \
	else \
		echo "No logs found at ~/.tyn/daemon.log"; \
	fi
