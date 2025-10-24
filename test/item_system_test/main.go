package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/systems"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game represents the test game state
type Game struct {
	itemRegistry      *components.ItemRegistry
	consumableManager *systems.ConsumableManager
	testPlayer        *ecs.Entity
	testItems         []*components.Item
	currentItemIndex  int
	message           string
	lastKeyTime       int
}

// NewGame creates a new test game instance
func NewGame() *Game {
	// Initialize item system
	components.InitializeItemSystem()

	g := &Game{
		itemRegistry:      components.GlobalItemRegistry,
		consumableManager: systems.NewConsumableManager(),
		currentItemIndex:  0,
		message:           "Item System Test - Press SPACE to cycle items, ENTER to use consumable",
	}

	// Create test player
	g.createTestPlayer()

	// Get test items from registry
	g.getTestItems()

	return g
}

// createTestPlayer creates a test player entity
func (g *Game) createTestPlayer() {
	g.testPlayer = ecs.NewEntity("TestPlayer")

	// Add RPG stats component
	stats := components.NewRPGStatsComponent("Test Warrior", components.JobWarrior, 10)

	// Simulate some damage for testing healing
	stats.CurrentHP = 75 // Not at full health
	stats.CurrentMP = 40 // Not at full mana

	g.testPlayer.AddComponent(ecs.ComponentRPGStats, stats)
}

// getTestItems retrieves test items from the registry
func (g *Game) getTestItems() {
	g.testItems = []*components.Item{
		g.itemRegistry.CreateItem(200), // Health Potion
		g.itemRegistry.CreateItem(210), // Mana Potion
		g.itemRegistry.CreateItem(1),   // Iron Sword
		g.itemRegistry.CreateItem(301), // Magic Crystal
	}
}

// Update handles game logic and input
func (g *Game) Update() error {
	g.lastKeyTime++

	// Prevent key spam by requiring delay between inputs
	if g.lastKeyTime < 10 {
		return nil
	}

	// Cycle through items with SPACE
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.currentItemIndex = (g.currentItemIndex + 1) % len(g.testItems)
		g.message = fmt.Sprintf("Selected: %s", g.testItems[g.currentItemIndex].Name)
		g.lastKeyTime = 0
	}

	// Use consumable with ENTER
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.useCurrentItem()
		g.lastKeyTime = 0
	}

	// Reset player stats with R
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.resetPlayerStats()
		g.message = "Player stats reset!"
		g.lastKeyTime = 0
	}

	return nil
}

// useCurrentItem attempts to use the currently selected item
func (g *Game) useCurrentItem() {
	if g.currentItemIndex >= len(g.testItems) {
		return
	}

	item := g.testItems[g.currentItemIndex]

	// Check if item is consumable
	if item.Type == components.ItemTypeConsumable {
		// Check if player can use this consumable
		if g.consumableManager.CanUseConsumable(item, g.testPlayer, g.testPlayer) {
			err := g.consumableManager.UseConsumable(item, g.testPlayer, g.testPlayer)
			if err != nil {
				g.message = fmt.Sprintf("Failed to use %s: %v", item.Name, err)
			} else {
				g.message = fmt.Sprintf("Used %s successfully!", item.Name)
			}
		} else {
			g.message = fmt.Sprintf("Cannot use %s - requirements not met", item.Name)
		}
	} else {
		g.message = fmt.Sprintf("%s is not consumable - it's %s", item.Name, item.Type.String())
	}
}

