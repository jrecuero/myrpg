# System Reference

Complete reference documentation for all systems in MyRPG's ECS architecture.

## Overview

Systems contain the logic that operates on entities with specific components. They implement game behavior by querying entities and updating their state. This document catalogs all systems, their purpose, and how they work.

---

## Core Concept

**Systems follow these principles:**
1. **Query entities** with required components
2. **Process each entity** that matches criteria
3. **Update component data** based on game logic
4. **Never create entities** (only modify existing ones)
5. **Single responsibility** (one system, one job)

**Basic System Pattern:**
```go
type MySystem struct {
    world *ecs.World
}

func (s *MySystem) Update(deltaTime float64) {
    for _, entity := range s.world.GetEntities() {
        // Check for required component
        component := entity.MyComponent()
        if component == nil {
            continue
        }
        
        // Update component state
        component.Value += deltaTime
    }
}
```

---

## Battle Systems

### BattleSystem

**File:** `internal/engine/battle.go`  
**Purpose:** Manages turn-based combat flow between players and enemies

**Responsibilities:**
- Turn management (player turn, enemy turn)
- Attack resolution (physical and magical)
- Battle state transitions
- Attack animations
- Combat UI coordination

**Key Methods:**
```go
StartBattle(player, enemy)       // Initialize combat
Update(deltaTime)                 // Process battle logic
HandleInput()                     // Process player combat inputs
ExecuteAttack(attacker, target)   // Resolve attacks
CheckBattleEnd()                  // Determine victory/defeat
```

**Battle States:**
- `BattleStateNone`: No active battle
- `BattleStatePlayerTurn`: Awaiting player input
- `BattleStateEnemyTurn`: Processing enemy action
- `BattleStateBattleEnd`: Battle concluded

**Example:**
```go
battleSystem := engine.NewBattleSystem()
battleSystem.SetMessageCallback(func(msg string) {
    ui.DisplayMessage(msg)
})
battleSystem.StartBattle(player, enemy)
```

**Components Used:**
- `RPGStats`: Character stats and health
- `Animation`: Attack animations
- `CombatState`: Battle-specific state

---

### Classic Battle System

**File:** `internal/battle/classic/battle_manager.go`  
**Purpose:** Dragon Quest-style battle system with formations

**Responsibilities:**
- Party vs. enemy group battles
- Formation positioning (2x2 grid)
- Activity-based turn system (speed determines turn order)
- Target selection and navigation
- Battle actions (attack, defend, magic, item)

**Key Features:**
- **Speed-Based Turns**: Entities act based on speed stat
- **Formation System**: Front row and back row positioning
- **Target Navigation**: Arrow key navigation between targets
- **Action Queue**: Actions execute based on speed

**Key Methods:**
```go
StartBattle(playerParty, enemyParty)  // Initialize battle
Update(deltaTime)                      // Process activity queue
HandleActionSelect(action)             // Player selects action
HandleTargetNavigation(direction)      // Navigate targets
ConfirmTargetSelection()               // Execute action
```

**Battle Flow:**
1. Initialize formations and activity queue
2. Process entity turns based on speed
3. Player selects action and target
4. Execute action and apply effects
5. Check victory/defeat conditions

**Example:**
```go
battleMgr := classic.NewBattleManager()

playerParty := []*ecs.Entity{warrior, mage, cleric}
enemyParty := []*ecs.Entity{goblin1, goblin2, goblin3}

battleMgr.StartBattle(playerParty, enemyParty)
```

**Components Used:**
- `RPGStats`: Combat stats and HP/MP
- `Transform`: Formation positioning
- `Sprite`: Visual representation
- `CombatState`: Defending state

---

### Tactical Battle System

**File:** `internal/battle/tactical/turn_based_combat.go`  
**Purpose:** Grid-based tactical combat system

**Responsibilities:**
- Grid movement and positioning
- Action point management
- Turn-based tactical combat
- Line of sight and range calculations
- Movement undo functionality

