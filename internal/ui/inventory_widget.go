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

// Constants for visual styling
const (
	// Widget dimensions
	InventoryWidgetWidth       = 600
	InventoryWidgetHeight      = 500
	InventoryWidgetPadding     = 20
	InventoryWidgetBorderWidth = 3

	// Grid layout
	InventoryGridWidth   = 8
	InventoryGridHeight  = 6
	InventorySlotSize    = 48
	InventorySlotSpacing = 4
	InventoryGridStartX  = 40
	InventoryGridStartY  = 80

	// Header section
	InventoryHeaderHeight = 50
	InventoryTitleX       = 25
	InventoryTitleY       = 25
	InventoryStatsX       = 400
	InventoryStatsY       = 25

	// Tooltip system
	InventoryTooltipWidth       = 250
	InventoryTooltipMaxHeight   = 200
	InventoryTooltipPadding     = 8
	InventoryTooltipBorderWidth = 2
	InventoryTooltipOffsetX     = 10
	InventoryTooltipOffsetY     = -10

	// Action panel
	InventoryActionPanelWidth    = 140
	InventoryActionPanelHeight   = 300
	InventoryActionPanelX        = 440
	InventoryActionPanelY        = 80
	InventoryActionButtonHeight  = 30
	InventoryActionButtonSpacing = 5

	// Animation
	InventoryDragOpacity = 128
)

// Color definitions
var (
	// Widget colors
	InventoryWidgetBackground = color.RGBA{45, 45, 55, 220}
	InventoryWidgetBorder     = color.RGBA{120, 120, 140, 255}

	// Slot colors
	InventorySlotEmpty          = color.RGBA{60, 60, 70, 180}
	InventorySlotFilled         = color.RGBA{80, 80, 95, 200}
	InventorySlotSelected       = color.RGBA{100, 150, 200, 230}
	InventorySlotHover          = color.RGBA{90, 130, 180, 200}
	InventorySlotBorder         = color.RGBA{100, 100, 120, 255}
	InventorySlotSelectedBorder = color.RGBA{150, 200, 255, 255}

	// Rarity colors
	RarityColors = map[components.ItemRarity]color.RGBA{
		components.ItemRarityCommon:    {200, 200, 200, 255},
		components.ItemRarityUncommon:  {30, 255, 0, 255},
		components.ItemRarityRare:      {0, 112, 255, 255},
		components.ItemRarityEpic:      {163, 53, 238, 255},
		components.ItemRarityLegendary: {255, 128, 0, 255},
	}

	// Tooltip colors
	InventoryTooltipBackground  = color.RGBA{25, 25, 35, 240}
	InventoryTooltipBorderColor = color.RGBA{150, 150, 170, 255}

	// Action panel colors
	InventoryActionPanelBg     = color.RGBA{35, 35, 45, 200}
	InventoryActionButton      = color.RGBA{70, 80, 90, 255}
	InventoryActionButtonHover = color.RGBA{90, 100, 120, 255}

	// Shadow
	InventoryWidgetShadow = color.RGBA{0, 0, 0, 100}
)

// InventorySlot represents a visual slot in the inventory grid
type InventorySlot struct {
	X, Y          int
	Width, Height int
	Item          *components.Item
	Quantity      int
	IsHovered     bool
	IsSelected    bool
}

// InventoryWidget manages the inventory display and interaction
type InventoryWidget struct {
	// Widget properties
	X, Y          int
	Width, Height int
	Visible       bool

	// Entity and inventory data
	entity         *ecs.Entity
	inventory      *components.InventoryComponent
	slots          [][]InventorySlot
	selectedSlot   *InventorySlot
	draggedItem    *components.Item
	draggedSlot    *InventorySlot
	isDragging     bool
	mouseX, mouseY int
	showTooltip    bool
	tooltipItem    *components.Item
	tooltipX       int
	tooltipY       int
	sortMode       string
	filterType     components.ItemType
	showAllTypes   bool
}

