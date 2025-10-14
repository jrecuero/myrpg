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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/gfx"
	"github.com/jrecuero/myrpg/internal/logger"
	"github.com/jrecuero/myrpg/internal/tactical"
	"github.com/jrecuero/myrpg/internal/ui"
)

// Game represents the state of the game using an ECS architecture.
type Game struct {
	world              *ecs.World          // The game world containing all entities
	activePlayerIndex  int                 // Index of the currently active player
	tabKeyPressed      bool                // Track TAB key state to prevent multiple switches
	uiManager          *ui.UIManager       // UI system for panels and messages
	battleSystem       *BattleSystem       // Battle system for combat
	tacticalManager    *TacticalManager    // Tactical combat system
	partyManager       *PartyManager       // Party and team management
	enemyGroupManager  *EnemyGroupManager  // Enemy group formation
	tacticalDeployment *TacticalDeployment // Unit deployment for tactical combat
	currentMode        GameMode            // Current game mode (exploration/tactical)
}

// NewGame creates a new game instance with an empty world
func NewGame() *Game {
	world := ecs.NewWorld()
	uiManager := ui.NewUIManager()
	battleSystem := NewBattleSystem()
	tacticalManager := NewTacticalManager(constants.GridWidth, constants.GridHeight, constants.TileSize) // Grid with tile size from constants
	partyManager := NewPartyManager(constants.MaxPartyMembers)                                           // Max party members from constants
	enemyGroupManager := NewEnemyGroupManager(constants.EnemyGroupRange)                                 // Enemy group range from constants
	tacticalDeployment := NewTacticalDeployment(tacticalManager.Grid)

	game := &Game{
		world:              world,
		activePlayerIndex:  constants.DefaultActivePlayerIndex,
		tabKeyPressed:      false,
		uiManager:          uiManager,
		battleSystem:       battleSystem,
		tacticalManager:    tacticalManager,
		partyManager:       partyManager,
		enemyGroupManager:  enemyGroupManager,
		tacticalDeployment: tacticalDeployment,
		currentMode:        ModeExploration, // Start in exploration mode
	}

	// Set up battle system callbacks
	battleSystem.SetMessageCallback(uiManager.AddMessage)
	battleSystem.SetSwitchPlayerCallback(game.SwitchToNextPlayer)

	// Set up turn-based combat callbacks
	tacticalManager.GetTurnBasedCombat().SetMessageCallback(uiManager.AddMessage)   // For verbose/log messages
	tacticalManager.GetTurnBasedCombat().SetUIMessageCallback(uiManager.AddMessage) // For important UI messages

	return game
}

// AddEntity adds an entity to the game world and manages party system
func (g *Game) AddEntity(entity *ecs.Entity) {
	g.world.AddEntity(entity)

	// Automatically manage party system for player entities
	if entity.HasTag("player") && entity.RPGStats() != nil {
		// Add to party members
		g.partyManager.AddPartyMember(entity)

		// Set as party leader if this is the first player
		if g.partyManager.GetPartyLeader() == nil {
			g.partyManager.SetPartyLeader(entity)
		}
	}
}

// RemoveEntity removes an entity from the game world
func (g *Game) RemoveEntity(entity *ecs.Entity) {
	g.world.RemoveEntity(entity)
}

// SetAttackAnimationDuration configures how long attack animations should last
func (g *Game) SetAttackAnimationDuration(duration time.Duration) {
	g.battleSystem.SetAttackAnimationDuration(duration)
}

// GetCurrentMode returns the current game mode
func (g *Game) GetCurrentMode() GameMode {
	return g.currentMode
}

