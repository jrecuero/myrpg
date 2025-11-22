# Entity-Component-System (ECS) Architecture

## Overview

The MyRPG game uses an Entity-Component-System (ECS) architectural pattern for managing game objects. ECS is a powerful design pattern that promotes composition over inheritance, making the codebase flexible, maintainable, and easy to extend.

## Core Concepts

### Entity

An **Entity** is a container for components - it represents a game object but contains no behavior or data itself beyond an ID and name. Every game object (player, enemy, NPC, item, etc.) is an entity.

**Key Properties:**
- `ID`: Unique 64-bit identifier (auto-generated)
- `Name`: Human-readable name
- `components`: Map of component name to component data
- `tags`: Set of string tags for categorization

**Location:** `internal/ecs/entity.go`

### Component

A **Component** is a data structure that holds specific properties. Components have no behavior - they are pure data containers. Examples include position (Transform), health stats (RPGStats), or visual representation (Sprite).

**Location:** `internal/ecs/components/`

### System

A **System** contains the logic that operates on entities with specific components. Systems implement game behavior by querying for entities with required components and updating their state.

**Location:** Various locations (battle systems, world systems, etc.)

## ECS Architecture Benefits

1. **Composition over Inheritance**: Build complex entities by combining simple components
2. **Data-Oriented Design**: Components are pure data, improving cache locality
3. **Easy Extension**: Add new components without modifying existing code
4. **Flexible Entity Creation**: Any combination of components creates new entity types
5. **Reusability**: Components can be shared across different entity types

## Entity Lifecycle

### Creating an Entity

```go
// Create a new entity
entity := ecs.NewEntity("PlayerCharacter")

// Add components
entity.AddComponent(ecs.ComponentTransform, components.NewTransform(100, 100, 32, 32))
entity.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent(components.JobWarrior, "Conan"))

// Add tags for categorization
entity.AddTag(ecs.TagPlayer)
```

### Retrieving Components

Entities provide type-safe accessor methods for common components:

```go
// Direct accessors (returns nil if component doesn't exist)
transform := entity.Transform()
stats := entity.RPGStats()
sprite := entity.Sprite()

// Generic component access
if comp, exists := entity.GetComponent("custom_component"); exists {
    customComp := comp.(*CustomComponent)
    // Use customComp
}

// Check if component exists
if entity.HasComponent(ecs.ComponentTransform) {
    // Entity has Transform component
}
```

### Removing Components

```go
// Remove a specific component
entity.RemoveComponent(ecs.ComponentAnimation)

// Note: Removing a component doesn't delete the entity
```

## Component Types

### Built-in Components

| Component | Purpose | Key Fields |
|-----------|---------|------------|
| **Transform** | Position and size | X, Y, Width, Height |
| **Sprite** | Visual representation | Image, sprite sheet info |
| **Collider** | Physics collision | Bounds, collision layer |
| **RPGStats** | Character stats | HP, MP, Attack, Defense, Level |
| **Animation** | Sprite animation | Frames, timing, state |
| **Equipment** | Worn items | Weapon, armor slots |
| **Inventory** | Carried items | Item list, capacity |
| **Skills** | Character abilities | Skill list, levels |
| **QuestJournal** | Quest tracking | Active/completed quests |
| **Event** | Event handling | Event callbacks |
| **ActionPoints** | Tactical combat | AP available/max |
| **CombatState** | Battle status | State, effects |

### Component Constants

All component names are defined as constants in `internal/ecs/constants.go`:

```go
const (
    ComponentTransform    = "transform"
    ComponentSprite       = "sprite"
    ComponentCollider     = "collider"
    ComponentRPGStats     = "rpgstats"
    ComponentAnimation    = "animation"
    ComponentActionPoints = "actionpoints"
    ComponentCombatState  = "combatstate"
    ComponentEquipment    = "equipment"
    ComponentInventory    = "inventory"
    ComponentSkills       = "skills"
    ComponentQuestJournal = "questjournal"
    ComponentEvent        = "event"
)
```

## Entity Tags

Tags provide a lightweight way to categorize entities without adding components:

```go
const (
    TagPlayer     = "player"
    TagEnemy      = "enemy"
    TagBackground = "background"
    TagNPC        = "npc"
    TagProjectile = "projectile"
)

// Using tags
entity.AddTag(ecs.TagPlayer)
entity.AddTag("boss")  // Custom tags are allowed

if entity.HasTag(ecs.TagEnemy) {
    // Handle enemy logic
}

// Get all tags
tags := entity.GetTags()
```

## World Management

The `World` struct manages collections of entities:

```go
// Create a world
world := ecs.NewWorld()

// Add entities
world.AddEntity(playerEntity)
world.AddEntity(enemyEntity)

// Query entities by tag
players := world.GetEntitiesByTag(ecs.TagPlayer)
enemies := world.GetEntitiesByTag(ecs.TagEnemy)

// Get all entities
allEntities := world.GetEntities()

// Remove entities
world.RemoveEntity(entityToRemove)
```

## Common Entity Patterns

### Player Character

```go
player := ecs.NewEntity("Hero")
player.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
player.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(sprite))
player.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent(components.JobWarrior, "Hero"))
player.AddComponent(ecs.ComponentEquipment, components.NewEquipmentComponent())
player.AddComponent(ecs.ComponentInventory, components.NewInventoryComponent(20))
player.AddTag(ecs.TagPlayer)
```

### Enemy NPC

```go
enemy := ecs.NewEntity("Goblin")
enemy.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
enemy.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(sprite))
enemy.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent(components.JobWarrior, "Goblin"))
enemy.AddTag(ecs.TagEnemy)
```

### Static Background Object

```go
background := ecs.NewEntity("Tree")
background.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 64, 64))
background.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(sprite))
background.AddTag(ecs.TagBackground)
```

## System Design Pattern

Systems should follow these guidelines:

1. **Query for Required Components**: Only process entities that have the necessary components
2. **Single Responsibility**: Each system handles one aspect of game logic
3. **Update Method**: Implement an Update method that takes deltaTime
4. **No Component Creation**: Systems read/modify components but don't create entities

Example system pattern:

```go
type MovementSystem struct {
    world *ecs.World
}

func (s *MovementSystem) Update(deltaTime float64) {
    // Get all entities with Transform component
    for _, entity := range s.world.GetEntities() {
        transform := entity.Transform()
        if transform == nil {
            continue // Skip entities without Transform
        }
        
        // Update position based on velocity, input, etc.
        // transform.X += velocity.X * deltaTime
    }
}
```

## Best Practices

1. **Keep Components Data-Only**: No methods that change game state
2. **Use Type-Safe Accessors**: Prefer `entity.Transform()` over generic `GetComponent()`
3. **Check for nil**: Always check if component exists before using
4. **Immutable Component Names**: Use constants from `ecs.Component*`
5. **Meaningful Tags**: Use tags for broad categorization, components for specific data
6. **Don't Over-Component**: If data is always together, keep it in one component
7. **Document Dependencies**: Note which components are typically used together

## Performance Considerations

1. **Component Access**: Direct accessors are faster than generic GetComponent
2. **Tag Queries**: Tags are faster than component queries for broad categorization
3. **Cache Queries**: If querying frequently, cache the entity list
4. **Avoid String Keys**: Use constants for component/tag names
5. **Batch Updates**: Process similar entities together in systems

## See Also

- [Adding New Components](./adding-components.md) - Step-by-step guide
- [Component Reference](./component-reference.md) - Detailed component documentation
- [Battle System](../systems/battle-system.md) - Example of ECS in practice
