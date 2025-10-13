package constants

// Entity and Sprite Constants
const (
	// Entity Dimensions (standard sprite sizes)
	EntityWidth  = 32
	EntityHeight = 32
	
	// Collision and Bounds
	EntityColliderWidth  = 32
	EntityColliderHeight = 32
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
	MaxPartyMembers  = 6
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
	DefaultAnimationScale = 1.0
	DefaultAnimationOffsetX = 0
	DefaultAnimationOffsetY = 0
)