// NewInventoryWidget creates a new inventory widget
func NewInventoryWidget(x, y int, ent *ecs.Entity) *InventoryWidget {
	iw := &InventoryWidget{
		X:            x,
		Y:            y,
		Width:        InventoryWidgetWidth,
		Height:       InventoryWidgetHeight,
		Visible:      true,
		entity:       ent,
		sortMode:     "name",
		showAllTypes: true,
	}

	// Get inventory component
	if inv := ent.Inventory(); inv != nil {
		iw.inventory = inv
		iw.initializeSlots()
	}

	return iw
}

// initializeSlots creates the visual slot grid
func (iw *InventoryWidget) initializeSlots() {
	iw.slots = make([][]InventorySlot, InventoryGridHeight)

	for row := 0; row < InventoryGridHeight; row++ {
		iw.slots[row] = make([]InventorySlot, InventoryGridWidth)

		for col := 0; col < InventoryGridWidth; col++ {
			slot := &iw.slots[row][col]
			slot.X = iw.X + InventoryGridStartX + col*(InventorySlotSize+InventorySlotSpacing)
			slot.Y = iw.Y + InventoryGridStartY + row*(InventorySlotSize+InventorySlotSpacing)
			slot.Width = InventorySlotSize
			slot.Height = InventorySlotSize

			// Map to inventory data
			slotIndex := row*InventoryGridWidth + col
			if slotIndex < len(iw.inventory.Slots) {
				invSlot := &iw.inventory.Slots[slotIndex]
				slot.Item = invSlot.Item
				slot.Quantity = invSlot.Quantity
			}
		}
	}
}

// Update handles input and state changes
// Returns InputResult indicating what input was consumed
func (iw *InventoryWidget) Update() InputResult {
	result := NewInputResult()

	if !iw.Visible {
		return result
	}

	// Update mouse position
	iw.mouseX, iw.mouseY = ebiten.CursorPosition()

	// Handle escape key to close
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		iw.Visible = false
		result.EscConsumed = true
		return result
	}

	// Check if mouse is over the widget area
	isMouseOverWidget := iw.mouseX >= iw.X && iw.mouseX <= iw.X+iw.Width &&
		iw.mouseY >= iw.Y && iw.mouseY <= iw.Y+iw.Height

	if isMouseOverWidget {
		result.MouseConsumed = true
	}

	// Update hover states
	iw.updateHoverStates()

	// Handle mouse interactions
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		iw.handleLeftClick()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		iw.handleRightClick()
	}

	// Handle drag and drop
	if iw.isDragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// Continue dragging
		} else {
			// Drop item
			iw.handleDrop()
		}
	}

	// Handle keyboard shortcuts
	iw.handleKeyboardInput()

	return result
}

// updateHoverStates updates which slots are being hovered
func (iw *InventoryWidget) updateHoverStates() {
	iw.showTooltip = false

	for row := 0; row < InventoryGridHeight; row++ {
		for col := 0; col < InventoryGridWidth; col++ {
			slot := &iw.slots[row][col]

			// Check if mouse is over this slot
			if iw.mouseX >= slot.X && iw.mouseX < slot.X+slot.Width &&
				iw.mouseY >= slot.Y && iw.mouseY < slot.Y+slot.Height {
				slot.IsHovered = true

				// Show tooltip if slot has an item
				if slot.Item != nil {
					iw.showTooltip = true
					iw.tooltipItem = slot.Item
					iw.tooltipX = iw.mouseX + InventoryTooltipOffsetX
					iw.tooltipY = iw.mouseY + InventoryTooltipOffsetY
				}
			} else {
				slot.IsHovered = false
			}
		}
	}
}

// handleLeftClick processes left mouse clicks
func (iw *InventoryWidget) handleLeftClick() {
	// Find clicked slot
	for row := 0; row < InventoryGridHeight; row++ {
		for col := 0; col < InventoryGridWidth; col++ {
			slot := &iw.slots[row][col]

			if slot.IsHovered {
				if slot.Item != nil {
					// Start dragging if item exists
					iw.startDrag(slot)
				} else if iw.selectedSlot != nil {
					// Clear selection if clicking empty slot
					iw.selectedSlot.IsSelected = false
					iw.selectedSlot = nil
				}
				return
			}
		}
	}

	// Check action panel clicks
	iw.handleActionPanelClick()
}

