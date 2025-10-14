# Column 0 Movement Issue - Diagnostic Report

## Problem Description

**Issue**: Players can move to any column except column 0 (X=0). When trying to move within column 0, the system reports "Cannot move to occupied tile" even for apparently empty tiles.

## Behavior Analysis

### What Works ✅
- Players can move to columns 1, 2, 3, etc.
- Players can move TO column 0 from other columns 
- Players can be deployed initially at column 0 positions

### What Doesn't Work ❌
- Players in column 0 cannot move to other positions within column 0
- System reports "Cannot move to occupied tile" for column 0 movements

## Debugging Steps Implemented

### 1. Enhanced Coordinate Conversion
**Purpose**: Ensure symmetric world ↔ grid coordinate conversion

**Implementation**:
```go
func (g *Game) worldToGridPos(worldX, worldY float64) tactical.GridPos {
    offsetX, offsetY := 50.0, 120.0
    tileSize := float64(g.tacticalManager.Grid.TileSize)
    gridX := int((worldX - offsetX) / tileSize)
    gridZ := int((worldY - offsetY) / tileSize)
    return tactical.GridPos{X: gridX, Z: gridZ}
}
```

### 2. Detailed Debug Information
**Added Debug Messages**:
- World coordinates → Grid coordinates conversion
- Current position vs target position
- Occupancy clearing operations
- Unit IDs occupying blocked tiles

**Example Output**:
```
World pos: (50.0, 152.0) -> Grid pos: (0, 1)
Trying to move from (0, 1) to (0, 2) 
Clearing occupancy at (0, 1)
Cannot move to (0,2) - occupied by unit Player123
```

### 3. Deployment Diagnostics
**Added Deployment Logging**:
```go
fmt.Printf("DEPLOY: %s at grid(%d,%d) -> world(%.1f,%.1f) -> final(%.1f,%.1f)\n",
    member.Name, gridPos.X, gridPos.Z, worldX, worldY, transform.X, transform.Y)
```

## Coordinate System Analysis

### Grid Layout (20x15 tiles)
```
Column:  0   1   2   3   4  ...  19
Row 0:  (0,0)(1,0)(2,0)(3,0)...  (19,0)
Row 1:  (0,1)(1,1)(2,1)(3,1)...  (19,1)
...
```

### World Coordinates
- **Grid (0,0)** → **World (50.0, 120.0)** 
- **Grid (1,0)** → **World (82.0, 120.0)** (50 + 32*1)
- **Grid (0,1)** → **World (50.0, 152.0)** (120 + 32*1)

### Conversion Test Cases
```
Grid (0,0) → World (50, 120) → Grid (0,0) ✓
Grid (0,1) → World (50, 152) → Grid (0,1) ✓  
Grid (1,0) → World (82, 120) → Grid (1,0) ✓
```

**Mathematical Verification**: 
- `(50 - 50) / 32 = 0 / 32 = 0` ✓
- `(152 - 120) / 32 = 32 / 32 = 1` ✓

## Hypothesis Analysis

### ❌ Hypothesis 1: Coordinate Conversion Error
**Status**: UNLIKELY
**Reason**: Math is correct, works for other columns

### ❌ Hypothesis 2: Floating Point Precision
**Status**: UNLIKELY  
**Reason**: Uses exact integer multiples (32px)

### ✅ Hypothesis 3: Grid State Corruption
**Status**: LIKELY
**Reason**: Issue is movement-specific, not position-specific

### ✅ Hypothesis 4: Occupancy Clearing Bug
**Status**: MOST LIKELY
**Reason**: Can move TO column 0, but not WITHIN column 0

## Potential Root Causes

### 1. Race Condition in Grid State Updates
```go
// This sequence might have a bug:
g.tacticalManager.Grid.SetOccupied(currentPos, false, "") // Clear old
g.tacticalManager.Grid.SetOccupied(gridPos, true, player.GetID()) // Set new
```

### 2. Entity ID Mismatch
- Player deployed with one ID
- Movement trying to clear with different ID
- Grid tile remains "occupied" by phantom unit

### 3. Coordinate Conversion Edge Case
- Column 0 has special case: `worldX = 50.0` exactly
- Other columns: `worldX = 82.0, 114.0, etc.`
- Potential floating point edge case with exact offset match

### 4. Multiple Players in Same Column
- Initial deployment places multiple players in column 0
- Grid state gets confused about which unit occupies which tile
- Occupancy not cleared properly when switching between units

## Recommended Fixes

### Fix 1: Add Grid State Validation
```go
func (g *Game) validateGridState() {
    for x := 0; x < g.tacticalManager.Grid.Width; x++ {
        for z := 0; z < g.tacticalManager.Grid.Height; z++ {
            pos := tactical.GridPos{X: x, Z: z}
            tile := g.tacticalManager.Grid.GetTile(pos) 
            if tile.Occupied {
                // Verify the unit actually exists at this position
                // Log any mismatches
            }
        }
    }
}
```

### Fix 2: Enhanced Occupancy Management
```go
func (g *Game) movePlayerSafely(player *ecs.Entity, from, to tactical.GridPos) {
    // Verify current position matches expected position
    // Clear old position only if it's actually occupied by this player
    // Set new position with validation
}
```

### Fix 3: Reset Grid State on Tactical Mode Entry
```go
func (g *Game) SwitchToTacticalMode(participants []*ecs.Entity) {
    // Clear all grid occupancy
    g.tacticalManager.Grid.ClearAllOccupancy()
    
    // Redeploy all units with fresh state
    g.tacticalDeployment.DeployParty(partyMembers)
    g.tacticalDeployment.DeployEnemies(enemies)
}
```

## Next Steps for Diagnosis

### 1. Test Specific Scenarios
- Move player from (1,0) to (0,0) 
- Try to move from (0,0) to (0,1)
- Check debug output for occupancy details

### 2. Grid State Inspection
- Add command to dump entire grid state
- Verify which units occupy which tiles
- Check for phantom occupancy

### 3. Unit ID Verification  
- Ensure Entity.GetID() returns consistent values
- Verify deployed units use same ID for movement

## Expected Debug Output Pattern

### Normal Movement (Column 1→2)
```
World pos: (82.0, 120.0) -> Grid pos: (1, 0)
Trying to move from (1, 0) to (2, 0)
Clearing occupancy at (1, 0)
Conan moved to (2, 0)
```

### Broken Movement (Column 0)
```
World pos: (50.0, 120.0) -> Grid pos: (0, 0)  
Trying to move from (0, 0) to (0, 1)
Clearing occupancy at (0, 0)
Cannot move to (0,1) - occupied by unit Conan
```

**Key Question**: Why is (0,1) occupied by the same unit that's trying to move there?

## Status

**Current Phase**: Diagnostic information gathering
**Debug Tools**: Enhanced logging implemented
**Next Action**: Run game with debug output to identify specific failure pattern

The issue is almost certainly related to grid state management rather than coordinate conversion, given that movement TO column 0 works but movement WITHIN column 0 fails.