# Turn Counter and Movement Reset Fixes

## Issues Identified and Fixed

### 1. Turn Counter Not Updating in UI

**Problem**: The combat UI was not displaying the current round number, making it appear that the turn counter was stuck at "turn 1".

**Root Cause**: The UI panel didn't include round information display.

**Solution**:
- Added `GetCurrentRound()` method to `TurnBasedCombatManager`
- Enhanced the combat UI panel to display round information
- Expanded panel height to accommodate additional information

**Files Modified**:
- `internal/tactical/turn_based_combat.go`: Added `GetCurrentRound()` method
- `internal/ui/combat_ui.go`: Updated `drawTurnInfo()` to display round number

**Code Changes**:
```go
// Added to TurnBasedCombatManager
func (cbm *TurnBasedCombatManager) GetCurrentRound() int {
	return cbm.CurrentRound
}

// Updated UI to display round information
roundInfo := fmt.Sprintf("Round: %d", combatManager.GetCurrentRound())
ebitenutil.DebugPrintAt(screen, roundInfo, int(panelX+5), int(panelY+10))
```

### 2. Player Movement Not Reset Properly Between Turns

**Problem**: When a turn ended with "E" (End Turn), player movement range was not being restored for the next turn, preventing proper movement.

**Root Cause**: The `restoreTeamActionPoints()` function was only restoring Action Points and Combat State, but not calling `ResetMovement()` on the RPGStats component.

**Solution**:
- Enhanced `restoreTeamActionPoints()` to also reset the legacy movement system
- Added call to `stats.ResetMovement()` for each team member when their turn starts

**Files Modified**:
- `internal/tactical/turn_based_combat.go`: Updated `restoreTeamActionPoints()`

**Code Changes**:
```go
// Enhanced restoreTeamActionPoints function
func (cbm *TurnBasedCombatManager) restoreTeamActionPoints(team *TeamInfo) {
	for _, member := range team.Members {
		if actionPoints := member.ActionPoints(); actionPoints != nil {
			actionPoints.Restore()
		}
		if combatState := member.CombatState(); combatState != nil {
			combatState.StartTurn()
		}
		// Reset legacy movement system
		if stats := member.RPGStats(); stats != nil {
			stats.ResetMovement()
		}
	}
}
```

### 3. Enhanced UI Information Display

**Additional Improvements**: Enhanced the combat UI to provide better visibility into the combat system state.

**New UI Features**:
- **Round Counter**: Shows current round number
- **Movement Display**: Shows remaining moves from legacy system for debugging
- **Expanded Panel**: Increased panel height to fit more information

**UI Layout Changes**:
```go
// Updated panel layout
panelHeight := float32(110) // Increased from 80

// Added information displays
roundInfo := fmt.Sprintf("Round: %d", combatManager.GetCurrentRound())
moveInfo := fmt.Sprintf("Moves: %d/%d", stats.MovesRemaining, stats.MoveRange)
```

## System Integration

### Turn-Based Combat Flow
The fixes ensure proper integration between the Action Points system and the legacy movement system:

1. **Turn Start**: Both AP and movement are restored
2. **Action Execution**: Both systems are updated when movement occurs
3. **Turn End**: Proper cleanup and state reset for next turn

### Dual Movement Systems
The game currently supports both movement systems:
- **Action Points System**: Primary system using AP costs for movement
- **Legacy Movement System**: Tile-based movement with MoveRange/MovesRemaining
- **Integration**: Both systems are kept in sync for compatibility

### UI Information Hierarchy
The combat information panel now displays:
1. Round number (new)
2. Active unit name
3. Unit class/job
4. Action Points (current/maximum)
5. Movement points (current/maximum) - for debugging
6. Current UI state

## Testing and Validation

### Test Cases Verified
1. **Turn Progression**: Round counter increments correctly
2. **Movement Reset**: Full movement range available each turn
3. **Action Points**: Proper AP restoration between turns
4. **UI Updates**: All information displays update correctly
5. **System Sync**: Both movement systems stay synchronized

### Expected Behavior
- Round counter starts at 1 and increments each full round
- Each unit starts their turn with full AP and movement
- UI shows real-time information about combat state
- End Turn action properly advances to next unit/team

## Compatibility Notes

### Legacy System Support
The fixes maintain backward compatibility with the existing movement system while adding the new Action Points system. This dual approach ensures:
- Existing movement logic continues to work
- New AP-based combat can coexist
- Future migration to pure AP system is possible

### Future Considerations
1. **Movement System Unification**: Eventually consolidate to single system
2. **UI Enhancements**: Add more detailed combat information
3. **Animation Support**: Visual feedback for turn changes
4. **Performance**: Optimize UI updates and state management

## Conclusion

The fixes address both the turn counter display issue and the movement reset problem, ensuring that the turn-based combat system works correctly with proper state management between turns. Players will now see:

- Clear round progression in the UI
- Full movement capabilities restored each turn
- Accurate real-time combat information
- Proper turn-by-turn state management

The enhanced UI provides better visibility into the combat system state, making it easier to understand and debug the tactical combat flow.