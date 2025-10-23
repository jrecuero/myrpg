# Inventory Widget Constants

This document defines the visual and layout constants for the Inventory Widget, following the established popup widget architecture.

## Widget Dimensions

| Constant | Value | Description |
|----------|--------|-------------|
| `InventoryWidgetWidth` | 600 | Total widget width in pixels |
| `InventoryWidgetHeight` | 500 | Total widget height in pixels |
| `InventoryWidgetPadding` | 20 | Internal padding around content |
| `InventoryWidgetBorderWidth` | 3 | Border thickness in pixels |

## Grid Layout

| Constant | Value | Description |
|----------|--------|-------------|
| `InventoryGridWidth` | 8 | Number of columns in inventory grid |
| `InventoryGridHeight` | 6 | Number of rows in inventory grid |
| `InventorySlotSize` | 48 | Width/height of each inventory slot |
| `InventorySlotSpacing` | 4 | Space between inventory slots |
| `InventoryGridStartX` | 40 | X offset for grid from widget left |
| `InventoryGridStartY` | 80 | Y offset for grid from widget top |

## Header Section

| Constant | Value | Description |
|----------|--------|-------------|
| `InventoryHeaderHeight` | 50 | Height of header area |
| `InventoryTitleX` | 25 | X position for title text |
| `InventoryTitleY` | 25 | Y position for title text |
| `InventoryStatsX` | 400 | X position for capacity stats |
| `InventoryStatsY` | 25 | Y position for capacity stats |

## Tooltip System

| Constant | Value | Description |
|----------|--------|-------------|
| `InventoryTooltipWidth` | 250 | Maximum tooltip width |
| `InventoryTooltipMaxHeight` | 200 | Maximum tooltip height |
| `InventoryTooltipPadding` | 8 | Internal tooltip padding |
| `InventoryTooltipBorder` | 2 | Tooltip border thickness |
| `InventoryTooltipOffsetX` | 10 | Horizontal offset from cursor |
| `InventoryTooltipOffsetY` | -10 | Vertical offset from cursor |

## Action Panel

| Constant | Value | Description |
|----------|--------|-------------|
| `InventoryActionPanelWidth` | 140 | Width of action buttons panel |
| `InventoryActionPanelHeight` | 300 | Height of action buttons panel |
| `InventoryActionPanelX` | 440 | X position from widget left |
| `InventoryActionPanelY` | 80 | Y position from widget top |
| `InventoryActionButtonHeight` | 30 | Height of action buttons |
| `InventoryActionButtonSpacing` | 5 | Space between buttons |

## Colors - Background

| Constant | R | G | B | A | Description |
|----------|---|---|---|---|-------------|
| `InventoryWidgetBackgroundR` | 45 | 45 | 55 | 220 | Main background color |
| `InventoryWidgetBackgroundG` | 45 | | | |
| `InventoryWidgetBackgroundB` | 55 | | | |
| `InventoryWidgetBackgroundA` | 220 | | | |

## Colors - Border

| Constant | R | G | B | A | Description |
|----------|---|---|---|---|-------------|
| `InventoryWidgetBorderR` | 120 | 120 | 140 | 255 | Widget border color |
| `InventoryWidgetBorderG` | 120 | | | |
| `InventoryWidgetBorderB` | 140 | | | |
| `InventoryWidgetBorderA` | 255 | | | |

## Colors - Inventory Slots

| Constant | R | G | B | A | Description |
|----------|---|---|---|---|-------------|
| `InventorySlotEmptyR` | 60 | 60 | 70 | 180 | Empty slot background |
| `InventorySlotEmptyG` | 60 | | | |
| `InventorySlotEmptyB` | 70 | | | |
| `InventorySlotEmptyA` | 180 | | | |
| `InventorySlotFilledR` | 80 | 80 | 95 | 200 | Filled slot background |
| `InventorySlotFilledG` | 80 | | | |
| `InventorySlotFilledB` | 95 | | | |
| `InventorySlotFilledA` | 200 | | | |
| `InventorySlotSelectedR` | 100 | 150 | 200 | 230 | Selected slot highlight |
| `InventorySlotSelectedG` | 150 | | | |
| `InventorySlotSelectedB` | 200 | | | |
| `InventorySlotSelectedA` | 230 | | | |
| `InventorySlotHoverR` | 90 | 130 | 180 | 200 | Hover slot highlight |
| `InventorySlotHoverG` | 130 | | | |
| `InventorySlotHoverB` | 180 | | | |
| `InventorySlotHoverA` | 200 | | | |

## Colors - Slot Borders

