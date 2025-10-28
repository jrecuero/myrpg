// Test program for Event State Persistence system
// Demonstrates saving and loading event states with prerequisites and cooldowns
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/engine"
	"github.com/jrecuero/myrpg/internal/events"
	"github.com/jrecuero/myrpg/internal/logger"
)

func main() {
	fmt.Println("=== Event State Persistence Test ===")

	// Initialize logger
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Create a new game instance
	game := engine.NewGame()

	// Create test events with different states and properties
	fmt.Println("\n1. Creating test events...")

	// Event 1: Starter quest (no prerequisites)
	starterQuest := createTestEventEntity(
		"starter_quest", "Welcome Quest", 100, 100,
		components.EventQuest, components.TriggerOnInteract,
		[]string{}, // No prerequisites
		3,          // Max 3 triggers
		0,          // No cooldown
	)
	game.AddEntity(starterQuest)

	// Event 2: Advanced quest (requires starter quest)
	advancedQuest := createTestEventEntity(
		"advanced_quest", "Advanced Quest", 150, 100,
		components.EventQuest, components.TriggerOnInteract,
		[]string{"starter_quest"}, // Requires starter quest
		1,                         // Max 1 trigger
		0,                         // No cooldown
	)
	game.AddEntity(advancedQuest)

	// Event 3: Repeatable resource node (with cooldown)
	resourceNode := createTestEventEntity(
		"resource_node", "Magic Crystal", 200, 100,
		components.EventChest, components.TriggerOnTouch,
		[]string{},     // No prerequisites
		0,              // Unlimited triggers
		30*time.Second, // 30 second cooldown
	)
	game.AddEntity(resourceNode)

	// Event 4: Hidden trap event
	hiddenTrap := createTestEventEntity(
		"hidden_trap", "Pressure Plate Trap", 250, 150,
		components.EventTrap, components.TriggerOnTouch,
		[]string{},    // No prerequisites
		5,             // Max 5 triggers
		5*time.Second, // 5 second cooldown
	)
	// Make it hidden and set to tactical mode only
	if trapEvent := hiddenTrap.Event(); trapEvent != nil {
		trapEvent.SetVisible(false).SetActiveInMode(components.GameModeTactical)
	}
	game.AddEntity(hiddenTrap)

	fmt.Printf("Created %d test events\n", 4)

	// Simulate some event interactions
	fmt.Println("\n2. Simulating event interactions...")

	// Trigger starter quest twice
	if starterEvent := starterQuest.Event(); starterEvent != nil {
		fmt.Printf("Starter Quest - Initial state: %v, triggers: %d\n", starterEvent.State, starterEvent.TriggerCount)

		starterEvent.Trigger()
		fmt.Printf("Starter Quest - After trigger 1: %v, triggers: %d\n", starterEvent.State, starterEvent.TriggerCount)

		starterEvent.Trigger()
		fmt.Printf("Starter Quest - After trigger 2: %v, triggers: %d\n", starterEvent.State, starterEvent.TriggerCount)
	}

	// Complete the starter quest to unlock advanced quest
	if starterEvent := starterQuest.Event(); starterEvent != nil {
		starterEvent.Trigger() // This should complete it (3rd trigger)
		fmt.Printf("Starter Quest - After trigger 3: %v, triggers: %d (should be completed)\n",
			starterEvent.State, starterEvent.TriggerCount)
	}

	// Try to trigger resource node
	if resourceEvent := resourceNode.Event(); resourceEvent != nil {
		fmt.Printf("Resource Node - Initial state: %v, triggers: %d\n", resourceEvent.State, resourceEvent.TriggerCount)

		resourceEvent.Trigger()
		fmt.Printf("Resource Node - After trigger 1: %v, triggers: %d, last triggered: %v\n",
			resourceEvent.State, resourceEvent.TriggerCount, resourceEvent.LastTriggered)
	}

	// Try to trigger hidden trap
	if trapEvent := hiddenTrap.Event(); trapEvent != nil {
		fmt.Printf("Hidden Trap - Initial state: %v, triggers: %d, visible: %v\n",
			trapEvent.State, trapEvent.TriggerCount, trapEvent.Visible)

		trapEvent.Trigger()
		fmt.Printf("Hidden Trap - After trigger 1: %v, triggers: %d\n",
			trapEvent.State, trapEvent.TriggerCount)
	}

	// Update event manager completion tracking
	eventManager := events.NewEventManager()
	eventManager.SetEventCompleted("starter_quest", true)

	// Save the current state
	fmt.Println("\n3. Saving event state...")
	if err := game.SaveEventState(); err != nil {
		log.Fatalf("Failed to save event state: %v", err)
	}
	fmt.Println("Event state saved successfully!")

	// Show save statistics
	saveManager := game.GetSaveManager()
	if saveData := saveManager.GetCurrentEventSaveData(); saveData != nil {
		stats := saveData.GetStatistics()
		fmt.Printf("Save Statistics: %+v\n", stats)
	}

	// Reset all events to simulate new game
	fmt.Println("\n4. Clearing event state (simulating new game)...")
	if err := game.ClearEventState(); err != nil {
		log.Fatalf("Failed to clear event state: %v", err)
	}

	// Show events are reset
	if starterEvent := starterQuest.Event(); starterEvent != nil {
		fmt.Printf("Starter Quest after reset: %v, triggers: %d\n", starterEvent.State, starterEvent.TriggerCount)
	}
	if resourceEvent := resourceNode.Event(); resourceEvent != nil {
		fmt.Printf("Resource Node after reset: %v, triggers: %d\n", resourceEvent.State, resourceEvent.TriggerCount)
	}

	// Load the saved state back
	fmt.Println("\n5. Loading event state from save...")
	if err := game.LoadEventState(); err != nil {
		log.Fatalf("Failed to load event state: %v", err)
	}
	fmt.Println("Event state loaded successfully!")

	// Verify loaded state
	fmt.Println("\n6. Verifying loaded state...")
	if starterEvent := starterQuest.Event(); starterEvent != nil {
		fmt.Printf("Starter Quest after load: %v, triggers: %d\n", starterEvent.State, starterEvent.TriggerCount)
	}
	if resourceEvent := resourceNode.Event(); resourceEvent != nil {
		fmt.Printf("Resource Node after load: %v, triggers: %d, last triggered: %v\n",
			resourceEvent.State, resourceEvent.TriggerCount, resourceEvent.LastTriggered)
	}
	if trapEvent := hiddenTrap.Event(); trapEvent != nil {
		fmt.Printf("Hidden Trap after load: %v, triggers: %d, visible: %v, mode: %v\n",
			trapEvent.State, trapEvent.TriggerCount, trapEvent.Visible, trapEvent.ActiveInMode)
	}

	// Test prerequisite checking with advanced quest
	fmt.Println("\n7. Testing prerequisite system...")
	if advancedEvent := advancedQuest.Event(); advancedEvent != nil {
		fmt.Printf("Advanced Quest prerequisites: %v\n", advancedEvent.Prerequisites)
		fmt.Printf("Advanced Quest can trigger: %v (should be true since starter is completed)\n",
			advancedEvent.CanTrigger())

		if advancedEvent.CanTrigger() {
			advancedEvent.Trigger()
			fmt.Printf("Advanced Quest after trigger: %v, triggers: %d\n",
				advancedEvent.State, advancedEvent.TriggerCount)
		}
	}

	// Test cooldown system
	fmt.Println("\n8. Testing cooldown system...")
	if resourceEvent := resourceNode.Event(); resourceEvent != nil {
		fmt.Printf("Resource Node cooldown: %v\n", resourceEvent.Cooldown)
		fmt.Printf("Resource Node can trigger immediately: %v (should be false due to cooldown)\n",
			resourceEvent.CanTrigger())

		fmt.Println("Waiting 2 seconds and trying again...")
		time.Sleep(2 * time.Second)
		fmt.Printf("Resource Node can trigger after 2s: %v (still should be false)\n",
			resourceEvent.CanTrigger())
	}

	fmt.Println("\n9. Save file information...")
	if files, err := saveManager.ListSaveFiles(); err == nil {
		fmt.Printf("Available save files: %v\n", files)

		for _, file := range files {
			if info, err := saveManager.GetSaveFileInfo(file); err == nil {
				fmt.Printf("File %s info: %+v\n", file, info)
			}
		}
	}

	fmt.Println("\n=== Event State Persistence Test Complete ===")
}

// createTestEventEntity creates a test entity with an event component
func createTestEventEntity(id, name string, x, y float64, eventType components.EventType,
	triggerCondition components.TriggerCondition, prerequisites []string, maxTriggers int,
	cooldown time.Duration) *ecs.Entity {

	entity := ecs.NewEntity(name)

	// Add transform component
	entity.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))

	// Create event component with all the advanced features
	eventComp := components.NewEventComponent(id, name, triggerCondition, eventType)
	eventComp.SetPrerequisites(prerequisites...)
	eventComp.SetMaxTriggers(maxTriggers)
	eventComp.SetCooldown(cooldown)
	eventComp.SetRepeatable(maxTriggers == 0 || maxTriggers > 1) // Repeatable if unlimited or more than 1

	// Add event component
	entity.AddComponent(ecs.ComponentEvent, eventComp)

	return entity
}
