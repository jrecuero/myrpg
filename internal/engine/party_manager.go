// Package engine provides party management for exploration and tactical modes
package engine

import (
	"github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/logger"
	"github.com/jrecuero/myrpg/internal/tactical"
)

// PartyManager handles party composition and deployment
type PartyManager struct {
	PartyLeader    *ecs.Entity   // Single character in exploration mode
	PartyMembers   []*ecs.Entity // Full party for tactical mode
	ReserveMembers []*ecs.Entity // Inactive party members
	MaxPartySize   int           // Maximum party size for tactical
}

// NewPartyManager creates a new party manager
func NewPartyManager(maxSize int) *PartyManager {
	return &PartyManager{
		PartyMembers:   make([]*ecs.Entity, 0),
		ReserveMembers: make([]*ecs.Entity, 0),
		MaxPartySize:   maxSize,
	}
}

// SetPartyLeader designates the exploration mode representative
func (pm *PartyManager) SetPartyLeader(leader *ecs.Entity) {
	pm.PartyLeader = leader

	// Ensure leader is in party members if not already
	if !pm.isInParty(leader) {
		pm.AddPartyMember(leader)
	}
}

// AddPartyMember adds a character to the tactical party
func (pm *PartyManager) AddPartyMember(member *ecs.Entity) bool {
	if len(pm.PartyMembers) >= pm.MaxPartySize {
		return false // Party full
	}

	if pm.isInParty(member) {
		return false // Already in party
	}

	pm.PartyMembers = append(pm.PartyMembers, member)
	logger.Debug("➕ Added %s to party (now %d members)", member.GetID(), len(pm.PartyMembers))

	// Debug: Show current party order
	logger.Debug("🏃‍♂️ Current party order after adding %s:", member.GetID())
	for i, m := range pm.PartyMembers {
		logger.Debug("  [%d] %s", i, m.GetID())
	}

	return true
}

// RemovePartyMember removes a character from the party
func (pm *PartyManager) RemovePartyMember(member *ecs.Entity) bool {
	for i, partyMember := range pm.PartyMembers {
		if partyMember == member {
			// Remove from slice
			pm.PartyMembers = append(pm.PartyMembers[:i], pm.PartyMembers[i+1:]...)

			// If this was the party leader, assign new leader
			if pm.PartyLeader == member && len(pm.PartyMembers) > 0 {
				pm.PartyLeader = pm.PartyMembers[0]
			}

			return true
		}
	}
	return false
}

// isInParty checks if an entity is already in the party
func (pm *PartyManager) isInParty(entity *ecs.Entity) bool {
	for _, member := range pm.PartyMembers {
		if member == entity {
			return true
		}
	}
	return false
}

// GetPartyForTactical returns all party members for tactical combat
func (pm *PartyManager) GetPartyForTactical() []*ecs.Entity {
	return append([]*ecs.Entity{}, pm.PartyMembers...) // Return copy in original order
}

// GetPartyForClassic returns all party members for classic battles in consistent order
func (pm *PartyManager) GetPartyForClassic() []*ecs.Entity {
	// Debug: Show original party order for classic battles
	logger.Debug("🗡️  Classic battle party order:")
	for i, member := range pm.PartyMembers {
		if stats := member.RPGStats(); stats != nil {
			logger.Debug("   [%d] %s (ID: %s)", i, stats.Name, member.GetID())
		} else {
			logger.Debug("   [%d] %s (no stats)", i, member.GetID())
		}
	}

	// For classic battles, return exact original order (no sorting)
	return append([]*ecs.Entity{}, pm.PartyMembers...)
}

// GetPartyForTacticalDeployment returns party members sorted for consistent deployment
func (pm *PartyManager) GetPartyForTacticalDeployment() []*ecs.Entity {
	// Debug: Show original party order
	logger.Debug("🔍 Original PartyMembers order:")
	for i, member := range pm.PartyMembers {
		logger.Debug("  Original[%d]: %s", i, member.GetID())
	}

	// Create a copy and sort by ID to ensure consistent ordering
	partyCopy := append([]*ecs.Entity{}, pm.PartyMembers...)

	// Sort by entity ID to ensure deterministic positioning
	for i := 0; i < len(partyCopy)-1; i++ {
		for j := i + 1; j < len(partyCopy); j++ {
			if partyCopy[i].GetID() > partyCopy[j].GetID() {
				partyCopy[i], partyCopy[j] = partyCopy[j], partyCopy[i]
			}
		}
	}

	// Debug: Show sorted order
	logger.Debug("🔍 Sorted PartyMembers for deployment:")
	for i, member := range partyCopy {
		logger.Debug("  Sorted[%d]: %s", i, member.GetID())
	}

	return partyCopy
}

