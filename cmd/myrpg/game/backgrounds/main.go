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
	// Load background sprite for game world area only (800x408)
	// Current layout: TopPanel=110px + Separator=2px + GameWorld=408px + BottomPanel=80px = 600px
	backgroundSprite, err := gfx.NewSpriteFromFile("assets/backgrounds/background.png", 800, 408)
	if err != nil {
		log.Fatal("failed to load background sprite:", err)
	}

	// Create background entity positioned in game world area
	// Y=112 to account for TopPanel(110px) + Separator(2px)
	background := ecs.NewEntity(constants.BackgroundEntityName)
	background.AddComponent(ecs.ComponentTransform, components.NewTransform(0, 112, 800, 408))
	background.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(backgroundSprite, 1.0, 0, 0))
	background.AddTag(ecs.TagBackground)

	return background
}
