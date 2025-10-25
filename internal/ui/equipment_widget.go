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
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// Equipment Widget Constants
const (
	// Layout dimensions
	EquipmentWidgetWidth          = 500 // Widget width in pixels
	EquipmentWidgetHeight         = 580 // Widget height in pixels
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
	HeadSlotY       = 40
	ChestSlotX      = 200
	ChestSlotY      = 120
	LegsSlotX       = 200
	LegsSlotY       = 200
	FeetSlotX       = 200
	FeetSlotY       = 280
	WeaponSlotX     = 120
	WeaponSlotY     = 120
	ShieldSlotX     = 280
	ShieldSlotY     = 120
	AccessorySlot1X = 120
	AccessorySlot1Y = 40
	AccessorySlot2X = 280
	AccessorySlot2Y = 40

	// Stat comparison panel
	StatComparisonPanelX      = 350
	StatComparisonPanelY      = 40
	StatComparisonPanelWidth  = 120
	StatComparisonPanelHeight = 320
	StatComparisonLineHeight  = 16
	StatComparisonArrowWidth  = 12
	StatComparisonValueWidth  = 40

	// Equipment details panel
	EquipmentDetailsPanelX      = 20
	EquipmentDetailsPanelY      = 420
	EquipmentDetailsPanelWidth  = 460
	EquipmentDetailsPanelHeight = 120
	EquipmentDetailsLineHeight  = 14
	EquipmentDetailsPadding     = 8
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

	// Equipment details panel colors
	EquipmentDetailsPanelR = 30
	EquipmentDetailsPanelG = 30
	EquipmentDetailsPanelB = 40
	EquipmentDetailsPanelA = 255

	EquipmentDetailsBorderR = 80
	EquipmentDetailsBorderG = 80
	EquipmentDetailsBorderB = 100
	EquipmentDetailsBorderA = 255

	EquipmentDetailsTextR = 220
	EquipmentDetailsTextG = 220
	EquipmentDetailsTextB = 220
	EquipmentDetailsTextA = 255
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
	EquipmentComp      *components.EquipmentComponent
	CharacterStats     *components.RPGStatsComponent
	AvailableEquipment []*components.Equipment // Mock equipment available to equip
	Entity             *ecs.Entity             // Player entity for inventory access

	// Navigation and selection
	SelectedSlot  components.EquipmentSlot
	SlotPositions map[components.EquipmentSlot]SlotPosition
	SlotOrder     []components.EquipmentSlot // For tab navigation

	// Interaction state
	HoveredEquipment *components.Equipment // Equipment being hovered for stat preview
}

