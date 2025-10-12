package components

import (
	"time"

	"github.com/jrecuero/myrpg/internal/gfx"
)

// AnimationState represents different animation states a character can be in
type AnimationState int

const (
	AnimationIdle AnimationState = iota
	AnimationWalking
	AnimationAttacking
	AnimationCasting
	AnimationDeath
)

func (as AnimationState) String() string {
	switch as {
	case AnimationIdle:
		return "idle"
	case AnimationWalking:
		return "walking"
	case AnimationAttacking:
		return "attacking"
	case AnimationCasting:
		return "casting"
	case AnimationDeath:
		return "death"
	default:
		return "unknown"
	}
}

// Animation represents a sequence of frames from a sprite sheet
type Animation struct {
	Frames        []*gfx.Sprite // Array of sprite frames
	FrameDuration time.Duration // Duration each frame should be displayed
	Loop          bool          // Whether the animation should loop
	CurrentFrame  int           // Current frame index
	LastFrameTime time.Time     // Time when the current frame started
}

// NewAnimation creates a new animation from a sprite sheet
func NewAnimation(frames []*gfx.Sprite, frameDuration time.Duration, loop bool) *Animation {
	return &Animation{
		Frames:        frames,
		FrameDuration: frameDuration,
		Loop:          loop,
		CurrentFrame:  0,
		LastFrameTime: time.Now(),
	}
}

// Update advances the animation to the next frame if enough time has passed
func (a *Animation) Update() {
	if len(a.Frames) <= 1 {
		return // No animation needed for single frame
	}

	now := time.Now()
	if now.Sub(a.LastFrameTime) >= a.FrameDuration {
		a.CurrentFrame++

		if a.CurrentFrame >= len(a.Frames) {
			if a.Loop {
				a.CurrentFrame = 0 // Loop back to first frame
			} else {
				a.CurrentFrame = len(a.Frames) - 1 // Stay on last frame
			}
		}

		a.LastFrameTime = now
	}
}

// GetCurrentSprite returns the current frame's sprite
func (a *Animation) GetCurrentSprite() *gfx.Sprite {
	if len(a.Frames) == 0 {
		return nil
	}
	if a.CurrentFrame >= len(a.Frames) {
		return a.Frames[len(a.Frames)-1]
	}
	return a.Frames[a.CurrentFrame]
}

// Reset resets the animation to the first frame
func (a *Animation) Reset() {
	a.CurrentFrame = 0
	a.LastFrameTime = time.Now()
}

// IsFinished returns true if the animation has finished (only relevant for non-looping animations)
func (a *Animation) IsFinished() bool {
	return !a.Loop && a.CurrentFrame >= len(a.Frames)-1
}

// AnimationComponent manages multiple animations for an entity
type AnimationComponent struct {
	Animations   map[AnimationState]*Animation // Map of animation states to animations
	CurrentState AnimationState                // Current animation state
	Scale        float64                       // Scale factor for rendering
	OffsetX      float64                       // X offset for rendering
	OffsetY      float64                       // Y offset for rendering
}

// NewAnimationComponent creates a new AnimationComponent
func NewAnimationComponent(scale, offsetX, offsetY float64) *AnimationComponent {
	return &AnimationComponent{
		Animations:   make(map[AnimationState]*Animation),
		CurrentState: AnimationIdle,
		Scale:        scale,
		OffsetX:      offsetX,
		OffsetY:      offsetY,
	}
}

// AddAnimation adds an animation for a specific state
func (ac *AnimationComponent) AddAnimation(state AnimationState, animation *Animation) {
	ac.Animations[state] = animation
}

// SetState changes the current animation state
func (ac *AnimationComponent) SetState(state AnimationState) {
	if state != ac.CurrentState {
		// Reset the new animation when switching states
		if animation, exists := ac.Animations[state]; exists {
			animation.Reset()
		}
		ac.CurrentState = state
	}
}

// Update updates the current animation
func (ac *AnimationComponent) Update() {
	if animation, exists := ac.Animations[ac.CurrentState]; exists {
		animation.Update()
	}
}

// GetCurrentSprite returns the current frame's sprite for the current state
func (ac *AnimationComponent) GetCurrentSprite() *gfx.Sprite {
	if animation, exists := ac.Animations[ac.CurrentState]; exists {
		return animation.GetCurrentSprite()
	}
	return nil
}

// GetCurrentAnimation returns the current animation
func (ac *AnimationComponent) GetCurrentAnimation() *Animation {
	if animation, exists := ac.Animations[ac.CurrentState]; exists {
		return animation
	}
	return nil
}

// HasAnimation checks if the component has an animation for the given state
func (ac *AnimationComponent) HasAnimation(state AnimationState) bool {
	_, exists := ac.Animations[state]
	return exists
}

// GetAvailableStates returns a slice of all available animation states
func (ac *AnimationComponent) GetAvailableStates() []AnimationState {
	states := make([]AnimationState, 0, len(ac.Animations))
	for state := range ac.Animations {
		states = append(states, state)
	}
	return states
}

// SetStateIfAvailable sets the animation state only if that animation exists
func (ac *AnimationComponent) SetStateIfAvailable(state AnimationState) bool {
	if ac.HasAnimation(state) {
		ac.SetState(state)
		return true
	}
	return false
}

// GetCurrentState returns the current animation state
func (ac *AnimationComponent) GetCurrentState() AnimationState {
	return ac.CurrentState
}
