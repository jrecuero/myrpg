# Entity Types Reference

Complete reference documentation for all entity types in MyRPG.

## Overview

Entities in MyRPG are created through factory functions that combine appropriate components. This document catalogs all standard entity types, their components, and use cases.

---

## Character Entities

### Player Character

**Factory Functions:**
- `CreatePlayer()` - Basic player with default settings
- `CreatePlayerAtPosition(name, x, y)` - Player at specific location
- `CreatePlayerWithJob(name, x, y, job, level)` - Player with custom job/class
- `CreateAnimatedPlayerWithJob(name, x, y, job, level, animations)` - Animated player

**Components:**
- `Transform`: Position and size (32x32 pixels)
- `Sprite`/`Animation`: Visual representation
- `Collider`: Collision detection (32x32)
- `RPGStats`: Character statistics and combat data
- `Inventory`: Item storage (8x6 grid by default)
- `Equipment`: Worn gear and weapons
- `Skills`: Learned abilities
- `QuestJournal`: Quest tracking (for main player)

**Tags:**
- `player`

**Purpose:**
Main controllable characters in the game, including the protagonist and party members.

**Example:**
```go
// Create a Level 10 Warrior named Conan
player := entities.CreatePlayerWithJob("Conan", 100, 100, components.JobWarrior, 10)

// Create animated player with sprite sheet
animations := entities.CharacterAnimations{
    Animations: []AnimationConfig{
        {State: components.AnimationIdle, SpriteSheet: "assets/sprites/warrior_idle.png", ...},
        {State: components.AnimationWalking, SpriteSheet: "assets/sprites/warrior_walk.png", ...},
        {State: components.AnimationAttacking, SpriteSheet: "assets/sprites/warrior_attack.png", ...},
    },
}
player := entities.CreateAnimatedPlayerWithJob("Gandalf", 200, 200, components.JobMage, 15, animations)
```

**Typical Stats by Job:**
| Job | Base HP | Base MP | Attack | Defense | Magic Attack | Magic Defense | Speed |
|-----|---------|---------|---------|---------|--------------|---------------|-------|
| Warrior | 120 | 20 | 15 | 12 | 5 | 8 | 45 |
| Mage | 60 | 100 | 5 | 5 | 20 | 15 | 50 |
| Rogue | 80 | 40 | 12 | 8 | 8 | 8 | 65 |
| Cleric | 90 | 80 | 8 | 10 | 15 | 18 | 48 |
| Archer | 75 | 30 | 13 | 6 | 7 | 9 | 58 |

---

### Enemy Character

**Factory Functions:**
- `CreateEnemy(x, y)` - Basic enemy at position
- `CreateEnemyWithJob(name, x, y, job, level)` - Enemy with custom job/class
- `CreateAnimatedEnemyWithJob(name, x, y, job, level, spriteSheet)` - Animated enemy

**Components:**
- `Transform`: Position and size (32x32 pixels)
- `Sprite`/`Animation`: Visual representation
- `Collider`: Collision detection (32x32)
- `RPGStats`: Character statistics and combat data
- `CombatState`: Battle-specific state (optional)

**Tags:**
- `enemy`
- `boss` (for boss enemies)

**Purpose:**
Hostile entities that the player battles. Can range from basic monsters to complex bosses.

**Example:**
```go
// Create a Level 5 Goblin
goblin := entities.CreateEnemyWithJob("Goblin", 300, 200, components.JobWarrior, 5)

// Create animated dragon boss
dragon := entities.CreateAnimatedEnemyWithJob("Ancient Dragon", 400, 300, 
    components.JobWarrior, 50, "assets/sprites/dragon.png")
dragon.AddTag("boss")
```

---

## Event Entities

Event entities trigger scripted behaviors when the player interacts with them.

### Battle Event

**Factory Function:** `CreateBattleEvent(id, name, x, y, enemyIDs, battleMap)`

**Components:**
- `Transform`: Position and size
- `Event`: Event configuration and callbacks
- `Collider`: Trigger zone (optional)
- `Sprite`: Visual representation (optional)

**Tags:**
- `event`

**Purpose:**
Triggers combat encounters when the player enters the area.

