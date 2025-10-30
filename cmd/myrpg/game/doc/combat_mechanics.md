# Combat Mechanics Documentation

## Action Points System

### Job Classes and Action Points

| Job Class | Maximum AP | Attacks Per Turn | Movement Per Turn |
|-----------|------------|------------------|-------------------|
| **Warrior** | 4 AP | 2 attacks (4 AP) | 4 tiles |
| **Mage** | 3 AP | 1 attack + 1 move | 3 tiles |
| **Rogue** | 5 AP | 2 attacks + 1 move | 5 tiles |
| **Cleric** | 4 AP | 2 attacks (4 AP) | 4 tiles |
| **Archer** | 4 AP | 2 attacks (4 AP) | 4 tiles |

### Action Costs

| Action Type | AP Cost | Notes |
|-------------|---------|-------|
| **Attack** | 2 AP | Melee or ranged combat |
| **Movement** | 1 AP | Per tile moved |
| **Item Use** | 1 AP | Consumables, equipment |
| **End Turn** | 0 AP | Free action (pressing 'E') |
| **Wait** | 0 AP | Free action |

### Turn System

- **Turn Start**: All units get their full AP restored
- **Action Phase**: Players can act with any unit until AP is exhausted or turn is ended
- **Turn End**: Only when player presses 'E' key or chooses End Turn action
- **Team Turns**: Player team → Enemy team → Repeat

## Combat Example Scenarios

### Rogue (5 AP) Combinations:
- Attack → Attack → Move (2+2+1 = 5 AP)
- Move → Move → Attack → Move (1+1+2+1 = 5 AP) 
- Attack → Move → Move → Move (2+1+1+1 = 5 AP)
- Move → Attack → Move → Move (1+2+1+1 = 5 AP)

### Warrior (4 AP) Combinations:
- Attack → Attack (2+2 = 4 AP)
- Attack → Move → Move (2+1+1 = 4 AP)
- Move → Move → Move → Move (1+1+1+1 = 4 AP)
- Move → Attack → Move (1+2+1 = 4 AP)

### Mage (3 AP) Combinations:
- Attack → Move (2+1 = 3 AP)
- Move → Move → Move (1+1+1 = 3 AP)
- Move → Attack (1+2 = 3 AP)
- Item → Move → Move (1+1+1 = 3 AP)

## Combat Flow

1. **Initialize Combat**
   - All units receive AP based on their job class
   - Initiative order is calculated
   - Grid positions are synchronized

2. **Team Turn**
   - All team members can act until AP is exhausted
   - Multiple units can perform actions in any order
   - Turn only ends when explicitly ended or no units can act

3. **Action Execution**
   - Actions are validated (range, AP cost, etc.)
   - AP is consumed upon successful action
   - Game returns to team turn phase for more actions

4. **Turn End**
   - Triggered by 'E' key or End Turn action
   - Switches to next team
   - AP is restored for incoming team

## Death System

- Units reaching 0 HP are immediately removed from grid
- Dead units are moved off-screen (invisible)
- Grid positions are freed for other units
- Combat continues until victory conditions are met

## Victory Conditions

- **Player Victory**: All enemy units defeated
- **Enemy Victory**: All player units defeated
- **Ongoing**: Both teams have living units

## Controls

| Key | Action | Mode |
|-----|--------|------|
| **A** | Attack adjacent enemy | Tactical |
| **E** | End current team turn | Tactical |
| **Arrow Keys** | Move unit | Tactical |
| **T** | Toggle tactical mode | Overworld |
| **I** | Toggle inventory | Any |

## Technical Constants

Defined in `internal/constants/game.go`:

```go
// Job Class Maximum AP
WarriorMaxAP = 4
MageMaxAP    = 3
RogueMaxAP   = 5
ClericMaxAP  = 4
ArcherMaxAP  = 4

// Action Point Costs
MovementAPCost = 1 // 1 AP per tile moved
AttackAPCost   = 2 // 2 AP per attack action
ItemAPCost     = 1 // 1 AP per item used
EndTurnAPCost  = 0 // Free action
WaitAPCost     = 0 // Free action
```

---
*Last Updated: October 2025*