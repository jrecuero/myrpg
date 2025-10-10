// Package engine implements the core game loop and state management.
// It uses an Entity-Component-System (ECS) architecture to manage game entities and their behaviors.
// The engine handles player input, updates game state, and renders graphics using the Ebiten library.
// It demonstrates basic player movement and rendering.
// To run this code, ensure you have the Ebiten library installed and
// an 'assets/sprites/player.png' image for the player sprite.
// If the asset is missing, a placeholder will be used.
package engine

import (
	"image"
)

// CheckCollision checks if two rectangles overlap (collide).
// r1 and r2 are the rectangles to check for collision.
// returns true if the rectangles overlap, false otherwise.
func CheckCollision(r1, r2 image.Rectangle) bool {
	return r1.Overlaps(r2)
}
