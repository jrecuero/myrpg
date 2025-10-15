// Package ui provides equipment management widget for the RPG game.
// This widget shows character equipment in a paperdoll-style layout
// with drag-and-drop functionality and stat comparisons.
package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// Equipment Widget Constants
const (
	// Layout dimensions
	EquipmentWidgetWidth          = 500 // Widget width in pixels
	EquipmentWidgetHeight         = 600 // Widget height in pixels
	EquipmentWidgetBorderWidth    = 2   // Border thickness
	EquipmentWidgetShadowOffset   = 4   // Shadow offset distance
	EquipmentWidgetPadding        = 15  // Internal padding
	EquipmentWidgetTitleY         = 20  // Y offset for main title
	EquipmentWidgetContentStartY  = 50  // Y offset where content begins
	EquipmentWidgetBottomReserved = 40  // Space reserved at bottom for help text

	// Equipment slot layout
	EquipmentSlotSize         = 64 // Size of each equipment slot (square)
	EquipmentSlotBorder       = 2  // Border width for slots
	EquipmentSlotSpacing      = 10 // Space between adjacent slots
	EquipmentSlotCornerRadius = 6  // Corner radius for slot borders
	EquipmentIconSize         = 48 // Size of equipment icons within slots
	EquipmentIconPadding      = 8  // Padding around icons in slots

	// Paperdoll layout positions (relative to content area)
	HeadSlotX       = 200
	HeadSlotY       = 80
	ChestSlotX      = 200
	ChestSlotY      = 160
	LegsSlotX       = 200
	LegsSlotY       = 240
	FeetSlotX       = 200
	FeetSlotY       = 320
	WeaponSlotX     = 120
	WeaponSlotY     = 160
	ShieldSlotX     = 280
	ShieldSlotY     = 160
	AccessorySlot1X = 120
	AccessorySlot1Y = 80
	AccessorySlot2X = 280
	AccessorySlot2Y = 80

	// Stat comparison panel
	StatComparisonPanelX      = 350
	StatComparisonPanelY      = 80
	StatComparisonPanelWidth  = 120
	StatComparisonPanelHeight = 320
	StatComparisonLineHeight  = 16
	StatComparisonArrowWidth  = 12
	StatComparisonValueWidth  = 40
)

// Color constants - RGBA values for easy customization
const (
	// Background and border colors
	EquipmentWidgetBackgroundR = 25
	EquipmentWidgetBackgroundG = 25
	EquipmentWidgetBackgroundB = 35
	EquipmentWidgetBackgroundA = 255

	EquipmentWidgetBorderR = 100
	EquipmentWidgetBorderG = 100
	EquipmentWidgetBorderB = 120
	EquipmentWidgetBorderA = 255

	// Text colors
	EquipmentWidgetTitleR = 255
	EquipmentWidgetTitleG = 255
	EquipmentWidgetTitleB = 255
	EquipmentWidgetTitleA = 255

	EquipmentWidgetHelpR = 180
	EquipmentWidgetHelpG = 180
	EquipmentWidgetHelpB = 180
	EquipmentWidgetHelpA = 255

	// Equipment slot colors
	EquipmentSlotEmptyR = 40
	EquipmentSlotEmptyG = 40
	EquipmentSlotEmptyB = 50
	EquipmentSlotEmptyA = 255

	EquipmentSlotBorderR = 100
	EquipmentSlotBorderG = 100
	EquipmentSlotBorderB = 120
	EquipmentSlotBorderA = 255

	EquipmentSlotEquippedR = 60
	EquipmentSlotEquippedG = 80
	EquipmentSlotEquippedB = 100
	EquipmentSlotEquippedA = 255

	// Stat comparison colors
	StatIncreaseR = 100
	StatIncreaseG = 255
	StatIncreaseB = 100
	StatIncreaseA = 255

	StatDecreaseR = 255
	StatDecreaseG = 100
	StatDecreaseB = 100
	StatDecreaseA = 255

	StatNoChangeR = 200
	StatNoChangeG = 200
	StatNoChangeB = 200
	StatNoChangeA = 255
)

