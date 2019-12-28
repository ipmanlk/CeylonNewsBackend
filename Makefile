GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

BINARY_NAME=cnapi
BUILD_DIR=build

# Main build target
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

# Run tests
test:
	$(GOTEST) -v ./...

# Clean up
clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

# Build and run the application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)


# Docker specific
dkdev:
	docker-compose -f compose.dev.yml up -d

dkdevstop:
	docker-compose -f compose.dev.yml stop

dkdevrm:
	docker-compose -f compose.dev.yml rm --stop --force

dkdevwatch:
	docker-compose -f compose.dev.yml watch

dkdevbuild:
	docker-compose -f compose.dev.yml up --build -d

dkprod:
	docker-compose -f compose.prod.yml up -d

dkprodbuild:
	docker-compose -f compose.prod.yml up --build -d

dkprodstop:
	docker-compose -f compose.prod.yml stop

dkprodrm:
	docker-compose -f compose.prod.yml rm --stop --force
