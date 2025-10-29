// Package engine provides game-specific event handlers that integrate with the game engine systems
package engine

import (
	"fmt"
	"log"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/events"
)

// CreateGameEventHandlers creates event handlers that integrate with the game engine
func (g *Game) CreateGameEventHandlers() map[components.EventType]events.EventHandler {
	handlers := make(map[components.EventType]events.EventHandler)

	handlers[components.EventBattle] = g.handleBattleEvent
	handlers[components.EventDialog] = g.handleDialogEvent
	handlers[components.EventChest] = g.handleChestEvent
	handlers[components.EventDoor] = g.handleDoorEvent
	handlers[components.EventTrap] = g.handleTrapEvent
	handlers[components.EventInfo] = g.handleInfoEvent
	handlers[components.EventQuest] = g.handleQuestEvent
	handlers[components.EventCutscene] = g.handleCutsceneEvent
	handlers[components.EventShop] = g.handleShopEvent
	handlers[components.EventRest] = g.handleRestEvent

	return handlers
}

// handleBattleEvent integrates battle events with the tactical combat system
func (g *Game) handleBattleEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Battle event triggered: %s", eventComp.Name)

	// Get all combat participants (all party members + any enemies)
	participants := g.getAllCombatParticipants()

	// Add message about the encounter
	g.uiManager.AddMessage(fmt.Sprintf("A %s appears!", eventComp.Name))

	// Start tactical mode
	g.SwitchToTacticalMode(participants)

	return &events.EventResult{
		Success: true,
		Message: fmt.Sprintf("Battle started: %s", eventComp.Name),
		Data: map[string]interface{}{
			"enemy_ids":    eventComp.EventData.Enemies,
			"battle_map":   eventComp.EventData.BattleMap,
			"event_entity": entity.GetID(),
		},
	}
}

// handleDialogEvent switches to dialog view for NPC interactions
func (g *Game) handleDialogEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Dialog event triggered: %s", eventComp.Name)

	// Switch to dialog view with context data
	dialogData := map[string]interface{}{
		"npc_name":  eventComp.Name,
		"npc_id":    eventComp.EventData.NPCID,
		"dialog_id": eventComp.EventData.DialogID,
		"entity_id": entity.GetID(),
	}

	if err := g.SwitchToDialogView(dialogData); err != nil {
		log.Printf("Failed to switch to dialog view: %v", err)
		// Fallback to simple message
		message := fmt.Sprintf("You approach %s. They greet you warmly.", eventComp.Name)
		g.uiManager.AddMessage(message)
	}

	return &events.EventResult{
		Success: true,
		Message: fmt.Sprintf("Dialog with %s", eventComp.Name),
		Data:    dialogData,
	}
}

// handleChestEvent shows loot information using the info widget
func (g *Game) handleChestEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Chest event triggered: %s", eventComp.Name)

	// Check if locked
	if eventComp.EventData.IsLocked {
		message := fmt.Sprintf("The %s is locked!", eventComp.Name)
		g.uiManager.AddMessage(message)
		g.uiManager.ShowInfoWidget("Locked Chest", message, "")
		return &events.EventResult{
			Success: false,
			Message: message,
		}
	}

	// Build loot description
	lootDescription := "You found:\n\n"

	// Add items
	if len(eventComp.EventData.Items) > 0 {
		lootDescription += "Items:\n"
		for _, itemID := range eventComp.EventData.Items {
			lootDescription += fmt.Sprintf("• %s\n", itemID)
		}
		lootDescription += "\n"
	}

	// Add gold
	if eventComp.EventData.Gold > 0 {
		lootDescription += fmt.Sprintf("Gold: %d coins", eventComp.EventData.Gold)
	}

	// Show info widget with loot
	g.uiManager.ShowInfoWidget(fmt.Sprintf("Opened %s", eventComp.Name), lootDescription, "")
	g.uiManager.AddMessage(fmt.Sprintf("You opened a %s!", eventComp.Name))

	return &events.EventResult{
		Success: true,
		Message: fmt.Sprintf("Opened %s", eventComp.Name),
		Data: map[string]interface{}{
			"item_ids": eventComp.EventData.Items,
			"gold":     eventComp.EventData.Gold,
		},
	}
}

