package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jrecuero/myrpg/internal/ui"
)

type InputBlockingTestGame struct {
	uiManager     *ui.UIManager
	playerX       int
	playerY       int
	moveAttempts  int
	blockedInputs int
}

func NewInputBlockingTestGame() *InputBlockingTestGame {
	uiManager := ui.NewUIManager()
	return &InputBlockingTestGame{
		uiManager: uiManager,
		playerX:   5,
		playerY:   5,
	}
}

func (g *InputBlockingTestGame) Update() error {
	// Update UI manager first (processes popup input)
	g.uiManager.Update()

	// Show popup with P key
	if inpututil.IsKeyJustPressed(ebiten.KeyP) && !g.uiManager.IsPopupVisible() {
		options := []string{
			"Move North",
			"Move South",
			"Move East",
			"Move West",
			"Stay Here",
		}

		g.uiManager.ShowSelectionPopup("Player Actions", options,
			func(index int, option string) {
				fmt.Printf("✅ Selected: %s (index %d)\n", option, index)
				log.Printf("Player chose: %s", option)
			},
			func() {
				fmt.Printf("❌ Action cancelled\n")
				log.Printf("Player cancelled action")
			})

		fmt.Println("🎮 Popup shown! Try using arrow keys...")
		fmt.Println("   - Arrow keys should ONLY affect popup")
		fmt.Println("   - Player should NOT move while popup is open")
		return nil
	}

	// CRITICAL: Block game input when popup is visible
	if g.uiManager.IsPopupVisible() {
		// Count blocked input attempts
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) ||
			ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.blockedInputs++
			if g.blockedInputs%60 == 1 { // Log every second (assuming 60 FPS)
				fmt.Printf("🚫 BLOCKED: Game input blocked while popup is open (blocked %d times)\n", g.blockedInputs)
			}
		}
		return nil // Skip game logic entirely
	}

	// Normal game input (only when popup is closed)
	moved := false
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.playerY--
		moved = true
		fmt.Printf("⬆️  Player moved UP to (%d, %d)\n", g.playerX, g.playerY)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.playerY++
		moved = true
		fmt.Printf("⬇️  Player moved DOWN to (%d, %d)\n", g.playerX, g.playerY)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.playerX--
		moved = true
		fmt.Printf("⬅️  Player moved LEFT to (%d, %d)\n", g.playerX, g.playerY)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.playerX++
		moved = true
		fmt.Printf("➡️  Player moved RIGHT to (%d, %d)\n", g.playerX, g.playerY)
	}

	if moved {
		g.moveAttempts++
		fmt.Printf("✅ ALLOWED: Player movement #%d completed\n", g.moveAttempts)
	}

	return nil
}

func (g *InputBlockingTestGame) Draw(screen *ebiten.Image) {
	// Simple game state display
	if !g.uiManager.IsPopupVisible() {
		ebiten.SetWindowTitle(fmt.Sprintf("Input Blocking Test - Player at (%d,%d) - Press P for popup", g.playerX, g.playerY))
	} else {
		ebiten.SetWindowTitle("Input Blocking Test - POPUP ACTIVE - Arrow keys should only affect popup!")
	}

	// Draw popup if visible
	g.uiManager.DrawPopups(screen)

	// Instructions (always visible in background)
	instruction := fmt.Sprintf("Player: (%d,%d) | Moves: %d | Blocked: %d | Press P for popup",
		g.playerX, g.playerY, g.moveAttempts, g.blockedInputs)

	// We can't easily draw text without more complex setup, so we'll use the title
	if !g.uiManager.IsPopupVisible() {
		ebiten.SetWindowTitle("Input Test - " + instruction + " | Press P for popup")
	}
}

func (g *InputBlockingTestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	fmt.Println("🧪 Testing Input Blocking with Popup Widget...")
	fmt.Println()
	fmt.Println("📋 Test Procedure:")
	fmt.Println("1. Use arrow keys to move player (should work)")
	fmt.Println("2. Press P to show popup")
	fmt.Println("3. Use arrow keys while popup is open (should NOT move player)")
	fmt.Println("4. Navigate popup with arrows, select with Enter, or cancel with Esc")
	fmt.Println("5. After popup closes, arrow keys should move player again")
	fmt.Println()
	fmt.Println("✅ Expected Behavior:")
	fmt.Println("   - Player moves when popup is closed")
	fmt.Println("   - Player DOES NOT move when popup is open")
	fmt.Println("   - Arrow keys only control popup navigation when popup is open")
	fmt.Println()
	fmt.Println("Starting test game...")

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Input Blocking Test Game")

	game := NewInputBlockingTestGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
