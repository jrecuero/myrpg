// Package engine provides a comprehensive view management system
// that handles different game views, entity visibility, and view transitions
package engine

import (
	"fmt"
	"log"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/logger"
)

// ViewType represents different types of game views
type ViewType int

const (
	ViewExploration ViewType = iota // Free movement exploration
	ViewTactical                    // Grid-based tactical combat
	ViewDialog                      // Dialog/conversation view
	ViewInventory                   // Inventory management view
	ViewShop                        // Shop interface view
	ViewMenu                        // Game menu view
	ViewCutscene                    // Cutscene/story view
	ViewWorldMap                    // World map navigation
)

func (vt ViewType) String() string {
	switch vt {
	case ViewExploration:
		return "Exploration"
	case ViewTactical:
		return "Tactical"
	case ViewDialog:
		return "Dialog"
	case ViewInventory:
		return "Inventory"
	case ViewShop:
		return "Shop"
	case ViewMenu:
		return "Menu"
	case ViewCutscene:
		return "Cutscene"
	case ViewWorldMap:
		return "WorldMap"
	default:
		return "Unknown"
	}
}

// ViewConfiguration defines how a view should behave
type ViewConfiguration struct {
	Type                ViewType                              // Type of this view
	Name                string                                // Human-readable name
	AllowsPlayerControl bool                                  // Whether player can control characters
	ShowsUI             bool                                  // Whether to show UI elements
	PausesGame          bool                                  // Whether this view pauses underlying game
	EntityFilter        func(*ecs.Entity) bool                // Function to determine which entities to show
	EventFilter         func(*components.EventComponent) bool // Function to determine which events are active
	UpdateHandler       func(float64) error                   // Custom update logic for this view
	InputHandler        func() error                          // Custom input handling for this view
}

// ViewTransition defines how to transition between views
type ViewTransition struct {
	FromView     ViewType     // Source view (use -1 for any)
	ToView       ViewType     // Target view
	Condition    func() bool  // Condition to check for transition
	TransitionFn func() error // Function to execute during transition
	Priority     int          // Priority (higher = checked first)
	AllowReverse bool         // Whether reverse transition is allowed
}

// ViewManager manages all game views and transitions
type ViewManager struct {
	game           *Game                           // Reference to main game
	currentView    ViewType                        // Currently active view
	previousView   ViewType                        // Previous view for back navigation
	viewStack      []ViewType                      // Stack for nested views
	views          map[ViewType]*ViewConfiguration // View configurations
	transitions    []*ViewTransition               // Available transitions
	viewEntities   map[ViewType][]*ecs.Entity      // Entities specific to each view
	globalEntities []*ecs.Entity                   // Entities visible in all views
	transitionData map[string]interface{}          // Data passed between views
}

// NewViewManager creates a new view manager
func NewViewManager(game *Game) *ViewManager {
	vm := &ViewManager{
		game:           game,
		currentView:    ViewExploration,
		previousView:   ViewExploration,
		viewStack:      make([]ViewType, 0),
		views:          make(map[ViewType]*ViewConfiguration),
		transitions:    make([]*ViewTransition, 0),
		viewEntities:   make(map[ViewType][]*ecs.Entity),
		globalEntities: make([]*ecs.Entity, 0),
		transitionData: make(map[string]interface{}),
	}

	vm.initializeDefaultViews()
	vm.initializeDefaultTransitions()

	return vm
}

