# Popup Widgets Configuration Constants

This document describes the configurable constants used in the popup widget system for better maintainability and customization.

## PopupInfoWidget Constants

### Layout Dimensions
| Constant | Value | Description |
|----------|-------|-------------|
| `InfoWidgetLineHeight` | 16 | Height of each text line in pixels |
| `InfoWidgetTitleY` | 20 | Y offset for title from top of widget |
| `InfoWidgetTitleSpacing` | 30 | Space between title and content area |
| `InfoWidgetBottomPadding` | 30 | Space reserved for help text at bottom |
| `InfoWidgetSidePadding` | 10 | Left/right padding for text content |
| `InfoWidgetExtraPadding` | 10 | Additional padding for layout spacing |
| `InfoWidgetReservedSpace` | 90 | Total reserved space (calculated from above) |
| `InfoWidgetBorderWidth` | 2 | Border thickness in pixels |
| `InfoWidgetShadowOffset` | 4 | Shadow offset distance |

### Scrollbar Configuration
| Constant | Value | Description |
|----------|-------|-------------|
| `InfoWidgetScrollbarWidth` | 8 | Width of scrollbar in pixels |
| `InfoWidgetScrollbarRight` | 15 | Distance from right edge to scrollbar |
| `InfoWidgetMinThumbSize` | 20 | Minimum scrollbar thumb height |

### Color Defaults
| Component | R | G | B | A | Description |
|-----------|---|---|---|---|-------------|
| Background | 30 | 30 | 30 | 240 | Semi-transparent dark background |
| Border | 100 | 100 | 100 | 255 | Gray border |
| Text | 255 | 255 | 255 | 255 | White text |
| Title | 255 | 255 | 0 | 255 | Yellow title |
| Scrollbar | 160 | 160 | 160 | 255 | Light gray scrollbar |
| Shadow | 0 | 0 | 0 | 100 | Semi-transparent shadow |
| Scroll Background | 60 | 60 | 60 | 255 | Scrollbar track |

## PopupSelectionWidget Constants

### Layout Dimensions
| Constant | Value | Description |
|----------|-------|-------------|
| `SelectionWidgetItemHeight` | 20 | Height of each selectable option |
| `SelectionWidgetReservedSpace` | 40 | Space reserved for title and padding |
| `SelectionWidgetTitleY` | 15 | Y offset for title from top |
| `SelectionWidgetContentY` | 35 | Y offset for options content |
| `SelectionWidgetSidePadding` | 10 | Left/right padding for text |
| `SelectionWidgetItemPadding` | 2 | Padding between option items |
| `SelectionWidgetBorderWidth` | 2 | Border thickness |
| `SelectionWidgetShadowOffset` | 4 | Shadow offset distance |
| `SelectionWidgetTitlePadding` | 10 | Left padding for title text |
| `SelectionWidgetTitleMargin` | 10 | Top margin for title text |
| `SelectionWidgetNoTitleMargin` | 10 | Top margin when no title is present |
| `SelectionWidgetBottomMargin` | 10 | Bottom margin for instructions area |
| `SelectionWidgetInstrMargin` | 15 | Bottom margin for instruction text |
| `SelectionWidgetHighlightMargin` | 5 | Left margin for highlight rectangle |
| `SelectionWidgetHighlightPadding` | 2 | Vertical padding for highlight rectangle |
| `SelectionWidgetTextPadding` | 10 | Left padding for option text |

### Scrollbar Configuration
| Constant | Value | Description |
|----------|-------|-------------|
| `SelectionWidgetScrollbarWidth` | 10 | Width of scrollbar track |
| `SelectionWidgetScrollbarRight` | 15 | Distance from right edge |
| `SelectionWidgetScrollbarMargin` | 20 | Vertical margin for scrollbar area |
| `SelectionWidgetScrollbarMinThumb` | 10 | Minimum scrollbar thumb height |
| `SelectionWidgetScrollbarPadding` | 1 | Internal padding for scrollbar thumb |

### Color Defaults
| Component | R | G | B | A | Description |
|-----------|---|---|---|---|-------------|
| Background | 40 | 40 | 40 | 240 | Semi-transparent dark background |
| Border | 100 | 100 | 100 | 255 | Gray border |
| Text | 255 | 255 | 255 | 255 | White text |
| Title | 255 | 255 | 0 | 255 | Yellow title |
| Highlight | 70 | 130 | 180 | 200 | Steel blue selection highlight |
| Scrollbar | 160 | 160 | 160 | 255 | Light gray scrollbar |
| Shadow | 0 | 0 | 0 | 100 | Semi-transparent shadow |

## Usage Benefits

### Before Constants (Hardcoded Values)
```go
// Hard to maintain and customize
maxLines := (height - 60) / 16  // Magic numbers!
titleY := p.Y + 20
vector.StrokeRect(screen, x, y, width, height, 2, borderColor, false)
```

### After Constants (Configurable)
```go
// Clear, maintainable, and customizable
maxLines := (height - InfoWidgetReservedSpace) / InfoWidgetLineHeight
titleY := p.Y + InfoWidgetTitleY  
vector.StrokeRect(screen, x, y, width, height, InfoWidgetBorderWidth, borderColor, false)
```

## Customization Examples

### Changing Widget Appearance
```go
// Create widget with custom size
widget := NewPopupInfoWidget("Title", "Content", 100, 100, 500, 400)

// Customize colors after creation
widget.BackgroundColor = color.RGBA{20, 20, 20, 200}  // Darker background
widget.TitleColor = color.RGBA{0, 255, 255, 255}      // Cyan title
widget.BorderColor = color.RGBA{255, 0, 0, 255}       // Red border
```

### Modifying Layout Constants
To change the overall appearance, modify the constants in the source files:
```go
// Make widgets more compact
const InfoWidgetLineHeight = 14        // Smaller line height
const InfoWidgetTitleSpacing = 25      // Less space after title
const InfoWidgetBottomPadding = 25     // Less bottom padding
```

## Layout Calculations

### InfoWidget Space Allocation
```
Total Height: 300px
├── Title Area: 50px (TitleY + TitleSpacing)
├── Content Area: 160px (Height - ReservedSpace)
│   └── Lines: 10 lines (160px ÷ 16px per line)
└── Help Area: 40px (BottomPadding + ExtraPadding)
```

### SelectionWidget Space Allocation
```
Total Height: 200px  
├── Title Area: 40px (ReservedSpace)
├── Options Area: 160px (Height - ReservedSpace)
│   └── Items: 8 items (160px ÷ 20px per item)
```

## Implementation Notes

- **Constants are defined at package level** for easy access and modification
- **Color constants use separate R,G,B,A values** for maximum flexibility
- **Layout constants are used in multiple functions** ensuring consistency
- **Reserved space is calculated automatically** from component constants
- **All hardcoded values have been replaced** with named constants

This approach makes the popup widgets much more maintainable and allows for easy theming and customization without modifying core logic.