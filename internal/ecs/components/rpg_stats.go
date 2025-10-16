package components

import "fmt"

// MoveRecord represents a single movement for undo functionality
type MoveRecord struct {
	FromX, FromZ int // Starting position
	ToX, ToZ     int // Ending position
	Distance     int // Movement cost
}

// JobType represents different character classes/jobs in the RPG
type JobType int

const (
	JobWarrior JobType = iota // High HP and Defense, moderate Attack
	JobMage                   // High Magic Power, low HP and Defense
	JobRogue                  // High Speed and Crit, moderate HP
	JobCleric                 // High Magic Defense and Healing, moderate HP
	JobArcher                 // High Range Attack and Speed, low Defense
)

const (
	DefaultRPGStateSpeed    = 50
	DefaultRPGStateAccuracy = 85
	DefaultRPGStatCritRate  = 5
)

// String returns the string representation of a JobType
func (j JobType) String() string {
	switch j {
	case JobWarrior:
		return "Warrior"
	case JobMage:
		return "Mage"
	case JobRogue:
		return "Rogue"
	case JobCleric:
		return "Cleric"
	case JobArcher:
		return "Archer"
	default:
		return "Unknown"
	}
}

// RPGStatsComponent represents the RPG statistics of a character entity.
// It includes core stats like health, level, experience, and combat attributes.
type RPGStatsComponent struct {
	// Core Stats
	Level      int // Current character level (1-99)
	Experience int // Current experience points
	ExpToNext  int // Experience needed for next level

	// Health and Mana
	CurrentHP int // Current health points
	MaxHP     int // Maximum health points
	CurrentMP int // Current mana/magic points
	MaxMP     int // Maximum mana/magic points

	// Combat Stats
	Attack       int // Physical attack power
	Defense      int // Physical defense
	MagicAttack  int // Magic attack power
	MagicDefense int // Magic defense
	Speed        int // Speed/agility stat
	Accuracy     int // Hit accuracy percentage (0-100)
	CritRate     int // Critical hit rate percentage (0-100)

	// Tactical Stats
	MoveRange      int          // Movement range in tiles per turn
	MovesRemaining int          // Remaining moves this turn
	MoveHistory    []MoveRecord // History of moves this turn for undo functionality

	// Character Info
	Job  JobType // Character class/job
	Name string  // Character display name
}

// NewRPGStatsComponent creates a new RPG stats component with default values for a specific job
func NewRPGStatsComponent(name string, job JobType, level int) *RPGStatsComponent {
	stats := &RPGStatsComponent{
		Name:       name,
		Job:        job,
		Level:      level,
		Experience: 0,
		ExpToNext:  level * 100, // Simple EXP formula
		Speed:      DefaultRPGStateSpeed,
		Accuracy:   DefaultRPGStateAccuracy,
		CritRate:   DefaultRPGStatCritRate,
	}

	// Apply job-specific base stats
	switch job {
	case JobWarrior:
		stats.MaxHP = 120 + (level-1)*15
		stats.MaxMP = 30 + (level-1)*3
		stats.Attack = 18 + (level-1)*3
		stats.Defense = 15 + (level-1)*2
		stats.MagicAttack = 8 + (level-1)*1
		stats.MagicDefense = 12 + (level-1)*2
		stats.Speed = 40 + (level-1)*1
		stats.MoveRange = 3 // Warriors are slower but steady

	case JobMage:
		stats.MaxHP = 80 + (level-1)*8
		stats.MaxMP = 100 + (level-1)*12
		stats.Attack = 10 + (level-1)*1
		stats.Defense = 8 + (level-1)*1
		stats.MagicAttack = 25 + (level-1)*4
		stats.MagicDefense = 20 + (level-1)*3
		stats.Speed = 45 + (level-1)*2
		stats.MoveRange = 2 // Mages are slow and fragile

	case JobRogue:
		stats.MaxHP = 95 + (level-1)*10
		stats.MaxMP = 50 + (level-1)*5
		stats.Attack = 20 + (level-1)*3
		stats.Defense = 10 + (level-1)*1
		stats.MagicAttack = 12 + (level-1)*1
		stats.MagicDefense = 10 + (level-1)*1
		stats.Speed = 65 + (level-1)*3
		stats.CritRate = 15 + (level-1)/2
		stats.MoveRange = 5 // Rogues are very mobile

	case JobCleric:
		stats.MaxHP = 100 + (level-1)*12
		stats.MaxMP = 90 + (level-1)*10
		stats.Attack = 12 + (level-1)*1
		stats.Defense = 12 + (level-1)*2
		stats.MagicAttack = 20 + (level-1)*3
		stats.MagicDefense = 25 + (level-1)*4
		stats.Speed = 35 + (level-1)*1
		stats.MoveRange = 3 // Clerics have moderate mobility

	case JobArcher:
		stats.MaxHP = 85 + (level-1)*9
		stats.MaxMP = 40 + (level-1)*4
		stats.Attack = 22 + (level-1)*3
		stats.Defense = 9 + (level-1)*1
		stats.MagicAttack = 10 + (level-1)*1
		stats.MagicDefense = 8 + (level-1)*1
		stats.Speed = 60 + (level-1)*3
		stats.Accuracy = 95 + (level-1)/5
		stats.MoveRange = 4 // Archers need positioning flexibility
	}

	// Set current HP/MP to max
	stats.CurrentHP = stats.MaxHP
	stats.CurrentMP = stats.MaxMP

	// Initialize tactical stats
	stats.MovesRemaining = stats.MoveRange
	stats.MoveHistory = make([]MoveRecord, 0)

	return stats
}

