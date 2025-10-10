// Package entities provides functions to create game entities like players and enemies
// It also includes utility functions for managing these entities
package entities

import (
	"log"

	"github.com/jrecuero/myrpg/cmd/myrpg/game/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/gfx"
)

// CreatePlayer creates the player entity for the game
func CreatePlayer() *ecs.Entity {
	// Load player sprite
	playerSprite, err := gfx.NewSpriteFromFile("assets/sprites/player.png", 32, 32)
	if err != nil {
		log.Fatal("failed to load player sprite:", err)
	}

	// Create player entity
	player := ecs.NewEntity(constants.PlayerEntityName)
	player.AddComponent(ecs.ComponentTransform, components.NewTransform(100, 100, 32, 32))
	player.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(playerSprite, 1.0, 0, 0))
	player.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 32, 32, 0, 0))
	player.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent("Hero", components.JobWarrior, 1))
	player.AddTag(ecs.TagPlayer)
	
	return player
}

// CreateEnemy creates an enemy entity at the specified position
func CreateEnemy(x, y float64) *ecs.Entity {
	// Load enemy sprite
	enemySprite, err := gfx.NewSpriteFromFile("assets/sprites/enemy.png", 32, 32)
	if err != nil {
		log.Fatal("failed to load enemy sprite:", err)
	}

	// Create enemy entity
	enemy := ecs.NewEntity(constants.EnemyEntityName)
	enemy.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
	enemy.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(enemySprite, 1.0, 0, 0))
	enemy.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 32, 32, 0, 0))
	enemy.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent("Goblin", components.JobRogue, 1))
	enemy.AddTag(ecs.TagEnemy)
	
	return enemy
}

// CreatePlayerAtPosition creates a player entity at the specified position with a custom name
func CreatePlayerAtPosition(name string, x, y float64) *ecs.Entity {
	// Load player sprite
	playerSprite, err := gfx.NewSpriteFromFile("assets/sprites/player.png", 32, 32)
	if err != nil {
		log.Fatal("failed to load player sprite:", err)
	}

	// Create player entity with custom name
	player := ecs.NewEntity(name)
	player.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
	player.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(playerSprite, 1.0, 0, 0))
	player.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 32, 32, 0, 0))
	player.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent(name, components.JobWarrior, 1))
	player.AddTag(ecs.TagPlayer)
	
	return player
}

// CreatePlayerWithJob creates a player entity with a specific job/class
func CreatePlayerWithJob(name string, x, y float64, job components.JobType, level int) *ecs.Entity {
	// Load player sprite
	playerSprite, err := gfx.NewSpriteFromFile("assets/sprites/player.png", 32, 32)
	if err != nil {
		log.Fatal("failed to load player sprite:", err)
	}

	// Create player entity with specified job
	player := ecs.NewEntity(name)
	player.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
	player.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(playerSprite, 1.0, 0, 0))
	player.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 32, 32, 0, 0))
	player.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent(name, job, level))
	player.AddTag(ecs.TagPlayer)
	
	return player
}

// CreateEnemyWithJob creates an enemy entity with a specific job/class and level
func CreateEnemyWithJob(name string, x, y float64, job components.JobType, level int) *ecs.Entity {
	// Load enemy sprite
	enemySprite, err := gfx.NewSpriteFromFile("assets/sprites/enemy.png", 32, 32)
	if err != nil {
		log.Fatal("failed to load enemy sprite:", err)
	}

	// Create enemy entity with specified job
	enemy := ecs.NewEntity(name)
	enemy.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
	enemy.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(enemySprite, 1.0, 0, 0))
	enemy.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 32, 32, 0, 0))
	enemy.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent(name, job, level))
	enemy.AddTag(ecs.TagEnemy)
	
	return enemy
}