// initializeDefaultViews sets up the standard game views
func (vm *ViewManager) initializeDefaultViews() {
	// Exploration View
	vm.RegisterView(&ViewConfiguration{
		Type:                ViewExploration,
		Name:                "Exploration",
		AllowsPlayerControl: true,
		ShowsUI:             true,
		PausesGame:          false,
		EntityFilter: func(entity *ecs.Entity) bool {
			// Show backgrounds, party leader, enemies, and event entities
			if entity.HasTag("background") {
				return true
			}
			if entity.HasTag("player") {
				return entity == vm.game.partyManager.GetPartyLeader()
			}
			if entity.HasTag("enemy") || entity.Event() != nil {
				return true
			}
			return false
		},
		EventFilter: func(eventComp *components.EventComponent) bool {
			return eventComp.IsActiveInMode(components.GameModeExploration)
		},
	})

	// Tactical View
	vm.RegisterView(&ViewConfiguration{
		Type:                ViewTactical,
		Name:                "Tactical Combat",
		AllowsPlayerControl: true,
		ShowsUI:             true,
		PausesGame:          false,
		EntityFilter: func(entity *ecs.Entity) bool {
			// Show background entities
			if entity.HasTag("background") {
				return true
			}

			// Show all players in tactical mode (not just leader)
			if entity.HasTag("player") {
				return true
			}

			// Show all enemies
			if entity.HasTag("enemy") {
				return true
			}

			// Show tactical participants (covers any other combat entities)
			if vm.game.tacticalManager != nil && vm.game.tacticalManager.Participants != nil {
				for _, participant := range vm.game.tacticalManager.Participants {
					if entity == participant {
						return true
					}
				}
			}

			// Show tactical-specific entities
			if entity.HasTag("tactical") || entity.HasTag("combat") {
				return true
			}

			return false
		},
		EventFilter: func(eventComp *components.EventComponent) bool {
			return eventComp.IsActiveInMode(components.GameModeTactical)
		},
	})

	// Dialog View
	vm.RegisterView(&ViewConfiguration{
		Type:                ViewDialog,
		Name:                "Dialog",
		AllowsPlayerControl: false,
		ShowsUI:             true,
		PausesGame:          true,
		EntityFilter: func(entity *ecs.Entity) bool {
			// Show same entities as exploration but disable movement
			return vm.views[ViewExploration].EntityFilter(entity)
		},
		EventFilter: func(eventComp *components.EventComponent) bool {
			// No events active during dialog
			return false
		},
	})

	// Inventory View
	vm.RegisterView(&ViewConfiguration{
		Type:                ViewInventory,
		Name:                "Inventory",
		AllowsPlayerControl: false,
		ShowsUI:             true,
		PausesGame:          true,
		EntityFilter: func(entity *ecs.Entity) bool {
			// Show same entities as exploration
			return vm.views[ViewExploration].EntityFilter(entity)
		},
		EventFilter: func(eventComp *components.EventComponent) bool {
			return false // No events during inventory
		},
	})
}

// initializeDefaultTransitions sets up common view transitions
func (vm *ViewManager) initializeDefaultTransitions() {
	// Battle event transitions
	vm.RegisterTransition(&ViewTransition{
		FromView: ViewExploration,
		ToView:   ViewTactical,
		Condition: func() bool {
			// This will be triggered by battle events
			return false // Controlled by events
		},
		TransitionFn: func() error {
			logger.Debug("Transitioning from exploration to tactical combat")
			return nil
		},
		Priority: 10,
	})

	// Tactical to exploration (battle end)
	vm.RegisterTransition(&ViewTransition{
		FromView: ViewTactical,
		ToView:   ViewExploration,
		Condition: func() bool {
			// Check if tactical combat has ended
			return !vm.game.tacticalManager.IsActive
		},
		TransitionFn: func() error {
			logger.Debug("Transitioning from tactical combat back to exploration")
			vm.game.tacticalManager.EndTacticalCombat()
			return nil
		},
		Priority: 10,
	})

	// ESC key transitions (from any modal view back to previous)
	for _, viewType := range []ViewType{ViewDialog, ViewInventory, ViewShop, ViewMenu} {
		vm.RegisterTransition(&ViewTransition{
			FromView: viewType,
			ToView:   -1, // Will use previous view
			Condition: func() bool {
				// This will be checked by input handling
				return false
			},
			TransitionFn: func() error {
				return vm.PopView()
			},
			Priority: 5,
		})
	}
}

// RegisterView adds a new view configuration
func (vm *ViewManager) RegisterView(config *ViewConfiguration) {
	vm.views[config.Type] = config
	logger.Debug("Registered view: %s", config.Name)
}

// RegisterTransition adds a new view transition
func (vm *ViewManager) RegisterTransition(transition *ViewTransition) {
	vm.transitions = append(vm.transitions, transition)
	logger.Debug("Registered transition: %s -> %s",
		ViewType(transition.FromView).String(),
		transition.ToView.String())
}

