// Package classic implements formation positioning for Dragon Quest-style battles
package classic

import (
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/logger"
)

// Position represents a 2D position in the battle formation
type Position struct {
	X, Y  float64
	Row   int
	Index int
}

// Formation manages the positioning of entities in battle
type Formation struct {
	entities  []*ecs.Entity
	positions map[*ecs.Entity]Position
	rows      int
	maxPerRow int
	isEnemy   bool

	// Screen positioning
	screenX      float64
	screenY      float64
	entityWidth  float64
	entityHeight float64
	spacing      float64
	rowSpacing   float64
}

// NewPlayerFormation creates a formation for the player party
func NewPlayerFormation(players []*ecs.Entity) *Formation {
	formation := &Formation{
		entities:     players,
		positions:    make(map[*ecs.Entity]Position),
		rows:         2,
		maxPerRow:    2, // 2x2 formation for players
		isEnemy:      false,
		screenX:      100, // Left side for players
		screenY:      400, // Bottom of screen
		entityWidth:  64,
		entityHeight: 64,
		spacing:      80, // Space between entities in same row
		rowSpacing:   80, // Space between rows
	}

	formation.arrangeEntities()
	return formation
}

// NewEnemyFormation creates a formation for the enemy party
func NewEnemyFormation(enemies []*ecs.Entity) *Formation {
	formation := &Formation{
		entities:     enemies,
		positions:    make(map[*ecs.Entity]Position),
		rows:         2,
		maxPerRow:    3, // Up to 3 enemies per row
		isEnemy:      true,
		screenX:      400, // Right side for enemies
		screenY:      150, // Top area
		entityWidth:  64,
		entityHeight: 64,
		spacing:      70, // Slightly tighter spacing for enemies
		rowSpacing:   70,
	}

	formation.arrangeEntities()
	return formation
}

// arrangeEntities positions all entities in the formation
func (f *Formation) arrangeEntities() {
	logger.Debug("📐 Arranging %s formation with %d entities",
		f.getFormationType(), len(f.entities))

	entityIndex := 0

	for row := 0; row < f.rows && entityIndex < len(f.entities); row++ {
		entitiesInThisRow := f.getEntitiesForRow(row, entityIndex)

		// Center entities in the row
		startX := f.screenX - (float64(entitiesInThisRow-1) * f.spacing / 2)
		rowY := f.screenY + float64(row)*f.rowSpacing

		// If this is enemy formation, reverse Y direction (enemies at top)
		if f.isEnemy {
			rowY = f.screenY - float64(row)*f.rowSpacing
		}

		// Position entities in this row
		for i := 0; i < entitiesInThisRow && entityIndex < len(f.entities); i++ {
			entity := f.entities[entityIndex]

			x := startX + float64(i)*f.spacing
			y := rowY

			position := Position{
				X:     x,
				Y:     y,
				Row:   row,
				Index: i,
			}

			f.positions[entity] = position

			// Update entity transform if it has one
			if transform := entity.Transform(); transform != nil {
				transform.X = x
				transform.Y = y
			}

			logger.Debug("   %s positioned at (%.1f, %.1f) - Row %d, Index %d",
				entity.GetID(), x, y, row, i)

			entityIndex++
		}
	}
}

// getEntitiesForRow calculates how many entities should be in a specific row
func (f *Formation) getEntitiesForRow(row, startIndex int) int {
	remainingEntities := len(f.entities) - startIndex
	remainingRows := f.rows - row

	if remainingRows <= 0 {
		return remainingEntities
	}

	// Try to distribute entities evenly across remaining rows
	entitiesPerRow := (remainingEntities + remainingRows - 1) / remainingRows // Ceiling division

	if entitiesPerRow > f.maxPerRow {
		entitiesPerRow = f.maxPerRow
	}

	return entitiesPerRow
}

// getFormationType returns a string describing the formation type
func (f *Formation) getFormationType() string {
	if f.isEnemy {
		return "enemy"
	}
	return "player"
}

// GetPosition returns the position of a specific entity
func (f *Formation) GetPosition(entity *ecs.Entity) (Position, bool) {
	pos, exists := f.positions[entity]
	return pos, exists
}

// GetAllPositions returns all entity positions
func (f *Formation) GetAllPositions() map[*ecs.Entity]Position {
	return f.positions
}

// UpdateFormation recalculates positions (useful when entities are removed/added)
func (f *Formation) UpdateFormation(entities []*ecs.Entity) {
	f.entities = entities
	f.positions = make(map[*ecs.Entity]Position)
	f.arrangeEntities()
}

// RemoveEntity removes an entity from the formation and reorganizes
func (f *Formation) RemoveEntity(entity *ecs.Entity) {
	// Remove from entities slice
	for i, e := range f.entities {
		if e == entity {
			f.entities = append(f.entities[:i], f.entities[i+1:]...)
			break
		}
	}

	// Remove from positions
	delete(f.positions, entity)

	// Rearrange remaining entities
	f.arrangeEntities()

	logger.Debug("📐 Removed %s from %s formation, reorganized remaining %d entities",
		entity.GetID(), f.getFormationType(), len(f.entities))
}

// GetFrontRowEntities returns entities in the front row (closest to enemies)
func (f *Formation) GetFrontRowEntities() []*ecs.Entity {
	frontRow := make([]*ecs.Entity, 0)

	for entity, pos := range f.positions {
		if pos.Row == 0 { // Front row is row 0
			frontRow = append(frontRow, entity)
		}
	}

	return frontRow
}

// GetBackRowEntities returns entities in the back row
func (f *Formation) GetBackRowEntities() []*ecs.Entity {
	backRow := make([]*ecs.Entity, 0)

	for entity, pos := range f.positions {
		if pos.Row == f.rows-1 { // Back row is the last row
			backRow = append(backRow, entity)
		}
	}

	return backRow
}

// IsEntityInFrontRow checks if an entity is in the front row
func (f *Formation) IsEntityInFrontRow(entity *ecs.Entity) bool {
	pos, exists := f.positions[entity]
	return exists && pos.Row == 0
}

// GetEntitiesInRow returns all entities in a specific row
func (f *Formation) GetEntitiesInRow(row int) []*ecs.Entity {
	entities := make([]*ecs.Entity, 0)

	for entity, pos := range f.positions {
		if pos.Row == row {
			entities = append(entities, entity)
		}
	}

	return entities
}

// GetFormationBounds returns the bounding box of the formation
func (f *Formation) GetFormationBounds() (minX, minY, maxX, maxY float64) {
	if len(f.positions) == 0 {
		return 0, 0, 0, 0
	}

	first := true
	for _, pos := range f.positions {
		if first {
			minX, minY = pos.X, pos.Y
			maxX, maxY = pos.X+f.entityWidth, pos.Y+f.entityHeight
			first = false
		} else {
			if pos.X < minX {
				minX = pos.X
			}
			if pos.Y < minY {
				minY = pos.Y
			}
			if pos.X+f.entityWidth > maxX {
				maxX = pos.X + f.entityWidth
			}
			if pos.Y+f.entityHeight > maxY {
				maxY = pos.Y + f.entityHeight
			}
		}
	}

	return minX, minY, maxX, maxY
}

// GetEntityCount returns the number of entities in the formation
func (f *Formation) GetEntityCount() int {
	return len(f.entities)
}

// GetAliveEntityCount returns the number of living entities in the formation
func (f *Formation) GetAliveEntityCount() int {
	alive := 0
	for _, entity := range f.entities {
		if stats := entity.RPGStats(); stats != nil && stats.CurrentHP > 0 {
			alive++
		}
	}
	return alive
}
