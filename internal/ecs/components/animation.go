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
	Animations        map[AnimationState]*Animation // Map of animation states to animations
	CurrentState      AnimationState                // Current animation state
	PreviousState     AnimationState                // Previous animation state (for reverting)
	Scale             float64                       // Scale factor for rendering
	OffsetX           float64                       // X offset for rendering
	OffsetY           float64                       // Y offset for rendering
	TempStateTimer    time.Time                     // Timer for temporary animation states
	TempStateDuration time.Duration                 // Duration for temporary animation
	IsTemporaryState  bool                          // Whether current state is temporary
}

// NewAnimationComponent creates a new AnimationComponent
func NewAnimationComponent(scale, offsetX, offsetY float64) *AnimationComponent {
	return &AnimationComponent{
		Animations:       make(map[AnimationState]*Animation),
		CurrentState:     AnimationIdle,
		PreviousState:    AnimationIdle,
		Scale:            scale,
		OffsetX:          offsetX,
		OffsetY:          offsetY,
		IsTemporaryState: false,
	}
}

// AddAnimation adds an animation for a specific state
func (ac *AnimationComponent) AddAnimation(state AnimationState, animation *Animation) {
	ac.Animations[state] = animation
}

// SetState changes the current animation state
func (ac *AnimationComponent) SetState(state AnimationState) {
	if state != ac.CurrentState {
		// Store previous state for potential revert
		ac.PreviousState = ac.CurrentState

		// Reset the new animation when switching states
		if animation, exists := ac.Animations[state]; exists {
			animation.Reset()
		}
		ac.CurrentState = state
		ac.IsTemporaryState = false // Clear temporary state flag
	}
}

// SetTemporaryState changes to a temporary animation state for a specified duration
func (ac *AnimationComponent) SetTemporaryState(state AnimationState, duration time.Duration) {
	if ac.HasAnimation(state) {
		// Store previous state for reverting
		if !ac.IsTemporaryState {
			ac.PreviousState = ac.CurrentState
		}

		// Set new temporary state
		if animation, exists := ac.Animations[state]; exists {
			animation.Reset()
		}
		ac.CurrentState = state
		ac.IsTemporaryState = true
		ac.TempStateTimer = time.Now()
		ac.TempStateDuration = duration
	}
}

// SetTemporaryStateWithRevertTo changes to a temporary animation state and specifies what state to revert to
func (ac *AnimationComponent) SetTemporaryStateWithRevertTo(state AnimationState, duration time.Duration, revertToState AnimationState) {
	if ac.HasAnimation(state) {
		// Set specific revert state instead of current state
		ac.PreviousState = revertToState

		// Set new temporary state
		if animation, exists := ac.Animations[state]; exists {
			animation.Reset()
		}
		ac.CurrentState = state
		ac.IsTemporaryState = true
		ac.TempStateTimer = time.Now()
		ac.TempStateDuration = duration
	}
}

// Update updates the current animation and handles temporary state expiration
func (ac *AnimationComponent) Update() {
	// Check if temporary state has expired
	if ac.IsTemporaryState && time.Since(ac.TempStateTimer) >= ac.TempStateDuration {
		// Revert to previous state
		ac.SetState(ac.PreviousState)
	}

	// Update current animation
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

// IsCurrentAnimationFinished returns true if the current animation has finished
func (ac *AnimationComponent) IsCurrentAnimationFinished() bool {
	if animation, exists := ac.Animations[ac.CurrentState]; exists {
		return animation.IsFinished()
	}
	return false
}

// IsInTemporaryState returns true if currently in a temporary animation state
func (ac *AnimationComponent) IsInTemporaryState() bool {
	return ac.IsTemporaryState
}
