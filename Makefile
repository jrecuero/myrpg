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
	rm -f ./main
	rm -f ./character_stats_test
	rm -f ./component_test
	rm -f ./dialog_test
	rm -f ./equipment_test
	rm -f ./event_persistence_test
	rm -f ./event_test
	rm -f ./event_visual_test
	rm -f ./info_layout_test
	rm -f ./info_popup_test
	rm -f ./input_test
	rm -f ./inventory_test
	rm -f ./item_system_test
	rm -f ./logic_test
	rm -f ./popup_test
	rm -f ./quest_test
	rm -f ./skills_test
	rm -f ./verify_events
	@echo "✅ Clean complete"

# Clean all build artifacts including test binaries
clean-all: clean
	@echo "Cleaning all test binaries..."
	find . -maxdepth 1 -name "*_test" -type f -delete 2>/dev/null || true
	find . -maxdepth 1 -name "test_*" -type f -delete 2>/dev/null || true
	@echo "✅ All clean complete"

# Build test binaries to bin directory (prevent root clutter)
build-tests:
	@echo "Building test binaries to bin directory..."
	@mkdir -p bin
	go build -o ./bin/character_stats_test ./test/character_stats_test
	go build -o ./bin/component_test ./test/component_test
	go build -o ./bin/dialog_test ./test/dialog_test
	go build -o ./bin/equipment_test ./test/equipment_test
	go build -o ./bin/event_persistence_test ./test/event_persistence_test
	go build -o ./bin/event_test ./test/event_test
	go build -o ./bin/event_visual_test ./test/event_visual_test
	go build -o ./bin/info_layout_test ./test/info_layout_test
	go build -o ./bin/info_popup_test ./test/info_popup_test
	go build -o ./bin/input_test ./test/input_test
	go build -o ./bin/inventory_test ./test/inventory_test
	go build -o ./bin/item_system_test ./test/item_system_test
	go build -o ./bin/logic_test ./test/logic_test
	go build -o ./bin/popup_test ./test/popup_test
	go build -o ./bin/quest_test ./test/quest_test
	go build -o ./bin/skills_test ./test/skills_test
	go build -o ./bin/verify_events ./test/verify_events
	@echo "✅ Test binaries created in ./bin/"

# Run all tests
test: test-unit test-logic test-popup test-info-popup test-character-stats test-equipment test-dialog test-input

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

# Run equipment widget tests (requires display)
test-equipment:
	@echo "Running equipment widget tests..."
	@echo "Note: This test requires a display/GUI environment"
	go run test/equipment_test/main.go

# Run dialog widget tests (requires display)
test-dialog:
	@echo "Running dialog widget tests..."
	@echo "Note: This test requires a display/GUI environment"
	go run test/dialog_test/main.go

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
	@echo "  make clean-all - Remove all build artifacts including test binaries"
	@echo "  make build-tests - Build all test binaries to ./bin/"
	@echo "  make test     - Run all tests"
	@echo "  make test-unit    - Run unit tests only"
	@echo "  make test-logic   - Run logic verification tests"
	@echo "  make test-popup   - Run popup widget tests"
	@echo "  make test-info-popup - Run info popup widget tests"
	@echo "  make test-character-stats - Run character stats widget tests"
	@echo "  make test-equipment - Run equipment widget tests"
	@echo "  make test-dialog  - Run dialog widget tests"
	@echo "  make test-input   - Run interactive input tests"
	@echo "  make test-all     - Run all tests via shell script"
	@echo "  make release  - Build optimized release version"
	@echo "  make dev      - Build development version with race detection"
	@echo "  make help     - Show this help message"
	@echo ""
	@echo "Manual commands:"
	@echo "  go build -o ./bin/myrpg ./cmd/myrpg  - Correct build command"
	@echo "  go build ./cmd/myrpg                 - WRONG (creates binary in root)"