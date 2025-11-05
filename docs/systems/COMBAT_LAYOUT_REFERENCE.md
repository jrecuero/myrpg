# Combat Layout Quick Reference

## Panel Dimensions & Positions

| Panel | X | Y | Width | Height | End X | End Y |
|-------|---|---|-------|--------|-------|-------|
| Enemy Panel | 5 | 5 | 791 | 222 | 796 | 227 |
| Player Panel | 5 | 232 | 791 | 222 | 796 | 454 |
| Combined Log | 801 | 5 | 194 | 449 | 995 | 454 |
| Action Menu | 801 | 459 | 194 | 136 | 995 | 595 |
| Battle Log | 5 | 459 | 791 | 136 | 796 | 595 |

## Key Constants

```go
// Spacing
PanelMargin = 5       // Screen edges
PanelSeparator = 5    // Between panels

// Formation
FormationRowWidth = 787    // 791 - 2 - 2
FormationRowHeight = 105
EntitySpriteBoxSize = 95

// Margins within formations
FormationMarginLeft = 2
FormationMarginRight = 2  
FormationMarginTop = 2
FormationMarginBottom = 11
FormationRowSeparator = 2
```

## Layout Validation

- **Screen Size**: 1000×600
- **Max Width Used**: 995 (5px margin remaining)
- **Max Height Used**: 595 (5px margin remaining)
- **Bottom Panel Alignment**: Both start at Y=459, height=136
- **Formation Rows**: Back (Y=2), Front (Y=109) within each panel

## File Locations

- **Constants**: `internal/constants/display.go`
- **Renderer**: `internal/battle/classic/renderer.go`  
- **Documentation**: `docs/systems/COMBAT_PANEL_LAYOUT.md`