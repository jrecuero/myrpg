// Package gfx provides functions for drawing sprites and handling graphics operations.
// It includes functions to draw sprites at specified positions with scaling options,
// as well as functions to get the bounding rectangle of a sprite for collision detection
// and other purposes.
package gfx

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Sprite represents a drawable image in the game.
// It contains the image and its dimensions.
// Img is the image representing the sprite.
// W is the width of the sprite.
// H is the height of the sprite.
type Sprite struct {
	Img *ebiten.Image // The image representing the sprite
	W   int           // Width of the sprite
	H   int           // Height of the sprite
}

// LoadImage loads an image from the specified file path.
// path is the file path to the image.
// w and h are the desired width and height of the image (not used in this function).
// returns a pointer to the loaded ebiten.Image and an error if any occurs.
func LoadImage(path string, w int, h int) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// NewSpriteFromFile creates a new Sprite by loading an image from the specified file path.
// path is the file path to the image.
// w and h are the desired width and height of the sprite.
// returns a pointer to the newly created Sprite and an error if any occurs during image loading.
func NewSpriteFromFile(path string, w int, h int) (*Sprite, error) {
	img, err := LoadImage(path, w, h)
	if err != nil {
		return nil, err
	}
	return &Sprite{
		Img: img,
		W:   w,
		H:   h,
	}, nil
}

func (s *Sprite) Bounds(x float64, y float64) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: int(x), Y: int(y)},
		Max: image.Point{X: int(x) + s.W, Y: int(y) + s.H},
	}
}
