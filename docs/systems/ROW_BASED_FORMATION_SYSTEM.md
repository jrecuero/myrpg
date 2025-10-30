# Row-Based Formation System Enhancement

## Overview
Enhance the classic Dragon Quest-style battle system with strategic front/back row mechanics and attack reach limitations to add tactical depth while maintaining the classic JRPG feel.

## Current State
- **Player Formation**: 2x2 grid at bottom of screen (green area)  
- **Enemy Formation**: 1-2 rows at top of screen (red area)
- **Attack System**: All attacks target any available enemy/player
- **No positional strategy**: Position doesn't affect combat mechanics

## Enhanced Row System Design

### Row Layout Reorganization

#### Enemy Formation (Top Area)
```
┌─────────────────────────┐
│     BACK ROW (Top)      │ ← Mages, Archers, Support
│    🧙 🏹 🛡️ 🧙         │
├─────────────────────────┤
│    FRONT ROW (Bottom)   │ ← Warriors, Tanks, Melee
│    ⚔️ 🗡️ ⚔️ 🗡️       │
└─────────────────────────┘
```

#### Player Formation (Bottom Area)  
```
┌─────────────────────────┐
│    FRONT ROW (Top)      │ ← Tanks, Warriors, Melee
│    🛡️ ⚔️ 🗡️ ⚔️       │
├─────────────────────────┤
│    BACK ROW (Bottom)    │ ← Mages, Healers, Archers
│    🧙 💚 🏹 🧙         │
└─────────────────────────┘
```

### Attack Reach System

#### Close Combat Attacks (Melee Range)
- **Normal Rule**: Can only attack front row enemies
- **Exception Rule**: If front row is empty, can attack back row
- **Examples**: Sword Strike, Hammer Blow, Claw Attack

#### Ranged Combat Attacks (Unlimited Range)  
- **Rule**: Can attack any row at any time
- **Examples**: Bow Shot, Magic Missile, Lightning Bolt, Heal

#### Multi-Target Attacks (Area Effects)
- **Front Row Cleave**: Hits all front row enemies
- **Back Row Sweep**: Hits all back row enemies  
- **Full Formation**: Hits entire enemy formation (rare/powerful)

### Formation Mechanics

#### Row Switching
- **Mid-Battle**: Allow formation changes as an action (costs AP/turn)
- **Automatic**: When front row member dies, back row can advance
- **Strategic**: Players choose when to reposition

#### Row Bonuses
- **Front Row**: +20% physical damage dealt, +20% physical damage received  
- **Back Row**: +20% magic damage dealt, -20% physical damage received

#### Protection Mechanics
- **Covering**: Front row provides partial protection to back row
- **Flanking**: Empty front row exposes back row to full damage
- **Formation Integrity**: Maintaining formation provides team bonuses

## Technical Implementation

### Data Structures

#### Row Position Enum
```go
type RowPosition int

const (
    RowFront RowPosition = iota
    RowBack
)
```

#### Attack Reach Enum
```go
type AttackReach int

const (
    ReachMelee   AttackReach = iota // Front row only (with exception)
    ReachRanged                     // Any row
    ReachFrontAOE                   // All front row
    ReachBackAOE                    // All back row  
    ReachAllAOE                     // Full formation
)
```

#### Enhanced Formation Structure
```go
type EnhancedFormation struct {
    FrontRow []*ecs.Entity  // Front line combatants
    BackRow  []*ecs.Entity  // Support/ranged units
    
    // Position mapping for UI rendering
    FrontPositions []FormationPosition
    BackPositions  []FormationPosition
}

type FormationPosition struct {
    X, Y     float64       // Screen coordinates  
    Occupied bool         // Has entity at this position
    Entity   *ecs.Entity  // Entity reference (nil if empty)
}
```

#### Attack Validation
```go
func ValidateAttackTarget(attacker, target *ecs.Entity, attackReach AttackReach) bool {
    attackerRow := GetEntityRow(attacker)
    targetRow := GetEntityRow(target)
    
    switch attackReach {
    case ReachMelee:
        // Can attack front row, or back row if front is empty
        if targetRow == RowFront {
            return true
        }
        if targetRow == RowBack {
            return IsRowEmpty(GetOpposingFormation(attacker), RowFront)
        }
        return false
        
    case ReachRanged:
        return true // Can always attack any row
        
    default:
        return false
    }
}
```

