package systems

import (
	"fmt"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// ConsumableManager handles consumable item usage and effects
type ConsumableManager struct{}

// NewConsumableManager creates a new consumable manager
func NewConsumableManager() *ConsumableManager {
	return &ConsumableManager{}
}

// UseConsumable applies the effects of a consumable item to a target entity
func (cm *ConsumableManager) UseConsumable(item *components.Item, user *ecs.Entity, target *ecs.Entity) error {
	if item.Type != components.ItemTypeConsumable {
		return fmt.Errorf("item %s is not consumable", item.Name)
	}

	// Apply each effect
	for _, effect := range item.Effects {
		err := cm.applyEffect(effect, user, target)
		if err != nil {
			return fmt.Errorf("failed to apply effect %s: %v", effect.Type, err)
		}
	}

	return nil
}

// applyEffect applies a single consumable effect
func (cm *ConsumableManager) applyEffect(effect components.ConsumableEffect, user *ecs.Entity, target *ecs.Entity) error {
	// Determine the actual target based on effect target specification
	actualTarget := target
	if effect.Target == "self" {
		actualTarget = user
	}

	// Get target's RPG stats component
	statsComp, exists := actualTarget.GetComponent(ecs.ComponentRPGStats)
	if !exists {
		return fmt.Errorf("target entity has no RPG stats component")
	}

	stats := statsComp.(*components.RPGStatsComponent)

	switch effect.Type {
	case "heal_hp":
		return cm.healHP(stats, effect.Value)
	case "heal_mp":
		return cm.healMP(stats, effect.Value)
	case "buff_attack":
		return cm.buffStat(stats, "attack", effect.Value)
	case "buff_defense":
		return cm.buffStat(stats, "defense", effect.Value)
	case "buff_magic_attack":
		return cm.buffStat(stats, "magic_attack", effect.Value)
	case "buff_magic_defense":
		return cm.buffStat(stats, "magic_defense", effect.Value)
	case "buff_speed":
		return cm.buffStat(stats, "speed", effect.Value)
	case "cure_all":
		return cm.cureAllStatusEffects(stats)
	default:
		return fmt.Errorf("unknown effect type: %s", effect.Type)
	}
}

// healHP restores HP to the target
func (cm *ConsumableManager) healHP(stats *components.RPGStatsComponent, amount int) error {
	if amount == 9999 { // Full heal
		stats.CurrentHP = stats.MaxHP
	} else {
		stats.CurrentHP += amount
		if stats.CurrentHP > stats.MaxHP {
			stats.CurrentHP = stats.MaxHP
		}
	}
	return nil
}

// healMP restores MP to the target
func (cm *ConsumableManager) healMP(stats *components.RPGStatsComponent, amount int) error {
	if amount == 9999 { // Full restore
		stats.CurrentMP = stats.MaxMP
	} else {
		stats.CurrentMP += amount
		if stats.CurrentMP > stats.MaxMP {
			stats.CurrentMP = stats.MaxMP
		}
	}
	return nil
}

// buffStat applies a temporary stat buff (would need buff system for duration)
func (cm *ConsumableManager) buffStat(stats *components.RPGStatsComponent, statName string, amount int) error {
	// For now, just apply direct temporary bonus
	// In a full system, this would create a temporary buff with duration
	switch statName {
	case "attack":
		stats.Attack += amount
	case "defense":
		stats.Defense += amount
	case "magic_attack":
		stats.MagicAttack += amount
	case "magic_defense":
		stats.MagicDefense += amount
	case "speed":
		stats.Speed += amount
	default:
		return fmt.Errorf("unknown stat: %s", statName)
	}
	return nil
}

// cureAllStatusEffects removes all negative status effects
func (cm *ConsumableManager) cureAllStatusEffects(stats *components.RPGStatsComponent) error {
	// In a full system, this would clear status effects like poison, sleep, etc.
	// For now, just ensure HP/MP are not negative
	if stats.CurrentHP < 0 {
		stats.CurrentHP = 1
	}
	if stats.CurrentMP < 0 {
		stats.CurrentMP = 0
	}
	return nil
}

// CanUseConsumable checks if a consumable can be used on a target
func (cm *ConsumableManager) CanUseConsumable(item *components.Item, user *ecs.Entity, target *ecs.Entity) bool {
	if item.Type != components.ItemTypeConsumable {
		return false
	}

	// Check if user meets level requirements
	if userStats, exists := user.GetComponent(ecs.ComponentRPGStats); exists {
		stats := userStats.(*components.RPGStatsComponent)
		if stats.Level < item.LevelRequirement {
			return false
		}

		// Check job restrictions
		if !item.CanUse(stats.Level, stats.Job) {
			return false
		}
	}

	// Check if target has RPG stats (needed for effects)
	_, exists := target.GetComponent(ecs.ComponentRPGStats)
	return exists
}

// GetConsumableDescription returns a description of what the consumable does
func (cm *ConsumableManager) GetConsumableDescription(item *components.Item) string {
	if item.Type != components.ItemTypeConsumable {
		return "Not a consumable item"
	}

	description := "Effects:\n"
	for _, effect := range item.Effects {
		switch effect.Type {
		case "heal_hp":
			if effect.Value == 9999 {
				description += "- Fully restores HP\n"
			} else {
				description += fmt.Sprintf("- Restores %d HP\n", effect.Value)
			}
		case "heal_mp":
			if effect.Value == 9999 {
				description += "- Fully restores MP\n"
			} else {
				description += fmt.Sprintf("- Restores %d MP\n", effect.Value)
			}
		case "buff_attack":
			description += fmt.Sprintf("- Increases Attack by %d (temporary)\n", effect.Value)
		case "buff_defense":
			description += fmt.Sprintf("- Increases Defense by %d (temporary)\n", effect.Value)
		case "buff_magic_attack":
			description += fmt.Sprintf("- Increases Magic Attack by %d (temporary)\n", effect.Value)
		case "buff_magic_defense":
			description += fmt.Sprintf("- Increases Magic Defense by %d (temporary)\n", effect.Value)
		case "buff_speed":
			description += fmt.Sprintf("- Increases Speed by %d (temporary)\n", effect.Value)
		case "cure_all":
			description += "- Cures all status effects\n"
		default:
			description += fmt.Sprintf("- %s: %d\n", effect.Type, effect.Value)
		}
	}

	return description
}

var (
	GlobalConsumableManager *ConsumableManager
)

// InitializeConsumableSystem sets up the global consumable manager
func InitializeConsumableSystem() {
	GlobalConsumableManager = NewConsumableManager()
}
