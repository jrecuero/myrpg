# Test Programs

This directory contains standalone test programs for verifying specific functionality of the MyRPG game components.

## Directory Structure

```
test/
├── README.md              # This file
├── input_test/            # Interactive input blocking test 
├── logic_test/            # Non-interactive logic verification
├── popup_test/            # Popup selection widget functionality test
├── info_popup_test/       # Popup info widget functionality test
├── dialog_test/           # Dialog widget conversation system test
├── equipment_test/        # Equipment management widget test  
├── character_stats_test/  # Character statistics widget test
├── info_layout_test/      # Info layout widget test
└── inventory_test/        # Inventory management widget test
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

### 2. **`info_popup_test/`** - Info Popup Widget Test
- **File**: `test/info_popup_test/main.go`
- **Purpose**: Comprehensive testing of PopupInfoWidget functionality
- **Type**: Interactive test (requires display/graphics)
- **Features Tested**:

**Usage:**
```bash
# Info popup widget test
go run test/info_popup_test/main.go
```

**What it tests:**
- ✅ Info widget creation and display
- ✅ Multi-line text content display
- ✅ Scrollable content with arrow key navigation
- ✅ Title display and close functionality (ESC)
- ✅ Input blocking integration
- ✅ Long content handling and scrollbar display

**Controls:**
- Press 'I' to show info popup with test content
- Use ↑↓ arrows to scroll content (when popup is open)
- Press ESC to close popup
- ESC exits test when no popup is open

**Note:** May have graphics allocation issues on macOS 15.0 due to Ebiten/Metal compatibility.

---

### 9. **`inventory_test/`** - Inventory Management Widget Test
- **File**: `test/inventory_test/main.go`
- **Purpose**: Comprehensive testing of InventoryWidget functionality
- **Type**: Interactive test (requires display/graphics)
- **Features Tested**:

**Usage:**
```bash
# Inventory widget test
go run test/inventory_test/main.go
```

**What it tests:**
- ✅ Grid-based inventory display (8x6 = 48 slots)
- ✅ Item creation and display with different types and rarities
- ✅ Drag & drop functionality between inventory slots
- ✅ Item stacking for stackable items (potions, materials)
- ✅ Item tooltips on hover with detailed information
- ✅ Sorting functionality (by name, type, rarity)
- ✅ Filtering by item category (equipment, consumables, all)
- ✅ Keyboard shortcuts (Delete to remove, 1-9 to split stacks)
- ✅ Rarity-based color coding (Common: Gray, Rare: Blue, Epic: Purple, etc.)
- ✅ Action panel with sort/filter buttons
- ✅ ECS integration with Entity and InventoryComponent

**Controls:**
- Press 'I' to toggle inventory open/close
- Left-click to select items and start dragging
- Right-click to select items (context menu ready for future)
- Drag & drop items between slots
- Use action panel buttons for sorting/filtering
- Delete key to remove selected item
- Number keys (1-9) to split selected stack
- Press ESC to close inventory

**Test Items Included:**
- Equipment: Iron Sword (Common), Steel Helmet (Uncommon), Magic Ring (Epic)
- Consumables: Health Potions (stackable), Mana Potions (stackable), Greater Elixir (Legendary)
- Materials: Iron Ore (stackable x25)
- Quest Items: Mysterious Key (Rare, unique)

**Note:** Demonstrates complete inventory management system ready for integration.

---

### 3. **`logic_test/`** - Input Blocking Logic Test  
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

### 4. **`input_test/`** - Interactive Input Test
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
echo "=== Popup Selection Widget Test ==="
go run test/popup_test/main.go

echo -e "\n=== Info Popup Widget Test ==="
# Note: May fail due to graphics issues on macOS 15.0
go run test/info_popup_test/main.go

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