// handleDoorEvent handles area transitions
func (g *Game) handleDoorEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Door event triggered: %s", eventComp.Name)

	message := fmt.Sprintf("You approach %s. The path leads to %s.",
		eventComp.Name, eventComp.EventData.TargetMap)

	g.uiManager.AddMessage(message)
	g.uiManager.ShowInfoWidget("Travel Opportunity", message+"\n\n(Area transitions not yet implemented)", "")

	return &events.EventResult{
		Success: true,
		Message: fmt.Sprintf("Found passage: %s", eventComp.Name),
		Data: map[string]interface{}{
			"target_map": eventComp.EventData.TargetMap,
			"target_x":   eventComp.EventData.TargetX,
			"target_y":   eventComp.EventData.TargetY,
		},
	}
}

// handleTrapEvent handles trap activation
func (g *Game) handleTrapEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Trap event triggered: %s", eventComp.Name)

	damage := eventComp.EventData.Damage
	trapType := eventComp.EventData.TrapType

	message := fmt.Sprintf("You triggered a %s!", trapType)
	if damage > 0 {
		message += fmt.Sprintf(" You take %d damage!", damage)

		// Apply damage to player (if they have RPG stats)
		if playerStats := player.RPGStats(); playerStats != nil {
			playerStats.CurrentHP = max(0, playerStats.CurrentHP-damage)
			g.uiManager.AddMessage(fmt.Sprintf("%s takes %d damage! HP: %d/%d",
				playerStats.Name, damage, playerStats.CurrentHP, playerStats.MaxHP))
		}
	}

	g.uiManager.ShowInfoWidget("Trap!", message, "")
	g.uiManager.AddMessage(message)

	return &events.EventResult{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"trap_type": trapType,
			"damage":    damage,
		},
	}
}

// handleInfoEvent displays information using the info widget
func (g *Game) handleInfoEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Info event triggered: %s", eventComp.Name)

	title := eventComp.EventData.Title
	message := eventComp.EventData.Message

	if title == "" {
		title = eventComp.Name
	}

	g.uiManager.ShowInfoWidget(title, message, "")
	g.uiManager.AddMessage(fmt.Sprintf("You read: %s", title))

	return &events.EventResult{
		Success: true,
		Message: fmt.Sprintf("Read: %s", title),
		Data: map[string]interface{}{
			"title":   title,
			"message": message,
		},
	}
}

// handleQuestEvent handles quest-related events
func (g *Game) handleQuestEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Quest event triggered: %s", eventComp.Name)

	message := fmt.Sprintf("Quest event: %s", eventComp.Name)
	g.uiManager.AddMessage(message)
	g.uiManager.ShowInfoWidget("Quest", message+"\n\n(Quest integration coming soon)", "")

	return &events.EventResult{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"quest_id": eventComp.EventData.QuestID,
			// "objective_id": eventComp.EventData.ObjectiveID,  // TODO: Add back if needed
		},
	}
}

// handleCutsceneEvent handles cutscene events
func (g *Game) handleCutsceneEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Cutscene event triggered: %s", eventComp.Name)

	message := fmt.Sprintf("Cutscene: %s", eventComp.Name)
	g.uiManager.AddMessage(message)

	return &events.EventResult{
		Success: true,
		Message: message,
	}
}

// handleShopEvent handles shop events
func (g *Game) handleShopEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Shop event triggered: %s", eventComp.Name)

	message := fmt.Sprintf("Welcome to %s!", eventComp.Name)
	g.uiManager.AddMessage(message)
	g.uiManager.ShowInfoWidget("Shop", message+"\n\n(Shopping system coming soon)", "")

	return &events.EventResult{
		Success: true,
		Message: message,
	}
}

// handleRestEvent handles rest/save points
func (g *Game) handleRestEvent(entity *ecs.Entity, eventComp *components.EventComponent, player *ecs.Entity) *events.EventResult {
	log.Printf("Rest event triggered: %s", eventComp.Name)

	message := fmt.Sprintf("You rest at %s. Your health is restored!", eventComp.Name)

	// Heal the player
	if playerStats := player.RPGStats(); playerStats != nil {
		playerStats.CurrentHP = playerStats.MaxHP
		g.uiManager.AddMessage(fmt.Sprintf("%s's HP fully restored!", playerStats.Name))
	}

	g.uiManager.ShowInfoWidget("Rest Point", message, "")
	g.uiManager.AddMessage(message)

	return &events.EventResult{
		Success: true,
		Message: message,
	}
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
