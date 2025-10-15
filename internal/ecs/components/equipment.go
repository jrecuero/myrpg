package components

import "fmt"

// EquipmentSlot represents different equipment slots on a character
type EquipmentSlot int

const (
	SlotHead EquipmentSlot = iota
	SlotChest
	SlotLegs
	SlotFeet
	SlotWeapon
	SlotShield
	SlotAccessory1
	SlotAccessory2
)

// String returns the string representation of an EquipmentSlot
func (s EquipmentSlot) String() string {
	switch s {
	case SlotHead:
		return "Head"
	case SlotChest:
		return "Chest"
	case SlotLegs:
		return "Legs"
	case SlotFeet:
		return "Feet"
	case SlotWeapon:
		return "Weapon"
	case SlotShield:
		return "Shield"
	case SlotAccessory1:
		return "Accessory 1"
	case SlotAccessory2:
		return "Accessory 2"
	default:
		return "Unknown"
	}
}

// EquipmentRarity represents the rarity/quality of equipment
type EquipmentRarity int

const (
	RarityCommon EquipmentRarity = iota
	RarityUncommon
	RarityRare
	RarityEpic
	RarityLegendary
)

// String returns the string representation of EquipmentRarity
func (r EquipmentRarity) String() string {
	switch r {
	case RarityCommon:
		return "Common"
	case RarityUncommon:
		return "Uncommon"
	case RarityRare:
		return "Rare"
	case RarityEpic:
		return "Epic"
	case RarityLegendary:
		return "Legendary"
	default:
		return "Unknown"
	}
}

// EquipmentStats represents the stat bonuses provided by equipment
type EquipmentStats struct {
	// Core stat bonuses
	AttackBonus     int // Bonus to Attack stat
	DefenseBonus    int // Bonus to Defense stat
	MagicPowerBonus int // Bonus to Magic Power stat
	MagicDefBonus   int // Bonus to Magic Defense stat
	SpeedBonus      int // Bonus to Speed stat
	HPBonus         int // Bonus to Max HP
	MPBonus         int // Bonus to Max MP

	// Combat bonuses
	CritChanceBonus int // Bonus to critical hit chance (percentage)
	CritDamageBonus int // Bonus to critical damage multiplier (percentage)
	AccuracyBonus   int // Bonus to hit accuracy (percentage)
	EvasionBonus    int // Bonus to evasion chance (percentage)

	// Special bonuses
	MovementBonus int // Bonus to movement range
	APBonus       int // Bonus to Action Points per turn
}

// Equipment represents a piece of equipment that can be worn by characters
type Equipment struct {
	ID          int             // Unique equipment ID
	Name        string          // Display name of the equipment
	Description string          // Detailed description
	Slot        EquipmentSlot   // Which slot this equipment occupies
	Rarity      EquipmentRarity // Rarity/quality level
	Stats       EquipmentStats  // Stat bonuses provided by this equipment

	// Requirements and restrictions
	LevelRequirement int       // Minimum character level to equip
	JobRestrictions  []JobType // Jobs that can equip this item (empty = all jobs)

	// Item properties
	Value  int // Gold value of the equipment
	IconID int // ID for the equipment icon/sprite
	SetID  int // Equipment set ID (0 = no set)
}

// CanEquip checks if a character with given level and job can equip this equipment
func (e *Equipment) CanEquip(level int, job JobType) bool {
	// Check level requirement
	if level < e.LevelRequirement {
		return false
	}

	// Check job restrictions (empty list means all jobs can equip)
	if len(e.JobRestrictions) == 0 {
		return true
	}

	// Check if the character's job is in the allowed list
	for _, allowedJob := range e.JobRestrictions {
		if job == allowedJob {
			return true
		}
	}

	return false
}