// SwitchToTacticalMode transitions to tactical combat mode with full party deployment
func (g *Game) SwitchToTacticalMode(participants []*ecs.Entity) {
	if g.currentMode == ModeTactical {
		return // Already in tactical mode
	}

	g.currentMode = ModeTactical

	// Deploy full party instead of just the leader
	partyMembers := g.partyManager.GetPartyForTactical()

	// Get enemy entities from participants
	var enemies []*ecs.Entity
	for _, participant := range participants {
		if participant.RPGStats() != nil && participant.HasTag("enemy") {
			enemies = append(enemies, participant)
		}
	}

	// Create full tactical deployment
	allParticipants := append(partyMembers, enemies...)

	// Clear grid occupancy state before deployment to ensure clean state
	g.clearGridOccupancy()

	// Debug: Check unit positions before deployment
	logger.Debug("Unit positions before deployment:")
	for i, member := range partyMembers {
		if transform := member.Transform(); transform != nil {
			currentGridPos := g.worldToGridPos(transform.X, transform.Y)
			logger.Debug("Unit %s (index %d) - World: (%.1f,%.1f) -> Grid: (%d,%d)",
				member.GetID(), i, transform.X, transform.Y, currentGridPos.X, currentGridPos.Y)
		}
	}

	// Deploy entities to tactical grid
	g.tacticalDeployment.DeployParty(partyMembers)
	g.tacticalDeployment.DeployEnemies(enemies)

	// Debug: Check unit positions after deployment
	logger.Debug("Unit positions after deployment:")
	for i, member := range partyMembers {
		if transform := member.Transform(); transform != nil {
			currentGridPos := g.worldToGridPos(transform.X, transform.Y)
			logger.Debug("Unit %s (index %d) - World: (%.1f,%.1f) -> Grid: (%d,%d)",
				member.GetID(), i, transform.X, transform.Y, currentGridPos.X, currentGridPos.Y)
		}
	}

	// Debug: Validate grid state after deployment
	g.validateGridState()

	g.tacticalManager.StartTacticalCombat(allParticipants)
	g.uiManager.AddMessage("Entering tactical combat mode!")

	// Highlight movement range for initial active player
	activePlayer := g.GetActivePlayer()
	if activePlayer != nil {
		g.tacticalManager.HighlightMovementRangeForPlayer(activePlayer)
	}
}

// SwitchToExplorationMode transitions back to exploration mode
func (g *Game) SwitchToExplorationMode() {
	if g.currentMode == ModeExploration {
		return // Already in exploration mode
	}

	g.currentMode = ModeExploration
	g.tacticalManager.EndTacticalCombat()

	// Restore exploration positions for all entities
	g.restoreExplorationPositions()

	g.uiManager.AddMessage("Returning to exploration mode.")
}

// restoreExplorationPositions repositions entities for exploration mode
func (g *Game) restoreExplorationPositions() {
	partyLeader := g.partyManager.GetPartyLeader()
	if partyLeader == nil || partyLeader.Transform() == nil {
		return
	}

	leaderTransform := partyLeader.Transform()

	// Position all party members at or near the leader's position
	partyMembers := g.partyManager.GetPartyForTactical()
	for i, member := range partyMembers {
		if member != partyLeader && member.Transform() != nil {
			memberTransform := member.Transform()
			// Position party members near the leader (slightly offset to avoid exact overlap)
			offset := float64(i * 5) // Small offset for each member
			memberTransform.X = leaderTransform.X + offset
			memberTransform.Y = leaderTransform.Y + offset
		}
	}

	// Restore enemies to reasonable exploration positions
	// (They should be positioned where they were before tactical mode,
	//  but for simplicity, we'll spread them around the current area)
	enemies := make([]*ecs.Entity, 0)
	for _, entity := range g.world.GetEntities() {
		if entity.HasTag("enemy") && entity.Transform() != nil {
			enemies = append(enemies, entity)
		}
	}

	// Spread enemies around the leader area in exploration mode
	for i, enemy := range enemies {
		enemyTransform := enemy.Transform()
		// Position enemies at various offsets from the leader
		distance := 100.0 + float64(i*20)                   // Different distances
		enemyTransform.X = leaderTransform.X + distance*1.0 // Simple positioning
		enemyTransform.Y = leaderTransform.Y + distance*0.5 + float64(i*30)

		// Ensure they stay within reasonable bounds (screen area)
		if enemyTransform.X < 50 {
			enemyTransform.X = 50
		} else if enemyTransform.X > 700 {
			enemyTransform.X = 700
		}
		if enemyTransform.Y < 150 {
			enemyTransform.Y = 150
		} else if enemyTransform.Y > 450 {
			enemyTransform.Y = 450
		}
	}
}

