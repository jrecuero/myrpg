# Adding New Components

This guide walks you through the complete process of adding a new component to the ECS system.

## Overview

Adding a new component involves:
1. Creating the component struct
2. Adding a component constant
3. Adding an accessor method to Entity
4. Implementing constructor functions
5. Documenting the component
6. Using the component in systems

## Step-by-Step Guide

### Step 1: Create Component File

Create a new file in `internal/ecs/components/` for your component:

**File:** `internal/ecs/components/your_component.go`

```go
package components

// YourComponent represents [describe what this component does]
// It is used for [explain use cases]
type YourComponent struct {
    // Add your fields here with descriptive comments
    Field1 string  // Description of Field1
    Field2 int     // Description of Field2
    Field3 bool    // Description of Field3
}

// NewYourComponent creates a new YourComponent with the specified parameters
// param1: description of first parameter
// param2: description of second parameter
// returns a pointer to the newly created YourComponent
func NewYourComponent(param1 string, param2 int) *YourComponent {
    return &YourComponent{
        Field1: param1,
        Field2: param2,
        Field3: false, // default value
    }
}

// Add any helper methods if needed (but keep them minimal)
// Components should primarily be data containers

// Method1 performs [describe what it does]
// Returns [describe return value]
func (c *YourComponent) Method1() string {
    return c.Field1
}
```

### Step 2: Add Component Constant

Add your component name constant to `internal/ecs/constants.go`:

```go
const (
    ComponentTransform    = "transform"
    ComponentSprite       = "sprite"
    // ... existing components ...
    ComponentYourComponent = "yourcomponent"  // Add your new component
)
```

**Naming Convention:**
- Use `Component` prefix
- CamelCase for the constant name
- Lowercase string value
- String value should match component purpose

### Step 3: Add Entity Accessor Method

Add a type-safe accessor method to `internal/ecs/entity.go`:

```go
// YourComponent retrieves the YourComponent from the entity.
// returns a pointer to the YourComponent or nil if not found.
func (e *Entity) YourComponent() *components.YourComponent {
    if comp, exists := e.GetComponent(ComponentYourComponent); exists {
        if yourComp, ok := comp.(*components.YourComponent); ok {
            return yourComp
        }
    }
    return nil
}
```

**Location:** Add this method with the other component accessor methods (around line 100-200)

### Step 4: Implement Component Logic

#### Example: StatusEffect Component

Let's walk through a complete example - adding a StatusEffect component for tracking buffs/debuffs:

**File:** `internal/ecs/components/status_effects.go`

```go
package components

import "time"

// StatusEffectType represents different types of status effects
type StatusEffectType int

const (
    StatusPoison StatusEffectType = iota
    StatusBurn
    StatusFreeze
    StatusStun
    StatusRegen
    StatusShield
)

// StatusEffect represents a single status effect on an entity
type StatusEffect struct {
    Type       StatusEffectType // Type of status effect
    Duration   time.Duration    // How long the effect lasts
    Magnitude  int             // Strength of the effect
    StartTime  time.Time       // When the effect was applied
}

// StatusEffectsComponent manages all status effects on an entity
// It tracks buffs, debuffs, and other temporary effects
type StatusEffectsComponent struct {
    Effects []StatusEffect // List of active status effects
}

// NewStatusEffectsComponent creates a new StatusEffectsComponent
// returns a pointer to the newly created component
func NewStatusEffectsComponent() *StatusEffectsComponent {
    return &StatusEffectsComponent{
        Effects: make([]StatusEffect, 0),
    }
}

// AddEffect adds a new status effect to the entity
// effectType: the type of effect to add
// duration: how long the effect should last
// magnitude: the strength of the effect
func (s *StatusEffectsComponent) AddEffect(effectType StatusEffectType, duration time.Duration, magnitude int) {
    effect := StatusEffect{
        Type:      effectType,
        Duration:  duration,
        Magnitude: magnitude,
        StartTime: time.Now(),
    }
    s.Effects = append(s.Effects, effect)
}

// RemoveExpiredEffects removes effects that have expired
// currentTime: the current game time
// returns the number of effects removed
func (s *StatusEffectsComponent) RemoveExpiredEffects(currentTime time.Time) int {
    removed := 0
    newEffects := make([]StatusEffect, 0, len(s.Effects))
    
    for _, effect := range s.Effects {
        if currentTime.Sub(effect.StartTime) < effect.Duration {
            newEffects = append(newEffects, effect)
        } else {
            removed++
        }
    }
    
    s.Effects = newEffects
    return removed
}

// HasEffect checks if the entity has a specific status effect
// effectType: the type of effect to check for
// returns true if the effect is active
func (s *StatusEffectsComponent) HasEffect(effectType StatusEffectType) bool {
    for _, effect := range s.Effects {
        if effect.Type == effectType {
            return true
        }
    }
    return false
}

// GetEffect retrieves a specific status effect if it exists
// effectType: the type of effect to retrieve
// returns the effect and true if found, zero value and false otherwise
func (s *StatusEffectsComponent) GetEffect(effectType StatusEffectType) (StatusEffect, bool) {
    for _, effect := range s.Effects {
        if effect.Type == effectType {
            return effect, true
        }
    }
    return StatusEffect{}, false
}
```