// handleRightClick processes right mouse clicks
func (iw *InventoryWidget) handleRightClick() {
	// Find right-clicked slot
	for row := 0; row < InventoryGridHeight; row++ {
		for col := 0; col < InventoryGridWidth; col++ {
			slot := &iw.slots[row][col]

			if slot.IsHovered && slot.Item != nil {
				iw.showItemContextMenu(slot)
				return
			}
		}
	}
}

// startDrag begins dragging an item
func (iw *InventoryWidget) startDrag(slot *InventorySlot) {
	iw.isDragging = true
	iw.draggedItem = slot.Item
	iw.draggedSlot = slot

	// Select the slot
	if iw.selectedSlot != nil {
		iw.selectedSlot.IsSelected = false
	}
	slot.IsSelected = true
	iw.selectedSlot = slot
}

// handleDrop processes dropping an item
func (iw *InventoryWidget) handleDrop() {
	defer func() {
		iw.isDragging = false
		iw.draggedItem = nil
		iw.draggedSlot = nil
	}()

	// Find drop target slot
	for row := 0; row < InventoryGridHeight; row++ {
		for col := 0; col < InventoryGridWidth; col++ {
			slot := &iw.slots[row][col]

			if slot.IsHovered && slot != iw.draggedSlot {
				iw.performItemSwap(iw.draggedSlot, slot)
				return
			}
		}
	}

	// If dropped outside inventory, item returns to original position
	// (Item stays in draggedSlot)
}

// performItemSwap swaps items between two slots
func (iw *InventoryWidget) performItemSwap(slot1, slot2 *InventorySlot) {
	// Calculate inventory indices
	slot1Index := iw.getSlotIndex(slot1)
	slot2Index := iw.getSlotIndex(slot2)

	if slot1Index == -1 || slot2Index == -1 {
		return
	}

	// Perform swap manually using slot operations
	slot1Data := iw.inventory.GetSlotByIndex(slot1Index)
	slot2Data := iw.inventory.GetSlotByIndex(slot2Index)

	if slot1Data != nil && slot2Data != nil {
		// Store slot 1 data temporarily
		tempItem := slot1Data.Item
		tempQuantity := slot1Data.Quantity

		// Clear slot 1 and set to slot 2 data
		slot1Data.Item = slot2Data.Item
		slot1Data.Quantity = slot2Data.Quantity

		// Set slot 2 to original slot 1 data
		slot2Data.Item = tempItem
		slot2Data.Quantity = tempQuantity

		// Update visual slots
		slot1.Item, slot2.Item = slot2.Item, slot1.Item
		slot1.Quantity, slot2.Quantity = slot2.Quantity, slot1.Quantity
	}
}

// getSlotIndex returns the inventory index for a visual slot
func (iw *InventoryWidget) getSlotIndex(slot *InventorySlot) int {
	for row := 0; row < InventoryGridHeight; row++ {
		for col := 0; col < InventoryGridWidth; col++ {
			if &iw.slots[row][col] == slot {
				return row*InventoryGridWidth + col
			}
		}
	}
	return -1
}

// showItemContextMenu shows context menu for an item
func (iw *InventoryWidget) showItemContextMenu(slot *InventorySlot) {
	if slot.Item == nil {
		return
	}

	// Select the slot first
	if iw.selectedSlot != nil {
		iw.selectedSlot.IsSelected = false
	}
	slot.IsSelected = true
	iw.selectedSlot = slot

	// For equipment items, try to equip them directly
	if slot.Item.Type == components.ItemTypeEquipment && slot.Item.Equipment != nil {
		iw.tryEquipItem(slot)
	}
}

