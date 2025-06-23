.PHONY: all build run run-race test test-race lint format check install-hooks list tasks demo test-notification view-logs clear-logs update update-places update-due update-tags update-text

APP_NAME = tn
TYN_DEV = TYN_DEV=1

all: build

build:
	go build -o bin/tn .

run: build
	$(TYN_DEV) ./bin/$(APP_NAME)

run-race:
	$(TYN_DEV) go run -race main.go

test:
	$(TYN_DEV) go test ./...

test-race:
	$(TYN_DEV) go test -race ./...

lint:
	golangci-lint run ./...

format:
	go fmt ./...

check: format lint

install-hooks:
	ln -sf scripts/hooks/pre-commit .git/hooks/pre-commit
	chmod +x scripts/hooks/pre-commit

list:
	$(TYN_DEV) bin/tn list

tasks:
	$(TYN_DEV) bin/tn tasks

update-tags:
	@if [ -z "$(ID)" ]; then \
		echo "Error: Task ID is required. Usage: make update-tags ID=<task_id>"; \
		exit 1; \
	fi
	$(TYN_DEV) bin/tn tasks update $(ID) --tags "projectz,roadmap"

update-place:
	@if [ -z "$(ID)" ]; then \
		echo "Error: Task ID is required. Usage: make update-place ID=<task_id>"; \
		exit 1; \
	fi
	$(TYN_DEV) bin/tn tasks update $(ID) --places "office,home"

update:
	@if [ -z "$(ID)" ]; then \
		echo "Error: Task ID is required. Usage: make update ID=<task_id>"; \
		exit 1; \
	fi
	$(TYN_DEV) bin/tn tasks update $(ID) --tags "projectx,roadmap"

update-places:
	@if [ -z "$(ID)" ]; then \
		echo "Error: Task ID is required. Usage: make update-places ID=<task_id>"; \
		exit 1; \
	fi
	$(TYN_DEV) bin/tn tasks update $(ID) --places "office,home"

update-due:
	@if [ -z "$(ID)" ]; then \
		echo "Error: Task ID is required. Usage: make update-due ID=<task_id>"; \
		exit 1; \
	fi
	$(TYN_DEV) bin/tn tasks update $(ID) --due 2025-07-01

update-text:
	@if [ -z "$(ID)" ]; then \
		echo "Error: Task ID is required. Usage: make update-text ID=<task_id>"; \
		exit 1; \
	fi
	$(TYN_DEV) bin/tn tasks update $(ID) --text "Updated task text with new details"

demo:
	@echo "Creating a realistic set of notes, tasks, and links..."
	@current_hour=$$(date +"%H"); \
	previous_hour=$$(($$current_hour - 1)); \
	if [ $$previous_hour -lt 0 ]; then previous_hour=0; fi; \
	today_date="2025-06-19-$$previous_hour-00-00"; \
	echo "Using date: $$today_date"; \
	$(TYN_DEV) bin/tn capture "Sync with Alice and Bob #projectx :done ^2025-06-17 Discussed Q3 roadmap https://company.com/roadmap"; \
	$(TYN_DEV) bin/tn capture "Write project summary :todo #writing @home Due by end of week"; \
	$(TYN_DEV) bin/tn capture "Read about Go generics https://go.dev/doc/tutorial/generics #reading #golang ^$$today_date"; \
	$(TYN_DEV) bin/tn capture "Coffee with Carol #networking @cafe Great conversation about potential collaboration ^$$today_date"; \
	$(TYN_DEV) bin/tn capture "Submit tax report :done ^2025-04-15 #finance @office Filed electronically"; \
	$(TYN_DEV) bin/tn capture "Research cloud providers #infrastructure :wip Comparing AWS, GCP and Azure pricing models"; \
	$(TYN_DEV) bin/tn capture "Weekly team meeting notes #internal Team discussed sprint goals and blockers ^$$today_date"; \
	$(TYN_DEV) bin/tn capture "Schedule dentist appointment :todo ^2025-06-20 #health @personal"; \
	$(TYN_DEV) bin/tn capture "Review pull request #23 :todo #coding https://github.com/user/repo/pull/23"; \
	$(TYN_DEV) bin/tn capture "Design database schema :done #projectX @home Finalized user and product tables"; \
	$(TYN_DEV) bin/tn capture "Interesting article on CLI tools https://dev.to/cli-tools #reading Bookmark for weekend ^$$today_date"; \
	$(TYN_DEV) bin/tn capture "Order new laptop :wip #shopping @online Looking at developer-focused models"; \
	$(TYN_DEV) bin/tn capture "Fix critical bug #urgent :wip ^2025-06-10 Need to fix memory leak issue for release"; \
	$(TYN_DEV) bin/tn capture +code-echo "Network Security Alert: Identifying Echo-Pattern Vulnerabilities #security Our team recently discovered a critical vulnerability in proxy services that process high-volume request batching."; \
	$(TYN_DEV) bin/tn capture +code-echo "Key findings: The malformed packet signatures trigger a reflection pattern, duplicating header data across node boundaries. #networking This can be detected by implementing checksum verification at each network hop. @devops"
	@echo "Demo data created successfully. Run 'make list' to see results."

test-notification:
	@echo "Creating task with due date in 2 minutes..."
	@future_date=$$(date -d "+1 minutes" +"%Y-%m-%d-%H-%M-%S"); \
	$(TYN_DEV) bin/tn capture "Test notification task :todo ^$$future_date This task should generate a notification"
	@echo "The notification should appear soon."
	@echo "Press Ctrl+C when finished testing."

logs:
	@echo "Showing daemon logs (press Ctrl+C to exit):"
	@if [ -f ~/.tyn/daemon.log ]; then \
		tail -f ~/.tyn/daemon.log; \
	else \
		echo "No logs found at ~/.tyn/daemon.log"; \
	fi

logs-bg:
	@echo "Starting background log viewer..."
	@if [ -f ~/.tyn/daemon.log ]; then \
		xterm -T "tyn daemon logs" -e "tail -f ~/.tyn/daemon.log" & \
		echo "Log viewer opened in new terminal window."; \
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
