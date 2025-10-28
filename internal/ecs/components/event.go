package components

import (
	"time"

	"github.com/jrecuero/myrpg/internal/constants"
)

// TriggerCondition defines when an event should be triggered
type TriggerCondition int

const (
	TriggerOnTouch      TriggerCondition = iota // Triggered when player touches the event entity
	TriggerOnTimeout                            // Triggered after a specified duration
	TriggerOnRoomEntry                          // Triggered when player enters the room/area
	TriggerOnProximity                          // Triggered when player is within range
	TriggerOnInteract                           // Triggered when player explicitly interacts (e.g., press key)
	TriggerOnQuestState                         // Triggered based on quest completion state
	TriggerManual                               // Triggered manually by game logic
)

// EventType defines the type of event that will be executed
type EventType int

const (
	EventBattle   EventType = iota // Start a tactical battle
	EventDialog                    // Open dialog with NPC
	EventChest                     // Open a chest/container
	EventDoor                      // Move to another area/room
	EventTrap                      // Trigger a trap
	EventInfo                      // Display information to player
	EventQuest                     // Trigger quest-related actions
	EventCutscene                  // Play a cutscene or story sequence
	EventShop                      // Open shop interface
	EventRest                      // Trigger rest/save point
)

// EventState tracks the current state of an event
type EventState int

const (
	EventActive    EventState = iota // Event is active and can be triggered
	EventTriggered                   // Event has been triggered but may repeat
	EventCompleted                   // Event is completed and won't trigger again
	EventDisabled                    // Event is temporarily disabled
)

// GameMode represents the game mode context where events should be active
type GameMode int

const (
	GameModeExploration GameMode = iota // Event active in exploration mode
	GameModeTactical                    // Event active in tactical mode
	GameModeBoth                        // Event active in both modes
)

// EventData contains type-specific data for different event types
type EventData struct {
	// Battle event data
	Enemies   []string `json:"enemies,omitempty"`    // Enemy IDs to spawn in battle
	BattleMap string   `json:"battle_map,omitempty"` // Battle map to use

	// Dialog event data
	NPCID    string `json:"npc_id,omitempty"`    // NPC to talk to
	DialogID string `json:"dialog_id,omitempty"` // Dialog tree to start

	// Chest event data
	Items    []string `json:"items,omitempty"`     // Item IDs in chest
	Gold     int      `json:"gold,omitempty"`      // Gold amount
	IsLocked bool     `json:"is_locked,omitempty"` // Whether chest is locked

	// Door event data
	TargetMap string  `json:"target_map,omitempty"` // Map to transition to
	TargetX   float64 `json:"target_x,omitempty"`   // Target X position
	TargetY   float64 `json:"target_y,omitempty"`   // Target Y position

	// Trap event data
	TrapType string `json:"trap_type,omitempty"` // Type of trap
	Damage   int    `json:"damage,omitempty"`    // Damage amount

	// Info event data
	Title   string `json:"title,omitempty"`   // Info title
	Message string `json:"message,omitempty"` // Info message

	// Quest event data
	QuestID string `json:"quest_id,omitempty"` // Quest to start/complete

	// Shop event data
	ShopID string `json:"shop_id,omitempty"` // Shop inventory ID
}

// EventConditionData contains data for different trigger conditions
type EventConditionData struct {
	Duration   time.Duration `json:"duration,omitempty"`    // For timeout triggers
	Distance   float64       `json:"distance,omitempty"`    // For proximity triggers
	QuestState string        `json:"quest_state,omitempty"` // For quest state triggers
}

