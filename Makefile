# Makefile for MyRPG
# Simple build automation for the tactical RPG game

.PHONY: build run clean test help

# Default target
all: build

# Build the game binary to the bin directory
build:
	@echo "Building MyRPG..."
	@mkdir -p bin
	go build -o ./bin/myrpg ./cmd/myrpg
	@echo "✅ Binary created at: ./bin/myrpg"

# Run the game
run: build
	@echo "Running MyRPG..."
	./bin/myrpg

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f ./bin/myrpg
	rm -f ./myrpg
	@echo "✅ Clean complete"

# Run tests (when we have them)
test:
	@echo "Running tests..."
	go test ./...

# Build for release (with optimizations)
release:
	@echo "Building release version..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/myrpg ./cmd/myrpg
	@echo "✅ Release binary created at: ./bin/myrpg"

# Development build with race detection
dev:
	@echo "Building development version with race detection..."
	@mkdir -p bin
	go build -race -o ./bin/myrpg ./cmd/myrpg
	@echo "✅ Development binary created at: ./bin/myrpg"

# Show help
help:
	@echo "MyRPG Build Commands:"
	@echo "  make build    - Build the game binary to ./bin/myrpg"
	@echo "  make run      - Build and run the game"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make release  - Build optimized release version"
	@echo "  make dev      - Build development version with race detection"
	@echo "  make help     - Show this help message"
	@echo ""
	@echo "Manual commands:"
	@echo "  go build -o ./bin/myrpg ./cmd/myrpg  - Correct build command"
	@echo "  go build ./cmd/myrpg                 - WRONG (creates binary in root)"