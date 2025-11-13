// Package ui provides combat-specific UI components for the turn-based battle system
package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/logger"
	"github.com/jrecuero/myrpg/internal/tactical"
)

// CombatUIState represents the current state of combat UI interaction
type CombatUIState int

const (
	CombatUIStateNone CombatUIState = iota
	CombatUIStateSelectingAction
	CombatUIStateSelectingMoveTarget
	CombatUIStateSelectingAttackTarget
	CombatUIStateActionConfirmation
)

func (cuis CombatUIState) String() string {
	switch cuis {
	case CombatUIStateNone:
		return "None"
	case CombatUIStateSelectingAction:
		return "Selecting Action"
	case CombatUIStateSelectingMoveTarget:
		return "Selecting Move Target"
	case CombatUIStateSelectingAttackTarget:
		return "Selecting Attack Target"
	case CombatUIStateActionConfirmation:
		return "Action Confirmation"
	default:
		return "Unknown"
	}
}

// ActionButton represents a clickable action button
type ActionButton struct {
	X, Y, Width, Height float32
	Text                string
	Enabled             bool
	Visible             bool
	ActionType          tactical.ActionType
	HotKey              ebiten.Key
}

// CombatUI manages the combat-specific user interface
type CombatUI struct {
	State              CombatUIState
	ActionButtons      []*ActionButton
	SelectedAction     tactical.ActionType
	ValidMovePositions []tactical.GridPos
	ValidAttackTargets []*ecs.Entity
	HoveredPosition    *tactical.GridPos
	SelectedTarget     *ecs.Entity

	// Current context for action calculations
	CurrentCombatManager *tactical.TurnBasedCombatManager
	CurrentActiveUnit    *ecs.Entity

	// Stability flags to prevent constant recalculation
	TargetsCalculatedForUnit *ecs.Entity
	MovesCalculatedForUnit   *ecs.Entity

	// UI Layout
	ButtonAreaX, ButtonAreaY  float32
	ButtonWidth, ButtonHeight float32
	ButtonSpacing             float32

	// Colors
	ButtonColor         color.RGBA
	ButtonHoverColor    color.RGBA
	ButtonDisabledColor color.RGBA
	ButtonTextColor     color.RGBA

	// Callbacks
	OnActionSelected func(tactical.ActionType)
	OnMoveTarget     func(tactical.GridPos)
	OnAttackTarget   func(*ecs.Entity)
	OnCancel         func()
}

// NewCombatUI creates a new combat UI system
func NewCombatUI() *CombatUI {
	// Calculate button area (right side of screen)
	buttonAreaX := float32(constants.ScreenWidth - 200)
	buttonAreaY := float32(constants.GameWorldY + 10)
	buttonWidth := float32(180)
	buttonHeight := float32(30)
	buttonSpacing := float32(5)

	ui := &CombatUI{
		State:              CombatUIStateNone,
		ActionButtons:      make([]*ActionButton, 0),
		ValidMovePositions: make([]tactical.GridPos, 0),
		ValidAttackTargets: make([]*ecs.Entity, 0),

		ButtonAreaX:   buttonAreaX,
		ButtonAreaY:   buttonAreaY,
		ButtonWidth:   buttonWidth,
		ButtonHeight:  buttonHeight,
		ButtonSpacing: buttonSpacing,

		ButtonColor:         color.RGBA{60, 60, 100, 200},
		ButtonHoverColor:    color.RGBA{80, 80, 120, 220},
		ButtonDisabledColor: color.RGBA{40, 40, 40, 150},
		ButtonTextColor:     color.RGBA{255, 255, 255, 255},
	}

	// Create action buttons
	ui.createActionButtons()

	return ui
}

