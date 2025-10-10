// Package engine implements the core game loop and state management.
// It uses an Entity-Component-System (ECS) architecture to manage game entities and their behaviors.
// The engine handles player input, updates game state, and renders graphics using the Ebiten library.
// It demonstrates basic player movement and rendering.
// To run this code, ensure you have the Ebiten library installed and
// an 'assets/sprites/player.png' image for the player sprite.
// If the asset is missing, a placeholder will be used.
package engine

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/gfx"
)

// Game represents the state of the game using an ECS architecture.
type Game struct {
	world        *ecs.World   // The game world containing all entities
	playerEntity *ecs.Entity  // Reference to the player entity for input handling
}

// NewGame creates a new game instance with an empty world
func NewGame() *Game {
	world := ecs.NewWorld()
	return &Game{
		world:        world,
		playerEntity: nil,
	}
}

// AddEntity adds an entity to the game world
func (g *Game) AddEntity(entity *ecs.Entity) {
	g.world.AddEntity(entity)
}

// SetPlayerEntity sets which entity should be controlled by player input
func (g *Game) SetPlayerEntity(entity *ecs.Entity) {
	g.playerEntity = entity
}

func (g *Game) Update() error {
	if g.playerEntity == nil {
		return nil // No player entity set
	}
	playerT := g.playerEntity.Transform()
	if playerT == nil {
		return nil // Player has no transform component
	}
	oldX, oldY := playerT.X, playerT.Y
	speed := 2.0

	// Handle player movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		playerT.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		playerT.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		playerT.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		playerT.X += speed
	}

	// Check for collisions with other entities
	for _, entity := range g.world.GetEntities() {
		// Skip player
		if entity == g.playerEntity {
			continue
		}
		// Skip entities without a collider
		if entity.Collider() == nil {
			continue
		}
		// Simple AABB collision detection
		if CheckCollision(playerT.Bounds(), entity.Transform().Bounds()) {
			playerT.X, playerT.Y = oldX, oldY // Revert to old position on collision
			log.Printf("Collision detected between player and entity at (%.2f, %.2f)", playerT.X, playerT.Y)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw all entities in the order they were added to the world
	// (assuming background entities are added first for proper layering)
	for _, entity := range g.world.GetEntities() {
		spriteC := entity.Sprite()
		if spriteC == nil {
			continue // Skip entities without a sprite
		}
		transform := entity.Transform()
		if transform == nil {
			continue // Skip entities without a transform
		}
		gfx.DrawSprite(screen, spriteC.Sprite, transform.X, transform.Y, spriteC.Scale)
	}

	// Display instructions
	ebitenutil.DebugPrint(screen, "Use arrow keys to move the player")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