| Constant | R | G | B | A | Description |
|----------|---|---|---|---|-------------|
| `InventorySlotBorderR` | 100 | 100 | 120 | 255 | Normal slot border |
| `InventorySlotBorderG` | 100 | | | |
| `InventorySlotBorderB` | 120 | | | |
| `InventorySlotBorderA` | 255 | | | |
| `InventorySlotSelectedBorderR` | 150 | 200 | 255 | 255 | Selected border |
| `InventorySlotSelectedBorderG` | 200 | | | |
| `InventorySlotSelectedBorderB` | 255 | | | |
| `InventorySlotSelectedBorderA` | 255 | | | |

## Colors - Item Rarity

| Rarity | R | G | B | Description |
|--------|---|---|---|-------------|
| Common | 200 | 200 | 200 | Light Gray |
| Uncommon | 30 | 255 | 0 | Green |
| Rare | 0 | 112 | 255 | Blue |
| Epic | 163 | 53 | 238 | Purple |
| Legendary | 255 | 128 | 0 | Orange |

## Colors - Tooltip

| Constant | R | G | B | A | Description |
|----------|---|---|---|---|-------------|
| `InventoryTooltipBackgroundR` | 25 | 25 | 35 | 240 | Tooltip background |
| `InventoryTooltipBackgroundG` | 25 | | | |
| `InventoryTooltipBackgroundB` | 35 | | | |
| `InventoryTooltipBackgroundA` | 240 | | | |
| `InventoryTooltipBorderR` | 150 | 150 | 170 | 255 | Tooltip border |
| `InventoryTooltipBorderG` | 150 | | | |
| `InventoryTooltipBorderB` | 170 | | | |
| `InventoryTooltipBorderA` | 255 | | | |

## Colors - Action Panel

| Constant | R | G | B | A | Description |
|----------|---|---|---|---|-------------|
| `InventoryActionPanelR` | 35 | 35 | 45 | 200 | Action panel background |
| `InventoryActionPanelG` | 35 | | | |
| `InventoryActionPanelB` | 45 | | | |
| `InventoryActionPanelA` | 200 | | | |
| `InventoryActionButtonR` | 70 | 80 | 90 | 255 | Button background |
| `InventoryActionButtonG` | 80 | | | |
| `InventoryActionButtonB` | 90 | | | |
| `InventoryActionButtonA` | 255 | | | |
| `InventoryActionButtonHoverR` | 90 | 100 | 120 | 255 | Button hover |
| `InventoryActionButtonHoverG` | 100 | | | |
| `InventoryActionButtonHoverB` | 120 | | | |
| `InventoryActionButtonHoverA` | 255 | | | |

## Shadow and Effects

| Constant | Value | Description |
|----------|--------|-------------|
| `InventoryWidgetShadowOffset` | 5 | Shadow offset in pixels |
| `InventoryWidgetShadowR` | 0 | Shadow color red |
| `InventoryWidgetShadowG` | 0 | Shadow color green |
| `InventoryWidgetShadowB` | 0 | Shadow color blue |
| `InventoryWidgetShadowA` | 100 | Shadow transparency |

## Text and Fonts

| Constant | Value | Description |
|----------|--------|-------------|
| `InventoryTextColor` | White | Primary text color |
| `InventoryHeaderTextColor` | Light Yellow | Header text color |
| `InventoryQuantityTextSize` | Small | Size for quantity numbers |
| `InventoryTooltipTextSize` | Small | Size for tooltip text |

## Animation and Interaction

| Constant | Value | Description |
|----------|--------|-------------|
| `InventoryHoverAnimationSpeed` | 200ms | Hover effect transition |
| `InventorySelectionAnimationSpeed` | 150ms | Selection transition |
| `InventoryDragOpacity` | 128 | Opacity when dragging items |

## Sorting and Filtering

| Constant | Value | Description |
|----------|--------|-------------|
| `InventorySortButtonWidth` | 80 | Width of sort buttons |
| `InventorySortButtonHeight` | 25 | Height of sort buttons |
| `InventoryFilterButtonWidth` | 60 | Width of filter buttons |
| `InventoryFilterButtonHeight` | 20 | Height of filter buttons |

## Usage Notes

### Grid Calculations
- **Total Grid Width**: `(InventorySlotSize + InventorySlotSpacing) * InventoryGridWidth - InventorySlotSpacing`
- **Total Grid Height**: `(InventorySlotSize + InventorySlotSpacing) * InventoryGridHeight - InventorySlotSpacing`
- **Slot Position**: `StartX + (SlotSize + Spacing) * Column`, `StartY + (SlotSize + Spacing) * Row`

### Color Application
- **Empty Slots**: Use `InventorySlotEmpty` colors
- **Filled Slots**: Use `InventorySlotFilled` colors with rarity border
- **Interactive States**: Apply hover/selection overlays
- **Rarity Indication**: Use rarity colors for item borders or backgrounds

### Tooltip Positioning
- Position tooltips to avoid screen edges
- Use offset constants to separate from cursor
- Adjust position based on available space

This constants system provides comprehensive styling for a professional inventory interface that integrates seamlessly with the existing popup widget architecture!