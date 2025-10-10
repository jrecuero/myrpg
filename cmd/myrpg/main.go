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
	"github.com/jrecuero/myrpg/internal/engine"
)

func main() {
	// Create a new game instance
	game := engine.NewGame()

	// Create and add background (render first for proper layering)
	background := backgrounds.CreateMainBackground()
	game.AddEntity(background)

	// Create and add player entities (multiple players!)
	player1 := entities.CreatePlayer() // Default player at (100, 100)
	game.AddEntity(player1)
	
	player2 := entities.CreatePlayerAtPosition("Player2", 150, 100) // Second player
	game.AddEntity(player2)

	// Create and add enemy entity
	enemy := entities.CreateEnemy(200, 200)
	game.AddEntity(enemy)

	// Set window properties and run the game
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("My RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