Add constant to `constants.go`:

```go
ComponentStatusEffects = "statuseffects"
```

Add accessor to `entity.go`:

```go
// StatusEffects retrieves the StatusEffectsComponent from the entity.
// returns a pointer to the StatusEffectsComponent or nil if not found.
func (e *Entity) StatusEffects() *components.StatusEffectsComponent {
    if comp, exists := e.GetComponent(ComponentStatusEffects); exists {
        if effects, ok := comp.(*components.StatusEffectsComponent); ok {
            return effects
        }
    }
    return nil
}
```

### Step 5: Use the Component

#### Adding to Entities

```go
// Create entity with your new component
entity := ecs.NewEntity("Player")
entity.AddComponent(ecs.ComponentYourComponent, components.NewYourComponent("value", 42))

// Access the component
if comp := entity.YourComponent(); comp != nil {
    // Use the component
    value := comp.Field1
}
```

#### Creating a System

```go
// Create a system that uses your component
type YourSystem struct {
    world *ecs.World
}

func NewYourSystem(world *ecs.World) *YourSystem {
    return &YourSystem{world: world}
}

func (s *YourSystem) Update(deltaTime time.Duration) {
    for _, entity := range s.world.GetEntities() {
        // Only process entities with your component
        comp := entity.YourComponent()
        if comp == nil {
            continue
        }
        
        // Perform your system logic
        // comp.Field2 += 1
    }
}
```

#### Status Effect System Example

```go
type StatusEffectSystem struct {
    world *ecs.World
}

func (s *StatusEffectSystem) Update(deltaTime time.Duration) {
    currentTime := time.Now()
    
    for _, entity := range s.world.GetEntities() {
        effects := entity.StatusEffects()
        if effects == nil {
            continue
        }
        
        // Remove expired effects
        effects.RemoveExpiredEffects(currentTime)
        
        // Apply effect logic
        if effects.HasEffect(StatusPoison) {
            if stats := entity.RPGStats(); stats != nil {
                effect, _ := effects.GetEffect(StatusPoison)
                stats.CurrentHP -= effect.Magnitude
            }
        }
        
        if effects.HasEffect(StatusRegen) {
            if stats := entity.RPGStats(); stats != nil {
                effect, _ := effects.GetEffect(StatusRegen)
                stats.CurrentHP += effect.Magnitude
                if stats.CurrentHP > stats.MaxHP {
                    stats.CurrentHP = stats.MaxHP
                }
            }
        }
    }
}
```

## Component Design Principles

### 1. Single Responsibility
Each component should represent one logical aspect:
✅ Good: `TransformComponent` (position and size)
❌ Bad: `PlayerDataComponent` (position, stats, inventory, quests)

### 2. Data-Focused
Components should be primarily data containers:
```go
// ✅ Good - data with minimal helper methods
type HealthComponent struct {
    Current int
    Max     int
}

func (h *HealthComponent) IsDead() bool {
    return h.Current <= 0
}

// ❌ Bad - too much logic in component
type HealthComponent struct {
    Current int
    Max     int
}

func (h *HealthComponent) Update(deltaTime float64) {
    // Complex game logic doesn't belong here
}

func (h *HealthComponent) ApplyDamage(amount int, attacker *Entity) {
    // This logic should be in a system
}
```

### 3. Clear Naming
Use descriptive names that indicate the component's purpose:
- Use `Component` suffix for type names
- Use clear, domain-specific terminology
- Avoid abbreviations unless universally understood