// createActionButtons creates the main action buttons
func (cui *CombatUI) createActionButtons() {
	actions := []struct {
		actionType tactical.ActionType
		text       string
		hotKey     ebiten.Key
	}{
		{tactical.ActionMove, "Move (M)", ebiten.KeyM},
		{tactical.ActionAttack, "Attack (A)", ebiten.KeyA},
		{tactical.ActionWait, "End Turn (E)", ebiten.KeyE},
	}

	for i, action := range actions {
		button := &ActionButton{
			X:          cui.ButtonAreaX,
			Y:          cui.ButtonAreaY + float32(i)*(cui.ButtonHeight+cui.ButtonSpacing),
			Width:      cui.ButtonWidth,
			Height:     cui.ButtonHeight,
			Text:       action.text,
			Enabled:    true,
			Visible:    true,
			ActionType: action.actionType,
			HotKey:     action.hotKey,
		}
		cui.ActionButtons = append(cui.ActionButtons, button)
	}
}

// Update handles combat UI input and state updates
func (cui *CombatUI) Update(combatManager *tactical.TurnBasedCombatManager, activeUnit *ecs.Entity) error {
	if combatManager == nil || activeUnit == nil {
		cui.State = CombatUIStateNone
		return nil
	}

	// Only show UI during player turns
	if !combatManager.IsPlayerTurn() {
		cui.State = CombatUIStateNone
		cui.CurrentCombatManager = nil
		cui.CurrentActiveUnit = nil
		return nil
	}

	// Store current context
	cui.CurrentCombatManager = combatManager

	// Initialize targets when active unit changes
	if cui.CurrentActiveUnit != activeUnit {
		cui.CurrentActiveUnit = activeUnit
		cui.initializeValidActions()
		cui.State = CombatUIStateSelectingAction // Reset to action selection
	}

	// Update valid actions for current unit
	cui.updateValidActions(combatManager, activeUnit)

	// Handle input based on current state
	switch cui.State {
	case CombatUIStateNone, CombatUIStateSelectingAction:
		return cui.updateActionSelection()
	case CombatUIStateSelectingMoveTarget:
		return cui.updateMoveTargetSelection(combatManager, activeUnit)
	case CombatUIStateSelectingAttackTarget:
		return cui.updateAttackTargetSelection(combatManager, activeUnit)
	}

	return nil
}

// updateValidActions updates the valid actions for the current unit
func (cui *CombatUI) updateValidActions(combatManager *tactical.TurnBasedCombatManager, activeUnit *ecs.Entity) {
	actionPoints := activeUnit.ActionPoints()
	if actionPoints == nil {
		// Disable all buttons
		for _, button := range cui.ActionButtons {
			button.Enabled = false
		}
		return
	}

	// Update button states based on available AP (without recalculating targets every frame)
	for _, button := range cui.ActionButtons {
		switch button.ActionType {
		case tactical.ActionMove:
			button.Enabled = len(cui.ValidMovePositions) > 0 && actionPoints.Current >= constants.MovementAPCost
		case tactical.ActionAttack:
			button.Enabled = len(cui.ValidAttackTargets) > 0 && actionPoints.Current >= constants.AttackAPCost
		case tactical.ActionWait:
			button.Enabled = true // End turn is always available
		}
	}

	// Set state to action selection if not already in a selection mode
	if cui.State == CombatUIStateNone {
		cui.State = CombatUIStateSelectingAction
	}
}

// initializeValidActions calculates valid actions when a unit's turn starts
func (cui *CombatUI) initializeValidActions() {
	// Clear previous action state
	cui.SelectedAction = 0 // Reset to no action selected
	cui.ValidMovePositions = cui.ValidMovePositions[:0]
	cui.ValidAttackTargets = cui.ValidAttackTargets[:0]
	cui.HoveredPosition = nil
	cui.SelectedTarget = nil

	// Reset calculation flags for new unit
	cui.TargetsCalculatedForUnit = nil
	cui.MovesCalculatedForUnit = nil

	// Pre-calculate valid moves for efficiency
	if cui.CurrentCombatManager != nil && cui.CurrentActiveUnit != nil {
		cui.ValidMovePositions = cui.CurrentCombatManager.GetValidMovesForUnit(cui.CurrentActiveUnit)
		cui.MovesCalculatedForUnit = cui.CurrentActiveUnit
		logger.Debug("Initialized valid actions for %s: %d moves available",
			cui.CurrentActiveUnit.GetID(), len(cui.ValidMovePositions))
	}
}

