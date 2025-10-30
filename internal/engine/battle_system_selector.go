// Package engine provides battle system selection and management
package engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jrecuero/myrpg/internal/battle/classic"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/logger"
)

// BattleSystemType represents the different types of battle systems
type BattleSystemType int

const (
	BattleSystemTactical BattleSystemType = iota
	BattleSystemClassic
)

// BattleSystemSelector manages which battle system to use
type BattleSystemSelector struct {
	currentSystem BattleSystemType

	// Classic battle system
	classicManager  *classic.BattleManager
	classicRenderer *classic.BattleRenderer

	// Keep references to tactical system (already in Game)
	// We won't duplicate it here, just manage the selection
}

// NewBattleSystemSelector creates a new battle system selector
func NewBattleSystemSelector(screenWidth, screenHeight int) *BattleSystemSelector {
	// Create classic battle system components
	classicManager := classic.NewBattleManager()
	classicRenderer := classic.NewBattleRenderer(classicManager, screenWidth, screenHeight)

	// Set up callbacks
	classicManager.SetOnActionExecuted(func(action *classic.BattleAction) {
		// Add battle log messages
		if action.Target != nil {
			message := getBattleActionMessage(action)
			classicRenderer.AddBattleMessage(message)
		}
	})

	return &BattleSystemSelector{
		currentSystem:   BattleSystemClassic, // Default to classic for testing
		classicManager:  classicManager,
		classicRenderer: classicRenderer,
	}
}

// getBattleActionMessage formats a battle action into a readable message
func getBattleActionMessage(action *classic.BattleAction) string {
	attackerName := "Unknown"
	targetName := "Unknown"

	if attackerStats := action.Entity.RPGStats(); attackerStats != nil {
		attackerName = attackerStats.Name
	}

	if targetStats := action.Target.RPGStats(); targetStats != nil {
		targetName = targetStats.Name
	}

	switch action.ActionType {
	case classic.ActionAttack:
		return attackerName + " attacks " + targetName + "!"
	case classic.ActionMagic:
		return attackerName + " casts magic on " + targetName + "!"
	case classic.ActionDefend:
		return attackerName + " defends!"
	case classic.ActionItem:
		return attackerName + " uses an item!"
	case classic.ActionEscape:
		return attackerName + " tries to escape!"
	default:
		return attackerName + " acts!"
	}
}

// SetBattleSystem changes the current battle system
func (bss *BattleSystemSelector) SetBattleSystem(systemType BattleSystemType) {
	if bss.currentSystem != systemType {
		logger.Debug("🔄 Switching battle system from %s to %s",
			bss.getBattleSystemName(bss.currentSystem),
			bss.getBattleSystemName(systemType))

		bss.currentSystem = systemType
	}
}

// GetBattleSystem returns the current battle system type
func (bss *BattleSystemSelector) GetBattleSystem() BattleSystemType {
	return bss.currentSystem
}

// StartBattle initiates a battle using the selected system
func (bss *BattleSystemSelector) StartBattle(game *Game, playerParty, enemyParty []*ecs.Entity) error {
	logger.Debug("🗡️  Starting battle with %s system", bss.getBattleSystemName(bss.currentSystem))

	switch bss.currentSystem {
	case BattleSystemTactical:
		// Use existing tactical system
		participants := append(playerParty, enemyParty...)
		game.SwitchToTacticalMode(participants)
		return nil

	case BattleSystemClassic:
		// Use new classic system
		return bss.classicManager.StartBattle(playerParty, enemyParty)

	default:
		logger.Debug("❌ Unknown battle system type: %d", bss.currentSystem)
		return nil
	}
}

// IsClassicBattleActive returns true if a classic battle is currently running
func (bss *BattleSystemSelector) IsClassicBattleActive() bool {
	return bss.currentSystem == BattleSystemClassic && bss.classicManager.IsActive()
}

// IsTacticalBattleActive returns true if tactical battle should be used
func (bss *BattleSystemSelector) IsTacticalBattleActive() bool {
	return bss.currentSystem == BattleSystemTactical
}

// UpdateClassicBattle updates the classic battle system
func (bss *BattleSystemSelector) UpdateClassicBattle(deltaTime time.Duration) {
	if bss.currentSystem == BattleSystemClassic {
		bss.classicManager.Update(deltaTime)
	}
}

// RenderClassicBattle renders the classic battle system
func (bss *BattleSystemSelector) RenderClassicBattle(screen *ebiten.Image) {
	if bss.currentSystem == BattleSystemClassic {
		bss.classicRenderer.Render(screen)
	}
}

// GetClassicBattleManager returns the classic battle manager (for advanced control)
func (bss *BattleSystemSelector) GetClassicBattleManager() *classic.BattleManager {
	return bss.classicManager
}

// GetClassicBattleRenderer returns the classic battle renderer (for UI control)
func (bss *BattleSystemSelector) GetClassicBattleRenderer() *classic.BattleRenderer {
	return bss.classicRenderer
}

// ToggleBattleSystem switches between tactical and classic systems
func (bss *BattleSystemSelector) ToggleBattleSystem() {
	if bss.currentSystem == BattleSystemTactical {
		bss.SetBattleSystem(BattleSystemClassic)
	} else {
		bss.SetBattleSystem(BattleSystemTactical)
	}
}

// getBattleSystemName returns a readable name for the battle system type
func (bss *BattleSystemSelector) getBattleSystemName(systemType BattleSystemType) string {
	switch systemType {
	case BattleSystemTactical:
		return "Tactical"
	case BattleSystemClassic:
		return "Classic"
	default:
		return "Unknown"
	}
}
