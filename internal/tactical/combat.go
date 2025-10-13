// Package tactical provides turn-based combat system for FFT-style gameplay
package tactical

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// TurnPhase represents the current phase of a unit's turn
type TurnPhase int

const (
	TurnPhaseMove TurnPhase = iota
	TurnPhaseAction
	TurnPhaseEnd
)

// ActionType represents different actions a unit can take
type ActionType int

const (
	ActionMove ActionType = iota
	ActionAttack
	ActionSkill
	ActionItem
	ActionWait
)

// TurnOrder represents a unit's position in the turn order
type TurnOrder struct {
	Entity       *ecs.Entity
	Initiative   int
	HasActed     bool
	CurrentPhase TurnPhase
	GridPos      GridPos
	MoveRange    int
	CanMove      bool
	CanAct       bool
}

// TacticalCombat manages FFT-style turn-based combat
type TacticalCombat struct {
	Grid         *Grid
	TurnOrder    []*TurnOrder
	CurrentTurn  int
	IsActive     bool
	ActiveUnit   *ecs.Entity
	SelectedTile GridPos
	ValidMoves   []GridPos
	ValidTargets []GridPos
	Phase        TurnPhase
}

// NewTacticalCombat creates a new tactical combat system
func NewTacticalCombat(grid *Grid) *TacticalCombat {
	return &TacticalCombat{
		Grid:        grid,
		TurnOrder:   make([]*TurnOrder, 0),
		CurrentTurn: 0,
		IsActive:    false,
		Phase:       TurnPhaseMove,
	}
}

// StartCombat initializes combat with all entities
func (tc *TacticalCombat) StartCombat(entities []*ecs.Entity) {
	tc.TurnOrder = make([]*TurnOrder, 0)

	// Calculate initiative and create turn order
	for _, entity := range entities {
		if stats := entity.RPGStats(); stats != nil {
			// Always use random positioning for tactical combat
			var gridPos GridPos
			if entity.HasTag(ecs.TagPlayer) && entity.Transform() != nil {
				// For players, convert existing position to nearest grid position and align
				transform := entity.Transform()
				offsetX, offsetY := 50.0, 112.0 // Updated to match game world Y position (110px panel + 2px separator)
				tileSize := float64(tc.Grid.TileSize)
				gridX := int((transform.X - offsetX) / tileSize)
				gridY := int((transform.Y - offsetY) / tileSize)
				gridPos = GridPos{X: gridX, Y: gridY}

				// Align player to exact grid position to fix any misalignment
				oldX, oldY := transform.X, transform.Y
				worldX, worldY := tc.Grid.GridToWorld(gridPos)
				transform.X = worldX + offsetX
				transform.Y = worldY + offsetY

				fmt.Printf("DEBUG: Combat system aligned player %s: (%.1f,%.1f) -> (%.1f,%.1f) at grid (%d,%d)\n",
					entity.GetID(), oldX, oldY, transform.X, transform.Y, gridPos.X, gridPos.Y)
			} else {
				// For enemies, always find a random position
				gridPos = tc.FindStartingPosition(entity)
				fmt.Printf("DEBUG: Combat system found random position for enemy %s: (%d,%d)\n", entity.GetID(), gridPos.X, gridPos.Y)

				// Update the entity's world position to match the grid position
				if transform := entity.Transform(); transform != nil {
					offsetX, offsetY := 50.0, 112.0 // Updated to match game world Y position (110px panel + 2px separator)
					worldX, worldY := tc.Grid.GridToWorld(gridPos)
					transform.X = worldX + offsetX
					transform.Y = worldY + offsetY

					// Debug: Verify the positioning math
					fmt.Printf("DEBUG: Enemy %s positioning - Grid: (%d,%d) -> GridToWorld: (%.1f,%.1f) -> Final: (%.1f,%.1f)\n",
						entity.GetID(), gridPos.X, gridPos.Y, worldX, worldY, transform.X, transform.Y)

					// Double-check: convert back to grid to verify alignment
					backToGridX := int((transform.X - offsetX) / float64(tc.Grid.TileSize))
					backToGridZ := int((transform.Y - offsetY) / float64(tc.Grid.TileSize))
					fmt.Printf("DEBUG: Enemy %s round-trip verification: (%.1f,%.1f) -> (%d,%d) [should match (%d,%d)]\n",
						entity.GetID(), transform.X, transform.Y, backToGridX, backToGridZ, gridPos.X, gridPos.Y)
				}
			}

			turnOrder := &TurnOrder{
				Entity:       entity,
				Initiative:   tc.calculateInitiative(stats),
				HasActed:     false,
				CurrentPhase: TurnPhaseMove,
				GridPos:      gridPos,
				MoveRange:    tc.calculateMoveRange(stats),
				CanMove:      true,
				CanAct:       true,
			}

			tc.TurnOrder = append(tc.TurnOrder, turnOrder)
			// DON'T call SetOccupied again - entities are already positioned by deployment system
		}
	}

	// Sort by initiative (highest first)
	sort.Slice(tc.TurnOrder, func(i, j int) bool {
		return tc.TurnOrder[i].Initiative > tc.TurnOrder[j].Initiative
	})

	tc.IsActive = true
	tc.CurrentTurn = 0

	// Debug: Show combat summary
	fmt.Printf("DEBUG: Tactical combat started with %d participants:\n", len(tc.TurnOrder))
	playerCount := 0
	enemyCount := 0
	for i, turn := range tc.TurnOrder {
		entityType := "Enemy"
		if turn.Entity.HasTag(ecs.TagPlayer) {
			entityType = "Player"
			playerCount++
		} else {
			enemyCount++
		}
		stats := turn.Entity.RPGStats()
		if stats != nil {
			fmt.Printf("  %d. %s %s (%s) at (%d,%d) - Initiative: %d, Move: %d\n",
				i+1, entityType, turn.Entity.GetID(), stats.Job.String(),
				turn.GridPos.X, turn.GridPos.Y, turn.Initiative, turn.MoveRange)
		}
	}
	fmt.Printf("DEBUG: Combat setup complete - %d players vs %d enemies\n", playerCount, enemyCount)

	tc.startTurn()
}

