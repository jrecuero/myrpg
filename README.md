# MyRPG

A 2D RPG game built with Go and Ebiten featuring an Entity-Component-System (ECS) architecture, turn-based combat, and animated characters.

## ✨ Features

### 🎮 Core Gameplay
- **Multiple Characters**: Create warriors, mages, rogues, clerics, and archers with distinct stats
- **Dual Battle Systems**: Classic Dragon Quest-style combat and grid-based tactical battles
- **Party Management**: Full party system with exploration leader and tactical deployment
- **Character Movement**: Smooth movement with collision detection and mode transitions
- **RPG Stats**: Comprehensive stats including HP, Level, Experience, Attack, Defense, and job-based progression

### 🎨 Animation & Graphics
- **Multi-State Animations**: Characters support idle, walking, and attack animation states
- **Sprite Sheet Support**: Automatic frame extraction and animation sequencing
- **Attack Feedback**: Visual sword attack animations with customizable timing
- **Optimized UI**: Panel-based interface with damage tracking and battle logs
- **File-Based Logging**: Comprehensive debugging with timestamped log files

### ⚔️ Battle Systems
- **Classic Battle Mode**: Dragon Quest-style turn-based combat with formations and rows
- **Tactical Battle Mode**: Grid-based combat with movement ranges and strategic positioning  
- **Formation System**: Front/back row mechanics affecting attack reach and damage
- **Enhanced Combat UI**: Real-time damage display, health tracking, and detailed battle logs
- **Party vs Enemy Groups**: Full party deployment against organized enemy formations

## 🚀 Quick Start

### Prerequisites
- Go 1.19+
- Required assets in `assets/sprites/` directory

### Building and Running
```bash
# Quick build and run
make build && make run

# Or manually:
go build -o ./bin/myrpg ./cmd/myrpg
./bin/myrpg

# For detailed build instructions, see BUILD.md
```

**📖 See [BUILD.md](BUILD.md) for comprehensive build instructions and development setup.**

### Controls
#### Exploration Mode
- **Arrow Keys**: Move party leader
- **Tab**: Switch between party members (tactical mode only)
- **T**: Force tactical mode transition
- **Space**: Auto-enter tactical when enemies nearby

#### Tactical/Classic Battle Mode  
- **Arrow Keys**: Move selected character (tactical) or navigate UI (classic)
- **Space**: Confirm actions and selections
- **1**: Physical attack
- **2**: Magic attack  
- **3**: Defend
- **4**: Use item
- **Esc**: Return to exploration mode

## 📁 Project Structure

```
myrpg/
├── assets/                   # Game assets (sprites, audio, etc.)
│   └── sprites/
│       ├── hero/            # Hero animations (idle, walk, sword)
│       ├── player.png       # Player sprite
│       └── enemy.png        # Enemy sprite
├── cmd/myrpg/               # Main application entry point
│   └── game/
│       ├── entities/        # Entity creation functions
│       └── doc/             # Game-specific documentation
├── internal/                # Internal game packages
│   ├── ecs/                 # Entity-Component-System
│   ├── engine/              # Game engine and systems
│   ├── battle/              # Battle systems (classic & tactical)
│   ├── gfx/                 # Graphics and sprite management
│   ├── ui/                  # User interface systems
│   ├── tactical/            # Grid-based tactical systems
│   ├── logger/              # Logging system
│   └── save/                # Save/load functionality
├── docs/                    # Comprehensive documentation
│   ├── systems/             # Core system documentation
│   ├── debugging/           # Bug fixes and troubleshooting
│   ├── development/         # Planning and roadmaps
│   └── project/             # Project management
├── examples/                # Example implementations and custom views
├── logs/                    # Runtime log files
├── scripts/                 # Build and development scripts
├── test/                    # Test files and testing utilities
└── tools/                   # Development tools
```

## 📚 Documentation

Comprehensive documentation is available in the [`docs/`](docs/) directory:

### 🎮 Core Systems
- **[Animation System](docs/systems/ANIMATION_SYSTEM.md)** - Multi-state character animations with sprite sheet support
- **[Attack Animations](docs/systems/ATTACK_ANIMATIONS.md)** - Combat visual feedback and sword attack system  
- **[Classic Battle System](docs/systems/CLASSIC_BATTLE_SYSTEM.md)** - Dragon Quest-style turn-based combat (active)
- **[Tactical Battle System](docs/systems/TACTICAL_BATTLE_SYSTEM.md)** - Grid-based tactical combat system
- **[Row-Based Formation System](docs/systems/ROW_BASED_FORMATION_SYSTEM.md)** - Enhanced formation mechanics with front/back rows
- **[View System](docs/systems/VIEW_SYSTEM.md)** - Game view management and entity visibility framework

### 🐛 Debugging & Fixes
- **[Battle System Fixes](docs/debugging/BATTLE_FIXES.md)** - Technical battle system improvements and fixes
- **[Party System Fixes](docs/debugging/PARTY_SYSTEM_FIXES.md)** - Party management bug fixes and enhancements
- **[Logging System](docs/debugging/logging_system.md)** - File-based debugging and logging infrastructure

### 🚀 Development
- **[Party System Status](docs/development/PARTY_SYSTEM_STATUS.md)** - Complete party system implementation status
- **[Development Priorities](docs/development/DEVELOPMENT_PRIORITIES.md)** - Current focus and roadmap
- **[Documentation Index](docs/README.md)** - Complete documentation navigation guide

### 🎮 Game-Specific Documentation
- **[Game Mechanics](cmd/myrpg/game/doc/game_mechanics.md)** - Complete game systems overview
- **[Combat Mechanics](cmd/myrpg/game/doc/combat_mechanics.md)** - Detailed combat system guide
- **[Game Documentation Index](cmd/myrpg/game/doc/README.md)** - Quick reference and controls

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

### ✅ Battle System Enhancements
- **Enhanced Combat UI**: Optimized panel widths with 65%/35% split for better space utilization
- **Damage Tracking**: Real-time damage display with health information in battle logs
- **File-Based Logging**: Custom logging system with timestamped files instead of console output
- **Classic Battle System**: Full Dragon Quest-style combat with formation mechanics

### ✅ Party System Implementation
- **Exploration Mode**: Single party leader navigation with background party management
- **Tactical Deployment**: Full party vs enemy group combat with automatic grid positioning
- **Mode Transitions**: Seamless switching between exploration and tactical combat
- **Enhanced Party Management**: Complete party system supporting up to 6 members

### ✅ Technical Improvements
- **Documentation Organization**: Comprehensive docs structure with proper categorization
- **Build System**: Fixed compilation errors and improved .gitignore configuration
- **Architecture**: Clean ECS implementation with modular battle systems
- **Version Control**: Proper tracking of all source files including missing events.go

## 🚧 Development

Built with modern Go practices and actively maintained:
- **Clean Architecture**: ECS pattern with modular battle systems
- **Comprehensive Documentation**: Organized in `/docs/` with categorized guides
- **File-Based Logging**: Timestamped debug logs in `/logs/` directory  
- **Modular Components**: Extensible animation and battle system frameworks
- **Version Control**: Proper git tracking of all source files and documentation
- **Build System**: Makefile and build scripts for streamlined development

---

*A Go + Ebiten RPG Game with ECS Architecture*