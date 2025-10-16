// Package components provides combat-specific ECS components for the turn-based battle system
package components

// ActionPointsComponent manages action points for turn-based combat
type ActionPointsComponent struct {
	Current int // Current action points available this turn
	Maximum int // Maximum action points per turn (based on job class)
}

// NewActionPointsComponent creates a new action points component
func NewActionPointsComponent(maxAP int) *ActionPointsComponent {
	return &ActionPointsComponent{
		Current: maxAP,
		Maximum: maxAP,
	}
}

// CanAfford checks if the unit has enough AP for an action
func (ap *ActionPointsComponent) CanAfford(cost int) bool {
	return ap.Current >= cost
}

// Spend reduces current AP by the specified amount
func (ap *ActionPointsComponent) Spend(cost int) bool {
	if ap.CanAfford(cost) {
		ap.Current -= cost
		return true
	}
	return false
}

// Restore sets current AP back to maximum (start of turn)
func (ap *ActionPointsComponent) Restore() {
	ap.Current = ap.Maximum
}

// IsExhausted returns true if no AP remaining
func (ap *ActionPointsComponent) IsExhausted() bool {
	return ap.Current <= 0
}

// CombatStateComponent tracks combat-specific state for units
type CombatStateComponent struct {
	HasActed   bool // Has this unit completed its turn
	IsActive   bool // Is this unit currently taking their turn
	Team       Team // Which team does this unit belong to
	Initiative int  // Initiative value for turn order
	CanAct     bool // Can this unit still act (not stunned, etc.)
}

// Team represents which side a unit fights for
type Team int

const (
	TeamPlayer Team = iota
	TeamEnemy
)

func (t Team) String() string {
	switch t {
	case TeamPlayer:
		return "Player"
	case TeamEnemy:
		return "Enemy"
	default:
		return "Unknown"
	}
}

// NewCombatStateComponent creates a new combat state component
func NewCombatStateComponent(team Team, initiative int) *CombatStateComponent {
	return &CombatStateComponent{
		HasActed:   false,
		IsActive:   false,
		Team:       team,
		Initiative: initiative,
		CanAct:     true,
	}
}

// StartTurn marks this unit as active and resets turn flags
func (cs *CombatStateComponent) StartTurn() {
	cs.IsActive = true
	cs.HasActed = false
}

// EndTurn marks this unit as having completed their turn
func (cs *CombatStateComponent) EndTurn() {
	cs.IsActive = false
	cs.HasActed = true
}

// Reset prepares the unit for a new combat round
func (cs *CombatStateComponent) Reset() {
	cs.HasActed = false
	cs.IsActive = false
	cs.CanAct = true
}