// calculateInitiative determines turn order based on character stats
func (tc *TacticalCombat) calculateInitiative(stats *components.RPGStatsComponent) int {
	// FFT-style initiative: Speed + random factor
	baseInitiative := stats.Speed
	// Add some randomness (you might want to use a proper random function)
	randomFactor := time.Now().Nanosecond() % 10
	return baseInitiative + randomFactor
}

// calculateMoveRange determines how far a unit can move
func (tc *TacticalCombat) calculateMoveRange(stats *components.RPGStatsComponent) int {
	// Base movement range based on job/speed
	baseMove := 3 // Default movement range

	switch stats.Job {
	case components.JobRogue:
		baseMove = 4 // Rogues are faster
	case components.JobMage:
		baseMove = 2 // Mages are slower
	case components.JobWarrior:
		baseMove = 3 // Balanced
	}

	return baseMove
}

// FindStartingPosition finds an available starting position for an entity
func (tc *TacticalCombat) FindStartingPosition(entity *ecs.Entity) GridPos {
	// Random placement within deployment zones

	if entity.HasTag(ecs.TagPlayer) {
		// Place players on the left side with random positions
		return tc.findRandomPosition(0, tc.Grid.Width/2, 0, tc.Grid.Height)
	} else {
		// Place enemies randomly throughout the entire grid
		return tc.findRandomPosition(0, tc.Grid.Width, 0, tc.Grid.Height)
	}
}

// findRandomPosition finds a random available position within the specified bounds
func (tc *TacticalCombat) findRandomPosition(minX, maxX, minY, maxY int) GridPos {
	// Collect all available positions in the zone
	var availablePositions []GridPos

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			pos := GridPos{X: x, Y: y}
			if tc.Grid.IsPassable(pos) {
				availablePositions = append(availablePositions, pos)
			}
		}
	}

	// If no positions available, fallback to systematic search
	if len(availablePositions) == 0 {
		fmt.Printf("DEBUG: No available positions in zone (%d-%d, %d-%d), using systematic fallback\n",
			minX, maxX, minY, maxY)
		for x := minX; x < maxX; x++ {
			for y := minY; y < maxY; y++ {
				pos := GridPos{X: x, Y: y}
				if tc.Grid.IsValidPosition(pos) {
					return pos
				}
			}
		}
		// Ultimate fallback
		return GridPos{X: 0, Y: 0}
	}

	// Return random position from available ones
	randomIndex := rand.Intn(len(availablePositions))
	selectedPos := availablePositions[randomIndex]

	fmt.Printf("DEBUG: Selected random position (%d,%d) from %d available positions\n",
		selectedPos.X, selectedPos.Y, len(availablePositions))

	return selectedPos
}

