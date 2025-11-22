# Component Reference

Complete reference documentation for all ECS components in MyRPG.

## Core Components

### Transform

**File:** `internal/ecs/components/transform.go`  
**Constant:** `ecs.ComponentTransform`

Represents the position and size of an entity in 2D space.

**Fields:**
- `X` (float64): X position in world coordinates
- `Y` (float64): Y position in world coordinates
- `Width` (int): Width of the entity in pixels
- `Height` (int): Height of the entity in pixels

**Constructor:**
```go
transform := components.NewTransform(x, y, width, height)
```

**Methods:**
```go
bounds := transform.Bounds() // Returns image.Rectangle
```

**Usage:**
- Required for rendering entities
- Used for collision detection
- Necessary for positioning in world/battle scenes

**Example:**
```go
entity.AddComponent(ecs.ComponentTransform, components.NewTransform(100, 100, 32, 32))
transform := entity.Transform()
transform.X += 5 // Move entity right
```

---

### SpriteComponent

**File:** `internal/ecs/components/rendering.go`  
**Constant:** `ecs.ComponentSprite`

Manages the visual representation of an entity.

**Fields:**
- `Sprite` (*gfx.Sprite): The sprite image to render
- `Layer` (int): Rendering layer (higher = drawn on top)
- `Visible` (bool): Whether the sprite should be rendered
- Additional sprite sheet properties for animation frames

**Constructor:**
```go
sprite := components.NewSpriteComponent(spriteImg)
```

**Usage:**
- Required for visible entities
- Controls rendering order via Layer
- Can be toggled visible/invisible

**Example:**
```go
entity.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(playerSprite))
sprite := entity.Sprite()
sprite.Visible = false // Hide entity
```

---

### ColliderComponent

**File:** `internal/ecs/components/physics.go`  
**Constant:** `ecs.ComponentCollider`

Defines collision properties for physics interactions.

**Fields:**
- `Enabled` (bool): Whether collision is active
- `IsTrigger` (bool): If true, collision is detected but doesn't block movement
- `Layer` (int): Collision layer for filtering
- Collision bounds based on Transform

**Constructor:**
```go
collider := components.NewColliderComponent(enabled, isTrigger)
```

**Usage:**
- Required for entities that interact physically
- Used for collision detection and triggers
- Works in conjunction with Transform

**Example:**
```go
entity.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, false))
if collider := entity.Collider(); collider != nil {
    collider.IsTrigger = true // Make it a trigger zone
}
```

---

## RPG Components

### RPGStatsComponent

**File:** `internal/ecs/components/rpg_stats.go`  
**Constant:** `ecs.ComponentRPGStats`

Complete RPG statistics for characters.

**Core Stats:**
- `Level` (int): Character level (1-99)
- `Experience` (int): Current experience points
- `ExpToNext` (int): XP needed for next level

**Health & Mana:**
- `CurrentHP` (int): Current health points
- `MaxHP` (int): Maximum health points
- `CurrentMP` (int): Current mana points
- `MaxMP` (int): Maximum mana points

**Combat Stats:**
- `Attack` (int): Physical attack power
- `Defense` (int): Physical defense
- `MagicAttack` (int): Magic attack power
- `MagicDefense` (int): Magic defense
- `Speed` (int): Speed/agility stat (default: 50)
- `Accuracy` (int): Hit accuracy percentage (0-100, default: 85)
- `CritRate` (int): Critical hit rate percentage (0-100, default: 5)

**Tactical Stats:**
- `MoveRange` (int): Movement range in tiles per turn
- `MovesRemaining` (int): Remaining moves this turn
- `MoveHistory` ([]MoveRecord): History for undo functionality

**Character Info:**
- `Job` (JobType): Character class (Warrior, Mage, Rogue, Cleric, Archer)
- `Name` (string): Character display name

**Job Types:**
```go
const (
    JobWarrior  // High HP and Defense, moderate Attack
    JobMage     // High Magic Power, low HP and Defense
    JobRogue    // High Speed and Crit, moderate HP
    JobCleric   // High Magic Defense and Healing, moderate HP
    JobArcher   // High Range Attack and Speed, low Defense
)
```

**Constructor:**
```go
stats := components.NewRPGStatsComponent(components.JobWarrior, "Conan")
```

**Key Methods:**
```go
stats.TakeDamage(amount)           // Apply damage
stats.Heal(amount)                 // Restore HP
stats.IsDead()                     // Check if HP <= 0
stats.GainExperience(amount)       // Add XP and handle level up
stats.LevelUp()                    // Increase level and stats
stats.StartTurn()                  // Reset moves for new turn
stats.CanMove()                    // Check if has moves remaining
stats.RecordMove(fromX, fromZ, toX, toZ, distance)
stats.UndoLastMove()               // Undo last movement
stats.ClearMoveHistory()           // Clear movement history
```