### 4. Meaningful Defaults
Constructor functions should set sensible defaults:
```go
func NewHealthComponent(maxHP int) *HealthComponent {
    return &HealthComponent{
        Current: maxHP,  // Start at full health
        Max:     maxHP,
    }
}
```

### 5. Documentation
Every component should have:
- Package-level comment explaining purpose
- Field comments describing each field
- Constructor function documentation
- Method documentation for helper functions

## Common Patterns

### Optional Fields
Use pointers for optional complex fields:
```go
type QuestComponent struct {
    ActiveQuests    []Quest
    OptionalRewards *RewardData  // nil if no rewards
}
```

### Enumerations
Use typed constants for component states:
```go
type MovementState int

const (
    MovementIdle MovementState = iota
    MovementWalking
    MovementRunning
    MovementJumping
)
```

### Time-Based Data
Use `time.Time` and `time.Duration` for temporal data:
```go
type CooldownComponent struct {
    LastUseTime   time.Time
    CooldownTime  time.Duration
}

func (c *CooldownComponent) IsReady() bool {
    return time.Since(c.LastUseTime) >= c.CooldownTime
}
```

### Collections
Use slices for dynamic collections:
```go
type BufferComponent struct {
    ActiveBuffs []Buff
}

func (b *BufferComponent) AddBuff(buff Buff) {
    b.ActiveBuffs = append(b.ActiveBuffs, buff)
}
```

## Testing Your Component

Create a test file: `internal/ecs/components/your_component_test.go`

```go
package components

import (
    "testing"
)

func TestNewYourComponent(t *testing.T) {
    comp := NewYourComponent("test", 42)
    
    if comp.Field1 != "test" {
        t.Errorf("Expected Field1 to be 'test', got '%s'", comp.Field1)
    }
    
    if comp.Field2 != 42 {
        t.Errorf("Expected Field2 to be 42, got %d", comp.Field2)
    }
}

func TestYourComponentMethod(t *testing.T) {
    comp := NewYourComponent("test", 42)
    
    result := comp.Method1()
    if result != "test" {
        t.Errorf("Expected Method1() to return 'test', got '%s'", result)
    }
}
```

## Checklist

Before considering your component complete, verify:

- [ ] Component struct created in `internal/ecs/components/`
- [ ] Component constant added to `internal/ecs/constants.go`
- [ ] Accessor method added to `internal/ecs/entity.go`
- [ ] Constructor function implemented with sensible defaults
- [ ] All fields documented with comments
- [ ] Helper methods implemented (if needed)
- [ ] Test file created with basic tests
- [ ] Component follows single responsibility principle
- [ ] Component is primarily data-focused
- [ ] Documentation is clear and complete
- [ ] Component integrated into relevant systems
- [ ] Example usage documented

## Troubleshooting

### Component Not Found
**Problem:** `entity.YourComponent()` returns nil even after adding
**Solution:** Ensure you used the exact constant name when adding:
```go
// Wrong
entity.AddComponent("your_component", comp)

// Correct
entity.AddComponent(ecs.ComponentYourComponent, comp)
```

### Type Assertion Failure
**Problem:** Panic on type assertion in accessor
**Solution:** Check that you're storing the correct type:
```go
// Wrong
entity.AddComponent(ecs.ComponentYourComponent, &WrongComponent{})

// Correct
entity.AddComponent(ecs.ComponentYourComponent, components.NewYourComponent())
```

### Import Cycles
**Problem:** Cannot import entity package in components
**Solution:** Components should not import the entity package. Keep them independent.

## Examples from Existing Components

Study these well-designed components as references:

1. **Transform** (`transform.go`) - Simple, focused component
2. **RPGStats** (`rpg_stats.go`) - Complex component with many fields
3. **Equipment** (`equipment.go`) - Component with collection management
4. **Animation** (`animation.go`) - Component with state machine

## Next Steps

After adding your component:
1. Update the [Component Reference](./component-reference.md)
2. Create example entities using your component
3. Implement systems that utilize the component
4. Add integration tests
5. Update relevant game documentation

## See Also

- [ECS Architecture](./ecs-architecture.md) - Core ECS concepts
- [Component Reference](./component-reference.md) - All available components
- [System Design](./system-design.md) - Creating systems for components
