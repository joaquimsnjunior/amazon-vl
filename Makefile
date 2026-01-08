.PHONY: build run test clean lint docker-build docker-run help

# Variables
BINARY_NAME=amazon-vl
VERSION?=1.1.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"

# Default target
.DEFAULT_GOAL := help

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd

## build-static: Build static binary for containers
build-static:
	@echo "Building static $(BINARY_NAME)..."
	CGO_ENABLED=0 GOOS=linux go build -a $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd

## run: Run the application (usage: make run DIR=/path/to/logs PORT=9000)
run:
	@if [ -z "$(DIR)" ] || [ -z "$(PORT)" ]; then \
		echo "Usage: make run DIR=/path/to/logs PORT=9000"; \
		exit 1; \
	fi
	go run ./cmd $(DIR) $(PORT)

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## tidy: Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) -t $(BINARY_NAME):latest .

## docker-run: Run Docker container (usage: make docker-run DIR=/path/to/logs PORT=9000)
docker-run:
	@if [ -z "$(DIR)" ] || [ -z "$(PORT)" ]; then \
		echo "Usage: make docker-run DIR=/path/to/logs PORT=9000"; \
		exit 1; \
	fi
	docker run -d -p $(PORT):$(PORT) -v $(DIR):/logs:ro $(BINARY_NAME):latest /logs $(PORT)

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

## install: Install binary to /usr/local/bin
install: build
	@echo "Installing to /usr/local/bin..."
	sudo cp bin/$(BINARY_NAME) /usr/local/bin/

## help: Show this help
help:
	@echo "Amazon-VL Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'