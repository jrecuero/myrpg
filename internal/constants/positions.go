package constants

// Initial Entity Positions
// These define where entities spawn in exploration mode
const (
	// Player Starting Positions (X, Y coordinates)
	Player1StartX = 100.0
	Player1StartY = 150.0
	Player2StartX = 150.0
	Player2StartY = 150.0
	Player3StartX = 200.0
	Player3StartY = 150.0

	// Enemy Starting Positions
	Enemy1StartX = 300.0
	Enemy1StartY = 250.0
	Enemy2StartX = 350.0
	Enemy2StartY = 300.0
	Enemy3StartX = 400.0
	Enemy3StartY = 250.0
	Enemy4StartX = 450.0
	Enemy4StartY = 300.0
	Enemy5StartX = 500.0
	Enemy5StartY = 250.0
	Enemy6StartX = 550.0
	Enemy6StartY = 300.0
)

// Exploration Mode Positioning
const (
	// Party Member Spacing
	PartyMemberOffset = 5.0 // Small offset between party members

	// Enemy Positioning
	EnemyBaseDistance        = 100.0 // Base distance from player
	EnemySpacing             = 20.0  // Additional spacing per enemy
	EnemyPositionMultiplierX = 1.0
	EnemyPositionMultiplierY = 0.5
	EnemyVerticalSpacing     = 30.0

	// Boundary Constraints for Enemy Positioning
	EnemyMinX = 50.0
	EnemyMaxX = 700.0
	EnemyMinY = 150.0
	EnemyMaxY = 450.0
)
