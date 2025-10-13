package constants

// Screen and Layout Constants
const (
	// Screen Dimensions
	ScreenWidth  = 800
	ScreenHeight = 600

	// UI Panel Heights
	TopPanelHeight    = 110
	BottomPanelHeight = 80
	SeparatorHeight   = 2

	// Game World Area
	GameWorldY      = TopPanelHeight + SeparatorHeight
	GameWorldHeight = ScreenHeight - TopPanelHeight - BottomPanelHeight - SeparatorHeight
	GameWorldWidth  = ScreenWidth
	GameWorldLeft   = 0.0
	GameWorldRight  = float64(ScreenWidth)
	GameWorldTop    = float64(GameWorldY)
	GameWorldBottom = float64(GameWorldY + GameWorldHeight)
)

// Grid System Constants
const (
	// Tactical Grid Dimensions
	GridWidth  = 20 // Number of columns
	GridHeight = 10 // Number of rows
	TileSize   = 32 // Size of each grid tile in pixels

	// Grid Positioning
	GridOffsetX = 50.0
	GridOffsetY = float64(GameWorldY)
)