**Usage:**
- Required for all characters (players, enemies, NPCs)
- Manages combat calculations
- Handles leveling and progression
- Tactical battle movement tracking

**Example:**
```go
stats := components.NewRPGStatsComponent(components.JobWarrior, "Conan")
entity.AddComponent(ecs.ComponentRPGStats, stats)

// In combat
if stats := entity.RPGStats(); stats != nil {
    stats.TakeDamage(20)
    if stats.IsDead() {
        // Handle death
    }
}
```

---

### EquipmentComponent

**File:** `internal/ecs/components/equipment.go`  
**Constant:** `ecs.ComponentEquipment`

Manages equipped items and gear.

**Fields:**
- `Weapon` (*Item): Equipped weapon
- `Armor` (*Item): Equipped armor
- `Accessory` (*Item): Equipped accessory
- Equipment slots vary by implementation

**Constructor:**
```go
equipment := components.NewEquipmentComponent()
```

**Methods:**
```go
equipment.EquipWeapon(item)
equipment.UnequipWeapon()
equipment.GetTotalAttackBonus()
equipment.GetTotalDefenseBonus()
```

**Usage:**
- Tracks worn items
- Calculates stat bonuses from gear
- Manages equipment slots

---

### InventoryComponent

**File:** `internal/ecs/components/inventory.go`  
**Constant:** `ecs.ComponentInventory`

Manages carried items and inventory.

**Fields:**
- `Items` ([]Item): List of items
- `Capacity` (int): Maximum number of items
- `Gold` (int): Currency amount

**Constructor:**
```go
inventory := components.NewInventoryComponent(capacity)
```

**Methods:**
```go
inventory.AddItem(item) bool      // Returns false if full
inventory.RemoveItem(itemID)
inventory.HasItem(itemID) bool
inventory.IsFull() bool
inventory.GetItemCount()
```

**Usage:**
- Required for entities that carry items
- Manages item collection
- Handles inventory capacity

---

### SkillsComponent

**File:** `internal/ecs/components/skills.go`  
**Constant:** `ecs.ComponentSkills`

Manages character abilities and skills.

**Fields:**
- `Skills` ([]Skill): List of learned skills
- `SkillPoints` (int): Available skill points

**Constructor:**
```go
skills := components.NewSkillsComponent()
```

**Methods:**
```go
skills.LearnSkill(skill)
skills.HasSkill(skillID) bool
skills.CanUseSkill(skillID) bool
skills.GetSkillLevel(skillID) int
```

**Usage:**
- Tracks learned abilities
- Manages skill progression
- Handles skill usage requirements

---

### QuestJournalComponent

**File:** `internal/ecs/components/quest.go`  
**Constant:** `ecs.ComponentQuestJournal`

Tracks quest progress and completion.

**Fields:**
- `ActiveQuests` ([]Quest): Currently active quests
- `CompletedQuests` ([]string): IDs of completed quests
- Quest-related tracking data

**Constructor:**
```go
journal := components.NewQuestJournalComponent()
```

**Methods:**
```go
journal.AddQuest(quest)
journal.CompleteQuest(questID)
journal.GetQuestProgress(questID)
journal.HasCompletedQuest(questID) bool
```

**Usage:**
- Player-specific component
- Tracks story progression
- Manages quest objectives

---

## Combat Components

### ActionPointsComponent

**File:** `internal/ecs/components/combat.go`  
**Constant:** `ecs.ComponentActionPoints`

Manages action points for tactical combat.

**Fields:**
- `Current` (int): Current action points available
- `Max` (int): Maximum action points per turn
- `RegenRate` (int): AP regenerated per turn

**Constructor:**
```go
ap := components.NewActionPointsComponent(maxAP)
```

**Methods:**
```go
ap.Spend(amount) bool     // Returns false if insufficient
ap.Restore(amount)
ap.HasEnough(amount) bool
ap.RefillToMax()
```

**Usage:**
- Used in tactical battle system
- Limits actions per turn
- Manages resource economy

---

### CombatStateComponent

**File:** `internal/ecs/components/combat.go`  
**Constant:** `ecs.ComponentCombatState`

Tracks combat-specific state and effects.

**Fields:**
- `InCombat` (bool): Whether entity is in battle
- `StatusEffects` ([]StatusEffect): Active buffs/debuffs
- `IsDefending` (bool): Defensive stance
- Combat-specific flags

**Constructor:**
```go
combatState := components.NewCombatStateComponent()
```

**Methods:**
```go
combatState.AddStatusEffect(effect)
combatState.RemoveStatusEffect(effectType)
combatState.HasStatusEffect(effectType) bool
combatState.EnterCombat()
combatState.ExitCombat()
```

