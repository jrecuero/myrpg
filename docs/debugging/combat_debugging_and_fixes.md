# Combat System Issues - Debugging and Fixes

## Issues Addressed

### 1. **Round Counter Not Updating with "E" (End Turn)**

**Problem**: Round counter stayed at "Round 1" even after pressing "E" to end turn.

**Root Cause**: The `CreateEndTurnAction()` method was setting `APCost: 0`, meaning the unit wasn't actually exhausting all AP, so team turns weren't properly ending.

**Fix**: Modified `CreateEndTurnAction()` to consume all remaining Action Points:
```go
// Get current AP to consume all remaining
actionPoints := actor.ActionPoints()
apCost := 0
if actionPoints != nil {
    apCost = actionPoints.Current // Consume all remaining AP
}

action := &CombatAction{
    Type:      ActionWait,
    Actor:     actor,
    Target:    nil,
    TargetPos: GridPos{},
    APCost:    apCost, // Consume all remaining AP
    Validated: true,
    Message:   fmt.Sprintf("%s ends turn", cbm.getEntityName(actor)),
}
```

### 2. **Movement Range Not Updating After Turn End**

**Problem**: When turn ended with "E", movement range wasn't properly displayed until first movement.

**Root Cause**: UI updates valid moves but there was a timing issue with movement reset and UI refresh.

**Status**: The movement reset logic was already correct (`restoreTeamActionPoints()` calls `stats.ResetMovement()`), but added debugging to verify proper functioning.

### 3. **Attack System Issues**

**Problem**: Attack with "A" showed no feedback, no visual indication of attackable targets.

**Multiple Fixes Applied**:

#### A. Enhanced Attack Feedback
- Added detailed console logging for attack execution
- Added visual feedback in game messages
- Added debug output to track attack flow

```go
// Enhanced executeAttack with better feedback
cbm.logMessage(fmt.Sprintf("*** %s deals %d damage to %s (HP: %d/%d) ***",
    attackerStats.Name, damage, targetStats.Name,
    targetStats.CurrentHP, targetStats.MaxHP))

fmt.Printf("DEBUG: %s deals %d damage to %s (HP: %d/%d)\n",
    attackerStats.Name, damage, targetStats.Name,
    targetStats.CurrentHP, targetStats.MaxHP)
```

#### B. Improved Target Highlighting
- Enhanced coordinate conversion for attack target highlighting
- Added target count display in UI instructions
- Added visual target indicators (T1, T2, etc.)

```go
// Improved attack target highlighting
instruction := fmt.Sprintf("Click on a red tile to attack (%d targets), ESC to cancel", len(cui.ValidAttackTargets))
```

#### C. Added Comprehensive Debug Logging
- UI debug: Shows number of valid attack targets found
- Action debug: Tracks attack action creation and execution
- Combat debug: Logs attack resolution and damage

## Debugging Infrastructure Added

### 1. **Turn Management Debugging**
```go
fmt.Printf("DEBUG: End turn requested for %s (Round: %d)\n", 
    activeUnit.GetID(), tm.TurnBasedCombat.GetCurrentRound())
fmt.Printf("DEBUG: End turn executed. New Round: %d\n", 
    tm.TurnBasedCombat.GetCurrentRound())
```

### 2. **Attack System Debugging**
```go
fmt.Printf("DEBUG UI: Found %d valid attack targets for %s (AP: %d/%d)\n", 
    len(cui.ValidAttackTargets), activeUnit.GetID(), actionPoints.Current, actionPoints.Maximum)
fmt.Printf("DEBUG: Creating attack action from %s to %s\n", 
    activeUnit.GetID(), target.GetID())
```

### 3. **Combat Resolution Debugging**
```go
fmt.Printf("DEBUG: %s attacks %s\n", attackerStats.Name, targetStats.Name)
fmt.Printf("DEBUG: %s deals %d damage to %s (HP: %d/%d)\n",
    attackerStats.Name, damage, targetStats.Name,
    targetStats.CurrentHP, targetStats.MaxHP)
```

## Visual Improvements

### 1. **Enhanced Attack Target Display**
- Red highlighting for valid attack targets
- Target numbering (T1, T2, T3, etc.) for identification
- Target count in instruction text
- Improved coordinate conversion for accurate positioning

### 2. **Better UI Feedback**
- Enhanced game messages with `***` markers for important combat events
- Real-time target count display
- Clear instruction text for each interaction mode

## System Integration

### 1. **Action Point System**
- End Turn now properly consumes all remaining AP
- Ensures team turns end when all units are exhausted
- Maintains synchronization between AP and legacy movement systems

### 2. **Round Progression**
- Proper team turn completion triggers round advancement
- Round counter updates correctly in UI
- New rounds restore full AP and movement for all units

### 3. **Combat Flow**
```
Player Action → AP Consumption → Team Check → Round Progression
```

## Testing and Validation

### Expected Behavior After Fixes:

1. **Round Counter**: 
   - Starts at "Round: 1"
   - Increments when all teams complete their turns
   - Updates visible in UI panel

2. **Movement System**:
   - Full movement restored each turn
   - Movement counter shows correct values
   - Both AP and legacy systems synchronized

3. **Attack System**:
   - Red tiles highlight adjacent enemies
   - Target count shown in instructions
   - Clear feedback when attacks execute
   - Damage and HP changes logged to console and game messages

### Debug Output Guide:

When running the game, you should see console output like:
```
DEBUG UI: Found 2 valid attack targets for player_1 (AP: 4/4)
DEBUG: Creating attack action from player_1 to enemy_1
DEBUG: Executing attack action
DEBUG: player_1 attacks enemy_1
DEBUG: player_1 deals 3 damage to enemy_1 (HP: 7/10)
```

## Files Modified

1. **`internal/tactical/turn_based_combat.go`**:
   - Fixed `CreateEndTurnAction()` to consume all AP
   - Enhanced `executeAttack()` with better feedback
   - Added debug logging throughout combat system

2. **`internal/engine/tactical_manager.go`**:
   - Added debug logging to UI callbacks
   - Enhanced attack and end turn action handling

3. **`internal/ui/combat_ui.go`**:
   - Improved attack target highlighting
   - Added debug logging for UI updates
   - Enhanced visual feedback and instructions

## Next Steps

1. **Remove Debug Output**: Once testing confirms fixes work properly, remove debug `fmt.Printf()` statements
2. **Enhanced Visual Effects**: Consider adding animations for attacks
3. **Sound Effects**: Add audio feedback for combat actions
4. **UI Polish**: Replace debug text with proper font rendering

The system now provides comprehensive feedback for all combat actions and should resolve the reported issues with round progression, movement updates, and attack system functionality.