// IsTacticalMode returns true if currently in tactical combat mode
func (g *Game) IsTacticalMode() bool {
	return g.currentMode == ModeTactical
}

// getAllCombatParticipants returns all entities with RPG stats (players and enemies)
func (g *Game) getAllCombatParticipants() []*ecs.Entity {
	participants := make([]*ecs.Entity, 0)
	for _, entity := range g.world.GetEntities() {
		if entity.RPGStats() != nil {
			participants = append(participants, entity)
		}
	}
	return participants
}

// getNearbyEnemies returns enemies within the specified distance of the player
func (g *Game) getNearbyEnemies(player *ecs.Entity, distance float64) []*ecs.Entity {
	enemies := make([]*ecs.Entity, 0)
	playerT := player.Transform()
	if playerT == nil {
		return enemies
	}

	for _, entity := range g.world.GetEntities() {
		// Skip the player itself and entities without RPG stats
		if entity == player || entity.RPGStats() == nil {
			continue
		}

		// Skip other players - only get enemies
		if entity.HasTag(ecs.TagPlayer) {
			continue
		}

		entityT := entity.Transform()
		if entityT == nil {
			continue
		}

		// Calculate distance
		dx := playerT.X - entityT.X
		dy := playerT.Y - entityT.Y
		dist := dx*dx + dy*dy // Square distance (avoid sqrt for performance)

		if dist <= distance*distance {
			enemies = append(enemies, entity)
		}
	}

	return enemies
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
	g.uiManager.AddMessage("Press C near enemies for tactical combat")
	g.uiManager.AddMessage("In tactical mode: R to reset movement, ESC to exit")
	g.uiManager.AddMessage("Move back to previous positions to recover movement!")

	// Add message about the current active player
	activePlayer := g.GetActivePlayer()
	if activePlayer != nil {
		stats := activePlayer.RPGStats()
		if stats != nil {
			initMsg := fmt.Sprintf("Starting as %s (%s Level %d - %d movement)",
				stats.Name, stats.Job.String(), stats.Level, stats.MoveRange)
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
	if g.currentMode == ModeExploration {
		// In exploration mode, always return the party leader
		return g.partyManager.GetPartyLeader()
	} else {
		// In tactical mode, return active party member
		partyMembers := g.partyManager.GetPartyForTactical()
		if len(partyMembers) == 0 {
			return nil
		}
		if g.activePlayerIndex >= len(partyMembers) {
			g.activePlayerIndex = 0 // Reset if index is out of bounds
		}
		return partyMembers[g.activePlayerIndex]
	}
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

		// Update movement range highlight in tactical mode
		if g.currentMode == ModeTactical {
			g.tacticalManager.HighlightMovementRangeForPlayer(activePlayer)
		}
	}
}

func (g *Game) Update() error {
	// Update based on current game mode
	switch g.currentMode {
	case ModeExploration:
		return g.updateExploration()
	case ModeTactical:
		return g.updateTactical()
	}
	return nil
}

// updateExploration handles exploration mode updates (your current system)
func (g *Game) updateExploration() error {
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
		speed := constants.PlayerSpeed
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

		// Constrain player movement to game world boundaries
		if isMoving {
			// Constrain to game world bounds using constants
			if playerT.X < constants.GameWorldLeft {
				playerT.X = constants.GameWorldLeft
			} else if playerT.X > constants.GameWorldRight-constants.EntityWidth { // Account for player sprite width
				playerT.X = constants.GameWorldRight - constants.EntityWidth
			}
			if playerT.Y < constants.GameWorldTop {
				playerT.Y = constants.GameWorldTop
			} else if playerT.Y > constants.GameWorldBottom-constants.EntityHeight { // Account for player sprite height
				playerT.Y = constants.GameWorldBottom - constants.EntityHeight
			}
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
					// Enemy collision - start tactical combat with all participants
					participants := g.getAllCombatParticipants()
					g.SwitchToTacticalMode(participants)
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

		// Enhanced tactical combat triggers

		// Option 1: Manual tactical mode (T key) - for testing and manual control
		if ebiten.IsKeyPressed(ebiten.KeyT) {
			participants := g.getAllCombatParticipants()
			if len(participants) > 1 {
				g.SwitchToTacticalMode(participants)
			}
		}

		// Option 2: Smart tactical trigger (Spacebar) - when enemies are nearby
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			nearbyEnemies := g.getNearbyEnemies(activePlayer, 100.0) // Within 100 pixels
			if len(nearbyEnemies) > 0 {
				// Always include ALL combat participants, not just nearby ones
				participants := g.getAllCombatParticipants()
				g.SwitchToTacticalMode(participants)
			}
		}
	}

	return nil
}

