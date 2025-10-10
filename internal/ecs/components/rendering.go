package components

import (
	"github.com/jrecuero/myrpg/internal/gfx"
)

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

// NewSpriteComponent creates a new SpriteComponent with the specified properties.
// sprite is the gfx.Sprite to be used for rendering.
// scale is the scaling factor applied to the sprite when rendering.
// offsetX and offsetY are offsets applied to the sprite's position when rendering.
// returns a pointer to the newly created SpriteComponent.
func NewSpriteComponent(sprite *gfx.Sprite, scale, offsetX, offsetY float64) *SpriteComponent {
	return &SpriteComponent{
		Sprite:  sprite,
		Scale:   scale,
		OffsetX: offsetX,
		OffsetY: offsetY,
	}
}