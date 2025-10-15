# Equipment Widget Constants

This document defines all constants used by the Equipment Widget, following the established patterns from the Character Stats Widget.

## Layout Dimensions

| Constant | Value | Description |
|----------|--------|-------------|
| `EquipmentWidgetWidth` | 500 | Widget width in pixels |
| `EquipmentWidgetHeight` | 600 | Widget height in pixels |
| `EquipmentWidgetBorderWidth` | 2 | Border thickness |
| `EquipmentWidgetShadowOffset` | 4 | Shadow offset distance |
| `EquipmentWidgetPadding` | 15 | Internal padding |
| `EquipmentWidgetTitleY` | 20 | Y offset for main title |
| `EquipmentWidgetContentStartY` | 50 | Y offset where content begins |
| `EquipmentWidgetBottomReserved` | 40 | Space reserved at bottom for help text |

## Equipment Slot Layout

| Constant | Value | Description |
|----------|--------|-------------|
| `EquipmentSlotSize` | 64 | Size of each equipment slot (square) |
| `EquipmentSlotBorder` | 2 | Border width for slots |
| `EquipmentSlotSpacing` | 10 | Space between adjacent slots |
| `EquipmentSlotCornerRadius` | 6 | Corner radius for slot borders |
| `EquipmentIconSize` | 48 | Size of equipment icons within slots |
| `EquipmentIconPadding` | 8 | Padding around icons in slots |

## Paperdoll Layout Positions (Relative to content area)

| Slot Type | X Position | Y Position | Description |
|-----------|------------|------------|-------------|
| `HeadSlotX` | 200 | 80 | Head/helmet slot position |
| `HeadSlotY` | 80 | 80 | Head/helmet slot position |
| `ChestSlotX` | 200 | 160 | Chest/armor slot position |
| `ChestSlotY` | 160 | 160 | Chest/armor slot position |
| `LegsSlotX` | 200 | 240 | Legs/pants slot position |
| `LegsSlotY` | 240 | 240 | Legs/pants slot position |
| `FeetSlotX` | 200 | 320 | Feet/boots slot position |
| `FeetSlotY` | 320 | 320 | Feet/boots slot position |
| `WeaponSlotX` | 120 | 160 | Main weapon slot position |
| `WeaponSlotY` | 160 | 160 | Main weapon slot position |
| `ShieldSlotX` | 280 | 160 | Shield/off-hand slot position |
| `ShieldSlotY` | 160 | 160 | Shield/off-hand slot position |
| `AccessorySlot1X` | 120 | 80 | First accessory slot position |
| `AccessorySlot1Y` | 80 | 80 | First accessory slot position |
| `AccessorySlot2X` | 280 | 80 | Second accessory slot position |
| `AccessorySlot2Y` | 80 | 80 | Second accessory slot position |

## Stat Comparison Panel

| Constant | Value | Description |
|----------|--------|-------------|
| `StatComparisonPanelX` | 350 | X position of stats comparison panel |
| `StatComparisonPanelY` | 80 | Y position of stats comparison panel |
| `StatComparisonPanelWidth` | 120 | Width of stats comparison panel |
| `StatComparisonPanelHeight` | 320 | Height of stats comparison panel |
| `StatComparisonLineHeight` | 16 | Height of each stat comparison line |
| `StatComparisonArrowWidth` | 12 | Width of stat change arrows |
| `StatComparisonValueWidth` | 40 | Width allocated for stat values |

## Color Constants (RGBA values)

### Background and Border Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `EquipmentWidgetBackgroundR` | 25 | 25 | 35 | 255 | Widget background |
| `EquipmentWidgetBackgroundG` | 25 | 25 | 35 | 255 | Widget background |
| `EquipmentWidgetBackgroundB` | 35 | 25 | 35 | 255 | Widget background |
| `EquipmentWidgetBackgroundA` | 255 | 25 | 35 | 255 | Widget background |
| `EquipmentWidgetBorderR` | 100 | 100 | 120 | 255 | Widget border |
| `EquipmentWidgetBorderG` | 120 | 100 | 120 | 255 | Widget border |
| `EquipmentWidgetBorderB` | 120 | 100 | 120 | 255 | Widget border |
| `EquipmentWidgetBorderA` | 255 | 100 | 120 | 255 | Widget border |

### Text Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `EquipmentWidgetTitleR` | 255 | 255 | 255 | 255 | Widget title text |
| `EquipmentWidgetTitleG` | 255 | 255 | 255 | 255 | Widget title text |
| `EquipmentWidgetTitleB` | 255 | 255 | 255 | 255 | Widget title text |
| `EquipmentWidgetTitleA` | 255 | 255 | 255 | 255 | Widget title text |
| `EquipmentWidgetHelpR` | 180 | 180 | 180 | 255 | Help text |
| `EquipmentWidgetHelpG` | 180 | 180 | 180 | 255 | Help text |
| `EquipmentWidgetHelpB` | 180 | 180 | 180 | 255 | Help text |
| `EquipmentWidgetHelpA` | 255 | 180 | 180 | 255 | Help text |

