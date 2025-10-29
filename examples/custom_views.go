// Package examples demonstrates how to create custom views using the ViewManager system
package examples

import (
	"log"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/engine"
	"github.com/jrecuero/myrpg/internal/logger"
)

// CustomViewType demonstrates how to extend the ViewType enum
const (
	ViewBossBattle engine.ViewType = 100 + iota // Boss battle view
	ViewPuzzle                                  // Puzzle mini-game view
	ViewTown                                    // Town exploration view
	ViewDungeon                                 // Dungeon exploration view
	ViewArena                                   // Arena combat view
)

// SetupCustomViews demonstrates how to register custom views with the ViewManager
func SetupCustomViews(game *engine.Game) {
	viewManager := game.GetViewManager()

	// Boss Battle View - Special tactical combat with unique mechanics
	viewManager.RegisterView(&engine.ViewConfiguration{
		Type:                ViewBossBattle,
		Name:                "Boss Battle",
		AllowsPlayerControl: true,
		ShowsUI:             true,
		PausesGame:          false,
		EntityFilter: func(entity *ecs.Entity) bool {
			// Show boss battle participants and special effects
			if entity.HasTag("background") || entity.HasTag("boss_arena") {
				return true
			}
			// Show all combat participants (includes party and boss)
			if entity.HasTag("player") || entity.HasTag("boss") || entity.HasTag("minion") {
				return true
			}
			// Show special battle effects
			if entity.HasTag("boss_effect") || entity.HasTag("arena_hazard") {
				return true
			}
			return false
		},
		EventFilter: func(eventComp *components.EventComponent) bool {
			// Only boss battle specific events
			return engine.IsBossEvent(eventComp)
		},
		UpdateHandler: func(deltaTime float64) error {
			// Custom boss battle mechanics
			logger.Debug("Updating boss battle view")
			// Add boss-specific update logic here
			return nil
		},
		InputHandler: func() error {
			// Custom input handling for boss battles
			return nil
		},
	})

	// Puzzle View - For puzzle mini-games
	viewManager.RegisterView(&engine.ViewConfiguration{
		Type:                ViewPuzzle,
		Name:                "Puzzle",
		AllowsPlayerControl: false, // Controlled by puzzle mechanics
		ShowsUI:             true,
		PausesGame:          true,
		EntityFilter: func(entity *ecs.Entity) bool {
			// Only show puzzle elements
			return entity.HasTag("puzzle_piece") || entity.HasTag("puzzle_background")
		},
		EventFilter: func(eventComp *components.EventComponent) bool {
			// Only puzzle-related events
			return engine.IsPuzzleEvent(eventComp)
		},
		UpdateHandler: func(deltaTime float64) error {
			// Puzzle logic updates
			return nil
		},
	})

	// Town Exploration View - Different from regular exploration
	viewManager.RegisterView(&engine.ViewConfiguration{
		Type:                ViewTown,
		Name:                "Town Exploration",
		AllowsPlayerControl: true,
		ShowsUI:             true,
		PausesGame:          false,
		EntityFilter: func(entity *ecs.Entity) bool {
			// Show town-specific entities
			if entity.HasTag("town_background") || entity.HasTag("building") {
				return true
			}
			// Show NPCs and town events
			if entity.HasTag("npc") || entity.HasTag("shop_sign") {
				return true
			}
			// Show party leader
			if entity.HasTag("player") {
				// Could implement party following logic here
				return entity == game.GetActivePlayer()
			}
			return false
		},
		EventFilter: func(eventComp *components.EventComponent) bool {
			// Only town events (shops, NPCs, etc.)
			return engine.IsTownEvent(eventComp)
		},
	})

	// Arena Combat View - Structured PvP or tournament battles
	viewManager.RegisterView(&engine.ViewConfiguration{
		Type:                ViewArena,
		Name:                "Arena Combat",
		AllowsPlayerControl: true,
		ShowsUI:             true,
		PausesGame:          false,
		EntityFilter: func(entity *ecs.Entity) bool {
			// Show arena environment and participants
			if entity.HasTag("arena_background") || entity.HasTag("arena_decoration") {
				return true
			}
			// Show combatants
			if entity.HasTag("arena_participant") || entity.HasTag("spectator") {
				return true
			}
			return false
		},
		EventFilter: func(eventComp *components.EventComponent) bool {
			// Arena-specific events
			return engine.IsArenaEvent(eventComp)
		},
		UpdateHandler: func(deltaTime float64) error {
			// Arena-specific mechanics (crowd reactions, special rules)
			return nil
		},
	})

	logger.Info("Registered custom views: Boss Battle, Puzzle, Town, Arena")
}

