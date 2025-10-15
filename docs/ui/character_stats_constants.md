# Character Stats Widget Configuration Constants

This document describes the configurable constants used in the Character Stats Widget for better maintainability and customization.

## Layout Dimensions Constants

| Constant                      | Value | Description                             |
|-------------------------------|-------|-----------------------------------------|
| `StatsWidgetWidth`            | 600   | Widget width in pixels                  |
| `StatsWidgetHeight`           | 450   | Widget height in pixels                 |
| `StatsWidgetBorderWidth`      | 2     | Border thickness in pixels              |
| `StatsWidgetShadowOffset`     | 4     | Shadow offset distance                  |
| `StatsWidgetPadding`          | 15    | Internal padding for content            |
| `StatsWidgetLineHeight`       | 18    | Height of each text line                |
| `StatsWidgetSectionSpacing`   | 25    | Space between sections                  |
| `StatsWidgetColumnWidth`      | 280   | Width of each stat column               |
| `StatsWidgetColumnSpacing`    | 20    | Space between columns                   |
| `StatsWidgetHeaderHeight`     | 30    | Height of section headers               |
| `StatsWidgetTitleY`           | 20    | Y offset for main title                 |
| `StatsWidgetContentStartY`    | 50    | Y offset where content begins           |
| `StatsWidgetBottomReserved`   | 40    | Space reserved at bottom for help text  |

## Tab Navigation Constants

| Constant                      | Value | Description           |
|-------------------------------|-------|-----------------------|
| `StatsWidgetTabHeight`        | 25    | Height of tab buttons |
| `StatsWidgetTabSpacing`       | 5     | Space between tabs    |
| `StatsWidgetTabPadding`       | 10    | Internal tab padding  |

## Progress Bar Constants

| Constant                      | Value | Description             |
|-------------------------------|-------|-------------------------|
| `StatsWidgetBarWidth`         | 200   | Width of HP/MP/XP bars  |
| `StatsWidgetBarHeight`        | 12    | Height of progress bars |
| `StatsWidgetBarBorder`        | 1     | Border width for bars   |

## Color Constants (RGBA Values)

### Background and Border Colors
| Component  | R   | G   | B   | A   | Description                           |
|------------|-----|-----|-----|-----|---------------------------------------|
| Background | 25  | 25  | 35  | 245 | Semi-transparent dark blue background |
| Border     | 120 | 120 | 120 | 255 | Gray border                           |
| Shadow     | 0   | 0   | 0   | 120 | Semi-transparent shadow               |

### Text Colors
| Component | R   | G   | B   | A   | Description            |
|-----------|-----|-----|-----|-----|------------------------|
| Title     | 255 | 255 | 100 | 255 | Yellow title text      |
| Text      | 255 | 255 | 255 | 255 | White normal text      |
| Value     | 200 | 255 | 200 | 255 | Light green value text |

### Section Header Colors
| Component       | R   | G   | B   | A   | Description                      |
|-----------------|-----|-----|-----|-----|----------------------------------|
| Core Header     | 100 | 200 | 255 | 255 | Blue for core stats section      |
| Combat Header   | 255 | 150 | 100 | 255 | Orange for combat stats section  |
| Tactical Header | 150 | 255 | 150 | 255 | Green for tactical stats section |

### Progress Bar Colors
| Component      | R   | G   | B   | A   | Description                   |
|----------------|-----|-----|-----|-----|-------------------------------|
| HP Bar         | 200 | 50  | 50  | 255 | Red for health bars           |
| MP Bar         | 50  | 100 | 200 | 255 | Blue for mana bars            |
| XP Bar         | 255 | 200 | 50  | 255 | Gold for experience bars      |
| Bar Background | 40  | 40  | 40  | 255 | Dark gray for bar backgrounds |

### Tab Colors
| Component    | R  | G  | B   | A   | Description               |
|--------------|----|----|-----|-----|---------------------------|
| Active Tab   | 80 | 80 | 120 | 220 | Highlighted selected tab  |
| Inactive Tab | 50 | 50 | 70  | 180 | Dimmed unselected tabs    |

## Widget Categories

### StatCategory Enum Values
| Value                     | Description         | Content                                     |
|---------------------------|---------------------|---------------------------------------------|
| `StatCategoryOverview`    | General overview    | Progress bars, basic info, key stats        |
| `StatCategoryCore`        | Core statistics     | Level, XP, HP/MP details, progression       |
| `StatCategoryCombat`      | Combat statistics   | Attack, Defense, Magic stats, accuracy      |
| `StatCategoryTactical`    | Tactical statistics | Movement, initiative, combat positioning    |

