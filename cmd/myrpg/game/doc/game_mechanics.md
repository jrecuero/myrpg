# Game Mechanics Overview

## Core Systems

### Character System
- **Job Classes**: Warrior, Mage, Rogue, Cleric, Archer
- **Stats**: HP, MP, Attack, Defense, Speed
- **Progression**: Experience points, level advancement
- **Equipment**: Weapons, armor, accessories

### Inventory System
- **Grid-based storage**: Visual item arrangement
- **Item types**: Weapons, armor, consumables, materials
- **Item stacking**: Automatic stacking for consumables
- **Item tooltips**: Detailed item information
- **Drag & drop**: Intuitive item management
- **Sorting/filtering**: Organization tools

### Item System
- **Base items**: Iron Sword, Health Potion, Mana Potion, Magic Crystal
- **Rarity system**: Item quality levels
- **Equipment stats**: Attack/defense bonuses
- **Consumable effects**: Healing, mana restoration, buffs
- **Item registry**: Centralized item database

### Combat System
- **Turn-based**: Strategic action planning
- **Grid-based**: Tactical positioning
- **Action points**: Resource management per turn
- **Team combat**: Multiple units per side
- **Death handling**: Unit removal and cleanup

### World System
- **Overworld mode**: Free exploration
- **Tactical mode**: Grid-based combat
- **Mode switching**: Seamless transitions
- **Grid positioning**: Coordinate-based movement

### UI System
- **Widget-based**: Modular UI components
- **Message system**: Game feedback and notifications
- **Combat UI**: Turn indicators, action buttons
- **Inventory UI**: Grid-based item display
- **Debug UI**: Development tools and overlays

### Audio System
- **Background music**: Atmospheric audio
- **Sound effects**: Action feedback
- **Audio management**: Volume control, audio loading

### Skills System
- **Job-specific skill trees**: Unique abilities for each class
- **Skill point allocation**: Strategic character development
- **Skill types**: Passive bonuses, active abilities, traits, upgrades
- **Prerequisites**: Unlockable skill progression paths
- **Visual skill trees**: Interactive skill tree display
- **Active ability management**: Equippable combat abilities

### Quest System
- **Quest tracking**: Active and completed quest management
- **Objective system**: Kill enemies, collect items, talk to NPCs
- **Quest types**: Main story, side quests, daily tasks, tutorials
- **Progress tracking**: Visual progress indicators and completion status
- **Reward system**: Experience, gold, items, skill points, equipment
- **Quest giver NPCs**: Receive quests from various characters

#### Available Job Classes

**Warrior**: Tanky melee combatant focused on HP and physical damage
- Tough Skin (T1): +10 Maximum HP
- Power Strike (T1): +3 Attack Damage
- Iron Will (T2): +15 Maximum HP, +2 Defense
- Whirlwind (T2): +5 Attack Damage

**Mage**: Magical spellcaster with high MP and magical abilities
- Mana Pool (T1): +15 Maximum MP
- Arcane Knowledge (T1): +2 Magic Power
- Greater Mana (T2): +25 Maximum MP, +3 Magic Power
- Elemental Mastery (T2): +4 Magic Power, +10 Maximum MP

**Rogue**: Agile combatant focused on speed and precision
- Sneak (T1): +3 Speed
- Quick Reflexes (T1): +2 Defense, +2 Speed
- Precise Strike (T1): +4 Attack Damage
- Shadow Step (T2): +5 Speed, +3 Attack Damage
- Evasion (T2): +6 Defense, +8 Maximum HP
- Deadly Strike (T2): +7 Attack Damage

### Save System
- **Character persistence**: Stats, equipment, progress
- **World state**: Dialog progression, quest status
- **Settings**: Game configuration and preferences

## Game Modes

### Overworld Mode
- Free movement and exploration
- NPC interaction and dialogs
- World navigation
- Mode switching to tactical combat

### Tactical Mode
- Grid-based combat system
- Turn-based action selection
- Strategic positioning and movement
- Resource management (AP, HP, MP)

## Controls Reference

| Mode | Key | Action |
|------|-----|--------|
| **Any** | I | Toggle Inventory |
| **Overworld** | Arrow Keys | Move player |
| **Overworld** | Space | Interact with NPCs |
| **Overworld** | T | Enter Tactical Mode |
| **Any** | K | Toggle Skills Window |
| **Any** | J | Toggle Quest Journal |
| **Tactical** | A | Attack adjacent enemy |
| **Tactical** | E | End turn |
| **Tactical** | Esc | Exit Tactical Mode |

## File Structure

```
cmd/myrpg/           # Main application entry
├── main.go          # Application entry point
└── game/doc/        # Game documentation
internal/
├── audio/           # Audio system
├── ecs/             # Entity Component System
├── engine/          # Core game engine
├── gfx/             # Graphics and rendering
├── save/            # Save/load system
├── tactical/        # Combat system
├── ui/              # User interface
├── world/           # World management
└── constants/       # Game constants
assets/
├── audio/           # Audio files
├── sprites/         # Character and item sprites
├── tiles/           # World tile graphics
└── map/             # Map data
test/                # Test programs
tools/               # Development utilities
```

## Development Notes

### Current Status
- ✅ Basic combat system functional
- ✅ Inventory system complete
- ✅ Item system implemented
- ✅ Dialog system basic functionality
- 🔄 Quest system (planned)
- 🔄 Save/load system (planned)
- 🔄 Skills/abilities system (planned)

### Architecture
- **ECS Pattern**: Entity-Component-System for game objects
- **Component-based UI**: Modular widget system
- **Event-driven**: Message passing for system communication
- **Modular design**: Clear separation of concerns

---
*Last Updated: October 2025*