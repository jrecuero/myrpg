package components

import (
	"image"
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

// Bounds returns the bounding rectangle of the Transform.
// It uses the X, Y, Width, and Height fields of the Transform to create the rectangle.
// returns an image.Rectangle representing the bounding box of the Transform.
func (t *Transform) Bounds() image.Rectangle {
	return image.Rect(int(t.X), int(t.Y), int(t.X)+t.Width, int(t.Y)+t.Height)
}
