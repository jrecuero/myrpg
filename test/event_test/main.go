package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/events"
)

// Game represents our event system test game
type Game struct {
	eventManager  *events.EventManager
	player        *ecs.Entity
	eventEntities []*ecs.Entity
}

func main() {
	fmt.Println("Starting Event System Test...")

	game := &Game{}

	if err := game.init(); err != nil {
		log.Fatal(err)
	}

	// Test event system without Ebitengine for now
	game.testEventSystem()

	fmt.Println("Event System Test completed successfully!")
}

func (g *Game) init() error {
	// Create event manager
	g.eventManager = events.NewEventManager()

	// Register default handlers
	handlers := events.CreateDefaultHandlers()
	for eventType, handler := range handlers {
		g.eventManager.RegisterHandler(eventType, handler)
	}

	// Create player entity
	g.player = ecs.NewEntity("Player")
	g.player.AddComponent(ecs.ComponentTransform, components.NewTransform(100, 100, 32, 32))
	g.eventManager.SetPlayer(g.player)

	// Create test event entities
	g.createTestEvents()

	return nil
}

func (g *Game) createTestEvents() {
	g.eventEntities = make([]*ecs.Entity, 0)

	// 1. Battle Event (Touch trigger)
	battleEntity := ecs.NewEntity("Enemy Encounter")
	battleEntity.AddComponent(ecs.ComponentTransform, components.NewTransform(150, 100, 32, 32))

	battleEvent := components.NewEventComponent("battle_001", "Goblin Encounter", components.TriggerOnTouch, components.EventBattle)
	battleEvent.SetEventData(components.EventData{
		Enemies:   []string{"goblin_1", "goblin_2", "goblin_3"},
		BattleMap: "forest_clearing",
	})
	battleEntity.AddComponent(ecs.ComponentEvent, battleEvent)

	g.eventManager.RegisterEntity(battleEntity)
	g.eventEntities = append(g.eventEntities, battleEntity)

	// 2. Chest Event (Touch trigger)
	chestEntity := ecs.NewEntity("Treasure Chest")
	chestEntity.AddComponent(ecs.ComponentTransform, components.NewTransform(200, 150, 32, 32))

	chestEvent := components.NewEventComponent("chest_001", "Wooden Chest", components.TriggerOnTouch, components.EventChest)
	chestEvent.SetEventData(components.EventData{
		Items:    []string{"iron_sword", "health_potion"},
		Gold:     50,
		IsLocked: false,
	})
	chestEntity.AddComponent(ecs.ComponentEvent, chestEvent)

	g.eventManager.RegisterEntity(chestEntity)
	g.eventEntities = append(g.eventEntities, chestEntity)

	// 3. Dialog Event (Proximity trigger)
	npcEntity := ecs.NewEntity("Village Elder")
	npcEntity.AddComponent(ecs.ComponentTransform, components.NewTransform(300, 100, 32, 32))

	dialogEvent := components.NewEventComponent("dialog_001", "Talk to Elder", components.TriggerOnProximity, components.EventDialog)
	dialogEvent.SetConditionData(components.EventConditionData{
		Distance: 48, // 1.5 tiles
	})
	dialogEvent.SetEventData(components.EventData{
		NPCID:    "elder_001",
		DialogID: "elder_intro",
	})
	dialogEvent.SetRepeatable(true)
	npcEntity.AddComponent(ecs.ComponentEvent, dialogEvent)

	g.eventManager.RegisterEntity(npcEntity)
	g.eventEntities = append(g.eventEntities, npcEntity)

	// 4. Door Event (Touch trigger)
	doorEntity := ecs.NewEntity("Cave Entrance")
	doorEntity.AddComponent(ecs.ComponentTransform, components.NewTransform(400, 200, 32, 32))

	doorEvent := components.NewEventComponent("door_001", "Enter Cave", components.TriggerOnTouch, components.EventDoor)
	doorEvent.SetEventData(components.EventData{
		TargetMap: "cave_level_1",
		TargetX:   50,
		TargetY:   50,
	})
	doorEntity.AddComponent(ecs.ComponentEvent, doorEvent)

	g.eventManager.RegisterEntity(doorEntity)
	g.eventEntities = append(g.eventEntities, doorEntity)

	// 5. Info Event (Touch trigger)
	signEntity := ecs.NewEntity("Village Sign")
	signEntity.AddComponent(ecs.ComponentTransform, components.NewTransform(250, 250, 32, 32))

	infoEvent := components.NewEventComponent("info_001", "Village Sign", components.TriggerOnTouch, components.EventInfo)
	infoEvent.SetEventData(components.EventData{
		Title:   "Welcome to Riverside Village",
		Message: "Population: 127\nFounded: Year 892\nMayor: Eldric Stormwind",
	})
	infoEvent.SetRepeatable(true)
	signEntity.AddComponent(ecs.ComponentEvent, infoEvent)

	g.eventManager.RegisterEntity(signEntity)
	g.eventEntities = append(g.eventEntities, signEntity)

	fmt.Printf("Created %d test events:\n", len(g.eventEntities))
	for _, entity := range g.eventEntities {
		if event := entity.Event(); event != nil {
			fmt.Printf("- %s (%s): %s trigger\n", event.Name, event.GetEventTypeName(), event.GetTriggerConditionName())
		}
	}
}

