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
	"github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/logger"
)

// Use constants from the constants package
const (
	ScreenWidth       = constants.ScreenWidth
	ScreenHeight      = constants.ScreenHeight
	TopPanelHeight    = constants.TopPanelHeight
	BottomPanelHeight = constants.BottomPanelHeight

	// Game world area
	GameWorldY      = constants.GameWorldY
	GameWorldHeight = constants.GameWorldHeight
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
	// Also log UI messages to the file for debugging
	logger.UI("UI Message: %s", text)

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
	messageSystem  *MessageSystem
	popupSelection *PopupSelectionWidget // Reusable popup selection widget
	popupInfo      *PopupInfoWidget      // Reusable popup info widget
	characterStats *CharacterStatsWidget // Character statistics widget
	equipment      *EquipmentWidget      // Equipment management widget
	dialog         *DialogWidget         // Dialog conversation widget
	inventory      *InventoryWidget      // Inventory management widget
}

// NewUIManager creates a new UI manager
func NewUIManager() *UIManager {
	// Create popup widgets centered in screen
	popupX := (ScreenWidth - 300) / 2
	popupY := (ScreenHeight - 200) / 2

	popupSelection := NewPopupSelectionWidget("", []string{}, popupX, popupY, 300, 200)
	popupInfo := NewPopupInfoWidget("", "", popupX, popupY, 400, 300)

	// Create character stats widget centered
	statsX := (ScreenWidth - StatsWidgetWidth) / 2
	statsY := (ScreenHeight - StatsWidgetHeight) / 2
	characterStats := NewCharacterStatsWidget(statsX, statsY, nil)

	// Create equipment widget centered
	equipment := NewEquipmentWidget(ScreenWidth, ScreenHeight, nil, nil, nil)

	// Create dialog widget centered
	dialog := NewDialogWidget(ScreenWidth, ScreenHeight)

	// Inventory widget will be initialized when player entity is available

	return &UIManager{
		messageSystem:  NewMessageSystem(50), // Keep last 50 messages
		popupSelection: popupSelection,
		popupInfo:      popupInfo,
		characterStats: characterStats,
		equipment:      equipment,
		dialog:         dialog,
		inventory:      nil, // Will be set when player entity is available
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
	vector.FillRect(screen, 0, 0, constants.BackgroundWidth, TopPanelHeight, TopPanelColor, false)

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
	vector.FillRect(screen, 0, bottomY, constants.BackgroundWidth, BottomPanelHeight, BottomPanelColor, false)

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
	vector.FillRect(screen, 0, GameWorldY, constants.BackgroundWidth, 2, separatorColor, false)

	// Fill only the game area (800px) with background color, leave right area for UI
	gameWorldColor := color.RGBA{50, 70, 50, 255} // Dark green
	vector.FillRect(screen, 0, GameWorldY+2, constants.BackgroundWidth, GameWorldHeight-2, gameWorldColor, false)
}

// GetGameWorldBounds returns the bounds of the game world area
func (ui *UIManager) GetGameWorldBounds() (x, y, width, height int) {
	return 0, GameWorldY + 2, constants.BackgroundWidth, GameWorldHeight - 2
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

// Update handles input for UI components including popup widgets
// Returns true if a popup consumed the ESC key in this frame
func (ui *UIManager) Update() bool {
	escConsumed := false

	if ui.popupSelection != nil {
		if ui.popupSelection.Update() {
			escConsumed = true
		}
	}
	if ui.popupInfo != nil {
		if ui.popupInfo.Update() {
			escConsumed = true
		}
	}
	if ui.characterStats != nil {
		if ui.characterStats.Update() {
			escConsumed = true
		}
	}
	if ui.equipment != nil {
		if ui.equipment.Update() {
			escConsumed = true
		}
	}
	if ui.dialog != nil {
		if ui.dialog.Update() {
			escConsumed = true
		}
	}
	if ui.inventory != nil {
		if ui.inventory.Update() {
			escConsumed = true
		}
		// Check if inventory was closed
		if !ui.inventory.Visible {
			ui.inventory = nil
		}
	}

	return escConsumed
}

// DrawPopups renders any active popup widgets on top of other UI elements
func (ui *UIManager) DrawPopups(screen *ebiten.Image) {
	if ui.popupSelection != nil {
		ui.popupSelection.Draw(screen)
	}
	if ui.popupInfo != nil {
		ui.popupInfo.Draw(screen)
	}
	if ui.characterStats != nil {
		ui.characterStats.Draw(screen)
	}
	if ui.equipment != nil {
		ui.equipment.Draw(screen)
	}
	if ui.dialog != nil {
		ui.dialog.Draw(screen)
	}
	if ui.inventory != nil {
		ui.inventory.Draw(screen)
	}
}

// ShowSelectionPopup displays a popup with selectable options
func (ui *UIManager) ShowSelectionPopup(title string, options []string, onSelection func(int, string), onCancel func()) {
	if ui.popupSelection != nil {
		ui.popupSelection.OnSelection = onSelection
		ui.popupSelection.OnCancel = onCancel
		ui.popupSelection.Show(title, options)
	}
}

// HideSelectionPopup closes the selection popup
func (ui *UIManager) HideSelectionPopup() {
	if ui.popupSelection != nil {
		ui.popupSelection.Hide()
	}
}

// ShowInfoPopup displays a popup with information/text content
func (ui *UIManager) ShowInfoPopup(title string, content string, onClose func()) {
	if ui.popupInfo != nil {
		ui.popupInfo.OnClose = onClose
		ui.popupInfo.Show(title, content)
	}
}

// HideInfoPopup closes the info popup
func (ui *UIManager) HideInfoPopup() {
	if ui.popupInfo != nil {
		ui.popupInfo.Hide()
	}
}

// ShowCharacterStats displays the character statistics widget
func (ui *UIManager) ShowCharacterStats(character *components.RPGStatsComponent) {
	if ui.characterStats != nil {
		ui.characterStats.SetCharacter(character)
		ui.characterStats.Show()
	}
}

// HideCharacterStats closes the character statistics widget
func (ui *UIManager) HideCharacterStats() {
	if ui.characterStats != nil {
		ui.characterStats.Hide()
	}
}

// ShowEquipment displays the equipment widget
func (ui *UIManager) ShowEquipment(equipmentComp *components.EquipmentComponent, character *components.RPGStatsComponent, entity *ecs.Entity) {
	if ui.equipment != nil {
		ui.equipment.EquipmentComp = equipmentComp
		ui.equipment.CharacterStats = character
		ui.equipment.Entity = entity

		// Set up available equipment for testing (mock system)
		if len(ui.equipment.AvailableEquipment) == 0 {
			ui.equipment.SetAvailableEquipment(CreateMockEquipmentSet())
		}

		ui.equipment.Show()
	}
}

// HideEquipment closes the equipment widget
func (ui *UIManager) HideEquipment() {
	if ui.equipment != nil {
		ui.equipment.Hide()
	}
}

// ShowDialog starts a dialog conversation
func (ui *UIManager) ShowDialog(scriptsPath, charactersFile, dialogFile, startNode string) error {
	if ui.dialog == nil {
		return fmt.Errorf("dialog widget not initialized")
	}

	// Load characters if not already loaded
	if len(ui.dialog.Characters) == 0 {
		charactersPath := fmt.Sprintf("%s/%s", scriptsPath, charactersFile)
		if err := ui.dialog.LoadCharacters(charactersPath); err != nil {
			return fmt.Errorf("failed to load characters: %w", err)
		}
	}

	// Load dialog script
	dialogPath := fmt.Sprintf("%s/%s", scriptsPath, dialogFile)
	if err := ui.dialog.LoadDialog(dialogPath); err != nil {
		return fmt.Errorf("failed to load dialog: %w", err)
	}

	// Start dialog
	ui.dialog.StartDialog(startNode)
	return nil
}

// HideDialog closes the dialog widget
func (ui *UIManager) HideDialog() {
	if ui.dialog != nil {
		ui.dialog.Hide()
	}
}

// IsDialogVisible returns true if dialog widget is visible
func (ui *UIManager) IsDialogVisible() bool {
	return ui.dialog != nil && ui.dialog.IsVisible()
}

// GetDialogVariable gets a dialog variable value
func (ui *UIManager) GetDialogVariable(name string) (interface{}, bool) {
	if ui.dialog != nil {
		return ui.dialog.GetVariable(name)
	}
	return nil, false
}

// SetDialogVariable sets a dialog variable value
func (ui *UIManager) SetDialogVariable(name string, value interface{}) {
	if ui.dialog != nil {
		ui.dialog.SetVariable(name, value)
	}
}

// ShowInventory creates and shows the inventory widget for the given entity
func (ui *UIManager) ShowInventory(entity *ecs.Entity) error {
	if entity == nil {
		return fmt.Errorf("entity is nil")
	}

	// Check if entity has inventory component
	if entity.Inventory() == nil {
		return fmt.Errorf("entity does not have an inventory component")
	}

	// Create inventory widget if it doesn't exist or entity changed
	inventoryX := (ScreenWidth - 600) / 2  // Center horizontally
	inventoryY := (ScreenHeight - 500) / 2 // Center vertically

	// Create new inventory widget
	ui.inventory = NewInventoryWidget(inventoryX, inventoryY, entity)
	ui.inventory.Visible = true

	return nil
}

// HideInventory closes the inventory widget
func (ui *UIManager) HideInventory() {
	if ui.inventory != nil {
		ui.inventory.Close()
	}
}

// IsInventoryVisible returns true if inventory widget is visible
func (ui *UIManager) IsInventoryVisible() bool {
	return ui.inventory != nil && ui.inventory.IsOpen()
}

// IsPopupVisible returns true if any popup is currently visible
func (ui *UIManager) IsPopupVisible() bool {
	selectionVisible := ui.popupSelection != nil && ui.popupSelection.IsVisible
	infoVisible := ui.popupInfo != nil && ui.popupInfo.IsVisible
	statsVisible := ui.characterStats != nil && ui.characterStats.IsVisible()
	equipmentVisible := ui.equipment != nil && ui.equipment.IsVisible()
	dialogVisible := ui.dialog != nil && ui.dialog.IsVisible()
	inventoryVisible := ui.inventory != nil && ui.inventory.IsOpen()
	return selectionVisible || infoVisible || statsVisible || equipmentVisible || dialogVisible || inventoryVisible
}