// EventComponent represents an interactive event in the game world
type EventComponent struct {
	ID               string             // Unique identifier for this event
	Name             string             // Display name for this event
	TriggerCondition TriggerCondition   // When this event should be triggered
	EventType        EventType          // Type of event to execute
	State            EventState         // Current state of the event
	CanRepeat        bool               // Whether the event can be triggered multiple times
	EventData        EventData          // Type-specific event data
	ConditionData    EventConditionData // Data for trigger conditions

	// View context - controls which game mode this event is active in
	ActiveInMode GameMode // Which game mode(s) this event is active in

	// Prerequisite system for event dependencies
	Prerequisites []string      // List of event IDs that must be completed before this event can trigger
	TriggerCount  int           // Number of times this event has been triggered
	MaxTriggers   int           // Maximum number of times this event can trigger (0 = unlimited)
	LastTriggered time.Time     // When this event was last triggered
	Cooldown      time.Duration // Minimum time between triggers

	// Visual properties for event representation
	Visible       bool       // Whether the event is visible (false for hidden/traps)
	SpritePath    string     // Path to custom sprite (optional)
	FallbackColor [3]uint8   // RGB color for colored square fallback
	SpriteScale   float64    // Scale factor for sprite rendering
	SpriteOffset  [2]float64 // X,Y offset for sprite positioning
}

// NewEventComponent creates a new event component with default values
func NewEventComponent(id, name string, triggerCondition TriggerCondition, eventType EventType) *EventComponent {
	// Set default color based on event type using constants
	var defaultColor [3]uint8
	switch eventType {
	case EventBattle:
		defaultColor = constants.EventColorBattle
	case EventDialog:
		defaultColor = constants.EventColorDialog
	case EventChest:
		defaultColor = constants.EventColorChest
	case EventDoor:
		defaultColor = constants.EventColorDoor
	case EventTrap:
		defaultColor = constants.EventColorTrap
	case EventInfo:
		defaultColor = constants.EventColorInfo
	case EventQuest:
		defaultColor = constants.EventColorQuest
	case EventCutscene:
		defaultColor = constants.EventColorCutscene
	case EventShop:
		defaultColor = constants.EventColorShop
	case EventRest:
		defaultColor = constants.EventColorRest
	default:
		defaultColor = [3]uint8{128, 128, 128} // Gray for unknown
	}

	return &EventComponent{
		ID:               id,
		Name:             name,
		TriggerCondition: triggerCondition,
		EventType:        eventType,
		State:            EventActive,
		CanRepeat:        false,
		EventData:        EventData{},
		ConditionData:    EventConditionData{},
		ActiveInMode:     GameModeExploration, // Default to exploration mode

		// Prerequisites and cooldown system
		Prerequisites: make([]string, 0),
		TriggerCount:  0,
		MaxTriggers:   0,           // 0 means unlimited
		LastTriggered: time.Time{}, // Zero time
		Cooldown:      0,           // No cooldown by default

		// Visual properties
		Visible:       true, // Events are visible by default
		SpritePath:    "",   // No custom sprite by default
		FallbackColor: defaultColor,
		SpriteScale:   1.0, // Normal scaling
		SpriteOffset:  [2]float64{0, 0},
	}
}

// IsActiveInMode checks if the event should be active in the given game mode
func (e *EventComponent) IsActiveInMode(mode GameMode) bool {
	return e.ActiveInMode == GameModeBoth || e.ActiveInMode == mode
}

// SetActiveInMode sets which game mode(s) this event is active in
func (e *EventComponent) SetActiveInMode(mode GameMode) *EventComponent {
	e.ActiveInMode = mode
	return e
}

// CanTrigger checks if the event can be triggered based on its current state
func (e *EventComponent) CanTrigger() bool {
	// Check basic state
	if e.State == EventDisabled || e.State == EventCompleted {
		return false
	}

	// Check if not repeatable and already triggered
	if !e.CanRepeat && e.State == EventTriggered {
		return false
	}

	// Check trigger limits
	if e.MaxTriggers > 0 && e.TriggerCount >= e.MaxTriggers {
		return false
	}

	// Check cooldown
	if e.Cooldown > 0 && !e.LastTriggered.IsZero() {
		if time.Since(e.LastTriggered) < e.Cooldown {
			return false
		}
	}

	return true
}

// Trigger marks the event as triggered and updates relevant fields
func (e *EventComponent) Trigger() {
	if !e.CanTrigger() {
		return
	}

	e.State = EventTriggered
	e.TriggerCount++
	e.LastTriggered = time.Now()

	// If not repeatable or hit max triggers, mark as completed
	if !e.CanRepeat || (e.MaxTriggers > 0 && e.TriggerCount >= e.MaxTriggers) {
		e.State = EventCompleted
	}
}