**Example:**
```go
battleEvent := entities.CreateBattleEvent(
    "forest_battle_01",
    "Goblin Ambush",
    250, 300,
    []string{"goblin_01", "goblin_02"},
    "forest_battlefield",
)
```

**Configuration:**
```go
eventComp.SetEventData(components.EventData{
    Enemies:   []string{"goblin_01", "goblin_02"},  // Enemy IDs to spawn
    BattleMap: "forest_battlefield",                 // Battle background
})
```

---

### Dialog Event

**Factory Function:** `CreateDialogEvent(id, name, x, y, npcID, dialogID)`

**Components:**
- `Transform`: Position and size
- `Event`: Event configuration with dialog data
- `Sprite`: NPC appearance (optional)

**Tags:**
- `event`
- `npc`

**Purpose:**
Initiates conversations with NPCs when the player approaches.

**Example:**
```go
dialogEvent := entities.CreateDialogEvent(
    "merchant_greeting",
    "Village Merchant",
    400, 200,
    "merchant_npc",
    "merchant_shop_intro",
)
```

**Trigger Configuration:**
- Type: `TriggerOnProximity`
- Distance: 48 pixels (1.5 tiles)
- Repeatable: Yes

---

### Chest Event

**Factory Function:** `CreateChestEvent(id, name, x, y, itemIDs, gold, locked)`

**Components:**
- `Transform`: Position and size
- `Event`: Event configuration with loot data
- `Collider`: Interaction zone
- `Sprite`: Chest appearance (optional)

**Tags:**
- `event`

**Purpose:**
Contains items and gold that the player can loot.

**Example:**
```go
chest := entities.CreateChestEvent(
    "treasure_chest_01",
    "Old Chest",
    150, 250,
    []string{"potion", "iron_sword"},
    100,  // 100 gold
    false, // not locked
)
```

**State Tracking:**
- Once opened, event is marked as completed
- Not repeatable by default
- Can require keys if locked

---

### Door/Travel Event

**Factory Function:** `CreateDoorEvent(id, name, x, y, targetMap, targetX, targetY)`

**Components:**
- `Transform`: Position and size
- `Event`: Event configuration with travel data
- `Collider`: Trigger zone

**Tags:**
- `event`

**Purpose:**
Transitions the player between different maps or areas.

**Example:**
```go
door := entities.CreateDoorEvent(
    "cave_entrance",
    "Cave Entrance",
    500, 300,
    "dark_cave_map",
    100, 100,  // spawn position in target map
)
```

---

### Info Event

**Factory Function:** `CreateInfoEvent(id, name, x, y, title, message)`

**Components:**
- `Transform`: Position and size
- `Event`: Event configuration with info text
- `Collider`: Trigger zone
- `Sprite`: Sign or marker (optional)

**Tags:**
- `event`

**Purpose:**
Displays information to the player (signs, plaques, examine targets).

**Example:**
```go
sign := entities.CreateInfoEvent(
    "town_sign",
    "Town Entrance Sign",
    200, 150,
    "Welcome!",
    "Welcome to Riverside Village. Population: 342",
)
```

---

### Custom Event

**Factory Function:** `CreateEventEntity(id, name, x, y, eventComp)`

**Components:**
- `Transform`: Position and size
- `Event`: Custom event configuration
- Optional: `Collider`, `Sprite`

**Tags:**
- `event`
- Custom tags as needed

**Purpose:**
Flexible event entity for custom scripted behaviors.

**Example:**
```go
// Create custom event component
eventComp := components.NewEventComponent(
    "puzzle_crystal",
    "Glowing Crystal",
    components.TriggerOnTouch,
    components.EventCustom,
)

// Set custom callback
eventComp.OnTrigger = func(triggeringEntity *ecs.Entity) {
    // Custom behavior
    logger.Info("Player activated the crystal!")
    // Trigger puzzle logic...
}

crystal := entities.CreateEventEntity("crystal_01", "Glowing Crystal", 300, 300, eventComp)
```

---

## Environmental Entities

### Background Object

**Manual Creation** - No factory function (simple entity)

**Components:**
- `Transform`: Position and size
- `Sprite`: Visual representation

**Tags:**
- `background`

