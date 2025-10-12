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

// SpriteSheet represents a collection of sprites arranged in a grid
type SpriteSheet struct {
	Image        *ebiten.Image // The sprite sheet image
	SpriteWidth  int           // Width of each individual sprite
	SpriteHeight int           // Height of each individual sprite
	Columns      int           // Number of columns in the sprite sheet
	Rows         int           // Number of rows in the sprite sheet
}

// NewSpriteSheetFromFile creates a new sprite sheet from an image file
func NewSpriteSheetFromFile(path string, spriteWidth, spriteHeight int) (*SpriteSheet, error) {
	img, err := LoadImage(path, 0, 0) // Load the full image
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	columns := bounds.Dx() / spriteWidth
	rows := bounds.Dy() / spriteHeight

	return &SpriteSheet{
		Image:        img,
		SpriteWidth:  spriteWidth,
		SpriteHeight: spriteHeight,
		Columns:      columns,
		Rows:         rows,
	}, nil
}

// GetSprite extracts a specific sprite from the sprite sheet by index
func (ss *SpriteSheet) GetSprite(index int) (*Sprite, error) {
	if index < 0 || index >= ss.Columns*ss.Rows {
		return nil, image.ErrFormat // Invalid index
	}

	row := index / ss.Columns
	col := index % ss.Columns

	x := col * ss.SpriteWidth
	y := row * ss.SpriteHeight

	// Create a sub-image for the sprite
	subImg := ss.Image.SubImage(image.Rect(x, y, x+ss.SpriteWidth, y+ss.SpriteHeight)).(*ebiten.Image)

	return &Sprite{
		Img: subImg,
		W:   ss.SpriteWidth,
		H:   ss.SpriteHeight,
	}, nil
}

// GetSprites extracts a range of sprites from the sprite sheet
func (ss *SpriteSheet) GetSprites(startIndex, count int) ([]*Sprite, error) {
	sprites := make([]*Sprite, 0, count)

	for i := 0; i < count; i++ {
		sprite, err := ss.GetSprite(startIndex + i)
		if err != nil {
			return nil, err
		}
		sprites = append(sprites, sprite)
	}

	return sprites, nil
}

// GetAllSprites extracts all sprites from the sprite sheet
func (ss *SpriteSheet) GetAllSprites() ([]*Sprite, error) {
	return ss.GetSprites(0, ss.Columns*ss.Rows)
}
