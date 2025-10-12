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
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/gfx"
	"github.com/jrecuero/myrpg/internal/ui"
)

// Game represents the state of the game using an ECS architecture.
type Game struct {
	world             *ecs.World    // The game world containing all entities
	activePlayerIndex int           // Index of the currently active player
	tabKeyPressed     bool          // Track TAB key state to prevent multiple switches
	uiManager         *ui.UIManager // UI system for panels and messages
	battleSystem      *BattleSystem // Battle system for combat
}

// NewGame creates a new game instance with an empty world
func NewGame() *Game {
	world := ecs.NewWorld()
	uiManager := ui.NewUIManager()
	battleSystem := NewBattleSystem()

	game := &Game{
		world:             world,
		activePlayerIndex: 0,
		tabKeyPressed:     false,
		uiManager:         uiManager,
		battleSystem:      battleSystem,
	}

	// Set up battle system callbacks
	battleSystem.SetMessageCallback(uiManager.AddMessage)
	battleSystem.SetSwitchPlayerCallback(game.SwitchToNextPlayer)

	return game
}

// AddEntity adds an entity to the game world
func (g *Game) AddEntity(entity *ecs.Entity) {
	g.world.AddEntity(entity)
}

// RemoveEntity removes an entity from the game world
func (g *Game) RemoveEntity(entity *ecs.Entity) {
	g.world.RemoveEntity(entity)
}

// SetAttackAnimationDuration configures how long attack animations should last
func (g *Game) SetAttackAnimationDuration(duration time.Duration) {
	g.battleSystem.SetAttackAnimationDuration(duration)
}

// CheckAndRemoveDeadEntities removes entities with HP <= 0
func (g *Game) CheckAndRemoveDeadEntities() {
	entitiesToRemove := []*ecs.Entity{}

	for _, entity := range g.world.GetEntities() {
		stats := entity.RPGStats()
		if stats != nil && stats.CurrentHP <= 0 {
			g.uiManager.AddMessage(fmt.Sprintf("%s has been defeated and removed from battle!", stats.Name))
			entitiesToRemove = append(entitiesToRemove, entity)
		}
	}

	// Remove dead entities
	for _, entity := range entitiesToRemove {
		g.RemoveEntity(entity)

		// If it was the active player, switch to next player
		if entity.HasComponent("Player") {
			g.SwitchToNextPlayer()
		}
	}
}

// InitializeGame sets up the initial game state and messages
func (g *Game) InitializeGame() {
	g.uiManager.AddMessage("Welcome to MyRPG!")
	g.uiManager.AddMessage("Use arrow keys to move, TAB to switch between players")

	// Add message about the current active player
	activePlayer := g.GetActivePlayer()
	if activePlayer != nil {
		stats := activePlayer.RPGStats()
		if stats != nil {
			initMsg := fmt.Sprintf("Starting as %s (%s Level %d)",
				stats.Name, stats.Job.String(), stats.Level)
			g.uiManager.AddMessage(initMsg)
		}
	}
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

	// Add message about player switch
	activePlayer := g.GetActivePlayer()
	if activePlayer != nil {
		stats := activePlayer.RPGStats()
		if stats != nil {
			switchMsg := fmt.Sprintf("Switched to %s (%s Level %d)",
				stats.Name, stats.Job.String(), stats.Level)
			g.uiManager.AddMessage(switchMsg)
		}
	}
}

func (g *Game) Update() error {
	// Update battle system first
	g.battleSystem.Update()

	// Only handle movement and player switching if not in battle
	if !g.battleSystem.IsInBattle() {
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
		isMoving := false

		// Handle movement for ONLY the active player
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			playerT.Y -= speed
			isMoving = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			playerT.Y += speed
			isMoving = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			playerT.X -= speed
			isMoving = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			playerT.X += speed
			isMoving = true
		}

		// Update animation state based on movement
		if animationC := activePlayer.Animation(); animationC != nil {
			if isMoving {
				// Try to set walking animation, fallback to idle if not available
				if !animationC.SetStateIfAvailable(components.AnimationWalking) {
					animationC.SetStateIfAvailable(components.AnimationIdle)
				}
			} else {
				// Set idle animation when not moving
				animationC.SetStateIfAvailable(components.AnimationIdle)
			}
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

			// Check collision type
			if CheckCollision(playerT.Bounds(), entity.Transform().Bounds()) {
				playerT.X, playerT.Y = oldX, oldY // Revert to old position on collision

				// Determine if this is an enemy or regular collision
				if entity.HasTag("player") {
					// Player-to-player collision
					collisionMsg := fmt.Sprintf("Collision: %s bumped into %s", activePlayer.Name, entity.Name)
					g.uiManager.AddMessage(collisionMsg)
				} else if entity.RPGStats() != nil {
					// Enemy collision - has RPG stats but no Player tag
					g.battleSystem.StartBattle(activePlayer, entity)
					break // Only start one battle at a time
				} else {
					// Regular collision (wall, object, etc.)
					collisionMsg := fmt.Sprintf("Collision: %s hit an obstacle", activePlayer.Name)
					g.uiManager.AddMessage(collisionMsg)
				}
			}
		}

		// Check and remove dead entities after movement
		g.CheckAndRemoveDeadEntities()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with black background
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Draw UI panels first (top and bottom)
	activePlayer := g.GetActivePlayer()
	var activePlayerStats *components.RPGStatsComponent
	if activePlayer != nil {
		activePlayerStats = activePlayer.RPGStats()
	}

	// Draw top panel with player info
	g.uiManager.DrawTopPanel(screen, activePlayerStats)

	// Draw game world background
	g.uiManager.DrawGameWorldBackground(screen)

	// Draw all game entities in the game world area
	for _, entity := range g.world.GetEntities() {
		transform := entity.Transform()
		if transform == nil {
			continue // Skip entities without a transform
		}

		// Check for animated sprite first, fallback to static sprite
		animationC := entity.Animation()
		if animationC != nil {
			// Update animation
			animationC.Update()

			// Draw current animation frame
			currentSprite := animationC.GetCurrentSprite()
			if currentSprite != nil {
				gfx.DrawSprite(screen, currentSprite,
					transform.X+animationC.OffsetX,
					transform.Y+animationC.OffsetY,
					animationC.Scale)
			}
		} else {
			// Fallback to static sprite
			spriteC := entity.Sprite()
			if spriteC != nil {
				gfx.DrawSprite(screen, spriteC.Sprite, transform.X, transform.Y, spriteC.Scale)
			}
		}
	}

	// Highlight the active player with a green rectangle outline
	if activePlayer != nil {
		playerT := activePlayer.Transform()
		if playerT != nil {
			// Draw a green outline around the active player
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

	// Draw bottom panel with messages and commands
	g.uiManager.DrawBottomPanel(screen)

	// Draw battle menu if in battle
	if g.battleSystem.IsInBattle() {
		battleText := g.battleSystem.GetBattleMenuText()
		g.uiManager.DrawBattleMenu(screen, battleText)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
