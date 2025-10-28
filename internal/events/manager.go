// Package events provides event management and processing functionality.
// The EventManager handles registration, trigger detection, and execution of game events.
package events

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// EventResult represents the result of an event execution
type EventResult struct {
	Success    bool                   `json:"success"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
	NextAction string                 `json:"next_action,omitempty"` // For chaining events
}

// EventHandler is a function type for handling specific event types
type EventHandler func(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult

// EventManager manages all game events and their execution
type EventManager struct {
	entities        []*ecs.Entity                         // All entities that might have events
	handlers        map[components.EventType]EventHandler // Event type handlers
	activeEvents    map[string]*components.EventComponent // Currently active events by ID
	playerEntity    *ecs.Entity                           // Reference to player entity
	completedEvents map[string]bool                       // Track completed events
	eventHistory    []EventExecutionRecord                // History of executed events
	currentGameMode components.GameMode                   // Current game mode for event filtering
	collidingEvents map[string]bool                       // Track events currently colliding with player
}

// EventExecutionRecord tracks when and how events were executed
type EventExecutionRecord struct {
	EventID     string                      `json:"event_id"`
	EntityID    string                      `json:"entity_id"`
	Timestamp   time.Time                   `json:"timestamp"`
	Result      *EventResult                `json:"result"`
	TriggerType components.TriggerCondition `json:"trigger_type"`
}

// NewEventManager creates a new event manager
func NewEventManager() *EventManager {
	return &EventManager{
		entities:        make([]*ecs.Entity, 0),
		handlers:        make(map[components.EventType]EventHandler),
		activeEvents:    make(map[string]*components.EventComponent),
		completedEvents: make(map[string]bool),
		eventHistory:    make([]EventExecutionRecord, 0),
		collidingEvents: make(map[string]bool),
	}
}

// SetPlayer sets the player entity reference
func (em *EventManager) SetPlayer(player *ecs.Entity) {
	em.playerEntity = player
}

// SetGameMode sets the current game mode for event filtering
func (em *EventManager) SetGameMode(mode components.GameMode) {
	em.currentGameMode = mode
}

// RegisterEntity registers an entity that may contain events
func (em *EventManager) RegisterEntity(entity *ecs.Entity) {
	// Check if entity already registered
	for _, e := range em.entities {
		if e.ID == entity.ID {
			return
		}
	}

	em.entities = append(em.entities, entity)

	// If entity has an event component, add it to active events
	if eventComp := entity.Event(); eventComp != nil {
		em.activeEvents[eventComp.ID] = eventComp
	}
}

// UnregisterEntity removes an entity from event management
func (em *EventManager) UnregisterEntity(entity *ecs.Entity) {
	// Remove from entities slice
	for i, e := range em.entities {
		if e.ID == entity.ID {
			em.entities = append(em.entities[:i], em.entities[i+1:]...)
			break
		}
	}

	// Remove from active events if it has an event component
	if eventComp := entity.Event(); eventComp != nil {
		delete(em.activeEvents, eventComp.ID)
	}
}

// RegisterHandler registers an event handler for a specific event type
func (em *EventManager) RegisterHandler(eventType components.EventType, handler EventHandler) {
	em.handlers[eventType] = handler
}

// Update processes all events and checks for trigger conditions
func (em *EventManager) Update(deltaTime float64) {
	if em.playerEntity == nil {
		return
	}

	playerTransform := em.playerEntity.Transform()
	if playerTransform == nil {
		return
	}

	currentlyColliding := make(map[string]bool)

	// Check all entities with events
	for _, entity := range em.entities {
		eventComp := entity.Event()
		if eventComp == nil || !eventComp.CanTrigger() {
			continue
		}

		// Check if event is active in current game mode
		if !eventComp.IsActiveInMode(em.currentGameMode) {
			continue
		}

		// Check if prerequisites are met
		if !em.checkPrerequisites(eventComp) {
			continue
		}

		// Check trigger condition
		isTriggering := em.checkTriggerCondition(entity, eventComp, playerTransform, deltaTime)

		if isTriggering {
			currentlyColliding[eventComp.ID] = true

			// Only execute if this is a new collision or if it's a battle event (which should always trigger)
			wasAlreadyColliding := em.collidingEvents[eventComp.ID]

			if !wasAlreadyColliding || eventComp.EventType == components.EventBattle {
				em.executeEvent(entity, eventComp)
			}
		}
	}

	// Update collision state for next frame
	em.collidingEvents = currentlyColliding
}

// checkPrerequisites checks if all prerequisite events have been completed
func (em *EventManager) checkPrerequisites(eventComp *components.EventComponent) bool {
	// Check if the event has prerequisites defined
	if len(eventComp.Prerequisites) == 0 {
		return true
	}

	// Check if all prerequisite events have been completed
	for _, prereqID := range eventComp.Prerequisites {
		if !em.completedEvents[prereqID] {
			return false
		}
	}

	return true
}

// checkTriggerCondition evaluates whether an event's trigger condition is met
func (em *EventManager) checkTriggerCondition(entity *ecs.Entity, eventComp *components.EventComponent, playerTransform *components.Transform, deltaTime float64) bool {
	entityTransform := entity.Transform()

	switch eventComp.TriggerCondition {
	case components.TriggerOnTouch:
		return em.checkTouchTrigger(entityTransform, playerTransform)

	case components.TriggerOnProximity:
		return em.checkProximityTrigger(entityTransform, playerTransform, eventComp.ConditionData.Distance)

	case components.TriggerOnTimeout:
		return em.checkTimeoutTrigger(eventComp, deltaTime)

	case components.TriggerOnRoomEntry:
		// For now, we'll consider room entry as being within a certain area
		return em.checkRoomEntryTrigger(entityTransform, playerTransform)

	case components.TriggerOnInteract:
		// This would be triggered externally when player presses interact key
		return false

	case components.TriggerOnQuestState:
		return em.checkQuestStateTrigger(eventComp)

	case components.TriggerManual:
		// Manual triggers are handled externally
		return false
	}

	return false
}

// checkTouchTrigger checks if player is touching/colliding with the event entity
func (em *EventManager) checkTouchTrigger(entityTransform, playerTransform *components.Transform) bool {
	if entityTransform == nil || playerTransform == nil {
		return false
	}

	// Simple AABB collision detection
	return playerTransform.X < entityTransform.X+float64(entityTransform.Width) &&
		playerTransform.X+float64(playerTransform.Width) > entityTransform.X &&
		playerTransform.Y < entityTransform.Y+float64(entityTransform.Height) &&
		playerTransform.Y+float64(playerTransform.Height) > entityTransform.Y
}

// checkProximityTrigger checks if player is within range of the event entity
func (em *EventManager) checkProximityTrigger(entityTransform, playerTransform *components.Transform, triggerRange float64) bool {
	if entityTransform == nil || playerTransform == nil {
		return false
	}

	if triggerRange <= 0 {
		triggerRange = 32 // Default range
	}

	// Calculate distance between centers
	entityCenterX := entityTransform.X + float64(entityTransform.Width)/2
	entityCenterY := entityTransform.Y + float64(entityTransform.Height)/2
	playerCenterX := playerTransform.X + float64(playerTransform.Width)/2
	playerCenterY := playerTransform.Y + float64(playerTransform.Height)/2

	dx := entityCenterX - playerCenterX
	dy := entityCenterY - playerCenterY
	distance := math.Sqrt(dx*dx + dy*dy)

	return distance <= triggerRange
}

// checkTimeoutTrigger checks if the timeout duration has elapsed
func (em *EventManager) checkTimeoutTrigger(eventComp *components.EventComponent, deltaTime float64) bool {
	_ = deltaTime
	if eventComp.ConditionData.Duration <= 0 {
		return false
	}

	// For timeout triggers, we need to track when they started
	// This is a simplified version - in a full implementation, you might want to track start times
	// TODO: Add proper timeout tracking back if needed
	// For now, always return false for timeout triggers
	return false
}

// checkRoomEntryTrigger checks if player has entered the room/area
func (em *EventManager) checkRoomEntryTrigger(entityTransform, playerTransform *components.Transform) bool {
	// For now, treat room entry similar to proximity but with a larger range
	return em.checkProximityTrigger(entityTransform, playerTransform, 64)
}

// checkQuestStateTrigger checks if the required quest state is met
func (em *EventManager) checkQuestStateTrigger(eventComp *components.EventComponent) bool {
	_ = eventComp
	// This would need integration with the quest system
	// For now, return false - implement when integrating with quests
	return false
}

// executeEvent executes an event and calls the appropriate handler
func (em *EventManager) executeEvent(entity *ecs.Entity, eventComp *components.EventComponent) {
	handler, exists := em.handlers[eventComp.EventType]
	if !exists {
		log.Printf("No handler registered for event type: %s", eventComp.GetEventTypeName())
		return
	}

	// Execute the event handler
	result := handler(entity, eventComp, em.playerEntity)

	// Update event state
	eventComp.Trigger()

	// Record execution
	record := EventExecutionRecord{
		EventID:     eventComp.ID,
		EntityID:    entity.GetID(),
		Timestamp:   time.Now(),
		Result:      result,
		TriggerType: eventComp.TriggerCondition,
	}
	em.eventHistory = append(em.eventHistory, record)

	// Mark as completed if it's not repeatable
	if !eventComp.CanRepeat || eventComp.State == components.EventCompleted {
		em.completedEvents[eventComp.ID] = true
		delete(em.activeEvents, eventComp.ID)
	}

	log.Printf("Event executed: %s (%s) - Success: %t, Message: %s",
		eventComp.Name, eventComp.GetEventTypeName(), result.Success, result.Message)
}

// TriggerManualEvent manually triggers an event by ID
func (em *EventManager) TriggerManualEvent(eventID string) *EventResult {
	eventComp, exists := em.activeEvents[eventID]
	if !exists {
		return &EventResult{
			Success: false,
			Message: fmt.Sprintf("Event not found: %s", eventID),
		}
	}

	// Find the entity with this event
	var eventEntity *ecs.Entity
	for _, entity := range em.entities {
		if entity.Event() != nil && entity.Event().ID == eventID {
			eventEntity = entity
			break
		}
	}

	if eventEntity == nil {
		return &EventResult{
			Success: false,
			Message: fmt.Sprintf("Entity not found for event: %s", eventID),
		}
	}

	handler, exists := em.handlers[eventComp.EventType]
	if !exists {
		return &EventResult{
			Success: false,
			Message: fmt.Sprintf("No handler for event type: %s", eventComp.GetEventTypeName()),
		}
	}

	result := handler(eventEntity, eventComp, em.playerEntity)
	eventComp.Trigger()

	return result
}

// GetActiveEvents returns all currently active events
func (em *EventManager) GetActiveEvents() map[string]*components.EventComponent {
	return em.activeEvents
}

// GetEventHistory returns the event execution history
func (em *EventManager) GetEventHistory() []EventExecutionRecord {
	return em.eventHistory
}

// GetEventsInRange returns all events within a certain range of a position
func (em *EventManager) GetEventsInRange(x, y, radius float64) []*ecs.Entity {
	var eventsInRange []*ecs.Entity

	for _, entity := range em.entities {
		eventComp := entity.Event()
		transform := entity.Transform()
		if eventComp == nil || transform == nil {
			continue
		}

		// Calculate distance
		entityCenterX := transform.X + float64(transform.Width)/2
		entityCenterY := transform.Y + float64(transform.Height)/2

		dx := entityCenterX - x
		dy := entityCenterY - y
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance <= radius {
			eventsInRange = append(eventsInRange, entity)
		}
	}

	return eventsInRange
}

// GetEventsByType returns all events of a specific type
func (em *EventManager) GetEventsByType(eventType components.EventType) []*ecs.Entity {
	var events []*ecs.Entity

	for _, entity := range em.entities {
		if eventComp := entity.Event(); eventComp != nil && eventComp.EventType == eventType {
			events = append(events, entity)
		}
	}

	return events
}

// SortEventsByPriority sorts a slice of event entities by priority (highest first)
func (em *EventManager) SortEventsByPriority(events []*ecs.Entity) {
	sort.Slice(events, func(i, j int) bool {
		eventI := events[i].Event()
		eventJ := events[j].Event()
		if eventI == nil || eventJ == nil {
			return false
		}
		// TODO: Add priority support back if needed
		return false
	})
}

// ClearCompletedEvents removes all completed events from tracking
func (em *EventManager) ClearCompletedEvents() {
	em.completedEvents = make(map[string]bool)
	em.eventHistory = make([]EventExecutionRecord, 0)
}

// IsEventCompleted checks if an event has been completed
func (em *EventManager) IsEventCompleted(eventID string) bool {
	return em.completedEvents[eventID]
}

// GetEventStats returns statistics about events
func (em *EventManager) GetEventStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["total_entities"] = len(em.entities)
	stats["active_events"] = len(em.activeEvents)
	stats["completed_events"] = len(em.completedEvents)
	stats["execution_history"] = len(em.eventHistory)

	// Count events by type
	eventsByType := make(map[string]int)
	for _, entity := range em.entities {
		if eventComp := entity.Event(); eventComp != nil {
			typeName := eventComp.GetEventTypeName()
			eventsByType[typeName]++
		}
	}
	stats["events_by_type"] = eventsByType

	return stats
}

// GetGameMode returns the current game mode
func (em *EventManager) GetGameMode() components.GameMode {
	return em.currentGameMode
}

// SetEventCompleted marks an event as completed in the completion tracking
func (em *EventManager) SetEventCompleted(eventID string, completed bool) {
	em.completedEvents[eventID] = completed
	log.Printf("Event %s completion status set to %v", eventID, completed)
}

// GetCompletedEvents returns a copy of the completed events map
func (em *EventManager) GetCompletedEvents() map[string]bool {
	completed := make(map[string]bool)
	for id, status := range em.completedEvents {
		completed[id] = status
	}
	return completed
}

// LoadCompletedEvents sets the completed events from saved data
func (em *EventManager) LoadCompletedEvents(completedEvents map[string]bool) {
	em.completedEvents = make(map[string]bool)
	for id, status := range completedEvents {
		em.completedEvents[id] = status
	}
	log.Printf("Loaded %d completed event states", len(completedEvents))
}

// ClearCompletedEventsForReset removes all completed events from tracking (for new game)
func (em *EventManager) ClearCompletedEventsForReset() {
	em.completedEvents = make(map[string]bool)
	em.eventHistory = make([]EventExecutionRecord, 0)
	log.Printf("Cleared all event completion tracking for new game")
}