// updateActionSelection handles action button selection
func (cui *CombatUI) updateActionSelection() error {
	// Handle mouse input on buttons
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()

		for _, button := range cui.ActionButtons {
			if cui.isPointInButton(float32(mouseX), float32(mouseY), button) {
				if button.Enabled {
					cui.selectAction(button.ActionType)
					return nil
				} else if button.ActionType == tactical.ActionAttack {
					// Provide feedback when clicking disabled attack button
					logger.Warn("❌ Cannot attack: No enemies in range! Move closer to an enemy first.")
				}
			}
		}
	}

	// Handle hotkeys
	for _, button := range cui.ActionButtons {
		if button.Enabled && inpututil.IsKeyJustPressed(button.HotKey) {
			cui.selectAction(button.ActionType)
			return nil
		}
	}

	// Handle escape to cancel
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		cui.State = CombatUIStateNone
		if cui.OnCancel != nil {
			cui.OnCancel()
		}
	}

	return nil
}

// selectAction handles action selection
func (cui *CombatUI) selectAction(actionType tactical.ActionType) {
	cui.SelectedAction = actionType

	switch actionType {
	case tactical.ActionMove:
		// Calculate move positions only if not already calculated for this unit
		if cui.CurrentCombatManager != nil && cui.CurrentActiveUnit != nil && cui.MovesCalculatedForUnit != cui.CurrentActiveUnit {
			cui.ValidMovePositions = cui.CurrentCombatManager.GetValidMovesForUnit(cui.CurrentActiveUnit)
			cui.MovesCalculatedForUnit = cui.CurrentActiveUnit
			logger.Debug("✅ Calculated %d valid move positions for %s", len(cui.ValidMovePositions), cui.CurrentActiveUnit.GetID())
		}
		cui.State = CombatUIStateSelectingMoveTarget
	case tactical.ActionAttack:
		// Calculate valid attack targets only if not already calculated for this unit
		if cui.CurrentCombatManager != nil && cui.CurrentActiveUnit != nil && cui.TargetsCalculatedForUnit != cui.CurrentActiveUnit {
			cui.ValidAttackTargets = cui.CurrentCombatManager.GetValidAttackTargetsForUnit(cui.CurrentActiveUnit)
			cui.TargetsCalculatedForUnit = cui.CurrentActiveUnit
			logger.Debug("✅ Calculated %d valid attack targets for %s", len(cui.ValidAttackTargets), cui.CurrentActiveUnit.GetID())
		}

		// Only enter attack target selection if there are valid targets
		if len(cui.ValidAttackTargets) > 0 {
			cui.State = CombatUIStateSelectingAttackTarget
			message := fmt.Sprintf("🎯 Select an enemy to attack (found %d targets)", len(cui.ValidAttackTargets))
			logger.UI("%s", message)
		} else {
			// Provide feedback when no targets available
			message1 := "❌ No enemies in attack range! Move closer to an enemy first."
			message2 := "💡 Tip: Use the Move button to get within 1 tile of an enemy"
			logger.UI("%s", message1)
			logger.UI("%s", message2)
			logger.Warn("⚠️  No enemies in attack range when attack selected")
			// Stay in action selection mode
			cui.State = CombatUIStateSelectingAction
		}
	case tactical.ActionWait:
		// End turn immediately
		if cui.OnActionSelected != nil {
			cui.OnActionSelected(actionType)
		}
		cui.State = CombatUIStateNone
	}
}

