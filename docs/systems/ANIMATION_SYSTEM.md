# Animation System Documentation

## 🎬 **Enhanced Animation System**

Your RPG now has a flexible, multi-animation system that supports:
- **Multiple animations per character** (idle, walking, attacking, etc.)
- **Dynamic animation switching** based on game state
- **Flexible sprite sheet loading** with configurable frame ranges
- **Automatic fallback** to static sprites if animations fail to load

## 🔧 **Core Components**

### **AnimationState Enum**
```go
const (
    AnimationIdle AnimationState = iota
    AnimationWalking
    AnimationAttacking
    AnimationCasting
    AnimationDeath
)
```

### **AnimationConfig Structure**
```go
type AnimationConfig struct {
    State         AnimationState // Which animation this is (idle, walking, etc.)
    SpriteSheet   string         // Path to sprite sheet file
    StartFrame    int            // Starting frame index (0-based)
    FrameCount    int            // Number of frames (0 = use all frames)
    FrameDuration time.Duration  // How long each frame displays
    Loop          bool           // Whether animation repeats
}
```

### **CharacterAnimations Structure**
```go
type CharacterAnimations struct {
    Animations []AnimationConfig // List of all animations
    Scale      float64           // Rendering scale factor
    OffsetX    float64           // X offset for rendering
    OffsetY    float64           // Y offset for rendering
}
```

## 🎮 **Usage Examples**

### **Single Animation Character**
```go
// Simple way: just one animation
warrior := entities.CreateAnimatedPlayerWithSingleAnimation(
    "Conan", 100, 100, components.JobWarrior, 3,
    "assets/sprites/hero/hero-idle.png", 
    components.AnimationIdle
)
```

### **Multi-Animation Character**
```go
// Advanced way: multiple animations
heroAnimations := entities.CharacterAnimations{
    Animations: []entities.AnimationConfig{
        {
            State:         components.AnimationIdle,
            SpriteSheet:   "assets/sprites/hero/hero-idle.png",
            StartFrame:    0,
            FrameCount:    6, // Use first 6 frames
            FrameDuration: 200 * time.Millisecond,
            Loop:          true,
        },
        {
            State:         components.AnimationWalking,
            SpriteSheet:   "assets/sprites/hero/hero-walk.png",
            StartFrame:    0,
            FrameCount:    8, // Use first 8 frames
            FrameDuration: 150 * time.Millisecond,
            Loop:          true,
        },
        {
            State:         components.AnimationAttacking,
            SpriteSheet:   "assets/sprites/hero/hero-attack.png",
            StartFrame:    0,
            FrameCount:    4, // Attack has 4 frames
            FrameDuration: 100 * time.Millisecond,
            Loop:          false, // Don't loop attack animation
        },
    },
    Scale:   1.0,
    OffsetX: 0,
    OffsetY: 0,
}

warrior := entities.CreateAnimatedPlayerWithJob(
    "Conan", 100, 100, components.JobWarrior, 3, heroAnimations
)
```

## 🔄 **Animation State Management**

### **Runtime Animation Switching**
```go
// Get animation component
if animComp := player.Animation(); animComp != nil {
    // Switch to walking animation
    animComp.SetState(components.AnimationWalking)
    
    // Safe switching (only if animation exists)
    success := animComp.SetStateIfAvailable(components.AnimationAttacking)
    
    // Check what animations are available
    states := animComp.GetAvailableStates()
    
    // Check current state
    currentState := animComp.GetCurrentState()
}
```

### **Automatic Movement-Based Switching** ✨
The engine automatically switches between idle and walking animations:
- **Moving** (arrow keys pressed) → Walking animation (if available)
- **Stopped** → Idle animation
- **Fallback** → If walking animation doesn't exist, stays in idle

## 📁 **Sprite Sheet Organization**

### **Expected Format**
- **32x32 pixel sprites** arranged horizontally
- **PNG format** recommended
- **Transparent backgrounds** for best results

### **Example Directory Structure**
```
assets/sprites/
├── hero/
│   ├── hero-idle.png      # 6 frames of idle animation
│   ├── hero-walk.png      # 8 frames of walking
│   ├── hero-attack.png    # 4 frames of attack
│   └── hero-death.png     # 5 frames of death
├── enemies/
│   ├── goblin-idle.png
│   └── goblin-attack.png
└── player.png             # Fallback static sprite
```

## ⚡ **Performance Features**

### **Smart Loading**
- **Graceful Fallbacks**: If animation fails, falls back to static sprite
- **Partial Loading**: Can skip failed animations and load others
- **Memory Efficient**: Only loads frames that are actually used

### **Frame Management**
- **Configurable Timing**: Different frame rates for different animations
- **Loop Control**: Some animations loop (idle), others don't (attack)
- **State Caching**: Remembers animation progress when switching states

## 🎯 **Current Implementation**

Your hero "Conan" now uses:
- **6-frame idle animation** from `hero-idle.png` (200ms per frame, looping)
- **6-frame walking animation** from `hero-walk.png` (150ms per frame, looping)
- **Automatic state switching** based on movement:
  - Standing still → Idle animation
  - Moving with arrow keys → Walking animation
- **Fallback protection** if sprite sheets are missing

## 🚀 **Future Expansion Ideas**

You can easily add:
- **Directional animations** (walking left/right/up/down)
- **Combat animations** (different attacks per weapon type)
- **Emotional states** (happy, angry, hurt animations)
- **Environmental reactions** (swimming, climbing animations)
- **Spell casting** (different animations per magic school)

The system is designed to be highly extensible - just add more `AnimationConfig` entries to any character! 🎨