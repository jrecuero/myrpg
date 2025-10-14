// Package gfx provides functions for drawing sprites and handling graphics operations.
// It includes functions to draw sprites at specified positions with scaling options,
// as well as functions to get the bounding rectangle of a sprite for collision detection
// and other purposes.
package gfx

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawSprite draws the given sprite at the specified position with an optional scale.
// screen is the target ebiten image where the sprite will be drawn.
// s is the sprite to be drawn.
// x and y are the coordinates where the sprite will be drawn.
// scale is the scaling factor for the sprite (1.0 means no scaling).
func DrawSprite(screen *ebiten.Image, s *Sprite, x float64, y float64, scale float64) {
	if s == nil || s.Img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	if scale != 0 && scale != 1 {
		op.GeoM.Scale(scale, scale)
	}
	op.GeoM.Translate(x, y)
	screen.DrawImage(s.Img, op)
}

// DrawSpriteCentered draws the given sprite centered at the specified position with an optional scale.
// screen is the target ebiten image where the sprite will be drawn.
// s is the sprite to be drawn.
// x and y are the coordinates where the center of the sprite will be drawn.
// scale is the scaling factor for the sprite (1.0 means no scaling).
func DrawSpriteCentered(screen *ebiten.Image, s *Sprite, x float64, y float64, scale float64) {
	if s == nil || s.Img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	sw := float64(s.W) * scale
	sh := float64(s.H) * scale
	op.GeoM.Translate(-sw/2, -sh/2)
	if scale != 0 && scale != 1 {
		// Scale after translation to keep it centered
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Scale(scale, scale)
		op2.GeoM.Translate(x, y)
		// Compose transformations: GeoM is applied in order: scale then translate.
		// so we pre-scale the image and then translate it to (x,y)
		op.GeoM.Concat(op2.GeoM)
	} else {
		op.GeoM.Translate(x, y)
	}
	screen.DrawImage(s.Img, op)
}

// DrawSpriteClipped draws a sprite clipped to the specified bounds
// This is useful for constraining background sprites to the game world area
func DrawSpriteClipped(screen *ebiten.Image, s *Sprite, x float64, y float64, scale float64, clipX, clipY, clipWidth, clipHeight int) {
	if s == nil || s.Img == nil {
		return
	}

	// Create a sub-image of the screen for clipping
	clipRect := image.Rect(clipX, clipY, clipX+clipWidth, clipY+clipHeight)
	clippedScreen := screen.SubImage(clipRect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	if scale != 0 && scale != 1 {
		op.GeoM.Scale(scale, scale)
	}
	// Adjust translation for the clipped area
	op.GeoM.Translate(x-float64(clipX), y-float64(clipY))
	clippedScreen.DrawImage(s.Img, op)
}

// BoundsRect returns the bounding rectangle of a sprite at the specified position with an optional scale.
// x and y are the coordinates where the sprite is drawn.
// w and h are the width and height of the sprite.
// scale is the scaling factor for the sprite (1.0 means no scaling).
// returns an image.Rectangle representing the bounding box.
func BoundsRect(x float64, y float64, w int, h int, scale float64) image.Rectangle {
	sw := float64(w) * scale
	sh := float64(h) * scale
	return image.Rect(int(x), int(y), int(x+sw), int(y+sh))
}
