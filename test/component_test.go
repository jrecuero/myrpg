package test
package main

import (
	"fmt"
	"image/color"

	"github.com/jrecuero/myrpg/cmd/myrpg/game/entities"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

func main() {
	fmt.Println("Testing Event System Components...")
	fmt.Println("=====================================")

	// Test Battle Event
	fmt.Println("\n1. Testing Battle Event:")
	battleEvent := entities.CreateBattleEvent("battle1", "Test Battle", 100, 100, []string{"orc"}, "forest")
	
	if eventComp := battleEvent.GetComponent(&components.EventComponent{}); eventComp != nil {
		event := eventComp.(*components.EventComponent)
		fmt.Printf("   - ID: %s\n", event.ID)
		fmt.Printf("   - Name: %s\n", event.Name)
		fmt.Printf("   - Type: %s\n", event.EventType)
		fmt.Printf("   - Visible: %t\n", event.Visible)
		fmt.Printf("   - Fallback Color: %+v (should be red)\n", event.FallbackColor)
		fmt.Printf("   - Should trigger on collision: %t\n", event.ShouldTriggerOnCollision())
	}

	// Test Chest Event
	fmt.Println("\n2. Testing Chest Event:")
	chestEvent := entities.CreateChestEvent("chest1", "Treasure Chest", 200, 200, []string{"sword"}, 50, false)
	
	if eventComp := chestEvent.GetComponent(&components.EventComponent{}); eventComp != nil {
		event := eventComp.(*components.EventComponent)
		fmt.Printf("   - ID: %s\n", event.ID)
		fmt.Printf("   - Name: %s\n", event.Name)
		fmt.Printf("   - Type: %s\n", event.EventType)
		fmt.Printf("   - Visible: %t\n", event.Visible)
		fmt.Printf("   - Fallback Color: %+v (should be yellow)\n", event.FallbackColor)
		fmt.Printf("   - Should trigger on collision: %t\n", event.ShouldTriggerOnCollision())
	}

	// Test Dialog Event
	fmt.Println("\n3. Testing Dialog Event:")
	dialogEvent := entities.CreateDialogEvent("dialog1", "Village Elder", 300, 300, "elder", "greeting")
	
	if eventComp := dialogEvent.GetComponent(&components.EventComponent{}); eventComp != nil {
		event := eventComp.(*components.EventComponent)
		fmt.Printf("   - ID: %s\n", event.ID)
		fmt.Printf("   - Name: %s\n", event.Name)
		fmt.Printf("   - Type: %s\n", event.EventType)
		fmt.Printf("   - Visible: %t\n", event.Visible)
		fmt.Printf("   - Fallback Color: %+v (should be blue)\n", event.FallbackColor)
		fmt.Printf("   - Should trigger on collision: %t\n", event.ShouldTriggerOnCollision())
	}

	// Test Info Event
	fmt.Println("\n4. Testing Info Event:")
	infoEvent := entities.CreateInfoEvent("info1", "Ancient Stone", 400, 400, "Stone Marker", "An old marker stone")
	
	if eventComp := infoEvent.GetComponent(&components.EventComponent{}); eventComp != nil {
		event := eventComp.(*components.EventComponent)
		fmt.Printf("   - ID: %s\n", event.ID)
		fmt.Printf("   - Name: %s\n", event.Name)
		fmt.Printf("   - Type: %s\n", event.EventType)
		fmt.Printf("   - Visible: %t\n", event.Visible)
		fmt.Printf("   - Fallback Color: %+v (should be green)\n", event.FallbackColor)
		fmt.Printf("   - Should trigger on collision: %t\n", event.ShouldTriggerOnCollision())
	}

	// Test Trap Event (Hidden)
	fmt.Println("\n5. Testing Trap Event (should be hidden):")
	trapEvent := entities.CreateTrapEvent("trap1", "Spike Trap", 500, 500, "Spikes", 15)
	
	if eventComp := trapEvent.GetComponent(&components.EventComponent{}); eventComp != nil {
		event := eventComp.(*components.EventComponent)
		fmt.Printf("   - ID: %s\n", event.ID)
		fmt.Printf("   - Name: %s\n", event.Name)
		fmt.Printf("   - Type: %s\n", event.EventType)
		fmt.Printf("   - Visible: %t (should be false - hidden)\n", event.Visible)
		fmt.Printf("   - Fallback Color: %+v (should be red)\n", event.FallbackColor)
		fmt.Printf("   - Should trigger on collision: %t\n", event.ShouldTriggerOnCollision())
	}

	// Test Timer Event (Hidden)
	fmt.Println("\n6. Testing Timer Event (should be hidden):")
	timerEvent := entities.CreateTimerEvent("timer1", "Auto Event", 600, 600, 5000) // 5 second timer
	
	if eventComp := timerEvent.GetComponent(&components.EventComponent{}); eventComp != nil {
		event := eventComp.(*components.EventComponent)
		fmt.Printf("   - ID: %s\n", event.ID)
		fmt.Printf("   - Name: %s\n", event.Name)
		fmt.Printf("   - Type: %s\n", event.EventType)
		fmt.Printf("   - Visible: %t (should be false - hidden)\n", event.Visible)
		fmt.Printf("   - Fallback Color: %+v (should be purple)\n", event.FallbackColor)
		fmt.Printf("   - Should trigger on collision: %t (should be false - timer based)\n", event.ShouldTriggerOnCollision())
	}

	fmt.Println("\n=====================================")
	fmt.Println("Event System Component Test Complete!")
	
	// Test color values
	fmt.Println("\n7. Verifying Color Values:")
	expectedColors := map[string]color.RGBA{
		"battle": {255, 0, 0, 255},   // Red
		"chest":  {255, 255, 0, 255}, // Yellow
		"dialog": {0, 0, 255, 255},   // Blue
		"info":   {0, 255, 0, 255},   // Green
		"trap":   {255, 0, 0, 255},   // Red
		"timer":  {128, 0, 128, 255}, // Purple
	}
	
	for eventType, expectedColor := range expectedColors {
		fmt.Printf("   - %s events should be: R=%d, G=%d, B=%d, A=%d\n", 
			eventType, expectedColor.R, expectedColor.G, expectedColor.B, expectedColor.A)
	}
}