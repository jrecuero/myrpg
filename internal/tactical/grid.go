// Package tactical provides tile-based world representation for tactical RPG gameplay
// Simplified 2D grid system (no height complexity)
package tactical

import (
	"math"
)

// TileType represents different types of terrain
type TileType int

const (
	TileFloor TileType = iota
	TileWall
	TileWater
	TilePit
	TileElevated
)

// Tile represents a single grid cell in the tactical battlefield
type Tile struct {
	X        int      // Grid X coordinate
	Y        int      // Grid Y coordinate
	Type     TileType // Terrain type
	Passable bool     // Can units move through this tile
	Occupied bool     // Is there a unit on this tile
	UnitID   string   // ID of unit occupying this tile (if any)
}

// Grid represents the tactical battlefield
type Grid struct {
	Width    int               // Number of tiles horizontally
	Height   int               // Number of tiles vertically
	TileSize int               // Size of each tile in pixels
	Tiles    map[GridPos]*Tile // Map of grid positions to tiles
}

// GridPos represents a coordinate in the grid
type GridPos struct {
	X int
	Y int // Standard 2D Y coordinate for vertical axis
}

// NewGrid creates a new tactical grid
func NewGrid(width, height, tileSize int) *Grid {
	grid := &Grid{
		Width:    width,
		Height:   height,
		TileSize: tileSize,
		Tiles:    make(map[GridPos]*Tile),
	}

	// Initialize with basic floor tiles
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pos := GridPos{X: x, Y: y}
			grid.Tiles[pos] = &Tile{
				X:        x,
				Y:        y,
				Type:     TileFloor,
				Passable: true,
				Occupied: false,
			}
		}
	}

	return grid
}

// WorldToGrid converts world pixel coordinates to grid coordinates
func (g *Grid) WorldToGrid(worldX, worldY float64) GridPos {
	return GridPos{
		X: int(worldX) / g.TileSize,
		Y: int(worldY) / g.TileSize,
	}
}

// GridToWorld converts grid coordinates to world pixel coordinates
func (g *Grid) GridToWorld(pos GridPos) (float64, float64) {
	return float64(pos.X * g.TileSize), float64(pos.Y * g.TileSize)
}

// GetTile returns the tile at the given grid position
func (g *Grid) GetTile(pos GridPos) *Tile {
	return g.Tiles[pos]
}

// IsPassable checks if a unit can move to the given position
func (g *Grid) IsPassable(pos GridPos) bool {
	tile := g.GetTile(pos)
	return tile != nil && tile.Passable && !tile.Occupied
}

// SetOccupied marks a tile as occupied or free
func (g *Grid) SetOccupied(pos GridPos, occupied bool, unitID string) {
	if tile := g.GetTile(pos); tile != nil {
		tile.Occupied = occupied
		tile.UnitID = unitID
	}
}

// CalculateDistance calculates movement distance between two grid positions
func (g *Grid) CalculateDistance(from, to GridPos) int {
	dx := int(math.Abs(float64(to.X - from.X)))
	dy := int(math.Abs(float64(to.Y - from.Y)))
	return dx + dy // Manhattan distance for grid-based movement
}

// GetNeighbors returns adjacent tiles (4-directional movement)
func (g *Grid) GetNeighbors(pos GridPos) []GridPos {
	neighbors := []GridPos{
		{X: pos.X + 1, Y: pos.Y}, // Right
		{X: pos.X - 1, Y: pos.Y}, // Left
		{X: pos.X, Y: pos.Y + 1}, // Down
		{X: pos.X, Y: pos.Y - 1}, // Up
	}

	// Filter out invalid positions
	valid := make([]GridPos, 0)
	for _, neighbor := range neighbors {
		if g.IsValidPosition(neighbor) {
			valid = append(valid, neighbor)
		}
	}

	return valid
}

// IsValidPosition checks if a grid position is within the grid bounds
func (g *Grid) IsValidPosition(pos GridPos) bool {
	return pos.X >= 0 && pos.X < g.Width && pos.Y >= 0 && pos.Y < g.Height
}

// CalculateMovementRange returns all tiles within movement range of a position
func (g *Grid) CalculateMovementRange(from GridPos, moveRange int) []GridPos {
	reachable := make([]GridPos, 0)
	visited := make(map[GridPos]bool)

	// Use BFS to find all reachable tiles within range
	queue := []GridPos{from}
	distances := map[GridPos]int{from: 0}
	visited[from] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		currentDist := distances[current]
		if currentDist >= moveRange {
			continue
		}

		for _, neighbor := range g.GetNeighbors(current) {
			if visited[neighbor] || !g.IsPassable(neighbor) {
				continue
			}

			newDist := currentDist + 1
			if newDist <= moveRange {
				visited[neighbor] = true
				distances[neighbor] = newDist
				queue = append(queue, neighbor)
				reachable = append(reachable, neighbor)
			}
		}
	}

	return reachable
}
