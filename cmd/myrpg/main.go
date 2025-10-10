// Command myrpg is a simple 2D RPG game using the Ebiten library.
// It demonstrates basic player movement and rendering.
// To run this code, ensure you have the Ebiten library installed and
// an 'assets/sprites/player.png' image for the player sprite.
// If the asset is missing, a placeholder will be used.

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jrecuero/myrpg/internal/engine"
)

func main() {
	game := engine.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
