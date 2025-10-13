# MyRPG Battle System Documentation

## Overview
The MyRPG battle system is a turn-based tactical combat system built on a grid-based battlefield. Combat emphasizes strategic positioning, resource management through action points, and team coordination.

## Core Design Principles
- **Turn-Based**: Combat progresses in discrete turns with clear phases
- **Team-Based**: All team members act within a single team turn
- **Grid-Based**: All positioning and movement occurs on a tactical grid
- **Resource Management**: Action points limit what can be accomplished per turn
- **Scalable**: Designed to accommodate future enhancements (equipment, status effects, etc.)

---

## Battle Flow

### 1. Combat Initialization
```
┌─────────────────────────────────────────────┐
│ 1. Deploy teams on tactical grid            │
│ 2. Calculate team initiative (speed totals) │
│ 3. Determine first team (highest initiative)│
│ 4. Initialize action points for all units   │
│ 5. Begin first team's turn                  │
└─────────────────────────────────────────────┘
```

### 2. Team Turn Structure
```
Team Turn Phase:
├── Action Point Allocation
├── Unit Action Selection (any order)
│   ├── Movement
│   ├── Attack
│   ├── Item Usage
│   └── End Turn
├── Action Resolution
└── Turn End Check
```

### 3. Combat Resolution
Combat continues alternating between teams until:
- All enemy units are defeated (player victory)
- All player units are defeated (player defeat)
- Special win/loss conditions are met

---

## Initiative System

### Team Initiative Calculation
```go
teamInitiative = sum(unit.speed for unit in team)
```

### Turn Order Rules
1. **Higher Initiative**: Team with higher total speed acts first
2. **Tie Breaking**: Player team acts first in case of ties
3. **Turn Alternation**: Teams alternate turns throughout combat
4. **Initiative Recalculation**: Initiative recalculated if team composition changes

### Example
```
Player Team: Warrior(5) + Mage(3) + Rogue(7) = 15 initiative
Enemy Team:  Orc(4) + Goblin(6) + Troll(2) = 12 initiative
Result: Player team acts first
```

---

## Action Point System

### Action Point Allocation
- Each unit receives action points at the start of their team's turn
- Action points determine what actions a unit can perform
- Unused action points do not carry over to the next turn

### Action Point Costs
| Action Type | AP Cost | Notes |
|-------------|---------|-------|
| Movement | 1 AP per tile | Moving to adjacent tile |
| Basic Attack | 2 AP | Melee attack on adjacent enemy |
| Item Usage | 1 AP | Using consumable items |
| End Turn | 0 AP | Voluntarily end unit's actions |

### Default Action Point Values
| Unit Class | Base AP | Movement Range | Special Notes |
|------------|---------|----------------|---------------|
| Warrior | 4 AP | 3 tiles | Balanced offense/mobility |
| Mage | 3 AP | 2 tiles | Lower mobility, spell focus |
| Rogue | 5 AP | 5 tiles | High mobility, hit-and-run |
| Cleric | 4 AP | 3 tiles | Support and healing focus |
| Archer | 4 AP | 4 tiles | Ranged combat specialist |

---

## Combat Actions

### Movement
- **Cost**: 1 AP per tile moved
- **Range**: Based on unit class and remaining AP
- **Restrictions**: Cannot move through occupied tiles or obstacles
- **Pathfinding**: Currently requires manual tile selection

### Attack Actions
- **Cost**: 2 AP per attack
- **Range**: Adjacent tiles only (range = 1)
- **Target Selection**: Required when multiple valid targets available
- **Damage Calculation**: Based on unit base attack stats
- **Critical Hits**: Not implemented (future enhancement)

### Item Usage
- **Cost**: 1 AP per item used
- **Types**: Currently not implemented (future enhancement)
- **Restrictions**: Must have item in inventory

---

## Positioning and Targeting

### Grid System
- **Size**: 20×10 tactical grid (configurable via constants)
- **Tile Size**: 32×32 pixels per tile
- **Coordinates**: Standard X,Y coordinate system
- **Boundaries**: Movement restricted to valid grid areas

### Attack Targeting
- **Range**: 1 tile (adjacent only)
- **Line of Sight**: Not implemented (all adjacent tiles valid)
- **Area of Effect**: Not implemented (single target only)
- **Friendly Fire**: Not applicable (cannot target allies)

### Target Selection UI
When multiple enemies are adjacent:
1. Highlight all valid targets
2. Player selects specific target
3. Attack resolves against selected target

---

## Unit Statistics

