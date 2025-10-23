package components

import "fmt"

// ItemType represents the category of an item
type ItemType int

const (
	ItemTypeEquipment ItemType = iota
	ItemTypeConsumable
	ItemTypeMaterial
	ItemTypeQuest
	ItemTypeMiscellaneous
)

// String returns the string representation of an ItemType
func (t ItemType) String() string {
	switch t {
	case ItemTypeEquipment:
		return "Equipment"
	case ItemTypeConsumable:
		return "Consumable"
	case ItemTypeMaterial:
		return "Material"
	case ItemTypeQuest:
		return "Quest Item"
	case ItemTypeMiscellaneous:
		return "Miscellaneous"
	default:
		return "Unknown"
	}
}

// ItemRarity represents the rarity/quality of an item
type ItemRarity int

const (
	ItemRarityCommon ItemRarity = iota
	ItemRarityUncommon
	ItemRarityRare
	ItemRarityEpic
	ItemRarityLegendary
)

// String returns the string representation of ItemRarity
func (r ItemRarity) String() string {
	switch r {
	case ItemRarityCommon:
		return "Common"
	case ItemRarityUncommon:
		return "Uncommon"
	case ItemRarityRare:
		return "Rare"
	case ItemRarityEpic:
		return "Epic"
	case ItemRarityLegendary:
		return "Legendary"
	default:
		return "Unknown"
	}
}

// GetRarityColor returns a color associated with the rarity for UI display
func (r ItemRarity) GetRarityColor() (uint8, uint8, uint8) {
	switch r {
	case ItemRarityCommon:
		return 200, 200, 200 // Light Gray
	case ItemRarityUncommon:
		return 30, 255, 0 // Green
	case ItemRarityRare:
		return 0, 112, 255 // Blue
	case ItemRarityEpic:
		return 163, 53, 238 // Purple
	case ItemRarityLegendary:
		return 255, 128, 0 // Orange
	default:
		return 255, 255, 255 // White
	}
}

// ConsumableEffect represents what happens when a consumable item is used
type ConsumableEffect struct {
	Type   string // "heal_hp", "heal_mp", "buff_stat", "cure_status", etc.
	Value  int    // Amount of healing, stat bonus, etc.
	Target string // "self", "ally", "enemy", "area"
}

// Item represents any item in the game (equipment, consumables, materials, etc.)
type Item struct {
	ID          int        // Unique item ID
	Name        string     // Display name
	Description string     // Detailed description
	Type        ItemType   // Category of item
	Rarity      ItemRarity // Quality level
	Value       int        // Gold value
	IconID      int        // ID for item icon/sprite
	Stackable   bool       // Can this item stack in inventory
	MaxStack    int        // Maximum stack size (0 = no limit)

	// Equipment-specific data (only used if Type == ItemTypeEquipment)
	Equipment *Equipment // Nil if not equipment

	// Consumable-specific data (only used if Type == ItemTypeConsumable)
	Effects []ConsumableEffect // Effects when consumed

	// Requirements and restrictions
	LevelRequirement int       // Minimum level to use
	JobRestrictions  []JobType // Jobs that can use this item (empty = all jobs)

	// Quest and special properties
	QuestItem bool // Is this a quest item (cannot be sold/dropped)
	Unique    bool // Only one can be owned at a time
	SetID     int  // Item set ID (0 = no set)
}

// CanUse checks if a character with given level and job can use this item
func (i *Item) CanUse(level int, job JobType) bool {
	// Check level requirement
	if level < i.LevelRequirement {
		return false
	}

	// Check job restrictions (empty list means all jobs can use)
	if len(i.JobRestrictions) == 0 {
		return true
	}

	// Check if the character's job is in the allowed list
	for _, allowedJob := range i.JobRestrictions {
		if job == allowedJob {
			return true
		}
	}

	return false
}