// tryEquipItem attempts to equip an item from inventory
func (iw *InventoryWidget) tryEquipItem(slot *InventorySlot) {
	if slot.Item == nil || slot.Item.Equipment == nil {
		return
	}

	// Get the player's equipment component
	equipmentComp := iw.entity.Equipment()
	if equipmentComp == nil {
		// Create equipment component if it doesn't exist
		equipmentComp = components.NewEquipmentComponent()
		iw.entity.AddComponent(ecs.ComponentEquipment, equipmentComp)
	}

	// Check if player can use this equipment
	if playerStats := iw.entity.RPGStats(); playerStats != nil {
		if !slot.Item.CanUse(playerStats.Level, playerStats.Job) {
			// Player can't use this item - could show a message here
			return
		}
	}

	equipment := slot.Item.Equipment
	equipSlot := equipment.Slot

	// Check if something is already equipped in that slot
	if currentEquip := equipmentComp.GetEquipped(equipSlot); currentEquip != nil {
		// Unequip current item and add it back to inventory
		unequipped := equipmentComp.Unequip(equipSlot)
		if unequipped != nil {
			// Create an inventory item for the unequipped equipment
			unequippedItem := &components.Item{
				ID:          unequipped.ID,
				Name:        unequipped.Name,
				Description: unequipped.Description,
				Type:        components.ItemTypeEquipment,
				Rarity:      components.ItemRarity(unequipped.Rarity),
				Value:       unequipped.Value,
				IconID:      unequipped.IconID,
				Equipment:   unequipped,
				Stackable:   false,
				MaxStack:    1,
			}

			// Try to add unequipped item back to inventory
			remaining := iw.inventory.AddItem(unequippedItem, 1)
			if remaining > 0 {
				// If inventory is full, re-equip the old item
				equipmentComp.Equip(unequipped)
				return
			}
		}
	}

	// Equip the new item
	equipmentComp.Equip(equipment)

	// Remove the item from inventory
	iw.inventory.RemoveItem(slot.Item.ID, 1)

	// Update the slot display to reflect changes
	iw.initializeSlots()
}

// handleActionPanelClick processes clicks in the action panel
func (iw *InventoryWidget) handleActionPanelClick() {
	panelX := iw.X + InventoryActionPanelX
	panelY := iw.Y + InventoryActionPanelY

	if iw.mouseX < panelX || iw.mouseX >= panelX+InventoryActionPanelWidth ||
		iw.mouseY < panelY || iw.mouseY >= panelY+InventoryActionPanelHeight {
		return
	}

	// Calculate which button was clicked
	buttonY := iw.mouseY - panelY
	buttonIndex := buttonY / (InventoryActionButtonHeight + InventoryActionButtonSpacing)

	switch buttonIndex {
	case 0: // Sort by Name
		iw.sortMode = "name"
		iw.sortInventory()
	case 1: // Sort by Type
		iw.sortMode = "type"
		iw.sortInventory()
	case 2: // Sort by Rarity
		iw.sortMode = "rarity"
		iw.sortInventory()
	case 3: // Filter Equipment
		iw.filterType = components.ItemTypeEquipment
		iw.showAllTypes = false
		iw.applyFilter()
	case 4: // Filter Consumables
		iw.filterType = components.ItemTypeConsumable
		iw.showAllTypes = false
		iw.applyFilter()
	case 5: // Show All
		iw.showAllTypes = true
		iw.applyFilter()
	}
}

// sortInventory sorts items based on current sort mode
func (iw *InventoryWidget) sortInventory() {
	if iw.inventory == nil {
		return
	}

	// TODO: Implement sorting in InventoryComponent
	// For now, just refresh slots
	iw.refreshSlots()
}

// applyFilter applies current filter settings
func (iw *InventoryWidget) applyFilter() {
	// TODO: Implement filtering
	// For now, just refresh slots
	iw.refreshSlots()
}

// refreshSlots updates visual slots from inventory data
func (iw *InventoryWidget) refreshSlots() {
	for row := 0; row < InventoryGridHeight; row++ {
		for col := 0; col < InventoryGridWidth; col++ {
			slot := &iw.slots[row][col]
			slotIndex := row*InventoryGridWidth + col

			if slotIndex < len(iw.inventory.Slots) {
				invSlot := &iw.inventory.Slots[slotIndex]
				slot.Item = invSlot.Item
				slot.Quantity = invSlot.Quantity
			} else {
				slot.Item = nil
				slot.Quantity = 0
			}
		}
	}
}

