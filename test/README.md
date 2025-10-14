# Test Programs

This directory contains standalone test programs for verifying specific functionality of the MyRPG game components.

## Directory Structure

```
test/
├── README.md           # This file
├── input_test/         # Interactive input blocking test 
├── logic_test/         # Non-interactive logic verification
└── popup_test/         # Popup widget functionality test
```

## Test Programs

### 1. **`popup_test/`** - Popup Widget Logic Test
- **File**: `test/popup_test/main.go`
- **Purpose**: Comprehensive testing of PopupSelectionWidget functionality
- **Type**: Logic verification (no graphics required)
- **Features Tested**:

**Usage:**
```bash
# Full popup widget test
go run test/popup_test/main.go

**What it tests:**
- ✅ Widget creation and initialization
- ✅ Show/hide functionality  
- ✅ Option management and selection
- ✅ Callback system (selection and cancel)
- ✅ Position and size adjustment
- ✅ Current selection retrieval

**Output:** Text-based verification with checkmarks showing all features work correctly.

---

### `logic_test/` - Input Blocking Logic Test  
**Purpose:** Verify that popup widgets properly block game input to prevent conflicts.

**Usage:**
```bash
go run test/logic_test/main.go
```

**What it tests:**
- ✅ Popup visibility detection
- ✅ Game input blocking when popup is active
- ✅ Input restoration when popup is closed
- ✅ Clean separation between UI and game input

**Output:** Logic flow verification showing input blocking mechanism works correctly.

---

### `input_test/` - Interactive Input Test
**Purpose:** Interactive test for verifying input blocking behavior with actual key presses.

**Usage:**
```bash
go run test/input_test/main.go
```

**What it tests:**
- 🎮 Real-time input handling
- 🎮 Player movement vs popup navigation
- 🎮 Input conflict resolution
- 🎮 User experience verification

**Note:** May have graphics allocation issues on macOS 15.0 due to Ebiten/Metal compatibility.

## Running All Tests

To run all test programs in sequence:

```bash
# Popup functionality test
echo "=== Popup Widget Test ==="
go run test/popup_test/main.go

echo -e "\n=== Input Blocking Logic Test ==="  
go run test/logic_test/main.go

echo -e "\n=== Interactive Input Test ==="
# Note: May fail due to graphics issues on macOS 15.0
go run test/input_test/main.go
```

## Test Results Summary

All tests verify the **Popup Selection Widget** implementation:

### ✅ Core Functionality Verified
- Widget creation, display, and interaction
- Scrollable option lists with navigation
- Selection and cancellation callbacks
- Position and size customization

### ✅ Input Handling Fixed
- **Problem:** Arrow keys controlled both popup navigation AND player movement simultaneously
- **Solution:** Added input blocking in engine when popup is visible
- **Result:** Clean separation - arrow keys only control active UI element

### ✅ Integration Complete
- Widget integrated into UIManager
- Engine respects popup input priority
- Test trigger (P key) implemented in main game

## Architecture Notes

### Input Processing Flow
```
1. Engine.Update()
2. ├── uiManager.Update()     // Process popup input first
3. ├── Check IsPopupVisible() // Determine if game input should be blocked  
4. ├── IF popup visible:      // Block game input
5. │   └── return nil         // Skip game logic entirely
6. └── ELSE:                  // Normal game input
7.     └── updateExploration() or updateTactical()
```

### Key Design Principles
- **UI Input Priority:** UI elements get first access to input events
- **Exclusive Focus:** When popup is active, only popup processes input
- **Clean State Management:** Game input resumes automatically when popup closes
- **Extensible Pattern:** Works for any future UI widgets needing input focus

## Future Tests

Potential additional test programs:

- **Combat System Test:** Verify attack/damage calculation
- **Grid Positioning Test:** Verify tactical movement and positioning  
- **Save/Load Test:** Verify game state persistence
- **Performance Test:** Measure rendering and update performance
- **UI Layout Test:** Verify responsive UI scaling and positioning

## Troubleshooting

### Graphics Allocation Errors
If you see `[CAMetalLayer nextDrawable] returning nil because allocation failed`:
- This is a macOS 15.0 + Ebiten compatibility issue
- Non-interactive tests (logic_test, popup_test) still work
- Interactive tests may fail to display graphics
- The underlying logic still functions correctly

### Build Errors
Ensure you're running from the project root:
```bash
cd /path/to/myrpg
go run test/[test_name]/main.go
```

### Import Errors
All test programs use the same import structure as the main game:
```go
import "github.com/jrecuero/myrpg/internal/ui"
```

Make sure your `go.mod` is properly configured in the project root.