// Package ui provides user interface components for the game.
// This includes panels, message systems, and layout management for organizing
// the game display into distinct areas for the game world, player information,
// and command/message output.
package ui

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// Layout constants define the screen areas
const (
	ScreenWidth       = 800
	ScreenHeight      = 600
	TopPanelHeight    = 110 // Increased to 110 to properly fit all content
	BottomPanelHeight = 80

	// Game world area
	GameWorldY      = TopPanelHeight
	GameWorldHeight = ScreenHeight - TopPanelHeight - BottomPanelHeight
)

// Colors for UI panels
var (
	TopPanelColor    = color.RGBA{30, 30, 30, 255}    // Dark gray
	BottomPanelColor = color.RGBA{20, 20, 20, 255}    // Darker gray
	TextColor        = color.RGBA{255, 255, 255, 255} // White
)

// Message represents a single message with timestamp
type Message struct {
	Text      string
	Timestamp time.Time
}

// MessageSystem manages game messages and command output
type MessageSystem struct {
	messages    []Message
	maxMessages int
}

// NewMessageSystem creates a new message system
func NewMessageSystem(maxMessages int) *MessageSystem {
	return &MessageSystem{
		messages:    make([]Message, 0, maxMessages),
		maxMessages: maxMessages,
	}
}

// AddMessage adds a new message to the system
func (ms *MessageSystem) AddMessage(text string) {
	message := Message{
		Text:      text,
		Timestamp: time.Now(),
	}

	ms.messages = append(ms.messages, message)

	// Keep only the last N messages
	if len(ms.messages) > ms.maxMessages {
		ms.messages = ms.messages[len(ms.messages)-ms.maxMessages:]
	}
}

// GetRecentMessages returns the most recent messages for display
func (ms *MessageSystem) GetRecentMessages(count int) []string {
	if len(ms.messages) == 0 {
		return []string{}
	}

	start := len(ms.messages) - count
	if start < 0 {
		start = 0
	}

	result := make([]string, 0, count)
	for i := start; i < len(ms.messages); i++ {
		result = append(result, ms.messages[i].Text)
	}

	return result
}

// UIManager manages all UI panels and rendering
type UIManager struct {
	messageSystem *MessageSystem
}

// NewUIManager creates a new UI manager
func NewUIManager() *UIManager {
	return &UIManager{
		messageSystem: NewMessageSystem(50), // Keep last 50 messages
	}
}

// AddMessage adds a message to the message system
func (ui *UIManager) AddMessage(text string) {
	ui.messageSystem.AddMessage(text)
}

// GameMode represents different game modes for UI rendering
type GameMode int

const (
	ModeExploration GameMode = iota // Free movement exploration
	ModeTactical                    // Grid-based tactical combat
)

// DrawTopPanel renders the player information panel based on game mode
func (ui *UIManager) DrawTopPanel(screen *ebiten.Image, activePlayer *components.RPGStatsComponent, gameMode GameMode, partyMembers []*components.RPGStatsComponent, gridPosition string) {
	// Draw background
	vector.FillRect(screen, 0, 0, ScreenWidth, TopPanelHeight, TopPanelColor, false)

	if gameMode == ModeExploration {
		ui.drawExplorationPanel(screen, partyMembers)
	} else {
		ui.drawTacticalPanel(screen, activePlayer, gridPosition)
	}
}

// drawExplorationPanel renders the exploration mode UI
func (ui *UIManager) drawExplorationPanel(screen *ebiten.Image, partyMembers []*components.RPGStatsComponent) {
	// Header
	ebitenutil.DebugPrintAt(screen, "=== EXPLORATION MODE ===", 10, 8)
	ebitenutil.DebugPrintAt(screen, "Keys: Arrow Keys=Move, TAB=Switch Player, T/Space=Combat Mode", 10, 22)

	// Party members info (simplified: name, class, level)
	if len(partyMembers) == 0 {
		ebitenutil.DebugPrintAt(screen, "No party members", 10, 40)
		return
	}

	ebitenutil.DebugPrintAt(screen, "Party Members:", 10, 40)
	for i, member := range partyMembers {
		if member != nil {
			memberInfo := fmt.Sprintf("  %d. %s (%s Level %d)",
				i+1, member.Name, member.Job.String(), member.Level)
			ebitenutil.DebugPrintAt(screen, memberInfo, 10, 54+i*15)
		}
	}
}

// drawTacticalPanel renders the tactical mode UI
func (ui *UIManager) drawTacticalPanel(screen *ebiten.Image, activePlayer *components.RPGStatsComponent, gridPosition string) {
	// Header
	ebitenutil.DebugPrintAt(screen, "=== TACTICAL COMBAT ===", 10, 8)
	ebitenutil.DebugPrintAt(screen, "Keys: Arrow Keys=Move, TAB=Switch Unit, U=Undo, Q=Return", 10, 22)

	if activePlayer == nil {
		ebitenutil.DebugPrintAt(screen, "No active player", 10, 40)
		return
	}

	// Active player info (left side)
	playerInfo := fmt.Sprintf("Active: %s (%s Level %d)",
		activePlayer.Name, activePlayer.Job.String(), activePlayer.Level)
	ebitenutil.DebugPrintAt(screen, playerInfo, 10, 40)

	// Grid position (left side, second line)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Position: %s", gridPosition), 10, 52)

	// Health and Mana (right side)
	hpMpInfo := fmt.Sprintf("HP: %d/%d  MP: %d/%d",
		activePlayer.CurrentHP, activePlayer.MaxHP,
		activePlayer.CurrentMP, activePlayer.MaxMP)
	ebitenutil.DebugPrintAt(screen, hpMpInfo, 400, 40)

	// Combat stats (right side, second line)
	combatInfo := fmt.Sprintf("ATK: %d  DEF: %d  SPD: %d",
		activePlayer.Attack, activePlayer.Defense, activePlayer.Speed)
	ebitenutil.DebugPrintAt(screen, combatInfo, 400, 52)

	// HP Bar (visual representation) - moved higher and made more compact
	ui.drawHealthBar(screen, 10, 68, 200, 10, activePlayer.CurrentHP, activePlayer.MaxHP)

	// MP Bar (visual representation) - moved higher
	ui.drawManaBar(screen, 220, 68, 150, 10, activePlayer.CurrentMP, activePlayer.MaxMP)
}

