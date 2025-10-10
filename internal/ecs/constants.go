// Package ecs provides an Entity-Component-System (ECS) framework for game development.
// It defines core components such as Transform, SpriteComponent, and ColliderComponent
// that can be attached to entities to define their properties and behaviors.
// The ECS architecture allows for flexible and modular game object management.
// Each component struct includes fields relevant to its purpose, along with constructor functions
// for easy instantiation. The Transform component handles position and size, the SpriteComponent
// manages visual representation, and the ColliderComponent defines collision properties.
package ecs

const (
	ComponentTransform = "transform"
	ComponentSprite    = "sprite"
	ComponentCollider  = "collider"
)

// Common entity tags
const (
	TagPlayer     = "player"
	TagEnemy      = "enemy"
	TagBackground = "background"
	TagNPC        = "npc"
	TagProjectile = "projectile"
)
