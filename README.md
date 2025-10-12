# MyRPG

A 2D RPG game built with Go and Ebiten featuring an Entity-Component-System (ECS) architecture, turn-based combat, and animated characters.

## ✨ Features

### 🎮 Core Gameplay
- **Multiple Characters**: Create warriors, mages, and rogues with different stats
- **Turn-Based Combat**: Strategic battle system with physical and magical attacks
- **Character Movement**: Smooth player movement with collision detection
- **RPG Stats**: HP, Level, Experience, Attack, Defense, and job-based progression

### 🎨 Animation System
- **Multi-State Animations**: Characters have idle, walking, and attack animations
- **Sprite Sheet Support**: Automatic frame extraction from sprite sheets
- **Attack Feedback**: Visual sword attack animations during combat
- **Configurable Timing**: Customizable animation speeds and durations

### ⚔️ Battle System
- **Player vs Enemy Combat**: Engage enemies by colliding with them
- **Attack Types**: Choose between physical and magical attacks
- **Damage Calculation**: Strategic combat with defense and attack stats
- **Visual Feedback**: Attack animations provide immediate combat feedback

## 🚀 Quick Start

### Prerequisites
- Go 1.19+
- Required assets in `assets/sprites/` directory

### Building and Running
```bash
# Build the game
go build -o ./bin/myrpg ./cmd/myrpg

# Run the game
./bin/myrpg
```

### Controls
- **Arrow Keys**: Move character
- **Tab**: Switch between characters
- **Space**: Confirm battle actions
- **1/2**: Select attack type in battle (Physical/Magical)
- **3**: Cancel attack

## 📁 Project Structure

```
myrpg/
├── assets/                   # Game assets (sprites, audio, etc.)
│   └── sprites/
│       ├── hero/            # Hero animations (idle, walk, sword)
│       ├── player.png       # Player sprite
│       └── enemy.png        # Enemy sprite
├── cmd/myrpg/               # Main application entry point
├── internal/                # Internal game packages
│   ├── ecs/                 # Entity-Component-System
│   ├── engine/              # Game engine and systems
│   ├── gfx/                 # Graphics and sprite management
│   └── ui/                  # User interface systems
├── doc/                     # Documentation
└── tools/                   # Development tools
```

## 📚 Documentation

Comprehensive documentation is available in the [`doc/`](doc/) directory:

- **[Animation System](doc/ANIMATION_SYSTEM.md)** - Multi-state character animations
- **[Attack Animations](doc/ATTACK_ANIMATIONS.md)** - Combat visual feedback system  
- **[Battle System](doc/BATTLE_FIXES.md)** - Turn-based combat mechanics

## 🛠️ Architecture

### Entity-Component-System (ECS)
- **Entities**: Game objects (players, enemies, backgrounds)
- **Components**: Data containers (Transform, Sprite, Animation, RPGStats)
- **Systems**: Logic processors (Movement, Rendering, Battle, Animation)

### Key Components
- `Transform`: Position and dimensions
- `Sprite`: Visual representation
- `Animation`: Multi-state animation management
- `RPGStats`: Character statistics and progression
- `Collider`: Collision detection and battle triggers

## 🎯 Recent Features

- ✅ **Attack Animation System**: Visual feedback during combat
- ✅ **Multi-State Animations**: Idle, walking, and attack states
- ✅ **Flexible Animation Timing**: Configurable frame rates and durations
- ✅ **Battle System Integration**: Seamless combat with visual feedback
- ✅ **Sprite Sheet Support**: Automatic frame extraction from sprite sheets

## 🚧 Development

Built with modern Go practices:
- Clean architecture with ECS pattern
- Comprehensive documentation
- Modular component system
- Extensible animation framework

---

*A Go + Ebiten RPG Game with ECS Architecture*