// Equipment rarity colors
var rarityColors = map[components.EquipmentRarity]color.RGBA{
	components.RarityCommon:    {200, 200, 200, 255},
	components.RarityUncommon:  {100, 255, 100, 255},
	components.RarityRare:      {100, 100, 255, 255},
	components.RarityEpic:      {200, 100, 255, 255},
	components.RarityLegendary: {255, 200, 50, 255},
}

// SlotPosition represents the position of an equipment slot
type SlotPosition struct {
	X, Y int
}

// getSlotPositions returns the positions for all equipment slots
func getSlotPositions() map[components.EquipmentSlot]SlotPosition {
	return map[components.EquipmentSlot]SlotPosition{
		components.SlotHead:       {HeadSlotX, HeadSlotY},
		components.SlotChest:      {ChestSlotX, ChestSlotY},
		components.SlotLegs:       {LegsSlotX, LegsSlotY},
		components.SlotFeet:       {FeetSlotX, FeetSlotY},
		components.SlotWeapon:     {WeaponSlotX, WeaponSlotY},
		components.SlotShield:     {ShieldSlotX, ShieldSlotY},
		components.SlotAccessory1: {AccessorySlot1X, AccessorySlot1Y},
		components.SlotAccessory2: {AccessorySlot2X, AccessorySlot2Y},
	}
}

// EquipmentWidget represents the equipment management interface
type EquipmentWidget struct {
	// Widget properties
	X, Y          int
	Width, Height int
	Visible       bool

	// Equipment data
	EquipmentComp  *components.EquipmentComponent
	CharacterStats *components.RPGStatsComponent

	// Navigation and selection
	SelectedSlot  components.EquipmentSlot
	SlotPositions map[components.EquipmentSlot]SlotPosition
	SlotOrder     []components.EquipmentSlot // For tab navigation

	// Interaction state
	HoveredEquipment *components.Equipment // Equipment being hovered for stat preview
}

// NewEquipmentWidget creates a new equipment widget
func NewEquipmentWidget(screenWidth, screenHeight int, equipmentComp *components.EquipmentComponent, stats *components.RPGStatsComponent) *EquipmentWidget {
	x := (screenWidth - EquipmentWidgetWidth) / 2
	y := (screenHeight - EquipmentWidgetHeight) / 2

	widget := &EquipmentWidget{
		X:              x,
		Y:              y,
		Width:          EquipmentWidgetWidth,
		Height:         EquipmentWidgetHeight,
		Visible:        false,
		EquipmentComp:  equipmentComp,
		CharacterStats: stats,
		SelectedSlot:   components.SlotHead,
		SlotPositions:  getSlotPositions(),
		SlotOrder: []components.EquipmentSlot{
			components.SlotHead,
			components.SlotAccessory1,
			components.SlotAccessory2,
			components.SlotWeapon,
			components.SlotChest,
			components.SlotShield,
			components.SlotLegs,
			components.SlotFeet,
		},
	}

	return widget
}

// Show makes the equipment widget visible
func (ew *EquipmentWidget) Show() {
	ew.Visible = true
}

// Hide makes the equipment widget invisible
func (ew *EquipmentWidget) Hide() {
	ew.Visible = false
}

// IsVisible returns whether the widget is visible
func (ew *EquipmentWidget) IsVisible() bool {
	return ew.Visible
}

// Update handles input and updates the equipment widget state
// Returns true if ESC key was consumed by this widget
func (ew *EquipmentWidget) Update() bool {
	if !ew.Visible {
		return false
	}

	escConsumed := false

	// Handle ESC key to close widget
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ew.Hide()
		escConsumed = true
	}

	// Handle TAB key for slot navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		ew.nextSlot()
	}

	// Handle arrow keys for slot navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		ew.navigateSlot(-1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		ew.navigateSlot(1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		ew.navigateSlot(0, -1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		ew.navigateSlot(0, 1)
	}

	// Handle ENTER key for equip/unequip (placeholder for now)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ew.toggleEquipment()
	}

	return escConsumed
}

