APP_NAME = cnapi
GO_CMD = go
GO_BUILD = $(GO_CMD) build
OUTPUT_DIR = bin
TMP_DIR = .tmp

.PHONY: all dev clean build run dev-prod

all: build

dev:
	@echo "Starting development mode with Air..."
	@ENV=development air

dev-prod:
	@echo "Starting in production mode..."
	@ENV=production make build
	@ENV=production make run

build:
	@echo "Building Go binary..."
	mkdir -p $(OUTPUT_DIR)
	$(GO_BUILD) -o $(OUTPUT_DIR)/$(APP_NAME) ./cmd/server

run:
	@echo "Running production binary..."
	./$(OUTPUT_DIR)/$(APP_NAME)

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(OUTPUT_DIR)
	rm -rf $(TMP_DIR)

migrate-up:
	goose -dir internal/database/migrations sqlite3 data/db.sqlite up

migrate-down:
	goose -dir internal/database/migrations sqlite3 data/db.sqlite down