// resetPlayerStats resets player to partial health/mana for testing
func (g *Game) resetPlayerStats() {
	if stats, exists := g.testPlayer.GetComponent(ecs.ComponentRPGStats); exists {
		playerStats := stats.(*components.RPGStatsComponent)
		playerStats.CurrentHP = 75
		playerStats.CurrentMP = 40
	}
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{20, 25, 30, 255})

	// Draw title
	ebitenutil.DebugPrintAt(screen, "=== ITEM SYSTEM TEST ===", 10, 10)

	// Draw controls
	ebitenutil.DebugPrintAt(screen, "Controls:", 10, 40)
	ebitenutil.DebugPrintAt(screen, "  SPACE - Cycle through items", 10, 60)
	ebitenutil.DebugPrintAt(screen, "  ENTER - Use consumable item", 10, 80)
	ebitenutil.DebugPrintAt(screen, "  R - Reset player stats", 10, 100)

	// Draw current item info
	if g.currentItemIndex < len(g.testItems) {
		item := g.testItems[g.currentItemIndex]

		y := 140
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Current Item [%d/%d]:", g.currentItemIndex+1, len(g.testItems)), 10, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Name: %s", item.Name), 10, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Type: %s", item.Type.String()), 10, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rarity: %s", item.Rarity.String()), 10, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Value: %d gold", item.Value), 10, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Description: %s", item.Description), 10, y)

		// Show effects for consumables
		if item.Type == components.ItemTypeConsumable {
			y += 30
			ebitenutil.DebugPrintAt(screen, "Effects:", 10, y)
			for _, effect := range item.Effects {
				y += 20
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s: %d (%s)", effect.Type, effect.Value, effect.Target), 10, y)
			}
		}

		// Show equipment stats
		if item.Type == components.ItemTypeEquipment && item.Equipment != nil {
			y += 30
			ebitenutil.DebugPrintAt(screen, "Equipment Stats:", 10, y)
			y += 20
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  Slot: %s", item.Equipment.Slot.String()), 10, y)
			y += 20
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  Attack Bonus: %d", item.Equipment.Stats.AttackBonus), 10, y)
			y += 20
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  Crit Chance Bonus: %d%%", item.Equipment.Stats.CritChanceBonus), 10, y)
		}
	}

	// Draw player stats
	if stats, exists := g.testPlayer.GetComponent(ecs.ComponentRPGStats); exists {
		playerStats := stats.(*components.RPGStatsComponent)

		y := 400
		ebitenutil.DebugPrintAt(screen, "=== PLAYER STATUS ===", 400, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Name: %s", playerStats.Name), 400, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Job: %s Level %d", playerStats.Job.String(), playerStats.Level), 400, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d/%d", playerStats.CurrentHP, playerStats.MaxHP), 400, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MP: %d/%d", playerStats.CurrentMP, playerStats.MaxMP), 400, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attack: %d", playerStats.Attack), 400, y)
		y += 20
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Defense: %d", playerStats.Defense), 400, y)
	}

	// Draw message
	ebitenutil.DebugPrintAt(screen, g.message, 10, screenHeight-40)

	// Draw item registry info
	ebitenutil.DebugPrintAt(screen, "=== ITEM REGISTRY ===", 400, 140)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Total Items: %d", g.itemRegistry.GetItemCount()), 400, 160)
	ebitenutil.DebugPrintAt(screen, "Available Items:", 400, 180)

	// Get all items and sort by ID to ensure stable display
	allItems := g.itemRegistry.GetAllItems()
	var sortedIDs []int
	for id := range allItems {
		sortedIDs = append(sortedIDs, id)
	}

	// Sort IDs to ensure consistent display order
	for i := 0; i < len(sortedIDs); i++ {
		for j := i + 1; j < len(sortedIDs); j++ {
			if sortedIDs[i] > sortedIDs[j] {
				sortedIDs[i], sortedIDs[j] = sortedIDs[j], sortedIDs[i]
			}
		}
	}

	y := 200
	for _, id := range sortedIDs {
		if y > 350 { // Prevent overflow
			ebitenutil.DebugPrintAt(screen, "... and more", 400, y)
			break
		}
		item := allItems[id]
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %d: %s (%s)", id, item.Name, item.Type.String()), 400, y)
		y += 15
	}
}

// Layout returns the screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Set window properties
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("MyRPG - Item System Test")

	log.Println("=== ITEM SYSTEM TEST ===")
	log.Println("Controls:")
	log.Println("  SPACE - Cycle through items")
	log.Println("  ENTER - Use consumable item")
	log.Println("  R - Reset player stats")
	log.Println()
	log.Println("Testing item registry, item creation, and consumable system")

	// Create and run game
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