### Equipment Slot Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `EquipmentSlotEmptyR` | 40 | 40 | 50 | 255 | Empty slot background |
| `EquipmentSlotEmptyG` | 50 | 40 | 50 | 255 | Empty slot background |
| `EquipmentSlotEmptyB` | 50 | 40 | 50 | 255 | Empty slot background |
| `EquipmentSlotEmptyA` | 255 | 40 | 50 | 255 | Empty slot background |
| `EquipmentSlotBorderR` | 100 | 100 | 120 | 255 | Slot border |
| `EquipmentSlotBorderG` | 120 | 100 | 120 | 255 | Slot border |
| `EquipmentSlotBorderB` | 120 | 100 | 120 | 255 | Slot border |
| `EquipmentSlotBorderA` | 255 | 100 | 120 | 255 | Slot border |
| `EquipmentSlotEquippedR` | 60 | 80 | 100 | 255 | Equipped slot background |
| `EquipmentSlotEquippedG` | 100 | 80 | 100 | 255 | Equipped slot background |
| `EquipmentSlotEquippedB` | 100 | 80 | 100 | 255 | Equipped slot background |
| `EquipmentSlotEquippedA` | 255 | 80 | 100 | 255 | Equipped slot background |

### Equipment Rarity Colors
| Rarity | R | G | B | A | Description |
|--------|---|---|---|---|-------------|
| `RarityCommonR` | 200 | 200 | 200 | 255 | Common equipment |
| `RarityCommonG` | 200 | 200 | 200 | 255 | Common equipment |
| `RarityCommonB` | 200 | 200 | 200 | 255 | Common equipment |
| `RarityCommonA` | 255 | 200 | 200 | 255 | Common equipment |
| `RarityUncommonR` | 100 | 255 | 100 | 255 | Uncommon equipment |
| `RarityUncommonG` | 255 | 255 | 100 | 255 | Uncommon equipment |
| `RarityUncommonB` | 100 | 255 | 100 | 255 | Uncommon equipment |
| `RarityUncommonA` | 255 | 255 | 100 | 255 | Uncommon equipment |
| `RarityRareR` | 100 | 100 | 255 | 255 | Rare equipment |
| `RarityRareG` | 100 | 100 | 255 | 255 | Rare equipment |
| `RarityRareB` | 255 | 100 | 255 | 255 | Rare equipment |
| `RarityRareA` | 255 | 100 | 255 | 255 | Rare equipment |
| `RarityEpicR` | 200 | 100 | 255 | 255 | Epic equipment |
| `RarityEpicG` | 100 | 100 | 255 | 255 | Epic equipment |
| `RarityEpicB` | 255 | 100 | 255 | 255 | Epic equipment |
| `RarityEpicA` | 255 | 100 | 255 | 255 | Epic equipment |
| `RarityLegendaryR` | 255 | 200 | 50 | 255 | Legendary equipment |
| `RarityLegendaryG` | 200 | 200 | 50 | 255 | Legendary equipment |
| `RarityLegendaryB` | 50 | 200 | 50 | 255 | Legendary equipment |
| `RarityLegendaryA` | 255 | 200 | 50 | 255 | Legendary equipment |

### Stat Comparison Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `StatIncreaseR` | 100 | 255 | 100 | 255 | Positive stat change |
| `StatIncreaseG` | 255 | 255 | 100 | 255 | Positive stat change |
| `StatIncreaseB` | 100 | 255 | 100 | 255 | Positive stat change |
| `StatIncreaseA` | 255 | 255 | 100 | 255 | Positive stat change |
| `StatDecreaseR` | 255 | 100 | 100 | 255 | Negative stat change |
| `StatDecreaseG` | 100 | 100 | 100 | 255 | Negative stat change |
| `StatDecreaseB` | 100 | 100 | 100 | 255 | Negative stat change |
| `StatDecreaseA` | 255 | 100 | 100 | 255 | Negative stat change |
| `StatNoChangeR` | 200 | 200 | 200 | 255 | No stat change |
| `StatNoChangeG` | 200 | 200 | 200 | 255 | No stat change |
| `StatNoChangeB` | 200 | 200 | 200 | 255 | No stat change |
| `StatNoChangeA` | 255 | 200 | 200 | 255 | No stat change |

## Equipment Slot Enum

```go
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
```

## Equipment Rarity Enum

```go
type EquipmentRarity int

const (
	RarityCommon EquipmentRarity = iota
	RarityUncommon
	RarityRare
	RarityEpic
	RarityLegendary
)

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
```

## Navigation Controls

| Control | Action | Description |
|---------|--------|-------------|
| `ESC` | Close | Close equipment widget |
| `TAB` | Switch Focus | Switch between equipment slots |
| `ENTER` | Equip/Unequip | Toggle equipment in selected slot |
| `Arrow Keys` | Navigate | Move selection between equipment slots |
| `E` | Toggle Widget | Open/close equipment widget from exploration mode |

## Equipment Widget Help Text

```
"TAB: Next Slot  ESC: Close  ENTER: Equip/Unequip  Arrows: Navigate"
```

## Implementation Notes

### Equipment Slot Layout Pattern
The equipment slots are arranged in a paperdoll-style layout:
```
[Acc1]    [Head]    [Acc2]
[Weapon]  [Chest]   [Shield]
          [Legs]
          [Feet]
```

### Stat Comparison System
- When hovering over or selecting equipment, show stat changes
- Use color coding: green for improvements, red for decreases, white for no change
- Display format: "Attack: 15 → 18 (+3)"

### Equipment Restrictions
- Job-based restrictions (e.g., Mages can't equip heavy armor)
- Level requirements for equipment
- Stat requirements (minimum Strength for weapons, etc.)

### Visual Feedback
- Empty slots show subtle background with slot type icon
- Equipped slots highlight with different background
- Rarity affects border glow/color around equipment
- Selected slot has highlight border animation

This constants file provides the foundation for implementing the Equipment Widget following the established patterns from other UI components.