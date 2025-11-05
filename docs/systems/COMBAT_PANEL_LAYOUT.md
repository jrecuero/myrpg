# Combat Panel Layout System Documentation

## Overview

This document provides complete specifications for the MyRPG combat interface panel layout system. The layout is designed for a 1000x600 pixel screen with precise panel positioning, consistent margins, and formation-based entity display.

## Layout Architecture

### Screen Dimensions
- **Total Screen**: 1000×600 pixels
- **Panel Margin**: 5px from screen edges
- **Panel Separator**: 5px between adjacent panels

### Panel Configuration

The combat interface consists of 5 main panels arranged in a structured layout:

```
┌─────────────────────────────────────────────┬─────────┐
│ Enemy Panel (791×222)                       │ Combined│
├─────────────────────────────────────────────┤ Log     │
│ Player Panel (791×222)                      │ Panel   │
├─────────────────────────────────────────────┤ (194×   │
│ Battle Log Panel (791×136)                  │ 449)    │
│                                             ├─────────┤
│                                             │ Action  │
│                                             │ Menu    │
│                                             │ (194×   │
│                                             │ 136)    │
└─────────────────────────────────────────────┴─────────┘
```

## Panel Specifications

### 1. Enemy Panel
- **Position**: (5, 5)
- **Dimensions**: 791×222 pixels
- **Purpose**: Display enemy formations in back/front rows
- **Layout**: Two-row formation system with 95×95 entity boxes

### 2. Player Panel  
- **Position**: (5, 232)
- **Dimensions**: 791×222 pixels
- **Purpose**: Display player party formations in back/front rows
- **Layout**: Two-row formation system with 95×95 entity boxes

### 3. Combined Log Panel
- **Position**: (801, 5)
- **Dimensions**: 194×449 pixels
- **Purpose**: Display activity queue and player information
- **Alignment**: Spans from top to align with bottom panel row

### 4. Action Menu Panel
- **Position**: (801, 459)
- **Dimensions**: 194×136 pixels
- **Purpose**: Battle action selection interface
- **Alignment**: Same height and Y-position as Battle Log Panel

### 5. Battle Log Panel
- **Position**: (5, 459)
- **Dimensions**: 791×136 pixels
- **Purpose**: Display combat messages and battle events
- **Alignment**: Same height and Y-position as Action Menu Panel

## Formation Layout System

### Panel Internal Structure
Both Enemy and Player panels use identical formation layouts:

```go
// Formation margins within each panel
FormationMarginLeft   = 2   // Left edge margin
FormationMarginRight  = 2   // Right edge margin  
FormationMarginTop    = 2   // Top edge margin
FormationMarginBottom = 11  // Bottom edge margin
FormationRowSeparator = 2   // Space between rows
```

### Row Specifications
- **Row Dimensions**: 787×105 pixels each
  - Width: 791 (panel width) - 2 (left) - 2 (right) = 787px
  - Height: 105px per row
- **Back Row Position**: (2, 2) within panel
- **Front Row Position**: (2, 109) within panel

### Entity Display Boxes
Each entity is contained within a standardized display box:

- **Box Size**: 95×95 pixels
- **Content Layout**:
  - **Name**: Top area (5px from top)
  - **Sprite**: Center area (32×32 pixels, centered horizontally)
  - **Health Bar**: Bottom area (below sprite + 5px)
- **Margins**: 5px top and bottom within each row
- **Spacing**: Automatically calculated based on row width and entity count

## Constants Reference

### File Location
`/Users/jorecuer/go/src/github.com/jrecuero/myrpg/internal/constants/display.go`

### Key Constants
```go
// Panel Layout
PanelMargin    = 5  // Screen edge margins
PanelSeparator = 5  // Inter-panel spacing

// Main Panels
EnemyPanelWidth   = 791
EnemyPanelHeight  = 222
PlayerPanelWidth  = 791  
PlayerPanelHeight = 222
CombinedLogPanelWidth  = 194
CombinedLogPanelHeight = 449
ActionMenuPanelWidth   = 194
ActionMenuPanelHeight  = 136
BattleLogPanelWidth    = 791
BattleLogPanelHeight   = 136

// Formation System
FormationRowWidth     = 787
FormationRowHeight    = 105
EntitySpriteBoxSize   = 95
```

## Positioning Calculations

### Vertical Layout
```
Y=5    : Enemy Panel starts
Y=227  : Enemy Panel ends (5 + 222)
Y=232  : Player Panel starts (227 + 5 separator)
Y=454  : Player Panel ends (232 + 222)
Y=459  : Bottom panels start (454 + 5 separator)
Y=595  : Bottom panels end (459 + 136)
Y=600  : Screen bottom (595 + 5 margin)
```

### Horizontal Layout
```
X=5    : Left panels start (Enemy, Player, Battle Log)
X=796  : Left panels end (5 + 791)
X=801  : Right panels start (796 + 5 separator)
X=995  : Right panels end (801 + 194)
X=1000 : Screen right edge (995 + 5 margin)
```

## Alignment Features

### Perfect Symmetry
- **Bottom Row Alignment**: Action Menu and Battle Log panels share identical Y-coordinates (459-595)
- **Height Consistency**: Both bottom panels are exactly 136 pixels tall
- **Margin Uniformity**: 5px margins maintained on all screen edges

### Visual Balance
- **Left Column**: 791px wide (79.1% of screen width)
- **Right Column**: 194px wide (19.4% of screen width)  
- **Separation**: 5px (0.5% of screen width)
- **Top Section**: 454px tall (75.7% of screen height)
- **Bottom Section**: 136px tall (22.7% of screen height)

## Implementation Notes

### Renderer Integration
The layout constants are consumed by:
- `internal/battle/classic/renderer.go` - Main rendering implementation
- Panel drawing methods use constants directly for positioning
- Formation system automatically distributes entities within defined areas

### Flexibility
The constant-based system allows easy adjustments:
1. **Margin Changes**: Modify `PanelMargin` and `PanelSeparator`
2. **Size Adjustments**: Update individual panel width/height constants
3. **Formation Tuning**: Adjust entity box size and row dimensions
4. **Alignment**: All positioning calculations are relative and automatic

### Future Modifications
When making layout changes:
1. Update constants in `display.go`
2. Verify total dimensions don't exceed screen bounds
3. Ensure margin and separator consistency
4. Test formation entity distribution
5. Update this documentation

## Version History

- **v1.0**: Initial implementation with 2px margins
- **v1.1**: Increased margins to 5px for better visual spacing
- **v1.2**: Aligned bottom panels to same height (136px) for symmetry
- **Current**: Optimized layout with perfect alignment and consistent spacing

---

*This layout provides an optimal balance of information density, visual clarity, and consistent spacing for the MyRPG combat interface.*