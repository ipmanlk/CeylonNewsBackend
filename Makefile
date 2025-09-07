APP_NAME = cnapi
GO_CMD = go
GO_BUILD = $(GO_CMD) build
OUTPUT_DIR = bin

.PHONY: dev-run dev build clean migrate-up migrate-down

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

migrate-up:
	@echo "Running database migrations up..."
	goose -dir internal/database/migrations sqlite3 data/db.sqlite up

migrate-down:
	@echo "Running database migrations down..."
	goose -dir internal/database/migrations sqlite3 data/db.sqlite down