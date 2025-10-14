package backgrounds

import (
	"log"

	"github.com/jrecuero/myrpg/cmd/myrpg/game/constants"
	displayConstants "github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/gfx"
)

// CreateMainBackground creates the main background entity for the game
func CreateMainBackground() *ecs.Entity {
	// Load background sprite for game world area using display constants
	// Current layout: TopPanel=110px + Separator=2px + GameWorld=408px + BottomPanel=80px = 600px
	backgroundSprite, err := gfx.NewSpriteFromFile("assets/backgrounds/background.png", displayConstants.BackgroundWidth, displayConstants.BackgroundHeight)
	if err != nil {
		log.Fatal("failed to load background sprite:", err)
	}

	// Create background entity positioned in game world area
	background := ecs.NewEntity(constants.BackgroundEntityName)
	background.AddComponent(ecs.ComponentTransform, components.NewTransform(
		displayConstants.BackgroundX,
		displayConstants.BackgroundY,
		displayConstants.BackgroundWidth,
		displayConstants.BackgroundHeight))
	background.AddComponent(ecs.ComponentSprite, components.NewSpriteComponent(backgroundSprite, 1.0, 0, 0))
	background.AddTag(ecs.TagBackground)

	return background
}