// GetTypeDescription returns a formatted description based on item type
func (i *Item) GetTypeDescription() string {
	switch i.Type {
	case ItemTypeEquipment:
		if i.Equipment != nil {
			return i.Equipment.GetStatDescription()
		}
		return "Equipment item"
	case ItemTypeConsumable:
		description := "Consumable effects:\n"
		for _, effect := range i.Effects {
			description += fmt.Sprintf("- %s: %d\n", effect.Type, effect.Value)
		}
		return description
	case ItemTypeMaterial:
		return "Crafting material"
	case ItemTypeQuest:
		return "Important quest item"
	case ItemTypeMiscellaneous:
		return "Miscellaneous item"
	default:
		return "Unknown item type"
	}
}

// GetTooltipText returns formatted tooltip text for UI display
func (i *Item) GetTooltipText() string {
	tooltip := fmt.Sprintf("%s\n", i.Name)
	tooltip += fmt.Sprintf("Type: %s\n", i.Type.String())
	tooltip += fmt.Sprintf("Rarity: %s\n", i.Rarity.String())
	tooltip += fmt.Sprintf("Value: %d gold\n", i.Value)

	if i.LevelRequirement > 1 {
		tooltip += fmt.Sprintf("Required Level: %d\n", i.LevelRequirement)
	}

	if len(i.JobRestrictions) > 0 {
		tooltip += "Job Restrictions: "
		for i, job := range i.JobRestrictions {
			if i > 0 {
				tooltip += ", "
			}
			tooltip += job.String()
		}
		tooltip += "\n"
	}

	tooltip += "\n" + i.Description

	typeDesc := i.GetTypeDescription()
	if typeDesc != "" {
		tooltip += "\n\n" + typeDesc
	}

	return tooltip
}

// InventorySlot represents a slot in the inventory containing an item and quantity
type InventorySlot struct {
	Item     *Item // The item in this slot (nil if empty)
	Quantity int   // How many of this item (0 if empty)
}

// IsEmpty returns true if this slot contains no items
func (slot *InventorySlot) IsEmpty() bool {
	return slot.Item == nil || slot.Quantity <= 0
}

// CanAddItem returns true if the specified item can be added to this slot
func (slot *InventorySlot) CanAddItem(item *Item, quantity int) bool {
	// Empty slot can accept any item
	if slot.IsEmpty() {
		return true
	}

	// Non-empty slot can only accept same stackable items
	if slot.Item.ID == item.ID && item.Stackable {
		// Check if adding quantity would exceed max stack
		if item.MaxStack > 0 {
			return slot.Quantity+quantity <= item.MaxStack
		}
		return true // No stack limit
	}

	return false
}

// AddItem adds the specified quantity of an item to this slot
// Returns the quantity that couldn't be added (overflow)
func (slot *InventorySlot) AddItem(item *Item, quantity int) int {
	if !slot.CanAddItem(item, quantity) {
		return quantity // Couldn't add any
	}

	if slot.IsEmpty() {
		// Empty slot - add the item
		slot.Item = item
		if item.MaxStack > 0 && quantity > item.MaxStack {
			slot.Quantity = item.MaxStack
			return quantity - item.MaxStack
		}
		slot.Quantity = quantity
		return 0
	} else {
		// Existing item - add to stack
		if item.MaxStack > 0 {
			maxAddable := item.MaxStack - slot.Quantity
			if quantity <= maxAddable {
				slot.Quantity += quantity
				return 0
			} else {
				slot.Quantity = item.MaxStack
				return quantity - maxAddable
			}
		} else {
			slot.Quantity += quantity
			return 0
		}
	}
}

// RemoveItem removes the specified quantity from this slot
// Returns the quantity that was actually removed
func (slot *InventorySlot) RemoveItem(quantity int) int {
	if slot.IsEmpty() {
		return 0
	}

	if quantity >= slot.Quantity {
		// Remove all items
		removed := slot.Quantity
		slot.Item = nil
		slot.Quantity = 0
		return removed
	} else {
		// Remove partial quantity
		slot.Quantity -= quantity
		return quantity
	}
}

