// Command myrpg is a simple 2D RPG game using the Ebiten library.
// It demonstrates basic player movement and rendering.
// To run this code, ensure you have the Ebiten library installed and
// an 'assets/sprites/player.png' image for the player sprite.
// If the asset is missing, a placeholder will be used.

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jrecuero/myrpg/cmd/myrpg/game/backgrounds"
	"github.com/jrecuero/myrpg/cmd/myrpg/game/entities"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/engine"
)

func main() {
	// Create a new game instance
	game := engine.NewGame()

	// Create and add background (render first for proper layering)
	background := backgrounds.CreateMainBackground()
	game.AddEntity(background)

	// Create and add player entities with different jobs and levels
	warrior := entities.CreatePlayerWithJob("Conan", 100, 100, components.JobWarrior, 3)
	game.AddEntity(warrior)
	
	mage := entities.CreatePlayerWithJob("Gandalf", 150, 100, components.JobMage, 2)
	game.AddEntity(mage)
	
	rogue := entities.CreatePlayerWithJob("Robin", 200, 100, components.JobRogue, 4)
	game.AddEntity(rogue)

	// Create and add enemy entities with different jobs and levels
	goblin := entities.CreateEnemyWithJob("Goblin Scout", 300, 200, components.JobRogue, 2)
	game.AddEntity(goblin)
	
	orcWarrior := entities.CreateEnemyWithJob("Orc Warrior", 350, 250, components.JobWarrior, 5)
	game.AddEntity(orcWarrior)

	// Initialize the game with welcome messages
	game.InitializeGame()

	// Set window properties and run the game
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("My RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
