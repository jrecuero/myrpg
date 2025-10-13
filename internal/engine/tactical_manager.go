// Package engine provides game mode management for switching between exploration and tactical combat
package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/tactical"
)

// GameMode represents different gameplay modes
type GameMode int

const (
	ModeExploration GameMode = iota // Free movement exploration
	ModeTactical                    // Grid-based tactical combat
)

func (gm GameMode) String() string {
	switch gm {
	case ModeExploration:
		return "Exploration"
	case ModeTactical:
		return "Tactical"
	default:
		return "Unknown"
	}
}

// TacticalManager handles tactical combat mode
type TacticalManager struct {
	Grid               *tactical.Grid
	GridRenderer       *tactical.GridRenderer
	Combat             *tactical.TacticalCombat         // Legacy combat system
	TurnBasedCombat    *tactical.TurnBasedCombatManager // New turn-based combat system
	IsActive           bool
	Participants       []*ecs.Entity // Entities involved in tactical combat
	UseTurnBasedCombat bool          // Flag to switch between old and new combat systems
}

// NewTacticalManager creates a new tactical manager
func NewTacticalManager(gridWidth, gridHeight, tileSize int) *TacticalManager {
	grid := tactical.NewGrid(gridWidth, gridHeight, tileSize)
	gridRenderer := tactical.NewGridRenderer(grid)
	combat := tactical.NewTacticalCombat(grid)
	turnBasedCombat := tactical.NewTurnBasedCombatManager(grid)

	return &TacticalManager{
		Grid:               grid,
		GridRenderer:       gridRenderer,
		Combat:             combat,
		TurnBasedCombat:    turnBasedCombat,
		IsActive:           false,
		Participants:       make([]*ecs.Entity, 0),
		UseTurnBasedCombat: true, // Enable new turn-based system by default
	}
}

// StartTacticalCombat switches to tactical mode with given entities
func (tm *TacticalManager) StartTacticalCombat(entities []*ecs.Entity) {
	tm.IsActive = true
	tm.Participants = entities

	if tm.UseTurnBasedCombat {
		// Initialize new turn-based combat system
		if err := tm.TurnBasedCombat.InitializeCombat(entities); err != nil {
			fmt.Printf("Failed to initialize turn-based combat: %v\n", err)
			// Fall back to old system
			tm.UseTurnBasedCombat = false
			tm.Combat.StartCombat(entities)
		}
	} else {
		// Use legacy combat system
		tm.Combat.StartCombat(entities)
	}

	// Clear any previous highlights
	tm.GridRenderer.ClearHighlights()
}

// EndTacticalCombat returns to exploration mode
func (tm *TacticalManager) EndTacticalCombat() {
	tm.IsActive = false
	tm.Participants = make([]*ecs.Entity, 0)
	tm.GridRenderer.ClearHighlights()
	tm.GridRenderer.SetShowGrid(false)
}

// Update handles tactical combat updates
func (tm *TacticalManager) Update() {
	if !tm.IsActive {
		return
	}

	if tm.UseTurnBasedCombat {
		// Update new turn-based combat system
		if err := tm.TurnBasedCombat.Update(); err != nil {
			fmt.Printf("Turn-based combat update error: %v\n", err)
		}

		// Check if combat has ended
		if !tm.TurnBasedCombat.IsActive {
			tm.handleCombatEnd()
		}
	} else {
		// Update legacy combat system
		// This will be expanded as we add more tactical features
	}
}

// DrawGrid renders the tactical grid overlay
func (tm *TacticalManager) DrawGrid(screen *ebiten.Image, offsetX, offsetY float64) {
	if tm.IsActive {
		tm.GridRenderer.Draw(screen, offsetX, offsetY)
	}
}

// GetTileAtScreenPos returns grid position for screen coordinates
func (tm *TacticalManager) GetTileAtScreenPos(screenX, screenY, offsetX, offsetY float64) (tactical.GridPos, bool) {
	if tm.IsActive {
		return tm.GridRenderer.GetTileAtScreenPos(screenX, screenY, offsetX, offsetY)
	}
	return tactical.GridPos{}, false
}