// InventoryComponent represents a character's inventory
type InventoryComponent struct {
	Slots    []InventorySlot // Array of inventory slots
	Width    int             // Grid width (for UI layout)
	Height   int             // Grid height (for UI layout)
	Capacity int             // Total number of slots
}

// NewInventoryComponent creates a new inventory with the specified dimensions
func NewInventoryComponent(width, height int) *InventoryComponent {
	capacity := width * height
	return &InventoryComponent{
		Slots:    make([]InventorySlot, capacity),
		Width:    width,
		Height:   height,
		Capacity: capacity,
	}
}

// GetSlot returns the inventory slot at the specified grid position
func (inv *InventoryComponent) GetSlot(x, y int) *InventorySlot {
	if x < 0 || x >= inv.Width || y < 0 || y >= inv.Height {
		return nil
	}
	index := y*inv.Width + x
	return &inv.Slots[index]
}

// GetSlotByIndex returns the inventory slot at the specified index
func (inv *InventoryComponent) GetSlotByIndex(index int) *InventorySlot {
	if index < 0 || index >= inv.Capacity {
		return nil
	}
	return &inv.Slots[index]
}

// FindItem searches for an item by ID and returns the first slot containing it
func (inv *InventoryComponent) FindItem(itemID int) *InventorySlot {
	for i := range inv.Slots {
		slot := &inv.Slots[i]
		if !slot.IsEmpty() && slot.Item.ID == itemID {
			return slot
		}
	}
	return nil
}

// FindEmptySlot returns the first empty slot, or nil if inventory is full
func (inv *InventoryComponent) FindEmptySlot() *InventorySlot {
	for i := range inv.Slots {
		slot := &inv.Slots[i]
		if slot.IsEmpty() {
			return slot
		}
	}
	return nil
}

// AddItem attempts to add an item to the inventory
// Returns the quantity that couldn't be added (0 if all added successfully)
func (inv *InventoryComponent) AddItem(item *Item, quantity int) int {
	remaining := quantity

	// First, try to add to existing stacks of the same item
	if item.Stackable {
		for i := range inv.Slots {
			slot := &inv.Slots[i]
			if !slot.IsEmpty() && slot.Item.ID == item.ID {
				overflow := slot.AddItem(item, remaining)
				remaining = overflow
				if remaining <= 0 {
					return 0 // All added successfully
				}
			}
		}
	}

	// Then, try to add to empty slots
	for remaining > 0 {
		emptySlot := inv.FindEmptySlot()
		if emptySlot == nil {
			break // No more empty slots
		}
		overflow := emptySlot.AddItem(item, remaining)
		remaining = overflow
	}

	return remaining
}

// RemoveItem removes the specified quantity of an item from inventory
// Returns the quantity that was actually removed
func (inv *InventoryComponent) RemoveItem(itemID int, quantity int) int {
	removed := 0
	remaining := quantity

	for i := range inv.Slots {
		slot := &inv.Slots[i]
		if !slot.IsEmpty() && slot.Item.ID == itemID {
			slotRemoved := slot.RemoveItem(remaining)
			removed += slotRemoved
			remaining -= slotRemoved
			if remaining <= 0 {
				break
			}
		}
	}

	return removed
}

// GetItemCount returns the total quantity of an item in the inventory
func (inv *InventoryComponent) GetItemCount(itemID int) int {
	count := 0
	for i := range inv.Slots {
		slot := &inv.Slots[i]
		if !slot.IsEmpty() && slot.Item.ID == itemID {
			count += slot.Quantity
		}
	}
	return count
}

// IsFull returns true if there are no empty slots
func (inv *InventoryComponent) IsFull() bool {
	return inv.FindEmptySlot() == nil
}

// GetUsedSlots returns the number of slots that contain items
func (inv *InventoryComponent) GetUsedSlots() int {
	used := 0
	for i := range inv.Slots {
		if !inv.Slots[i].IsEmpty() {
			used++
		}
	}
	return used
}
