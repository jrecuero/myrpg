// Package entities provides functions to create game entities like players and enemies
// It also includes utility functions for managing these entities
package entities

import (
	"fmt"
	"log"
	"time"

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

// CreateAnimatedPlayerWithJob creates a player entity with animated sprites
func CreateAnimatedPlayerWithJob(name string, x, y float64, job components.JobType, level int, animations CharacterAnimations) *ecs.Entity {
	// Create player entity with specified job
	player := ecs.NewEntity(name)
	player.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
	player.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 32, 32, 0, 0))
	player.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent(name, job, level))
	player.AddTag(ecs.TagPlayer)

	// Load animations
	err := AddAnimationsToEntity(player, animations)
	if err != nil {
		log.Printf("failed to load animated sprites for %s: %v", name, err)
		// Fallback to static sprite
		playerSprite, err := gfx.NewSpriteFromFile("assets/sprites/player.png", 32, 32)
		if err != nil {
			log.Fatal("failed to load fallback player sprite:", err)
		}
		player.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(playerSprite, 1.0, 0, 0))
	}

	return player
}

// CreateAnimatedPlayerWithSingleAnimation creates a player with a single animation (convenience function)
func CreateAnimatedPlayerWithSingleAnimation(name string, x, y float64, job components.JobType, level int, spriteSheetPath string, animState components.AnimationState) *ecs.Entity {
	animations := CharacterAnimations{
		Animations: []AnimationConfig{
			{
				State:         animState,
				SpriteSheet:   spriteSheetPath,
				StartFrame:    0,
				FrameCount:    0, // Use all frames
				FrameDuration: 200 * time.Millisecond,
				Loop:          true,
			},
		},
		Scale:   1.0,
		OffsetX: 0,
		OffsetY: 0,
	}

	return CreateAnimatedPlayerWithJob(name, x, y, job, level, animations)
}

// AnimationConfig defines configuration for a single animation
type AnimationConfig struct {
	State         components.AnimationState // Animation state (idle, walking, etc.)
	SpriteSheet   string                    // Path to sprite sheet file
	StartFrame    int                       // Starting frame index in sprite sheet
	FrameCount    int                       // Number of frames in animation
	FrameDuration time.Duration             // Duration per frame
	Loop          bool                      // Whether animation should loop
}

// CharacterAnimations defines all animations for a character
type CharacterAnimations struct {
	Animations []AnimationConfig // List of animation configurations
	Scale      float64           // Rendering scale
	OffsetX    float64           // X offset for rendering
	OffsetY    float64           // Y offset for rendering
}

// AddAnimationsToEntity adds multiple animations to an entity from configuration
func AddAnimationsToEntity(entity *ecs.Entity, animConfig CharacterAnimations) error {
	// Create animation component
	animationComponent := components.NewAnimationComponent(animConfig.Scale, animConfig.OffsetX, animConfig.OffsetY)

	// Load each animation
	for _, config := range animConfig.Animations {
		err := addSingleAnimation(animationComponent, config)
		if err != nil {
			log.Printf("Failed to load animation %s: %v", config.State.String(), err)
			continue // Skip failed animations but don't fail completely
		}
	}

	// Add animation component to entity (only if we have at least one animation)
	if len(animationComponent.Animations) > 0 {
		entity.AddComponent(ecs.ComponentAnimation, animationComponent)
		return nil
	}

	return fmt.Errorf("no animations could be loaded")
}

// addSingleAnimation loads a single animation from configuration
func addSingleAnimation(animComponent *components.AnimationComponent, config AnimationConfig) error {
	// Load sprite sheet
	spriteSheet, err := gfx.NewSpriteSheetFromFile(config.SpriteSheet, 32, 32)
	if err != nil {
		return err
	}

	// Get specific sprites for this animation
	var sprites []*gfx.Sprite
	if config.FrameCount == 0 {
		// Use all sprites from sheet
		sprites, err = spriteSheet.GetAllSprites()
	} else {
		// Use specific range of sprites
		sprites, err = spriteSheet.GetSprites(config.StartFrame, config.FrameCount)
	}

	if err != nil {
		return err
	}

	// Create animation
	animation := components.NewAnimation(sprites, config.FrameDuration, config.Loop)
	animComponent.AddAnimation(config.State, animation)

	return nil
}

// AddIdleAnimation adds an idle animation to an entity from a sprite sheet (legacy function for backwards compatibility)
func AddIdleAnimation(entity *ecs.Entity, spriteSheetPath string) error {
	config := CharacterAnimations{
		Animations: []AnimationConfig{
			{
				State:         components.AnimationIdle,
				SpriteSheet:   spriteSheetPath,
				StartFrame:    0,
				FrameCount:    0, // Use all frames
				FrameDuration: 200 * time.Millisecond,
				Loop:          true,
			},
		},
		Scale:   1.0,
		OffsetX: 0,
		OffsetY: 0,
	}

	return AddAnimationsToEntity(entity, config)
}

// CreateAnimatedEnemyWithJob creates an enemy entity with animated sprites
func CreateAnimatedEnemyWithJob(name string, x, y float64, job components.JobType, level int, spriteSheetPath string) *ecs.Entity {
	// Create enemy entity with specified job
	enemy := ecs.NewEntity(name)
	enemy.AddComponent(ecs.ComponentTransform, components.NewTransform(x, y, 32, 32))
	enemy.AddComponent(ecs.ComponentCollider, components.NewColliderComponent(true, 32, 32, 0, 0))
	enemy.AddComponent(ecs.ComponentRPGStats, components.NewRPGStatsComponent(name, job, level))
	enemy.AddTag(ecs.TagEnemy)

	// Load sprite sheet and create animations
	err := AddIdleAnimation(enemy, spriteSheetPath)
	if err != nil {
		log.Printf("failed to load animated sprites for %s: %v", name, err)
		// Fallback to static sprite
		enemySprite, err := gfx.NewSpriteFromFile("assets/sprites/enemy.png", 32, 32)
		if err != nil {
			log.Fatal("failed to load fallback enemy sprite:", err)
		}
		enemy.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(enemySprite, 1.0, 0, 0))
	}

	return enemy
}