// IsAlive returns true if the character has HP remaining
func (r *RPGStatsComponent) IsAlive() bool {
	return r.CurrentHP > 0
}

// TakeDamage reduces current HP by the specified amount
func (r *RPGStatsComponent) TakeDamage(damage int) {
	r.CurrentHP -= damage
	if r.CurrentHP < 0 {
		r.CurrentHP = 0
	}
}

// Heal increases current HP by the specified amount, not exceeding max HP
func (r *RPGStatsComponent) Heal(amount int) {
	r.CurrentHP += amount
	if r.CurrentHP > r.MaxHP {
		r.CurrentHP = r.MaxHP
	}
}

// UseMana reduces current MP by the specified amount
func (r *RPGStatsComponent) UseMana(amount int) bool {
	if r.CurrentMP >= amount {
		r.CurrentMP -= amount
		return true
	}
	return false
}

// RestoreMana increases current MP by the specified amount, not exceeding max MP
func (r *RPGStatsComponent) RestoreMana(amount int) {
	r.CurrentMP += amount
	if r.CurrentMP > r.MaxMP {
		r.CurrentMP = r.MaxMP
	}
}

// GainExperience adds experience and handles level ups
func (r *RPGStatsComponent) GainExperience(exp int) bool {
	r.Experience += exp
	if r.Experience >= r.ExpToNext {
		r.LevelUp()
		return true // Level up occurred
	}
	return false
}

// LevelUp increases the character's level and recalculates stats
func (r *RPGStatsComponent) LevelUp() {
	r.Level++
	r.Experience -= r.ExpToNext
	r.ExpToNext = r.Level * 100

	// Recalculate stats based on new level (this is a simple approach)
	// In a real game, you might want more sophisticated stat growth
	oldMaxHP := r.MaxHP
	oldMaxMP := r.MaxMP

	// Create a new stats component with updated level to get new max values
	newStats := NewRPGStatsComponent(r.Name, r.Job, r.Level)

	// Update max stats
	r.MaxHP = newStats.MaxHP
	r.MaxMP = newStats.MaxMP
	r.Attack = newStats.Attack
	r.Defense = newStats.Defense
	r.MagicAttack = newStats.MagicAttack
	r.MagicDefense = newStats.MagicDefense
	r.Speed = newStats.Speed

	// Restore HP/MP proportionally or fully heal on level up
	r.CurrentHP += (r.MaxHP - oldMaxHP) // Add the HP increase
	r.CurrentMP += (r.MaxMP - oldMaxMP) // Add the MP increase
}

// ConsumeMovement reduces remaining movement by the specified distance
func (r *RPGStatsComponent) ConsumeMovement(distance int) bool {
	if r.MovesRemaining >= distance {
		r.MovesRemaining -= distance
		return true
	}
	return false
}

// ResetMovement restores movement for a new turn
func (r *RPGStatsComponent) ResetMovement() {
	r.MovesRemaining = r.MoveRange
	r.MoveHistory = make([]MoveRecord, 0) // Clear movement history
}

// CanMove checks if the unit can move the specified distance
func (r *RPGStatsComponent) CanMove(distance int) bool {
	return r.MovesRemaining >= distance
}

// RecordMove records a movement for potential undo
func (r *RPGStatsComponent) RecordMove(fromX, fromZ, toX, toZ, distance int) {
	move := MoveRecord{
		FromX:    fromX,
		FromZ:    fromZ,
		ToX:      toX,
		ToZ:      toZ,
		Distance: distance,
	}
	r.MoveHistory = append(r.MoveHistory, move)
}

// IsUndoMove checks if moving to this position would be an undo (without performing it)
func (r *RPGStatsComponent) IsUndoMove(toX, toZ int) bool {
	// Check if this position matches any previous position in our move history
	for i := len(r.MoveHistory) - 1; i >= 0; i-- {
		move := r.MoveHistory[i]
		if move.FromX == toX && move.FromZ == toZ {
			return true
		}
	}
	return false
}

// TryUndoMove attempts to recover movement if moving back to a previous position
func (r *RPGStatsComponent) TryUndoMove(toX, toZ int) (bool, int) {
	// Check if this position matches any previous position in our move history
	for i := len(r.MoveHistory) - 1; i >= 0; i-- {
		move := r.MoveHistory[i]
		if move.FromX == toX && move.FromZ == toZ {
			// Moving back to a previous position - calculate recovery
			movesToRecover := 0

			// Recover moves from this point forward
			for j := i; j < len(r.MoveHistory); j++ {
				movesToRecover += r.MoveHistory[j].Distance
			}

			// Truncate history to this point
			r.MoveHistory = r.MoveHistory[:i]

			// Restore movement points
			r.MovesRemaining += movesToRecover
			if r.MovesRemaining > r.MoveRange {
				r.MovesRemaining = r.MoveRange
			}

			return true, movesToRecover
		}
	}
	return false, 0
} // GetMoveHistoryString returns a debug string of the movement history
func (r *RPGStatsComponent) GetMoveHistoryString() string {
	if len(r.MoveHistory) == 0 {
		return "No moves"
	}

	result := "Moves: "
	for i, move := range r.MoveHistory {
		if i > 0 {
			result += " -> "
		}
		result += fmt.Sprintf("(%d,%d)", move.ToX, move.ToZ)
	}
	return result
}