// nextSlot advances to the next equipment slot in tab order
func (ew *EquipmentWidget) nextSlot() {
	currentIndex := -1
	for i, slot := range ew.SlotOrder {
		if slot == ew.SelectedSlot {
			currentIndex = i
			break
		}
	}

	if currentIndex != -1 {
		nextIndex := (currentIndex + 1) % len(ew.SlotOrder)
		ew.SelectedSlot = ew.SlotOrder[nextIndex]
	}
}

// navigateSlot moves selection based on spatial layout using arrow keys
func (ew *EquipmentWidget) navigateSlot(deltaRow, deltaCol int) {
	// For now, use simple tab-style navigation
	// TODO: Implement proper spatial navigation based on slot positions
	if deltaRow != 0 || deltaCol != 0 {
		ew.nextSlot()
	}
}

// toggleEquipment handles equipping/unequipping items (placeholder)
func (ew *EquipmentWidget) toggleEquipment() {
	// TODO: Implement equipment toggle logic
	// For now, this is a placeholder that could:
	// - Open inventory selection if slot is empty
	// - Unequip item if slot is occupied
	// - Show confirmation dialog for valuable items
}

// Draw renders the equipment widget to the screen
func (ew *EquipmentWidget) Draw(screen *ebiten.Image) {
	if !ew.Visible {
		return
	}

	// Draw shadow
	shadowColor := color.RGBA{0, 0, 0, 100}
	vector.DrawFilledRect(screen,
		float32(ew.X+EquipmentWidgetShadowOffset), float32(ew.Y+EquipmentWidgetShadowOffset),
		float32(ew.Width), float32(ew.Height),
		shadowColor, false)

	// Draw main widget background
	bgColor := color.RGBA{EquipmentWidgetBackgroundR, EquipmentWidgetBackgroundG, EquipmentWidgetBackgroundB, EquipmentWidgetBackgroundA}
	vector.DrawFilledRect(screen,
		float32(ew.X), float32(ew.Y),
		float32(ew.Width), float32(ew.Height),
		bgColor, false)

	// Draw border
	borderColor := color.RGBA{EquipmentWidgetBorderR, EquipmentWidgetBorderG, EquipmentWidgetBorderB, EquipmentWidgetBorderA}
	for i := 0; i < EquipmentWidgetBorderWidth; i++ {
		vector.StrokeRect(screen,
			float32(ew.X+i), float32(ew.Y+i),
			float32(ew.Width-i*2), float32(ew.Height-i*2),
			1, borderColor, false)
	}

	// Draw title
	title := fmt.Sprintf("Equipment - %s", ew.CharacterStats.Job.String())
	ebitenutil.DebugPrintAt(screen, title, ew.X+EquipmentWidgetPadding, ew.Y+EquipmentWidgetTitleY)

	// Draw equipment slots
	ew.drawEquipmentSlots(screen)

	// Draw stat comparison panel (if hovering over equipment)
	if ew.HoveredEquipment != nil {
		ew.drawStatComparison(screen)
	}

	// Draw help text
	ew.drawHelpText(screen)
}