// SetupCustomTransitions demonstrates how to create custom view transitions
func SetupCustomTransitions(game *engine.Game) {
	viewManager := game.GetViewManager()

	// Battle to Boss Battle transition
	viewManager.RegisterTransition(&engine.ViewTransition{
		FromView: engine.ViewTactical,
		ToView:   ViewBossBattle,
		Condition: func() bool {
			// Check if current battle involves a boss
			data, exists := viewManager.GetTransitionData("battle_type")
			return exists && data == "boss_battle"
		},
		TransitionFn: func() error {
			logger.Info("Transitioning to boss battle mode")
			// Custom transition logic for boss battles
			return nil
		},
		Priority: 15, // Higher than normal battle transitions
	})

	// Town entry transition
	viewManager.RegisterTransition(&engine.ViewTransition{
		FromView: engine.ViewExploration,
		ToView:   ViewTown,
		Condition: func() bool {
			// Check if player entered a town area
			data, exists := viewManager.GetTransitionData("area_type")
			return exists && data == "town"
		},
		TransitionFn: func() error {
			logger.Info("Entering town")
			// Setup town-specific UI, load town data, etc.
			return nil
		},
		Priority: 8,
	})

	// Puzzle activation transition
	viewManager.RegisterTransition(&engine.ViewTransition{
		FromView: -1, // From any view
		ToView:   ViewPuzzle,
		Condition: func() bool {
			// Check if puzzle was activated
			data, exists := viewManager.GetTransitionData("puzzle_activated")
			return exists && data.(bool)
		},
		TransitionFn: func() error {
			logger.Info("Starting puzzle")
			// Initialize puzzle state
			return nil
		},
		Priority: 20, // High priority to interrupt other activities
	})

	logger.Info("Registered custom view transitions")
}

// DemonstrateBossBattleEvent shows how to trigger a boss battle
func DemonstrateBossBattleEvent(game *engine.Game, bossEntity *ecs.Entity) {
	viewManager := game.GetViewManager()

	// Set transition data to indicate this is a boss battle
	transitionData := map[string]interface{}{
		"battle_type": "boss_battle",
		"boss_entity": bossEntity,
		"arena_type":  "volcano_arena",
	}

	// Switch to boss battle view
	if err := viewManager.SwitchToView(ViewBossBattle, transitionData); err != nil {
		log.Printf("Failed to start boss battle: %v", err)
		// Fallback to regular tactical battle
		game.SwitchToTacticalMode([]*ecs.Entity{bossEntity})
	}
}

// DemonstrateTownEntry shows how to transition to town view
func DemonstrateTownEntry(game *engine.Game, townName string) {
	viewManager := game.GetViewManager()

	transitionData := map[string]interface{}{
		"area_type":   "town",
		"town_name":   townName,
		"entry_point": "main_gate",
	}

	if err := viewManager.SwitchToView(ViewTown, transitionData); err != nil {
		log.Printf("Failed to enter town: %v", err)
	}
}

// DemonstratePuzzleActivation shows how to start a puzzle
func DemonstratePuzzleActivation(game *engine.Game, puzzleType string) {
	viewManager := game.GetViewManager()

	transitionData := map[string]interface{}{
		"puzzle_activated": true,
		"puzzle_type":      puzzleType,
		"difficulty":       "medium",
	}

	if err := viewManager.SwitchToView(ViewPuzzle, transitionData); err != nil {
		log.Printf("Failed to start puzzle: %v", err)
	}
}

// ViewSystemUsageExample demonstrates typical usage patterns
func ViewSystemUsageExample(game *engine.Game) {
	// Setup custom views and transitions
	SetupCustomViews(game)
	SetupCustomTransitions(game)

	// Example: Check current view and react accordingly
	currentView := game.GetCurrentViewType()
	switch currentView {
	case engine.ViewExploration:
		logger.Info("Currently exploring")
	case engine.ViewTactical:
		logger.Info("In tactical combat")
	case ViewBossBattle:
		logger.Info("Fighting a boss!")
	case ViewTown:
		logger.Info("In town")
	case ViewPuzzle:
		logger.Info("Solving a puzzle")
	}

	// Example: Register a view-specific entity
	if game.IsInView(ViewTown) {
		// Create town-specific entities
		shopKeeper := ecs.NewEntity("ShopKeeper_001")
		shopKeeper.AddTag("npc")
		shopKeeper.AddTag("shop_keeper")
		// Add to view-specific entities
		game.GetViewManager().AddViewEntity(ViewTown, shopKeeper)
	}

	// Example: Conditional event handling based on view
	viewManager := game.GetViewManager()
	if viewManager.GetCurrentView() == ViewArena {
		// Arena-specific event handling
		logger.Info("Arena events are active")
	}
}