**Usage:**
- Tracks battle state
- Manages temporary combat effects
- Controls combat behavior

---

## Visual Components

### AnimationComponent

**File:** `internal/ecs/components/animation.go`  
**Constant:** `ecs.ComponentAnimation`

Controls sprite animation playback.

**Fields:**
- `Animations` (map[string]Animation): Named animations
- `CurrentAnimation` (string): Active animation name
- `FrameIndex` (int): Current frame in animation
- `FrameTime` (float64): Time accumulator for frame timing
- `IsPlaying` (bool): Whether animation is playing
- `Loop` (bool): Whether to loop animation

**Constructor:**
```go
anim := components.NewAnimationComponent()
```

**Methods:**
```go
anim.AddAnimation(name, frames, frameRate)
anim.Play(animName, loop bool)
anim.Stop()
anim.Update(deltaTime)
anim.GetCurrentFrame() *ebiten.Image
```

**Usage:**
- Animates sprite entities
- Controls animation playback
- Manages multiple animation states

**Example:**
```go
anim := components.NewAnimationComponent()
anim.AddAnimation("walk", walkFrames, 0.1)
anim.AddAnimation("idle", idleFrames, 0.2)
anim.Play("walk", true)
entity.AddComponent(ecs.ComponentAnimation, anim)
```

---

## Event Components

### EventComponent

**File:** `internal/ecs/components/event.go`  
**Constant:** `ecs.ComponentEvent`

Handles event-driven behavior and callbacks.

**Fields:**
- `OnInteract` (func()): Callback for interaction events
- `OnCollide` (func(*Entity)): Callback for collision events
- `OnDeath` (func()): Callback when entity dies
- Additional event handlers

**Constructor:**
```go
event := components.NewEventComponent()
```

**Usage:**
- Implements event-driven behavior
- Connects entities to game events
- Enables scripted interactions

**Example:**
```go
event := components.NewEventComponent()
event.OnInteract = func() {
    logger.Info("Player interacted with entity")
}
entity.AddComponent(ecs.ComponentEvent, event)
```

---

## Component Combinations

### Common Entity Types

**Player Character:**
```go
// Full-featured player
entity.AddComponent(ecs.ComponentTransform, ...)
entity.AddComponent(ecs.ComponentSprite, ...)
entity.AddComponent(ecs.ComponentRPGStats, ...)
entity.AddComponent(ecs.ComponentEquipment, ...)
entity.AddComponent(ecs.ComponentInventory, ...)
entity.AddComponent(ecs.ComponentSkills, ...)
entity.AddComponent(ecs.ComponentQuestJournal, ...)
entity.AddComponent(ecs.ComponentAnimation, ...)
entity.AddTag(ecs.TagPlayer)
```

**Enemy:**
```go
// Combat enemy
entity.AddComponent(ecs.ComponentTransform, ...)
entity.AddComponent(ecs.ComponentSprite, ...)
entity.AddComponent(ecs.ComponentRPGStats, ...)
entity.AddComponent(ecs.ComponentCombatState, ...)
entity.AddTag(ecs.TagEnemy)
```

**Interactive NPC:**
```go
// Questgiver or merchant
entity.AddComponent(ecs.ComponentTransform, ...)
entity.AddComponent(ecs.ComponentSprite, ...)
entity.AddComponent(ecs.ComponentEvent, ...)
entity.AddTag(ecs.TagNPC)
```

**Static Object:**
```go
// Background decoration
entity.AddComponent(ecs.ComponentTransform, ...)
entity.AddComponent(ecs.ComponentSprite, ...)
entity.AddTag(ecs.TagBackground)
```

---

## Component Dependencies

Some components work best together:

| Component | Typically Used With |
|-----------|---------------------|
| Transform | Sprite, Collider (required for positioning) |
| Sprite | Transform (required for rendering) |
| Collider | Transform (required for collision bounds) |
| RPGStats | Equipment, Inventory, Skills (for full character) |
| Animation | Sprite (for animated visuals) |
| ActionPoints | RPGStats, CombatState (for tactical combat) |
| Event | Transform, Collider (for interaction zones) |

---

## Performance Notes

**Frequently Accessed Components:**
- Transform: Cached by rendering systems
- Sprite: Cached by rendering systems
- RPGStats: Used every combat frame

**Memory Considerations:**
- Large collections (Inventory, Skills) use slices
- Animation frames share sprite atlas when possible
- Equipment references items rather than copying

**Update Frequency:**
- Animation: Every frame
- CombatState: During combat only
- QuestJournal: On events only

---

## See Also

- [ECS Architecture](./ecs-architecture.md) - Core concepts
- [Adding Components](./adding-components.md) - Creating new components
- [System Design](./system-design.md) - Using components in systems
