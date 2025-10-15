// Character Stats Widget Test Program
// This program tests the character statistics widget functionality
// including category navigation, visual layout, and input handling.
package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/ui"
)

// Game represents the test game state
type Game struct {
	uiManager *ui.UIManager
	character *components.RPGStatsComponent
}

// NewGame creates a new test game instance
func NewGame() *Game {
	uiManager := ui.NewUIManager()

	// Create a test character with sample data
	character := createTestCharacter()

	return &Game{
		uiManager: uiManager,
		character: character,
	}
}

// createTestCharacter creates a sample character for testing
func createTestCharacter() *components.RPGStatsComponent {
	// Create a level 15 Warrior with some combat experience
	character := components.NewRPGStatsComponent("Thorin Ironforge", components.JobWarrior, 15)

	// Simulate some damage and MP usage
	character.TakeDamage(25) // Not at full health
	character.UseMana(10)    // Used some mana

	// Add some experience (not quite ready to level up)
	character.GainExperience(1200)

	// Simulate tactical combat state
	character.ConsumeMovement(2) // Moved 2 tiles

	return character
}

// Update updates the game state
func (g *Game) Update() error {
	// Update UI Manager (handles all widget updates)
	g.uiManager.Update()

	// Handle 'C' key to show/hide character stats
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		if g.uiManager.IsPopupVisible() {
			g.uiManager.HideCharacterStats()
		} else {
			g.uiManager.ShowCharacterStats(g.character)
		}
	}

	// Handle 'R' key to reset character (for testing different states)
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.character = createTestCharacter()
		// Update the widget if it's currently visible
		if g.uiManager.IsPopupVisible() {
			g.uiManager.ShowCharacterStats(g.character)
		}
	}

	// Handle 'H' key to heal character (for testing health bars)
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.character.Heal(50)
		g.character.RestoreMana(25)
	}

	// Handle 'D' key to damage character (for testing health bars)
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.character.TakeDamage(30)
		g.character.UseMana(15)
	}

	// Handle 'L' key to add experience (for testing level progression)
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.character.GainExperience(200)
	}

	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{26, 26, 51, 255}) // Dark blue background

	// Draw test instructions
	instructions := []string{
		"Character Stats Widget Test Program",
		"",
		"Controls:",
		"  C - Show/Hide Character Stats",
		"  R - Reset Character (new random stats)",
		"  H - Heal Character (+50 HP, +25 MP)",
		"  D - Damage Character (-30 HP, -15 MP)",
		"  L - Gain Experience (+200 XP)",
		"",
		"When stats widget is open:",
		"  ← → - Switch between stat categories",
		"  TAB - Next category",
		"  Shift+TAB - Previous category",
		"  ESC - Close widget",
		"",
		"Character: " + g.character.Name,
		"Level: " + fmt.Sprintf("%d", g.character.Level),
		"HP: " + fmt.Sprintf("%d/%d", g.character.CurrentHP, g.character.MaxHP),
		"MP: " + fmt.Sprintf("%d/%d", g.character.CurrentMP, g.character.MaxMP),
		"XP: " + fmt.Sprintf("%d/%d", g.character.Experience, g.character.ExpToNext),
	}

	for i, instruction := range instructions {
		y := 20 + i*16
		ebitenutil.DebugPrintAt(screen, instruction, 20, y)
	}

	// Draw status of popup visibility
	if g.uiManager.IsPopupVisible() {
		ebitenutil.DebugPrintAt(screen, "Status: Character Stats OPEN", 20, 350)
	} else {
		ebitenutil.DebugPrintAt(screen, "Status: Character Stats CLOSED (Press C to open)", 20, 350)
	}

	// Draw all UI elements (including character stats widget)
	g.uiManager.DrawPopups(screen)
}

// Layout returns the game screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 600 // Same as your main game size
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(1000, 600)
	ebiten.SetWindowTitle("Character Stats Widget Test")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	log.Println("Starting Character Stats Widget test...")
	log.Println("Press 'C' to open character stats widget")
	log.Println("Use arrow keys or TAB to navigate between categories")
	log.Println("Press ESC to close the widget")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
