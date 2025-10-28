// Package events provides default event handlers for different event types.
package events

import (
	"fmt"
	"log"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// CreateDefaultHandlers creates and returns default event handlers for all event types
func CreateDefaultHandlers() map[components.EventType]EventHandler {
	handlers := make(map[components.EventType]EventHandler)

	handlers[components.EventBattle] = HandleBattleEvent
	handlers[components.EventDialog] = HandleDialogEvent
	handlers[components.EventChest] = HandleChestEvent
	handlers[components.EventDoor] = HandleDoorEvent
	handlers[components.EventTrap] = HandleTrapEvent
	handlers[components.EventInfo] = HandleInfoEvent
	handlers[components.EventQuest] = HandleQuestEvent
	handlers[components.EventCutscene] = HandleCutsceneEvent
	handlers[components.EventShop] = HandleShopEvent
	handlers[components.EventRest] = HandleRestEvent

	return handlers
}

// HandleBattleEvent handles battle-related events
func HandleBattleEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Battle event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})

	// Extract battle data
	if len(eventComp.EventData.Enemies) > 0 {
		data["enemy_ids"] = eventComp.EventData.Enemies
	}

	if eventComp.EventData.BattleMap != "" {
		data["battle_map"] = eventComp.EventData.BattleMap
	}

	data["ambush"] = false // TODO: Add ambush support if needed
	data["event_entity_id"] = entity.GetID()

	return &EventResult{
		Success:    true,
		Message:    fmt.Sprintf("Battle started: %s", eventComp.Name),
		Data:       data,
		NextAction: "start_battle",
	}
}

// HandleDialogEvent handles dialog-related events
func HandleDialogEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Dialog event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})

	if eventComp.EventData.NPCID != "" {
		data["npc_id"] = eventComp.EventData.NPCID
	}

	if eventComp.EventData.DialogID != "" {
		data["dialog_id"] = eventComp.EventData.DialogID
	}

	data["event_entity_id"] = entity.GetID()

	return &EventResult{
		Success:    true,
		Message:    fmt.Sprintf("Dialog started: %s", eventComp.Name),
		Data:       data,
		NextAction: "start_dialog",
	}
}

// HandleChestEvent handles chest/container-related events
func HandleChestEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Chest event triggered: %s", eventComp.Name)

	// Check if chest is locked and player has key
	if eventComp.EventData.IsLocked {
		// TODO: Check player inventory for required key
		// For now, assume player doesn't have key
		return &EventResult{
			Success: false,
			Message: fmt.Sprintf("The %s is locked. You need a key to open it.", eventComp.Name),
		}
	}

	data := make(map[string]interface{})

	if len(eventComp.EventData.Items) > 0 {
		data["item_ids"] = eventComp.EventData.Items
	}

	if eventComp.EventData.Gold > 0 {
		data["gold"] = eventComp.EventData.Gold
	}

	data["event_entity_id"] = entity.GetID()

	// TODO: Add items to player inventory
	// TODO: Add gold to player

	return &EventResult{
		Success:    true,
		Message:    fmt.Sprintf("Opened %s", eventComp.Name),
		Data:       data,
		NextAction: "show_loot",
	}
}

// HandleDoorEvent handles door/travel-related events
func HandleDoorEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Door event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})

	if eventComp.EventData.TargetMap != "" {
		data["target_map"] = eventComp.EventData.TargetMap
		data["target_x"] = eventComp.EventData.TargetX
		data["target_y"] = eventComp.EventData.TargetY
	}

	data["event_entity_id"] = entity.GetID()

	return &EventResult{
		Success:    true,
		Message:    fmt.Sprintf("Moving through %s", eventComp.Name),
		Data:       data,
		NextAction: "change_map",
	}
}

// HandleTrapEvent handles trap-related events
func HandleTrapEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Trap event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})
	data["trap_type"] = eventComp.EventData.TrapType
	data["damage"] = eventComp.EventData.Damage
	data["status_effect"] = eventComp.EventData.TrapType
	data["event_entity_id"] = entity.GetID()

	// TODO: Apply damage to player
	// TODO: Apply status effect to player

	message := fmt.Sprintf("You triggered a %s!", eventComp.EventData.TrapType)
	if eventComp.EventData.Damage > 0 {
		message += fmt.Sprintf(" You take %d damage!", eventComp.EventData.Damage)
	}

	return &EventResult{
		Success:    true,
		Message:    message,
		Data:       data,
		NextAction: "apply_trap_effects",
	}
}

// HandleInfoEvent handles information display events
func HandleInfoEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Info event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})
	data["title"] = eventComp.EventData.Title
	data["message"] = eventComp.EventData.Message

	// TODO: Add image path support if needed
	// if eventComp.EventData.ImagePath != "" {
	//     data["image_path"] = eventComp.EventData.ImagePath
	// }

	data["event_entity_id"] = entity.GetID()

	return &EventResult{
		Success:    true,
		Message:    eventComp.EventData.Title,
		Data:       data,
		NextAction: "show_info",
	}
}

// HandleQuestEvent handles quest-related events
func HandleQuestEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Quest event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})

	if eventComp.EventData.QuestID != "" {
		data["quest_id"] = eventComp.EventData.QuestID
	}

	// TODO: Add objective ID support if needed
	// if eventComp.EventData.ObjectiveID != "" {
	//     data["objective_id"] = eventComp.EventData.ObjectiveID
	// }

	data["event_entity_id"] = entity.GetID()

	// TODO: Integrate with quest system
	// TODO: Update quest progress or start new quest

	return &EventResult{
		Success:    true,
		Message:    fmt.Sprintf("Quest event: %s", eventComp.Name),
		Data:       data,
		NextAction: "update_quest",
	}
}

// HandleCutsceneEvent handles cutscene/story events
func HandleCutsceneEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Cutscene event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})
	data["cutscene_id"] = eventComp.ID
	data["event_entity_id"] = entity.GetID()

	// TODO: Add custom data support for cutscenes if needed
	// if eventComp.EventData.CustomData != nil {
	//     for key, value := range eventComp.EventData.CustomData {
	//         data[key] = value
	//     }
	// }

	return &EventResult{
		Success:    true,
		Message:    fmt.Sprintf("Cutscene started: %s", eventComp.Name),
		Data:       data,
		NextAction: "play_cutscene",
	}
}

// HandleShopEvent handles shop-related events
func HandleShopEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Shop event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})
	data["shop_id"] = eventComp.ID
	data["event_entity_id"] = entity.GetID()

	// TODO: Add custom shop data support if needed
	// if eventComp.EventData.CustomData != nil {
	//     for key, value := range eventComp.EventData.CustomData {
	//         data[key] = value
	//     }
	// }

	return &EventResult{
		Success:    true,
		Message:    fmt.Sprintf("Welcome to %s!", eventComp.Name),
		Data:       data,
		NextAction: "open_shop",
	}
}

// HandleRestEvent handles rest/save point events
func HandleRestEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *EventResult {
	log.Printf("Rest event triggered: %s", eventComp.Name)

	data := make(map[string]interface{})
	data["rest_point_id"] = eventComp.ID
	data["event_entity_id"] = entity.GetID()

	// TODO: Implement rest functionality
	// TODO: Heal player, save game, etc.

	return &EventResult{
		Success:    true,
		Message:    fmt.Sprintf("You rest at %s. Your health is restored.", eventComp.Name),
		Data:       data,
		NextAction: "rest_player",
	}
}
