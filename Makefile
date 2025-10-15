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

# Run all tests
test: test-unit test-logic test-popup test-info-popup test-character-stats test-input

# Run unit tests (if any exist)
test-unit:
	@echo "Running unit tests..."
	go test ./internal/... ./cmd/...

# Run logic verification tests
test-logic:
	@echo "Running popup logic tests..."
	go run test/logic_test/main.go

# Run popup widget functionality tests
test-popup:
	@echo "Running popup widget tests..."
	go run test/popup_test/main.go

# Run info popup widget tests (requires display)
test-info-popup:
	@echo "Running info popup widget tests..."
	@echo "Note: This test requires a display/GUI environment"
	go run test/info_popup_test/main.go

# Run character stats widget tests (requires display)
test-character-stats:
	@echo "Running character stats widget tests..."
	@echo "Note: This test requires a display/GUI environment"
	go run test/character_stats_test/main.go

# Run interactive input tests (requires display)
test-input:
	@echo "Running interactive input tests..."
	@echo "Note: This test requires a display/GUI environment"
	go run test/input_test/main.go

# Run all tests using shell script
test-all:
	@echo "Running all tests via shell script..."
	./scripts/run_tests.sh

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
	@echo "  make test     - Run all tests"
	@echo "  make test-unit    - Run unit tests only"
	@echo "  make test-logic   - Run logic verification tests"
	@echo "  make test-popup   - Run popup widget tests"
	@echo "  make test-info-popup - Run info popup widget tests"
	@echo "  make test-character-stats - Run character stats widget tests"
	@echo "  make test-input   - Run interactive input tests"
	@echo "  make test-all     - Run all tests via shell script"
	@echo "  make release  - Build optimized release version"
	@echo "  make dev      - Build development version with race detection"
	@echo "  make help     - Show this help message"
	@echo ""
	@echo "Manual commands:"
	@echo "  go build -o ./bin/myrpg ./cmd/myrpg  - Correct build command"
	@echo "  go build ./cmd/myrpg                 - WRONG (creates binary in root)"