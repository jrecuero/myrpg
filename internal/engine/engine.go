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
	"github.com/jrecuero/myrpg/internal/events"
	"github.com/jrecuero/myrpg/internal/gfx"
	"github.com/jrecuero/myrpg/internal/logger"
	"github.com/jrecuero/myrpg/internal/quests"
	"github.com/jrecuero/myrpg/internal/save"
	"github.com/jrecuero/myrpg/internal/skills"
	"github.com/jrecuero/myrpg/internal/tactical"
	"github.com/jrecuero/myrpg/internal/ui"
)

// Game represents the state of the game using an ECS architecture.
type Game struct {
	world              *ecs.World           // The game world containing all entities
	activePlayerIndex  int                  // Index of the currently active player
	tabKeyPressed      bool                 // Track TAB key state to prevent multiple switches
	uiManager          *ui.UIManager        // UI system for panels and messages
	battleSystem       *BattleSystem        // Battle system for combat
	tacticalManager    *TacticalManager     // Tactical combat system
	partyManager       *PartyManager        // Party and team management
	enemyGroupManager  *EnemyGroupManager   // Enemy group formation
	tacticalDeployment *TacticalDeployment  // Unit deployment for tactical combat
	eventManager       *events.EventManager // Event system for interactive world elements
	saveManager        *save.SaveManager    // Save system for game state persistence
	currentMode        GameMode             // Current game mode (exploration/tactical)
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

	// Initialize event system
	eventManager := events.NewEventManager()

	// Initialize save system
	saveManager := save.NewSaveManager("saves") // Save files in "saves" directory

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
		eventManager:       eventManager,
		saveManager:        saveManager,
		currentMode:        ModeExploration, // Start in exploration mode
	}

	// Initialize skills system
	skills.InitializeSkillRegistry()

	// Initialize quest system
	quests.InitializeQuestRegistry()

	// Initialize event system handlers (will be set up after game creation)
	// handlers will be registered in SetupGameEventHandlers() method

	// Set up battle system callbacks
	battleSystem.SetMessageCallback(uiManager.AddMessage)
	battleSystem.SetSwitchPlayerCallback(game.SwitchToNextPlayer)

	// Set up turn-based combat callbacks
	tacticalManager.GetTurnBasedCombat().SetMessageCallback(uiManager.AddMessage)   // For verbose/log messages
	tacticalManager.GetTurnBasedCombat().SetUIMessageCallback(uiManager.AddMessage) // For important UI messages

	// Set up state change callback to refresh movement highlighting when player turn starts
	tacticalManager.GetTurnBasedCombat().SetStateChangeCallback(func(phase tactical.CombatPhase) {
		if phase == tactical.CombatPhaseTeamTurn {
			// Check if it's now a player turn and refresh movement highlighting for the engine's active player
			if activePlayer := game.GetActivePlayer(); activePlayer != nil {
				if activePlayer.HasTag("player") {
					tacticalManager.HighlightMovementRangeForPlayer(activePlayer)
				}
			}
		}
	})

	// Initialize item system
	components.InitializeItemSystem()

	return game
}

