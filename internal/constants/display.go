package constants

// Screen and Layout Constants
const (
	// Screen Dimensions
	ScreenWidth  = 1000
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

	// Background Constants
	BackgroundWidth  = 800             // Game world background width
	BackgroundHeight = GameWorldHeight // Same as game world height (408px)
	BackgroundX      = 0               // Background X position
	BackgroundY      = GameWorldY      // Background Y position (112px)
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

// Combat Panel Layout Constants
const (
	// Panel Separations - 5 pixels everywhere
	PanelMargin    = 5 // Margin from screen edges
	PanelSeparator = 5 // Separation between panels

	// Main Panel Specifications (adjusted for 5px margins/separators)
	EnemyPanelWidth  = 791 // 1000 - 5 - 5 - 5 - 194 = 791
	EnemyPanelHeight = 222 // Reduced slightly to fit height
	EnemyPanelX      = PanelMargin
	EnemyPanelY      = PanelMargin

	PlayerPanelWidth  = 791
	PlayerPanelHeight = 222
	PlayerPanelX      = PanelMargin
	PlayerPanelY      = EnemyPanelY + EnemyPanelHeight + PanelSeparator

	CombinedLogPanelWidth  = 194
	CombinedLogPanelHeight = 449 // Aligned to match Battle Log Panel Y position
	CombinedLogPanelX      = EnemyPanelX + EnemyPanelWidth + PanelSeparator
	CombinedLogPanelY      = PanelMargin

	ActionMenuPanelWidth  = 194
	ActionMenuPanelHeight = 136 // Same height as Battle Log Panel
	ActionMenuPanelX      = EnemyPanelX + EnemyPanelWidth + PanelSeparator
	ActionMenuPanelY      = CombinedLogPanelY + CombinedLogPanelHeight + PanelSeparator

	BattleLogPanelWidth  = 791
	BattleLogPanelHeight = 136 // Same height as Action Menu Panel
	BattleLogPanelX      = PanelMargin
	BattleLogPanelY      = PlayerPanelY + PlayerPanelHeight + PanelSeparator

	// Formation Layout within Enemy/Player Panels
	// Two rows: back row and front row
	FormationMarginLeft   = 2
	FormationMarginRight  = 2
	FormationMarginTop    = 2
	FormationMarginBottom = 11
	FormationRowSeparator = 2

	// Row dimensions
	FormationRowWidth  = 787 // 791 - 2 (left) - 2 (right)
	FormationRowHeight = 105

	// Back row positioning
	BackRowX = FormationMarginLeft
	BackRowY = FormationMarginTop

	// Front row positioning (below back row)
	FrontRowX = FormationMarginLeft
	FrontRowY = BackRowY + FormationRowHeight + FormationRowSeparator

	// Entity display within rows
	EntityDisplayMarginTop    = 5
	EntityDisplayMarginBottom = 5
	EntitySpriteBoxSize       = 95 // 95x95 box for sprite, health bar, and name
)