// Reset resets the event to its initial state
func (e *EventComponent) Reset() {
	e.State = EventActive
	e.TriggerCount = 0
	e.LastTriggered = time.Time{}
}

// SetRepeatable sets whether the event can be triggered multiple times
func (e *EventComponent) SetRepeatable(repeatable bool) *EventComponent {
	e.CanRepeat = repeatable
	return e
}

// SetEventData sets the event-specific data
func (e *EventComponent) SetEventData(data EventData) *EventComponent {
	e.EventData = data
	return e
}

// SetConditionData sets the condition-specific data
func (e *EventComponent) SetConditionData(data EventConditionData) *EventComponent {
	e.ConditionData = data
	return e
}

// GetTriggerConditionName returns a human-readable name for the trigger condition
func (e *EventComponent) GetTriggerConditionName() string {
	switch e.TriggerCondition {
	case TriggerOnTouch:
		return "On Touch"
	case TriggerOnTimeout:
		return "On Timeout"
	case TriggerOnRoomEntry:
		return "On Room Entry"
	case TriggerOnProximity:
		return "On Proximity"
	case TriggerOnInteract:
		return "On Interact"
	case TriggerOnQuestState:
		return "On Quest State"
	case TriggerManual:
		return "Manual"
	default:
		return "Unknown"
	}
}

// GetEventTypeName returns a human-readable name for the event type
func (e *EventComponent) GetEventTypeName() string {
	switch e.EventType {
	case EventBattle:
		return "Battle"
	case EventDialog:
		return "Dialog"
	case EventChest:
		return "Chest"
	case EventDoor:
		return "Door"
	case EventTrap:
		return "Trap"
	case EventInfo:
		return "Info"
	case EventQuest:
		return "Quest"
	case EventCutscene:
		return "Cutscene"
	case EventShop:
		return "Shop"
	case EventRest:
		return "Rest"
	default:
		return "Unknown"
	}
}

// SetVisible sets whether the event should be visible
func (e *EventComponent) SetVisible(visible bool) *EventComponent {
	e.Visible = visible
	return e
}

// SetSprite sets a custom sprite for the event
func (e *EventComponent) SetSprite(spritePath string) *EventComponent {
	e.SpritePath = spritePath
	return e
}

// SetFallbackColor sets the fallback color for when no sprite is available
func (e *EventComponent) SetFallbackColor(r, g, b uint8) *EventComponent {
	e.FallbackColor = [3]uint8{r, g, b}
	return e
}

// SetSpriteScale sets the scaling factor for the sprite
func (e *EventComponent) SetSpriteScale(scale float64) *EventComponent {
	e.SpriteScale = scale
	return e
}

// SetSpriteOffset sets the sprite positioning offset
func (e *EventComponent) SetSpriteOffset(x, y float64) *EventComponent {
	e.SpriteOffset = [2]float64{x, y}
	return e
}

// SetPrerequisites sets the event prerequisites
func (e *EventComponent) SetPrerequisites(prerequisites ...string) *EventComponent {
	e.Prerequisites = make([]string, len(prerequisites))
	copy(e.Prerequisites, prerequisites)
	return e
}

// SetMaxTriggers sets the maximum number of times this event can be triggered
func (e *EventComponent) SetMaxTriggers(maxTriggers int) *EventComponent {
	e.MaxTriggers = maxTriggers
	return e
}

// SetCooldown sets the cooldown duration between triggers
func (e *EventComponent) SetCooldown(cooldown time.Duration) *EventComponent {
	e.Cooldown = cooldown
	return e
}

// IsVisible returns whether the event is currently visible
func (e *EventComponent) IsVisible() bool {
	return e.Visible
}

// HasCustomSprite returns whether the event has a custom sprite set
func (e *EventComponent) HasCustomSprite() bool {
	return e.SpritePath != ""
}

// ShouldTriggerOnCollision returns true if this event should trigger on collision
func (e *EventComponent) ShouldTriggerOnCollision() bool {
	switch e.TriggerCondition {
	case TriggerOnTouch, TriggerOnProximity, TriggerOnInteract:
		return true
	case TriggerOnTimeout, TriggerOnRoomEntry, TriggerOnQuestState, TriggerManual:
		return false
	default:
		return false
	}
}
