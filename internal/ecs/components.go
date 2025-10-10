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