// GetCurrentView returns the currently active view type
func (vm *ViewManager) GetCurrentView() ViewType {
	return vm.currentView
}

// GetCurrentViewConfig returns the current view's configuration
func (vm *ViewManager) GetCurrentViewConfig() *ViewConfiguration {
	return vm.views[vm.currentView]
}

// SwitchToView transitions to a specific view
func (vm *ViewManager) SwitchToView(viewType ViewType, data map[string]interface{}) error {
	if vm.currentView == viewType {
		return nil // Already in target view
	}

	config, exists := vm.views[viewType]
	if !exists {
		return fmt.Errorf("view type %s not registered", viewType.String())
	}

	logger.Debug("Switching from %s to %s", vm.currentView.String(), viewType.String())

	// Store transition data
	for key, value := range data {
		vm.transitionData[key] = value
	}

	// Execute view-specific transition logic
	if config.UpdateHandler != nil {
		if err := config.UpdateHandler(0); err != nil {
			return fmt.Errorf("view transition failed: %v", err)
		}
	}

	vm.previousView = vm.currentView
	vm.currentView = viewType

	// Update game systems based on new view
	vm.updateGameSystems()

	logger.Info("Switched to %s view", config.Name)
	return nil
}

// PushView pushes current view to stack and switches to new view
func (vm *ViewManager) PushView(viewType ViewType, data map[string]interface{}) error {
	vm.viewStack = append(vm.viewStack, vm.currentView)
	return vm.SwitchToView(viewType, data)
}

// PopView returns to the previous view from the stack
func (vm *ViewManager) PopView() error {
	if len(vm.viewStack) == 0 {
		return fmt.Errorf("no views on stack to pop")
	}

	// Get previous view from stack
	previousView := vm.viewStack[len(vm.viewStack)-1]
	vm.viewStack = vm.viewStack[:len(vm.viewStack)-1]

	return vm.SwitchToView(previousView, nil)
}

// IsEntityVisible determines if an entity should be visible in the current view
func (vm *ViewManager) IsEntityVisible(entity *ecs.Entity) bool {
	config := vm.GetCurrentViewConfig()
	if config == nil || config.EntityFilter == nil {
		return true // Default to visible if no filter
	}
	return config.EntityFilter(entity)
}

// IsEventActive determines if an event should be active in the current view
func (vm *ViewManager) IsEventActive(eventComp *components.EventComponent) bool {
	config := vm.GetCurrentViewConfig()
	if config == nil || config.EventFilter == nil {
		return true // Default to active if no filter
	}
	return config.EventFilter(eventComp)
}

// Update processes view transitions and updates the current view
func (vm *ViewManager) Update(deltaTime float64) error {
	// Handle input for current view
	config := vm.GetCurrentViewConfig()
	if config != nil && config.InputHandler != nil {
		if err := config.InputHandler(); err != nil {
			return fmt.Errorf("view input handling failed: %v", err)
		}
	}

	// Check for automatic transitions
	vm.checkTransitions()

	// Update current view
	if config != nil && config.UpdateHandler != nil {
		if err := config.UpdateHandler(deltaTime); err != nil {
			return fmt.Errorf("view update failed: %v", err)
		}
	}

	return nil
}

// checkTransitions evaluates all registered transitions
func (vm *ViewManager) checkTransitions() {
	for _, transition := range vm.transitions {
		// Check if transition applies to current view
		if transition.FromView != -1 && ViewType(transition.FromView) != vm.currentView {
			continue
		}

		// Check transition condition
		if transition.Condition != nil && transition.Condition() {
			var targetView ViewType
			if transition.ToView == -1 {
				targetView = vm.previousView
			} else {
				targetView = transition.ToView
			}

			// Execute transition
			if transition.TransitionFn != nil {
				if err := transition.TransitionFn(); err != nil {
					log.Printf("Transition failed: %v", err)
					continue
				}
			}

			vm.SwitchToView(targetView, nil)
			break // Only execute one transition per frame
		}
	}
}