// SetupGameEventHandlers initializes the game-specific event handlers
func (g *Game) SetupGameEventHandlers() {
	handlers := g.CreateGameEventHandlers()
	for eventType, handler := range handlers {
		g.eventManager.RegisterHandler(eventType, handler)
	}
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

		// Set player entity reference in event manager (handles multiple calls gracefully)
		g.eventManager.SetPlayer(entity)
	}

	// Register event entities with the event manager
	if entity.HasTag("event") || entity.Event() != nil {
		g.eventManager.RegisterEntity(entity)
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

// SaveEventState saves the current event state to disk
func (g *Game) SaveEventState() error {
	if g.saveManager == nil {
		return fmt.Errorf("save manager not initialized")
	}
	return g.saveManager.SaveEventState(g.world, g.eventManager)
}

// LoadEventState loads event state from disk and applies it to the game
func (g *Game) LoadEventState() error {
	if g.saveManager == nil {
		return fmt.Errorf("save manager not initialized")
	}

	saveData, err := g.saveManager.LoadEventState()
	if err != nil {
		return fmt.Errorf("failed to load event state: %v", err)
	}

	// Apply the loaded state to the world and event manager
	if err := g.saveManager.ApplyEventState(g.world, g.eventManager, saveData); err != nil {
		return fmt.Errorf("failed to apply event state: %v", err)
	}

	// Load completed events into the event manager
	g.eventManager.LoadCompletedEvents(saveData.CompletedEvents)

	return nil
}

// ClearEventState resets all events to their default state (for new game)
func (g *Game) ClearEventState() error {
	if g.saveManager == nil {
		return fmt.Errorf("save manager not initialized")
	}

	// Clear save data
	if err := g.saveManager.ClearEventState(g.world); err != nil {
		return fmt.Errorf("failed to clear save state: %v", err)
	}

	// Clear event manager tracking
	g.eventManager.ClearCompletedEventsForReset()

	return nil
}

// GetSaveManager returns the save manager for external use
func (g *Game) GetSaveManager() *save.SaveManager {
	return g.saveManager
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
	g.uiManager.AddMessage("Press I for inventory, K for skills, J for quests, Q for equipment, H for help, C near enemies for combat")
	g.uiManager.AddMessage("In inventory: Right-click equipment to equip, drag items to move")
	g.uiManager.AddMessage("In tactical mode: A to attack, E to end turn, R to reset movement, ESC to exit")
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
	// Update UI manager (handles popups and other UI interactions)
	uiInputResult := g.uiManager.Update()

	// Test popup widgets (demo/testing)
	if inpututil.IsKeyJustPressed(ebiten.KeyP) && !g.uiManager.IsPopupVisible() {
		g.showTestPopup()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) && !g.uiManager.IsPopupVisible() {
		g.toggleInventory()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) && !g.uiManager.IsPopupVisible() {
		g.toggleSkills()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) && !g.uiManager.IsPopupVisible() {
		g.toggleQuestJournal()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) && !g.uiManager.IsPopupVisible() {
		g.showCharacterStats()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) && !g.uiManager.IsPopupVisible() {
		g.showEquipment()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) && !g.uiManager.IsPopupVisible() {
		g.showDialog()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyH) && !g.uiManager.IsPopupVisible() {
		g.showTestInfoPopup()
	}

	// Block game input processing when popup is visible OR when UI consumed ESC
	if g.uiManager.IsPopupVisible() || uiInputResult.EscConsumed {
		return nil // Only process UI input, skip game logic
	}

	// Update based on current game mode
	switch g.currentMode {
	case ModeExploration:
		return g.updateExploration()
	case ModeTactical:
		return g.updateTactical(uiInputResult)
	}
	return nil
}

// updateExploration handles exploration mode updates (your current system)
func (g *Game) updateExploration() error {
	// Update battle system first
	g.battleSystem.Update()

	// Update event system
	if g.currentMode == ModeExploration {
		g.eventManager.SetGameMode(components.GameModeExploration)
	} else {
		g.eventManager.SetGameMode(components.GameModeTactical)
	}

	// Set the current active player as the player entity for events
	activePlayer := g.GetActivePlayer()
	if activePlayer != nil {
		g.eventManager.SetPlayer(activePlayer)
	}

	g.eventManager.Update(1.0 / 60.0) // Assuming 60 FPS delta time

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
			// Skip hidden party members in exploration mode
			if g.currentMode == ModeExploration && entity.HasTag("player") {
				if entity != g.partyManager.GetPartyLeader() {
					continue // Skip non-leader party members that aren't visible
				}
			}

			// Check collision type
			if CheckCollision(playerT.Bounds(), entity.Transform().Bounds()) {
				// Check if this is an event entity first - don't block movement but let EventManager handle it
				if entity.Event() != nil {
					continue // Event entities are walkable, EventManager will handle triggering
				}

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
func (g *Game) updateTactical(uiInputResult ui.InputResult) error {
	// Update tactical manager
	g.tacticalManager.Update()

	// Handle tactical input
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
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

	// Handle A key to attack in tactical mode
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		activePlayer := g.GetActivePlayer()
		if activePlayer != nil {

			// Find adjacent enemies to attack
			target := g.tacticalManager.GetTurnBasedCombat().FindAdjacentTarget(activePlayer, components.TeamEnemy)
			if target != nil {
				// Create and execute attack action
				action := &tactical.CombatAction{
					Type:    tactical.ActionAttack,
					Actor:   activePlayer,
					Target:  target,
					APCost:  constants.AttackAPCost,
					Message: fmt.Sprintf("%s attacks %s", activePlayer.RPGStats().Name, target.RPGStats().Name),
				}

				err := g.tacticalManager.GetTurnBasedCombat().ExecuteAction(action)
				if err != nil {
					g.uiManager.AddMessage(fmt.Sprintf("Attack failed: %v", err))
					logger.Error("Failed to execute attack action: %v", err)
				} else {
					g.uiManager.AddMessage(fmt.Sprintf("%s attacks %s!", activePlayer.RPGStats().Name, target.RPGStats().Name))
					logger.Info("Player %s attacked %s", activePlayer.RPGStats().Name, target.RPGStats().Name)
				}
			} else {
				g.uiManager.AddMessage("No adjacent enemies to attack")
				logger.Info("Player %s tried to attack but no adjacent enemies found", activePlayer.RPGStats().Name)
			}
		} else {
			g.uiManager.AddMessage("No active player to attack with")
		}
	}

	// Handle E key to end turn in tactical mode
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		activePlayer := g.GetActivePlayer()
		if activePlayer != nil {
			// Create and execute end turn action
			action, err := g.tacticalManager.GetTurnBasedCombat().CreateEndTurnAction(activePlayer)
			if err != nil {
				g.uiManager.AddMessage(fmt.Sprintf("Failed to end turn: %v", err))
				logger.Error("Failed to create end turn action: %v", err)
			} else {
				err = g.tacticalManager.GetTurnBasedCombat().ExecuteAction(action)
				if err != nil {
					g.uiManager.AddMessage(fmt.Sprintf("Failed to execute end turn: %v", err))
					logger.Error("Failed to execute end turn action: %v", err)
				} else {
					g.uiManager.AddMessage(fmt.Sprintf("%s ended their turn", activePlayer.RPGStats().Name))
					logger.Info("Player %s ended their turn", activePlayer.RPGStats().Name)
				}
			}
		} else {
			g.uiManager.AddMessage("No active player to end turn for")
		}
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
	// Only process mouse clicks if UI didn't consume them
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !uiInputResult.MouseConsumed {
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
		if stats := activePlayer.RPGStats(); stats != nil {
			activePlayerStats = stats
		}
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
			if member != nil {
				if stats := member.RPGStats(); stats != nil {
					partyStats = append(partyStats, stats)
				}
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

		// Check if this is a background entity that needs clipping
		isBackground := entity.HasTag("background") || entity.Name == "Background"

		// Check for animated sprite first, fallback to static sprite
		animationC := entity.Animation()
		if animationC != nil {
			// Update animation
			animationC.Update()

			// Draw current animation frame
			currentSprite := animationC.GetCurrentSprite()
			if currentSprite != nil {
				if isBackground {
					// Clip background sprites to game world area
					gfx.DrawSpriteClipped(screen, currentSprite,
						transform.X+animationC.OffsetX,
						transform.Y+animationC.OffsetY,
						animationC.Scale,
						0, constants.GameWorldY, constants.BackgroundWidth, constants.GameWorldHeight)
				} else {
					gfx.DrawSprite(screen, currentSprite,
						transform.X+animationC.OffsetX,
						transform.Y+animationC.OffsetY,
						animationC.Scale)
				}
			}
		} else {
			// Fallback to static sprite
			spriteC := entity.Sprite()
			if spriteC != nil {
				// Check if this is an event with mode restrictions
				if eventC := entity.Event(); eventC != nil {
					// Check if event is active in current game mode
					var eventGameMode components.GameMode
					if g.currentMode == ModeExploration {
						eventGameMode = components.GameModeExploration
					} else {
						eventGameMode = components.GameModeTactical
					}

					// Only draw events that are active in the current game mode
					if !eventC.IsActiveInMode(eventGameMode) {
						continue // Skip this event entity
					}
				}

				if isBackground {
					// Clip background sprites to game world area
					gfx.DrawSpriteClipped(screen, spriteC.Sprite,
						transform.X, transform.Y, spriteC.Scale,
						0, constants.GameWorldY, constants.BackgroundWidth, constants.GameWorldHeight)
				} else {
					gfx.DrawSprite(screen, spriteC.Sprite, transform.X, transform.Y, spriteC.Scale)
				}
			} else {
				// Check if this is a visible event entity without a sprite
				if eventC := entity.Event(); eventC != nil && eventC.IsVisible() {
					// Check if event is active in current game mode
					var eventGameMode components.GameMode
					if g.currentMode == ModeExploration {
						eventGameMode = components.GameModeExploration
					} else {
						eventGameMode = components.GameModeTactical
					}

					// Only draw events that are active in the current game mode
					if eventC.IsActiveInMode(eventGameMode) {
						// Draw colored square for events without sprites
						color := color.RGBA{
							eventC.FallbackColor[0],
							eventC.FallbackColor[1],
							eventC.FallbackColor[2],
							255,
						}
						vector.FillRect(screen,
							float32(transform.X), float32(transform.Y),
							float32(transform.Width), float32(transform.Height),
							color, false)
					}
				}
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

	// Draw popups on top of everything else
	g.uiManager.DrawPopups(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
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
			logger.Warn("WARNING: Unit %s occupies multiple positions: %v\n", unitID, positions)
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
		logger.Warn("WARNING: Unit %s was occupying %d positions!\n", unitID, clearedCount)
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

// showTestPopup displays a test popup to demonstrate the popup selection widget
func (g *Game) showTestPopup() {
	options := []string{
		"Attack Enemy",
		"Cast Spell",
		"Use Item",
		"Move Unit",
		"End Turn",
		"View Stats",
		"Cancel Action",
	}

	g.uiManager.ShowSelectionPopup(
		"Combat Actions",
		options,
		func(index int, option string) {
			// Handle selection
			g.uiManager.AddMessage(fmt.Sprintf("Selected: %s (index %d)", option, index))
			logger.Info("Player selected popup option: %s", option)
		},
		func() {
			// Handle cancel
			g.uiManager.AddMessage("Selection cancelled")
			logger.Info("Player cancelled popup selection")
		},
	)
}

// showTestInfoPopup displays a test info popup to demonstrate the popup info widget
func (g *Game) showTestInfoPopup() {
	helpContent := `GAME HELP & INFORMATION

BASIC CONTROLS:
• Arrow Keys: Move player/units
• TAB: Switch between Exploration and Tactical modes
• SPACE: Enter tactical combat (in exploration)
• T: Enter tactical mode manually
• E: End turn (in tactical mode)
• Q: Open equipment menu
• A: Attack (when adjacent to enemy)
• ESC: Exit to desktop

EXPLORATION MODE:
- Move freely around the world
- Search for items and encounters
- Switch to tactical mode for combat

TACTICAL MODE:
- Grid-based movement system
- Turn-based combat mechanics
- Plan strategic unit positioning
- Use attacks and abilities

POPUP WIDGETS:
• P: Test selection popup
• I: Show this help information
• C: Character stats
• E: Equipment
• D: Dialog with NPCs
• ESC: Close any open popup

CHARACTER STATS:
HP: Health Points - When this reaches 0, the character is defeated
MP: Magic Points - Used for casting spells and abilities
STR: Strength - Affects physical damage and carrying capacity
DEF: Defense - Reduces incoming physical damage
AGI: Agility - Affects turn order and evasion chance

COMBAT SYSTEM:
1. Position units strategically on the grid
2. Use terrain and positioning advantages
3. Manage resources (HP/MP) carefully
4. Plan multi-turn strategies

Press ESC to close this help window.`

	g.uiManager.ShowInfoPopup(
		"Game Help & Information",
		helpContent,
		func() {
			// Handle close
			g.uiManager.AddMessage("Help closed")
			logger.Info("Player closed help popup")
		},
	)
}

// showCharacterStats displays the character statistics widget for the active player
func (g *Game) showCharacterStats() {
	activePlayer := g.GetActivePlayer()
	if activePlayer == nil || activePlayer.RPGStats() == nil {
		g.uiManager.AddMessage("No active player to display stats for")
		logger.Info("Attempted to show character stats but no active player available")
		return
	}

	g.uiManager.ShowCharacterStats(activePlayer.RPGStats())
	g.uiManager.AddMessage(fmt.Sprintf("Displaying stats for %s", activePlayer.RPGStats().Name))
	logger.Info("Showing character stats for player: %s", activePlayer.RPGStats().Name)
}

// showEquipment displays the equipment widget for the active player
func (g *Game) showEquipment() {
	activePlayer := g.GetActivePlayer()
	if activePlayer == nil || activePlayer.RPGStats() == nil {
		g.uiManager.AddMessage("No active player to display equipment for")
		logger.Info("Attempted to show equipment but no active player available")
		return
	}

	// Get or create equipment component for the active player
	equipmentComp := activePlayer.Equipment()
	if equipmentComp == nil {
		// Create a new equipment component if none exists
		equipmentComp = components.NewEquipmentComponent()
		// Add equipment component to the player entity so it persists
		activePlayer.AddComponent(ecs.ComponentEquipment, equipmentComp)
		logger.Info("Created and attached new equipment component for player: %s", activePlayer.RPGStats().Name)
	}

	g.uiManager.ShowEquipment(equipmentComp, activePlayer.RPGStats(), activePlayer)
	g.uiManager.AddMessage(fmt.Sprintf("Displaying equipment for %s", activePlayer.RPGStats().Name))
	logger.Info("Showing equipment for player: %s", activePlayer.RPGStats().Name)
}

func (g *Game) showDialog() {
	activePlayer := g.GetActivePlayer()
	if activePlayer == nil || activePlayer.RPGStats() == nil {
		g.uiManager.AddMessage("No active player to show dialog for")
		logger.Info("Attempted to show dialog but no active player available")
		return
	}

	// Initialize dialog with town elder as example
	err := g.uiManager.ShowDialog("assets/dialogs", "characters.json", "town_elder.json", "start")
	if err != nil {
		g.uiManager.AddMessage(fmt.Sprintf("Failed to load dialog: %v", err))
		logger.Error("Failed to show dialog: %v", err)
		return
	}

	g.uiManager.AddMessage("Started dialog with town elder")
	logger.Info("Showing dialog for player: %s", activePlayer.RPGStats().Name)
}

// populatePlayerInventoryWithTestItems adds sample items to a player's inventory using the item registry
func (g *Game) populatePlayerInventoryWithTestItems(player *ecs.Entity) {
	if player == nil || player.Inventory() == nil {
		return
	}

	inventory := player.Inventory()

	// Get the item registry
	registry := components.GlobalItemRegistry
	if registry == nil {
		logger.Error("Item registry not initialized")
		return
	}

	// Create items from the registry
	testItems := []struct {
		itemID   int
		quantity int
	}{
		{1, 1}, // Iron Sword
		{2, 3}, // Health Potion x3
		{3, 2}, // Mana Potion x2
		{4, 1}, // Magic Crystal
	}

	itemsAdded := 0
	for _, testItem := range testItems {
		item := registry.CreateItem(testItem.itemID)
		if item != nil {
			inventory.AddItem(item, testItem.quantity)
			itemsAdded++
		} else {
			logger.Warn("Failed to create item with ID %d from registry", testItem.itemID)
		}
	}

	logger.Info("Added %d test items to %s's inventory using item registry", itemsAdded, player.RPGStats().Name)
}

// PopulateAllPlayerInventoriesWithTestItems adds sample items to all player inventories
func (g *Game) PopulateAllPlayerInventoriesWithTestItems() {
	players := g.world.FindWithTag(ecs.TagPlayer)
	for _, player := range players {
		g.populatePlayerInventoryWithTestItems(player)
	}
	logger.Info("Populated inventories for %d players with test items", len(players))
}

// toggleInventory shows/hides the inventory for the active player
func (g *Game) toggleInventory() {
	activePlayer := g.GetActivePlayer()
	if activePlayer == nil || activePlayer.RPGStats() == nil {
		g.uiManager.AddMessage("No active player to show inventory for")
		logger.Info("Attempted to toggle inventory but no active player available")
		return
	}

	// Check if player has inventory component
	if activePlayer.Inventory() == nil {
		g.uiManager.AddMessage(fmt.Sprintf("%s does not have an inventory", activePlayer.RPGStats().Name))
		logger.Info("Player %s does not have inventory component", activePlayer.RPGStats().Name)
		return
	}

	// Toggle inventory visibility
	if g.uiManager.IsInventoryVisible() {
		g.uiManager.HideInventory()
		g.uiManager.AddMessage("Inventory closed")
		logger.Info("Inventory closed for player: %s", activePlayer.RPGStats().Name)
	} else {
		err := g.uiManager.ShowInventory(activePlayer)
		if err != nil {
			g.uiManager.AddMessage(fmt.Sprintf("Failed to show inventory: %v", err))
			logger.Error("Failed to show inventory for player %s: %v", activePlayer.RPGStats().Name, err)
			return
		}
		g.uiManager.AddMessage(fmt.Sprintf("Opened %s's inventory", activePlayer.RPGStats().Name))
		logger.Info("Inventory opened for player: %s", activePlayer.RPGStats().Name)
	}
}

// toggleSkills shows/hides the skills window for the active player
func (g *Game) toggleSkills() {
	activePlayer := g.GetActivePlayer()
	if activePlayer == nil || activePlayer.RPGStats() == nil {
		g.uiManager.AddMessage("No active player to show skills for")
		logger.Info("Attempted to toggle skills but no active player available")
		return
	}

	// Check if player has skills component, if not create one
	if activePlayer.Skills() == nil {
		// Create skills component for the player
		skillsComp := components.NewSkillsComponent(activePlayer.RPGStats().Job)
		skillsComp.AddSkillPoints(5) // Give some initial skill points
		activePlayer.AddComponent(ecs.ComponentSkills, skillsComp)
		logger.Info("Created skills component for player %s", activePlayer.RPGStats().Name)
	}

	// Toggle skills visibility
	if g.uiManager.IsSkillsVisible() {
		g.uiManager.HideSkills()
		g.uiManager.AddMessage("Skills window closed")
		logger.Info("Skills closed for player: %s", activePlayer.RPGStats().Name)
	} else {
		err := g.uiManager.ShowSkills(activePlayer)
		if err != nil {
			g.uiManager.AddMessage(fmt.Sprintf("Failed to show skills: %v", err))
			logger.Error("Failed to show skills for player %s: %v", activePlayer.RPGStats().Name, err)
			return
		}
		g.uiManager.AddMessage(fmt.Sprintf("Opened %s's skills", activePlayer.RPGStats().Name))
		logger.Info("Skills opened for player: %s", activePlayer.RPGStats().Name)
	}
}

// toggleQuestJournal shows/hides the quest journal for the active player
func (g *Game) toggleQuestJournal() {
	activePlayer := g.GetActivePlayer()
	if activePlayer == nil || activePlayer.RPGStats() == nil {
		g.uiManager.AddMessage("No active player to show quest journal for")
		logger.Info("Attempted to toggle quest journal but no active player available")
		return
	}

	// Check if player has quest journal component, if not create one
	if activePlayer.QuestJournal() == nil {
		// Create quest journal component for the player
		questJournal := components.NewQuestJournalComponent()
		activePlayer.AddComponent(ecs.ComponentQuestJournal, questJournal)
		logger.Info("Created quest journal component for player %s", activePlayer.RPGStats().Name)
	}

	// Toggle quest journal visibility
	if g.uiManager.IsQuestJournalVisible() {
		g.uiManager.HideQuestJournal()
		g.uiManager.AddMessage("Quest journal closed")
		logger.Info("Quest journal closed for player: %s", activePlayer.RPGStats().Name)
	} else {
		err := g.uiManager.ShowQuestJournal(activePlayer)
		if err != nil {
			g.uiManager.AddMessage(fmt.Sprintf("Failed to show quest journal: %v", err))
			logger.Error("Failed to show quest journal for player %s: %v", activePlayer.RPGStats().Name, err)
			return
		}
		g.uiManager.AddMessage(fmt.Sprintf("Opened %s's quest journal", activePlayer.RPGStats().Name))
		logger.Info("Quest journal opened for player: %s", activePlayer.RPGStats().Name)
	}
}
