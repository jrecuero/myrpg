# Popup Selection Widget Documentation

## Overview

The `PopupSelectionWidget` is a reusable UI component that displays a modal popup with a scrollable list of selectable options. It provides keyboard navigation, customizable styling, and callback-based interaction.

## Features

### Core Functionality
- ✅ **Modal popup display** with semi-transparent background
- ✅ **Scrollable option list** with automatic scrollbar when needed
- ✅ **Keyboard navigation** (Arrow keys, W/S keys)
- ✅ **Selection handling** (Enter/Space keys)
- ✅ **Cancellation support** (Escape key)
- ✅ **Callback system** for selection and cancel events
- ✅ **Customizable positioning and sizing**
- ✅ **Visual styling** with highlights and shadows

### Visual Elements
- Title bar with customizable text
- Scrollable options list with selection highlighting
- Scrollbar indicator for long lists
- Semi-transparent shadow effect
- Instruction text at bottom
- Color-coded styling system

## API Reference

### Creation
```go
popup := ui.NewPopupSelectionWidget(title, options, x, y, width, height)
```

### Display Control
```go
// Show popup with new content
popup.Show(title, options)

// Hide popup
popup.Hide()

// Check visibility
isVisible := popup.IsVisible
```

### Position & Size
```go
// Update position
popup.SetPosition(x, y)

// Update size 
popup.SetSize(width, height)
```

### Callbacks
```go
// Selection callback (called when user presses Enter/Space)
popup.OnSelection = func(index int, option string) {
    // Handle selection
    fmt.Printf("Selected: %s at index %d", option, index)
}

// Cancel callback (called when user presses Escape)
popup.OnCancel = func() {
    // Handle cancellation
    fmt.Println("Selection cancelled")
}
```

### State Access
```go
// Get current selection
index, option := popup.GetSelectedOption()

// Access properties directly
fmt.Printf("Position: (%d, %d)", popup.X, popup.Y)
fmt.Printf("Size: %dx%d", popup.Width, popup.Height)
fmt.Printf("Selected: %d", popup.SelectedIndex)
```

## Integration with UIManager

The widget is integrated into the main UI system:

```go
// Show popup through UIManager
uiManager.ShowSelectionPopup(title, options, onSelection, onCancel)

// Hide popup
uiManager.HideSelectionPopup()

// Check if any popup is visible
if uiManager.IsPopupVisible() {
    // Handle popup-specific input blocking
}
```

## Usage Examples

### Combat Action Menu
```go
options := []string{
    "Attack Enemy",
    "Cast Spell", 
    "Use Item",
    "Move Unit",
    "End Turn",
}

uiManager.ShowSelectionPopup("Combat Actions", options, 
    func(index int, option string) {
        switch index {
        case 0: // Attack Enemy
            startAttackMode()
        case 1: // Cast Spell
            showSpellMenu()
        case 2: // Use Item
            showInventory()
        // ... handle other options
        }
    },
    func() {
        // User cancelled - return to normal input mode
        resumeNormalInput()
    })
```

### Unit Selection
```go
unitNames := []string{}
for _, unit := range availableUnits {
    unitNames = append(unitNames, unit.Name)
}

uiManager.ShowSelectionPopup("Select Unit", unitNames,
    func(index int, name string) {
        selectedUnit := availableUnits[index]
        setActiveUnit(selectedUnit)
    },
    func() {
        // Cancel unit selection
    })
```

### Item/Spell Selection
```go
spellNames := []string{}
for _, spell := range playerSpells {
    spellNames = append(spellNames, fmt.Sprintf("%s (MP: %d)", spell.Name, spell.ManaCost))
}

uiManager.ShowSelectionPopup("Cast Spell", spellNames,
    func(index int, spellText string) {
        spell := playerSpells[index]
        castSpell(spell)
    },
    nil) // No cancel callback needed
```

## Keyboard Controls

| Key | Action |
|-----|--------|
| ↑ / W | Move selection up |
| ↓ / S | Move selection down |
| Enter / Space | Select current option |
| Escape | Cancel and close popup |

## Customization

### Styling Properties
```go
// Colors can be customized
popup.BackgroundColor = color.RGBA{40, 40, 40, 240}   // Semi-transparent dark
popup.BorderColor = color.RGBA{100, 100, 100, 255}     // Gray border
popup.TextColor = color.RGBA{255, 255, 255, 255}       // White text
popup.HighlightColor = color.RGBA{70, 130, 180, 200}   // Blue highlight
popup.TitleColor = color.RGBA{255, 255, 0, 255}        // Yellow title
popup.ScrollbarColor = color.RGBA{160, 160, 160, 255}  // Light gray scrollbar
popup.ShadowColor = color.RGBA{0, 0, 0, 100}          // Semi-transparent shadow
```

### Automatic Sizing
- The widget automatically calculates how many items can be displayed
- Scrollbar appears when options exceed visible area
- Item height is fixed at 20 pixels
- Minimum displayable items is 1

## Game Integration

### Engine Integration
The popup system is integrated into the main game engine:

```go
// In Update() method
uiManager.Update() // Updates popup input handling

// Test trigger (P key)
if inpututil.IsKeyJustPressed(ebiten.KeyP) && !uiManager.IsPopupVisible() {
    showTestPopup()
}

// In Draw() method  
uiManager.DrawPopups(screen) // Renders popups on top
```

### Input Blocking
When a popup is visible, you should typically block other game input:

```go
func (g *Game) Update() error {
    // Always update UI first
    g.uiManager.Update()
    
    // Block game input when popup is visible
    if g.uiManager.IsPopupVisible() {
        return nil // Don't process game input
    }
    
    // Normal game input processing
    return g.updateGameplay()
}
```

## File Structure

```
internal/ui/
├── popup_selection_widget.go  # Main widget implementation
├── ui_manager.go              # Integration with UI system
└── ...

test/
├── popup_test/               # Popup widget logic test
├── logic_test/               # Input blocking verification  
└── input_test/               # Interactive input test
```

## Testing

### Logic Test
Run the logic verification test:
```bash
cd /path/to/myrpg
go run test/popup_test/main.go
```

### In-Game Test
1. Run the main game
2. Press 'P' key to show test popup
3. Navigate with arrow keys
4. Select with Enter or cancel with Escape

## Future Enhancements

Potential improvements for the popup widget:

- **Multi-column layouts** for wide popups
- **Icons/images** alongside text options
- **Nested menus** with sub-options
- **Search/filter** functionality
- **Mouse support** for clicking
- **Animation effects** for show/hide
- **Keyboard shortcuts** for quick selection (1-9 keys)
- **Option descriptions** with details panel
- **Checkbox/radio button** style selections

## Performance Notes

- Widget uses immediate mode rendering (redraws every frame when visible)
- Scrolling is virtualized (only visible items are rendered)
- Input handling includes debouncing to prevent rapid key repeats
- Memory usage is minimal (no texture caching, uses vector graphics)

## Troubleshooting

### Common Issues

**Popup not appearing:**
- Check if `IsVisible` is true after calling `Show()`
- Ensure `DrawPopups()` is called after other rendering
- Verify popup position is within screen bounds

**Navigation not working:**
- Ensure `Update()` is called before input processing
- Check if other input handlers are consuming key events
- Verify popup has focus (no other UI elements blocking input)

**Callbacks not firing:**
- Ensure callbacks are set before showing popup
- Check if popup is properly hidden after selection
- Verify callback functions are not nil

**Visual issues:**
- Check popup dimensions are reasonable (not too small)
- Ensure colors have proper alpha values for visibility
- Verify z-order (popups should render last/on top)