### UI Changes

#### Formation Display
- **Visual Separation**: Clear row divisions with different background tints
- **Position Indicators**: Show available positions vs occupied positions  
- **Row Labels**: "Front" and "Back" indicators for clarity

#### Target Selection
- **Reach Highlighting**: Highlight valid targets based on attack reach
- **Invalid Targets**: Gray out unreachable targets with explanation
- **Row Status**: Show if front row is empty (enabling back row attacks)

#### Formation Management
- **Drag & Drop**: Allow repositioning entities between rows
- **Formation Presets**: Quick formation setups (Defensive, Offensive, Balanced)
- **Auto-Arrange**: Intelligent formation suggestions based on party composition

## Combat Flow Integration

### Battle Start
1. **Formation Setup**: Players choose initial row positions
2. **Enemy Formation**: AI determines optimal enemy positioning  
3. **Display Update**: Render formations with clear row separation

### Action Selection  
1. **Attack Choice**: Player selects attack type (melee/ranged/special)
2. **Reach Validation**: System highlights valid targets based on reach
3. **Target Selection**: Player chooses from available targets
4. **Execution**: Attack resolves with row-based modifiers

### Dynamic Changes
1. **Death Handling**: Remove entity, check if row becomes empty
2. **Row Exposure**: Update targeting rules if front row emptied
3. **Formation Shift**: Allow automatic or manual repositioning
4. **Visual Update**: Reflect changes in formation display

## Balance Considerations

### Melee vs Ranged Balance
- **Melee Advantage**: Higher base damage, critical hit chances
- **Ranged Advantage**: Targeting flexibility, safety from counter-attacks
- **Row Bonuses**: Offset positioning disadvantages

### Formation Strategy
- **Tank Front**: High HP/Defense units protect fragile back row
- **DPS Front**: High damage units for aggressive playstyle  
- **Mixed Formation**: Balanced approach with moderate protection

### AI Adaptation
- **Smart Positioning**: AI places tanks front, mages back
- **Target Priority**: Focus on exposed back row when possible
- **Formation Breaking**: AI attempts to eliminate front row protection

## Implementation Priority

### Phase 1: Core Mechanics (1 week)
- [ ] Row position tracking and management
- [ ] Basic reach validation (melee vs ranged)
- [ ] Formation display reorganization  
- [ ] Target highlighting system

### Phase 2: Advanced Features (1 week)
- [ ] Row bonuses and modifiers
- [ ] Formation change actions
- [ ] AI formation intelligence
- [ ] Multi-target attack types

### Phase 3: Polish & Balance (3-5 days)
- [ ] Visual enhancements and animations
- [ ] Formation presets and auto-arrange
- [ ] Balance testing and tuning
- [ ] Tutorial integration

## Testing Scenarios

### Functional Tests
- [ ] Melee attacks blocked by empty front row
- [ ] Ranged attacks always succeed  
- [ ] Front row elimination exposes back row
- [ ] Formation changes work correctly

### Balance Tests  
- [ ] Tank-heavy front row effectiveness
- [ ] All-ranged formation viability
- [ ] Mixed formation balance
- [ ] AI formation intelligence

### Edge Cases
- [ ] Single entity formations
- [ ] All entities in one row
- [ ] Mid-battle formation shifts
- [ ] Row bonus stacking

## Success Metrics

### Player Engagement
- **Formation Usage**: Players actively manage row positions
- **Strategic Depth**: Different formations for different encounters  
- **Learning Curve**: New players understand mechanics quickly

### Combat Depth
- **Tactical Variety**: Multiple viable formation strategies
- **Decision Complexity**: Meaningful choices in positioning
- **Replay Value**: Different approaches to same encounters

---

*Last Updated: October 30, 2025*  
*Status: Design Phase - Ready for Implementation*