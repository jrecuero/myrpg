package constants

// Entity and Sprite Constants
const (
	// Entity Dimensions (standard sprite sizes)
	EntityWidth  = 32
	EntityHeight = 32

	// Collision and Bounds
	EntityColliderWidth   = 32
	EntityColliderHeight  = 32
	EntityColliderOffsetX = 0
	EntityColliderOffsetY = 0
)

// Movement Constants
const (
	// Exploration Mode Movement
	PlayerSpeed = 2.0 // Pixels per frame movement speed

	// Tactical Movement Ranges by Job Class
	WarriorMoveRange = 3
	MageMoveRange    = 2
	RogueMoveRange   = 5
	ClericMoveRange  = 3
	ArcherMoveRange  = 4
)

// Party and Combat Constants
const (
	// Party Management
	MaxPartyMembers          = 6
	DefaultActivePlayerIndex = 0

	// Enemy Group Detection
	EnemyGroupRange = 150.0 // Pixel distance for forming enemy groups

	// Deployment Zones (as fraction of grid width)
	PlayerZoneWidth = 3 // Left 1/3 of grid width
	EnemyZoneStart  = 2 // Right 1/3 starts at 2/3 width
)

// Animation Constants
const (
	// Animation Speeds (milliseconds)
	IdleAnimationDuration   = 200
	WalkAnimationDuration   = 150
	AttackAnimationDuration = 100

	// Animation Behavior
	DefaultAnimationLoop = true
	AttackAnimationLoop  = true

	// Animation Scaling
	DefaultAnimationScale   = 1.0
	DefaultAnimationOffsetX = 0
	DefaultAnimationOffsetY = 0
)

// Action Point Constants for Turn-Based Combat
const (
	// Action Point Maximums by Job Class
	WarriorMaxAP = 4
	MageMaxAP    = 3
	RogueMaxAP   = 5
	ClericMaxAP  = 4
	ArcherMaxAP  = 4

	// Action Point Costs
	MovementAPCost = 1 // 1 AP per tile moved
	AttackAPCost   = 2 // 2 AP per attack action
	ItemAPCost     = 1 // 1 AP per item used
	EndTurnAPCost  = 0 // Free action
	WaitAPCost     = 0 // Free action
)

// Event System Color Constants
// These colors are used for event entities when no custom sprite is provided
var (
	// Event Type Colors (RGB values)
	EventColorBattle   = [3]uint8{200, 50, 50}   // Red-ish color for battle events
	EventColorDialog   = [3]uint8{50, 150, 200}  // Blue-ish color for dialog events
	EventColorChest    = [3]uint8{200, 200, 50}  // Yellow-ish color for chest events
	EventColorDoor     = [3]uint8{150, 100, 50}  // Brown-ish color for door events
	EventColorTrap     = [3]uint8{150, 50, 150}  // Purple-ish color for trap events
	EventColorInfo     = [3]uint8{100, 200, 100} // Green-ish color for info events
	EventColorQuest    = [3]uint8{255, 165, 0}   // Orange color for quest events
	EventColorCutscene = [3]uint8{100, 100, 200} // Purple-blue for cutscene events
	EventColorShop     = [3]uint8{200, 200, 200} // Silver-ish color for shop events
	EventColorRest     = [3]uint8{100, 150, 255} // Light blue for rest events
)
