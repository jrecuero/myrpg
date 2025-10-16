package main

import (
	"fmt"
	"image/color"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jrecuero/myrpg/internal/ui"
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 600
)

// Game represents the test game state
type Game struct {
	uiManager *ui.UIManager
	assetPath string
}

// NewGame creates a new test game instance
func NewGame() *Game {
	uiManager := ui.NewUIManager()

	// Get the path to assets directory
	assetPath, err := filepath.Abs("assets/dialogs")
	if err != nil {
		log.Printf("Warning: Could not find assets directory: %v", err)
		assetPath = "assets/dialogs"
	}

	return &Game{
		uiManager: uiManager,
		assetPath: assetPath,
	}
}

// Update handles the game logic
func (g *Game) Update() error {
	// Update UI Manager (handles all widget updates)
	g.uiManager.Update()

	// Handle 'D' key to show dialog
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if g.uiManager.IsDialogVisible() {
			g.uiManager.HideDialog()
		} else {
			// Start Elder dialog
			err := g.uiManager.ShowDialog(g.assetPath, "characters.json", "town_elder.json", "start")
			if err != nil {
				log.Printf("Failed to start dialog: %v", err)
			}
		}
	}

	// Handle '1' key for merchant dialog
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		if !g.uiManager.IsDialogVisible() {
			err := g.uiManager.ShowDialog(g.assetPath, "characters.json", "merchant.json", "start")
			if err != nil {
				log.Printf("Failed to start merchant dialog: %v", err)
			}
		}
	}

	// Handle '2' key for testing variables
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.uiManager.SetDialogVariable("player_name", "TestHero")
		g.uiManager.SetDialogVariable("elder_trust", 10)
		log.Println("Set test variables: player_name=TestHero, elder_trust=10")
	}

	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw instructions
	g.drawInstructions(screen)

	// Draw dialog variables for debugging
	g.drawDialogDebugInfo(screen)

	// Draw UI (including dialog)
	g.uiManager.DrawPopups(screen)
}

// Layout returns the game screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// drawInstructions shows the test instructions
func (g *Game) drawInstructions(screen *ebiten.Image) {
	instructions := []string{
		"=== DIALOG WIDGET TEST ===",
		"",
		"Controls:",
		"D - Show/Hide Elder Dialog",
		"1 - Show Merchant Dialog",
		"2 - Set Test Variables",
		"",
		"Dialog Controls (when dialog is open):",
		"SPACE/ENTER - Continue/Select choice",
		"UP/DOWN - Navigate choices",
		"TAB - Toggle typewriter speed",
		"ESC - Close dialog",
		"",
		"Features Tested:",
		"✓ External JSON dialog scripts",
		"✓ Character definitions with portraits",
		"✓ Branching dialog trees",
		"✓ Multiple choice selections",
		"✓ Variable system with substitution",
		"✓ Typewriter text effect",
		"✓ Conditional dialog branches",
		"✓ Dialog action execution",
	}

	for i, line := range instructions {
		ebitenutil.DebugPrintAt(screen, line, 10, 10+i*15)
	}
}

// drawDialogDebugInfo shows current dialog variables
func (g *Game) drawDialogDebugInfo(screen *ebiten.Image) {
	debugY := 350

	ebitenutil.DebugPrintAt(screen, "=== Dialog Variables ===", 10, debugY)
	debugY += 20

	// Check some common variables
	variables := []string{
		"met_elder", "elder_trust", "knows_prophecy", "player_name",
		"merchant_met", "merchant_reputation", "bandit_quest_offered",
	}

	for _, varName := range variables {
		if value, exists := g.uiManager.GetDialogVariable(varName); exists {
			ebitenutil.DebugPrintAt(screen,
				fmt.Sprintf("%s: %v", varName, value),
				10, debugY)
			debugY += 15
		}
	}

	// Show dialog state
	debugY += 10
	if g.uiManager.IsDialogVisible() {
		ebitenutil.DebugPrintAt(screen, "Dialog Status: VISIBLE", 10, debugY)
	} else {
		ebitenutil.DebugPrintAt(screen, "Dialog Status: Hidden", 10, debugY)
	}
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Dialog Widget Test")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	log.Println("Starting Dialog Widget test...")
	log.Println("Controls:")
	log.Println("  D - Show/Hide Elder Dialog")
	log.Println("  1 - Show Merchant Dialog")
	log.Println("  2 - Set Test Variables")
	log.Println("")
	log.Println("Dialog Files:")
	log.Println("  assets/dialogs/characters.json - Character definitions")
	log.Println("  assets/dialogs/town_elder.json - Elder conversation")
	log.Println("  assets/dialogs/merchant.json - Merchant conversation")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