// updateMoveTargetSelection handles movement target selection
func (cui *CombatUI) updateMoveTargetSelection(combatManager *tactical.TurnBasedCombatManager, activeUnit *ecs.Entity) error {
	_ = activeUnit
	// Handle mouse click on grid
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		screenX, screenY := float64(x), float64(y)
		offsetX, offsetY := float64(constants.GridOffsetX), float64(constants.GridOffsetY)

		// Check if click is within the grid bounds
		if screenX >= offsetX && screenY >= offsetY {
			gridPos := combatManager.Grid.WorldToGrid(screenX-offsetX, screenY-offsetY)

			// Validate grid bounds
			if gridPos.X >= 0 && gridPos.Y >= 0 && gridPos.X < constants.GridWidth && gridPos.Y < constants.GridHeight {
				// Check if clicked position is a valid move target
				for _, validPos := range cui.ValidMovePositions {
					if validPos.X == gridPos.X && validPos.Y == gridPos.Y {
						if cui.OnMoveTarget != nil {
							cui.OnMoveTarget(validPos)
						}
						cui.State = CombatUIStateNone
						return nil
					}
				}
			}
		}
	}

	// Handle escape to cancel
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		cui.State = CombatUIStateSelectingAction
	}

	return nil
}

// updateAttackTargetSelection handles attack target selection
func (cui *CombatUI) updateAttackTargetSelection(combatManager *tactical.TurnBasedCombatManager, activeUnit *ecs.Entity) error {
	_ = activeUnit
	// Handle mouse click on grid
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		screenX, screenY := float64(x), float64(y)
		offsetX, offsetY := float64(constants.GridOffsetX), float64(constants.GridOffsetY)

		// Check if click is within the grid bounds
		if screenX >= offsetX && screenY >= offsetY {
			gridPos := combatManager.Grid.WorldToGrid(screenX-offsetX, screenY-offsetY)

			// Validate grid bounds
			if gridPos.X >= 0 && gridPos.Y >= 0 && gridPos.X < constants.GridWidth && gridPos.Y < constants.GridHeight {
				// Check if there's a valid target at this position
				targetFound := false
				for _, target := range cui.ValidAttackTargets {
					if transform := target.Transform(); transform != nil {
						targetPos := combatManager.Grid.WorldToGrid(transform.X-offsetX, transform.Y-offsetY)
						if targetPos.X == gridPos.X && targetPos.Y == gridPos.Y {
							logger.Debug("🗡️  Attacking %s!\n", target.GetID())
							if cui.OnAttackTarget != nil {
								cui.OnAttackTarget(target)
							}
							cui.State = CombatUIStateNone
							targetFound = true
							return nil
						}
					}
				}

				// If no valid target was clicked, provide feedback
				if !targetFound && len(cui.ValidAttackTargets) > 0 {
					logger.Info("❌ Click on a highlighted enemy to attack (found %d valid targets)\n", len(cui.ValidAttackTargets))
				} else if len(cui.ValidAttackTargets) == 0 {
					logger.Warn("❌ No valid attack targets! Move closer to enemies first.\n")
					cui.State = CombatUIStateSelectingAction
				}
			}
		}
	}

	// Handle escape to cancel
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		cui.State = CombatUIStateSelectingAction
	}

	return nil
}

// Draw renders the combat UI
func (cui *CombatUI) Draw(screen *ebiten.Image, combatManager *tactical.TurnBasedCombatManager, activeUnit *ecs.Entity) {
	if combatManager == nil || !combatManager.IsPlayerTurn() {
		return
	}

	// Draw action buttons
	cui.drawActionButtons(screen)

	// Draw state-specific UI
	switch cui.State {
	case CombatUIStateSelectingMoveTarget:
		cui.drawMoveTargetSelection(screen)
	case CombatUIStateSelectingAttackTarget:
		cui.drawAttackTargetSelection(screen)
	}

	// Draw turn information
	cui.drawTurnInfo(screen, combatManager, activeUnit)
}

