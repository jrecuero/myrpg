# Attack Animation System

The RPG now includes a flexible attack animation system that provides visual feedback when heroes attack enemies during battle. This system supports customizable durations and automatic animation switching.

## Features

### Sword Attack Animation
- **Sprite Sheet**: `assets/sprites/hero/hero-sword.png` 
- **Frames**: 6 frames at 32x32 pixels each
- **Animation**: Fast-paced attack sequence (100ms per frame)
- **Looping**: Continuously loops during the attack period then returns to previous state

### Configurable Duration
The attack animation can be configured to last for any duration, allowing you to customize the visual feedback:

```go
// Short attack feedback (0.8 seconds)
game.SetAttackAnimationDuration(800 * time.Millisecond)

// Default duration (1.5 seconds)
game.SetAttackAnimationDuration(1500 * time.Millisecond)

// Long attack feedback (2.5 seconds)
game.SetAttackAnimationDuration(2500 * time.Millisecond)
```

### Automatic State Management
- **Triggers**: Attack animation starts when player attacks an enemy
- **Duration**: Stays in attack animation for the configured duration
- **Smart Revert**: Automatically returns to idle state after battle (never walking)
- **Battle Logic**: Forces idle state after combat regardless of movement during collision
- **Non-Interrupting**: Animation switching based on movement still works after attack animation ends

## How It Works

1. **Battle Trigger**: When a player selects and executes an attack in battle
2. **Animation Switch**: Hero switches to attack animation using the sword sprite sheet
3. **Timed Duration**: Animation stays active for the configured duration (default 1.5 seconds)
4. **Smart Return**: After the duration expires, hero returns to idle state (not walking)
5. **Battle Logic**: Ensures proper post-battle state regardless of movement during collision
6. **Continued Gameplay**: Normal idle/walking animation switching resumes

## Technical Implementation

### Animation Component Enhancement
- **Temporary States**: Support for time-limited animation states
- **State Memory**: Remembers previous state for auto-revert
- **Timer System**: Built-in timing for temporary animations
- **Smart Revert**: `SetTemporaryStateWithRevertTo()` method for controlled state transitions
- **Battle Logic**: Forces idle state after attack regardless of collision-time animation

### Battle System Integration
- **Attack Trigger**: Automatically triggers attack animation on player attacks
- **Visual Feedback**: Provides immediate visual confirmation of attack actions
- **Customizable Timing**: Game developers can adjust animation duration

## Usage Examples

### Basic Setup (Already configured in main.go)
```go
// Hero with attack animation
heroAnimations := entities.CharacterAnimations{
    Animations: []entities.AnimationConfig{
        // ... idle and walking animations ...
        {
            State:         components.AnimationAttacking,
            SpriteSheet:   "assets/sprites/hero/hero-sword.png",
            StartFrame:    0,
            FrameCount:    0,  // Use all 6 frames
            FrameDuration: 100 * time.Millisecond,
            Loop:          true, // Loop attack animation during attack period
        },
    },
    // ... other config ...
}
```

### Customizing Animation Duration
```go
// Configure different durations for different gameplay styles
game.SetAttackAnimationDuration(1200 * time.Millisecond) // 1.2 seconds
```

## Asset Requirements

To add attack animations to other characters:
1. Create a sprite sheet with attack frames (32x32 pixels recommended)
2. Add AnimationAttacking configuration to character animations
3. Sprite sheet should contain sequential attack frames
4. Set Loop: false for attack animations

## Benefits

- **Visual Feedback**: Players get immediate confirmation of their attacks
- **Enhanced Experience**: More engaging combat with animated feedback  
- **Flexible Timing**: Developers can adjust animation duration to match game pacing
- **Seamless Integration**: Works with existing animation and battle systems
- **Easy Configuration**: Simple API for customizing behavior