.PHONY: dev build test clean help

help:
	@echo "Available commands:"
	@echo "  make dev        - Run server in development mode"
	@echo "  make build      - Build the binary"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts and database"

dev:
	@echo "Starting development server..."
	@go run ./cmd/server

build:
	@echo "Building binary..."
	@go build -o ./bin/server ./cmd/server

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning up..."
	@rm -rf ./bin ./data
	@go clean

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
