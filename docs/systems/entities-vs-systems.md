# Entities vs Systems: A Practical Guide

## The Confusion

It's common to be confused about **what should be an entity** versus **what should be a system**. This guide will clarify the distinction with practical examples and decision-making rules.

---

## The Core Distinction

### Entity = "A THING"
**Entities are NOUNS** - they represent game objects that exist in your game world.

Examples:
- A player character
- An enemy goblin
- A treasure chest
- A door
- A fireball projectile
- An NPC merchant

**Key Question:** Can you point at it in the game world?
- ✅ Yes → It's probably an entity
- ❌ No → It's probably NOT an entity

### System = "A PROCESS"
**Systems are VERBS** - they represent behaviors or processes that act on entities.

Examples:
- Moving entities around (MovementSystem)
- Resolving combat (CombatSystem)
- Displaying things on screen (RenderingSystem)
- Checking collisions (CollisionSystem)
- Managing battles (BattleSystem)

**Key Question:** Does it DO something to game objects?
- ✅ Yes → It's probably a system
- ❌ No → It's probably NOT a system

---

## The Simple Rule

```
ENTITY = DATA (What it IS)
SYSTEM = LOGIC (What it DOES)
```

### Entities Hold Information
```go
// Entity: A Goblin (data only)
goblin := ecs.NewEntity("Goblin")
goblin.AddComponent(ecs.ComponentTransform, components.NewTransform(100, 100, 32, 32))  // WHERE it is
goblin.AddComponent(ecs.ComponentRPGStats, stats)                                        // WHAT its stats are
goblin.AddComponent(ecs.ComponentSprite, sprite)                                         // HOW it looks
goblin.AddTag(ecs.TagEnemy)                                                             // WHAT type it is
```

### Systems Process Information
```go
// System: Movement Logic (behavior)
type MovementSystem struct {
    world *ecs.World
}

func (s *MovementSystem) Update(deltaTime float64) {
    // DOES something: moves entities around
    for _, entity := range s.world.GetEntities() {
        transform := entity.Transform()
        velocity := entity.Velocity()
        
        if transform != nil && velocity != nil {
            transform.X += velocity.X * deltaTime  // CHANGES position
            transform.Y += velocity.Y * deltaTime
        }
    }
}
```

---

## Real-World Examples

### Example 1: Combat

**❌ WRONG - Making Combat an Entity:**
```go
// This is WRONG - combat is not a "thing"
combatEntity := ecs.NewEntity("Combat")
combatEntity.AddComponent("attacker", player)
combatEntity.AddComponent("defender", enemy)
// How do you even represent a process as data??
```

**✅ RIGHT - Combat as a System:**
```go
// Combat is a PROCESS that acts on entities
type CombatSystem struct {
    world *ecs.World
}

func (cs *CombatSystem) ExecuteAttack(attacker, defender *ecs.Entity) {
    // Logic: calculate damage, apply it, show effects
    attackerStats := attacker.RPGStats()
    defenderStats := defender.RPGStats()
    
    damage := attackerStats.Attack - defenderStats.Defense
    defenderStats.CurrentHP -= damage
}
```

**The Entities in Combat:**
- Player entity (the attacker)
- Enemy entity (the defender)
- Damage numbers entity (visual effect, optional)

**The System:**
- CombatSystem (processes the attack)

---

### Example 2: A Fireball

**Question:** Is a fireball an entity or a system?

**Answer:** It depends on WHAT you mean!