// GetPartyLeader returns the current party leader
func (pm *PartyManager) GetPartyLeader() *ecs.Entity {
	return pm.PartyLeader
}

// GetPartySize returns the current party size
func (pm *PartyManager) GetPartySize() int {
	return len(pm.PartyMembers)
}

// UpdatePartyLeaderPosition updates the party leader's position for tracking
func (pm *PartyManager) UpdatePartyLeaderPosition(x, y float64) {
	if pm.PartyLeader != nil && pm.PartyLeader.Transform() != nil {
		transform := pm.PartyLeader.Transform()
		transform.X = x
		transform.Y = y
	}
}

// EnemyGroupManager handles enemy group formation for tactical combat
type EnemyGroupManager struct {
	GroupRange float64 // Distance to include enemies in tactical combat
}

// NewEnemyGroupManager creates a new enemy group manager
func NewEnemyGroupManager(groupRange float64) *EnemyGroupManager {
	return &EnemyGroupManager{
		GroupRange: groupRange,
	}
}

// FormEnemyGroup creates an enemy group for tactical combat
func (egm *EnemyGroupManager) FormEnemyGroup(triggerEnemy *ecs.Entity, allEntities []*ecs.Entity) []*ecs.Entity {
	enemyGroup := []*ecs.Entity{triggerEnemy}

	triggerT := triggerEnemy.Transform()
	if triggerT == nil {
		return enemyGroup
	}

	// Find nearby enemies to include in the group
	for _, entity := range allEntities {
		// Skip the trigger enemy (already included)
		if entity == triggerEnemy {
			continue
		}

		// Only include enemies (entities with RPG stats but no Player tag)
		if entity.RPGStats() == nil || entity.HasTag(ecs.TagPlayer) {
			continue
		}

		entityT := entity.Transform()
		if entityT == nil {
			continue
		}

		// Calculate distance to trigger enemy
		dx := triggerT.X - entityT.X
		dy := triggerT.Y - entityT.Y
		distance := dx*dx + dy*dy // Square distance

		if distance <= egm.GroupRange*egm.GroupRange {
			enemyGroup = append(enemyGroup, entity)
		}
	}

	return enemyGroup
}

// TacticalDeployment handles positioning units at tactical combat start
type TacticalDeployment struct {
	Grid       *tactical.Grid
	PlayerZone DeploymentZone
	EnemyZone  DeploymentZone
}

// DeploymentZone defines where units can be placed
type DeploymentZone struct {
	StartX, StartY int // Starting grid coordinates
	Width, Height  int // Zone dimensions
}

// NewTacticalDeployment creates a new deployment manager
func NewTacticalDeployment(grid *tactical.Grid) *TacticalDeployment {
	return &TacticalDeployment{
		Grid: grid,
		PlayerZone: DeploymentZone{
			StartX: 0,
			StartY: 0,
			Width:  grid.Width / constants.PlayerZoneWidth, // Left third of map
			Height: grid.Height,
		},
		EnemyZone: DeploymentZone{
			StartX: (grid.Width * constants.EnemyZoneStart) / constants.PlayerZoneWidth, // Right third of map
			StartY: 0,
			Width:  grid.Width / constants.PlayerZoneWidth,
			Height: grid.Height,
		},
	}
}

