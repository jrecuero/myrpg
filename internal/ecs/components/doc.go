// Package components provides all the component types for the ECS (Entity-Component-System) framework.
// Components define the data and properties that can be attached to entities.
//
// This package is organized into logical groups:
// - transform.go: Position, size, and spatial components
// - rendering.go: Visual representation and sprite components
// - physics.go: Collision detection and physics components
// - rpg_stats.go: RPG-specific stats like HP, level, job classes, etc.
//
// Each component should be a simple data structure with minimal logic.
// Business logic should be handled by systems that operate on these components.
package components