func (g *Game) testEventSystem() {
	fmt.Println("\n=== Event System Test ===")

	// Test 1: Player doesn't trigger anything initially
	fmt.Println("\nTest 1: Initial state (no events should trigger)")
	g.eventManager.Update(0.016) // 60 FPS delta

	// Test 2: Move player to battle event (touch trigger)
	fmt.Println("\nTest 2: Moving player to battle event position")
	playerTransform := g.player.Transform()
	playerTransform.X = 150 // Same X as battle entity
	playerTransform.Y = 100 // Same Y as battle entity
	g.eventManager.Update(0.016)

	// Test 3: Move player near NPC (proximity trigger)
	fmt.Println("\nTest 3: Moving player near NPC (proximity trigger)")
	playerTransform.X = 280 // Near NPC
	playerTransform.Y = 100
	g.eventManager.Update(0.016)

	// Test 4: Move player to chest
	fmt.Println("\nTest 4: Moving player to chest")
	playerTransform.X = 200
	playerTransform.Y = 150
	g.eventManager.Update(0.016)

	// Test 5: Move player to door
	fmt.Println("\nTest 5: Moving player to door")
	playerTransform.X = 400
	playerTransform.Y = 200
	g.eventManager.Update(0.016)

	// Test 6: Move player to info sign
	fmt.Println("\nTest 6: Moving player to info sign")
	playerTransform.X = 250
	playerTransform.Y = 250
	g.eventManager.Update(0.016)

	// Test 7: Test repeatable event (dialog again)
	fmt.Println("\nTest 7: Testing repeatable dialog event")
	playerTransform.X = 280
	playerTransform.Y = 100
	g.eventManager.Update(0.016)

	// Display event statistics
	fmt.Println("\n=== Event Statistics ===")
	stats := g.eventManager.GetEventStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	// Display event history
	fmt.Println("\n=== Event Execution History ===")
	history := g.eventManager.GetEventHistory()
	for i, record := range history {
		fmt.Printf("%d. Event: %s, Entity: %s, Time: %s, Success: %t\n",
			i+1, record.EventID, record.EntityID, record.Timestamp.Format("15:04:05"), record.Result.Success)
		fmt.Printf("   Message: %s\n", record.Result.Message)
		if record.Result.NextAction != "" {
			fmt.Printf("   Next Action: %s\n", record.Result.NextAction)
		}
	}
}

// Ebitengine interface methods (for future visual testing)
func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Event System Test - See console for results")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
