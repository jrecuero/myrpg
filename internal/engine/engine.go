// Package engine implements the core game loop and state management.
// It uses an Entity-Component-System (ECS) architecture to manage game entities and their behaviors.
// The engine handles player input, updates game state, and renders graphics using the Ebiten library.
// It demonstrates basic player movement and rendering.
// To run this code, ensure you have the Ebiten library installed and
// an 'assets/sprites/player.png' image for the player sprite.
// If the asset is missing, a placeholder will be used.
package engine

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/gfx"
)

// Game represents the state of the game using an ECS architecture.
type Game struct {
	world             *ecs.World // The game world containing all entities
	activePlayerIndex int        // Index of the currently active player
	tabKeyPressed     bool       // Track TAB key state to prevent multiple switches
}

// NewGame creates a new game instance with an empty world
func NewGame() *Game {
	world := ecs.NewWorld()
	return &Game{
		world:             world,
		activePlayerIndex: 0,
		tabKeyPressed:     false,
	}
}

// AddEntity adds an entity to the game world
func (g *Game) AddEntity(entity *ecs.Entity) {
	g.world.AddEntity(entity)
}

// GetPlayerEntities returns all player entities
func (g *Game) GetPlayerEntities() []*ecs.Entity {
	return g.world.FindWithTag(ecs.TagPlayer)
}

// GetActivePlayer returns the currently active player entity
func (g *Game) GetActivePlayer() *ecs.Entity {
	players := g.GetPlayerEntities()
	if len(players) == 0 {
		return nil
	}
	if g.activePlayerIndex >= len(players) {
		g.activePlayerIndex = 0 // Reset if index is out of bounds
	}
	return players[g.activePlayerIndex]
}

// SwitchToNextPlayer cycles to the next player
func (g *Game) SwitchToNextPlayer() {
	players := g.GetPlayerEntities()
	if len(players) <= 1 {
		return // No switching needed with 0 or 1 player
	}
	g.activePlayerIndex = (g.activePlayerIndex + 1) % len(players)
}

func (g *Game) Update() error {
	// Handle TAB key for player switching
	if ebiten.IsKeyPressed(ebiten.KeyTab) {
		if !g.tabKeyPressed {
			g.SwitchToNextPlayer()
			g.tabKeyPressed = true
		}
	} else {
		g.tabKeyPressed = false
	}

	// Get the currently active player
	activePlayer := g.GetActivePlayer()
	if activePlayer == nil {
		return nil // No active player
	}

	playerT := activePlayer.Transform()
	if playerT == nil {
		return nil // Active player has no transform component
	}

	oldX, oldY := playerT.X, playerT.Y
	speed := 2.0

	// Handle movement for ONLY the active player
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
		// Skip the active player itself
		if entity == activePlayer {
			continue
		}
		// Skip entities without a collider
		if entity.Collider() == nil {
			continue
		}
		// Simple AABB collision detection
		if CheckCollision(playerT.Bounds(), entity.Transform().Bounds()) {
			playerT.X, playerT.Y = oldX, oldY // Revert to old position on collision
			log.Printf("Collision detected between active player '%s' and entity '%s' at (%.2f, %.2f)",
				activePlayer.Name, entity.Name, playerT.X, playerT.Y)
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

	// Highlight the active player with a green rectangle outline
	activePlayer := g.GetActivePlayer()
	if activePlayer != nil {
		playerT := activePlayer.Transform()
		if playerT != nil {
			// Draw a green outline around the active player
			// Use direct coordinates instead of bounds for simpler conversion
			x := playerT.X
			y := playerT.Y
			width := 32.0 // Assuming standard player sprite size
			height := 32.0

			// Draw rectangle outline using vector.FillRect
			vector.FillRect(screen,
				float32(x-2), float32(y-2),
				float32(width+4), float32(2), // Top border
				color.RGBA{0, 255, 0, 255}, false)
			vector.FillRect(screen,
				float32(x-2), float32(y+height),
				float32(width+4), float32(2), // Bottom border
				color.RGBA{0, 255, 0, 255}, false)
			vector.FillRect(screen,
				float32(x-2), float32(y-2),
				float32(2), float32(height+4), // Left border
				color.RGBA{0, 255, 0, 255}, false)
			vector.FillRect(screen,
				float32(x+width), float32(y-2),
				float32(2), float32(height+4), // Right border
				color.RGBA{0, 255, 0, 255}, false)
		}
	}

	// Display instructions and active player stats
	instructions := "Use arrow keys to move active player, TAB to switch"
	
	// Add active player stats to display
	if activePlayer := g.GetActivePlayer(); activePlayer != nil {
		stats := activePlayer.RPGStats()
		if stats != nil {
			instructions += fmt.Sprintf("\n\nActive Player: %s (%s Level %d)", 
				stats.Name, stats.Job.String(), stats.Level)
			instructions += fmt.Sprintf("\nHP: %d/%d  MP: %d/%d", 
				stats.CurrentHP, stats.MaxHP, stats.CurrentMP, stats.MaxMP)
			instructions += fmt.Sprintf("\nAttack: %d  Defense: %d  Speed: %d", 
				stats.Attack, stats.Defense, stats.Speed)
			instructions += fmt.Sprintf("\nEXP: %d/%d", 
				stats.Experience, stats.ExpToNext)
		}
	}
	
	ebitenutil.DebugPrint(screen, instructions)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
