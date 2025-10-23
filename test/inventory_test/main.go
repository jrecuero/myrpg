package main

// Test program for the Inventory Widget system
// This program creates a test scenario with items and demonstrates
// the inventory widget functionality including drag & drop, tooltips,
// and item management.

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/ui"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game represents the test game state
type Game struct {
	player          *ecs.Entity
	inventoryWidget *ui.InventoryWidget
	isInventoryOpen bool
	wasIPressed     bool
}

// NewGame creates a new test game instance
func NewGame() *Game {
	g := &Game{}

	// Create player entity with inventory
	g.player = ecs.NewEntity("TestPlayer")

	// Add inventory component (8x6 grid = 48 slots)
	inventoryComp := components.NewInventoryComponent(8, 6)
	g.player.AddComponent(ecs.ComponentInventory, inventoryComp)

	// Add some test items to inventory
	g.addTestItems(inventoryComp)

	// Create inventory widget
	widgetX := (screenWidth - 600) / 2  // Center horizontally
	widgetY := (screenHeight - 500) / 2 // Center vertically
	g.inventoryWidget = ui.NewInventoryWidget(widgetX, widgetY, g.player)

	// Start with inventory closed
	g.inventoryWidget.Visible = false
	g.isInventoryOpen = false

	return g
}

// addTestItems adds sample items to the inventory for testing
func (g *Game) addTestItems(inv *components.InventoryComponent) {
	// Create test items with different types and rarities

	// Equipment items
	sword := &components.Item{
		ID:          1,
		Name:        "Iron Sword",
		Description: "A sturdy iron sword. Increases attack damage.",
		Type:        components.ItemTypeEquipment,
		Rarity:      components.ItemRarityCommon,
		Value:       100,
		IconID:      1,
		Stackable:   false,
		MaxStack:    1,
	}

	helmet := &components.Item{
		ID:          2,
		Name:        "Steel Helmet",
		Description: "A protective steel helmet. Provides armor.",
		Type:        components.ItemTypeEquipment,
		Rarity:      components.ItemRarityUncommon,
		Value:       150,
		IconID:      2,
		Stackable:   false,
		MaxStack:    1,
	}

	// Consumable items
	healthPotion := &components.Item{
		ID:          3,
		Name:        "Health Potion",
		Description: "Restores 50 HP when consumed.",
		Type:        components.ItemTypeConsumable,
		Rarity:      components.ItemRarityCommon,
		Value:       25,
		IconID:      3,
		Stackable:   true,
		MaxStack:    10,
		Effects: []components.ConsumableEffect{
			{Type: "heal_hp", Value: 50, Target: "self"},
		},
	}

	manaPotion := &components.Item{
		ID:          4,
		Name:        "Mana Potion",
		Description: "Restores 30 MP when consumed.",
		Type:        components.ItemTypeConsumable,
		Rarity:      components.ItemRarityCommon,
		Value:       20,
		IconID:      4,
		Stackable:   true,
		MaxStack:    10,
		Effects: []components.ConsumableEffect{
			{Type: "heal_mp", Value: 30, Target: "self"},
		},
	}

	// Material items
	ironOre := &components.Item{
		ID:          5,
		Name:        "Iron Ore",
		Description: "Raw iron ore. Can be smelted into ingots.",
		Type:        components.ItemTypeMaterial,
		Rarity:      components.ItemRarityCommon,
		Value:       5,
		IconID:      5,
		Stackable:   true,
		MaxStack:    50,
	}

	// Quest item
	mysteriousKey := &components.Item{
		ID:          6,
		Name:        "Mysterious Key",
		Description: "An ancient ornate key discovered in the depths of the forgotten temple ruins. Its intricate engravings suggest it unlocks something of great importance, though its exact purpose remains unknown to scholars.",
		Type:        components.ItemTypeQuest,
		Rarity:      components.ItemRarityRare,
		Value:       0,
		IconID:      6,
		Stackable:   false,
		MaxStack:    1,
		QuestItem:   true,
		Unique:      true,
	}

	// Add items to inventory
	inv.AddItem(sword, 1)
	inv.AddItem(helmet, 1)
	inv.AddItem(healthPotion, 5) // Stack of 5 health potions
	inv.AddItem(manaPotion, 3)   // Stack of 3 mana potions
	inv.AddItem(ironOre, 25)     // Stack of 25 iron ore
	inv.AddItem(mysteriousKey, 1)

	// Add some additional items for testing drag & drop
	magicRing := &components.Item{
		ID:          7,
		Name:        "Magic Ring",
		Description: "A ring imbued with magical properties.",
		Type:        components.ItemTypeEquipment,
		Rarity:      components.ItemRarityEpic,
		Value:       500,
		IconID:      7,
		Stackable:   false,
		MaxStack:    1,
	}

	elixir := &components.Item{
		ID:          8,
		Name:        "Greater Elixir",
		Description: "A powerful magical elixir brewed by ancient alchemists that fully restores both health and mana when consumed. This legendary concoction is extremely rare and valuable.",
		Type:        components.ItemTypeConsumable,
		Rarity:      components.ItemRarityLegendary,
		Value:       200,
		IconID:      8,
		Stackable:   true,
		MaxStack:    5,
		Effects: []components.ConsumableEffect{
			{Type: "heal_hp", Value: 100, Target: "self"},
			{Type: "heal_mp", Value: 100, Target: "self"},
		},
	}

	inv.AddItem(magicRing, 1)
	inv.AddItem(elixir, 2)
}

