// Package classic implements a Dragon Quest-style battle system
// with enemy formations, speed-based turn order, and activity queues
package classic

import (
	"fmt"
	"sort"
	"time"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/logger"
)

// BattleState represents the current state of the battle
type BattleState int

const (
	BattleStateIdle BattleState = iota
	BattleStatePlayerTurn
	BattleStateEnemyTurn
	BattleStateAnimation
	BattleStateVictory
	BattleStateDefeat
	BattleStateEscaped
)

// BattleAction represents an action in the battle queue
type BattleAction struct {
	Entity     *ecs.Entity
	ActionType ActionType
	Target     *ecs.Entity
	Speed      int
	Timestamp  time.Time
}

// ActionType represents different types of battle actions
type ActionType int

const (
	ActionAttack ActionType = iota
	ActionMagic
	ActionItem
	ActionDefend
	ActionEscape
)

// ActivityEntry represents an entity's position in the activity queue
type ActivityEntry struct {
	Entity       *ecs.Entity
	NextActionAt time.Time
	Speed        int
	ActionDelay  time.Duration
}

// BattleManager manages the Dragon Quest-style battle system
type BattleManager struct {
	// Battle state
	state         BattleState
	battleStarted bool

	// Participants
	playerParty []*ecs.Entity
	enemyParty  []*ecs.Entity

	// Activity queue and turn management
	activityQueue []*ActivityEntry
	actionQueue   []*BattleAction

	// Formation positioning
	playerFormation *Formation
	enemyFormation  *Formation

	// Battle timing
	battleTime    time.Time
	speedModifier float64

	// Callbacks
	onBattleEnd      func(victory bool)
	onActionExecuted func(action *BattleAction)
}

// NewBattleManager creates a new Dragon Quest-style battle manager
func NewBattleManager() *BattleManager {
	return &BattleManager{
		state:         BattleStateIdle,
		battleStarted: false,
		playerParty:   make([]*ecs.Entity, 0),
		enemyParty:    make([]*ecs.Entity, 0),
		activityQueue: make([]*ActivityEntry, 0),
		actionQueue:   make([]*BattleAction, 0),
		speedModifier: 1000.0, // milliseconds per speed point
		battleTime:    time.Now(),
	}
}

// StartBattle initializes a new battle with the given parties
func (bm *BattleManager) StartBattle(playerParty, enemyParty []*ecs.Entity) error {
	logger.Debug("🗡️  Starting Dragon Quest-style battle")
	logger.Debug("   Players: %d, Enemies: %d", len(playerParty), len(enemyParty))

	if len(playerParty) == 0 || len(enemyParty) == 0 {
		return fmt.Errorf("both parties must have at least one member")
	}

	bm.playerParty = playerParty
	bm.enemyParty = enemyParty
	bm.battleStarted = true
	bm.state = BattleStatePlayerTurn
	bm.battleTime = time.Now()

	// Initialize formations
	bm.playerFormation = NewPlayerFormation(playerParty)
	bm.enemyFormation = NewEnemyFormation(enemyParty)

	// Initialize activity queue
	bm.initializeActivityQueue()

	logger.Debug("✅ Battle started successfully")
	return nil
}

// initializeActivityQueue sets up the speed-based activity queue
func (bm *BattleManager) initializeActivityQueue() {
	bm.activityQueue = make([]*ActivityEntry, 0)

	// Add all participants to activity queue
	allParticipants := append(bm.playerParty, bm.enemyParty...)

	for _, entity := range allParticipants {
		if stats := entity.RPGStats(); stats != nil {
			speed := bm.calculateEntitySpeed(entity)
			delay := time.Duration(bm.speedModifier/float64(speed)) * time.Millisecond

			entry := &ActivityEntry{
				Entity:       entity,
				NextActionAt: bm.battleTime.Add(delay),
				Speed:        speed,
				ActionDelay:  delay,
			}

			bm.activityQueue = append(bm.activityQueue, entry)

			logger.Debug("   Added to queue: %s (speed: %d, delay: %v)",
				entity.GetID(), speed, delay)
		}
	}

	// Sort by next action time
	bm.sortActivityQueue()
}

