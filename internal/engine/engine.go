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
	Player *ecs.Entity   // The player entity
	Enemy  *ecs.Entity   // The enemy entity
	BG     *ebiten.Image // The background image
}

func NewGame() *Game {
	// Load background image
	bg, _, err := ebitenutil.NewImageFromFile("assets/backgrounds/background.png")
	if err != nil {
		log.Fatal(err)
	}

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
	player := ecs.NewEntity()
	player.AddComponent(ecs.ComponentTransform, ecs.NewTransform(100, 100, 32, 32))
	player.AddComponent(ecs.ComponentSprite, ecs.NewSpriteComponent(playerSprite, 1.0, 0, 0))
	player.AddComponent(ecs.ComponentCollider, ecs.NewColliderComponent(true, 32, 32, 0, 0))

	// Create enemy entity
	enemy := ecs.NewEntity()
	enemy.AddComponent(ecs.ComponentTransform, ecs.NewTransform(200, 200, 32, 32))
	enemy.AddComponent(ecs.ComponentSprite, ecs.NewSpriteComponent(enemySprite, 1.0, 0, 0))
	enemy.AddComponent(ecs.ComponentCollider, ecs.NewColliderComponent(true, 32, 32, 0, 0))

	return &Game{
		Player: player,
		Enemy:  enemy,
		BG:     bg,
	}
}

func (g *Game) Update() error {
	playerT := g.Player.Transform()
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

	// Simple collision detection with enemy
	if CheckCollision(g.Player.Transform().Bounds(), g.Enemy.Transform().Bounds()) {
		playerT.X, playerT.Y = oldX, oldY // Revert to old position on collision
		log.Printf("Collision detected between player and enemy at (%.2f, %.2f)", playerT.X, playerT.Y)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.BG, op)

	// Draw player
	playerT := g.Player.Transform()
	playerS, ok := g.Player.GetComponent("sprite")
	if !ok {
		log.Fatal("player entity missing sprite component")
	}
	gfx.DrawSprite(screen, playerS.(*ecs.SpriteComponent).Sprite, playerT.X, playerT.Y, playerS.(*ecs.SpriteComponent).Scale)

	// Draw enemy
	enemyT := g.Enemy.Transform()
	enemyS, ok := g.Enemy.GetComponent("sprite")
	if !ok {
		log.Fatal("enemy entity missing sprite component")
	}
	gfx.DrawSprite(screen, enemyS.(*ecs.SpriteComponent).Sprite, enemyT.X, enemyT.Y, enemyS.(*ecs.SpriteComponent).Scale)

	// Display instructions
	ebitenutil.DebugPrint(screen, "Use arrow keys to move the player")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
