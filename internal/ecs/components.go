// Package ecs provides an Entity-Component-System (ECS) framework for game development.
// It defines core components such as Transform, SpriteComponent, and ColliderComponent
// that can be attached to entities to define their properties and behaviors.
// The ECS architecture allows for flexible and modular game object management.
// Each component struct includes fields relevant to its purpose, along with constructor functions
// for easy instantiation. The Transform component handles position and size, the SpriteComponent
// manages visual representation, and the ColliderComponent defines collision properties.
package ecs

import (
	"image"

	"github.com/jrecuero/myrpg/internal/gfx"
)

// Transform represents the position and size of an entity.
// It is a common component used for rendering and collision detection.
// X and Y represent the position of the entity.
// Width and Height represent the size of the entity.
type Transform struct {
	X      float64 // X position
	Y      float64 // Y position
	Width  int     // Width of the entity
	Height int     // Height of the entity
}

// SpriteComponent represents the visual representation of an entity using a sprite.
// It includes the sprite itself and rendering properties such as scale and offset.
// Sprite is a pointer to the gfx.Sprite used for rendering.
// Scale is a scaling factor applied to the sprite when rendering.
// OffsetX and OffsetY are offsets applied to the sprite's position when rendering.
type SpriteComponent struct {
	Sprite  *gfx.Sprite // The sprite associated with the entity
	Scale   float64     // Scale factor for rendering the sprite
	OffsetX float64     // X offset for rendering the sprite
	OffsetY float64     // Y offset for rendering the sprite
}

// ColliderComponent represents the collision properties of an entity.
// It defines the size and behavior of the collider used in collision detection.
// Solid indicates if the collider is solid (affects collision response).
// Width and Height define the size of the collider.
// OffsetX and OffsetY are offsets applied to the collider's position relative to the entity's position.
type ColliderComponent struct {
	Solid   bool // Indicates if the collider is solid (affects collision response)
	Width   int  // Width of the collider
	Height  int  // Height of the collider
	OffsetX int  // X offset of the collider relative to the entity's position
	OffsetY int  // Y offset of the collider relative to the entity's position
}

