package engine

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// AttackType represents the type of attack in battle
type AttackType int

const (
	AttackPhysical AttackType = iota
	AttackMagical
	AttackCancel
)

func (at AttackType) String() string {
	switch at {
	case AttackPhysical:
		return "Physical"
	case AttackMagical:
		return "Magical"
	case AttackCancel:
		return "Cancel"
	default:
		return "Unknown"
	}
}

// BattleState represents the current state of battle
type BattleState int

const (
	BattleStateNone BattleState = iota
	BattleStatePlayerTurn
	BattleStateEnemyTurn
	BattleStateBattleEnd
)

// BattleSystem manages combat between players and enemies
type BattleSystem struct {
	State                   BattleState
	CurrentPlayer           *ecs.Entity
	CurrentEnemy            *ecs.Entity
	SelectedAttack          AttackType
	UIVisible               bool
	messageCallback         func(string)  // Callback to send messages to UI
	switchPlayerCallback    func()        // Callback to switch to next player
	AttackAnimationDuration time.Duration // Duration for attack animation feedback
}

// NewBattleSystem creates a new battle system
func NewBattleSystem() *BattleSystem {
	return &BattleSystem{
		State:                   BattleStateNone,
		CurrentPlayer:           nil,
		CurrentEnemy:            nil,
		SelectedAttack:          AttackPhysical,
		UIVisible:               false,
		AttackAnimationDuration: 1000 * time.Millisecond, // Default 1 second attack animation
	}
}

// SetMessageCallback sets the callback function for sending messages to UI
func (bs *BattleSystem) SetMessageCallback(callback func(string)) {
	bs.messageCallback = callback
}

// SetSwitchPlayerCallback sets the callback function for switching to next player
func (bs *BattleSystem) SetSwitchPlayerCallback(callback func()) {
	bs.switchPlayerCallback = callback
}

// SetAttackAnimationDuration sets the duration for attack animation feedback
func (bs *BattleSystem) SetAttackAnimationDuration(duration time.Duration) {
	bs.AttackAnimationDuration = duration
}

// triggerAttackAnimation triggers the attack animation for the current player
func (bs *BattleSystem) triggerAttackAnimation() {
	if bs.CurrentPlayer == nil {
		return
	}

	// Get animation component from the current player
	if animComponent := bs.CurrentPlayer.Animation(); animComponent != nil {
		// Check if the player has an attack animation
		if animComponent.HasAnimation(components.AnimationAttacking) {
			// Trigger temporary attack animation that reverts to idle after battle
			// This ensures the player returns to idle state after battle, not walking
			animComponent.SetTemporaryStateWithRevertTo(components.AnimationAttacking, bs.AttackAnimationDuration, components.AnimationIdle)
		}
	}
}

// addMessage sends a message to the UI if callback is set, otherwise prints to console
func (bs *BattleSystem) addMessage(msg string) {
	if bs.messageCallback != nil {
		bs.messageCallback(msg)
	} else {
		fmt.Println(msg)
	}
}

// StartBattle initiates a battle between a player and enemy
func (bs *BattleSystem) StartBattle(player, enemy *ecs.Entity) {
	bs.State = BattleStatePlayerTurn
	bs.CurrentPlayer = player
	bs.CurrentEnemy = enemy
	bs.SelectedAttack = AttackPhysical
	bs.UIVisible = true

	// Log battle start
	playerStats := player.RPGStats()
	enemyStats := enemy.RPGStats()
	if playerStats != nil && enemyStats != nil {
		bs.addMessage(fmt.Sprintf("Battle started: %s vs %s!", playerStats.Name, enemyStats.Name))
	}
}

// Update handles battle input and state transitions
func (bs *BattleSystem) Update() {
	if bs.State != BattleStatePlayerTurn || !bs.UIVisible {
		return
	}

	// Handle attack selection input
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		bs.SelectedAttack = AttackPhysical
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		bs.SelectedAttack = AttackMagical
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		bs.SelectedAttack = AttackCancel
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// Execute the selected attack
		bs.executePlayerAttack()
	}
}