**Key Features:**
- **Grid-Based Movement**: Entities move on tile grid
- **Action Points**: Limited actions per turn
- **Tactical Positioning**: Range and positioning matters
- **Undo System**: Revert movement decisions

**Key Methods:**
```go
StartCombat(entities, gridWidth, gridHeight)
CreateMoveAction(actor, targetPos)
CreateAttackAction(actor, target)
ExecuteAction(action)
UndoLastMove(actor)
EndTurn(actor)
```

**Action Types:**
- `ActionMove`: Move to grid position
- `ActionAttack`: Attack target entity
- `ActionSkill`: Use special ability
- `ActionEndTurn`: End current turn

**Example:**
```go
tacMgr := tactical.NewTurnBasedCombatManager()
tacMgr.StartCombat(combatants, 12, 8)

// Player moves
action := tacMgr.CreateMoveAction(player, GridPos{X: 5, Y: 3})
tacMgr.ExecuteAction(action)

// Player attacks
attackAction := tacMgr.CreateAttackAction(player, enemy)
tacMgr.ExecuteAction(attackAction)
```

**Components Used:**
- `RPGStats`: Stats, HP, movement range
- `ActionPoints`: AP costs and availability
- `CombatState`: Turn state
- `Transform`: Grid positioning

---

## Item & Inventory Systems

### ConsumableManager

**File:** `internal/systems/consumable_system.go`  
**Purpose:** Handles consumable item usage and effects

**Responsibilities:**
- Apply consumable effects to targets
- HP/MP restoration
- Stat buffs
- Status effect curing

**Effect Types:**
- `heal_hp`: Restore HP
- `heal_mp`: Restore MP
- `buff_attack`: Increase attack
- `buff_defense`: Increase defense
- `buff_magic_attack`: Increase magic attack
- `buff_magic_defense`: Increase magic defense
- `buff_speed`: Increase speed
- `cure_all`: Remove status effects

**Example:**
```go
consumableMgr := systems.NewConsumableManager()

// Use potion on player
err := consumableMgr.UseConsumable(potionItem, player, player)
if err != nil {
    logger.Error("Failed to use potion:", err)
}
```

**Components Used:**
- `RPGStats`: Target stats to modify
- `Inventory`: Item storage
- `Item`: Item properties and effects

---

## View & State Management

### ViewManager

**File:** `internal/engine/view_system.go`  
**Purpose:** Manages game view states and transitions

**Responsibilities:**
- View state management (Exploration, Battle, Menu, etc.)
- View transitions with conditions
- View-specific entity management
- Event filtering by view
- Transition data passing

**View Types:**
```go
ViewExploration  // Overworld exploration
ViewBattle       // Combat encounters
ViewMenu         // Pause/main menu
ViewDialog       // NPC conversations
ViewShop         // Merchant interface
ViewInventory    // Inventory management
ViewTactical     // Tactical battle mode
```

**Key Methods:**
```go
SetCurrentView(viewType)
RegisterTransition(from, to, condition, priority)
Update(deltaTime)
IsEventActive(eventComp)
AddViewEntity(viewType, entity)
GetTransitionData(key)
```

**Transition System:**
```go
// Register transition from Exploration to Battle
viewMgr.RegisterTransition(
    engine.ViewExploration,
    engine.ViewBattle,
    func() bool { return battleTriggered },
    1, // priority
)
```

**Example:**
```go
viewMgr := engine.NewViewManager()
viewMgr.SetCurrentView(engine.ViewExploration)

// Add exploration-only entities
viewMgr.AddViewEntity(engine.ViewExploration, npcEntity)

// Transition to battle
viewMgr.SetCurrentView(engine.ViewBattle)
```

**Components Used:**
- `Event`: View-specific event filtering
- All components (manages which are active per view)

---

### BattleSystemSelector

**File:** `internal/engine/battle_system_selector.go`  
**Purpose:** Routes to appropriate battle system based on configuration

