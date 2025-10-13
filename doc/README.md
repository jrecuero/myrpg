# MyRPG Documentation

Welcome to the MyRPG documentation! This directory contains comprehensive guides for all the game's systems and features.

## 📚 Available Documentation

### Core Systems
- **[Animation System](ANIMATION_SYSTEM.md)** - Complete guide to the flexible multi-state animation system supporting idle, walking, attack, and other character animations
- **[Attack Animations](ATTACK_ANIMATIONS.md)** - Detailed documentation for the sword attack animation system with customizable duration and visual feedback
- **[Battle System Fixes](BATTLE_FIXES.md)** - Technical details about battle system improvements and bug fixes

### Tactical RPG Development
- **[FFT Implementation Guide](FFT_IMPLEMENTATION.md)** - Comprehensive guide for converting to Final Fantasy Tactics-style gameplay
- **[Tactical Roadmap](TACTICAL_ROADMAP.md)** - Complete 6-phase development plan for tactical RPG features
- **[Development Priorities](DEVELOPMENT_PRIORITIES.md)** - Quick reference for development priorities and success metrics
- **[Phase 1 Complete](PHASE1_COMPLETE.md)** - Summary of completed grid rendering and mode switching implementation

### Quick Reference

#### Animation System
- Multi-state character animations (idle, walking, attacking, etc.)
- Sprite sheet support with automatic frame extraction
- Configurable frame timing and looping
- Automatic state switching based on character behavior

#### Attack Animations
- Sword attack visual feedback during combat
- Customizable animation duration (default: 1.5 seconds)
- Looping animation during attack period
- Seamless integration with battle system

#### Battle System
- Turn-based combat with physical and magical attacks
- RPG stats integration (HP, Attack, Defense, etc.)
- UI panels for battle menus and messages
- Player vs Enemy collision detection and combat initiation

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
doc/
├── README.md                    # This file - documentation index
├── ANIMATION_SYSTEM.md          # Animation system guide
├── ATTACK_ANIMATIONS.md         # Attack animation system
├── BATTLE_FIXES.md             # Battle system improvements
├── FFT_IMPLEMENTATION.md        # Final Fantasy Tactics implementation guide
├── TACTICAL_ROADMAP.md          # Complete 6-phase tactical RPG roadmap
├── DEVELOPMENT_PRIORITIES.md    # Development priorities and success metrics
└── PHASE1_COMPLETE.md          # Phase 1 completion summary
```

---

*Last updated: October 11, 2025*