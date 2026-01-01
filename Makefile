APP_NAME = cnapi
GO_CMD = go
GO_BUILD = $(GO_CMD) build
OUTPUT_DIR = bin

.PHONY: dev-run dev build clean test test-sources test-source test-store help

help:
	@echo "Ceylon News Backend - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev-run         - Run the application directly"
	@echo "  make dev             - Run with Air hot-reload"
	@echo "  make build           - Build production binary"
	@echo "  make clean           - Remove build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run all tests"
	@echo "  make test-sources    - Run all scraper source tests"
	@echo "  make test-source s=X - Run specific source test (e.g., s=bbc)"
	@echo "  make test-store      - Run database store tests"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up      - Run database migrations"
	@echo "  make migrate-down    - Rollback database migrations"

dev-run:
	@echo "Running development server..."
	$(GO_CMD) run --tags "fts5" cmd/server/main.go

dev:
	@echo "Starting development mode with Air..."
	@ENV=development air

build:
	@echo "Building Go binary for production..."
	mkdir -p $(OUTPUT_DIR)
	$(GO_BUILD) --tags "fts5" -o $(OUTPUT_DIR)/$(APP_NAME) ./cmd/server

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(OUTPUT_DIR)

test:
	@echo "Running all tests..."
	$(GO_CMD) test -v --tags "fts5" ./...

migrate-up:
	@echo "Running database migrations up..."
	goose -dir internal/database/migrations sqlite3 data/db.sqlite up

migrate-down:
	@echo "Running database migrations down..."
	goose -dir internal/database/migrations sqlite3 data/db.sqlite down

test-sources:
	@echo "Running scraper source tests individually..."
	@for file in internal/scraper/source/*_test.go; do \
		# Extract the clean filename (e.g., 'bbc' from 'bbc_test.go') \
		name=$$(basename $$file _test.go); \
		echo ""; \
		echo "=================================================="; \
		echo "Testing source: $$name"; \
		echo "=================================================="; \
		# Run tests in the package that match the filename (case-insensitive) \
		$(GO_CMD) test -v -count=1 -timeout 240s ./internal/scraper/source -run "(?i)$$name" || exit 1; \
	done

test-source:
	@# Check if the 's' variable was provided
	@if [ -z "$(s)" ]; then \
		echo "Error: Please specify a source name. Example: make test-source s=bbc"; \
		exit 1; \
	fi
	@echo "Targeting tests matching: $(s)"
	# Run tests in the package matching the input 's' (case-insensitive)
	$(GO_CMD) test -v -count=1 -timeout 240s ./internal/scraper/source -run "(?i)$(s)"

test-store:
	@echo "Running database store tests..."
	$(GO_CMD) test -v -count=1 --tags "fts5" -timeout 30s ./internal/database/store