// Update handles game logic and input
func (g *Game) Update() error {
	// Toggle inventory with 'I' key
	if inpututil.IsKeyJustPressed(ebiten.KeyI) && !g.wasIPressed {
		g.isInventoryOpen = !g.isInventoryOpen
		g.inventoryWidget.Visible = g.isInventoryOpen
		g.wasIPressed = true
	} else if !inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.wasIPressed = false
	}

	// Update inventory widget if open
	if g.isInventoryOpen {
		err := g.inventoryWidget.Update()
		if err != nil {
			return err
		}

		// Check if widget was closed via escape key
		if !g.inventoryWidget.Visible {
			g.isInventoryOpen = false
		}
	}

	return nil
}

// drawInstructions shows the test instructions
func (g *Game) drawInstructions(screen *ebiten.Image) {
	// Draw a more visible header
	ebitenutil.DebugPrintAt(screen, "================================================", 10, 10)
	ebitenutil.DebugPrintAt(screen, "         INVENTORY WIDGET TEST", 10, 25)
	ebitenutil.DebugPrintAt(screen, "================================================", 10, 40)

	instructions := []string{
		"",
		"MAIN CONTROL:",
		"   I Key - Open/Close Inventory",
		"",
		"INVENTORY CONTROLS (when open):",
		"   HOVER - Show detailed item tooltip",
		"   LEFT CLICK - Select and drag items",
		"   RIGHT CLICK - Select items",
		"   DRAG & DROP - Move items between slots",
		"   DELETE Key - Remove selected item",
		"   1-9 Keys - Split selected stack",
		"   ESC Key - Close inventory",
		"",
		"ACTION PANEL BUTTONS:",
		"   Sort Name/Type/Rarity - Sort items",
		"   Equipment/Consumable - Filter by type",
		"   Show All - Remove filters",
		"",
		"TEST ITEMS LOADED:",
		"   • Equipment: Sword, Helmet, Ring",
		"   • Consumables: Potions, Elixir",
		"   • Materials: Iron Ore (25x)",
		"   • Quest: Mysterious Key",
		"",
		"FEATURES TO TEST:",
		"   ✓ Grid layout (8x6 = 48 slots)",
		"   ✓ Rarity colors (Common=Gray, Epic=Purple, etc)",
		"   ✓ Item stacking and splitting",
		"   ✓ Hover tooltips with full item details",
		"   ✓ Sorting and filtering",
		"",
		">>> PRESS 'I' TO OPEN INVENTORY <<<",
		">>> HOVER OVER ITEMS FOR TOOLTIPS <<<",
	}

	startY := 55
	for i, line := range instructions {
		ebitenutil.DebugPrintAt(screen, line, 15, startY+i*14)
	}
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{20, 25, 30, 255})

	// Draw instruction text if inventory is closed
	if !g.isInventoryOpen {
		g.drawInstructions(screen)
	}

	// Draw inventory widget if open
	if g.isInventoryOpen {
		g.inventoryWidget.Draw(screen)
	}
}

// Layout returns the screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Set window properties
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("MyRPG - Inventory Widget Test")
	ebiten.SetWindowResizable(true)

	// Print instructions to console as well
	log.Println("=== INVENTORY WIDGET TEST ===")
	log.Println("Controls:")
	log.Println("  I - Open/Close Inventory")
	log.Println("When inventory is open:")
	log.Println("  HOVER - Show detailed item tooltips")
	log.Println("  LEFT CLICK - Select and drag items")
	log.Println("  DELETE - Remove selected item")
	log.Println("  1-9 - Split selected stack")
	log.Println("  ESC - Close inventory")
	log.Println("Action panel buttons for sorting/filtering")
	log.Println("")
	log.Println("Test includes various items with different rarities and types.")
	log.Println("Press 'I' in the game window to open inventory!")
	log.Println("Hover over items to see detailed tooltips with item information!")

	// Create and run game
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