// GetStatDescription returns a formatted string describing the equipment's stat bonuses
func (e *Equipment) GetStatDescription() string {
	description := ""
	stats := e.Stats

	// Core stats
	if stats.AttackBonus != 0 {
		description += fmt.Sprintf("Attack: %+d\n", stats.AttackBonus)
	}
	if stats.DefenseBonus != 0 {
		description += fmt.Sprintf("Defense: %+d\n", stats.DefenseBonus)
	}
	if stats.MagicPowerBonus != 0 {
		description += fmt.Sprintf("Magic Power: %+d\n", stats.MagicPowerBonus)
	}
	if stats.MagicDefBonus != 0 {
		description += fmt.Sprintf("Magic Defense: %+d\n", stats.MagicDefBonus)
	}
	if stats.SpeedBonus != 0 {
		description += fmt.Sprintf("Speed: %+d\n", stats.SpeedBonus)
	}
	if stats.HPBonus != 0 {
		description += fmt.Sprintf("Max HP: %+d\n", stats.HPBonus)
	}
	if stats.MPBonus != 0 {
		description += fmt.Sprintf("Max MP: %+d\n", stats.MPBonus)
	}

	// Combat bonuses
	if stats.CritChanceBonus != 0 {
		description += fmt.Sprintf("Crit Chance: %+d%%\n", stats.CritChanceBonus)
	}
	if stats.CritDamageBonus != 0 {
		description += fmt.Sprintf("Crit Damage: %+d%%\n", stats.CritDamageBonus)
	}
	if stats.AccuracyBonus != 0 {
		description += fmt.Sprintf("Accuracy: %+d%%\n", stats.AccuracyBonus)
	}
	if stats.EvasionBonus != 0 {
		description += fmt.Sprintf("Evasion: %+d%%\n", stats.EvasionBonus)
	}

	// Special bonuses
	if stats.MovementBonus != 0 {
		description += fmt.Sprintf("Movement: %+d\n", stats.MovementBonus)
	}
	if stats.APBonus != 0 {
		description += fmt.Sprintf("Action Points: %+d\n", stats.APBonus)
	}

	return description
}

// EquipmentComponent manages the equipment worn by a character
type EquipmentComponent struct {
	Equipment map[EquipmentSlot]*Equipment // Currently equipped items by slot
}

// NewEquipmentComponent creates a new equipment component with empty slots
func NewEquipmentComponent() *EquipmentComponent {
	return &EquipmentComponent{
		Equipment: make(map[EquipmentSlot]*Equipment),
	}
}

// Equip puts an equipment item in the specified slot
func (ec *EquipmentComponent) Equip(equipment *Equipment) {
	ec.Equipment[equipment.Slot] = equipment
}

// Unequip removes equipment from the specified slot
func (ec *EquipmentComponent) Unequip(slot EquipmentSlot) *Equipment {
	equipment := ec.Equipment[slot]
	delete(ec.Equipment, slot)
	return equipment
}

// GetEquipped returns the equipment in the specified slot, or nil if empty
func (ec *EquipmentComponent) GetEquipped(slot EquipmentSlot) *Equipment {
	return ec.Equipment[slot]
}

// IsEquipped checks if there's equipment in the specified slot
func (ec *EquipmentComponent) IsEquipped(slot EquipmentSlot) bool {
	_, exists := ec.Equipment[slot]
	return exists
}

// GetTotalStats calculates the combined stat bonuses from all equipped items
func (ec *EquipmentComponent) GetTotalStats() EquipmentStats {
	total := EquipmentStats{}

	for _, equipment := range ec.Equipment {
		if equipment != nil {
			stats := equipment.Stats
			total.AttackBonus += stats.AttackBonus
			total.DefenseBonus += stats.DefenseBonus
			total.MagicPowerBonus += stats.MagicPowerBonus
			total.MagicDefBonus += stats.MagicDefBonus
			total.SpeedBonus += stats.SpeedBonus
			total.HPBonus += stats.HPBonus
			total.MPBonus += stats.MPBonus
			total.CritChanceBonus += stats.CritChanceBonus
			total.CritDamageBonus += stats.CritDamageBonus
			total.AccuracyBonus += stats.AccuracyBonus
			total.EvasionBonus += stats.EvasionBonus
			total.MovementBonus += stats.MovementBonus
			total.APBonus += stats.APBonus
		}
	}

	return total
}

// GetEquipmentList returns a list of all currently equipped items
func (ec *EquipmentComponent) GetEquipmentList() []*Equipment {
	var equipped []*Equipment

	// Return equipment in a consistent order
	slots := []EquipmentSlot{
		SlotHead, SlotChest, SlotLegs, SlotFeet,
		SlotWeapon, SlotShield, SlotAccessory1, SlotAccessory2,
	}

	for _, slot := range slots {
		if equipment := ec.GetEquipped(slot); equipment != nil {
			equipped = append(equipped, equipment)
		}
	}

	return equipped
}
