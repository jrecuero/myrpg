// Utility functions for entity management
// These functions help in retrieving and counting entities by their tags
package entities

import (
	"github.com/jrecuero/myrpg/internal/ecs"
)

// EntityUtils provides utility functions for working with tagged entities
type EntityUtils struct{}

// GetAllPlayers returns all player entities from a world
func (EntityUtils) GetAllPlayers(world *ecs.World) []*ecs.Entity {
	return world.FindWithTag(ecs.TagPlayer)
}

// GetAllEnemies returns all enemy entities from a world
func (EntityUtils) GetAllEnemies(world *ecs.World) []*ecs.Entity {
	return world.FindWithTag(ecs.TagEnemy)
}

// GetAllBackgrounds returns all background entities from a world
func (EntityUtils) GetAllBackgrounds(world *ecs.World) []*ecs.Entity {
	return world.FindWithTag(ecs.TagBackground)
}

// CountEntitiesWithTag counts how many entities have a specific tag
func (EntityUtils) CountEntitiesWithTag(world *ecs.World, tag string) int {
	return len(world.FindWithTag(tag))
}

// Global instance for convenience
var Utils EntityUtils
