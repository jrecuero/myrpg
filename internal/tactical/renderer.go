// Package tactical provides grid rendering for tactical combat mode
package tactical

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TileHighlight represents different ways to highlight tiles
type TileHighlight int

const (
	HighlightNone     TileHighlight = iota
	HighlightMovement               // Blue - valid movement tiles
	HighlightAttack                 // Red - valid attack tiles
	HighlightSelected               // Yellow - currently selected tile
	HighlightPath                   // Green - movement path preview
)

// GridRenderer handles visual representation of the tactical grid
type GridRenderer struct {
	Grid             *Grid
	TileSize         int
	GridColor        color.Color
	HighlightColors  map[TileHighlight]color.Color
	ShowGrid         bool
	HighlightedTiles map[GridPos]TileHighlight
}

// NewGridRenderer creates a new grid renderer
func NewGridRenderer(grid *Grid) *GridRenderer {
	return &GridRenderer{
		Grid:      grid,
		TileSize:  grid.TileSize,
		GridColor: color.RGBA{R: 128, G: 128, B: 128, A: 128}, // Semi-transparent gray
		HighlightColors: map[TileHighlight]color.Color{
			HighlightMovement: color.RGBA{R: 100, G: 150, B: 255, A: 100}, // Light blue
			HighlightAttack:   color.RGBA{R: 255, G: 100, B: 100, A: 100}, // Light red
			HighlightSelected: color.RGBA{R: 255, G: 255, B: 100, A: 150}, // Yellow
			HighlightPath:     color.RGBA{R: 100, G: 255, B: 100, A: 100}, // Light green
		},
		ShowGrid:         true,
		HighlightedTiles: make(map[GridPos]TileHighlight),
	}
}

// SetShowGrid toggles grid line visibility
func (gr *GridRenderer) SetShowGrid(show bool) {
	gr.ShowGrid = show
}

// HighlightTile highlights a specific tile with the given highlight type
func (gr *GridRenderer) HighlightTile(pos GridPos, highlight TileHighlight) {
	if highlight == HighlightNone {
		delete(gr.HighlightedTiles, pos)
	} else {
		gr.HighlightedTiles[pos] = highlight
	}
}

// HighlightTiles highlights multiple tiles with the same highlight type
func (gr *GridRenderer) HighlightTiles(positions []GridPos, highlight TileHighlight) {
	for _, pos := range positions {
		gr.HighlightTile(pos, highlight)
	}
}

// ClearHighlights removes all tile highlights
func (gr *GridRenderer) ClearHighlights() {
	gr.HighlightedTiles = make(map[GridPos]TileHighlight)
}

// ClearHighlightType removes highlights of a specific type
func (gr *GridRenderer) ClearHighlightType(highlight TileHighlight) {
	for pos, highlightType := range gr.HighlightedTiles {
		if highlightType == highlight {
			delete(gr.HighlightedTiles, pos)
		}
	}
}

// Draw renders the grid and highlights to the screen
func (gr *GridRenderer) Draw(screen *ebiten.Image, offsetX, offsetY float64) {
	// Draw tile highlights first (behind grid lines)
	gr.drawHighlights(screen, offsetX, offsetY)

	// Draw grid lines on top
	if gr.ShowGrid {
		gr.drawGridLines(screen, offsetX, offsetY)
	}
}

// drawHighlights renders tile highlighting
func (gr *GridRenderer) drawHighlights(screen *ebiten.Image, offsetX, offsetY float64) {
	for pos, highlight := range gr.HighlightedTiles {
		if highlightColor, exists := gr.HighlightColors[highlight]; exists {
			// Calculate tile position on screen
			tileX := float32(pos.X*gr.TileSize) + float32(offsetX)
			tileY := float32(pos.Y*gr.TileSize) + float32(offsetY) // Draw filled rectangle for tile highlight
			vector.DrawFilledRect(screen,
				tileX, tileY,
				float32(gr.TileSize), float32(gr.TileSize),
				highlightColor, false)
		}
	}
}

// drawGridLines renders the grid overlay
func (gr *GridRenderer) drawGridLines(screen *ebiten.Image, offsetX, offsetY float64) {
	// Draw vertical lines
	for x := 0; x <= gr.Grid.Width; x++ {
		lineX := float32(x*gr.TileSize) + float32(offsetX)
		vector.StrokeLine(screen,
			lineX, float32(offsetY),
			lineX, float32(offsetY)+float32(gr.Grid.Height*gr.TileSize),
			1, gr.GridColor, false)
	}

	// Draw horizontal lines
	for y := 0; y <= gr.Grid.Height; y++ {
		lineY := float32(y*gr.TileSize) + float32(offsetY)
		vector.StrokeLine(screen,
			float32(offsetX), lineY,
			float32(offsetX)+float32(gr.Grid.Width*gr.TileSize), lineY,
			1, gr.GridColor, false)
	}
}

// GetTileAtScreenPos returns the grid position for a screen coordinate
func (gr *GridRenderer) GetTileAtScreenPos(screenX, screenY, offsetX, offsetY float64) (GridPos, bool) {
	// Convert screen coordinates to grid coordinates
	gridX := int((screenX - offsetX) / float64(gr.TileSize))
	gridY := int((screenY - offsetY) / float64(gr.TileSize))

	pos := GridPos{X: gridX, Y: gridY}

	// Check if position is within grid bounds
	if gr.Grid.IsValidPosition(pos) {
		return pos, true
	}

	return GridPos{}, false
}

// GetTileScreenPos returns the screen coordinates for a grid position
func (gr *GridRenderer) GetTileScreenPos(pos GridPos, offsetX, offsetY float64) (float64, float64) {
	screenX := float64(pos.X*gr.TileSize) + offsetX + float64(gr.TileSize)/2
	screenY := float64(pos.Y*gr.TileSize) + offsetY + float64(gr.TileSize)/2
	return screenX, screenY
}

// IsGridVisible returns true if the grid should be displayed
func (gr *GridRenderer) IsGridVisible() bool {
	return gr.ShowGrid
}

// SetGridColor changes the color of grid lines
func (gr *GridRenderer) SetGridColor(c color.Color) {
	gr.GridColor = c
}

// SetHighlightColor changes the color for a specific highlight type
func (gr *GridRenderer) SetHighlightColor(highlight TileHighlight, c color.Color) {
	gr.HighlightColors[highlight] = c
}