// drawEquipmentSlots renders all equipment slots in paperdoll layout
func (ew *EquipmentWidget) drawEquipmentSlots(screen *ebiten.Image) {
	contentX := ew.X + EquipmentWidgetPadding
	contentY := ew.Y + EquipmentWidgetContentStartY

	for slot, pos := range ew.SlotPositions {
		slotX := contentX + pos.X
		slotY := contentY + pos.Y

		// Check if this slot is selected
		isSelected := slot == ew.SelectedSlot

		// Check if equipment is equipped in this slot
		equipment := ew.EquipmentComp.GetEquipped(slot)
		isEquipped := equipment != nil

		// Draw slot background
		var slotBgColor color.RGBA
		if isEquipped {
			slotBgColor = color.RGBA{EquipmentSlotEquippedR, EquipmentSlotEquippedG, EquipmentSlotEquippedB, EquipmentSlotEquippedA}
		} else {
			slotBgColor = color.RGBA{EquipmentSlotEmptyR, EquipmentSlotEmptyG, EquipmentSlotEmptyB, EquipmentSlotEmptyA}
		}

		vector.DrawFilledRect(screen,
			float32(slotX), float32(slotY),
			float32(EquipmentSlotSize), float32(EquipmentSlotSize),
			slotBgColor, false)

		// Draw slot border (with rarity color if equipped)
		var borderColor color.RGBA
		if isEquipped {
			borderColor = rarityColors[equipment.Rarity]
		} else {
			borderColor = color.RGBA{EquipmentSlotBorderR, EquipmentSlotBorderG, EquipmentSlotBorderB, EquipmentSlotBorderA}
		}

		// Draw thicker border if selected
		borderWidth := EquipmentSlotBorder
		if isSelected {
			borderWidth *= 2
		}

		for i := 0; i < borderWidth; i++ {
			vector.StrokeRect(screen,
				float32(slotX+i), float32(slotY+i),
				float32(EquipmentSlotSize-i*2), float32(EquipmentSlotSize-i*2),
				1, borderColor, false)
		}

		// Draw equipment icon or slot label
		if isEquipped {
			// TODO: Draw equipment icon/sprite
			// For now, draw equipment name abbreviated
			iconX := slotX + EquipmentIconPadding
			iconY := slotY + EquipmentSlotSize/2
			ebitenutil.DebugPrintAt(screen, equipment.Name[:min(3, len(equipment.Name))], iconX, iconY)
		} else {
			// Draw slot type label
			labelX := slotX + EquipmentIconPadding
			labelY := slotY + EquipmentSlotSize/2
			slotName := slot.String()
			ebitenutil.DebugPrintAt(screen, slotName[:min(4, len(slotName))], labelX, labelY)
		}
	}
}

// drawStatComparison renders the stat comparison panel for hovered equipment
func (ew *EquipmentWidget) drawStatComparison(screen *ebiten.Image) {
	panelX := ew.X + StatComparisonPanelX
	panelY := ew.Y + StatComparisonPanelY

	// Draw panel background
	panelBg := color.RGBA{40, 40, 50, 200}
	vector.DrawFilledRect(screen,
		float32(panelX), float32(panelY),
		float32(StatComparisonPanelWidth), float32(StatComparisonPanelHeight),
		panelBg, false)

	// Draw panel border
	panelBorder := color.RGBA{100, 100, 120, 255}
	vector.StrokeRect(screen,
		float32(panelX), float32(panelY),
		float32(StatComparisonPanelWidth), float32(StatComparisonPanelHeight),
		1, panelBorder, false)

	// TODO: Implement detailed stat comparison display
	// For now, show equipment name and rarity
	textY := panelY + 10
	ebitenutil.DebugPrintAt(screen, ew.HoveredEquipment.Name, panelX+5, textY)
	textY += StatComparisonLineHeight
	ebitenutil.DebugPrintAt(screen, ew.HoveredEquipment.Rarity.String(), panelX+5, textY)
}

// drawHelpText renders the help text at the bottom of the widget
func (ew *EquipmentWidget) drawHelpText(screen *ebiten.Image) {
	helpText := "TAB: Next Slot  ESC: Close  ENTER: Equip/Unequip  Arrows: Navigate"

	helpY := ew.Y + ew.Height - EquipmentWidgetBottomReserved + 10
	ebitenutil.DebugPrintAt(screen, helpText, ew.X+EquipmentWidgetPadding, helpY)
}

// Helper function for min since Go doesn't have built-in min for int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