// startTurn begins a new unit's turn
func (tc *TacticalCombat) startTurn() {
	if tc.CurrentTurn >= len(tc.TurnOrder) {
		tc.nextRound()
		return
	}

	currentUnit := tc.TurnOrder[tc.CurrentTurn]
	tc.ActiveUnit = currentUnit.Entity
	tc.Phase = TurnPhaseMove
	currentUnit.CurrentPhase = TurnPhaseMove
	currentUnit.CanMove = true
	currentUnit.CanAct = true

	// Calculate valid moves for this unit
	tc.ValidMoves = tc.Grid.CalculateMovementRange(currentUnit.GridPos, currentUnit.MoveRange)
}

// nextRound starts a new combat round
func (tc *TacticalCombat) nextRound() {
	// Reset all units for new round
	for _, unit := range tc.TurnOrder {
		unit.HasActed = false
		unit.CurrentPhase = TurnPhaseMove
	}

	tc.CurrentTurn = 0
	tc.startTurn()
}

// ProcessMove handles unit movement
func (tc *TacticalCombat) ProcessMove(targetPos GridPos) bool {
	if tc.CurrentTurn >= len(tc.TurnOrder) {
		return false
	}

	currentUnit := tc.TurnOrder[tc.CurrentTurn]

	// Check if move is valid
	if !tc.isValidMove(targetPos) {
		return false
	}

	// Move the unit
	oldPos := currentUnit.GridPos
	tc.Grid.SetOccupied(oldPos, false, "")
	tc.Grid.SetOccupied(targetPos, true, currentUnit.Entity.GetID())
	currentUnit.GridPos = targetPos
	currentUnit.CanMove = false

	// Update entity transform component
	if transform := currentUnit.Entity.Transform(); transform != nil {
		offsetX, offsetY := 50.0, 112.0 // Updated to match game world Y position (110px panel + 2px separator)
		worldX, worldY := tc.Grid.GridToWorld(targetPos)
		transform.X = worldX + offsetX
		transform.Y = worldY + offsetY
	}

	// Advance to action phase
	tc.Phase = TurnPhaseAction
	currentUnit.CurrentPhase = TurnPhaseAction

	return true
}

// isValidMove checks if a move to the target position is valid
func (tc *TacticalCombat) isValidMove(targetPos GridPos) bool {
	for _, validPos := range tc.ValidMoves {
		if validPos.X == targetPos.X && validPos.Y == targetPos.Y {
			return true
		}
	}
	return false
}

// EndTurn ends the current unit's turn
func (tc *TacticalCombat) EndTurn() {
	if tc.CurrentTurn < len(tc.TurnOrder) {
		tc.TurnOrder[tc.CurrentTurn].HasActed = true
	}

	tc.CurrentTurn++
	tc.startTurn()
}

// GetCurrentUnit returns the currently active unit
func (tc *TacticalCombat) GetCurrentUnit() *ecs.Entity {
	return tc.ActiveUnit
}

// GetValidMoves returns valid movement positions for current unit
func (tc *TacticalCombat) GetValidMoves() []GridPos {
	return tc.ValidMoves
}

// IsPositionOccupied checks if a grid position has a unit
func (tc *TacticalCombat) IsPositionOccupied(pos GridPos) bool {
	tile := tc.Grid.GetTile(pos)
	return tile != nil && tile.Occupied
}

// GetUnitAtPosition returns the entity at a grid position
func (tc *TacticalCombat) GetUnitAtPosition(pos GridPos) *ecs.Entity {
	tile := tc.Grid.GetTile(pos)
	if tile == nil || !tile.Occupied {
		return nil
	}

	// Find entity by ID
	for _, unit := range tc.TurnOrder {
		if unit.Entity.GetID() == tile.UnitID {
			return unit.Entity
		}
	}

	return nil
}
