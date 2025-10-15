package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/ui"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Game represents the test game state
type Game struct {
	equipmentWidget *ui.EquipmentWidget
	equipmentComp   *components.EquipmentComponent
	statsComp       *components.RPGStatsComponent
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	// Update equipment widget
	g.equipmentWidget.Update()

	return nil
}

// Draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw equipment widget
	g.equipmentWidget.Draw(screen)
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// createMockEquipment creates sample equipment items for testing
func createMockEquipment() []*components.Equipment {
	return []*components.Equipment{
		{
			ID:          1,
			Name:        "Iron Sword",
			Description: "A sturdy iron sword with a sharp edge.",
			Slot:        components.SlotWeapon,
			Rarity:      components.RarityCommon,
			Stats: components.EquipmentStats{
				AttackBonus:     12,
				CritChanceBonus: 5,
			},
			LevelRequirement: 1,
			JobRestrictions:  []components.JobType{components.JobWarrior, components.JobArcher},
			Value:            100,
			IconID:           1,
		},
		{
			ID:          2,
			Name:        "Leather Armor",
			Description: "Basic leather armor providing moderate protection.",
			Slot:        components.SlotChest,
			Rarity:      components.RarityCommon,
			Stats: components.EquipmentStats{
				DefenseBonus: 8,
				HPBonus:      15,
			},
			LevelRequirement: 1,
			JobRestrictions:  []components.JobType{},
			Value:            75,
			IconID:           2,
		},
		{
			ID:          3,
			Name:        "Mystic Crown",
			Description: "A magical crown that enhances mental abilities.",
			Slot:        components.SlotHead,
			Rarity:      components.RarityRare,
			Stats: components.EquipmentStats{
				MagicPowerBonus: 20,
				MPBonus:         25,
				DefenseBonus:    5,
			},
			LevelRequirement: 10,
			JobRestrictions:  []components.JobType{components.JobMage, components.JobCleric},
			Value:            500,
			IconID:           3,
		},
		{
			ID:          4,
			Name:        "Swift Boots",
			Description: "Enchanted boots that increase movement speed.",
			Slot:        components.SlotFeet,
			Rarity:      components.RarityUncommon,
			Stats: components.EquipmentStats{
				SpeedBonus:    8,
				MovementBonus: 1,
				EvasionBonus:  10,
			},
			LevelRequirement: 5,
			JobRestrictions:  []components.JobType{},
			Value:            250,
			IconID:           4,
		},
		{
			ID:          5,
			Name:        "Dragon Scale Shield",
			Description: "A legendary shield made from ancient dragon scales.",
			Slot:        components.SlotShield,
			Rarity:      components.RarityLegendary,
			Stats: components.EquipmentStats{
				DefenseBonus:    25,
				MagicDefBonus:   20,
				HPBonus:         50,
				CritDamageBonus: -10, // Reduces critical damage taken
			},
			LevelRequirement: 20,
			JobRestrictions:  []components.JobType{components.JobWarrior, components.JobCleric},
			Value:            2000,
			IconID:           5,
		},
		{
			ID:          6,
			Name:        "Ring of Power",
			Description: "An ancient ring that amplifies magical abilities.",
			Slot:        components.SlotAccessory1,
			Rarity:      components.RarityEpic,
			Stats: components.EquipmentStats{
				MagicPowerBonus: 15,
				MPBonus:         20,
				CritChanceBonus: 8,
			},
			LevelRequirement: 15,
			JobRestrictions:  []components.JobType{},
			Value:            800,
			IconID:           6,
		},
		{
			ID:          7,
			Name:        "Assassin's Leggings",
			Description: "Dark leggings that enhance stealth and agility.",
			Slot:        components.SlotLegs,
			Rarity:      components.RarityRare,
			Stats: components.EquipmentStats{
				SpeedBonus:      12,
				EvasionBonus:    15,
				CritChanceBonus: 10,
				DefenseBonus:    6,
			},
			LevelRequirement: 12,
			JobRestrictions:  []components.JobType{components.JobRogue, components.JobArcher},
			Value:            600,
			IconID:           7,
		},
		{
			ID:          8,
			Name:        "Amulet of Vitality",
			Description: "A protective amulet that enhances life force.",
			Slot:        components.SlotAccessory2,
			Rarity:      components.RarityUncommon,
			Stats: components.EquipmentStats{
				HPBonus:       40,
				DefenseBonus:  5,
				MagicDefBonus: 5,
			},
			LevelRequirement: 8,
			JobRestrictions:  []components.JobType{},
			Value:            300,
			IconID:           8,
		},
	}
}

// setupMockCharacter creates a character with sample equipment
func setupMockCharacter() (*components.EquipmentComponent, *components.RPGStatsComponent) {
	// Create equipment component
	equipmentComp := components.NewEquipmentComponent()

	// Create character stats
	statsComp := &components.RPGStatsComponent{
		Name:           "Test Warrior",
		Level:          15,
		Experience:     12500,
		ExpToNext:      2500,
		CurrentHP:      85,
		MaxHP:          100,
		CurrentMP:      60,
		MaxMP:          80,
		Attack:         35,
		Defense:        25,
		MagicAttack:    20,
		MagicDefense:   18,
		Speed:          28,
		Job:            components.JobWarrior,
		MoveRange:      3,
		MovesRemaining: 2,
		Accuracy:       85,
		CritRate:       5,
	}

	// Create mock equipment
	mockEquipment := createMockEquipment()

	// Equip some items for demonstration
	equipmentComp.Equip(mockEquipment[0]) // Iron Sword in weapon slot
	equipmentComp.Equip(mockEquipment[1]) // Leather Armor in chest slot
	equipmentComp.Equip(mockEquipment[3]) // Swift Boots in feet slot
	equipmentComp.Equip(mockEquipment[5]) // Ring of Power in accessory1 slot

	return equipmentComp, statsComp
}

func main() {
	// Setup mock character data
	equipmentComp, statsComp := setupMockCharacter()

	// Create equipment widget
	equipmentWidget := ui.NewEquipmentWidget(ScreenWidth, ScreenHeight, equipmentComp, statsComp)
	equipmentWidget.Show() // Show immediately for testing

	// Create game instance
	game := &Game{
		equipmentWidget: equipmentWidget,
		equipmentComp:   equipmentComp,
		statsComp:       statsComp,
	}

	// Set window properties
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Equipment Widget Test")
	ebiten.SetWindowResizable(false)

	// Run game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