// drawActionButtons renders the action buttons
func (cui *CombatUI) drawActionButtons(screen *ebiten.Image) {
	mouseX, mouseY := ebiten.CursorPosition()

	for _, button := range cui.ActionButtons {
		if !button.Visible {
			continue
		}

		// Determine button color
		buttonColor := cui.ButtonColor
		if !button.Enabled {
			buttonColor = cui.ButtonDisabledColor
		} else if cui.isPointInButton(float32(mouseX), float32(mouseY), button) {
			buttonColor = cui.ButtonHoverColor
		}

		// Draw button background
		vector.FillRect(screen, button.X, button.Y, button.Width, button.Height, buttonColor, false)

		// Draw button border
		borderColor := color.RGBA{255, 255, 255, 100}
		if !button.Enabled {
			borderColor = color.RGBA{100, 100, 100, 100}
		}
		vector.StrokeRect(screen, button.X, button.Y, button.Width, button.Height, 1, borderColor, false)

		// Draw button text
		// Note: ebitenutil.DebugPrintAt uses a default color, we'll enhance this later
		// Center text in button
		textX := int(button.X + 10)
		textY := int(button.Y + button.Height/2 - 5)
		ebitenutil.DebugPrintAt(screen, button.Text, textX, textY)

		// Show tooltip for disabled attack button
		if !button.Enabled && button.ActionType == tactical.ActionAttack &&
			cui.isPointInButton(float32(mouseX), float32(mouseY), button) {
			tooltipText := "No enemies in range (must be adjacent)"
			tooltipX := int(button.X + button.Width + 5)
			tooltipY := int(button.Y)

			// Draw tooltip background
			tooltipWidth := float32(len(tooltipText)*6 + 10)
			tooltipHeight := float32(20)
			vector.FillRect(screen, float32(tooltipX), float32(tooltipY), tooltipWidth, tooltipHeight,
				color.RGBA{0, 0, 0, 200}, false)
			vector.StrokeRect(screen, float32(tooltipX), float32(tooltipY), tooltipWidth, tooltipHeight,
				1, color.RGBA{255, 255, 0, 255}, false)

			// Draw tooltip text
			ebitenutil.DebugPrintAt(screen, tooltipText, tooltipX+5, tooltipY+5)
		}
	}
}

// drawMoveTargetSelection renders move target highlighting
func (cui *CombatUI) drawMoveTargetSelection(screen *ebiten.Image) {
	// Highlight valid movement positions
	moveColor := color.RGBA{0, 150, 255, 100} // Blue

	for _, pos := range cui.ValidMovePositions {
		screenX := float32(pos.X*constants.TileSize) + float32(constants.GridOffsetX)
		screenY := float32(pos.Y*constants.TileSize) + float32(constants.GridOffsetY)

		vector.FillRect(screen, screenX, screenY, constants.TileSize, constants.TileSize, moveColor, false)
	}

	// Draw instruction text
	instructionY := float32(constants.GameWorldY + constants.GameWorldHeight - 30)
	ebitenutil.DebugPrintAt(screen, "Click on a blue tile to move, ESC to cancel", 10, int(instructionY))
}

// drawAttackTargetSelection renders attack target highlighting
func (cui *CombatUI) drawAttackTargetSelection(screen *ebiten.Image) {
	// Highlight valid attack targets
	attackColor := color.RGBA{255, 100, 100, 150} // Red

	offsetX, offsetY := float64(constants.GridOffsetX), float64(constants.GridOffsetY)

	for i, target := range cui.ValidAttackTargets {
		if transform := target.Transform(); transform != nil {
			// Convert world coordinates to grid coordinates properly
			gridX := int((transform.X - offsetX) / constants.TileSize)
			gridY := int((transform.Y - offsetY) / constants.TileSize)

			// Convert back to screen coordinates for rendering
			screenX := float32(gridX*constants.TileSize) + float32(offsetX)
			screenY := float32(gridY*constants.TileSize) + float32(offsetY)

			vector.FillRect(screen, screenX, screenY, constants.TileSize, constants.TileSize, attackColor, false)

			// Debug: show target info
			if i < 3 { // Only show first 3 to avoid clutter
				targetInfo := fmt.Sprintf("T%d", i+1)
				ebitenutil.DebugPrintAt(screen, targetInfo, int(screenX+2), int(screenY+2))
			}
		}
	}

	// Draw instruction text with target count
	instructionY := float32(constants.GameWorldY + constants.GameWorldHeight - 30)
	instruction := fmt.Sprintf("Click on a red tile to attack (%d targets), ESC to cancel", len(cui.ValidAttackTargets))
	ebitenutil.DebugPrintAt(screen, instruction, 10, int(instructionY))
}

