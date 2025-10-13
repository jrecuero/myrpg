// Package tactical provides turn-based combat management system
package tactical

import (
	"fmt"
	"sort"

	"github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// CombatPhase represents the current phase of combat
type CombatPhase int

const (
	CombatPhaseInitialization CombatPhase = iota
	CombatPhaseTeamTurn
	CombatPhaseActionExecution
	CombatPhaseEndTurn
	CombatPhaseVictoryCheck
	CombatPhaseEnded
)

func (cp CombatPhase) String() string {
	switch cp {
	case CombatPhaseInitialization:
		return "Initialization"
	case CombatPhaseTeamTurn:
		return "Team Turn"
	case CombatPhaseActionExecution:
		return "Action Execution"
	case CombatPhaseEndTurn:
		return "End Turn"
	case CombatPhaseVictoryCheck:
		return "Victory Check"
	case CombatPhaseEnded:
		return "Combat Ended"
	default:
		return "Unknown"
	}
}

// CombatResult represents the outcome of combat
type CombatResult int

const (
	CombatResultNone CombatResult = iota
	CombatResultPlayerVictory
	CombatResultEnemyVictory
	CombatResultOngoing
)

func (cr CombatResult) String() string {
	switch cr {
	case CombatResultNone:
		return "None"
	case CombatResultPlayerVictory:
		return "Player Victory"
	case CombatResultEnemyVictory:
		return "Enemy Victory"
	case CombatResultOngoing:
		return "Ongoing"
	default:
		return "Unknown"
	}
}

// TeamInfo holds information about a combat team
type TeamInfo struct {
	Team         components.Team
	Members      []*ecs.Entity
	TotalSpeed   int
	IsActive     bool
	HasCompleted bool // All team members have acted this round
}

// TurnBasedCombatManager orchestrates turn-based tactical combat
type TurnBasedCombatManager struct {
	// Combat State
	Phase        CombatPhase
	Result       CombatResult
	IsActive     bool
	CurrentRound int

	// Teams and Initiative
	Teams           []*TeamInfo
	ActiveTeam      *TeamInfo
	InitiativeOrder []*TeamInfo

	// Current Action
	ActiveUnit    *ecs.Entity
	PendingAction *CombatAction

	// Grid and Systems
	Grid *Grid

	// Callbacks for UI communication
	MessageCallback     func(string)
	StateChangeCallback func(CombatPhase)

	// Debug and Logging
	DebugMode bool
}

// CombatAction represents a combat action to be executed
type CombatAction struct {
	Type      ActionType
	Actor     *ecs.Entity
	Target    *ecs.Entity
	TargetPos GridPos
	APCost    int
	Validated bool
	Message   string
}

// NewTurnBasedCombatManager creates a new combat manager
func NewTurnBasedCombatManager(grid *Grid) *TurnBasedCombatManager {
	return &TurnBasedCombatManager{
		Phase:        CombatPhaseInitialization,
		Result:       CombatResultNone,
		IsActive:     false,
		CurrentRound: 0,
		Teams:        make([]*TeamInfo, 0),
		Grid:         grid,
		DebugMode:    true, // Enable debug logging initially
	}
}

// SetMessageCallback sets the callback for sending messages to UI
func (cbm *TurnBasedCombatManager) SetMessageCallback(callback func(string)) {
	cbm.MessageCallback = callback
}

// SetStateChangeCallback sets the callback for phase changes
func (cbm *TurnBasedCombatManager) SetStateChangeCallback(callback func(CombatPhase)) {
	cbm.StateChangeCallback = callback
}

// InitializeCombat sets up combat with the given entities
func (cbm *TurnBasedCombatManager) InitializeCombat(entities []*ecs.Entity) error {
	cbm.logMessage("Initializing turn-based combat...")

	// Reset combat state
	cbm.Phase = CombatPhaseInitialization
	cbm.Result = CombatResultOngoing
	cbm.IsActive = true
	cbm.CurrentRound = 1
	cbm.Teams = make([]*TeamInfo, 0)
	cbm.ActiveTeam = nil
	cbm.ActiveUnit = nil
	cbm.PendingAction = nil

	// Add combat components to all entities
	for _, entity := range entities {
		if err := cbm.initializeEntityForCombat(entity); err != nil {
			return fmt.Errorf("failed to initialize entity %s for combat: %v", entity.GetID(), err)
		}
	}

	// Create teams based on entity tags
	if err := cbm.createTeams(entities); err != nil {
		return fmt.Errorf("failed to create teams: %v", err)
	}

	// Calculate initiative and set turn order
	cbm.calculateInitiative()

	// Start first team's turn
	cbm.changePhase(CombatPhaseTeamTurn)
	cbm.startNextTeamTurn()

	cbm.logMessage(fmt.Sprintf("Combat initialized with %d teams, %d total units",
		len(cbm.Teams), len(entities)))

	return nil
}

// initializeEntityForCombat adds required combat components to an entity
func (cbm *TurnBasedCombatManager) initializeEntityForCombat(entity *ecs.Entity) error {
	stats := entity.RPGStats()
	if stats == nil {
		return fmt.Errorf("entity has no RPG stats")
	}

	// Determine team based on entity tags
	var team components.Team
	if entity.HasTag(ecs.TagPlayer) {
		team = components.TeamPlayer
	} else if entity.HasTag(ecs.TagEnemy) {
		team = components.TeamEnemy
	} else {
		return fmt.Errorf("entity has no valid team tag")
	}

	// Add ActionPoints component based on job class
	maxAP := cbm.getMaxAPForJob(stats.Job)
	actionPoints := components.NewActionPointsComponent(maxAP)
	entity.AddComponent(ecs.ComponentActionPoints, actionPoints)

	// Add CombatState component
	initiative := cbm.calculateEntityInitiative(stats)
	combatState := components.NewCombatStateComponent(team, initiative)
	entity.AddComponent(ecs.ComponentCombatState, combatState)

	cbm.logMessage(fmt.Sprintf("Initialized %s (%s) - Team: %s, AP: %d, Initiative: %d",
		stats.Name, stats.Job.String(), team.String(), maxAP, initiative))

	return nil
}

// getMaxAPForJob returns the maximum AP for a given job class
func (cbm *TurnBasedCombatManager) getMaxAPForJob(job components.JobType) int {
	switch job {
	case components.JobWarrior:
		return constants.WarriorMaxAP
	case components.JobMage:
		return constants.MageMaxAP
	case components.JobRogue:
		return constants.RogueMaxAP
	case components.JobCleric:
		return constants.ClericMaxAP
	case components.JobArcher:
		return constants.ArcherMaxAP
	default:
		return constants.WarriorMaxAP // Default fallback
	}
}

// calculateEntityInitiative calculates initiative for a single entity
func (cbm *TurnBasedCombatManager) calculateEntityInitiative(stats *components.RPGStatsComponent) int {
	// For now, just use speed stat
	// Future: could add randomization or other factors
	return stats.Speed
}

// createTeams organizes entities into teams
func (cbm *TurnBasedCombatManager) createTeams(entities []*ecs.Entity) error {
	teamMap := make(map[components.Team][]*ecs.Entity)

	// Group entities by team
	for _, entity := range entities {
		combatState := entity.CombatState()
		if combatState == nil {
			continue // Skip entities without combat state
		}

		teamMap[combatState.Team] = append(teamMap[combatState.Team], entity)
	}

	// Create TeamInfo objects
	for team, members := range teamMap {
		if len(members) == 0 {
			continue
		}

		teamInfo := &TeamInfo{
			Team:         team,
			Members:      members,
			TotalSpeed:   cbm.calculateTeamSpeed(members),
			IsActive:     false,
			HasCompleted: false,
		}

		cbm.Teams = append(cbm.Teams, teamInfo)

		cbm.logMessage(fmt.Sprintf("Created %s team with %d members, total speed: %d",
			team.String(), len(members), teamInfo.TotalSpeed))
	}

	if len(cbm.Teams) < 2 {
		return fmt.Errorf("need at least 2 teams for combat")
	}

	return nil
}

// calculateTeamSpeed calculates the total speed for a team
func (cbm *TurnBasedCombatManager) calculateTeamSpeed(members []*ecs.Entity) int {
	totalSpeed := 0
	for _, member := range members {
		if stats := member.RPGStats(); stats != nil {
			totalSpeed += stats.Speed
		}
	}
	return totalSpeed
}

// calculateInitiative determines turn order based on team speeds
func (cbm *TurnBasedCombatManager) calculateInitiative() {
	// Sort teams by total speed (highest first)
	sort.Slice(cbm.Teams, func(i, j int) bool {
		return cbm.Teams[i].TotalSpeed > cbm.Teams[j].TotalSpeed
	})

	// Set initiative order
	cbm.InitiativeOrder = make([]*TeamInfo, len(cbm.Teams))
	copy(cbm.InitiativeOrder, cbm.Teams)

	cbm.logMessage("Initiative order determined:")
	for i, team := range cbm.InitiativeOrder {
		cbm.logMessage(fmt.Sprintf("  %d. %s Team (Speed: %d)",
			i+1, team.Team.String(), team.TotalSpeed))
	}
}

// startNextTeamTurn begins the next team's turn
func (cbm *TurnBasedCombatManager) startNextTeamTurn() {
	// Find next team that can act
	var nextTeam *TeamInfo

	for _, team := range cbm.InitiativeOrder {
		if !team.HasCompleted {
			nextTeam = team
			break
		}
	}

	// If no team can act, start new round
	if nextTeam == nil {
		cbm.startNewRound()
		return
	}

	// Deactivate previous team
	if cbm.ActiveTeam != nil {
		cbm.ActiveTeam.IsActive = false
	}

	// Activate new team
	cbm.ActiveTeam = nextTeam
	cbm.ActiveTeam.IsActive = true
	cbm.ActiveUnit = nil

	// Restore AP for all team members
	cbm.restoreTeamActionPoints(nextTeam)

	cbm.logMessage(fmt.Sprintf("Starting %s team turn (Round %d)",
		nextTeam.Team.String(), cbm.CurrentRound))

	// Notify UI of team change
	if cbm.StateChangeCallback != nil {
		cbm.StateChangeCallback(CombatPhaseTeamTurn)
	}
}

// restoreTeamActionPoints restores AP for all team members at turn start
func (cbm *TurnBasedCombatManager) restoreTeamActionPoints(team *TeamInfo) {
	for _, member := range team.Members {
		if actionPoints := member.ActionPoints(); actionPoints != nil {
			actionPoints.Restore()
		}
		if combatState := member.CombatState(); combatState != nil {
			combatState.StartTurn()
		}
	}
}

// startNewRound begins a new combat round
func (cbm *TurnBasedCombatManager) startNewRound() {
	cbm.CurrentRound++

	// Reset all teams for new round
	for _, team := range cbm.Teams {
		team.HasCompleted = false
		team.IsActive = false
	}

	cbm.logMessage(fmt.Sprintf("Starting Round %d", cbm.CurrentRound))

	// Start first team's turn
	cbm.startNextTeamTurn()
}

// Update processes combat logic each frame
func (cbm *TurnBasedCombatManager) Update() error {
	if !cbm.IsActive {
		return nil
	}

	switch cbm.Phase {
	case CombatPhaseTeamTurn:
		return cbm.updateTeamTurn()
	case CombatPhaseActionExecution:
		return cbm.updateActionExecution()
	case CombatPhaseEndTurn:
		return cbm.updateEndTurn()
	case CombatPhaseVictoryCheck:
		return cbm.updateVictoryCheck()
	}

	return nil
}

// updateTeamTurn handles team turn logic
func (cbm *TurnBasedCombatManager) updateTeamTurn() error {
	if cbm.ActiveTeam == nil {
		return fmt.Errorf("no active team")
	}

	// Check if team has any units that can still act
	canAct := false
	for _, member := range cbm.ActiveTeam.Members {
		if cbm.canUnitAct(member) {
			canAct = true
			break
		}
	}

	if !canAct {
		// Team is done, end their turn
		cbm.endTeamTurn()
		return nil
	}

	// Handle AI for enemy teams
	if cbm.ActiveTeam.Team == components.TeamEnemy {
		return cbm.processEnemyAI()
	}

	// For player teams, wait for player input
	return nil
}

// canUnitAct checks if a unit can perform actions
func (cbm *TurnBasedCombatManager) canUnitAct(entity *ecs.Entity) bool {
	// Check if unit is alive
	if stats := entity.RPGStats(); stats == nil || !stats.IsAlive() {
		return false
	}

	// Check if unit has action points
	if actionPoints := entity.ActionPoints(); actionPoints == nil || actionPoints.IsExhausted() {
		return false
	}

	// Check combat state
	if combatState := entity.CombatState(); combatState == nil || !combatState.CanAct {
		return false
	}

	return true
}

// processEnemyAI handles AI for enemy units
func (cbm *TurnBasedCombatManager) processEnemyAI() error {
	// Simple AI: Find first enemy that can act and try to attack adjacent players
	for _, enemy := range cbm.ActiveTeam.Members {
		if !cbm.canUnitAct(enemy) {
			continue
		}

		// Try to find adjacent player to attack
		if target := cbm.findAdjacentTarget(enemy, components.TeamPlayer); target != nil {
			// Create and execute attack action
			action := &CombatAction{
				Type:      ActionAttack,
				Actor:     enemy,
				Target:    target,
				APCost:    constants.AttackAPCost,
				Validated: true,
				Message:   fmt.Sprintf("%s attacks %s", cbm.getEntityName(enemy), cbm.getEntityName(target)),
			}

			return cbm.ExecuteAction(action)
		}

		// No adjacent targets, end turn for this enemy
		if actionPoints := enemy.ActionPoints(); actionPoints != nil {
			actionPoints.Current = 0 // Exhaust AP
		}
	}

	// No enemies can act, end team turn
	cbm.endTeamTurn()
	return nil
}

// findAdjacentTarget finds an adjacent enemy target
func (cbm *TurnBasedCombatManager) findAdjacentTarget(actor *ecs.Entity, targetTeam components.Team) *ecs.Entity {
	actorTransform := actor.Transform()
	if actorTransform == nil {
		return nil
	}

	actorPos := cbm.Grid.WorldToGrid(actorTransform.X, actorTransform.Y)
	neighbors := cbm.Grid.GetNeighbors(actorPos)

	for _, neighborPos := range neighbors {
		if unit := cbm.getUnitAtPosition(neighborPos); unit != nil {
			if combatState := unit.CombatState(); combatState != nil && combatState.Team == targetTeam {
				if stats := unit.RPGStats(); stats != nil && stats.IsAlive() {
					return unit
				}
			}
		}
	}

	return nil
}

// getUnitAtPosition returns the unit at a specific grid position
func (cbm *TurnBasedCombatManager) getUnitAtPosition(pos GridPos) *ecs.Entity {
	tile := cbm.Grid.GetTile(pos)
	if tile == nil || !tile.Occupied {
		return nil
	}

	// Find entity by checking all combat participants
	for _, team := range cbm.Teams {
		for _, member := range team.Members {
			if member.GetID() == tile.UnitID {
				return member
			}
		}
	}

	return nil
}

// endTeamTurn ends the current team's turn
func (cbm *TurnBasedCombatManager) endTeamTurn() {
	if cbm.ActiveTeam != nil {
		cbm.ActiveTeam.HasCompleted = true
		cbm.ActiveTeam.IsActive = false

		cbm.logMessage(fmt.Sprintf("%s team turn ended", cbm.ActiveTeam.Team.String()))
	}

	cbm.changePhase(CombatPhaseVictoryCheck)
}

// changePhase transitions to a new combat phase
func (cbm *TurnBasedCombatManager) changePhase(newPhase CombatPhase) {
	oldPhase := cbm.Phase
	cbm.Phase = newPhase

	cbm.logMessage(fmt.Sprintf("Combat phase: %s -> %s", oldPhase.String(), newPhase.String()))

	if cbm.StateChangeCallback != nil {
		cbm.StateChangeCallback(newPhase)
	}
}

// ExecuteAction executes a combat action
func (cbm *TurnBasedCombatManager) ExecuteAction(action *CombatAction) error {
	if action == nil {
		return fmt.Errorf("nil action")
	}

	cbm.PendingAction = action
	cbm.changePhase(CombatPhaseActionExecution)

	return nil
}

// updateActionExecution processes the pending action
func (cbm *TurnBasedCombatManager) updateActionExecution() error {
	if cbm.PendingAction == nil {
		cbm.changePhase(CombatPhaseEndTurn)
		return nil
	}

	action := cbm.PendingAction
	cbm.PendingAction = nil

	// Validate and execute the action
	if err := cbm.executeAction(action); err != nil {
		cbm.logMessage(fmt.Sprintf("Action failed: %v", err))
		cbm.changePhase(CombatPhaseTeamTurn)
		return err
	}

	cbm.changePhase(CombatPhaseEndTurn)
	return nil
}

// executeAction performs the actual action execution
func (cbm *TurnBasedCombatManager) executeAction(action *CombatAction) error {
	// Check if actor can afford the action
	actionPoints := action.Actor.ActionPoints()
	if actionPoints == nil || !actionPoints.CanAfford(action.APCost) {
		return fmt.Errorf("insufficient action points")
	}

	// Execute based on action type
	switch action.Type {
	case ActionMove:
		if err := cbm.executeMovement(action); err != nil {
			return err
		}
	case ActionAttack:
		if err := cbm.executeAttack(action); err != nil {
			return err
		}
	case ActionWait:
		// Wait action - just consume AP
	default:
		return fmt.Errorf("unsupported action type: %d", action.Type)
	}

	// Consume action points
	actionPoints.Spend(action.APCost)

	// Log the action
	cbm.logMessage(action.Message)

	return nil
}

// executeMovement handles movement actions
func (cbm *TurnBasedCombatManager) executeMovement(action *CombatAction) error {
	// Get actor's transform component
	transform := action.Actor.Transform()
	if transform == nil {
		return fmt.Errorf("actor has no transform component")
	}

	// Calculate current grid position using the same method as the main engine
	currentPos := cbm.worldToGridPos(transform.X, transform.Y)
	targetPos := action.TargetPos

	// Validate the movement
	if err := cbm.validateMovement(currentPos, targetPos, action.Actor); err != nil {
		return fmt.Errorf("movement validation failed: %v", err)
	}

	// Calculate movement distance for AP cost validation
	distance := cbm.Grid.CalculateDistance(currentPos, targetPos)
	expectedAPCost := distance * constants.MovementAPCost
	
	if action.APCost != expectedAPCost {
		return fmt.Errorf("AP cost mismatch: expected %d, got %d", expectedAPCost, action.APCost)
	}

	// Clear occupancy at current position
	cbm.Grid.SetOccupied(currentPos, false, "")
	
	// Convert target grid position to world coordinates
	worldX, worldY := cbm.Grid.GridToWorld(targetPos)
	
	// Add grid offset to match the coordinate system used throughout the game
	transform.X = worldX + constants.GridOffsetX
	transform.Y = worldY + constants.GridOffsetY
	
	// Set occupancy at new position
	cbm.Grid.SetOccupied(targetPos, true, action.Actor.GetID())

	// Update RPG stats if the actor has movement tracking
	if stats := action.Actor.RPGStats(); stats != nil {
		// Consume moves from the legacy movement system if it exists
		if stats.MovesRemaining >= distance {
			stats.MovesRemaining -= distance
		}
		
		// Add move to history for potential undo functionality
		moveRecord := components.MoveRecord{
			FromX: currentPos.X,
			FromZ: currentPos.Y,
			ToX:   targetPos.X,
			ToZ:   targetPos.Y,
			Distance: distance,
		}
		stats.MoveHistory = append(stats.MoveHistory, moveRecord)
	}

	// Log the movement
	actorName := cbm.getEntityName(action.Actor)
	cbm.logMessage(fmt.Sprintf("%s moved from (%d,%d) to (%d,%d) [Distance: %d, AP Cost: %d]",
		actorName, currentPos.X, currentPos.Y, targetPos.X, targetPos.Y, distance, action.APCost))

	return nil
}

// worldToGridPos converts world coordinates to grid position using the same logic as the main engine
func (cbm *TurnBasedCombatManager) worldToGridPos(worldX, worldY float64) GridPos {
	offsetX, offsetY := constants.GridOffsetX, constants.GridOffsetY
	tileSize := float64(cbm.Grid.TileSize)

	// Remove offset and convert to grid coordinates
	// This is the exact inverse of: worldX = gridX * tileSize + offsetX
	gridX := int((worldX - offsetX) / tileSize)
	gridY := int((worldY - offsetY) / tileSize)

	return GridPos{X: gridX, Y: gridY}
}

// validateMovement checks if a movement is valid
func (cbm *TurnBasedCombatManager) validateMovement(currentPos, targetPos GridPos, actor *ecs.Entity) error {
	// Check if target position is within grid bounds
	if !cbm.Grid.IsValidPosition(targetPos) {
		return fmt.Errorf("target position (%d,%d) is out of bounds", targetPos.X, targetPos.Y)
	}

	// Check if we're trying to move to the same position
	if currentPos.X == targetPos.X && currentPos.Y == targetPos.Y {
		return fmt.Errorf("already at target position (%d,%d)", targetPos.X, targetPos.Y)
	}

	// Check if target position is passable and not occupied
	if !cbm.Grid.IsPassable(targetPos) {
		// Get more specific error information
		tile := cbm.Grid.GetTile(targetPos)
		if tile == nil {
			return fmt.Errorf("target tile (%d,%d) does not exist", targetPos.X, targetPos.Y)
		}
		if !tile.Passable {
			return fmt.Errorf("target tile (%d,%d) is not passable", targetPos.X, targetPos.Y)
		}
		if tile.Occupied {
			return fmt.Errorf("target tile (%d,%d) is occupied by unit %s", targetPos.X, targetPos.Y, tile.UnitID)
		}
	}

	// Check if actor has enough movement range (if using legacy movement system)
	if stats := actor.RPGStats(); stats != nil {
		distance := cbm.Grid.CalculateDistance(currentPos, targetPos)
		if stats.MovesRemaining < distance {
			return fmt.Errorf("insufficient movement range: need %d, have %d", distance, stats.MovesRemaining)
		}
	}

	return nil
}

// executeAttack handles attack actions
func (cbm *TurnBasedCombatManager) executeAttack(action *CombatAction) error {
	attackerStats := action.Actor.RPGStats()
	targetStats := action.Target.RPGStats()

	if attackerStats == nil || targetStats == nil {
		return fmt.Errorf("missing stats for combat")
	}

	// Calculate damage (simple for now)
	damage := attackerStats.Attack - targetStats.Defense
	if damage < 1 {
		damage = 1 // Minimum damage
	}

	// Apply damage
	targetStats.TakeDamage(damage)

	cbm.logMessage(fmt.Sprintf("%s deals %d damage to %s (HP: %d/%d)",
		attackerStats.Name, damage, targetStats.Name,
		targetStats.CurrentHP, targetStats.MaxHP))

	// Check if target died
	if !targetStats.IsAlive() {
		cbm.logMessage(fmt.Sprintf("%s has been defeated!", targetStats.Name))
		// TODO: Remove from grid and mark as dead
	}

	return nil
}

// updateEndTurn handles end of turn processing
func (cbm *TurnBasedCombatManager) updateEndTurn() error {
	cbm.changePhase(CombatPhaseVictoryCheck)
	return nil
}

// updateVictoryCheck checks for victory conditions
func (cbm *TurnBasedCombatManager) updateVictoryCheck() error {
	result := cbm.checkVictoryConditions()

	if result != CombatResultOngoing {
		cbm.Result = result
		cbm.changePhase(CombatPhaseEnded)
		cbm.IsActive = false

		cbm.logMessage(fmt.Sprintf("Combat ended: %s", result.String()))
		return nil
	}

	// Combat continues, check if we need to start next team turn
	cbm.changePhase(CombatPhaseTeamTurn)
	cbm.startNextTeamTurn()

	return nil
}

// checkVictoryConditions determines if combat should end
func (cbm *TurnBasedCombatManager) checkVictoryConditions() CombatResult {
	playerAlive := false
	enemyAlive := false

	for _, team := range cbm.Teams {
		for _, member := range team.Members {
			if stats := member.RPGStats(); stats != nil && stats.IsAlive() {
				if team.Team == components.TeamPlayer {
					playerAlive = true
				} else if team.Team == components.TeamEnemy {
					enemyAlive = true
				}
			}
		}
	}

	if !playerAlive {
		return CombatResultEnemyVictory
	} else if !enemyAlive {
		return CombatResultPlayerVictory
	}

	return CombatResultOngoing
}

// Helper methods

// getEntityName safely gets an entity's name
func (cbm *TurnBasedCombatManager) getEntityName(entity *ecs.Entity) string {
	if stats := entity.RPGStats(); stats != nil {
		return stats.Name
	}
	return entity.GetID()
}

// logMessage sends a message to the UI if callback is set
func (cbm *TurnBasedCombatManager) logMessage(message string) {
	if cbm.DebugMode {
		fmt.Printf("[Combat] %s\n", message)
	}

	if cbm.MessageCallback != nil {
		cbm.MessageCallback(message)
	}
}

// GetActiveTeam returns the currently active team
func (cbm *TurnBasedCombatManager) GetActiveTeam() *TeamInfo {
	return cbm.ActiveTeam
}

// GetPhase returns the current combat phase
func (cbm *TurnBasedCombatManager) GetPhase() CombatPhase {
	return cbm.Phase
}

// IsPlayerTurn returns true if it's the player team's turn
func (cbm *TurnBasedCombatManager) IsPlayerTurn() bool {
	return cbm.ActiveTeam != nil && cbm.ActiveTeam.Team == components.TeamPlayer
}

// GetResult returns the combat result
func (cbm *TurnBasedCombatManager) GetResult() CombatResult {
	return cbm.Result
}

// Action Creation Helper Methods

// CreateMoveAction creates a validated movement action
func (cbm *TurnBasedCombatManager) CreateMoveAction(actor *ecs.Entity, targetPos GridPos) (*CombatAction, error) {
	if actor == nil {
		return nil, fmt.Errorf("actor is nil")
	}

	// Get current position
	transform := actor.Transform()
	if transform == nil {
		return nil, fmt.Errorf("actor has no transform component")
	}

	currentPos := cbm.worldToGridPos(transform.X, transform.Y)
	
	// Calculate distance and AP cost
	distance := cbm.Grid.CalculateDistance(currentPos, targetPos)
	apCost := distance * constants.MovementAPCost

	// Create action
	action := &CombatAction{
		Type:      ActionMove,
		Actor:     actor,
		Target:    nil,
		TargetPos: targetPos,
		APCost:    apCost,
		Validated: false,
		Message:   fmt.Sprintf("%s moves to (%d,%d)", cbm.getEntityName(actor), targetPos.X, targetPos.Y),
	}

	// Validate the action
	if err := cbm.validateMovement(currentPos, targetPos, actor); err != nil {
		action.Validated = false
		action.Message = fmt.Sprintf("Invalid move: %v", err)
		return action, err
	}

	action.Validated = true
	return action, nil
}

// CreateAttackAction creates a validated attack action
func (cbm *TurnBasedCombatManager) CreateAttackAction(actor *ecs.Entity, target *ecs.Entity) (*CombatAction, error) {
	if actor == nil {
		return nil, fmt.Errorf("actor is nil")
	}
	if target == nil {
		return nil, fmt.Errorf("target is nil")
	}

	// Create action
	action := &CombatAction{
		Type:      ActionAttack,
		Actor:     actor,
		Target:    target,
		TargetPos: GridPos{}, // Not used for attacks
		APCost:    constants.AttackAPCost,
		Validated: false,
		Message:   fmt.Sprintf("%s attacks %s", cbm.getEntityName(actor), cbm.getEntityName(target)),
	}

	// Validate attack (check adjacency, etc.)
	if err := cbm.validateAttack(actor, target); err != nil {
		action.Validated = false
		action.Message = fmt.Sprintf("Invalid attack: %v", err)
		return action, err
	}

	action.Validated = true
	return action, nil
}

// CreateEndTurnAction creates an end turn action
func (cbm *TurnBasedCombatManager) CreateEndTurnAction(actor *ecs.Entity) (*CombatAction, error) {
	if actor == nil {
		return nil, fmt.Errorf("actor is nil")
	}

	action := &CombatAction{
		Type:      ActionWait, // Using ActionWait to end turn and consume remaining AP
		Actor:     actor,
		Target:    nil,
		TargetPos: GridPos{},
		APCost:    0, // End turn is free, but will exhaust all remaining AP
		Validated: true,
		Message:   fmt.Sprintf("%s ends turn", cbm.getEntityName(actor)),
	}

	return action, nil
}

// validateAttack checks if an attack action is valid
func (cbm *TurnBasedCombatManager) validateAttack(actor *ecs.Entity, target *ecs.Entity) error {
	// Check if target has stats and is alive
	targetStats := target.RPGStats()
	if targetStats == nil {
		return fmt.Errorf("target has no stats")
	}
	if !targetStats.IsAlive() {
		return fmt.Errorf("target is already dead")
	}

	// Check if units are on different teams
	actorCombat := actor.CombatState()
	targetCombat := target.CombatState()
	if actorCombat == nil || targetCombat == nil {
		return fmt.Errorf("missing combat state components")
	}
	if actorCombat.Team == targetCombat.Team {
		return fmt.Errorf("cannot attack ally")
	}

	// Get positions of actor and target
	actorTransform := actor.Transform()
	targetTransform := target.Transform()
	if actorTransform == nil || targetTransform == nil {
		return fmt.Errorf("missing transform components")
	}

	// Check if targets are adjacent (range = 1 for now)
	actorGridPos := cbm.worldToGridPos(actorTransform.X, actorTransform.Y)
	targetGridPos := cbm.worldToGridPos(targetTransform.X, targetTransform.Y)
	
	distance := cbm.Grid.CalculateDistance(actorGridPos, targetGridPos)
	if distance > 1 {
		return fmt.Errorf("target out of range (distance: %d, max range: 1)", distance)
	}

	return nil
}

// GetValidMovesForUnit returns all valid movement positions for a unit
func (cbm *TurnBasedCombatManager) GetValidMovesForUnit(actor *ecs.Entity) []GridPos {
	if actor == nil {
		return []GridPos{}
	}

	transform := actor.Transform()
	actionPoints := actor.ActionPoints()
	if transform == nil || actionPoints == nil {
		return []GridPos{}
	}

	currentPos := cbm.worldToGridPos(transform.X, transform.Y)
	maxDistance := actionPoints.Current / constants.MovementAPCost

	if maxDistance <= 0 {
		return []GridPos{}
	}

	// Calculate all positions within movement range
	validMoves := []GridPos{}
	
	for x := 0; x < cbm.Grid.Width; x++ {
		for y := 0; y < cbm.Grid.Height; y++ {
			targetPos := GridPos{X: x, Y: y}
			distance := cbm.Grid.CalculateDistance(currentPos, targetPos)
			
			// Check if position is within range and passable
			if distance <= maxDistance && distance > 0 {
				if cbm.Grid.IsPassable(targetPos) {
					validMoves = append(validMoves, targetPos)
				}
			}
		}
	}

	return validMoves
}

// GetValidAttackTargetsForUnit returns all valid attack targets for a unit
func (cbm *TurnBasedCombatManager) GetValidAttackTargetsForUnit(actor *ecs.Entity) []*ecs.Entity {
	if actor == nil {
		return []*ecs.Entity{}
	}

	actionPoints := actor.ActionPoints()
	if actionPoints == nil || actionPoints.Current < constants.AttackAPCost {
		return []*ecs.Entity{}
	}

	validTargets := []*ecs.Entity{}

	// Check all combat participants
	for _, team := range cbm.Teams {
		for _, member := range team.Members {
			if err := cbm.validateAttack(actor, member); err == nil {
				validTargets = append(validTargets, member)
			}
		}
	}

	return validTargets
}