// updateTactical handles tactical combat mode updates
func (g *Game) updateTactical() error {
	// Update tactical manager
	g.tacticalManager.Update()

	// Handle tactical input
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.SwitchToExplorationMode()
	}

	// Handle TAB key for player switching in tactical mode
	if ebiten.IsKeyPressed(ebiten.KeyTab) {
		if !g.tabKeyPressed {
			g.SwitchToNextTacticalPlayer()
			g.tabKeyPressed = true
		}
	} else {
		g.tabKeyPressed = false
	}

	// Handle R key to reset all movement (for testing/turn management)
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.tacticalManager.ResetAllMovement()
		g.uiManager.AddMessage("All movement reset!")
		// Update movement range for current player
		activePlayer := g.GetActivePlayer()
		if activePlayer != nil {
			g.tacticalManager.HighlightMovementRangeForPlayer(activePlayer)
		}
	}

	// Handle V key to toggle verbose logging
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		logger.SetVerbose(!logger.Verbose)
		if logger.Verbose {
			g.uiManager.AddMessage("Verbose logging enabled")
			logger.Info("Verbose logging enabled")
		} else {
			g.uiManager.AddMessage("Verbose logging disabled")
			logger.Info("Verbose logging disabled")
		}
	}

	// Handle mouse input for tile selection and movement
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		screenX, screenY := float64(x), float64(y)
		offsetX, offsetY := constants.GridOffsetX, constants.GridOffsetY

		if gridPos, valid := g.tacticalManager.GetTileAtScreenPos(screenX, screenY, offsetX, offsetY); valid {
			g.handleTacticalClick(gridPos)
		}
	}

	// Handle arrow keys for grid-based movement in tactical mode
	activePlayer := g.GetActivePlayer()
	if activePlayer != nil {
		g.handleTacticalArrowKeys(activePlayer)
	} else {
		// Debug: Check if arrow keys are pressed when no active player
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
			inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			logger.Debug("Arrow key pressed but no active player found")
		}
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

	// Prepare UI data based on current mode
	var uiMode ui.GameMode
	var partyStats []*components.RPGStatsComponent
	var gridPosition string

	if g.currentMode == ModeExploration {
		uiMode = ui.ModeExploration
		// Get all party members' stats for exploration view
		partyMembers := g.partyManager.GetPartyForTactical()
		for _, member := range partyMembers {
			if member != nil && member.RPGStats() != nil {
				partyStats = append(partyStats, member.RPGStats())
			}
		}
	} else {
		uiMode = ui.ModeTactical
		// Get grid position for tactical view
		if activePlayer != nil {
			if transform := activePlayer.Transform(); transform != nil {
				gridPos := g.worldToGridPos(transform.X, transform.Y)
				gridPosition = fmt.Sprintf("(%d, %d)", gridPos.X, gridPos.Y)
			} else {
				gridPosition = "Unknown"
			}
		}
	}

	// Draw top panel with mode-specific info
	g.uiManager.DrawTopPanel(screen, activePlayerStats, uiMode, partyStats, gridPosition)

	// Draw game world background
	g.uiManager.DrawGameWorldBackground(screen)

	// Draw game entities based on current mode
	for _, entity := range g.world.GetEntities() {
		transform := entity.Transform()
		if transform == nil {
			continue // Skip entities without a transform
		}

		// In exploration mode, only show party leader and enemies/objects
		if g.currentMode == ModeExploration {
			if entity.HasTag("player") {
				// Only show the party leader in exploration mode
				if entity != g.partyManager.GetPartyLeader() {
					continue // Skip non-leader party members
				}
			}
		} else if g.currentMode == ModeTactical {
			// In tactical mode, all entities (including tactical participants)
			// should be rendered at their current positions
			// The tactical system will have repositioned tactical participants
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

	// Draw tactical grid overlay if in tactical mode
	if g.IsTacticalMode() {
		// Update offset to match game world Y position (110px panel + 2px separator = 112px)
		offsetX, offsetY := constants.GridOffsetX, constants.GridOffsetY
		g.tacticalManager.DrawGrid(screen, offsetX, offsetY)

		// Draw combat UI for turn-based combat
		g.tacticalManager.DrawCombatUI(screen)
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

// clearGridOccupancy resets all grid tiles to unoccupied state
func (g *Game) clearGridOccupancy() {
	for x := 0; x < g.tacticalManager.Grid.Width; x++ {
		for y := 0; y < g.tacticalManager.Grid.Height; y++ {
			pos := tactical.GridPos{X: x, Y: y}
			g.tacticalManager.Grid.SetOccupied(pos, false, "")
		}
	}
}

// validateGridState checks for grid state inconsistencies
func (g *Game) validateGridState() {
	logger.Debug("Validating grid state...")
	occupiedTiles := 0
	unitPositions := make(map[string][]tactical.GridPos) // Map of unit ID to positions

	for x := 0; x < g.tacticalManager.Grid.Width; x++ {
		for y := 0; y < g.tacticalManager.Grid.Height; y++ {
			pos := tactical.GridPos{X: x, Y: y}
			tile := g.tacticalManager.Grid.GetTile(pos)
			if tile != nil && tile.Occupied {
				occupiedTiles++
				logger.Debug("Grid pos (%d,%d) occupied by unit %s", x, y, tile.UnitID)
				unitPositions[tile.UnitID] = append(unitPositions[tile.UnitID], pos)
			}
		}
	}
	logger.Debug("Total occupied tiles: %d", occupiedTiles)

	// Check for units occupying multiple positions
	for unitID, positions := range unitPositions {
		if len(positions) > 1 {
			fmt.Printf("WARNING: Unit %s occupies multiple positions: %v\n", unitID, positions)
			// Fix: Clear all positions except the first one
			for i := 1; i < len(positions); i++ {
				g.tacticalManager.Grid.SetOccupied(positions[i], false, "")
				logger.Debug("Cleared duplicate position (%d,%d) for unit %s",
					positions[i].X, positions[i].Y, unitID)
			}
		}
	}
}

// clearUnitFromAllGridPositions removes a unit from all grid positions it might occupy
func (g *Game) clearUnitFromAllGridPositions(unitID string) {
	clearedCount := 0
	for x := 0; x < g.tacticalManager.Grid.Width; x++ {
		for y := 0; y < g.tacticalManager.Grid.Height; y++ {
			pos := tactical.GridPos{X: x, Y: y}
			tile := g.tacticalManager.Grid.GetTile(pos)
			if tile != nil && tile.Occupied && tile.UnitID == unitID {
				g.tacticalManager.Grid.SetOccupied(pos, false, "")
				logger.Debug("Cleared unit %s from grid pos (%d,%d)", unitID, x, y)
				clearedCount++
			}
		}
	}
	if clearedCount > 1 {
		fmt.Printf("WARNING: Unit %s was occupying %d positions!\n", unitID, clearedCount)
	}
}

// worldToGridPos converts world coordinates to grid position - exact inverse of GridToWorld
func (g *Game) worldToGridPos(worldX, worldY float64) tactical.GridPos {
	offsetX, offsetY := constants.GridOffsetX, constants.GridOffsetY
	tileSize := float64(g.tacticalManager.Grid.TileSize)

	// Remove offset and convert to grid coordinates
	// This is the exact inverse of: worldX = gridX * tileSize + offsetX
	gridX := int((worldX - offsetX) / tileSize)
	gridY := int((worldY - offsetY) / tileSize)

	return tactical.GridPos{X: gridX, Y: gridY}
}

// SwitchToNextTacticalPlayer switches to the next player in tactical mode
func (g *Game) SwitchToNextTacticalPlayer() {
	// Only switch between party members in tactical mode
	if g.currentMode != ModeTactical {
		return
	}

	partyMembers := g.partyManager.GetPartyForTactical()
	if len(partyMembers) <= 1 {
		return // No other players to switch to
	}

	// Find current active player in party members
	currentActive := g.GetActivePlayer()
	currentIndex := -1
	for i, member := range partyMembers {
		if member == currentActive {
			currentIndex = i
			break
		}
	}

	// Switch to next party member
	if currentIndex >= 0 {
		nextIndex := (currentIndex + 1) % len(partyMembers)
		g.activePlayerIndex = nextIndex
		g.uiManager.AddMessage(fmt.Sprintf("Switched to %s", partyMembers[nextIndex].Name))

		// Update the combat manager's active unit
		if g.tacticalManager.UseTurnBasedCombat {
			g.tacticalManager.GetTurnBasedCombat().SetActiveUnit(partyMembers[nextIndex])
		}

		// Highlight movement range for new active player
		g.tacticalManager.HighlightMovementRangeForPlayer(partyMembers[nextIndex])
	}
}

// handleTacticalClick handles mouse clicks in tactical mode (MOVEMENT DISABLED)
func (g *Game) handleTacticalClick(gridPos tactical.GridPos) {
	// Debug: Show click information
	g.uiManager.AddMessage(fmt.Sprintf("Clicked on tile (%d, %d)", gridPos.X, gridPos.Y))

	// Select the clicked tile
	g.tacticalManager.SelectTile(gridPos)

	// Check if there's a unit at this position
	tile := g.tacticalManager.Grid.GetTile(gridPos)
	if tile != nil && tile.Occupied {
		g.uiManager.AddMessage(fmt.Sprintf("Selected tile (%d, %d) - Unit: %s", gridPos.X, gridPos.Y, tile.UnitID))
	} else {
		g.uiManager.AddMessage(fmt.Sprintf("Selected empty tile (%d, %d) - Use arrow keys to move", gridPos.X, gridPos.Y))
	}

	// MOUSE MOVEMENT DISABLED - Use arrow keys only for movement
} // handleTacticalArrowKeys handles arrow key movement in tactical mode
func (g *Game) handleTacticalArrowKeys(player *ecs.Entity) {
	if player.Transform() == nil {
		logger.Debug("Player has no transform in arrow key handler")
		return
	}

	// Get current grid position using consistent conversion
	transform := player.Transform()
	currentPos := g.worldToGridPos(transform.X, transform.Y)

	var newPos tactical.GridPos
	moved := false

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		newPos = tactical.GridPos{X: currentPos.X, Y: currentPos.Y - 1}
		moved = true
		logger.Debug("Arrow UP pressed - moving from (%d,%d) to (%d,%d)", currentPos.X, currentPos.Y, newPos.X, newPos.Y)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		newPos = tactical.GridPos{X: currentPos.X, Y: currentPos.Y + 1}
		moved = true
		logger.Debug("Arrow DOWN pressed - moving from (%d,%d) to (%d,%d)", currentPos.X, currentPos.Y, newPos.X, newPos.Y)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		newPos = tactical.GridPos{X: currentPos.X - 1, Y: currentPos.Y}
		moved = true
		logger.Debug("Arrow LEFT pressed - moving from (%d,%d) to (%d,%d)", currentPos.X, currentPos.Y, newPos.X, newPos.Y)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		newPos = tactical.GridPos{X: currentPos.X + 1, Y: currentPos.Y}
		moved = true
		logger.Debug("Arrow RIGHT pressed - moving from (%d,%d) to (%d,%d)", currentPos.X, currentPos.Y, newPos.X, newPos.Y)
	}

	if moved {
		logger.Debug("About to call tryMovePlayerToTile for player %s", player.GetID())

		// Debug: Show current grid state before attempting move
		logger.Debug("Current grid occupancy:")
		occupiedCount := 0
		for x := 0; x < g.tacticalManager.Grid.Width && x < 5; x++ { // Limit to first 5 columns for readability
			for y := 0; y < g.tacticalManager.Grid.Height && y < 5; y++ { // Limit to first 5 rows
				pos := tactical.GridPos{X: x, Y: y}
				tile := g.tacticalManager.Grid.GetTile(pos)
				if tile != nil && tile.Occupied {
					logger.Debug("  (%d,%d): %s", x, y, tile.UnitID)
					occupiedCount++
				}
			}
		}
		logger.Debug("Total occupied tiles (partial): %d", occupiedCount)

		g.tryMovePlayerToTile(player, newPos)
	}
}

// tryMovePlayerToTile attempts to move a player to the specified grid tile
func (g *Game) tryMovePlayerToTile(player *ecs.Entity, gridPos tactical.GridPos) {
	// Update player position
	transform := player.Transform()
	if transform == nil {
		g.uiManager.AddMessage("Player has no transform component")
		return
	}

	// Calculate current grid position using consistent conversion
	currentPos := g.worldToGridPos(transform.X, transform.Y)

	// Debug: Show movement attempt
	logger.Debug("Moving from (%d, %d) to (%d, %d)",
		currentPos.X, currentPos.Y, gridPos.X, gridPos.Y)

	// Debug: Check current position state
	currentTile := g.tacticalManager.Grid.GetTile(currentPos)
	if currentTile != nil {
		logger.Debug("Current pos (%d,%d) - Occupied: %t, UnitID: %s",
			currentPos.X, currentPos.Y, currentTile.Occupied, currentTile.UnitID)
	}

	// Debug: Check target position state
	targetTile := g.tacticalManager.Grid.GetTile(gridPos)
	if targetTile != nil {
		logger.Debug("Target pos (%d,%d) - Occupied: %t, UnitID: %s",
			gridPos.X, gridPos.Y, targetTile.Occupied, targetTile.UnitID)
	}

	// Check if we're trying to move to the same position
	if currentPos.X == gridPos.X && currentPos.Y == gridPos.Y {
		g.uiManager.AddMessage("Already at that position")
		return
	}

	// Check if the target position is valid
	if !g.tacticalManager.Grid.IsValidPosition(gridPos) {
		g.uiManager.AddMessage("Invalid position")
		return
	}

	// Check if the move is within the player's remaining movement (unless it's an undo move)
	stats := player.RPGStats()
	if stats != nil {
		distance := g.tacticalManager.Grid.CalculateDistance(currentPos, gridPos)

		// First check if this would be an undo move
		isUndo := stats.IsUndoMove(gridPos.X, gridPos.Y)

		// Only apply movement restriction if it's NOT an undo move
		if !isUndo && !stats.CanMove(distance) {
			g.uiManager.AddMessage(fmt.Sprintf("Not enough movement! %s has %d moves left (need %d)",
				stats.Job.String(), stats.MovesRemaining, distance))
			return
		}

		if isUndo {
			logger.Debug("Undo move detected - Distance: %d, allowing move despite %d remaining", distance, stats.MovesRemaining)
		} else {
			logger.Debug("Normal move - Distance: %d, Moves Remaining: %d", distance, stats.MovesRemaining)
		}
	}

	// Debug: Show which unit is trying to move
	logger.Debug("Unit %s attempting to move", player.GetID())

	// Clear ALL positions occupied by this unit FIRST - prevents multi-occupancy
	g.clearUnitFromAllGridPositions(player.GetID())

	// Also log what we thought was the current position for debugging
	oldTile := g.tacticalManager.Grid.GetTile(currentPos)
	if oldTile != nil {
		logger.Debug("Expected current pos (%d,%d) - Occupied: %t, UnitID: %s",
			currentPos.X, currentPos.Y, oldTile.Occupied, oldTile.UnitID)
	}

	// Debug: Check target tile state after clearing
	targetTileAfterClear := g.tacticalManager.Grid.GetTile(gridPos)
	if targetTileAfterClear != nil {
		logger.Debug("Target pos (%d,%d) after clearing - Occupied: %t, UnitID: %s",
			gridPos.X, gridPos.Y, targetTileAfterClear.Occupied, targetTileAfterClear.UnitID)
	}

	// Now check if target position is passable (after clearing our old position)
	if !g.tacticalManager.Grid.IsPassable(gridPos) {
		// Debug: Check what's occupying the tile
		tile := g.tacticalManager.Grid.GetTile(gridPos)
		if tile != nil && tile.Occupied {
			g.uiManager.AddMessage(fmt.Sprintf("Cannot move to (%d,%d) - occupied by unit %s",
				gridPos.X, gridPos.Y, tile.UnitID))
			// Restore our position at the expected current location since move failed
			g.tacticalManager.Grid.SetOccupied(currentPos, true, player.GetID())
		} else {
			g.uiManager.AddMessage("Cannot move to occupied tile")
		}
		return
	}

	// Set new position
	offsetX, offsetY := constants.GridOffsetX, constants.GridOffsetY
	worldX, worldY := g.tacticalManager.Grid.GridToWorld(gridPos)
	transform.X = worldX + offsetX
	transform.Y = worldY + offsetY
	g.tacticalManager.Grid.SetOccupied(gridPos, true, player.GetID())

	// Log the actual position update for debugging
	logger.Debug("POSITION UPDATE: Unit %s moved to World(%.1f,%.1f) Grid(%d,%d)",
		player.GetID(), transform.X, transform.Y, gridPos.X, gridPos.Y)

	// Handle movement consumption and tracking
	if stats != nil {
		distance := g.tacticalManager.Grid.CalculateDistance(currentPos, gridPos)

		// Try to undo move if returning to a previous position
		if undoSuccessful, recoveredMoves := stats.TryUndoMove(gridPos.X, gridPos.Y); undoSuccessful {
			g.uiManager.AddMessage(fmt.Sprintf("%s returned to previous position - recovered %d moves (%d moves left)",
				player.Name, recoveredMoves, stats.MovesRemaining))
			logger.Debug("Undo successful! Recovered %d moves, %d remaining", recoveredMoves, stats.MovesRemaining)
		} else {
			// Normal move - consume movement and record it
			stats.ConsumeMovement(distance)
			stats.RecordMove(currentPos.X, currentPos.Y, gridPos.X, gridPos.Y, distance)
			g.uiManager.AddMessage(fmt.Sprintf("%s moved to (%d, %d) - %d moves left",
				player.Name, gridPos.X, gridPos.Y, stats.MovesRemaining))
			logger.Debug("Normal move from (%d,%d) to (%d,%d), cost %d, %d remaining",
				currentPos.X, currentPos.Y, gridPos.X, gridPos.Y, distance, stats.MovesRemaining)
			logger.Debug("Move history: %s", stats.GetMoveHistoryString())
		}
	} else {
		g.uiManager.AddMessage(fmt.Sprintf("%s moved to (%d, %d)", player.Name, gridPos.X, gridPos.Y))
	}

	// Update movement range display
	g.tacticalManager.HighlightMovementRangeForPlayer(player)
}
