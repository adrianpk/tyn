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
	@echo "Creating a realistic set of notes, tasks, and links..."
	bin/tn capture "Sync with Alice and Bob #projectx :done ^2025-06-17 Discussed Q3 roadmap https://company.com/roadmap"
	bin/tn capture "Write project summary :todo #writing @home Due by end of week"
	bin/tn capture "Read about Go generics https://go.dev/doc/tutorial/generics #reading #golang"
	bin/tn capture "Coffee with Carol #networking @cafe Great conversation about potential collaboration"
	bin/tn capture "Submit tax report :done ^2025-04-15 #finance @office Filed electronically"
	bin/tn capture "Research cloud providers #infrastructure :wip Comparing AWS, GCP and Azure pricing models"
	bin/tn capture "Weekly team meeting notes #internal Team discussed sprint goals and blockers"
	bin/tn capture "Schedule dentist appointment :todo ^2025-06-20 #health @personal"
	bin/tn capture "Review pull request #23 :todo #coding https://github.com/user/repo/pull/23"
	bin/tn capture "Design database schema :done #projectX @home Finalized user and product tables"
	bin/tn capture "Interesting article on CLI tools https://dev.to/cli-tools #reading Bookmark for weekend"
	bin/tn capture "Order new laptop :wip #shopping @online Looking at developer-focused models"
	@echo "Demo data created successfully. Run 'bin/tn list' to see results."

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