// handleKeyboardInput processes keyboard shortcuts
func (iw *InventoryWidget) handleKeyboardInput() {
	// Delete key to remove selected item
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) && iw.selectedSlot != nil {
		slotIndex := iw.getSlotIndex(iw.selectedSlot)
		if slotIndex != -1 {
			slotData := iw.inventory.GetSlotByIndex(slotIndex)
			if slotData != nil {
				slotData.RemoveItem(slotData.Quantity) // Remove all items from slot
			}
			iw.refreshSlots()
		}
	}

	// Number keys to split stacks
	for i := ebiten.Key1; i <= ebiten.Key9; i++ {
		if inpututil.IsKeyJustPressed(i) && iw.selectedSlot != nil {
			splitAmount := int(i - ebiten.Key1 + 1)
			iw.splitStack(iw.selectedSlot, splitAmount)
		}
	}
}

// splitStack splits a stack of items
func (iw *InventoryWidget) splitStack(slot *InventorySlot, amount int) {
	if slot.Item == nil || slot.Quantity <= 1 || amount >= slot.Quantity {
		return
	}

	slotIndex := iw.getSlotIndex(slot)
	if slotIndex == -1 {
		return
	}

	// Try to add split amount to an empty slot
	newItem := *slot.Item // Copy item
	addedAmount := iw.inventory.AddItem(&newItem, amount)
	if addedAmount > 0 {
		// Reduce original stack by the amount successfully added
		slotData := iw.inventory.GetSlotByIndex(slotIndex)
		if slotData != nil {
			slotData.RemoveItem(addedAmount)
		}
		iw.refreshSlots()
	}
}

// Draw renders the inventory widget
func (iw *InventoryWidget) Draw(screen *ebiten.Image) {
	if !iw.Visible {
		return
	}

	// Draw widget shadow
	iw.drawShadow(screen)

	// Draw widget background
	iw.drawBackground(screen)

	// Draw widget border
	iw.drawBorder(screen)

	// Draw header
	iw.drawHeader(screen)

	// Draw inventory grid
	iw.drawInventoryGrid(screen)

	// Draw action panel
	iw.drawActionPanel(screen)

	// Draw dragged item
	if iw.isDragging && iw.draggedItem != nil {
		iw.drawDraggedItem(screen)
	}

	// Draw tooltip
	if iw.showTooltip && iw.tooltipItem != nil {
		iw.drawTooltip(screen)
	}
}

// drawShadow draws the widget drop shadow
func (iw *InventoryWidget) drawShadow(screen *ebiten.Image) {
	shadowX := float32(iw.X + 5)
	shadowY := float32(iw.Y + 5)
	shadowW := float32(iw.Width)
	shadowH := float32(iw.Height)

	vector.FillRect(screen, shadowX, shadowY, shadowW, shadowH, InventoryWidgetShadow, false)
}

// drawBackground draws the widget background
func (iw *InventoryWidget) drawBackground(screen *ebiten.Image) {
	x := float32(iw.X)
	y := float32(iw.Y)
	w := float32(iw.Width)
	h := float32(iw.Height)

	vector.FillRect(screen, x, y, w, h, InventoryWidgetBackground, false)
}

// drawBorder draws the widget border
func (iw *InventoryWidget) drawBorder(screen *ebiten.Image) {
	x := float32(iw.X)
	y := float32(iw.Y)
	w := float32(iw.Width)
	h := float32(iw.Height)

	// Draw border lines
	for i := 0; i < InventoryWidgetBorderWidth; i++ {
		fi := float32(i)
		vector.StrokeRect(screen, x+fi, y+fi, w-fi*2, h-fi*2, 1, InventoryWidgetBorder, false)
	}
}

// drawHeader draws the header section
func (iw *InventoryWidget) drawHeader(screen *ebiten.Image) {
	_ = screen
	// Title
	titleText := "Inventory"
	titleX := float32(iw.X + InventoryTitleX)
	titleY := float32(iw.Y + InventoryTitleY)

	// TODO: Add text rendering when font system is available
	_ = titleText
	_ = titleX
	_ = titleY

	// Capacity stats
	if iw.inventory != nil {
		usedSlots := iw.inventory.GetUsedSlots()
		totalSlots := len(iw.inventory.Slots)
		statsText := fmt.Sprintf("Capacity: %d/%d", usedSlots, totalSlots)

		// TODO: Add text rendering for stats
		_ = statsText
	}
}

