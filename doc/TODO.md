# MyRPG Development TODO List

## Project Overview
A 2D tactical RPG game built with Go and Ebiten engine, featuring grid-based combat and party management.

## ✅ Completed Tasks

### 1. Player Positioning Fixes
- **Issue**: Player positioned too high (Y=150 instead of Y=112)
- **Issue**: Players could move outside game world boundaries
- **Solution**: Updated all Y offsets to use `constants.GameWorldY` (112px)
- **Solution**: Added boundary constraints using `GameWorldTop` and `GameWorldBottom` constants
- **Files Modified**: `engine.go`, `ui_manager.go`, `party_manager.go`
- **Status**: ✅ COMPLETED

### 2. Coordinate System Refactoring (Z → Y)
- **Issue**: Confusing use of "Z" coordinates in a 2D game system
- **Solution**: Systematically renamed all `GridPos.Z` to `GridPos.Y` throughout codebase
- **Files Modified**: All files using GridPos struct
- **Approach**: Used sed commands for bulk replacements with validation
- **Status**: ✅ COMPLETED

### 3. Constants Centralization
- **Issue**: Hardcoded values scattered throughout codebase making maintenance difficult
- **Solution**: Created comprehensive constants package with organized structure
- **Files Created**:
  - `internal/constants/display.go` - Screen dimensions, UI panels, game world boundaries
  - `internal/constants/game.go` - Game mechanics, movement ranges, entity properties
  - `internal/constants/positions.go` - Initial entity placement coordinates
- **Files Updated**: `engine.go`, `ui_manager.go`, `party_manager.go`, `tactical_manager.go`, `combat.go`, `main.go`
- **Constants Added**:
  - Screen dimensions (800x600)
  - UI layout (110px top + 408px game + 80px bottom)
  - Grid system (20x10 tiles, 32px each, offset 50,112)
  - Movement ranges by job class
  - Player speeds and entity dimensions
  - Starting positions for all entities
- **Status**: ✅ COMPLETED

## 🔄 In Progress / Identified Issues

### Additional Hardcoded Values
- **Status**: User mentioned finding additional hardcoded numbers
- **Action Needed**: Identify and catalog remaining hardcoded values
- **Priority**: Medium

## 📋 Pending Features & Improvements

### Core Game Systems
- [🔄] **Turn-Based Combat System** 🎯 *CURRENT PRIORITY*
  - [✅] **Combat Components**: ActionPointsComponent, CombatStateComponent with Team management
  - [✅] **Combat Manager**: TurnBasedCombatManager with phase-based combat flow
  - [✅] **Initiative System**: Team-based turns, team with highest total speed goes first
  - [✅] **Action Points Economy**: Movement/actions cost action points, free-form spending
  - [✅] **Combat Phases**: Initialization → Team Turn → Action Execution → Victory Check
  - [✅] **Basic Combat**: Adjacent tile attacks with damage calculation (Attack - Defense)
  - [✅] **Death Handling**: HP=0 detection and victory condition checking
  - [✅] **Enemy AI**: Simple adjacent attack behavior (no enemy movement yet)
  - [✅] **Integration**: Combat manager wired into existing tactical system with callbacks
  - [✅] **Movement Execution**: Full grid-based movement with AP consumption, validation, and occupancy tracking
  - [✅] **Action Creation**: Helper methods for creating and validating Move/Attack/EndTurn actions
  - [ ] **UI Integration**: Add combat UI for action selection and turn indicators
  - [ ] **Player Input**: Mouse/keyboard input for selecting actions and targets

- [ ] **Enhanced Tactical Positioning**
  - Implement line of sight mechanics
  - Add terrain effects on movement
  - Create cover and concealment system
  - Add area of effect abilities

- [ ] **Action System Redesign**
  - Implement action points system
  - Add different action types (move, attack, defend, special)
  - Create action queue and validation
  - Add action preview system

### User Interface
- [ ] **UI Overhaul**
  - Improve visual design of panels
  - Add better feedback for player actions
  - Implement proper HUD elements
  - Add tooltips and help system

- [ ] **Menu Systems**
  - Main menu implementation
  - In-game pause menu
  - Settings/options menu
  - Save/Load game interface

### Game Features
- [ ] **Combat Enhancements** 🔄 *SHORT-TERM PRIORITY*
  - Add range to different attacks (beyond adjacent tiles)
  - Add equipment system affecting damage/defense calculations
  - Add inventory system to manage equipment and items
  - Add status effects (poison, bleeding, burning, stun, etc.)

