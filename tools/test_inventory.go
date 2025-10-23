package main

// Test program for the Inventory Widget system
// This program creates a test scenario with items and demonstrates
// the inventory widget functionality including drag & drop, tooltips,
// and item management.

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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
		Description: "An ancient key with unknown purpose. Quest item.",
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
		Description: "A powerful elixir that fully restores health and mana.",
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
	if ebiten.IsKeyPressed(ebiten.KeyI) && !g.wasIPressed {
		g.isInventoryOpen = !g.isInventoryOpen
		g.inventoryWidget.Visible = g.isInventoryOpen
		g.wasIPressed = true
	} else if !ebiten.IsKeyPressed(ebiten.KeyI) {
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

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{20, 25, 30, 255})

	// Draw instruction text
	if !g.isInventoryOpen {
		// TODO: Add text rendering when font system is available
		// For now, inventory widget provides visual feedback
		// Instructions: Press 'I' to open inventory
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

	// Create and run game
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