// drawInventoryGrid draws the inventory slot grid
func (iw *InventoryWidget) drawInventoryGrid(screen *ebiten.Image) {
	for row := 0; row < InventoryGridHeight; row++ {
		for col := 0; col < InventoryGridWidth; col++ {
			slot := &iw.slots[row][col]
			iw.drawSlot(screen, slot)
		}
	}
}

// drawSlot draws a single inventory slot
func (iw *InventoryWidget) drawSlot(screen *ebiten.Image, slot *InventorySlot) {
	x := float32(slot.X)
	y := float32(slot.Y)
	w := float32(slot.Width)
	h := float32(slot.Height)

	// Determine slot background color
	bgColor := InventorySlotEmpty
	if slot.Item != nil {
		bgColor = InventorySlotFilled
	}

	if slot.IsSelected {
		bgColor = InventorySlotSelected
	} else if slot.IsHovered {
		bgColor = InventorySlotHover
	}

	// Draw slot background
	vector.FillRect(screen, x, y, w, h, bgColor, false)

	// Determine border color
	borderColor := InventorySlotBorder
	if slot.IsSelected {
		borderColor = InventorySlotSelectedBorder
	} else if slot.Item != nil {
		// Use rarity color for border
		if rarityColor, exists := RarityColors[slot.Item.Rarity]; exists {
			borderColor = rarityColor
		}
	}

	// Draw slot border
	vector.StrokeRect(screen, x, y, w, h, 2, borderColor, false)

	// Draw item icon (placeholder)
	if slot.Item != nil && !iw.isDragging {
		iw.drawItemInSlot(screen, slot)
	}

	// Draw quantity if > 1
	if slot.Quantity > 1 {
		iw.drawQuantity(screen, slot)
	}
}

// drawItemInSlot draws an item within a slot
func (iw *InventoryWidget) drawItemInSlot(screen *ebiten.Image, slot *InventorySlot) {
	// TODO: Draw actual item icon when asset system is available
	// For now, draw a colored rectangle as placeholder

	itemX := float32(slot.X + 4)
	itemY := float32(slot.Y + 4)
	itemW := float32(slot.Width - 8)
	itemH := float32(slot.Height - 8)

	// Use item type color as placeholder
	var itemColor color.RGBA
	switch slot.Item.Type {
	case components.ItemTypeEquipment:
		itemColor = color.RGBA{200, 150, 100, 255}
	case components.ItemTypeConsumable:
		itemColor = color.RGBA{100, 200, 100, 255}
	case components.ItemTypeMaterial:
		itemColor = color.RGBA{150, 150, 200, 255}
	case components.ItemTypeQuest:
		itemColor = color.RGBA{255, 215, 0, 255}
	default:
		itemColor = color.RGBA{150, 150, 150, 255}
	}

	vector.FillRect(screen, itemX, itemY, itemW, itemH, itemColor, false)

	// Draw item name as text fallback (when sprites are not available)
	itemName := slot.Item.Name

	// Calculate max characters that fit in slot (6 pixels per character, with padding)
	maxChars := int((itemW - 4) / 6)
	if maxChars < 3 {
		maxChars = 3 // Minimum 3 characters
	}

	if len(itemName) > maxChars {
		if maxChars > 3 {
			itemName = itemName[:maxChars-3] + "..."
		} else {
			itemName = itemName[:maxChars]
		}
	}

	// Calculate proper text centering
	textWidth := len(itemName) * 6 // 6 pixels per character
	textX := int(itemX + (itemW-float32(textWidth))/2)
	textY := int(itemY + itemH/2 - 4) // Center vertically

	// Ensure text stays within slot bounds
	if textX < int(itemX)+2 {
		textX = int(itemX) + 2
	}

	ebitenutil.DebugPrintAt(screen, itemName, textX, textY)
}

// drawQuantity draws the item quantity number
func (iw *InventoryWidget) drawQuantity(screen *ebiten.Image, slot *InventorySlot) {
	if slot.Quantity > 1 {
		quantityText := fmt.Sprintf("%d", slot.Quantity)
		// Draw quantity in bottom-right corner of slot
		textX := slot.X + slot.Width - len(quantityText)*6 - 2
		textY := slot.Y + slot.Height - 10
		ebitenutil.DebugPrintAt(screen, quantityText, textX, textY)
	}
}

