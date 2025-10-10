package backgrounds

import (
	"log"

	"github.com/jrecuero/myrpg/cmd/myrpg/game/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/gfx"
)

// CreateMainBackground creates the main background entity for the game
func CreateMainBackground() *ecs.Entity {
	// Load background sprite
	backgroundSprite, err := gfx.NewSpriteFromFile("assets/backgrounds/background.png", 800, 600)
	if err != nil {
		log.Fatal("failed to load background sprite:", err)
	}

	// Create background entity
	background := ecs.NewEntity(constants.BackgroundEntityName)
	background.AddComponent(ecs.ComponentTransform, ecs.NewTransform(0, 0, 800, 600))
	background.AddComponent(ecs.ComponentSprite, ecs.NewSpriteComponent(backgroundSprite, 1.0, 0, 0))
	
	return background
}