**Responsibilities:**
- Select between Classic and Tactical battle modes
- Initialize correct battle system
- Handle battle system transitions
- Coordinate with ViewManager

**Battle Modes:**
- `BattleModeClassic`: Dragon Quest-style battles
- `BattleModeTactical`: Grid-based tactical battles

**Example:**
```go
selector := engine.NewBattleSystemSelector(engine.BattleModeClassic)

// Start battle with selected system
selector.StartBattle(players, enemies)
```

---

## UI & Message Systems

### MessageSystem

**File:** `internal/ui/ui_manager.go`  
**Purpose:** Displays messages and dialog to the player

**Responsibilities:**
- Message queue management
- Dialog box rendering
- Text display timing
- Message auto-dismissal

**Example:**
```go
msgSystem := ui.NewMessageSystem()
msgSystem.AddMessage("The enemy attacks!")
msgSystem.AddMessage("You take 15 damage!")

// In game loop
msgSystem.Update(deltaTime)
msgSystem.Draw(screen)
```

---

## Implicit Systems

These are patterns and functionalities that act like systems but aren't formal structs.

### Event Processing

**Purpose:** Triggers events based on player interaction

**Logic Flow:**
1. Check player position vs. event positions
2. Evaluate trigger conditions (proximity, touch, etc.)
3. Check view-specific activation
4. Execute event callbacks
5. Mark events as completed if one-time

**Trigger Types:**
- `TriggerOnTouch`: Collision-based
- `TriggerOnProximity`: Distance-based
- `TriggerOnInteract`: Button press required
- `TriggerOnCondition`: Custom logic

**Example Implementation:**
```go
// In game update loop
for _, entity := range world.GetEntities() {
    eventComp := entity.Event()
    if eventComp == nil {
        continue
    }
    
    // Check if event is active in current view
    if !viewMgr.IsEventActive(eventComp) {
        continue
    }
    
    // Check trigger condition
    if eventComp.TriggerType == components.TriggerOnTouch {
        if checkCollision(player, entity) {
            eventComp.Execute(player)
        }
    }
}
```

---

### Animation System

**Purpose:** Updates sprite animations based on entity state

**Logic Flow:**
1. Update frame timers
2. Advance to next frame when timer expires
3. Handle animation loops
4. Manage animation state transitions
5. Apply temporary animations

**Example Implementation:**
```go
// In game update loop
for _, entity := range world.GetEntities() {
    animComp := entity.Animation()
    if animComp == nil {
        continue
    }
    
    animComp.Update(deltaTime)
}
```

**Animation States:**
- `AnimationIdle`: Standing still
- `AnimationWalking`: Moving
- `AnimationRunning`: Running
- `AnimationAttacking`: Combat action
- `AnimationHurt`: Taking damage
- `AnimationDying`: Death animation
- `AnimationCasting`: Magic casting

---

### Collision System

**Purpose:** Detects and resolves entity collisions

**Logic Flow:**
1. Query entities with Collider component
2. Check for overlapping bounding boxes
3. Resolve collisions (blocking or trigger)
4. Trigger collision events

**Example Implementation:**
```go
func CheckCollisions(world *ecs.World) {
    entities := world.FindWithComponent(ecs.ComponentCollider)
    
    for i, entity1 := range entities {
        collider1 := entity1.Collider()
        transform1 := entity1.Transform()
        
        for j := i + 1; j < len(entities); j++ {
            entity2 := entities[j]
            collider2 := entity2.Collider()
            transform2 := entity2.Transform()
            
            if checkOverlap(transform1, collider1, transform2, collider2) {
                handleCollision(entity1, entity2)
            }
        }
    }
}
```

---

### Rendering System

**Purpose:** Draws all visible entities to the screen

**Logic Flow:**
1. Query entities with Sprite/Animation component
2. Sort by rendering layer
3. Draw each entity at its Transform position
4. Handle sprite effects (scale, rotation, alpha)