// NewEquipmentWidget creates a new equipment widget
func NewEquipmentWidget(screenWidth, screenHeight int, equipmentComp *components.EquipmentComponent, stats *components.RPGStatsComponent, entity *ecs.Entity) *EquipmentWidget {
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
		Entity:         entity,
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

// SetAvailableEquipment sets the list of available equipment for equipping
func (ew *EquipmentWidget) SetAvailableEquipment(equipment []*components.Equipment) {
	ew.AvailableEquipment = equipment
}

// AddAvailableEquipment adds a single piece of equipment to the available list
func (ew *EquipmentWidget) AddAvailableEquipment(equipment *components.Equipment) {
	if ew.AvailableEquipment == nil {
		ew.AvailableEquipment = make([]*components.Equipment, 0)
	}
	ew.AvailableEquipment = append(ew.AvailableEquipment, equipment)
}

// Update handles input and updates the equipment widget state
// Returns InputResult indicating what input was consumed
func (ew *EquipmentWidget) Update() InputResult {
	result := NewInputResult()
	
	if !ew.Visible {
		return result
	}

	// Handle ESC key to close widget
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ew.Hide()
		result.EscConsumed = true
	}

	// Check if mouse is over the widget area
	mouseX, mouseY := ebiten.CursorPosition()
	isMouseOverWidget := mouseX >= ew.X && mouseX <= ew.X+ew.Width &&
	                    mouseY >= ew.Y && mouseY <= ew.Y+ew.Height

	if isMouseOverWidget {
		result.MouseConsumed = true
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

	return result
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

// toggleEquipment handles equipping/unequipping items
func (ew *EquipmentWidget) toggleEquipment() {
	if ew.EquipmentComp == nil || ew.CharacterStats == nil {
		return
	}

	// Check if there's already equipment in this slot
	currentEquipment := ew.EquipmentComp.GetEquipped(ew.SelectedSlot)

	if currentEquipment != nil {
		// Unequip the current item
		ew.unequipItem(ew.SelectedSlot)
	} else {
		// Try to equip the first compatible item from available equipment
		ew.equipFirstCompatibleItem(ew.SelectedSlot)
	}
}

// unequipItem removes equipment from the specified slot
func (ew *EquipmentWidget) unequipItem(slot components.EquipmentSlot) {
	if ew.EquipmentComp == nil {
		return
	}

	unequippedItem := ew.EquipmentComp.Unequip(slot)
	if unequippedItem != nil {
		// Return the unequipped item to player's inventory
		if ew.Entity != nil && ew.Entity.Inventory() != nil {
			// Create an inventory item for the unequipped equipment
			inventoryItem := &components.Item{
				ID:          unequippedItem.ID,
				Name:        unequippedItem.Name,
				Description: unequippedItem.Description,
				Type:        components.ItemTypeEquipment,
				Rarity:      components.ItemRarity(unequippedItem.Rarity),
				Value:       unequippedItem.Value,
				IconID:      unequippedItem.IconID,
				Equipment:   unequippedItem,
				Stackable:   false,
				MaxStack:    1,
			}

			// Try to add to inventory
			remaining := ew.Entity.Inventory().AddItem(inventoryItem, 1)
			if remaining > 0 {
				// If inventory is full, add back to available equipment as fallback
				ew.AvailableEquipment = append(ew.AvailableEquipment, unequippedItem)
			}
		} else {
			// Fallback: Add the item back to available equipment for re-equipping
			ew.AvailableEquipment = append(ew.AvailableEquipment, unequippedItem)
		}
	}
}

// equipFirstCompatibleItem tries to equip the first compatible item for the given slot
func (ew *EquipmentWidget) equipFirstCompatibleItem(slot components.EquipmentSlot) {
	if ew.EquipmentComp == nil || ew.CharacterStats == nil {
		return
	}

	// Find the first compatible item for this slot
	for i, equipment := range ew.AvailableEquipment {
		if equipment.Slot == slot && ew.canEquipItem(equipment) {
			// Equip the item
			ew.EquipmentComp.Equip(equipment)

			// Remove from available equipment (move to end and slice off)
			ew.AvailableEquipment[i] = ew.AvailableEquipment[len(ew.AvailableEquipment)-1]
			ew.AvailableEquipment = ew.AvailableEquipment[:len(ew.AvailableEquipment)-1]

			return
		}
	}
}

// canEquipItem checks if the character can equip the given item
func (ew *EquipmentWidget) canEquipItem(equipment *components.Equipment) bool {
	if ew.CharacterStats == nil {
		return false
	}

	// Check level requirement
	if ew.CharacterStats.Level < equipment.LevelRequirement {
		return false
	}

	// Check job restrictions
	return equipment.CanEquip(ew.CharacterStats.Level, ew.CharacterStats.Job)
}

// Draw renders the equipment widget to the screen
func (ew *EquipmentWidget) Draw(screen *ebiten.Image) {
	if !ew.Visible {
		return
	}

	// Draw shadow
	shadowColor := color.RGBA{0, 0, 0, 100}
	vector.FillRect(screen,
		float32(ew.X+EquipmentWidgetShadowOffset), float32(ew.Y+EquipmentWidgetShadowOffset),
		float32(ew.Width), float32(ew.Height),
		shadowColor, false)

	// Draw main widget background
	bgColor := color.RGBA{EquipmentWidgetBackgroundR, EquipmentWidgetBackgroundG, EquipmentWidgetBackgroundB, EquipmentWidgetBackgroundA}
	vector.FillRect(screen,
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

	// Draw equipment details panel
	ew.drawEquipmentDetails(screen)

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

		vector.FillRect(screen,
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
	vector.FillRect(screen,
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

// drawEquipmentDetails renders detailed information about the currently selected equipment
func (ew *EquipmentWidget) drawEquipmentDetails(screen *ebiten.Image) {
	panelX := ew.X + EquipmentDetailsPanelX
	panelY := ew.Y + EquipmentDetailsPanelY

	// Draw panel background
	panelBg := color.RGBA{EquipmentDetailsPanelR, EquipmentDetailsPanelG, EquipmentDetailsPanelB, EquipmentDetailsPanelA}
	vector.FillRect(screen,
		float32(panelX), float32(panelY),
		float32(EquipmentDetailsPanelWidth), float32(EquipmentDetailsPanelHeight),
		panelBg, false)

	// Draw panel border
	panelBorder := color.RGBA{EquipmentDetailsBorderR, EquipmentDetailsBorderG, EquipmentDetailsBorderB, EquipmentDetailsBorderA}
	vector.StrokeRect(screen,
		float32(panelX), float32(panelY),
		float32(EquipmentDetailsPanelWidth), float32(EquipmentDetailsPanelHeight),
		1, panelBorder, false)

	// Get the equipment in the currently selected slot
	var selectedEquipment *components.Equipment
	if ew.EquipmentComp != nil {
		selectedEquipment = ew.EquipmentComp.GetEquipped(ew.SelectedSlot)
	}

	textX := panelX + EquipmentDetailsPadding
	textY := panelY + EquipmentDetailsPadding

	if selectedEquipment != nil {
		// Show detailed information about the equipped item
		ew.drawEquipmentInfo(screen, selectedEquipment, textX, textY)
	} else {
		// Show information about the selected slot
		slotName := ew.SelectedSlot.String()
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%s Slot]", slotName), textX, textY)
		textY += EquipmentDetailsLineHeight + 2

		// Find the first compatible item that could be equipped in this slot
		var compatibleItem *components.Equipment
		for _, equipment := range ew.AvailableEquipment {
			if equipment.Slot == ew.SelectedSlot && ew.canEquipItem(equipment) {
				compatibleItem = equipment
				break
			}
		}

		if compatibleItem != nil {
			ebitenutil.DebugPrintAt(screen, "Press ENTER to equip:", textX, textY)
			textY += EquipmentDetailsLineHeight + 2
			ew.drawEquipmentInfo(screen, compatibleItem, textX, textY)
		} else {
			ebitenutil.DebugPrintAt(screen, "No compatible equipment available", textX, textY)
		}
	}
}

// drawEquipmentInfo renders detailed information about a specific piece of equipment
func (ew *EquipmentWidget) drawEquipmentInfo(screen *ebiten.Image, equipment *components.Equipment, startX, startY int) {
	textY := startY

	// Equipment name
	ebitenutil.DebugPrintAt(screen, equipment.Name, startX, textY)

	// Rarity (on same line, offset to the right)
	rarityText := fmt.Sprintf("(%s)", equipment.Rarity.String())
	ebitenutil.DebugPrintAt(screen, rarityText, startX+200, textY)
	textY += EquipmentDetailsLineHeight + 2

	// Description
	if equipment.Description != "" {
		ebitenutil.DebugPrintAt(screen, equipment.Description, startX, textY)
		textY += EquipmentDetailsLineHeight + 2
	}

	// Requirements
	if equipment.LevelRequirement > 1 || len(equipment.JobRestrictions) > 0 {
		reqText := "Req: "
		if equipment.LevelRequirement > 1 {
			reqText += fmt.Sprintf("Lv %d", equipment.LevelRequirement)
		}
		if len(equipment.JobRestrictions) > 0 {
			if equipment.LevelRequirement > 1 {
				reqText += ", "
			}
			jobNames := make([]string, len(equipment.JobRestrictions))
			for i, job := range equipment.JobRestrictions {
				jobNames[i] = job.String()
			}
			reqText += fmt.Sprintf("Jobs: %v", jobNames)
		}
		ebitenutil.DebugPrintAt(screen, reqText, startX, textY)
		textY += EquipmentDetailsLineHeight + 2
	}

	// Stats (show only non-zero stats)
	stats := equipment.Stats
	statsText := "Stats: "
	statCount := 0

	if stats.AttackBonus != 0 {
		if statCount > 0 {
			statsText += ", "
		}
		statsText += fmt.Sprintf("ATK %+d", stats.AttackBonus)
		statCount++
	}
	if stats.DefenseBonus != 0 {
		if statCount > 0 {
			statsText += ", "
		}
		statsText += fmt.Sprintf("DEF %+d", stats.DefenseBonus)
		statCount++
	}
	if stats.MagicPowerBonus != 0 {
		if statCount > 0 {
			statsText += ", "
		}
		statsText += fmt.Sprintf("MAG %+d", stats.MagicPowerBonus)
		statCount++
	}
	if stats.SpeedBonus != 0 {
		if statCount > 0 {
			statsText += ", "
		}
		statsText += fmt.Sprintf("SPD %+d", stats.SpeedBonus)
		statCount++
	}
	if stats.HPBonus != 0 {
		if statCount > 0 {
			statsText += ", "
		}
		statsText += fmt.Sprintf("HP %+d", stats.HPBonus)
		statCount++
	}
	if stats.MPBonus != 0 {
		if statCount > 0 {
			statsText += ", "
		}
		statsText += fmt.Sprintf("MP %+d", stats.MPBonus)
		statCount++
	}

	if statCount > 0 {
		ebitenutil.DebugPrintAt(screen, statsText, startX, textY)
	} else {
		ebitenutil.DebugPrintAt(screen, "Stats: None", startX, textY)
	}
}

// drawHelpText renders the help text at the bottom of the widget
func (ew *EquipmentWidget) drawHelpText(screen *ebiten.Image) {
	baseHelp := "TAB: Next Slot  ESC: Close  Arrows: Navigate"

	// Add contextual ENTER action
	if ew.EquipmentComp != nil {
		currentEquipment := ew.EquipmentComp.GetEquipped(ew.SelectedSlot)
		if currentEquipment != nil {
			baseHelp += "  ENTER: Unequip"
		} else {
			baseHelp += "  ENTER: Equip"
		}
	} else {
		baseHelp += "  ENTER: Equip/Unequip"
	}

	helpY := ew.Y + ew.Height - EquipmentWidgetBottomReserved + 10
	ebitenutil.DebugPrintAt(screen, baseHelp, ew.X+EquipmentWidgetPadding, helpY)
}

// Helper function for min since Go doesn't have built-in min for int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CreateMockEquipmentSet creates a basic set of equipment for testing
func CreateMockEquipmentSet() []*components.Equipment {
	return []*components.Equipment{
		// Weapons
		{
			ID:          1,
			Name:        "Iron Sword",
			Description: "A basic iron sword.",
			Slot:        components.SlotWeapon,
			Rarity:      components.RarityCommon,
			Stats: components.EquipmentStats{
				AttackBonus:     10,
				CritChanceBonus: 5,
			},
			LevelRequirement: 1,
			JobRestrictions:  []components.JobType{components.JobWarrior},
			Value:            50,
		},
		{
			ID:          2,
			Name:        "Magic Staff",
			Description: "A staff crackling with magical energy.",
			Slot:        components.SlotWeapon,
			Rarity:      components.RarityUncommon,
			Stats: components.EquipmentStats{
				MagicPowerBonus: 15,
				MPBonus:         20,
			},
			LevelRequirement: 3,
			JobRestrictions:  []components.JobType{components.JobMage, components.JobCleric},
			Value:            120,
		},

		// Armor
		{
			ID:          3,
			Name:        "Leather Vest",
			Description: "Light leather armor.",
			Slot:        components.SlotChest,
			Rarity:      components.RarityCommon,
			Stats: components.EquipmentStats{
				DefenseBonus: 5,
				HPBonus:      15,
			},
			LevelRequirement: 1,
			JobRestrictions:  []components.JobType{},
			Value:            30,
		},
		{
			ID:          4,
			Name:        "Iron Helm",
			Description: "A sturdy iron helmet.",
			Slot:        components.SlotHead,
			Rarity:      components.RarityCommon,
			Stats: components.EquipmentStats{
				DefenseBonus:  3,
				MagicDefBonus: 2,
			},
			LevelRequirement: 2,
			JobRestrictions:  []components.JobType{components.JobWarrior},
			Value:            40,
		},

		// Accessories
		{
			ID:          5,
			Name:        "Ring of Strength",
			Description: "Increases physical power.",
			Slot:        components.SlotAccessory1,
			Rarity:      components.RarityRare,
			Stats: components.EquipmentStats{
				AttackBonus: 8,
				HPBonus:     10,
			},
			LevelRequirement: 5,
			JobRestrictions:  []components.JobType{},
			Value:            200,
		},
		{
			ID:          6,
			Name:        "Swift Boots",
			Description: "Boots that enhance movement speed.",
			Slot:        components.SlotFeet,
			Rarity:      components.RarityUncommon,
			Stats: components.EquipmentStats{
				SpeedBonus:    10,
				MovementBonus: 1,
			},
			LevelRequirement: 3,
			JobRestrictions:  []components.JobType{},
			Value:            80,
		},
	}
}