// drawActionPanel draws the action buttons panel
func (iw *InventoryWidget) drawActionPanel(screen *ebiten.Image) {
	panelX := float32(iw.X + InventoryActionPanelX)
	panelY := float32(iw.Y + InventoryActionPanelY)
	panelW := float32(InventoryActionPanelWidth)
	panelH := float32(InventoryActionPanelHeight)

	// Draw panel background
	vector.FillRect(screen, panelX, panelY, panelW, panelH, InventoryActionPanelBg, false)

	// Draw action buttons
	buttons := []string{"Sort Name", "Sort Type", "Sort Rarity", "Equipment", "Consumable", "Show All"}

	for i, buttonText := range buttons {
		buttonX := panelX + 5
		buttonY := panelY + 5 + float32(i*(InventoryActionButtonHeight+InventoryActionButtonSpacing))
		buttonW := panelW - 10
		buttonH := float32(InventoryActionButtonHeight)

		// Check if button is hovered
		isHovered := iw.mouseX >= int(buttonX) && iw.mouseX < int(buttonX+buttonW) &&
			iw.mouseY >= int(buttonY) && iw.mouseY < int(buttonY+buttonH)

		buttonColor := InventoryActionButton
		if isHovered {
			buttonColor = InventoryActionButtonHover
		}

		// Draw button
		vector.FillRect(screen, buttonX, buttonY, buttonW, buttonH, buttonColor, false)
		vector.StrokeRect(screen, buttonX, buttonY, buttonW, buttonH, 1, InventorySlotBorder, false)

		// TODO: Draw button text when font system is available
		_ = buttonText
	}
}

// drawDraggedItem draws the item being dragged
func (iw *InventoryWidget) drawDraggedItem(screen *ebiten.Image) {
	if iw.draggedItem == nil {
		return
	}

	// Draw semi-transparent item at mouse position
	dragX := float32(iw.mouseX - InventorySlotSize/2)
	dragY := float32(iw.mouseY - InventorySlotSize/2)
	dragW := float32(InventorySlotSize)
	dragH := float32(InventorySlotSize)

	// Draw slot background with reduced opacity
	dragBg := InventorySlotSelected
	dragBg.A = InventoryDragOpacity
	vector.FillRect(screen, dragX, dragY, dragW, dragH, dragBg, false)

	// Draw item (simplified for now)
	itemX := dragX + 4
	itemY := dragY + 4
	itemW := dragW - 8
	itemH := dragH - 8

	var itemColor color.RGBA
	switch iw.draggedItem.Type {
	case components.ItemTypeEquipment:
		itemColor = color.RGBA{200, 150, 100, InventoryDragOpacity}
	case components.ItemTypeConsumable:
		itemColor = color.RGBA{100, 200, 100, InventoryDragOpacity}
	case components.ItemTypeMaterial:
		itemColor = color.RGBA{150, 150, 200, InventoryDragOpacity}
	case components.ItemTypeQuest:
		itemColor = color.RGBA{255, 215, 0, InventoryDragOpacity}
	default:
		itemColor = color.RGBA{150, 150, 150, InventoryDragOpacity}
	}

	vector.FillRect(screen, itemX, itemY, itemW, itemH, itemColor, false)
}

// wrapText wraps text to fit within a specified character width
func (iw *InventoryWidget) wrapText(text string, maxWidth int) []string {
	if len(text) <= maxWidth {
		return []string{text}
	}

	var lines []string

	// Simple character-based wrapping for now
	for len(text) > maxWidth {
		// Find the last space before maxWidth
		breakPoint := maxWidth
		for breakPoint > 0 && text[breakPoint] != ' ' {
			breakPoint--
		}

		// If no space found, break at maxWidth
		if breakPoint == 0 {
			breakPoint = maxWidth
		}

		// Add the line
		lines = append(lines, text[:breakPoint])

		// Move to next part, skip space if we broke at a space
		text = text[breakPoint:]
		if len(text) > 0 && text[0] == ' ' {
			text = text[1:]
		}
	}

	// Add remaining text
	if len(text) > 0 {
		lines = append(lines, text)
	}

	return lines
}