**Purpose:**
Non-interactive decorative elements (trees, rocks, buildings).

**Example:**
```go
tree := ecs.NewEntity("Oak Tree")
tree.AddComponent(ecs.ComponentTransform, components.NewTransform(150, 200, 64, 96))
tree.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(treeSprite, 1.0, 0, 0))
tree.AddTag(ecs.TagBackground)
```

---

### Collider Object

**Manual Creation** - No factory function

**Components:**
- `Transform`: Position and size
- `Collider`: Blocks movement

**Tags:**
- `obstacle`

**Purpose:**
Invisible collision boundaries, walls, and impassable areas.

**Example:**
```go
wall := ecs.NewEntity("Wall")
wall.AddComponent(ecs.ComponentTransform, components.NewTransform(100, 100, 200, 32))
wall.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 200, 32, 0, 0))
wall.AddTag("obstacle")
```

---

## Component Combinations by Entity Type

### Minimal Entity (Background)
```go
- Transform
- Sprite
```

### Interactive Object
```go
- Transform
- Sprite
- Event
- Collider (optional)
```

### Combat Character
```go
- Transform
- Sprite/Animation
- Collider
- RPGStats
- CombatState (in battle)
```

### Full Player Character
```go
- Transform
- Sprite/Animation
- Collider
- RPGStats
- Equipment
- Inventory
- Skills
- QuestJournal
- CombatState (in battle)
```

---

## Entity Creation Patterns

### Pattern 1: Factory Functions

Most common entities have dedicated factory functions:

```go
// Use factory for standard entities
player := entities.CreatePlayer()
enemy := entities.CreateEnemy(x, y)
chest := entities.CreateChestEvent(id, name, x, y, items, gold, locked)
```

**Advantages:**
- Consistent configuration
- Less error-prone
- Easy to use
- Well-tested

---

### Pattern 2: Manual Construction

For custom or unique entities:

```go
customEntity := ecs.NewEntity("Custom")
customEntity.AddComponent(ecs.ComponentTransform, ...)
customEntity.AddComponent(ecs.ComponentSprite, ...)
// Add more components as needed
customEntity.AddTag("custom")
```

**When to Use:**
- One-off unique entities
- Testing/prototyping
- Special boss mechanics
- Custom gameplay elements

---

### Pattern 3: Component Composition

Build complex entities by adding components incrementally:

```go
// Start with basic entity
entity := ecs.NewEntity("Elite Enemy")
entity.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 48, 48))
entity.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(sprite, 1.5, 0, 0))
entity.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent("Elite Guard", components.JobWarrior, 20))
entity.AddTag(ecs.TagEnemy)

// Add special abilities
entity.AddComponent(ecs.ComponentSkills, components.NewSkillsComponent())
entity.Skills().LearnSkill(powerAttackSkill)

// Add boss features
entity.AddTag("boss")
entity.AddComponent(ecs.ComponentBossAbilities, bossAbilities)
```

---

## Entity Lifecycle

### Creation
```go
entity := entities.CreatePlayer()
world.AddEntity(entity)
```

### Modification
```go
// Add components at runtime
entity.AddComponent(ecs.ComponentAnimation, animComp)

// Remove components
entity.RemoveComponent(ecs.ComponentCollider)

// Add/remove tags
entity.AddTag("cursed")
entity.RemoveTag("blessed")
```

### Querying
```go
// Find by tag
players := world.FindWithTag(ecs.TagPlayer)

// Find by component
animatedEntities := world.FindWithComponent(ecs.ComponentAnimation)

// Find by name
player := world.FindByName("Hero")

// Find by ID
entity := world.FindByID(entityID)
```

### Removal
```go
// Remove specific entity
world.RemoveEntity(entity)

// Remove by name
world.RemoveByName("Temporary Enemy")
```

---

## Entity Tags Reference

### Standard Tags

| Tag | Purpose | Example Entities |
|-----|---------|------------------|
| `player` | Player-controlled characters | Hero, party members |
| `enemy` | Hostile combatants | Goblins, dragons, bosses |
| `npc` | Non-player characters | Merchants, villagers |
| `boss` | Boss enemies | Elite enemies, dungeon bosses |
| `event` | Event trigger entities | Chests, doors, battle triggers |
| `background` | Non-interactive decoration | Trees, rocks, buildings |
| `obstacle` | Physical barriers | Walls, collision boundaries |
| `projectile` | Projectiles and effects | Arrows, fireballs |