// calculateEntitySpeed determines an entity's speed based on stats and job
func (bm *BattleManager) calculateEntitySpeed(entity *ecs.Entity) int {
	stats := entity.RPGStats()
	if stats == nil {
		return 10 // default speed
	}

	baseSpeed := 10

	// Job-based speed modifiers
	switch stats.Job {
	case components.JobRogue:
		baseSpeed = 15 // Rogues are fast
	case components.JobMage:
		baseSpeed = 8 // Mages are slower
	case components.JobWarrior:
		baseSpeed = 12 // Warriors are moderate
	default:
		baseSpeed = 10
	}

	// Level scaling
	levelBonus := stats.Level * 2

	// Add some randomness (±20%)
	variation := (baseSpeed + levelBonus) / 5
	if variation < 1 {
		variation = 1
	}

	finalSpeed := baseSpeed + levelBonus // + random variation in actual implementation

	return finalSpeed
}

// sortActivityQueue sorts the queue by next action time
func (bm *BattleManager) sortActivityQueue() {
	sort.Slice(bm.activityQueue, func(i, j int) bool {
		return bm.activityQueue[i].NextActionAt.Before(bm.activityQueue[j].NextActionAt)
	})
}

// Update processes the battle logic each frame
func (bm *BattleManager) Update(deltaTime time.Duration) {
	if !bm.battleStarted || bm.state == BattleStateIdle {
		return
	}

	bm.battleTime = bm.battleTime.Add(deltaTime)

	// Check if any entities are ready to act
	bm.processActivityQueue()

	// Execute queued actions
	bm.processActionQueue()

	// Check for battle end conditions
	bm.checkBattleEndConditions()
}

// processActivityQueue checks for entities ready to act
func (bm *BattleManager) processActivityQueue() {
	for len(bm.activityQueue) > 0 {
		entry := bm.activityQueue[0]

		if bm.battleTime.Before(entry.NextActionAt) {
			break // No one is ready to act yet
		}

		// This entity is ready to act
		bm.scheduleEntityAction(entry)

		// Update the entity's next action time
		entry.NextActionAt = bm.battleTime.Add(entry.ActionDelay)

		// Re-sort the queue
		bm.sortActivityQueue()
	}
}

// scheduleEntityAction determines what action an entity will take
func (bm *BattleManager) scheduleEntityAction(entry *ActivityEntry) {
	entity := entry.Entity

	// Determine if this is a player or enemy
	isPlayer := bm.isPlayerEntity(entity)

	if isPlayer {
		// For now, auto-attack (later this will be player choice)
		target := bm.selectRandomTarget(bm.enemyParty)
		if target != nil {
			action := &BattleAction{
				Entity:     entity,
				ActionType: ActionAttack,
				Target:     target,
				Speed:      entry.Speed,
				Timestamp:  bm.battleTime,
			}
			bm.actionQueue = append(bm.actionQueue, action)

			logger.Debug("🎯 Player %s scheduled attack on %s", entity.GetID(), target.GetID())
		}
	} else {
		// Enemy AI - simple attack for now
		target := bm.selectRandomTarget(bm.playerParty)
		if target != nil {
			action := &BattleAction{
				Entity:     entity,
				ActionType: ActionAttack,
				Target:     target,
				Speed:      entry.Speed,
				Timestamp:  bm.battleTime,
			}
			bm.actionQueue = append(bm.actionQueue, action)

			logger.Debug("👹 Enemy %s scheduled attack on %s", entity.GetID(), target.GetID())
		}
	}
}

// processActionQueue executes queued actions
func (bm *BattleManager) processActionQueue() {
	for len(bm.actionQueue) > 0 {
		action := bm.actionQueue[0]
		bm.actionQueue = bm.actionQueue[1:]

		bm.executeAction(action)

		if bm.onActionExecuted != nil {
			bm.onActionExecuted(action)
		}
	}
}

