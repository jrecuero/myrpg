# View System Documentation

The View System provides a comprehensive framework for managing different game views, entity visibility, and transitions between game states. This system replaces hardcoded mode switching with a flexible, extensible architecture.

## Core Components

### ViewType
Defines the different types of game views available:
- `ViewExploration` - Free movement exploration
- `ViewTactical` - Grid-based tactical combat  
- `ViewDialog` - Dialog/conversation view
- `ViewInventory` - Inventory management view
- `ViewShop` - Shop interface view
- `ViewMenu` - Game menu view
- `ViewCutscene` - Cutscene/story view
- `ViewWorldMap` - World map navigation

### ViewConfiguration
Defines how each view behaves:
```go
type ViewConfiguration struct {
    Type                ViewType                              // Type of this view
    Name                string                                // Human-readable name
    AllowsPlayerControl bool                                  // Whether player can control characters
    ShowsUI             bool                                  // Whether to show UI elements
    PausesGame          bool                                  // Whether this view pauses underlying game
    EntityFilter        func(*ecs.Entity) bool               // Function to determine which entities to show
    EventFilter         func(*components.EventComponent) bool // Function to determine which events are active
    UpdateHandler       func(float64) error                   // Custom update logic for this view
    InputHandler        func() error                          // Custom input handling for this view
}
```

### ViewManager
Central management system that handles:
- View registration and configuration
- View transitions and state management
- Entity visibility filtering
- Event activation filtering
- View stack management for modal views

## Key Features

### 1. Entity Visibility Management
Each view can define which entities should be visible:
```go
// Example: Exploration view shows only party leader and interactive entities
EntityFilter: func(entity *ecs.Entity) bool {
    if entity.HasTag("background") {
        return true
    }
    if entity.HasTag("player") {
        return entity == vm.game.partyManager.GetPartyLeader()
    }
    if entity.HasTag("enemy") || entity.Event() != nil {
        return true
    }
    return false
}
```

### 2. Event Filtering
Views can control which events are active:
```go
// Example: Only exploration events are active during exploration
EventFilter: func(eventComp *components.EventComponent) bool {
    return eventComp.IsActiveInMode(components.GameModeExploration)
}
```

### 3. View Transitions
Automatic and manual transitions between views:
```go
// Automatic transition from tactical to exploration when combat ends
vm.RegisterTransition(&ViewTransition{
    FromView: ViewTactical,
    ToView:   ViewExploration,
    Condition: func() bool {
        return !vm.game.tacticalManager.IsActive
    },
    TransitionFn: func() error {
        vm.game.tacticalManager.EndTacticalCombat()
        return nil
    },
    Priority: 10,
})
```

### 4. View Stack Management
Support for modal views that can be stacked:
```go
// Push a dialog view (can return to previous view with ESC)
game.SwitchToDialogView(dialogData)

// Pop back to previous view
game.ReturnToPreviousView()
```

## Usage Examples

### Basic View Switching
```go
// Switch to tactical mode
game.SwitchToTacticalMode(participants)

// Switch to inventory
game.SwitchToInventoryView()

// Switch to dialog with context data
dialogData := map[string]interface{}{
    "npc_name": "Merchant",
    "dialog_id": "shop_greeting",
}
game.SwitchToDialogView(dialogData)
```

### Custom View Registration
```go
// Register a custom boss battle view
game.RegisterCustomView(&engine.ViewConfiguration{
    Type: ViewBossBattle,
    Name: "Boss Battle",
    AllowsPlayerControl: true,
    ShowsUI: true,
    EntityFilter: func(entity *ecs.Entity) bool {
        return entity.HasTag("boss") || entity.HasTag("player") || entity.HasTag("arena_effect")
    },
    EventFilter: func(eventComp *components.EventComponent) bool {
        return eventComp.HasTag("boss_battle")
    },
})
```

### View-Specific Entity Management
```go
// Add entities specific to a view
game.GetViewManager().AddViewEntity(ViewTown, shopKeeperEntity)

// Add entities visible in all views  
game.GetViewManager().AddGlobalEntity(backgroundEntity)
```

### Conditional Logic Based on Current View
```go
currentView := game.GetCurrentViewType()
switch currentView {
case engine.ViewExploration:
    // Exploration-specific logic
case engine.ViewTactical:
    // Combat-specific logic
case ViewBossBattle:
    // Boss battle-specific logic
}

// Check if in specific view
if game.IsInView(ViewTown) {
    // Town-specific behavior
}
```

## Integration with Game Systems

### Event System Integration
The View System integrates seamlessly with the event system:
```go
// Events can trigger view transitions
func (g *Game) handleBattleEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
    // Battle events automatically switch to tactical view
    participants := g.getAllCombatParticipants()
    g.SwitchToTacticalMode(participants)
    return &events.EventResult{Success: true}
}
```

### UI System Integration
Views control UI visibility and behavior:
```go
// Views can pause game input while showing modal dialogs
if g.viewManager.GetCurrentViewConfig().PausesGame {
    return nil // Skip game input processing
}
```

### Rendering Integration
The rendering system uses view filtering:
```go
// Only render entities visible in current view
visibleEntities := g.viewManager.GetVisibleEntities(allEntities)
for _, entity := range visibleEntities {
    if g.viewManager.IsEntityVisible(entity) {
        // Render entity
    }
}
```

## Benefits

### 1. Flexibility
- Easy to add new view types
- Configurable entity visibility per view
- Custom update and input logic per view

### 2. Maintainability  
- Centralized view management
- Clear separation of concerns
- Reduced code duplication

### 3. Extensibility
- Support for custom views
- Pluggable transition conditions
- View-specific entity management

### 4. Robustness
- Automatic transition handling
- View stack for modal dialogs
- Error handling and fallbacks

## Migration from Legacy System

The View System is designed to work alongside the existing mode system during migration:

1. **Phase 1**: View System runs parallel to existing mode system
2. **Phase 2**: Gradually migrate functionality to use View System
3. **Phase 3**: Remove legacy mode switching code

Current state supports both systems to ensure compatibility while migrating.

## Future Enhancements

- **View Scripting**: Lua/JavaScript scripting for view behavior
- **View Templates**: Common view patterns for rapid development  
- **View Analytics**: Performance and usage tracking per view
- **View Serialization**: Save/load view state for complex scenarios
- **Conditional View Visibility**: Dynamic entity filters based on game state

## Example: Custom Battle Types

The View System makes it easy to create specialized battle types:

```go
// Arena battle with spectators
game.RegisterCustomView(&engine.ViewConfiguration{
    Type: ViewArena,
    EntityFilter: func(entity *ecs.Entity) bool {
        return entity.HasTag("arena_participant") || entity.HasTag("spectator")
    },
    UpdateHandler: func(deltaTime float64) error {
        // Handle crowd reactions, special arena rules
        return nil
    },
})

// Stealth mission view
game.RegisterCustomView(&engine.ViewConfiguration{
    Type: ViewStealth, 
    EntityFilter: func(entity *ecs.Entity) bool {
        // Show only entities within stealth vision range
        return isWithinStealthRange(entity)
    },
    EventFilter: func(eventComp *components.EventComponent) bool {
        // Only stealth-compatible events
        return eventComp.HasTag("stealth_compatible")
    },
})
```

This system provides the foundation for complex, view-specific gameplay mechanics while maintaining clean code architecture.