## Usage Examples

### Creating a Character Stats Widget
```go
// Create widget centered on screen
statsX := (ScreenWidth - StatsWidgetWidth) / 2
statsY := (ScreenHeight - StatsWidgetHeight) / 2
widget := NewCharacterStatsWidget(statsX, statsY, characterData)
```

### Customizing Widget Colors
```go
// Change background to darker theme
widget.BackgroundColor = color.RGBA{15, 15, 25, 245}

// Use different title color
widget.TitleColor = color.RGBA{100, 255, 255, 255}  // Cyan title

// Customize section headers
// Core section: Blue theme
coreHeaderColor := color.RGBA{100, 200, 255, 255}

// Combat section: Red theme  
combatHeaderColor := color.RGBA{255, 100, 100, 255}

// Tactical section: Green theme
tacticalHeaderColor := color.RGBA{100, 255, 100, 255}
```

### Modifying Layout Constants
```go
// Create more compact widget
const StatsWidgetLineHeight = 16        // Smaller line spacing
const StatsWidgetSectionSpacing = 20    // Less space between sections
const StatsWidgetPadding = 10           // Tighter padding

// Create larger progress bars
const StatsWidgetBarWidth = 250         // Wider bars
const StatsWidgetBarHeight = 15         // Taller bars
```

## Widget Integration

### UIManager Integration
```go
// UIManager automatically includes character stats widget
ui := NewUIManager()  // Creates widget internally

// Show character stats for active player
ui.ShowCharacterStats(player.RPGStats())

// Hide widget
ui.HideCharacterStats()

// Check if any popup (including stats) is visible
if ui.IsPopupVisible() {
    // Handle blocked input
}
```

### Engine Integration
```go
// In main game engine (engine.go):
if inpututil.IsKeyJustPressed(ebiten.KeyC) && !g.uiManager.IsPopupVisible() {
    g.showCharacterStats()  // Shows active player's stats
}

// showCharacterStats method:
func (g *Game) showCharacterStats() {
    activePlayer := g.GetActivePlayer()
    if activePlayer != nil && activePlayer.RPGStats() != nil {
        g.uiManager.ShowCharacterStats(activePlayer.RPGStats())
    }
}
```

## Input Controls

### Navigation
| Key           | Action            |
|---------------|-------------------|
| `←`           | Previous category |
| `→`           | Next category     |
| `TAB`         | Next category     |
| `Shift + TAB` | Previous category |
| `ESC`         | Close widget      |

### Category Navigation Flow
```
Overview → Core Stats → Combat Stats → Tactical Stats → (loops back to Overview)
```

## Layout Calculations

### Content Area Distribution
```
Widget Height: 450px
├── Title Area: 50px (TitleY + spacing)
├── Tab Area: 35px (TabHeight + spacing)
├── Content Area: 325px (Height - TitleY - ContentStartY - BottomReserved)
└── Help Area: 40px (BottomReserved)
```

### Two-Column Layout (Overview)
```
Widget Width: 600px
├── Left Column: 280px (StatsWidgetColumnWidth)
├── Column Spacing: 20px (StatsWidgetColumnSpacing)  
└── Right Column: 280px (remaining space - padding)
```

## Testing

### Test Program Usage
```bash
# Run character stats widget test
make test-character-stats

# Test controls:
# C - Show/Hide Character Stats
# R - Reset Character (new random stats)
# H - Heal Character (+50 HP, +25 MP)
# D - Damage Character (-30 HP, -15 MP)
# L - Gain Experience (+200 XP)
```

## Implementation Notes

- **Category System**: Widget uses enum-based category system for organized display
- **Progress Bars**: Custom progress bar implementation with configurable colors
- **Input Blocking**: Integrates with UIManager's input blocking system
- **Dynamic Data**: Widget updates automatically when character data changes
- **Error Handling**: Gracefully handles missing character data
- **Responsive Design**: Layout adapts to different content lengths

## Benefits of Constants System

### Before Constants (Hardcoded)
```go
// Hard to maintain and customize
vector.DrawFilledRect(screen, x, y, 600, 450, bgColor, false)
ebitenutil.DebugPrintAt(screen, title, x+15, y+20)
```

### After Constants (Configurable)  
```go
// Clear, maintainable, and customizable
vector.DrawFilledRect(screen, x, y, float32(StatsWidgetWidth), float32(StatsWidgetHeight), bgColor, false)
ebitenutil.DebugPrintAt(screen, title, x+StatsWidgetPadding, y+StatsWidgetTitleY)
```

This constants system makes the Character Stats Widget highly customizable and maintainable, allowing easy theming and layout adjustments without modifying core logic.