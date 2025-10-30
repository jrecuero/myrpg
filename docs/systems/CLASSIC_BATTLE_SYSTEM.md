# Battle System Documentation

The game features two battle systems, though only the Dragon Quest-style classic system is currently active in production.

## Current Active System: Dragon Quest-Style Classic Battles

### Overview
The game uses a classic JRPG battle system inspired by Dragon Quest, featuring:
- Formation-based combat with enemy and player parties
- Speed-based turn order and activity queues  
- Visual formations with enemies at top, players at bottom
- Automatic combat progression with battle logs

### Battle Triggers
- **Battle Events**: Touch red battle event squares in exploration mode
- **Automatic**: Battles start immediately when triggered by events
- **No Manual Triggers**: Battles cannot be started manually

### Battle Features
- **Enemy Formation**: 1-2 rows at top of screen (red area)
- **Player Formation**: 2x2 formation at bottom (green area) 
- **Activity Queue**: Right panel shows turn order based on speed
- **Battle Log**: Real-time combat actions displayed
- **Health Bars**: Color-coded HP visualization for all participants
- **Speed System**: Faster entities (Rogue > Warrior > Mage) act more frequently

### Battle Flow
1. Battle event triggered by player collision
2. Dragon Quest-style battle screen appears
3. Formations automatically positioned
4. Speed-based combat with activity queue
5. Victory/defeat state displayed for 3 seconds
6. Automatic return to exploration mode

## Preserved Legacy System: Tactical Grid-Based Combat

### Status: **DISABLED IN PRODUCTION**
The original tactical combat system has been preserved in the codebase but is not accessible through normal gameplay.

### Location
- **Code**: `internal/battle/tactical/` - Complete tactical system preserved
- **Integration**: Battle system selector maintains both systems
- **Access**: No UI controls or key bindings active

### Features (Preserved)
- Grid-based movement and positioning
- Turn-based tactical combat
- Manual unit positioning and control  
- Strategic terrain and positioning mechanics
- Direct player control of all party members

### Technical Implementation
- All tactical code remains intact and functional
- Battle system selector can switch between systems programmatically
- No user-facing controls or documentation references
- Can be reactivated in future development if needed

## Battle System Architecture

### Battle System Selector
The `BattleSystemSelector` manages which battle system is used:
```go
// Current configuration (production)
currentSystem: BattleSystemClassic  // Default and only active system

// Available systems
BattleSystemTactical // Preserved but disabled
BattleSystemClassic  // Active system
```

### Key Classes
- **`BattleManager`**: Core Dragon Quest battle logic
- **`Formation`**: Enemy/player positioning system  
- **`BattleRenderer`**: Visual rendering and UI
- **`BattleSystemSelector`**: System selection management

### Integration Points
- **Event Handler**: Battle events trigger classic battles only
- **Game Engine**: Classic battle updates and rendering
- **UI System**: Dragon Quest battle UI overlays

## Developer Notes

### Reactivating Tactical System
To reactivate tactical battles in development:
1. Uncomment battle system toggle in `engine.go` (around line 640)
2. Re-add key binding documentation in help system
3. Update initialization messages to include tactical controls
4. Test tactical system integration

### Adding New Battle Systems
The architecture supports additional battle systems:
1. Create new battle type in `BattleSystemType` enum
2. Implement battle manager and renderer
3. Add case in `BattleSystemSelector.StartBattle()`
4. Update UI and documentation accordingly

### Current Battle Balance
- **Enemy Stats**: Level 5 Warriors (strong, high HP)
- **Battle Duration**: ~10-30 seconds depending on party strength
- **Victory Display**: 3-second result screen before return to exploration
- **Speed Scaling**: Job-based + level bonuses for turn frequency

## User Experience

### What Players See
✅ **Dragon Quest Battles**: Classic formation-based combat  
✅ **Battle Events**: Red squares trigger automatic battles  
✅ **Formation Display**: Clear enemy/player positioning  
✅ **Activity Queue**: Turn order visualization  
✅ **Battle Results**: Victory/defeat feedback  

❌ **No Tactical Mode**: Grid-based combat not accessible  
❌ **No Manual Triggers**: Cannot start battles manually  
❌ **No Battle System Choice**: Only Dragon Quest style available  

### Controls
- **Movement**: Arrow keys in exploration
- **Battle Trigger**: Touch red battle events
- **No Battle Controls**: Battles are fully automatic
- **Return**: Automatic after battle completion

This configuration provides a streamlined classic JRPG experience while preserving the tactical system for potential future use.