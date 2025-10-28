# Event System Documentation

The Event System provides a comprehensive framework for creating interactive elements in the game world. Events can trigger different actions based on player interaction, timing, or other conditions.

## Table of Contents
- [Overview](#overview)
- [Event Types](#event-types)
- [Trigger Conditions](#trigger-conditions)
- [Game Mode Context](#game-mode-context)
- [Visual Representation](#visual-representation)
- [Creating Events](#creating-events)
- [Event Handlers](#event-handlers)
- [Configuration](#configuration)

## Overview

The Event System consists of several key components:

- **EventComponent**: Core component that defines event properties and behavior
- **EventManager**: Handles event registration, trigger detection, and execution
- **Event Handlers**: Functions that execute when events are triggered
- **Event Entities**: Game objects that represent interactive elements in the world

## Event Types

The system supports multiple event types, each with specific behavior:

### EventBattle
Triggers tactical combat with specified enemies.
- **Use Case**: Enemy encounters, boss battles
- **Data**: Enemy list, battle map
- **Color**: Red square when no sprite
- **Example**: Monster encounters in exploration

### EventDialog
Opens dialog with NPCs.
- **Use Case**: Conversations, story exposition
- **Data**: NPC ID, dialog tree ID
- **Color**: Blue square when no sprite
- **Trigger**: Usually proximity-based for natural conversation flow

### EventChest
Opens containers with items and gold.
- **Use Case**: Treasure chests, loot containers
- **Data**: Item list, gold amount, locked status
- **Color**: Yellow square when no sprite
- **Example**: Treasure chests, supply caches

### EventDoor
Transitions between areas or maps.
- **Use Case**: Doorways, portals, area transitions
- **Data**: Target map, target coordinates
- **Color**: Brown square when no sprite

### EventTrap
Hidden events that trigger damage or effects.
- **Use Case**: Spike traps, poison darts, hidden dangers
- **Data**: Trap type, damage amount
- **Color**: Purple square (when visible)
- **Special**: Hidden by default (Visible = false)

### EventInfo
Displays information to the player.
- **Use Case**: Signs, plaques, information boards
- **Data**: Title, message text
- **Color**: Green square when no sprite
- **Example**: Town signs, instruction boards

### EventQuest
Triggers quest-related actions.
- **Use Case**: Quest givers, objective markers
- **Data**: Quest ID
- **Color**: Orange square when no sprite

### EventCutscene
Plays cutscenes or story sequences.
- **Use Case**: Story moments, cinematic sequences
- **Data**: Cutscene ID
- **Color**: Purple-blue square when no sprite

### EventShop
Opens shop interfaces.
- **Use Case**: Merchants, item vendors
- **Data**: Shop inventory ID
- **Color**: Silver square when no sprite

### EventRest
Provides rest points and save functionality.
- **Use Case**: Inn beds, campfires, save points
- **Data**: Rest type, costs
- **Color**: Light blue square when no sprite

## Trigger Conditions

Events can trigger based on different conditions:

### TriggerOnTouch
**Behavior**: Activates when player directly collides with event entity
**Use Cases**: 
- Chests (player walks into them)
- Doors (player steps on them)
- Battle events (enemy encounters)
- Info signs (player examines)

**Example**:
```go
eventComp := components.NewEventComponent("chest1", "Treasure Chest", components.TriggerOnTouch, components.EventChest)
```

### TriggerOnProximity
**Behavior**: Activates when player gets within specified range
**Use Cases**:
- Dialog events (natural conversation distance)
- Area-of-effect traps
- Automatic story triggers

**Configuration**:
```go
eventComp.SetConditionData(components.EventConditionData{
    Distance: 48, // Pixels (1.5 tiles at 32px each)
})
```

### TriggerOnTimeout
**Behavior**: Activates after specified time duration
**Use Cases**:
- Timed events
- Automatic story progression
- Environmental changes

**Configuration**:
```go
eventComp.SetConditionData(components.EventConditionData{
    Duration: 5 * time.Second,
})
```

### TriggerOnRoomEntry
**Behavior**: Activates when player enters specific area
**Use Cases**:
- Area-based story triggers
- Environmental effects
- Area unlock rewards

### TriggerOnInteract
**Behavior**: Requires explicit player interaction (key press)
**Use Cases**:
- Optional interactions
- Detailed examinations
- Confirmation-required actions

### TriggerOnQuestState
**Behavior**: Activates based on quest completion status
**Use Cases**:
- Quest-dependent events
- Story progression gates
- Conditional content

### TriggerManual
**Behavior**: Only triggered by game logic, not player action
**Use Cases**:
- Scripted sequences
- Conditional story events
- System-triggered actions

## Game Mode Context

Events can be configured to only appear in specific game modes:

### GameModeExploration (Default)
- Events active during free-roaming exploration
- Most common mode for world interactions
- Events don't appear during tactical combat

### GameModeTactical
- Events active during tactical combat
- Used for battle-specific interactions
- Rare but useful for combat mechanics

### GameModeBoth
- Events active in all game modes
- Use sparingly for global systems

**Example Configuration**:
```go
// Event only active in exploration mode
eventComp.SetActiveInMode(components.GameModeExploration)

// Event active in all modes
eventComp.SetActiveInMode(components.GameModeBoth)
```

## Visual Representation

Events support flexible visual representation:

### Custom Sprites
When sprite path is provided:
```go
eventComp.SetSprite("assets/sprites/chest.png")
eventComp.SetSpriteScale(1.5)
eventComp.SetSpriteOffset(0, -4)
```

### Fallback Colors
When no sprite is available, colored squares are used:
- Colors defined in `internal/constants/game.go`
- Type-specific colors for easy identification
- Consistent visual language across the game

### Visibility Control
```go
eventComp.SetVisible(false)  // Hidden events (traps, timers)
eventComp.SetVisible(true)   // Visible events (default)
```

## Creating Events

### Using Factory Functions

The recommended way to create events is using the factory functions in `cmd/myrpg/game/entities/events.go`:

```go
// Battle event
battleEvent := entities.CreateBattleEvent("battle1", "Orc Encounter", 100, 150, 
    []string{"orc_warrior", "orc_archer"}, "forest_map")

// Chest event  
chestEvent := entities.CreateChestEvent("chest1", "Treasure Chest", 200, 200,
    []string{"iron_sword", "health_potion"}, 50, false)

// Dialog event
dialogEvent := entities.CreateDialogEvent("npc1", "Village Elder", 300, 100,
    "elder_npc", "greeting_dialog")

// Info event
signEvent := entities.CreateInfoEvent("sign1", "Town Sign", 150, 50,
    "Welcome to Greenwood", "A peaceful village nestled in the forest.")

// Trap event (hidden)
trapEvent := entities.CreateTrapEvent("trap1", "Spike Trap", 250, 300,
    "spikes", 15)
```

### Manual Creation

For advanced customization:

```go
// Create base component
eventComp := components.NewEventComponent("custom1", "Custom Event", 
    components.TriggerOnTouch, components.EventInfo)

// Configure behavior
eventComp.SetRepeatable(true)
eventComp.SetEventData(components.EventData{
    Title:   "Custom Sign",
    Message: "This is a custom message",
})

// Configure visuals
eventComp.SetFallbackColor(255, 128, 0) // Orange color
eventComp.SetActiveInMode(components.GameModeExploration)

// Create entity
entity := entities.CreateEventEntity("custom1", "Custom Event", 100, 100, eventComp)
```

## Event Handlers

Event handlers define what happens when events trigger. The system includes default handlers for all event types.

### Default Handlers Location
- **Generic Handlers**: `internal/events/handlers.go`
- **Game-Specific Handlers**: `internal/engine/event_handlers.go`

### Custom Handler Registration

```go
// In game initialization
game.SetupGameEventHandlers()

// Or register custom handlers
eventManager.RegisterHandler(components.EventCustom, func(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
    // Custom logic here
    return &events.EventResult{
        Success: true,
        Message: "Custom event triggered",
        Data:    map[string]interface{}{"custom": "data"},
    }
})
```

## Configuration

### Event Colors

Colors are configured in `internal/constants/game.go`:

```go
var (
    EventColorBattle   = [3]uint8{200, 50, 50}   // Red-ish
    EventColorDialog   = [3]uint8{50, 150, 200}  // Blue-ish  
    EventColorChest    = [3]uint8{200, 200, 50}  // Yellow-ish
    EventColorInfo     = [3]uint8{100, 200, 100} // Green-ish
    // ... more colors
)
```

### Adding New Event Types

1. **Add EventType constant**:
```go
// In internal/ecs/components/event.go
const (
    // ... existing types
    EventCustom EventType = iota // Your new type
)
```

2. **Add color constant**:
```go
// In internal/constants/game.go
EventColorCustom = [3]uint8{255, 0, 255} // Magenta
```

3. **Update NewEventComponent**:
```go
// Add case in NewEventComponent color switch
case EventCustom:
    defaultColor = constants.EventColorCustom
```

4. **Create factory function**:
```go
// In cmd/myrpg/game/entities/events.go
func CreateCustomEvent(id, name string, x, y float64, customData string) *ecs.Entity {
    eventComp := components.NewEventComponent(id, name, components.TriggerOnTouch, components.EventCustom)
    // Configure as needed
    return CreateEventEntity(id, name, x, y, eventComp)
}
```

5. **Add event handler**:
```go
// Register handler in game initialization
eventManager.RegisterHandler(components.EventCustom, HandleCustomEvent)
```

## Best Practices

### Event Placement
- Place events logically in the world
- Consider player flow and discoverability
- Use appropriate trigger conditions for context

### Visual Design
- Provide custom sprites for important events
- Use consistent visual language
- Consider visibility based on event purpose

### Performance
- Events are lightweight but avoid excessive numbers
- Use appropriate trigger conditions to reduce checks
- Consider game mode restrictions for optimization

### Game Mode Usage
- Default to GameModeExploration for world events
- Only use GameModeTactical for combat-specific events
- Avoid GameModeBoth unless truly necessary

### Error Handling
- Events gracefully handle missing data
- Provide meaningful feedback to players
- Log important event triggers for debugging

## Debugging

### Event Information
Events provide debug information through their methods:
- `eventComp.GetEventTypeName()` - Human-readable type name
- `eventComp.GetTriggerConditionName()` - Trigger condition name
- `eventComp.CanTrigger()` - Whether event can currently trigger

### Logging
The event system logs important actions:
- Event registration
- Trigger detection
- Handler execution
- Errors and failures

### Testing
Use the event verification test:
```bash
go run test/verify_events.go
```

## Integration Examples

### Story-Driven Event Chain
```go
// Initial quest giver (repeatable dialog)
questGiver := entities.CreateDialogEvent("elder", "Village Elder", 100, 100, "elder", "quest_start")
questGiver.Event().SetRepeatable(true)

// Quest objective (battle event, triggers once)
questBattle := entities.CreateBattleEvent("quest_battle", "Bandit Camp", 300, 200, 
    []string{"bandit_leader", "bandit_grunt"}, "camp_map")

// Reward chest (unlocks after battle)
rewardChest := entities.CreateChestEvent("reward_chest", "Quest Reward", 350, 200,
    []string{"magic_sword", "gold_pouch"}, 100, true) // Initially locked
```

### Environmental Storytelling
```go
// Information signs provide world building
sign1 := entities.CreateInfoEvent("sign_town", "Town Entrance", 50, 150,
    "Welcome to Riverside", "Population: 847\nFounded: 1847\nMayor: Helena Riverstone")

// Hidden trap adds danger
trap1 := entities.CreateTrapEvent("trap_bridge", "Weak Planks", 200, 150, "fall", 10)
trap1.Event().SetVisible(false) // Hidden until triggered

// Timer event for atmosphere
timer1 := entities.CreateTimerEvent("bell", "Town Bell", 100, 50, components.EventInfo, 60) // Every minute
timer1.Event().SetEventData(components.EventData{
    Title: "Town Bell",
    Message: "The town bell chimes, marking the hour.",
})
```

This event system provides a powerful and flexible foundation for creating rich, interactive game worlds with minimal code complexity.