### Custom Tags

You can add custom tags for game-specific logic:

```go
entity.AddTag("quest_target")
entity.AddTag("merchant")
entity.AddTag("guard")
entity.AddTag("summoned")
entity.AddTag("cursed")
```

---

## Best Practices

### Entity Creation

1. **Use Factory Functions**: Prefer `entities.Create*()` functions for standard entities
2. **Validate Components**: Ensure required components are present
3. **Set Appropriate Tags**: Tag entities for easy querying
4. **Consistent Naming**: Use descriptive entity names

### Component Management

1. **Add Required Components First**: Transform, Sprite, etc.
2. **Check Before Access**: Always check if component exists
3. **Don't Over-Component**: Only add what's needed
4. **Remove Unused Components**: Clean up when features are disabled

### Performance

1. **Reuse Sprites**: Load sprites once, share across entities
2. **Batch Creation**: Create multiple entities together when possible
3. **Tag Efficiently**: Use tags for broad categorization
4. **Clean Up**: Remove entities when no longer needed

### Organization

1. **Group Related Entities**: Keep entity factories in logical modules
2. **Document Entity Types**: Explain purpose and configuration
3. **Version Entities**: Track entity definition changes
4. **Test Factories**: Unit test entity creation functions

---

## Common Entity Recipes

### Boss Enemy with Special Abilities
```go
boss := entities.CreateEnemyWithJob("Dragon Lord", 400, 300, components.JobWarrior, 50)
boss.AddTag("boss")
boss.RPGStats().MaxHP = 5000
boss.RPGStats().CurrentHP = 5000
boss.AddComponent(ecs.ComponentSkills, components.NewSkillsComponent())
boss.Skills().LearnSkill(breathFireSkill)
boss.Skills().LearnSkill(tailWhipSkill)
```

### Quest NPC
```go
questGiver := ecs.NewEntity("Village Elder")
questGiver.AddComponent(ecs.ComponentTransform, components.NewTransform(500, 300, 32, 32))
questGiver.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(elderSprite, 1.0, 0, 0))

eventComp := components.NewEventComponent("elder_quest", "Village Elder", 
    components.TriggerOnProximity, components.EventDialog)
eventComp.SetEventData(components.EventData{
    NPCID:    "elder",
    DialogID: "main_quest_start",
})
questGiver.AddComponent(ecs.ComponentEvent, eventComp)
questGiver.AddTag(ecs.TagNPC)
questGiver.AddTag("quest_giver")
```

### Shopkeeper NPC
```go
merchant := entities.CreateDialogEvent("merchant", "Weapon Merchant", 200, 200, "merchant", "shop_greeting")
merchant.AddTag("merchant")
merchant.AddComponent(ecs.ComponentInventory, components.NewInventoryComponent(10, 10))

// Stock the shop
merchant.Inventory().AddItem(swordItem)
merchant.Inventory().AddItem(armorItem)
merchant.Inventory().AddItem(potionItem)
```

### Animated Combat Entity
```go
animations := entities.CharacterAnimations{
    Animations: []AnimationConfig{
        {State: components.AnimationIdle, SpriteSheet: "assets/knight_idle.png", ...},
        {State: components.AnimationWalking, SpriteSheet: "assets/knight_walk.png", ...},
        {State: components.AnimationAttacking, SpriteSheet: "assets/knight_attack.png", ...},
        {State: components.AnimationHurt, SpriteSheet: "assets/knight_hurt.png", ...},
        {State: components.AnimationDying, SpriteSheet: "assets/knight_death.png", ...},
    },
}
knight := entities.CreateAnimatedPlayerWithJob("Knight", 150, 150, components.JobWarrior, 15, animations)
```

---

## See Also

- [ECS Architecture](./ecs-architecture.md) - Core ECS concepts
- [Component Reference](./component-reference.md) - Available components
- [System Reference](./system-reference.md) - Available systems
- [Adding Components](./adding-components.md) - Creating new components
