// Command myrpg is a simple 2D RPG game using the Ebiten library.
// It demonstrates basic player movement and rendering.
// To run this code, ensure you have the Ebiten library installed and
// an 'assets/sprites/player.png' image for the player sprite.
// If the asset is missing, a placeholder will be used.

package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jrecuero/myrpg/cmd/myrpg/game/backgrounds"
	"github.com/jrecuero/myrpg/cmd/myrpg/game/entities"
	"github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/engine"
)

func main() {
	// Create a new game instance
	// Note: As of Go 1.20+, the global random generator is automatically seeded
	game := engine.NewGame()

	// Create and add background (render first for proper layering)
	background := backgrounds.CreateMainBackground()
	game.AddEntity(background)

	// Create and add player entities with different jobs and levels
	// Create animated hero with multiple animations (if available)
	heroAnimations := entities.CharacterAnimations{
		Animations: []entities.AnimationConfig{
			{
				State:         components.AnimationIdle,
				SpriteSheet:   "assets/sprites/hero/hero-idle.png",
				StartFrame:    0,
				FrameCount:    0, // Use all 6 frames from sprite sheet
				FrameDuration: 200 * time.Millisecond,
				Loop:          true,
			},
			{
				State:         components.AnimationWalking,
				SpriteSheet:   "assets/sprites/hero/hero-walk.png",
				StartFrame:    0,
				FrameCount:    0,                      // Use all 6 frames from sprite sheet
				FrameDuration: 150 * time.Millisecond, // Slightly faster for walking
				Loop:          true,
			},
			{
				State:         components.AnimationAttacking,
				SpriteSheet:   "assets/sprites/hero/hero-sword.png",
				StartFrame:    0,
				FrameCount:    0,                      // Use all 6 frames from sprite sheet
				FrameDuration: 100 * time.Millisecond, // Fast attack animation
				Loop:          true,                   // Loop attack animation during attack period
			},
		},
		Scale:   1.0,
		OffsetX: 0,
		OffsetY: 0,
	}
	warrior := entities.CreateAnimatedPlayerWithJob("Conan", constants.Player1StartX, constants.Player1StartY, components.JobWarrior, 3, heroAnimations)
	game.AddEntity(warrior)

	mage := entities.CreatePlayerWithJob("Gandalf", constants.Player2StartX, constants.Player2StartY, components.JobMage, 2)
	game.AddEntity(mage)

	rogue := entities.CreatePlayerWithJob("Robin", constants.Player3StartX, constants.Player3StartY, components.JobRogue, 4)
	game.AddEntity(rogue)

	// Create and add enemy entities with different jobs and levels
	goblin := entities.CreateEnemyWithJob("Goblin Scout", constants.Enemy1StartX, constants.Enemy1StartY, components.JobRogue, 2)
	game.AddEntity(goblin)

	orcWarrior := entities.CreateEnemyWithJob("Orc Warrior", constants.Enemy2StartX, constants.Enemy2StartY, components.JobWarrior, 5)
	game.AddEntity(orcWarrior)

	// Add more enemies with varied jobs and levels
	koboldMage := entities.CreateEnemyWithJob("Kobold Mage", constants.Enemy3StartX, constants.Enemy3StartY, components.JobMage, 3)
	game.AddEntity(koboldMage)

	banditArcher := entities.CreateEnemyWithJob("Bandit Archer", constants.Enemy4StartX, constants.Enemy4StartY, components.JobArcher, 4)
	game.AddEntity(banditArcher)

	orcShaman := entities.CreateEnemyWithJob("Orc Shaman", constants.Enemy5StartX, constants.Enemy5StartY, components.JobCleric, 6)
	game.AddEntity(orcShaman)

	goblinThief := entities.CreateEnemyWithJob("Goblin Thief", constants.Enemy6StartX, constants.Enemy6StartY, components.JobRogue, 1)
	game.AddEntity(goblinThief)

	// Configure attack animation duration (customizable)
	game.SetAttackAnimationDuration(1500 * time.Millisecond) // 1.5 seconds for attack animation

	// Initialize the game with welcome messages
	game.InitializeGame()

	// Set window properties and run the game
	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("My RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
