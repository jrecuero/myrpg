package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jrecuero/myrpg/internal/ui"
)

// TestInfoPopupGame demonstrates the PopupInfoWidget functionality
type TestInfoPopupGame struct {
	uiManager *ui.UIManager
}

// NewTestInfoPopupGame creates a new test game for info popup
func NewTestInfoPopupGame() *TestInfoPopupGame {
	return &TestInfoPopupGame{
		uiManager: ui.NewUIManager(),
	}
}

// Update handles the game logic
func (g *TestInfoPopupGame) Update() error {
	g.uiManager.Update()

	// Test trigger - press I to show info popup
	if ebiten.IsKeyPressed(ebiten.KeyI) && !g.uiManager.IsPopupVisible() {
		g.showTestInfoPopup()
	}

	// Exit test with ESC when no popup is open
	if ebiten.IsKeyPressed(ebiten.KeyEscape) && !g.uiManager.IsPopupVisible() {
		return ebiten.Termination
	}

	return nil
}

// Draw renders the game
func (g *TestInfoPopupGame) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{40, 40, 40, 255})

	// Draw instructions
	if !g.uiManager.IsPopupVisible() {
		g.drawInstructions(screen)
	}

	// Draw popups on top
	g.uiManager.DrawPopups(screen)
}

// Layout returns the game screen size
func (g *TestInfoPopupGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

// drawInstructions shows the test instructions
func (g *TestInfoPopupGame) drawInstructions(screen *ebiten.Image) {
	instructions := []string{
		"🧪 Testing PopupInfoWidget",
		"",
		"📋 Test Instructions:",
		"• Press 'I' to show info popup",
		"• Use ↑↓ arrows to scroll (if content is long)",
		"• Press ESC to close popup",
		"• ESC also exits the test when no popup is open",
		"",
		"✅ Expected Behavior:",
		"  - Info popup displays with title and scrollable content",
		"  - Arrow keys scroll content when popup is open",
		"  - Arrow keys do NOT affect anything when popup is closed",
		"  - ESC closes the popup",
		"  - Input is blocked from this test program while popup is open",
	}

	y := 50
	for _, line := range instructions {
		if line == "" {
			y += 20
			continue
		}

		ebitenutil.DebugPrintAt(screen, line, 20, y)
		y += 20
	}
}

// showTestInfoPopup displays a comprehensive test info popup
func (g *TestInfoPopupGame) showTestInfoPopup() {
	testContent := `POPUP INFO WIDGET TEST

This is a test of the PopupInfoWidget component!

FEATURES BEING TESTED:
• Multi-line text display
• Automatic text wrapping (future enhancement)
• Scrollable content when text exceeds widget height
• Title display at the top
• Close functionality with ESC key
• Input blocking (no external input while open)

SCROLLING TEST:
This content is intentionally long to test the scrolling
functionality of the widget. You should be able to use
the UP and DOWN arrow keys to scroll through this content.

Line 1 of scrolling test content...
Line 2 of scrolling test content...
Line 3 of scrolling test content...
Line 4 of scrolling test content...
Line 5 of scrolling test content...
Line 6 of scrolling test content...
Line 7 of scrolling test content...
Line 8 of scrolling test content...
Line 9 of scrolling test content...
Line 10 of scrolling test content...
Line 11 of scrolling test content...
Line 12 of scrolling test content...
Line 13 of scrolling test content...
Line 14 of scrolling test content...
Line 15 of scrolling test content...

STYLING TEST:
The widget should display with:
✓ Semi-transparent dark background
✓ Gray border around the popup
✓ Yellow title text
✓ White content text
✓ Shadow effect behind popup
✓ Scrollbar (if content overflows)

INPUT BLOCKING TEST:
While this popup is open, pressing any keys should only
affect the popup (scrolling with arrows, closing with ESC).
No other input should be processed by the background.

USAGE SCENARIOS:
This widget is perfect for:
• Game help and instructions
• Character attribute displays  
• Item descriptions and lore
• System messages and notifications
• Settings and configuration info

Press ESC to close this popup and return to the test.`

	g.uiManager.ShowInfoPopup(
		"Info Widget Test - Comprehensive Demo",
		testContent,
		func() {
			log.Println("✅ Info popup closed successfully")
		},
	)
}

func main() {
	// Configure window
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("MyRPG - PopupInfoWidget Test")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	log.Println("🧪 Starting PopupInfoWidget test...")
	log.Println("📋 Press 'I' to show info popup")
	log.Println("📋 Use arrow keys to scroll when popup is open")
	log.Println("📋 Press ESC to close popup or exit test")

	// Create and run test game
	game := NewTestInfoPopupGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
