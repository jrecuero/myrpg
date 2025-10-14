# Logging System Implementation for Combat Debugging

## Overview

Implemented a comprehensive file-based logging system to replace console output debugging. This makes it much easier to analyze combat system behavior and debug issues without needing to copy console output.

## Logging System Features

### 1. **File-Based Output**
- Logs automatically saved to `logs/` directory
- Timestamped log files: `myrpg_2025-10-13_15-30-45.log`
- Automatic directory creation if it doesn't exist

### 2. **Multiple Log Levels**
- **DEBUG**: Most verbose, includes detailed system operation
- **INFO**: General information and important events
- **ERROR**: Error conditions and failures

### 3. **Specialized Combat Logging Categories**
- **[COMBAT]**: Attack execution, damage calculation, unit status
- **[TURN]**: Turn management, round progression, team completion
- **[ACTION]**: Action point consumption, action execution
- **[UI]**: User interface state changes and interactions

### 4. **Detailed Log Format**
```
[DEBUG] 2025/10/13 15:30:45.123456 combat_ui.go:192: [UI] Found 2 valid attack targets for player_1 (AP: 4/4)
[DEBUG] 2025/10/13 15:30:45.123789 turn_based_combat.go:445: [TURN] Unit player_1 can act (AP: 4/4)
[DEBUG] 2025/10/13 15:30:45.124012 tactical_manager.go:115: [ACTION] End turn requested for player_1 (Round: 1)
```

## Log Categories and Usage

### Combat Logs (`[COMBAT]`)
- Attack validation results
- Damage calculation and application
- Unit death/defeat notifications
- Attack success/failure reasons

```go
logger.Combat("Attack validation failed for %s -> %s: %v", actor.GetID(), member.GetID(), err)
logger.Combat("%s attacks %s", attackerStats.Name, targetStats.Name)
logger.Combat("%s deals %d damage to %s (HP: %d/%d)", attackerStats.Name, damage, targetStats.Name, targetStats.CurrentHP, targetStats.MaxHP)
```

### Turn Management Logs (`[TURN]`)
- Round progression tracking
- Team turn completion
- Unit action capability status
- Team reset operations

```go
logger.Turn("Starting new round: %d -> %d", oldRound, cbm.CurrentRound)
logger.Turn("Looking for next team to act:")
logger.Turn("Unit %s can act (AP: %d/%d)", entity.GetID(), actionPoints.Current, actionPoints.Maximum)
logger.Turn("No units can act for %s team, ending turn", cbm.ActiveTeam.Team.String())
```

### Action Logs (`[ACTION]`)
- Action point consumption tracking
- Action creation and execution
- End turn processing

```go
logger.Action("End turn requested for %s (Round: %d)", activeUnit.GetID(), tm.TurnBasedCombat.GetCurrentRound())
logger.Action("Spending %d AP for %s (before: %d/%d)", action.APCost, action.Actor.GetID(), actionPoints.Current, actionPoints.Maximum)
logger.Action("After spending AP: %d/%d", actionPoints.Current, actionPoints.Maximum)
```

### UI Logs (`[UI]`)
- Target detection changes
- UI state transitions
- Button enable/disable status

```go
logger.UI("Found %d valid attack targets for %s (AP: %d/%d)", len(newTargets), activeUnit.GetID(), actionPoints.Current, actionPoints.Maximum)
```

## Files Modified

### New Files:
- `internal/logger/logger.go`: Complete logging system implementation

### Modified Files:
- `cmd/myrpg/main.go`: Logger initialization and cleanup
- `internal/tactical/turn_based_combat.go`: Replaced all debug prints with appropriate logging calls
- `internal/engine/tactical_manager.go`: Updated action callbacks to use logging
- `internal/ui/combat_ui.go`: UI-specific logging with spam reduction

## Usage Instructions

### 1. **Running the Game with Logging**
```bash
cd /Users/jorecuer/go/src/github.com/jrecuero/myrpg
go run ./cmd/myrpg
```

### 2. **Finding the Log Files**
- Logs are automatically created in the `logs/` directory
- File name includes timestamp: `myrpg_2025-10-13_15-30-45.log`
- Console output shows the log file name when the game starts

### 3. **Reading Log Files**
```bash
# View the latest log file
tail -f logs/myrpg_*.log

# Search for specific categories
grep "\[COMBAT\]" logs/myrpg_*.log
grep "\[TURN\]" logs/myrpg_*.log
grep "\[ACTION\]" logs/myrpg_*.log
```

## Expected Log Output for Issues

### Issue 1: Round Counter Not Updating
Look for these log patterns:
```
[ACTION] End turn requested for player_1 (Round: 1)
[ACTION] Spending 4 AP for player_1 (before: 4/4)
[ACTION] After spending AP: 0/4
[TURN] Unit player_1 cannot act - AP exhausted
[TURN] No units can act for Player team, ending turn
[TURN] Ending turn for Player team (Round 1)
[TURN] Looking for next team to act:
[TURN]   Team Player - HasCompleted: true
[TURN]   Team Enemy - HasCompleted: false
```

If round doesn't advance, we should see:
```
[TURN] No team can act, starting new round
[TURN] Starting new round: 1 -> 2
[TURN] Resetting team Player (was completed: true)
[TURN] Resetting team Enemy (was completed: true)
```

### Issue 2: Attack Target Detection
Look for these patterns:
```
[UI] Found 0 valid attack targets for player_1 (AP: 4/4)
[COMBAT] Attack validation failed for player_1 -> enemy_1: target out of range (distance: 3, max range: 1)
[COMBAT] Attack validation failed for player_1 -> player_2: cannot attack ally
```

### Issue 3: Movement Range Updates
Look for turn transitions:
```
[TURN] Starting Player team turn (Round 1)
[TURN] Unit player_1 can act (AP: 4/4)
```

## Debugging Workflow

### 1. **Reproduce the Issue**
- Start the game
- Note the log file name from console output
- Perform the actions that trigger the bug

### 2. **Analyze the Log File**
```bash
# View real-time logs
tail -f logs/myrpg_*.log

# Filter by category
grep "\[TURN\]" logs/myrpg_*.log | tail -20
```

### 3. **Share Log Contents**
- Copy relevant sections from the log file
- Much easier than copying console output
- Includes precise timestamps and source code locations

## Benefits

1. **Persistent Records**: All debug information saved to files
2. **Easy Analysis**: Can grep, filter, and search log contents
3. **No Console Spam**: Clean console output with detailed file logs
4. **Timestamped**: Precise timing information for debugging
5. **Categorized**: Easy to focus on specific subsystems
6. **Professional**: Production-ready logging system

## Next Steps

1. **Test the Issues**: Run the game and check the log files for the reported problems
2. **Analyze Patterns**: Use the categorized logs to identify root causes
3. **Clean Up**: Once debugging is complete, can reduce log verbosity if needed
4. **Extend**: Add more specialized logging categories as needed

The logging system will now provide comprehensive insight into exactly what's happening in the combat system, making it much easier to identify and fix the remaining issues!