**The Fireball Projectile = Entity** (it's a thing)
```go
fireball := ecs.NewEntity("Fireball")
fireball.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 16, 16))
fireball.AddComponent(ecs.ComponentSprite, fireballSprite)
fireball.AddComponent(ecs.ComponentVelocity, velocity)
fireball.AddComponent(ecs.ComponentProjectile, projectileData)
fireball.AddTag(ecs.TagProjectile)
```

**The Fireball Movement = System** (it's a process)
```go
type ProjectileSystem struct {
    world *ecs.World
}

func (ps *ProjectileSystem) Update(deltaTime float64) {
    // DOES something: moves all projectiles
    projectiles := ps.world.FindWithTag(ecs.TagProjectile)
    
    for _, projectile := range projectiles {
        transform := projectile.Transform()
        velocity := projectile.Velocity()
        
        // Move the projectile
        transform.X += velocity.X * deltaTime
        transform.Y += velocity.Y * deltaTime
        
        // Check if it hit something
        // Remove if out of bounds
    }
}
```

---

### Example 3: A Shop

**Question:** Is a shop an entity or a system?

**The Shop Building = Entity** (it's a thing in the world)
```go
shopBuilding := ecs.NewEntity("Weapon Shop")
shopBuilding.AddComponent(ecs.ComponentTransform, components.NewTransform(400, 300, 64, 64))
shopBuilding.AddComponent(ecs.ComponentSprite, shopSprite)
shopBuilding.AddComponent(ecs.ComponentEvent, shopEvent)  // Triggers shop UI
shopBuilding.AddTag("shop")
```

**The Shopkeeper = Entity** (it's a person/NPC)
```go
shopkeeper := ecs.NewEntity("Merchant")
shopkeeper.AddComponent(ecs.ComponentTransform, components.NewTransform(420, 320, 32, 32))
shopkeeper.AddComponent(ecs.ComponentSprite, merchantSprite)
shopkeeper.AddComponent(ecs.ComponentInventory, shopInventory)  // Items for sale
shopkeeper.AddTag(ecs.TagNPC)
```

**The Shopping Process = System** (it's a process)
```go
type ShopSystem struct {
    world *ecs.World
}

func (ss *ShopSystem) BuyItem(player, merchant *ecs.Entity, itemID string) error {
    // Logic: check price, transfer item, deduct gold
    playerInventory := player.Inventory()
    merchantInventory := merchant.Inventory()
    
    item := merchantInventory.GetItem(itemID)
    if playerInventory.Gold < item.Price {
        return errors.New("not enough gold")
    }
    
    playerInventory.Gold -= item.Price
    merchantInventory.RemoveItem(itemID)
    playerInventory.AddItem(item)
    return nil
}
```

---

## Common Mistakes

### Mistake 1: Making Processes into Entities

**❌ WRONG:**
```go
// Don't make abstract concepts into entities
battleEntity := ecs.NewEntity("Battle")
dialogEntity := ecs.NewEntity("Dialog")
menuEntity := ecs.NewEntity("Menu")
```

**✅ RIGHT:**
```go
// These are systems/managers
type BattleSystem struct { ... }
type DialogSystem struct { ... }
type MenuSystem struct { ... }
```

**Why?** Battles, dialogs, and menus are PROCESSES, not objects in the game world. You can't see a "battle" - you see the player and enemy entities fighting, which is managed by the battle system.

---

### Mistake 2: Putting Logic in Entities

**❌ WRONG:**
```go
type Player struct {
    ecs.Entity
    x, y float64
}

func (p *Player) Update(deltaTime float64) {
    // Handle input
    if ebiten.IsKeyPressed(ebiten.KeyW) {
        p.y -= 5 * deltaTime
    }
    // Calculate damage
    // Update animations
    // Check collisions
    // etc... ALL THE LOGIC
}
```

**✅ RIGHT:**
```go
// Entity: just data
player := ecs.NewEntity("Player")
player.AddComponent(ecs.ComponentTransform, ...)
player.AddComponent(ecs.ComponentVelocity, ...)
player.AddTag(ecs.TagPlayer)

// Systems: logic separated by concern
type InputSystem struct {}      // Handles input → sets velocity
type MovementSystem struct {}   // Handles velocity → updates position
type AnimationSystem struct {}  // Handles animation state
type CollisionSystem struct {}  // Handles collision detection
```

**Why?** Entities are data containers. Systems contain the logic. This separation makes code testable, maintainable, and flexible.

---

### Mistake 3: Making Everything a System

**❌ WRONG:**
```go
// Don't make game objects into systems
type PlayerSystem struct {
    x, y float64
    hp int
    sprite *Sprite
}

type GoblinSystem struct {
    x, y float64
    hp int
    sprite *Sprite
}
```

**✅ RIGHT:**
```go
// Game objects are entities
player := entities.CreatePlayer()
goblin := entities.CreateEnemy(x, y)

// Systems operate on ALL entities of a type
type MovementSystem struct {
    world *ecs.World  // Works on ALL entities with Transform
}

type RenderSystem struct {
    world *ecs.World  // Works on ALL entities with Sprite
}
```

**Why?** You'd need infinite systems (PlayerSystem, Goblin1System, Goblin2System, etc.). Instead, make them all entities and have ONE system that processes them.

---

## The Decision Tree

Use this flowchart to decide:

```
Is it a concrete object in the game world?
    └─ YES → Entity
        Examples: Player, Enemy, Chest, Door, Potion
    
    └─ NO → Is it a process that acts on game objects?
        └─ YES → System
            Examples: Combat, Movement, Rendering, AI
        
        └─ NO → Is it a property/capability?
            └─ YES → Component
                Examples: Position, Health, Sprite, Inventory
            
            └─ NO → Is it configuration or constants?
                └─ YES → Just data/constants
                    Examples: GameSettings, Constants
```

---

## Practical Test: What Am I?

### Test 1: "A Health Bar"

**Question:** Entity or System?

**Think about it:**
- Is it a thing? → Kinda... it's UI?
- Does it DO something? → It displays health

**Answer:** 
- **The visual UI element** → Could be an entity (UI entities)
- **The logic to display health** → System (RenderingSystem/UISystem)
- **The health value itself** → Component (RPGStats component on character entity)

**Reality Check:**
```go
// The character is an entity
player := entities.CreatePlayer()

// Health is stored in a component
stats := player.RPGStats()
currentHP := stats.CurrentHP

// The UI system reads health and draws the bar
type UISystem struct {}
func (ui *UISystem) DrawHealthBar(entity *ecs.Entity, screen *ebiten.Image) {
    stats := entity.RPGStats()
    percentage := float64(stats.CurrentHP) / float64(stats.MaxHP)
    // Draw bar based on percentage
}
```

---

### Test 2: "Level Up"

**Question:** Entity or System?

**Think about it:**
- Can you point at it? → No, it's an event/action
- Does it DO something? → Yes, it increases stats

**Answer:** System (or method in RPGStats component)

**Reality:**
```go
// Not an entity! Level up is a process
type ProgressionSystem struct {}

func (ps *ProgressionSystem) LevelUp(entity *ecs.Entity) {
    stats := entity.RPGStats()
    stats.Level++
    stats.MaxHP += 10
    stats.Attack += 2
    // ... update other stats
}

// Or as a method in the component
func (stats *RPGStatsComponent) LevelUp() {
    stats.Level++
    stats.MaxHP += 10
    // ...
}
```

---

### Test 3: "A Quest"

**Question:** Entity or System?

**Think about it:**
- Is it a thing? → It's abstract... but has properties
- Can you interact with it? → Yes, in the journal

**Answer:** Component data + System logic

**Reality:**
```go
// Quest data stored in component
type Quest struct {
    ID          string
    Name        string
    Description string
    Objectives  []Objective
    Completed   bool
}

// Component on player entity
type QuestJournalComponent struct {
    ActiveQuests    []Quest
    CompletedQuests []string
}

// System manages quest logic
type QuestSystem struct {}
func (qs *QuestSystem) UpdateQuest(player *ecs.Entity, questID string) {
    journal := player.QuestJournal()
    quest := journal.GetQuest(questID)
    
    // Check objectives
    // Update progress
    // Trigger completion
}
```

The quest data lives in a component, but the quest management logic is a system.

---

## The Relationship

### Entities Need Systems
```go
// Entity without system = static data, nothing happens
player := entities.CreatePlayer()
// Player just sits there... doing nothing

// Add systems to make things happen
movementSystem.Update(deltaTime)    // Now player can move
renderSystem.Draw(screen)           // Now player is visible
combatSystem.Update(deltaTime)      // Now player can fight
```

### Systems Need Entities
```go
// System without entities = logic with nothing to operate on
movementSystem := NewMovementSystem(world)
movementSystem.Update(deltaTime)
// Nothing happens... no entities to move

// Add entities
world.AddEntity(player)
world.AddEntity(enemy)
// Now system has entities to process
```

### Components Connect Them
```go
// Entity holds components
player.AddComponent(ecs.ComponentTransform, transform)

// System reads/modifies components
type MovementSystem struct {}
func (ms *MovementSystem) Update(deltaTime float64) {
    for _, entity := range ms.world.GetEntities() {
        transform := entity.Transform()  // Read component
        transform.X += 5 * deltaTime     // Modify component
    }
}
```

---

## Real MyRPG Examples

### Example: Classic Battle

**Entities in Battle:**
```go
// The fighters are entities
warrior := entities.CreatePlayerWithJob("Conan", x, y, components.JobWarrior, 10)
goblin := entities.CreateEnemy(x, y)

// Each has components
warrior.RPGStats()  // Health, attack, defense
warrior.Transform() // Position in formation
```

**The System:**
```go
// The battle process is a system
battleSystem := classic.NewBattleManager()
battleSystem.StartBattle([]*ecs.Entity{warrior}, []*ecs.Entity{goblin})

// System processes the fight
battleSystem.Update(deltaTime)
battleSystem.ExecuteAttack(warrior, goblin)
```

**Breakdown:**
- **Entities:** Warrior, Goblin (the combatants)
- **Components:** RPGStats (health, stats), Transform (position)
- **System:** BattleManager (combat logic, turn management)

---

### Example: Opening a Chest

**Entities:**
```go
// Player entity
player := entities.CreatePlayer()

// Chest entity
chest := entities.CreateChestEvent("chest_01", "Treasure", 300, 200, 
    []string{"potion", "sword"}, 100, false)
```

**The System:**
```go
// Event system checks for interaction
type EventSystem struct {}

func (es *EventSystem) Update(deltaTime float64) {
    player := es.world.FindWithTag(ecs.TagPlayer)[0]
    events := es.world.FindWithTag("event")
    
    for _, eventEntity := range events {
        if checkCollision(player, eventEntity) {
            eventComp := eventEntity.Event()
            if eventComp.EventType == components.EventChest {
                es.openChest(player, eventEntity)
            }
        }
    }
}

func (es *EventSystem) openChest(player, chest *ecs.Entity) {
    eventComp := chest.Event()
    playerInv := player.Inventory()
    
    // Transfer items
    for _, itemID := range eventComp.EventData.Items {
        playerInv.AddItem(itemID)
    }
    
    // Transfer gold
    playerInv.Gold += eventComp.EventData.Gold
    
    // Mark as opened
    eventComp.IsCompleted = true
}
```

**Breakdown:**
- **Entities:** Player, Chest (the objects)
- **Components:** Inventory (stores items), Event (chest data)
- **System:** EventSystem (interaction logic)

---

## Summary Table

| Aspect | Entity | Component | System |
|--------|--------|-----------|--------|
| **What is it?** | A thing in game world | A property/capability | A process/behavior |
| **Represents** | NOUN | ADJECTIVE | VERB |
| **Contains** | Components | Data | Logic |
| **Examples** | Player, Enemy, Chest | Transform, Health, Sprite | Movement, Combat, Rendering |
| **Question** | "What is it?" | "What can it do?" | "What happens?" |
| **Lifetime** | Created/destroyed | Added/removed | Always running |
| **Code Pattern** | `NewEntity("Name")` | `struct { fields }` | `Update(deltaTime)` |

---

## Key Takeaways

1. **Entities are THINGS** (nouns) - objects you can point at in the game
2. **Systems are PROCESSES** (verbs) - behaviors that act on entities
3. **Components are PROPERTIES** (adjectives) - data that describes entities
4. **Entities = Data**, **Systems = Logic**
5. **When in doubt**: If it exists in the game world → Entity. If it processes things → System.

---

## Quick Reference

**Make it an Entity if:**
- ✅ It exists in the game world
- ✅ You can point at it
- ✅ It has a position
- ✅ It can be created/destroyed
- ✅ Multiple instances can exist

**Make it a System if:**
- ✅ It's a process or behavior
- ✅ It operates on multiple entities
- ✅ It runs every frame/update
- ✅ It implements game logic
- ✅ It coordinates between entities

---

## Still Confused?

### Ask These Questions:

1. **Can I see it in the game?**
   - Yes → Probably entity
   - No → Probably system

2. **Does it have a position?**
   - Yes → Probably entity
   - No → Probably system

3. **How many exist?**
   - Many → Probably entities
   - One → Probably system

4. **Does it DO things or IS it a thing?**
   - IS → Entity
   - DOES → System

---

## See Also

- [ECS Architecture](./ecs-architecture.md) - Core ECS concepts
- [Component Reference](./component-reference.md) - All components
- [Entity Reference](./entity-reference.md) - All entity types
- [System Reference](./system-reference.md) - All systems
- [ECS Design Rationale](./ecs-design-rationale.md) - Why ECS?