// updateGameSystems updates game systems based on current view
func (vm *ViewManager) updateGameSystems() {
	config := vm.GetCurrentViewConfig()
	if config == nil {
		return
	}

	// Update event manager game mode
	var gameMode components.GameMode
	switch vm.currentView {
	case ViewExploration:
		gameMode = components.GameModeExploration
	case ViewTactical:
		gameMode = components.GameModeTactical
	default:
		gameMode = components.GameModeExploration // Default to exploration
	}
	vm.game.eventManager.SetGameMode(gameMode)

	// Update other systems as needed
	// This is where you can add more system-specific updates
}

// GetTransitionData retrieves data passed during view transitions
func (vm *ViewManager) GetTransitionData(key string) (interface{}, bool) {
	value, exists := vm.transitionData[key]
	return value, exists
}

// ClearTransitionData clears all transition data
func (vm *ViewManager) ClearTransitionData() {
	vm.transitionData = make(map[string]interface{})
}

// AddViewEntity adds an entity specific to a view
func (vm *ViewManager) AddViewEntity(viewType ViewType, entity *ecs.Entity) {
	if vm.viewEntities[viewType] == nil {
		vm.viewEntities[viewType] = make([]*ecs.Entity, 0)
	}
	vm.viewEntities[viewType] = append(vm.viewEntities[viewType], entity)
}

// AddGlobalEntity adds an entity visible in all views
func (vm *ViewManager) AddGlobalEntity(entity *ecs.Entity) {
	vm.globalEntities = append(vm.globalEntities, entity)
}

// GetVisibleEntities returns all entities that should be visible in the current view
func (vm *ViewManager) GetVisibleEntities(allEntities []*ecs.Entity) []*ecs.Entity {
	visible := make([]*ecs.Entity, 0)

	// Add global entities
	visible = append(visible, vm.globalEntities...)

	// Add view-specific entities
	if viewEntities, exists := vm.viewEntities[vm.currentView]; exists {
		visible = append(visible, viewEntities...)
	}

	// Filter world entities based on current view
	for _, entity := range allEntities {
		if vm.IsEntityVisible(entity) {
			visible = append(visible, entity)
		}
	}

	return visible
}

// Event Helper Functions for View Filtering

// IsEventOfType checks if an event component is of a specific type
func IsEventOfType(eventComp *components.EventComponent, eventType components.EventType) bool {
	return eventComp.EventType == eventType
}

// IsEventNamed checks if an event component has a specific name
func IsEventNamed(eventComp *components.EventComponent, name string) bool {
	return eventComp.Name == name
}

// IsEventID checks if an event component has a specific ID
func IsEventID(eventComp *components.EventComponent, id string) bool {
	return eventComp.ID == id
}

// IsBattleEvent checks if this is a battle event
func IsBattleEvent(eventComp *components.EventComponent) bool {
	return eventComp.EventType == components.EventBattle
}

// IsDialogEvent checks if this is a dialog event
func IsDialogEvent(eventComp *components.EventComponent) bool {
	return eventComp.EventType == components.EventDialog
}

// IsShopEvent checks if this is a shop event
func IsShopEvent(eventComp *components.EventComponent) bool {
	return eventComp.EventType == components.EventShop
}

// IsBossEvent checks if this is a boss battle (by name or ID convention)
func IsBossEvent(eventComp *components.EventComponent) bool {
	return IsBattleEvent(eventComp) &&
		(IsEventNamed(eventComp, "Boss Battle") ||
			IsEventID(eventComp, "boss_encounter") ||
			IsEventID(eventComp, "boss_battle"))
}

// IsArenaEvent checks if this is an arena event (by name or ID convention)
func IsArenaEvent(eventComp *components.EventComponent) bool {
	return IsBattleEvent(eventComp) &&
		(IsEventNamed(eventComp, "Arena Battle") ||
			IsEventID(eventComp, "arena_event") ||
			IsEventID(eventComp, "arena_battle"))
}

// IsTownEvent checks if this is a town-related event
func IsTownEvent(eventComp *components.EventComponent) bool {
	return IsShopEvent(eventComp) || IsDialogEvent(eventComp) ||
		IsEventID(eventComp, "town_event")
}

// IsPuzzleEvent checks if this is a puzzle event
func IsPuzzleEvent(eventComp *components.EventComponent) bool {
	return IsEventNamed(eventComp, "Puzzle") ||
		IsEventID(eventComp, "puzzle_event")
}