// drawHealthBar draws a visual health bar
func (ui *UIManager) drawHealthBar(screen *ebiten.Image, x, y, width, height float32, current, max int) {
	// Background (dark red)
	vector.FillRect(screen, x, y, width, height, color.RGBA{60, 20, 20, 255}, false)

	// Health bar (bright red to green based on percentage)
	if max > 0 {
		percentage := float32(current) / float32(max)
		barWidth := width * percentage

		// Color changes from red to yellow to green
		var barColor color.RGBA
		if percentage > 0.6 {
			barColor = color.RGBA{0, 200, 0, 255} // Green
		} else if percentage > 0.3 {
			barColor = color.RGBA{200, 200, 0, 255} // Yellow
		} else {
			barColor = color.RGBA{200, 0, 0, 255} // Red
		}

		vector.FillRect(screen, x, y, barWidth, height, barColor, false)
	}

	// Border
	vector.StrokeRect(screen, x, y, width, height, 1, color.RGBA{255, 255, 255, 255}, false)
}

// drawManaBar draws a visual mana bar
func (ui *UIManager) drawManaBar(screen *ebiten.Image, x, y, width, height float32, current, max int) {
	// Background (dark blue)
	vector.FillRect(screen, x, y, width, height, color.RGBA{20, 20, 60, 255}, false)

	// Mana bar (blue)
	if max > 0 {
		percentage := float32(current) / float32(max)
		barWidth := width * percentage
		vector.FillRect(screen, x, y, barWidth, height, color.RGBA{0, 100, 255, 255}, false)
	}

	// Border
	vector.StrokeRect(screen, x, y, width, height, 1, color.RGBA{255, 255, 255, 255}, false)
}

// DrawBottomPanel renders the command output and messages panel
func (ui *UIManager) DrawBottomPanel(screen *ebiten.Image) {
	// Draw background
	bottomY := float32(ScreenHeight - BottomPanelHeight)
	vector.FillRect(screen, 0, bottomY, ScreenWidth, BottomPanelHeight, BottomPanelColor, false)

	// Get recent messages (up to 4 lines)
	messages := ui.messageSystem.GetRecentMessages(4)

	// If no messages, show default instructions
	if len(messages) == 0 {
		messages = []string{"Use arrow keys to move active player, TAB to switch between players"}
	}

	// Display messages
	for i, message := range messages {
		y := int(bottomY) + 10 + (i * 15)
		ebitenutil.DebugPrintAt(screen, message, 10, y)
	}
}

// DrawGameWorldBackground fills the game world area with a background color
func (ui *UIManager) DrawGameWorldBackground(screen *ebiten.Image) {
	// Draw a thin separator line between top panel and game world
	separatorColor := color.RGBA{0, 0, 0, 255} // Black line
	vector.FillRect(screen, 0, GameWorldY, ScreenWidth, 2, separatorColor, false)

	// Fill game world area with a neutral background color
	gameWorldColor := color.RGBA{50, 70, 50, 255} // Dark green
	// Start game world background 2 pixels below top panel
	vector.FillRect(screen, 0, GameWorldY+2, ScreenWidth, GameWorldHeight-2, gameWorldColor, false)
}

// GetGameWorldBounds returns the bounds of the game world area
func (ui *UIManager) GetGameWorldBounds() (x, y, width, height int) {
	return 0, GameWorldY + 2, ScreenWidth, GameWorldHeight - 2
}

// DrawBattleMenu renders the battle selection menu overlay
func (ui *UIManager) DrawBattleMenu(screen *ebiten.Image, battleText string) {
	if battleText == "" {
		return
	}

	// Draw semi-transparent overlay
	overlayColor := color.RGBA{0, 0, 0, 180}
	vector.FillRect(screen, 0, 0, ScreenWidth, ScreenHeight, overlayColor, false)

	// Calculate menu dimensions
	menuWidth := float32(300)
	menuHeight := float32(200)
	menuX := (ScreenWidth - menuWidth) / 2
	menuY := (ScreenHeight - menuHeight) / 2

	// Draw menu background
	menuBgColor := color.RGBA{40, 40, 40, 255}
	vector.FillRect(screen, menuX, menuY, menuWidth, menuHeight, menuBgColor, false)

	// Draw menu border
	borderColor := color.RGBA{200, 200, 200, 255}
	vector.StrokeRect(screen, menuX, menuY, menuWidth, menuHeight, 2, borderColor, false)

	// Draw battle text
	lines := []string{}
	current := ""
	for _, char := range battleText {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}

	// Render text lines
	for i, line := range lines {
		textX := int(menuX) + 20
		textY := int(menuY) + 30 + (i * 20)
		ebitenutil.DebugPrintAt(screen, line, textX, textY)
	}
}
