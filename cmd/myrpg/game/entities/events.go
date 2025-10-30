// Package entities provides functions to create event entities for the game
package entities

import (
	"log"
	"time"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/gfx"
)

// CreateEventEntity creates an event entity with the specified parameters
func CreateEventEntity(id, name string, x, y float64, eventComp *components.EventComponent) *ecs.Entity {
	entity := ecs.NewEntity(name)

	// Add basic components
	entity.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
	entity.AddComponent(ecs.ComponentEvent, eventComp)

	// Add sprite or fallback rendering
	if eventComp.HasCustomSprite() {
		// Try to load custom sprite
		sprite, err := gfx.NewSpriteFromFile(eventComp.SpritePath, 32, 32)
		if err != nil {
			log.Printf("Failed to load event sprite %s: %v. Using fallback color.", eventComp.SpritePath, err)
			// Don't add sprite component, will use fallback color rendering
		} else {
			spriteComp := components.NewSpriteComponent(sprite, eventComp.SpriteScale,
				eventComp.SpriteOffset[0], eventComp.SpriteOffset[1])
			entity.AddComponent(ecs.ComponentSprite, spriteComp)
		}
	}
	// If no custom sprite or loading failed, we'll use fallback color rendering

	// Add collider only for events that should trigger on collision
	if eventComp.ShouldTriggerOnCollision() {
		entity.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 32, 32, 0, 0))
	}

	// Add event tag
	entity.AddTag("event")

	return entity
}

// CreateBattleEvent creates a battle event entity
func CreateBattleEvent(id, name string, x, y float64, enemyIDs []string, battleMap string) *ecs.Entity {
	eventComp := components.NewEventComponent(id, name, components.TriggerOnTouch, components.EventBattle)
	eventComp.SetEventData(components.EventData{
		Enemies:   enemyIDs,
		BattleMap: battleMap,
	})
	eventComp.SetActiveInMode(components.GameModeExploration) // Only active in exploration mode
	return CreateEventEntity(id, name, x, y, eventComp)
}

// CreateDialogEvent creates a dialog event entity
func CreateDialogEvent(id, name string, x, y float64, npcID, dialogID string) *ecs.Entity {
	eventComp := components.NewEventComponent(id, name, components.TriggerOnProximity, components.EventDialog)
	eventComp.SetConditionData(components.EventConditionData{
		Distance: 48, // 1.5 tiles
	})
	eventComp.SetEventData(components.EventData{
		NPCID:    npcID,
		DialogID: dialogID,
	})
	eventComp.SetRepeatable(true)                             // Dialogs can be repeated
	eventComp.SetActiveInMode(components.GameModeExploration) // Only active in exploration mode
	return CreateEventEntity(id, name, x, y, eventComp)
}

// CreateChestEvent creates a chest event entity
func CreateChestEvent(id, name string, x, y float64, itemIDs []string, gold int, locked bool) *ecs.Entity {
	eventComp := components.NewEventComponent(id, name, components.TriggerOnTouch, components.EventChest)
	eventComp.SetEventData(components.EventData{
		Items:    itemIDs,
		Gold:     gold,
		IsLocked: locked,
	})
	eventComp.SetActiveInMode(components.GameModeExploration) // Only active in exploration mode
	return CreateEventEntity(id, name, x, y, eventComp)
}

// CreateDoorEvent creates a door/travel event entity
func CreateDoorEvent(id, name string, x, y float64, targetMap string, targetX, targetY int) *ecs.Entity {
	eventComp := components.NewEventComponent(id, name, components.TriggerOnTouch, components.EventDoor)
	eventComp.SetEventData(components.EventData{
		TargetMap: targetMap,
		TargetX:   float64(targetX),
		TargetY:   float64(targetY),
	})
	eventComp.SetActiveInMode(components.GameModeExploration) // Only active in exploration mode
	return CreateEventEntity(id, name, x, y, eventComp)
}

// CreateInfoEvent creates an information event entity (signs, etc.)
func CreateInfoEvent(id, name string, x, y float64, title, message string) *ecs.Entity {
	eventComp := components.NewEventComponent(id, name, components.TriggerOnTouch, components.EventInfo)
	eventComp.SetEventData(components.EventData{
		Title:   title,
		Message: message,
	})
	eventComp.SetRepeatable(true)                             // Info can be read multiple times
	eventComp.SetActiveInMode(components.GameModeExploration) // Only active in exploration mode
	return CreateEventEntity(id, name, x, y, eventComp)
}

// CreateTrapEvent creates a hidden trap event entity
func CreateTrapEvent(id, name string, x, y float64, trapType string, damage int) *ecs.Entity {
	eventComp := components.NewEventComponent(id, name, components.TriggerOnTouch, components.EventTrap)
	eventComp.SetEventData(components.EventData{
		TrapType: trapType,
		Damage:   damage,
	})
	eventComp.SetVisible(false)                               // Traps are hidden by default
	eventComp.SetActiveInMode(components.GameModeExploration) // Only active in exploration mode
	return CreateEventEntity(id, name, x, y, eventComp)
}

// CreateTimerEvent creates a timer-based event (not visible, triggers after delay)
func CreateTimerEvent(id, name string, x, y float64, eventType components.EventType, delay int) *ecs.Entity {
	eventComp := components.NewEventComponent(id, name, components.TriggerOnTimeout, eventType)
	eventComp.SetConditionData(components.EventConditionData{
		Duration: time.Duration(delay) * time.Second,
	})
	eventComp.SetVisible(false)                               // Timer events are not visible
	eventComp.SetActiveInMode(components.GameModeExploration) // Only active in exploration mode
	return CreateEventEntity(id, name, x, y, eventComp)
}

// CreateQuestEvent creates a quest-related event entity
func CreateQuestEvent(id, name string, x, y float64, questID, objectiveID string) *ecs.Entity {
	eventComp := components.NewEventComponent(id, name, components.TriggerOnTouch, components.EventQuest)
	eventComp.SetEventData(components.EventData{
		QuestID: questID,
		// Note: ObjectiveID not in current EventData structure, can be added if needed
	})
	eventComp.SetActiveInMode(components.GameModeExploration) // Only active in exploration mode
	return CreateEventEntity(id, name, x, y, eventComp)
}