// DeployParty positions party members in the player zone
func (td *TacticalDeployment) DeployParty(party []*ecs.Entity) map[*ecs.Entity]tactical.GridPos {
	positions := make(map[*ecs.Entity]tactical.GridPos)

	// Debug: Show party order before deployment
	logger.Debug("🏃‍♂️ Deploying %d party members in order:", len(party))
	for i, member := range party {
		logger.Debug("  [%d] %s", i, member.GetID())
	}

	deployedCount := 0
	for _, member := range party {
		if deployedCount >= td.PlayerZone.Width*td.PlayerZone.Height {
			break // No more space
		}

		// Calculate grid position
		row := deployedCount / td.PlayerZone.Width
		col := deployedCount % td.PlayerZone.Width

		gridPos := tactical.GridPos{
			X: td.PlayerZone.StartX + col,
			Y: td.PlayerZone.StartY + row,
		}

		// Ensure position is valid and not occupied
		if td.Grid.IsValidPosition(gridPos) && td.Grid.IsPassable(gridPos) {
			// First, clear any existing occupancy for this unit (prevent double occupancy)
			td.clearUnitFromGrid(member.GetID())

			positions[member] = gridPos
			td.Grid.SetOccupied(gridPos, true, member.GetID())

			// Debug: Log deployment with detailed info
			logger.Debug("🚀 DEPLOYED: %s at Grid(%d,%d) [count=%d, row=%d, col=%d]",
				member.GetID(), gridPos.X, gridPos.Y, deployedCount, row, col)

			// Update entity transform to match grid position with offset
			if transform := member.Transform(); transform != nil {
				worldX, worldY := td.Grid.GridToWorld(gridPos)
				// Add the grid offset (same as used in DrawGrid) - Updated to match game world Y position
				transform.X = worldX + constants.GridOffsetX
				transform.Y = worldY + constants.GridOffsetY // 110px top panel + 2px separator

				// Debug: Log world coordinates with more detail
				logger.Debug("📍 Unit %s world coords: (%.1f,%.1f) from grid (%d,%d)",
					member.GetID(), transform.X, transform.Y, gridPos.X, gridPos.Y)
			}

			deployedCount++
		} else {
			// Debug: Log failed deployment
			logger.Debug("Failed to deploy unit %s at grid pos (%d,%d) - Valid: %t, Passable: %t",
				member.GetID(), gridPos.X, gridPos.Y, td.Grid.IsValidPosition(gridPos), td.Grid.IsPassable(gridPos))
		}
	}

	return positions
}

// clearUnitFromGrid removes a unit from all grid positions
func (td *TacticalDeployment) clearUnitFromGrid(unitID string) {
	for x := 0; x < td.Grid.Width; x++ {
		for y := 0; y < td.Grid.Height; y++ {
			pos := tactical.GridPos{X: x, Y: y}
			tile := td.Grid.GetTile(pos)
			if tile != nil && tile.Occupied && tile.UnitID == unitID {
				td.Grid.SetOccupied(pos, false, "")
				logger.Debug("Cleared unit %s from grid pos (%d,%d)", unitID, x, y)
			}
		}
	}
}

// DeployEnemies positions enemy entities in the enemy zone
func (td *TacticalDeployment) DeployEnemies(enemies []*ecs.Entity) map[*ecs.Entity]tactical.GridPos {
	positions := make(map[*ecs.Entity]tactical.GridPos)

	deployedCount := 0
	for _, enemy := range enemies {
		if deployedCount >= td.EnemyZone.Width*td.EnemyZone.Height {
			break // No more space
		}

		// Calculate grid position
		row := deployedCount / td.EnemyZone.Width
		col := deployedCount % td.EnemyZone.Width

		gridPos := tactical.GridPos{
			X: td.EnemyZone.StartX + col,
			Y: td.EnemyZone.StartY + row,
		}

		// Ensure position is valid and not occupied
		if td.Grid.IsValidPosition(gridPos) && td.Grid.IsPassable(gridPos) {
			// First, clear any existing occupancy for this unit (prevent double occupancy)
			td.clearUnitFromGrid(enemy.GetID())

			positions[enemy] = gridPos
			td.Grid.SetOccupied(gridPos, true, enemy.GetID())

			// Update entity transform to match grid position with offset
			if transform := enemy.Transform(); transform != nil {
				worldX, worldY := td.Grid.GridToWorld(gridPos)
				// Add the grid offset (same as used in DrawGrid) - Updated to match game world Y position
				transform.X = worldX + constants.GridOffsetX
				transform.Y = worldY + constants.GridOffsetY // 110px top panel + 2px separator
			}

			deployedCount++
		}
	}

	return positions
}