// NewTransform creates a new Transform component with the specified position and size.
// x and y are the initial position of the entity.
// width and height are the size of the entity.
// returns a pointer to the newly created Transform component.
func NewTransform(x, y float64, width, height int) *Transform {
	return &Transform{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// NewColliderComponent creates a new ColliderComponent with the specified properties.
// solid indicates if the collider is solid (affects collision response).
// width and height define the size of the collider.
// offsetX and offsetY are offsets applied to the collider's position relative to the entity's position.
// returns a pointer to the newly created ColliderComponent.
func NewSpriteComponent(sprite *gfx.Sprite, scale, offsetX, offsetY float64) *SpriteComponent {
	return &SpriteComponent{
		Sprite:  sprite,
		Scale:   scale,
		OffsetX: offsetX,
		OffsetY: offsetY,
	}
}

// NewColliderComponent creates a new ColliderComponent with the specified properties.
// solid indicates if the collider is solid (affects collision response).
// width and height define the size of the collider.
// offsetX and offsetY are offsets applied to the collider's position relative to the entity's position.
// returns a pointer to the newly created ColliderComponent.
func NewColliderComponent(solid bool, width, height, offsetX, offsetY int) *ColliderComponent {
	return &ColliderComponent{
		Solid:   solid,
		Width:   width,
		Height:  height,
		OffsetX: offsetX,
		OffsetY: offsetY,
	}
}

// Bounds returns the bounding rectangle of the Transform.
// It uses the X, Y, Width, and Height fields of the Transform to create the rectangle.
// returns an image.Rectangle representing the bounding box of the Transform.
func (t *Transform) Bounds() image.Rectangle {
	return image.Rect(int(t.X), int(t.Y), int(t.X)+t.Width, int(t.Y)+t.Height)
}

// JobType represents different character classes/jobs in the RPG
type JobType int

const (
	JobWarrior JobType = iota // High HP and Defense, moderate Attack
	JobMage                   // High Magic Power, low HP and Defense
	JobRogue                  // High Speed and Crit, moderate HP
	JobCleric                 // High Magic Defense and Healing, moderate HP
	JobArcher                 // High Range Attack and Speed, low Defense
)

// String returns the string representation of a JobType
func (j JobType) String() string {
	switch j {
	case JobWarrior:
		return "Warrior"
	case JobMage:
		return "Mage"
	case JobRogue:
		return "Rogue"
	case JobCleric:
		return "Cleric"
	case JobArcher:
		return "Archer"
	default:
		return "Unknown"
	}
}

// RPGStatsComponent represents the RPG statistics of a character entity.
// It includes core stats like health, level, experience, and combat attributes.
type RPGStatsComponent struct {
	// Core Stats
	Level      int     // Current character level (1-99)
	Experience int     // Current experience points
	ExpToNext  int     // Experience needed for next level
	
	// Health and Mana
	CurrentHP  int     // Current health points
	MaxHP      int     // Maximum health points
	CurrentMP  int     // Current mana/magic points
	MaxMP      int     // Maximum mana/magic points
	
	// Combat Stats
	Attack     int     // Physical attack power
	Defense    int     // Physical defense
	MagicAttack int    // Magic attack power
	MagicDefense int   // Magic defense
	Speed      int     // Speed/agility stat
	Accuracy   int     // Hit accuracy percentage (0-100)
	CritRate   int     // Critical hit rate percentage (0-100)
	
	// Character Info
	Job        JobType // Character class/job
	Name       string  // Character display name
}

// NewRPGStatsComponent creates a new RPG stats component with default values for a specific job
func NewRPGStatsComponent(name string, job JobType, level int) *RPGStatsComponent {
	stats := &RPGStatsComponent{
		Name:         name,
		Job:          job,
		Level:        level,
		Experience:   0,
		ExpToNext:    level * 100, // Simple EXP formula
		Speed:        50,
		Accuracy:     85,
		CritRate:     5,
	}
	
	// Apply job-specific base stats
	switch job {
	case JobWarrior:
		stats.MaxHP = 120 + (level-1)*15
		stats.MaxMP = 30 + (level-1)*3
		stats.Attack = 18 + (level-1)*3
		stats.Defense = 15 + (level-1)*2
		stats.MagicAttack = 8 + (level-1)*1
		stats.MagicDefense = 12 + (level-1)*2
		stats.Speed = 40 + (level-1)*1
		
	case JobMage:
		stats.MaxHP = 80 + (level-1)*8
		stats.MaxMP = 100 + (level-1)*12
		stats.Attack = 10 + (level-1)*1
		stats.Defense = 8 + (level-1)*1
		stats.MagicAttack = 25 + (level-1)*4
		stats.MagicDefense = 20 + (level-1)*3
		stats.Speed = 45 + (level-1)*2
		
	case JobRogue:
		stats.MaxHP = 95 + (level-1)*10
		stats.MaxMP = 50 + (level-1)*5
		stats.Attack = 20 + (level-1)*3
		stats.Defense = 10 + (level-1)*1
		stats.MagicAttack = 12 + (level-1)*1
		stats.MagicDefense = 10 + (level-1)*1
		stats.Speed = 65 + (level-1)*3
		stats.CritRate = 15 + (level-1)/2
		
	case JobCleric:
		stats.MaxHP = 100 + (level-1)*12
		stats.MaxMP = 90 + (level-1)*10
		stats.Attack = 12 + (level-1)*1
		stats.Defense = 12 + (level-1)*2
		stats.MagicAttack = 20 + (level-1)*3
		stats.MagicDefense = 25 + (level-1)*4
		stats.Speed = 35 + (level-1)*1
		
	case JobArcher:
		stats.MaxHP = 85 + (level-1)*9
		stats.MaxMP = 40 + (level-1)*4
		stats.Attack = 22 + (level-1)*3
		stats.Defense = 9 + (level-1)*1
		stats.MagicAttack = 10 + (level-1)*1
		stats.MagicDefense = 8 + (level-1)*1
		stats.Speed = 60 + (level-1)*3
		stats.Accuracy = 95 + (level-1)/5
	}
	
	// Set current HP/MP to max
	stats.CurrentHP = stats.MaxHP
	stats.CurrentMP = stats.MaxMP
	
	return stats
}

// IsAlive returns true if the character has HP remaining
func (r *RPGStatsComponent) IsAlive() bool {
	return r.CurrentHP > 0
}

// TakeDamage reduces current HP by the specified amount
func (r *RPGStatsComponent) TakeDamage(damage int) {
	r.CurrentHP -= damage
	if r.CurrentHP < 0 {
		r.CurrentHP = 0
	}
}

// Heal increases current HP by the specified amount, not exceeding max HP
func (r *RPGStatsComponent) Heal(amount int) {
	r.CurrentHP += amount
	if r.CurrentHP > r.MaxHP {
		r.CurrentHP = r.MaxHP
	}
}

// UseMana reduces current MP by the specified amount
func (r *RPGStatsComponent) UseMana(amount int) bool {
	if r.CurrentMP >= amount {
		r.CurrentMP -= amount
		return true
	}
	return false
}

// RestoreMana increases current MP by the specified amount, not exceeding max MP
func (r *RPGStatsComponent) RestoreMana(amount int) {
	r.CurrentMP += amount
	if r.CurrentMP > r.MaxMP {
		r.CurrentMP = r.MaxMP
	}
}

// GainExperience adds experience and handles level ups
func (r *RPGStatsComponent) GainExperience(exp int) bool {
	r.Experience += exp
	if r.Experience >= r.ExpToNext {
		r.LevelUp()
		return true // Level up occurred
	}
	return false
}

// LevelUp increases the character's level and recalculates stats
func (r *RPGStatsComponent) LevelUp() {
	r.Level++
	r.Experience -= r.ExpToNext
	r.ExpToNext = r.Level * 100
	
	// Recalculate stats based on new level (this is a simple approach)
	// In a real game, you might want more sophisticated stat growth
	oldMaxHP := r.MaxHP
	oldMaxMP := r.MaxMP
	
	// Create a new stats component with updated level to get new max values
	newStats := NewRPGStatsComponent(r.Name, r.Job, r.Level)
	
	// Update max stats
	r.MaxHP = newStats.MaxHP
	r.MaxMP = newStats.MaxMP
	r.Attack = newStats.Attack
	r.Defense = newStats.Defense
	r.MagicAttack = newStats.MagicAttack
	r.MagicDefense = newStats.MagicDefense
	r.Speed = newStats.Speed
	
	// Restore HP/MP proportionally or fully heal on level up
	r.CurrentHP += (r.MaxHP - oldMaxHP) // Add the HP increase
	r.CurrentMP += (r.MaxMP - oldMaxMP) // Add the MP increase
}
