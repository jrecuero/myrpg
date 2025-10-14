# Combat System Debugging - Round 2

## Issues Identified from Debug Output

Based on the debug messages provided:
```
DEBUG: Moving from (12, 0) to (13, 0)
DEBUG: Current pos (12,0) - Occupied: true, UnitID: 4
DEBUG: Target pos (13,0) - Occupied: true, UnitID: 5
959
DEBUG UI: Found 0 valid attack targets for 2 (AP: 4/4)
```

### Issue Analysis:

1. **Movement is Working**: Units are successfully moving between positions
2. **Grid Occupancy is Correct**: Both positions show proper unit IDs (4 and 5)
3. **Attack Target Detection Failing**: 0 targets found despite adjacent units
4. **Round Counter Stuck**: Still showing Round 1

## Enhanced Debugging Added

### 1. **Action Point Tracking**
Added detailed AP consumption logging:
```go
fmt.Printf("DEBUG: Spending %d AP for %s (before: %d/%d)\n", 
    action.APCost, action.Actor.GetID(), actionPoints.Current, actionPoints.Maximum)
fmt.Printf("DEBUG: After spending AP: %d/%d\n", actionPoints.Current, actionPoints.Maximum)
```

### 2. **Unit Action Capability**
Added `canUnitAct()` debugging:
```go
fmt.Printf("DEBUG: Unit %s can act (AP: %d/%d)\n", 
    entity.GetID(), actionPoints.Current, actionPoints.Maximum)
```

### 3. **Team Turn Progression**
Added comprehensive team turn debugging:
```go
fmt.Printf("DEBUG: No units can act for %s team, ending turn\n", cbm.ActiveTeam.Team.String())
fmt.Printf("DEBUG: Ending turn for %s team (Round %d)\n", cbm.ActiveTeam.Team.String(), cbm.CurrentRound)
```

### 4. **Round Progression Tracking**
Added round advancement debugging:
```go
fmt.Printf("DEBUG: Starting new round: %d -> %d\n", oldRound, cbm.CurrentRound)
fmt.Printf("DEBUG: Looking for next team to act:\n")
fmt.Printf("DEBUG:   Team %s - HasCompleted: %v\n", team.Team.String(), team.HasCompleted)
```

### 5. **Attack Validation Debugging**
Added attack failure reason logging:
```go
fmt.Printf("DEBUG: Attack validation failed for %s -> %s: %v\n", 
    actor.GetID(), member.GetID(), err)
```

### 6. **UI Debug Spam Reduction**
Reduced UI logging to only show changes:
```go
// Only log when target count changes to reduce spam
if len(newTargets) != len(cui.ValidAttackTargets) {
    fmt.Printf("DEBUG UI: Found %d valid attack targets for %s (AP: %d/%d)\n",
        len(newTargets), activeUnit.GetID(), actionPoints.Current, actionPoints.Maximum)
}
```

## Suspected Root Causes

### 1. **Attack Target Issue**
**Hypothesis**: The `validateAttack()` method is rejecting all targets for one of these reasons:
- Units are on the same team (should be different teams)
- Distance calculation is wrong (should be adjacent = distance 1)
- Combat state components missing or invalid
- Units marked as dead when they shouldn't be

**Debug Output**: Will now show exactly why each potential target is rejected

### 2. **Round Progression Issue**
**Hypothesis**: End Turn action isn't properly exhausting AP, so teams never complete
- AP consumption might not be working correctly
- Team completion logic might not trigger
- Round advancement might not be called

**Debug Output**: Will now show AP before/after consumption and team completion status

### 3. **Movement Reset Issue**
**Hypothesis**: Movement reset happens but UI doesn't reflect it immediately
- Movement reset is called but timing issue with UI update
- Legacy movement system vs AP system synchronization issue

**Debug Output**: Will show unit action capability including AP status

## Testing Strategy

### Expected Debug Output for Working System:

1. **End Turn Action**:
```
DEBUG: End turn requested for player_1 (Round: 1)
DEBUG: Executing end turn action (APCost: 4)
DEBUG: Spending 4 AP for player_1 (before: 4/4)
DEBUG: After spending AP: 0/4
DEBUG: Unit player_1 cannot act - AP exhausted (AP: 0/4)
DEBUG: No units can act for Player team, ending turn
DEBUG: Ending turn for Player team (Round 1)
```

2. **Team Turn Progression**:
```
DEBUG: Looking for next team to act:
DEBUG:   Team Player - HasCompleted: true
DEBUG:   Team Enemy - HasCompleted: false
Starting Enemy team turn (Round 1)
```

3. **Round Advancement**:
```
DEBUG: Looking for next team to act:
DEBUG:   Team Player - HasCompleted: true
DEBUG:   Team Enemy - HasCompleted: true
DEBUG: No team can act, starting new round
DEBUG: Starting new round: 1 -> 2
DEBUG: Resetting team Player (was completed: true)
DEBUG: Resetting team Enemy (was completed: true)
Starting Round 2
```

4. **Attack Target Detection**:
```
DEBUG: Attack validation failed for player_1 -> player_2: cannot attack ally
DEBUG: Attack validation succeeded for player_1 -> enemy_1
DEBUG UI: Found 1 valid attack targets for player_1 (AP: 4/4)
```

### Files Modified:
- `internal/tactical/turn_based_combat.go`: Enhanced debugging throughout combat system
- `internal/ui/combat_ui.go`: Reduced debug spam, improved target detection logging

### Next Steps:
1. Run the game and collect debug output
2. Identify which validation is failing for attacks
3. Confirm AP consumption is working for End Turn
4. Verify team completion and round progression logic
5. Remove debug output once issues are resolved

The enhanced debugging should now provide clear visibility into exactly where the combat system is failing and why the round counter isn't advancing properly.