**Example Implementation:**
```go
func RenderEntities(screen *ebiten.Image, world *ecs.World) {
    entities := world.FindWithComponent(ecs.ComponentSprite)
    
    // Sort by layer
    sort.Slice(entities, func(i, j int) bool {
        return entities[i].Sprite().Layer < entities[j].Sprite().Layer
    })
    
    // Draw each entity
    for _, entity := range entities {
        sprite := entity.Sprite()
        transform := entity.Transform()
        
        if !sprite.Visible {
            continue
        }
        
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(transform.X, transform.Y)
        screen.DrawImage(sprite.Sprite.Img, op)
    }
}
```

---

## System Interaction Patterns

### Battle Initiation Flow

```
EventProcessing → ViewManager → BattleSystemSelector → BattleSystem
       ↓                                                      ↓
   Event Component                                     RPGStats Components
```

1. Event system detects battle trigger
2. ViewManager transitions to Battle view
3. BattleSystemSelector chooses battle mode
4. BattleSystem initializes with entities
5. Battle loop updates RPGStats

---

### Item Usage Flow

```
UI Input → ConsumableManager → RPGStats/StatusEffects → MessageSystem
                                        ↓
                                 Inventory Component
```

1. Player selects item from inventory
2. ConsumableManager applies effects
3. Components updated (HP restored, etc.)
4. MessageSystem displays result

---

### Animation State Flow

```
GameLogic → Animation Component → Rendering System
              ↓
       State Transitions
```

1. Game logic changes entity state (moving, attacking, etc.)
2. Animation component switches animation state
3. Rendering system displays current frame

---

## Creating Custom Systems

### Step 1: Define System Struct

```go
type MyCustomSystem struct {
    world *ecs.World
    // Additional fields as needed
}

func NewMyCustomSystem(world *ecs.World) *MyCustomSystem {
    return &MyCustomSystem{
        world: world,
    }
}
```

### Step 2: Implement Update Method

```go
func (s *MyCustomSystem) Update(deltaTime float64) {
    // Query relevant entities
    entities := s.world.FindWithComponent(ecs.ComponentMyComponent)
    
    for _, entity := range entities {
        // Get required components
        myComp := entity.MyComponent()
        if myComp == nil {
            continue
        }
        
        // Additional component checks
        transform := entity.Transform()
        if transform == nil {
            continue
        }
        
        // System logic
        myComp.Value += deltaTime
        transform.X += myComp.Velocity * deltaTime
    }
}
```

### Step 3: Integrate with Game Loop

```go
// In game initialization
mySystem := NewMyCustomSystem(world)

// In game update loop
mySystem.Update(deltaTime)
```

---

## System Best Practices

### 1. Single Responsibility
Each system should handle ONE aspect of game logic:

✅ **Good:**
```go
type MovementSystem struct {}  // Only handles movement
type CombatSystem struct {}    // Only handles combat
type AnimationSystem struct {} // Only handles animation
```

❌ **Bad:**
```go
type GameSystem struct {}  // Handles everything
```

### 2. Query Efficiently
Cache entity queries when appropriate:

```go
type MySystem struct {
    world           *ecs.World
    cachedEntities  []*ecs.Entity
    lastUpdateTime  time.Time
}

func (s *MySystem) Update(deltaTime float64) {
    // Refresh cache periodically
    if time.Since(s.lastUpdateTime) > time.Second {
        s.cachedEntities = s.world.FindWithComponent(ecs.ComponentMyComponent)
        s.lastUpdateTime = time.Now()
    }
    
    // Process cached entities
    for _, entity := range s.cachedEntities {
        // ...
    }
}
```

### 3. Handle Missing Components Gracefully

```go
func (s *MySystem) Update(deltaTime float64) {
    for _, entity := range s.world.GetEntities() {
        // Always check for nil
        comp := entity.MyComponent()
        if comp == nil {
            continue  // Skip entities without required component
        }
        
        // Now safe to use comp
        comp.Update(deltaTime)
    }
}
```

### 4. Don't Modify Entity Structure
Systems should modify component data, not entity structure:

✅ **Good:**
```go
func (s *System) Update(deltaTime float64) {
    entity.RPGStats().CurrentHP -= 10  // Modify data
}
```

❌ **Bad:**
```go
func (s *System) Update(deltaTime float64) {
    entity.AddComponent(...)  // Don't modify structure
    entity.RemoveComponent(...)
}
```

### 5. Use Callbacks for Cross-System Communication

```go
type BattleSystem struct {
    onBattleEnd func(won bool)
}

func (bs *BattleSystem) SetBattleEndCallback(callback func(bool)) {
    bs.onBattleEnd = callback
}

func (bs *BattleSystem) endBattle(won bool) {
    if bs.onBattleEnd != nil {
        bs.onBattleEnd(won)
    }
}
```

### 6. Document System Dependencies

```go
// MovementSystem processes entity movement
// 
// Required Components:
//   - Transform: Position to update
//   - Velocity: Movement speed and direction
//
// Optional Components:
//   - Collider: For collision detection
//   - Animation: For movement animation
//
// Dependencies:
//   - CollisionSystem: Must run after collision resolution
type MovementSystem struct {
    // ...
}
```

---

## System Execution Order

Order matters! Some systems depend on others:

```go
// Recommended update order
func (g *Game) Update(deltaTime float64) error {
    // 1. Input processing
    inputSystem.Update(deltaTime)
    
    // 2. Game logic systems
    movementSystem.Update(deltaTime)
    animationSystem.Update(deltaTime)
    
    // 3. Physics and collision
    collisionSystem.Update(deltaTime)
    
    // 4. Combat and interactions
    battleSystem.Update(deltaTime)
    eventSystem.Update(deltaTime)
    
    // 5. UI and view management
    viewManager.Update(deltaTime)
    uiManager.Update(deltaTime)
    
    return nil
}
```

---

## Performance Considerations

### 1. Avoid Over-Querying
```go
// ❌ Bad: Query every frame
func (s *System) Update(deltaTime float64) {
    entities := s.world.FindWithComponent(ecs.ComponentTransform)  // Slow!
}

// ✅ Good: Cache query results
type System struct {
    cachedEntities []*ecs.Entity
}

func (s *System) RefreshCache() {
    s.cachedEntities = s.world.FindWithComponent(ecs.ComponentTransform)
}
```

### 2. Early Exit for Inactive Systems
```go
func (s *BattleSystem) Update(deltaTime float64) {
    if !s.active {
        return  // Skip processing if not in battle
    }
    // Process battle logic...
}
```

### 3. Use Spatial Partitioning for Collision
For large worlds, divide space into regions:

```go
type CollisionSystem struct {
    spatialGrid *SpatialGrid
}

func (s *CollisionSystem) Update(deltaTime float64) {
    // Only check entities in nearby cells
    nearbyEntities := s.spatialGrid.GetNearby(playerPosition)
    // Check collisions only with nearby entities
}
```

---

## Testing Systems

### Unit Test Example

```go
func TestMovementSystem(t *testing.T) {
    // Setup
    world := ecs.NewWorld()
    entity := ecs.NewEntity("Test")
    entity.AddComponent(ecs.ComponentTransform, components.NewTransform(0, 0, 32, 32))
    entity.AddComponent(ecs.ComponentVelocity, &Velocity{X: 5, Y: 0})
    world.AddEntity(entity)
    
    system := NewMovementSystem(world)
    
    // Execute
    system.Update(1.0)  // 1 second
    
    // Assert
    transform := entity.Transform()
    if transform.X != 5.0 {
        t.Errorf("Expected X=5.0, got X=%f", transform.X)
    }
}
```

---

## See Also

- [ECS Architecture](./ecs-architecture.md) - Core ECS concepts
- [Component Reference](./component-reference.md) - Available components
- [Entity Reference](./entity-reference.md) - Entity types
- [ECS Design Rationale](./ecs-design-rationale.md) - Why ECS?