- [ ] **Character Progression**
  - Level up system
  - Skill trees
  - Equipment system expansion
  - Character stats improvement

- [ ] **Inventory Management**
  - Item system expansion
  - Equipment slots management
  - Item effects and properties
  - Loot and rewards system

- [ ] **Save/Load System**
  - Game state serialization
  - Save file management
  - Quick save/load functionality
  - Save game metadata

### Graphics & Audio
- [ ] **Animation System**
  - Sprite animations for characters
  - Combat effect animations
  - UI transition animations
  - Particle effects

- [ ] **Audio Integration**
  - Background music system
  - Sound effects for actions
  - Audio settings and controls
  - Dynamic audio based on game state

### Advanced Features
- [ ] **AI System**
  - Enemy AI for tactical combat
  - Pathfinding algorithms
  - Behavior trees for different enemy types
  - Difficulty scaling

- [ ] **Map System**
  - Multiple battle maps
  - Map editor functionality
  - Dynamic map generation
  - Environmental hazards

## 🔧 Technical Debt & Code Quality

### Code Organization
- [ ] **Error Handling**
  - Implement consistent error handling patterns
  - Add proper logging system
  - Create error recovery mechanisms

- [ ] **Testing**
  - Unit tests for core systems
  - Integration tests for game flow
  - Performance benchmarks
  - Test coverage reporting

- [ ] **Documentation**
  - Code documentation and comments
  - Architecture documentation
  - API documentation
  - User manual/guide

### Performance & Optimization
- [ ] **Rendering Optimization**
  - Sprite batching
  - Culling off-screen entities
  - Texture atlas optimization
  - Memory usage optimization

## 🎯 Turn-Based Combat Design Specification

### Combat Flow Overview
1. **Combat Initialization**
   - Calculate team initiative (sum of all team member speeds)
   - Team with higher total speed goes first
   - Initialize action points for all units

2. **Turn Structure**
   - **Team Turn**: All units in active team can spend action points
   - **Action Points**: Units spend AP on movement, attacks, items (free-form order)
   - **Turn End**: When team has no more AP or chooses to end turn
   - **Next Team**: Switch to other team, repeat until combat ends

3. **Combat Actions**
   - **Movement**: Cost AP based on distance moved
   - **Attack**: Cost AP, requires adjacent target (range=1)
   - **Target Selection**: If multiple enemies adjacent, player selects target
   - **Death**: Units with HP≤0 removed from tactical grid

4. **Enemy AI (Initial)**
   - No enemy movement (stationary)
   - Attack adjacent player units if available
   - Simple damage calculation using base stats

### Technical Implementation Plan
- Extend existing `TacticalManager` for combat state
- Add `ActionPoint` system to unit components
- Create `CombatAction` types (Move, Attack, Item, EndTurn)
- Implement turn management and initiative calculation
- Add target selection UI for attack actions
- Create unit death/removal system

## 🎯 Current Architecture

### Project Structure
```
myrpg/
├── cmd/myrpg/main.go           # Entry point
├── internal/
│   ├── constants/              # ✅ Centralized constants
│   │   ├── display.go         # Screen/UI constants
│   │   ├── game.go            # Game mechanics constants
│   │   └── positions.go       # Entity positioning constants
│   ├── engine/                 # Core game engine
│   ├── ecs/                    # Entity Component System
│   ├── gfx/                    # Graphics utilities
│   ├── ui/                     # User Interface
│   ├── world/                  # World/Map management
│   ├── tactical/               # Tactical combat system
│   ├── audio/                  # Audio system
│   └── save/                   # Save/Load functionality
├── assets/                     # Game assets
└── tools/                      # Development tools
```

### Key Systems
- **ECS Architecture**: Entity Component System for game objects
- **Grid-based Movement**: 20x10 tile tactical grid (32px tiles)
- **Party Management**: Multi-character party with different job classes
- **UI Layout**: Three-panel system (top HUD, game world, bottom controls)

## 📝 Notes

### Development Practices
- Use constants from `internal/constants` package for all configurable values
- Maintain consistent coordinate system (X, Y for 2D positioning)
- Follow Go conventions for naming and structure
- Test changes with `go build` before committing

### Next Priority Items
1. Identify and replace remaining hardcoded values
2. Implement basic turn-based combat mechanics
3. Enhance UI feedback and visual design
4. Add comprehensive error handling

---

*Last Updated: October 13, 2025*
*Status: Post Constants Refactoring - Ready for Next Phase*