// HighlightMovementRange highlights valid movement tiles for current unit
func (tm *TacticalManager) HighlightMovementRange() {
	if !tm.IsActive {
		return
	}

	// Clear previous highlights
	tm.GridRenderer.ClearHighlightType(tactical.HighlightMovement)

	// Get valid moves from combat system
	validMoves := tm.Combat.GetValidMoves()
	tm.GridRenderer.HighlightTiles(validMoves, tactical.HighlightMovement)
}

// HighlightMovementRangeForPlayer highlights movement range for a specific player
func (tm *TacticalManager) HighlightMovementRangeForPlayer(player *ecs.Entity) {
	if !tm.IsActive || player == nil {
		return
	}

	// Clear previous highlights
	tm.GridRenderer.ClearHighlightType(tactical.HighlightMovement)

	// Get player's current position
	transform := player.Transform()
	if transform == nil {
		return
	}

	// Get player's RPG stats to determine movement range
	stats := player.RPGStats()
	if stats == nil {
		return
	}

	// Convert world coordinates to grid position
	offsetX, offsetY := constants.GridOffsetX, constants.GridOffsetY // Updated to match game world Y position (110px panel + 2px separator)
	tileSize := float64(tm.Grid.TileSize)
	gridX := int((transform.X - offsetX) / tileSize)
	gridY := int((transform.Y - offsetY) / tileSize)
	currentPos := tactical.GridPos{X: gridX, Y: gridY}

	// Use player's remaining movement from RPG stats
	moveRange := stats.MovesRemaining
	validMoves := tm.Grid.CalculateMovementRange(currentPos, moveRange)

	fmt.Printf("DEBUG: Highlighting movement range for player %s (%s) at (%d,%d) with %d moves remaining (max: %d)\n",
		player.GetID(), stats.Job.String(), currentPos.X, currentPos.Y, moveRange, stats.MoveRange)

	tm.GridRenderer.HighlightTiles(validMoves, tactical.HighlightMovement)
}

// SelectTile highlights a tile as selected
func (tm *TacticalManager) SelectTile(pos tactical.GridPos) {
	tm.GridRenderer.ClearHighlightType(tactical.HighlightSelected)
	tm.GridRenderer.HighlightTile(pos, tactical.HighlightSelected)
}

// ResetAllMovement resets movement for all participants
func (tm *TacticalManager) ResetAllMovement() {
	for _, entity := range tm.Participants {
		stats := entity.RPGStats()
		if stats != nil {
			stats.ResetMovement()
			fmt.Printf("DEBUG: Reset movement for %s (%s) - %d moves available\n",
				entity.GetID(), stats.Job.String(), stats.MoveRange)
		}
	}
}

// ResetPlayerMovement resets movement for a specific player (useful for turn-based systems)
func (tm *TacticalManager) ResetPlayerMovement(player *ecs.Entity) {
	if stats := player.RPGStats(); stats != nil {
		stats.ResetMovement()
		fmt.Printf("DEBUG: Reset movement for player %s (%s) - %d moves available\n",
			player.GetID(), stats.Job.String(), stats.MoveRange)
	}
}

// handleCombatEnd processes the end of combat
func (tm *TacticalManager) handleCombatEnd() {
	if tm.UseTurnBasedCombat {
		result := tm.TurnBasedCombat.GetResult()
		fmt.Printf("Combat ended with result: %s\n", result.String())

		// TODO: Add proper victory/defeat handling
		// This could trigger UI changes, experience gain, loot, etc.
	}

	// End tactical mode
	tm.EndTacticalCombat()
}

// GetTurnBasedCombat returns the turn-based combat manager for external access
func (tm *TacticalManager) GetTurnBasedCombat() *tactical.TurnBasedCombatManager {
	return tm.TurnBasedCombat
}

// IsUsingTurnBasedCombat returns true if using the new turn-based combat system
func (tm *TacticalManager) IsUsingTurnBasedCombat() bool {
	return tm.UseTurnBasedCombat
}
