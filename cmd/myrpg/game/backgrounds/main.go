package backgrounds

import (
	"log"

	"github.com/jrecuero/myrpg/cmd/myrpg/game/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/gfx"
)

// CreateMainBackground creates the main background entity for the game
func CreateMainBackground() *ecs.Entity {
	// Load background sprite for game world area only (800x440)
	backgroundSprite, err := gfx.NewSpriteFromFile("assets/backgrounds/background.png", 800, 440)
	if err != nil {
		log.Fatal("failed to load background sprite:", err)
	}

	// Create background entity positioned in game world area
	background := ecs.NewEntity(constants.BackgroundEntityName)
	background.AddComponent(ecs.ComponentTransform, components.NewTransform(0, 80, 800, 440)) // Y=80 for top panel offset
	background.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(backgroundSprite, 1.0, 0, 0))
	background.AddTag(ecs.TagBackground)
	
	return background
}