// drawTurnInfo renders turn and AP information
func (cui *CombatUI) drawTurnInfo(screen *ebiten.Image, combatManager *tactical.TurnBasedCombatManager, activeUnit *ecs.Entity) {
	if activeUnit == nil {
		return
	}

	// Get unit information
	stats := activeUnit.RPGStats()
	actionPoints := activeUnit.ActionPoints()
	if stats == nil || actionPoints == nil {
		return
	}

	// Draw turn info panel (top right area)
	panelX := float32(constants.ScreenWidth - 200)
	panelY := float32(10)
	panelWidth := float32(190)
	panelHeight := float32(110) // Increased height to fit more info

	// Panel background
	panelColor := color.RGBA{30, 30, 30, 200}
	vector.FillRect(screen, panelX, panelY, panelWidth, panelHeight, panelColor, false)

	// Panel border
	borderColor := color.RGBA{100, 100, 100, 255}
	vector.StrokeRect(screen, panelX, panelY, panelWidth, panelHeight, 1, borderColor, false)

	// Round info
	roundInfo := fmt.Sprintf("Round: %d", combatManager.GetCurrentRound())
	ebitenutil.DebugPrintAt(screen, roundInfo, int(panelX+5), int(panelY+10))

	// Active unit info
	unitInfo := fmt.Sprintf("Active: %s", stats.Name)
	ebitenutil.DebugPrintAt(screen, unitInfo, int(panelX+5), int(panelY+25))

	jobInfo := fmt.Sprintf("Class: %s", stats.Job.String())
	ebitenutil.DebugPrintAt(screen, jobInfo, int(panelX+5), int(panelY+40))

	// AP info
	apInfo := fmt.Sprintf("AP: %d/%d", actionPoints.Current, actionPoints.Maximum)
	ebitenutil.DebugPrintAt(screen, apInfo, int(panelX+5), int(panelY+55))

	// Movement info (from legacy system for verification)
	moveInfo := fmt.Sprintf("Moves: %d/%d", stats.MovesRemaining, stats.MoveRange)
	ebitenutil.DebugPrintAt(screen, moveInfo, int(panelX+5), int(panelY+70))

	// Current state
	stateInfo := fmt.Sprintf("State: %s", cui.State.String())
	ebitenutil.DebugPrintAt(screen, stateInfo, int(panelX+5), int(panelY+85))
}

// Helper methods

// isPointInButton checks if a point is inside a button
func (cui *CombatUI) isPointInButton(x, y float32, button *ActionButton) bool {
	return x >= button.X && x <= button.X+button.Width &&
		y >= button.Y && y <= button.Y+button.Height
}

// SetCallbacks sets the UI callback functions
func (cui *CombatUI) SetCallbacks(
	onActionSelected func(tactical.ActionType),
	onMoveTarget func(tactical.GridPos),
	onAttackTarget func(*ecs.Entity),
	onCancel func(),
) {
	cui.OnActionSelected = onActionSelected
	cui.OnMoveTarget = onMoveTarget
	cui.OnAttackTarget = onAttackTarget
	cui.OnCancel = onCancel
}

// Reset resets the UI state
func (cui *CombatUI) Reset() {
	cui.State = CombatUIStateNone
	cui.SelectedAction = tactical.ActionMove // Default
	cui.ValidMovePositions = cui.ValidMovePositions[:0]
	cui.ValidAttackTargets = cui.ValidAttackTargets[:0]
	cui.HoveredPosition = nil
	cui.SelectedTarget = nil
}
