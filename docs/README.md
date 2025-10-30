# MyRPG Documentation

Welcome to the MyRPG documentation! This directory contains comprehensive guides for all the game's systems and features.

## 📚 Documentation Organization

### 🎮 Systems Documentation (`/systems/`)
Core game systems, implementations, and technical architecture:
- **[Animation System](systems/ANIMATION_SYSTEM.md)** - Complete guide to the flexible multi-state animation system supporting idle, walking, attack, and other character animations
- **[Attack Animations](systems/ATTACK_ANIMATIONS.md)** - Detailed documentation for the sword attack animation system with customizable duration and visual feedback
- **[Tactical Battle System](systems/TACTICAL_BATTLE_SYSTEM.md)** - Turn-based tactical combat system architecture (grid-based)
- **[Classic Battle System](systems/CLASSIC_BATTLE_SYSTEM.md)** - Dragon Quest-style classic battle system (currently active)
- **[Row-Based Formation System](systems/ROW_BASED_FORMATION_SYSTEM.md)** - Enhanced formation mechanics with front/back rows and attack reach
- **[View System](systems/VIEW_SYSTEM.md)** - Comprehensive framework for managing game views, entity visibility, and state transitions
- **[Combat Implementation Status](systems/COMBAT_IMPLEMENTATION_STATUS.md)** - Current state of combat system implementation
- **[Combat UI Integration](systems/combat_ui_integration.md)** - Complete UI system for turn-based combat
- **[Event System](systems/EVENT_SYSTEM.md)** - Event handling and game interaction system
- **[FFT Implementation](systems/FFT_IMPLEMENTATION.md)** - Final Fantasy Tactics-inspired tactical system
- **[Movement Implementation](systems/MOVEMENT_IMPLEMENTATION_COMPLETE.md)** - Grid-based movement system

### 🐛 Debugging Documentation (`/debugging/`)
Bug fixes, diagnostics, troubleshooting, and logging systems:
- **[Battle System Fixes](debugging/BATTLE_FIXES.md)** - Technical details about battle system improvements and bug fixes
- **[Column 0 Diagnostic](debugging/COLUMN_0_DIAGNOSTIC.md)** - Debugging grid column issues
- **[Column 0 Fix Final](debugging/COLUMN_0_FIX_FINAL.md)** - Final resolution for grid positioning
- **[Grid Movement Fix](debugging/GRID_MOVEMENT_FIX.md)** - Movement system debugging
- **[Party System Fixes](debugging/PARTY_SYSTEM_FIXES.md)** - Party management bug fixes
- **[Combat Debugging Round 1](debugging/combat_debugging_and_fixes.md)** - Initial combat system debugging
- **[Combat Debugging Round 2](debugging/combat_debugging_round2.md)** - Advanced combat system debugging
- **[Turn Counter & Movement Fixes](debugging/turn_counter_and_movement_fixes.md)** - Turn progression debugging
- **[Logging System](debugging/logging_system.md)** - File-based debugging and logging infrastructure

### 🚀 Development Documentation (`/development/`)
Planning, roadmaps, priorities, and project status:
- **[Development Priorities](development/DEVELOPMENT_PRIORITIES.md)** - Current development focus and priorities
- **[Party System Status](development/PARTY_SYSTEM_STATUS.md)** - Party system development status
- **[Phase 1 Complete](development/PHASE1_COMPLETE.md)** - Milestone completion documentation
- **[Tactical Roadmap](development/TACTICAL_ROADMAP.md)** - Strategic development roadmap
- **[TODO](development/TODO.md)** - Current task list and upcoming work

### 📋 Project Documentation (`/project/`)
Project management, setup, and organizational changes:
- **[Project Cleanup Completed](project/project_cleanup_completed.md)** - Documentation reorganization and logging setup

## 🎯 Quick Navigation

### By Feature Area
- **Animation & Visual Effects**: See `/systems/` for animation systems and attack effects
- **Combat & Battle**: See `/systems/` for battle mechanics and `/debugging/` for combat fixes  
- **Movement & Grid**: See `/systems/` for movement implementation and `/debugging/` for grid fixes
- **UI & Interface**: See `/systems/` for combat UI integration
- **Project Setup**: See `/project/` for logging and organizational improvements

### By Development Phase
- **Current Work**: Check `/development/TODO.md` and `/development/DEVELOPMENT_PRIORITIES.md`
- **Completed Features**: See `/development/PHASE1_COMPLETE.md` and system docs in `/systems/`
- **Bug Fixes & Debugging**: Browse `/debugging/` for comprehensive troubleshooting guides
- **Future Plans**: Review `/development/TACTICAL_ROADMAP.md` for long-term strategy

## 🚀 Getting Started

1. **Animation System**: Start with `ANIMATION_SYSTEM.md` to understand how character animations work
2. **Battle Mechanics**: Learn about combat in the battle system documentation
3. **Attack Feedback**: Explore `ATTACK_ANIMATIONS.md` for visual combat enhancements

## 🛠️ Development

All documentation is kept up-to-date with the latest code changes and includes:
- Technical implementation details
- Code examples and configuration options
- Usage patterns and best practices
- Asset requirements and file structure

## 📁 Project Structure

```
docs/
├── README.md                          # This file - documentation index
├── systems/                           # Core game systems documentation
│   ├── CLASSIC_BATTLE_SYSTEM.md       # Dragon Quest-style battles (active)
│   ├── TACTICAL_BATTLE_SYSTEM.md      # Tactical grid-based battles (preserved)
│   ├── ROW_BASED_FORMATION_SYSTEM.md  # Enhanced formation mechanics design
│   ├── VIEW_SYSTEM.md                 # Game state and view management
│   ├── ANIMATION_SYSTEM.md            # Character animation system
│   ├── ATTACK_ANIMATIONS.md           # Attack visual effects
│   └── EVENT_SYSTEM.md                # Game events and interactions
├── debugging/                         # Bug fixes and diagnostics
│   ├── BATTLE_FIXES.md                # Battle system improvements
│   ├── combat_debugging_and_fixes.md  # Combat system debugging
│   └── logging_system.md              # Development logging
├── development/                       # Planning and roadmaps
│   ├── TACTICAL_ROADMAP.md            # Development roadmap
│   ├── DEVELOPMENT_PRIORITIES.md      # Current priorities
│   └── TODO.md                        # Task lists
└── project/                          # Project management
    └── project_cleanup_completed.md   # Documentation organization
```

---

*Last updated: October 11, 2025*