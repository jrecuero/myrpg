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

const (
	PlayerEntityName     = "Player"     // Name constant for the player entity
	BackgroundEntityName = "Background" // Name constant for the background entity
)

// Game represents the state of the game using an ECS architecture.
type Game struct {
	world *ecs.World // The game world containing all entities
}

func NewGame() *Game {
	world := ecs.NewWorld()

	// Load background sprite
	backgroundSprite, err := gfx.NewSpriteFromFile("assets/backgrounds/background.png", 800, 600)
	if err != nil {
		log.Fatal("failed to load background sprite:", err)
	}

	// Create background entity
	background := ecs.NewEntity(BackgroundEntityName)
	background.AddComponent(ecs.ComponentTransform, ecs.NewTransform(0, 0, 800, 600))
	background.AddComponent(ecs.ComponentSprite, ecs.NewSpriteComponent(backgroundSprite, 1.0, 0, 0))
	world.AddEntity(background)

	// Load player sprite
	playerSprite, err := gfx.NewSpriteFromFile("assets/sprites/player.png", 32, 32)
	if err != nil {
		log.Fatal("failed to load player sprite:", err)
	}

	// Load enemy sprite
	enemySprite, err := gfx.NewSpriteFromFile("assets/sprites/enemy.png", 32, 32)
	if err != nil {
		log.Fatal("failed to load enemy sprite:", err)
	}

	// Create player entity
	player := ecs.NewEntity(PlayerEntityName)
	player.AddComponent(ecs.ComponentTransform, ecs.NewTransform(100, 100, 32, 32))
	player.AddComponent(ecs.ComponentSprite, ecs.NewSpriteComponent(playerSprite, 1.0, 0, 0))
	player.AddComponent(ecs.ComponentCollider, ecs.NewColliderComponent(true, 32, 32, 0, 0))
	world.AddEntity(player)

	// Create enemy entity
	enemy := ecs.NewEntity("Enemy")
	enemy.AddComponent(ecs.ComponentTransform, ecs.NewTransform(200, 200, 32, 32))
	enemy.AddComponent(ecs.ComponentSprite, ecs.NewSpriteComponent(enemySprite, 1.0, 0, 0))
	enemy.AddComponent(ecs.ComponentCollider, ecs.NewColliderComponent(true, 32, 32, 0, 0))
	world.AddEntity(enemy)

	return &Game{
		world: world,
	}
}

func (g *Game) Update() error {
	player := g.world.FindByName(PlayerEntityName)
	if player == nil {
		return nil // No player found, skip update
	}
	playerT := player.Transform()
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
		if entity.Name == PlayerEntityName {
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
	// Draw background first
	background := g.world.FindByName(BackgroundEntityName)
	if background != nil {
		spriteC := background.Sprite()
		if spriteC != nil {
			transform := background.Transform()
			if transform != nil {
				gfx.DrawSprite(screen, spriteC.Sprite, transform.X, transform.Y, spriteC.Scale)
			}
		}
	}

	// Draw all other entities
	for _, entity := range g.world.GetEntities() {
		// Skip background as it's already drawn
		if entity.Name == BackgroundEntityName {
			continue
		}
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
