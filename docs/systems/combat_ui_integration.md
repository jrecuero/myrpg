# Combat UI Integration Documentation

## Overview

This document describes the complete implementation of the Combat UI system for the turn-based tactical combat in MyRPG. The system provides a comprehensive player interface for managing combat actions, target selection, and turn management.

## System Architecture

### Components

1. **CombatUI** (`internal/ui/combat_ui.go`)
   - Main UI controller for combat interface
   - Manages action buttons, target selection, and visual feedback
   - Handles mouse and keyboard input during combat

2. **TacticalManager** (`internal/engine/tactical_manager.go`)
   - Integration layer between CombatUI and TurnBasedCombatManager
   - Sets up UI callbacks and coordinates combat execution
   - Manages UI lifecycle during tactical mode

3. **Main Game Engine** (`internal/engine/engine.go`)
   - Renders CombatUI during tactical mode
   - Integrates UI drawing with existing rendering pipeline

## Features Implemented

### Action Selection Interface
- **Move Button**: Initiates movement mode with valid tile highlighting
- **Attack Button**: Initiates attack mode with target highlighting  
- **End Turn Button**: Completes the current unit's turn
- **Hotkeys**: M (Move), A (Attack), E (End Turn)

### Visual Feedback System
- **Movement Highlighting**: Blue tiles show valid movement positions
- **Attack Highlighting**: Red tiles show valid attack targets
- **Button States**: Buttons are disabled when actions are unavailable
- **Turn Information**: Displays active unit, class, and action points

### Input Handling
- **Mouse Support**: Click on action buttons and grid positions
- **Keyboard Support**: Hotkeys for quick action selection
- **Cancel Support**: ESC key to cancel current action selection

### State Management
- **Combat UI States**: None, Selecting Action, Selecting Move Target, Selecting Attack Target
- **Automatic Updates**: UI refreshes based on current unit's capabilities
- **Clean Transitions**: Proper state reset when combat ends

## Usage Instructions

### Starting Combat UI
The Combat UI automatically activates when:
1. Tactical mode is enabled (`g.IsTacticalMode()` returns true)
2. Turn-based combat is active (`tm.UseTurnBasedCombat` is true)
3. It's the player team's turn (`combatManager.IsPlayerTurn()` returns true)

### Player Interaction Flow
1. **Action Selection**: Player sees available action buttons
2. **Target Selection**: After selecting Move/Attack, player clicks on highlighted tiles
3. **Action Execution**: Selected actions are automatically executed
4. **Turn Progression**: Combat continues to next unit/team

### Visual Layout
- **Action Buttons**: Right side panel (800-200px width)
- **Turn Info Panel**: Top right corner showing unit details
- **Grid Overlays**: Highlighted tiles for movement/attack ranges
- **Instruction Text**: Bottom area with current action prompts

## Technical Implementation

### UI State Machine
```
CombatUIStateNone → CombatUIStateSelectingAction
    ↓ (Move selected)
CombatUIStateSelectingMoveTarget → Execute → Reset to None
    ↓ (Attack selected)  
CombatUIStateSelectingAttackTarget → Execute → Reset to None
```

### Action Callback System
```go
// Action selection callback
func(actionType tactical.ActionType) {
    // Handle Move/Attack/EndTurn selection
}

// Move target callback  
func(gridPos tactical.GridPos) {
    // Execute movement to selected position
}

// Attack target callback
func(target *ecs.Entity) {
    // Execute attack on selected target
}
```

### Coordinate System Integration
- **Screen to Grid**: Converts mouse clicks to grid positions
- **Grid Validation**: Ensures clicks are within valid grid bounds
- **Offset Handling**: Accounts for UI panel offsets (GridOffsetX, GridOffsetY)

## Integration Points

### TacticalManager Integration
- **Combat UI Creation**: UI instance created in `NewTacticalManager()`
- **Callback Setup**: Callbacks configured to execute combat actions
- **Update Loop**: UI updated each frame during tactical mode
- **State Reset**: UI reset when combat ends

### Main Game Integration
- **Rendering**: CombatUI drawn after grid in main Draw() method
- **Mode Detection**: UI only active during tactical mode
- **Input Priority**: Combat UI input processed before other systems

### Combat Manager Integration
- **Action Creation**: Uses `CreateMoveAction()`, `CreateAttackAction()`, `CreateEndTurnAction()`
- **Action Execution**: All actions executed through `ExecuteAction()`
- **Valid Targets**: Uses `GetValidMovesForUnit()` and `GetValidAttackTargetsForUnit()`

## Error Handling

### Graceful Degradation
- Missing active unit → UI disabled
- Invalid action creation → Error logged, action cancelled  
- Action execution failure → Error logged, UI reset
- Combat end during UI interaction → UI properly reset

### User Feedback
- Disabled buttons for unavailable actions
- Clear instruction text for current mode
- Visual highlighting for valid selections
- ESC key support for cancelling actions

## Configuration

### UI Layout Constants
```go
ButtonAreaX: 600px (right panel)
ButtonWidth: 180px
ButtonHeight: 30px  
ButtonSpacing: 5px
```

### Color Scheme
```go
ButtonColor: RGBA{60, 60, 100, 200}        // Normal state
ButtonHoverColor: RGBA{80, 80, 120, 220}   // Mouse hover
ButtonDisabledColor: RGBA{40, 40, 40, 150} // Disabled state
MoveHighlight: RGBA{0, 150, 255, 100}      // Blue for movement
AttackHighlight: RGBA{255, 100, 100, 150}  // Red for attacks
```

### Input Mappings
- **M Key**: Move action
- **A Key**: Attack action  
- **E Key**: End turn action
- **ESC Key**: Cancel current selection
- **Mouse Left Click**: Select button or grid position

## Testing and Validation

### Test Scenarios
1. **Action Button Functionality**: Verify all buttons respond to clicks and hotkeys
2. **Target Selection**: Confirm movement and attack highlighting works correctly
3. **Action Execution**: Ensure selected actions are properly executed
4. **State Transitions**: Verify UI state changes correctly between actions
5. **Combat Flow**: Test complete turn cycle from action selection to execution

### Known Limitations
- UI uses `ebitenutil.DebugPrintAt()` for text rendering (no custom fonts yet)
- Single active unit selection per turn (no multi-unit actions)
- Fixed UI layout (not responsive to different screen sizes)

## Future Enhancements

### Planned Improvements
1. **Custom Font Rendering**: Replace debug text with proper font system
2. **Animation Support**: Add smooth transitions between UI states
3. **Multi-Unit Selection**: Support for coordinated team actions
4. **Tooltips**: Hover information for actions and targets
5. **Sound Integration**: Audio feedback for UI interactions

### Extensibility Points
- **Additional Actions**: New action types can be easily added
- **Custom Highlighting**: Different highlight colors for special abilities
- **UI Themes**: Color schemes can be made configurable
- **Input Methods**: Support for game controllers/touch input

## Conclusion

The Combat UI system provides a complete, intuitive interface for turn-based tactical combat. It integrates seamlessly with the existing combat management system and provides clear visual feedback for all player actions. The modular design allows for easy extension and customization as the game evolves.