// executePlayerAttack performs the player's attack and handles enemy retaliation
func (bs *BattleSystem) executePlayerAttack() {
	if bs.SelectedAttack == AttackCancel {
		bs.addMessage("Attack cancelled.")
		bs.endBattle()
		return
	}

	playerStats := bs.CurrentPlayer.RPGStats()
	enemyStats := bs.CurrentEnemy.RPGStats()

	if playerStats == nil || enemyStats == nil {
		bs.endBattle()
		return
	}

	// Calculate damage based on attack type
	var damage int
	var attackName string

	switch bs.SelectedAttack {
	case AttackPhysical:
		damage = bs.calculatePhysicalDamage(playerStats.Attack, enemyStats.Defense)
		attackName = "Physical Attack"
	case AttackMagical:
		damage = bs.calculateMagicalDamage(playerStats.MagicAttack, enemyStats.MagicDefense)
		attackName = "Magical Attack"
	}

	// Trigger attack animation for visual feedback
	bs.triggerAttackAnimation()

	// Apply damage to enemy
	enemyStats.CurrentHP -= damage
	if enemyStats.CurrentHP < 0 {
		enemyStats.CurrentHP = 0
	}

	bs.addMessage(fmt.Sprintf("%s uses %s on %s for %d damage! (%d HP remaining)",
		playerStats.Name, attackName, enemyStats.Name, damage, enemyStats.CurrentHP))

	// Check if enemy is defeated
	if enemyStats.CurrentHP <= 0 {
		bs.addMessage(fmt.Sprintf("%s is defeated!", enemyStats.Name))
		bs.endBattle()
		return
	}

	// Enemy retaliates with the same attack type
	bs.executeEnemyRetaliation()
}

// executeEnemyRetaliation handles enemy counter-attack
func (bs *BattleSystem) executeEnemyRetaliation() {
	playerStats := bs.CurrentPlayer.RPGStats()
	enemyStats := bs.CurrentEnemy.RPGStats()

	if playerStats == nil || enemyStats == nil {
		bs.endBattle()
		return
	}

	// Calculate retaliation damage using the same attack type
	var damage int
	var attackName string

	switch bs.SelectedAttack {
	case AttackPhysical:
		damage = bs.calculatePhysicalDamage(enemyStats.Attack, playerStats.Defense)
		attackName = "Physical Attack"
	case AttackMagical:
		damage = bs.calculateMagicalDamage(enemyStats.MagicAttack, playerStats.MagicDefense)
		attackName = "Magical Attack"
	}

	// Apply damage to player
	playerStats.CurrentHP -= damage
	if playerStats.CurrentHP < 0 {
		playerStats.CurrentHP = 0
	}

	bs.addMessage(fmt.Sprintf("%s retaliates with %s on %s for %d damage! (%d HP remaining)",
		enemyStats.Name, attackName, playerStats.Name, damage, playerStats.CurrentHP))

	// Check if player is defeated
	if playerStats.CurrentHP <= 0 {
		bs.addMessage(fmt.Sprintf("%s is defeated!", playerStats.Name))
	}

	bs.endBattle()
}

// calculatePhysicalDamage computes physical attack damage
func (bs *BattleSystem) calculatePhysicalDamage(attack, defense int) int {
	damage := attack - (defense / 2)
	if damage < 1 {
		damage = 1 // Minimum damage
	}
	return damage
}

// calculateMagicalDamage computes magical attack damage
func (bs *BattleSystem) calculateMagicalDamage(magicAttack, magicDefense int) int {
	damage := magicAttack - (magicDefense / 2)
	if damage < 1 {
		damage = 1 // Minimum damage
	}
	return damage
}

// endBattle concludes the current battle
func (bs *BattleSystem) endBattle() {
	bs.State = BattleStateNone
	bs.CurrentPlayer = nil
	bs.CurrentEnemy = nil
	bs.UIVisible = false
	bs.addMessage("Battle ended.")

	// Switch to next player after battle
	if bs.switchPlayerCallback != nil {
		bs.switchPlayerCallback()
	}
}

// IsInBattle returns true if a battle is currently active
func (bs *BattleSystem) IsInBattle() bool {
	return bs.State == BattleStatePlayerTurn && bs.UIVisible && bs.CurrentPlayer != nil && bs.CurrentEnemy != nil
}

// GetBattleMenuText returns the text to display for the battle menu
func (bs *BattleSystem) GetBattleMenuText() string {
	// Only show battle menu during active player turn
	if bs.State != BattleStatePlayerTurn || !bs.UIVisible || bs.CurrentPlayer == nil || bs.CurrentEnemy == nil {
		return ""
	}

	playerStats := bs.CurrentPlayer.RPGStats()
	enemyStats := bs.CurrentEnemy.RPGStats()

	if playerStats == nil || enemyStats == nil {
		return ""
	}

	text := fmt.Sprintf("BATTLE: %s vs %s\n", playerStats.Name, enemyStats.Name)
	text += "Select attack:\n"

	// Highlight selected option
	if bs.SelectedAttack == AttackPhysical {
		text += "> [1] Physical Attack\n"
	} else {
		text += "  [1] Physical Attack\n"
	}

	if bs.SelectedAttack == AttackMagical {
		text += "> [2] Magical Attack\n"
	} else {
		text += "  [2] Magical Attack\n"
	}

	if bs.SelectedAttack == AttackCancel {
		text += "> [ESC] Cancel\n"
	} else {
		text += "  [ESC] Cancel\n"
	}

	text += "\nPress ENTER to confirm"

	return text
}