### Core Stats (Current Implementation)
```go
type UnitStats struct {
    Health    int     // Hit points (HP)
    Attack    int     // Base attack damage
    Defense   int     // Damage reduction (future)
    Speed     int     // Initiative contribution
    MaxAP     int     // Action points per turn
}
```

### Health and Death
- **Death Condition**: HP ≤ 0
- **Death Resolution**: Unit immediately removed from tactical grid
- **Resurrection**: Not implemented (future enhancement)
- **Damage Types**: Generic damage only (no elemental types)

### Stat Modifiers
- **Equipment**: Not implemented (future enhancement)
- **Status Effects**: Not implemented (future enhancement)
- **Temporary Buffs**: Not implemented (future enhancement)

---

## Enemy AI Behavior

### Current AI Implementation
**Phase**: Basic Static AI
- **Movement**: Enemies do not move
- **Targeting**: Attack any adjacent player unit
- **Decision Making**: Simple proximity-based attacks
- **Coordination**: No team coordination

### AI Action Priority
1. **Attack**: If player unit adjacent → attack
2. **Wait**: If no valid targets → end turn

### Future AI Enhancements (Planned)
- Pathfinding and movement toward players
- Target prioritization (low health, high threat)
- Ability usage and resource management
- Team coordination and positioning

---

## Technical Architecture

### Key Components
```go
// Combat state management
type CombatManager struct {
    CurrentTeam     Team
    TurnNumber      int
    InitiativeOrder []Team
    CombatPhase     CombatPhase
}

// Action system
type ActionType int
const (
    ActionMove ActionType = iota
    ActionAttack
    ActionItem
    ActionEndTurn
)

// Unit action points
type ActionPoints struct {
    Current int
    Maximum int
}
```

### Integration Points
- **ECS System**: Combat components integrated with entity system
- **Grid Manager**: Positioning and movement validation
- **UI System**: Action selection and feedback
- **Animation System**: Combat visual effects (future)

---

## Future Enhancements

### Phase 2: Enhanced Combat
- [ ] **Attack Ranges**: Different ranges per weapon/spell type
- [ ] **Equipment System**: Weapons and armor affecting stats
- [ ] **Status Effects**: Poison, stun, buffs, debuffs
- [ ] **Spell System**: Magical abilities with varied effects

### Phase 3: Advanced Tactics
- [ ] **Line of Sight**: Vision and concealment mechanics
- [ ] **Terrain Effects**: Movement costs and cover
- [ ] **Area Attacks**: Multi-tile damage abilities
- [ ] **Opportunity Attacks**: Reaction-based combat

### Phase 4: Strategic Depth
- [ ] **Elemental System**: Damage types and resistances
- [ ] **Critical Hits**: Chance-based damage bonuses
- [ ] **Formation Bonuses**: Positioning-based team benefits
- [ ] **Environmental Hazards**: Interactive battlefield elements

---

## Configuration Constants

### Combat Constants (internal/constants/game.go)
```go
// Action Point defaults
const (
    WarriorMaxAP = 4
    MageMaxAP    = 3
    RogueMaxAP   = 5
    ClericMaxAP  = 4
    ArcherMaxAP  = 4
)

// Action costs
const (
    MovementAPCost = 1
    AttackAPCost   = 2
    ItemAPCost     = 1
)

// Combat ranges
const (
    MeleeRange = 1  // Adjacent tiles only
)
```

### Display Constants (internal/constants/display.go)
```go
// Grid configuration
const (
    GridWidth  = 20
    GridHeight = 10
    TileSize   = 32
)
```

---

## Testing and Validation

### Unit Test Coverage
- [ ] Initiative calculation
- [ ] Action point management
- [ ] Movement validation
- [ ] Attack resolution
- [ ] Death handling

### Integration Test Scenarios
- [ ] Full combat scenario (player victory)
- [ ] Full combat scenario (player defeat)
- [ ] Edge cases (single units, no AP)
- [ ] UI interaction testing

### Performance Considerations
- Action calculation optimization
- Grid update efficiency
- Memory management for large battles

---

## Changelog

### Version 1.0 (Initial Implementation)
- Basic turn-based combat system
- Team initiative calculation
- Action point economy
- Adjacent-tile combat only
- Simple enemy AI (static, attack-only)
- Unit death and removal

### Future Versions
- Version 1.1: Attack ranges and equipment
- Version 1.2: Status effects and advanced AI
- Version 2.0: Full tactical enhancement suite

---

*Last Updated: October 13, 2025*
*Status: Design Phase - Ready for Implementation*