# MyRPG Build & Test Scripts

This directory contains automation scripts for building and testing the MyRPG project.

## Scripts

### `run_tests.sh`
Comprehensive test runner that executes all test suites and provides detailed reporting.

**Features:**
- ✅ Color-coded output with detailed progress tracking
- ✅ Environment validation (checks for required files)
- ✅ Display detection for interactive tests
- ✅ Error handling and test result summary
- ✅ Build verification

**Usage:**
```bash
# Run all tests via shell script
./scripts/run_tests.sh

# Or via Makefile
make test-all
```

## Test Phases

### 1. **Unit Tests** 🧪
- Runs standard Go unit tests on `internal/` and `cmd/` packages
- Command: `go test ./internal/... ./cmd/...`

### 2. **Logic Verification** 🧠
- Tests popup input blocking logic without graphics
- Command: `go run test/logic_test/main.go`

### 3. **Widget Functionality** 🎯
- Comprehensive PopupSelectionWidget testing
- Command: `go run test/popup_test/main.go`

### 4. **Interactive Tests** 🎮
- GUI-based testing (requires display)
- Command: `go run test/input_test/main.go`
- Auto-terminated after 10 seconds

### 5. **Build Verification** 🔨
- Ensures main game builds successfully
- Command: `go build -o /tmp/myrpg_test ./cmd/myrpg`

## Makefile Integration

The test system is fully integrated with the project Makefile:

```bash
# Run all tests (includes shell script)
make test

# Run individual test phases
make test-unit     # Unit tests only
make test-logic    # Logic verification
make test-popup    # Widget functionality
make test-input    # Interactive tests (GUI required)
make test-all      # Comprehensive shell script

# Other build targets
make build         # Build game binary
make run           # Build and run game
make clean         # Remove build artifacts
make help          # Show all available commands
```

## Exit Codes

The test script uses proper exit codes for CI/CD integration:
- **0**: All tests passed
- **>0**: Number of failed tests (for detailed error reporting)

## Environment Requirements

### All Tests
- Go 1.19+ with proper module support
- Project must be run from root directory

### Interactive Tests
- **macOS**: Native display support
- **Linux**: X11 (`$DISPLAY`) or Wayland (`$WAYLAND_DISPLAY`)
- **Headless**: Interactive tests automatically skipped

## Output Examples

### Successful Run
```
🚀 MyRPG Test Suite Runner
==================================================
✅ Test environment verified

🧪 Phase 1: Unit Tests
✅ PASSED: Go Unit Tests

🧠 Phase 2: Logic Verification  
✅ PASSED: Popup Logic Tests

🎉 ALL TESTS PASSED! 🎉
Total Tests: 5, Passed: 5, Failed: 0
```

### Test Failure
```
❌ FAILED: Interactive Input Tests
📊 TEST SUMMARY
Total Tests: 5, Passed: 4, Failed: 1
⚠️  Some tests failed. Please review the output above.
```

## Integration with CI/CD

The test system is designed to work seamlessly in automated environments:

1. **Unit Tests**: Always run (no graphics required)
2. **Logic Tests**: Always run (no graphics required) 
3. **Widget Tests**: Always run (no graphics required)
4. **Interactive Tests**: Automatically skipped in headless environments
5. **Build Tests**: Always run (validates compilation)

This ensures maximum test coverage while maintaining compatibility with headless CI/CD systems.