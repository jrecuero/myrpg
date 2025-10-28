package save

// Package save provides game state persistence functionality.

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// EventState represents the persistent state of an event that needs to be saved
type EventState struct {
	ID            string                `json:"id"`             // Event ID
	State         components.EventState `json:"state"`          // Current event state (Active, Triggered, Completed, Disabled)
	TriggerCount  int                   `json:"trigger_count"`  // Number of times triggered
	MaxTriggers   int                   `json:"max_triggers"`   // Maximum trigger limit
	LastTriggered *time.Time            `json:"last_triggered"` // Last trigger timestamp (pointer for nil handling)
	Cooldown      time.Duration         `json:"cooldown"`       // Cooldown duration
	Prerequisites []string              `json:"prerequisites"`  // Required prerequisite events
	CanRepeat     bool                  `json:"can_repeat"`     // Whether event can repeat
	ActiveInMode  components.GameMode   `json:"active_in_mode"` // Which game mode event is active in

	// Additional metadata for debugging and analytics
	FirstTriggered *time.Time `json:"first_triggered,omitempty"` // When first triggered
	SavedAt        time.Time  `json:"saved_at"`                  // When this state was saved
}

// EventStateSaveData contains all event states and completion tracking
type EventStateSaveData struct {
	Version         int                    `json:"version"`          // Save format version for compatibility
	SavedAt         time.Time              `json:"saved_at"`         // When the save was created
	EventStates     map[string]*EventState `json:"event_states"`     // Event ID -> State mapping
	CompletedEvents map[string]bool        `json:"completed_events"` // Event completion tracking
	GameMode        components.GameMode    `json:"game_mode"`        // Current game mode when saved
}

// NewEventState creates a new EventState from an EventComponent
func NewEventState(eventComp *components.EventComponent) *EventState {
	es := &EventState{
		ID:            eventComp.ID,
		State:         eventComp.State,
		TriggerCount:  eventComp.TriggerCount,
		MaxTriggers:   eventComp.MaxTriggers,
		Cooldown:      eventComp.Cooldown,
		Prerequisites: make([]string, len(eventComp.Prerequisites)),
		CanRepeat:     eventComp.CanRepeat,
		ActiveInMode:  eventComp.ActiveInMode,
		SavedAt:       time.Now(),
	}

	// Copy prerequisites slice
	copy(es.Prerequisites, eventComp.Prerequisites)

	// Handle LastTriggered (use pointer to handle zero time)
	if !eventComp.LastTriggered.IsZero() {
		es.LastTriggered = &eventComp.LastTriggered
	}

	// Set FirstTriggered if this event has been triggered at least once
	if eventComp.TriggerCount > 0 && !eventComp.LastTriggered.IsZero() {
		// For simplicity, use LastTriggered as FirstTriggered if we don't have better data
		es.FirstTriggered = &eventComp.LastTriggered
	}

	return es
}

// ApplyToEventComponent applies saved state back to an EventComponent
func (es *EventState) ApplyToEventComponent(eventComp *components.EventComponent) {
	eventComp.State = es.State
	eventComp.TriggerCount = es.TriggerCount
	eventComp.MaxTriggers = es.MaxTriggers
	eventComp.Cooldown = es.Cooldown
	eventComp.CanRepeat = es.CanRepeat
	eventComp.ActiveInMode = es.ActiveInMode

	// Copy prerequisites
	eventComp.Prerequisites = make([]string, len(es.Prerequisites))
	copy(eventComp.Prerequisites, es.Prerequisites)

	// Handle LastTriggered
	if es.LastTriggered != nil {
		eventComp.LastTriggered = *es.LastTriggered
	} else {
		eventComp.LastTriggered = time.Time{} // Zero time
	}
}

// NewEventStateSaveData creates a new save data structure
func NewEventStateSaveData() *EventStateSaveData {
	return &EventStateSaveData{
		Version:         1, // Current save format version
		SavedAt:         time.Now(),
		EventStates:     make(map[string]*EventState),
		CompletedEvents: make(map[string]bool),
		GameMode:        components.GameModeExploration, // Default mode
	}
}

// ToJSON serializes the event save data to JSON
func (esd *EventStateSaveData) ToJSON() ([]byte, error) {
	return json.MarshalIndent(esd, "", "  ")
}

// FromJSON deserializes event save data from JSON
func (esd *EventStateSaveData) FromJSON(data []byte) error {
	return json.Unmarshal(data, esd)
}

// AddEventState adds or updates an event state
func (esd *EventStateSaveData) AddEventState(eventComp *components.EventComponent) {
	eventState := NewEventState(eventComp)
	esd.EventStates[eventComp.ID] = eventState

	// Update completion tracking
	if eventState.State == components.EventCompleted {
		esd.CompletedEvents[eventComp.ID] = true
	}

	esd.SavedAt = time.Now()
}

// GetEventState retrieves an event state by ID
func (esd *EventStateSaveData) GetEventState(eventID string) (*EventState, bool) {
	state, exists := esd.EventStates[eventID]
	return state, exists
}

// IsEventCompleted checks if an event is marked as completed
func (esd *EventStateSaveData) IsEventCompleted(eventID string) bool {
	return esd.CompletedEvents[eventID]
}

// SetEventCompleted marks an event as completed
func (esd *EventStateSaveData) SetEventCompleted(eventID string, completed bool) {
	esd.CompletedEvents[eventID] = completed
	esd.SavedAt = time.Now()
}

// GetCompletedEventCount returns the number of completed events
func (esd *EventStateSaveData) GetCompletedEventCount() int {
	return len(esd.CompletedEvents)
}

// GetTotalEventCount returns the total number of events tracked
func (esd *EventStateSaveData) GetTotalEventCount() int {
	return len(esd.EventStates)
}

// Validate checks the integrity of the save data
func (esd *EventStateSaveData) Validate() error {
	if esd.Version <= 0 {
		return fmt.Errorf("invalid save version: %d", esd.Version)
	}

	if esd.EventStates == nil {
		return fmt.Errorf("event states map is nil")
	}

	if esd.CompletedEvents == nil {
		return fmt.Errorf("completed events map is nil")
	}

	// Validate that completed events are consistent
	for eventID := range esd.CompletedEvents {
		if state, exists := esd.EventStates[eventID]; exists {
			if state.State != components.EventCompleted {
				return fmt.Errorf("event %s marked as completed but state is %v", eventID, state.State)
			}
		}
	}

	return nil
}

// GetStatistics returns useful statistics about the event save data
func (esd *EventStateSaveData) GetStatistics() map[string]interface{} {
	stats := map[string]interface{}{
		"version":          esd.Version,
		"saved_at":         esd.SavedAt,
		"total_events":     len(esd.EventStates),
		"completed_events": len(esd.CompletedEvents),
		"game_mode":        esd.GameMode.String(),
	}

	// Count events by state
	stateCounts := make(map[string]int)
	for _, state := range esd.EventStates {
		stateStr := fmt.Sprintf("%d", int(state.State)) // Convert to string for JSON
		stateCounts[stateStr]++
	}
	stats["state_counts"] = stateCounts

	// Count events with cooldowns
	eventsWithCooldowns := 0
	eventsWithPrereqs := 0
	for _, state := range esd.EventStates {
		if state.Cooldown > 0 {
			eventsWithCooldowns++
		}
		if len(state.Prerequisites) > 0 {
			eventsWithPrereqs++
		}
	}
	stats["events_with_cooldowns"] = eventsWithCooldowns
	stats["events_with_prerequisites"] = eventsWithPrereqs

	return stats
}