// drawTooltip draws the item tooltip
func (iw *InventoryWidget) drawTooltip(screen *ebiten.Image) {
	if iw.tooltipItem == nil {
		return
	}

	// Prepare tooltip content with text wrapping
	item := iw.tooltipItem
	maxLineChars := 28 // Approximate characters that fit in tooltip width

	var contentLines []string

	// Item name (no wrapping, keep it short)
	contentLines = append(contentLines, item.Name)
	contentLines = append(contentLines, "")

	// Wrap description
	descriptionLines := iw.wrapText(item.Description, maxLineChars)
	contentLines = append(contentLines, descriptionLines...)
	contentLines = append(contentLines, "")

	// Basic info
	contentLines = append(contentLines, fmt.Sprintf("Type: %s", item.Type.String()))
	contentLines = append(contentLines, fmt.Sprintf("Rarity: %s", item.Rarity.String()))
	contentLines = append(contentLines, fmt.Sprintf("Value: %d gold", item.Value))

	// Add consumable effects if applicable
	if item.Type == components.ItemTypeConsumable && len(item.Effects) > 0 {
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Effects:")
		for _, effect := range item.Effects {
			contentLines = append(contentLines, fmt.Sprintf("  %s: %d", effect.Type, effect.Value))
		}
	}

	// Add stackable information
	if item.Stackable {
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, fmt.Sprintf("Stackable (Max: %d)", item.MaxStack))
	}

	// Add quest item notice
	if item.QuestItem {
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "* Quest Item *")
	}

	// Calculate tooltip dimensions based on content (with tighter spacing)
	tooltipW := float32(InventoryTooltipWidth)
	tooltipH := float32(len(contentLines)*12 + 16) // 12px per line + padding

	// Adjust position to stay on screen
	screenW, screenH := screen.Bounds().Dx(), screen.Bounds().Dy()
	tooltipX := float32(iw.tooltipX)
	tooltipY := float32(iw.tooltipY)

	if tooltipX+tooltipW > float32(screenW) {
		tooltipX = float32(screenW) - tooltipW - 10
	}
	if tooltipY+tooltipH > float32(screenH) {
		tooltipY = float32(screenH) - tooltipH - 10
	}

	// Draw tooltip background
	vector.FillRect(screen, tooltipX, tooltipY, tooltipW, tooltipH,
		InventoryTooltipBackground, false)

	// Draw tooltip border
	vector.StrokeRect(screen, tooltipX, tooltipY, tooltipW, tooltipH,
		float32(InventoryTooltipBorderWidth), InventoryTooltipBorderColor, false)

	// Draw tooltip content
	startX := int(tooltipX) + 8
	startY := int(tooltipY) + 8

	for i, line := range contentLines {
		y := startY + i*12 // Tighter line spacing for better fit

		if i == 0 {
			// Item name - make it stand out
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("=== %s ===", line), startX, y)
		} else {
			// All other content (wrapped descriptions will display naturally)
			ebitenutil.DebugPrintAt(screen, line, startX, y)
		}
	}
}

// Close hides the inventory widget
func (iw *InventoryWidget) Close() {
	iw.Visible = false
	iw.isDragging = false
	iw.draggedItem = nil
	iw.draggedSlot = nil
	iw.selectedSlot = nil
	iw.showTooltip = false
}

// IsOpen returns true if the widget is visible
func (iw *InventoryWidget) IsOpen() bool {
	return iw.Visible
}

// GetSelectedItem returns the currently selected item
func (iw *InventoryWidget) GetSelectedItem() *components.Item {
	if iw.selectedSlot != nil {
		return iw.selectedSlot.Item
	}
	return nil
}

// GetSelectedQuantity returns the quantity of the selected item
func (iw *InventoryWidget) GetSelectedQuantity() int {
	if iw.selectedSlot != nil {
		return iw.selectedSlot.Quantity
	}
	return 0
}

// RefreshInventory updates the widget from the entity's inventory
func (iw *InventoryWidget) RefreshInventory() {
	if inv := iw.entity.Inventory(); inv != nil {
		iw.inventory = inv
		iw.refreshSlots()
	}
}
