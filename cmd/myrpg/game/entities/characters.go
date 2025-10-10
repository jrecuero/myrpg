package entities

import (
	"log"

	"github.com/jrecuero/myrpg/cmd/myrpg/game/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
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
	player.AddComponent(ecs.ComponentTransform, ecs.NewTransform(100, 100, 32, 32))
	player.AddComponent(ecs.ComponentSprite, ecs.NewSpriteComponent(playerSprite, 1.0, 0, 0))
	player.AddComponent(ecs.ComponentCollider, ecs.NewColliderComponent(true, 32, 32, 0, 0))
	
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
	enemy.AddComponent(ecs.ComponentTransform, ecs.NewTransform(x, y, 32, 32))
	enemy.AddComponent(ecs.ComponentSprite, ecs.NewSpriteComponent(enemySprite, 1.0, 0, 0))
	enemy.AddComponent(ecs.ComponentCollider, ecs.NewColliderComponent(true, 32, 32, 0, 0))
	
	return enemy
}