// executeAction performs the actual battle action
func (bm *BattleManager) executeAction(action *BattleAction) {
	logger.Debug("⚔️  Executing action: %s attacks %s", action.Entity.GetID(), action.Target.GetID())

	// Simple damage calculation for now
	attacker := action.Entity.RPGStats()
	defender := action.Target.RPGStats()

	if attacker == nil || defender == nil {
		return
	}

	// Basic damage formula
	damage := attacker.Level*5 + 10 // Simple calculation

	// Apply damage
	defender.CurrentHP -= damage
	if defender.CurrentHP < 0 {
		defender.CurrentHP = 0
	}

	logger.Debug("💥 %s takes %d damage! HP: %d/%d",
		action.Target.GetID(), damage, defender.CurrentHP, defender.MaxHP)

	// Check if target is defeated
	if defender.CurrentHP <= 0 {
		logger.Debug("💀 %s is defeated!", action.Target.GetID())
		bm.removeFromActivityQueue(action.Target)
	}
}

// Helper methods
func (bm *BattleManager) isPlayerEntity(entity *ecs.Entity) bool {
	for _, player := range bm.playerParty {
		if player == entity {
			return true
		}
	}
	return false
}

func (bm *BattleManager) selectRandomTarget(targets []*ecs.Entity) *ecs.Entity {
	aliveTargets := make([]*ecs.Entity, 0)

	for _, target := range targets {
		if stats := target.RPGStats(); stats != nil && stats.CurrentHP > 0 {
			aliveTargets = append(aliveTargets, target)
		}
	}

	if len(aliveTargets) == 0 {
		return nil
	}

	// For now, just return the first alive target
	// In a real implementation, you'd use random selection
	return aliveTargets[0]
}

func (bm *BattleManager) removeFromActivityQueue(entity *ecs.Entity) {
	for i, entry := range bm.activityQueue {
		if entry.Entity == entity {
			bm.activityQueue = append(bm.activityQueue[:i], bm.activityQueue[i+1:]...)
			break
		}
	}
}

func (bm *BattleManager) checkBattleEndConditions() {
	// Check if all enemies are defeated
	aliveEnemies := 0
	for _, enemy := range bm.enemyParty {
		if stats := enemy.RPGStats(); stats != nil && stats.CurrentHP > 0 {
			aliveEnemies++
		}
	}

	if aliveEnemies == 0 {
		bm.endBattle(true) // Victory
		return
	}

	// Check if all players are defeated
	alivePlayers := 0
	for _, player := range bm.playerParty {
		if stats := player.RPGStats(); stats != nil && stats.CurrentHP > 0 {
			alivePlayers++
		}
	}

	if alivePlayers == 0 {
		bm.endBattle(false) // Defeat
		return
	}
}

func (bm *BattleManager) endBattle(victory bool) {
	logger.Debug("🏁 Battle ended! Victory: %t", victory)

	if victory {
		bm.state = BattleStateVictory
		logger.Debug("🎉 Victory! Battle will remain active for 3 seconds to show results")
	} else {
		bm.state = BattleStateDefeat
		logger.Debug("💀 Defeat! Battle will remain active for 3 seconds to show results")
	}

	// Don't immediately end the battle - keep it active in victory/defeat state
	// The battle will be ended after a delay or player input

	// Schedule battle end after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		bm.battleStarted = false
		logger.Debug("⏰ Battle timeout - ending battle")
		if bm.onBattleEnd != nil {
			bm.onBattleEnd(victory)
		}
	}()
}

// Getters
func (bm *BattleManager) GetState() BattleState {
	return bm.state
}

func (bm *BattleManager) IsActive() bool {
	return bm.battleStarted
}

func (bm *BattleManager) GetPlayerFormation() *Formation {
	return bm.playerFormation
}

func (bm *BattleManager) GetEnemyFormation() *Formation {
	return bm.enemyFormation
}

func (bm *BattleManager) GetActivityQueue() []*ActivityEntry {
	return bm.activityQueue
}

// Setters for callbacks
func (bm *BattleManager) SetOnBattleEnd(callback func(victory bool)) {
	bm.onBattleEnd = callback
}

func (bm *BattleManager) SetOnActionExecuted(callback func(action *BattleAction)) {
	bm.onActionExecuted = callback
}
