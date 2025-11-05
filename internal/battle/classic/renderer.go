package classic

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/constants"
	"github.com/jrecuero/myrpg/internal/ecs"
)

// BattleRenderer handles rendering with new panel layout
type BattleRenderer struct {
	battleManager *BattleManager
	screenWidth   int
	screenHeight  int

	// Panel areas
	enemyPanelX, enemyPanelY, enemyPanelWidth, enemyPanelHeight                         int
	playerPanelX, playerPanelY, playerPanelWidth, playerPanelHeight                     int
	actionPanelX, actionPanelY, actionPanelWidth, actionPanelHeight                     int
	combinedLogPanelX, combinedLogPanelY, combinedLogPanelWidth, combinedLogPanelHeight int
	battleLogPanelX, battleLogPanelY, battleLogPanelWidth, battleLogPanelHeight         int

	// Colors
	backgroundColor, textColor, panelBorderColor, playerNameColor, enemyNameColor color.Color
	battleMessages                                                                []string
	maxLogMessages                                                                int
	showActivityQueue                                                             bool
	showBattleLog                                                                 bool
}

// NewBattleRenderer creates a new renderer with panel layout using constants
func NewBattleRenderer(battleManager *BattleManager, screenWidth, screenHeight int) *BattleRenderer {
	return &BattleRenderer{
		battleManager: battleManager,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,

		// Use exact panel dimensions from constants
		enemyPanelX:      constants.EnemyPanelX,
		enemyPanelY:      constants.EnemyPanelY,
		enemyPanelWidth:  constants.EnemyPanelWidth,
		enemyPanelHeight: constants.EnemyPanelHeight,

		playerPanelX:      constants.PlayerPanelX,
		playerPanelY:      constants.PlayerPanelY,
		playerPanelWidth:  constants.PlayerPanelWidth,
		playerPanelHeight: constants.PlayerPanelHeight,

		actionPanelX:      constants.ActionMenuPanelX,
		actionPanelY:      constants.ActionMenuPanelY,
		actionPanelWidth:  constants.ActionMenuPanelWidth,
		actionPanelHeight: constants.ActionMenuPanelHeight,

		combinedLogPanelX:      constants.CombinedLogPanelX,
		combinedLogPanelY:      constants.CombinedLogPanelY,
		combinedLogPanelWidth:  constants.CombinedLogPanelWidth,
		combinedLogPanelHeight: constants.CombinedLogPanelHeight,

		battleLogPanelX:      constants.BattleLogPanelX,
		battleLogPanelY:      constants.BattleLogPanelY,
		battleLogPanelWidth:  constants.BattleLogPanelWidth,
		battleLogPanelHeight: constants.BattleLogPanelHeight,

		backgroundColor:   color.RGBA{20, 20, 40, 255},
		textColor:         color.RGBA{255, 255, 255, 255},
		panelBorderColor:  color.RGBA{100, 100, 100, 255},
		playerNameColor:   color.RGBA{100, 255, 100, 255},
		enemyNameColor:    color.RGBA{255, 100, 100, 255},
		battleMessages:    make([]string, 0),
		maxLogMessages:    5,
		showActivityQueue: true,
		showBattleLog:     true,
	}
}

// Render draws the entire battle scene with panel layout
func (br *BattleRenderer) Render(screen *ebiten.Image) {
	if !br.battleManager.IsActive() {
		return
	}

	// Draw background
	screen.Fill(br.backgroundColor)

	// Draw all panels
	br.drawEnemyPanel(screen)
	br.drawPlayerPanel(screen)
	br.drawActionPanel(screen)
	br.drawCombinedLogPanel(screen)
	br.drawBattleLogPanel(screen)
}

// Helper to draw panel with border
func (br *BattleRenderer) drawPanel(screen *ebiten.Image, x, y, width, height int, bgColor color.RGBA, title string) {
	// Background
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(width), float64(height), bgColor)

	// Border
	border := 2.0
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(width), border, br.panelBorderColor)
	ebitenutil.DrawRect(screen, float64(x), float64(y+height-2), float64(width), border, br.panelBorderColor)
	ebitenutil.DrawRect(screen, float64(x), float64(y), border, float64(height), br.panelBorderColor)
	ebitenutil.DrawRect(screen, float64(x+width-2), float64(y), border, float64(height), br.panelBorderColor)

	// Title
	if title != "" {
		ebitenutil.DebugPrintAt(screen, title, x+5, y+5)
	}
}

// drawEnemyPanel draws enemies in two-row formation
func (br *BattleRenderer) drawEnemyPanel(screen *ebiten.Image) {
	br.drawPanel(screen, br.enemyPanelX, br.enemyPanelY, br.enemyPanelWidth, br.enemyPanelHeight,
		color.RGBA{40, 20, 20, 100}, "Enemies")

	// Draw enemy formations in back and front rows
	if formation := br.battleManager.GetEnemyFormation(); formation != nil {
		br.drawFormationRows(screen, formation, br.enemyPanelX, br.enemyPanelY, true)
	}
}

// drawPlayerPanel draws players in two-row formation
func (br *BattleRenderer) drawPlayerPanel(screen *ebiten.Image) {
	br.drawPanel(screen, br.playerPanelX, br.playerPanelY, br.playerPanelWidth, br.playerPanelHeight,
		color.RGBA{20, 40, 20, 100}, "Players")

	// Draw player formations in back and front rows
	if formation := br.battleManager.GetPlayerFormation(); formation != nil {
		br.drawFormationRows(screen, formation, br.playerPanelX, br.playerPanelY, false)
	}
}

// drawActionPanel draws action menu in bottom-left
func (br *BattleRenderer) drawActionPanel(screen *ebiten.Image) {
	br.drawPanel(screen, br.actionPanelX, br.actionPanelY, br.actionPanelWidth, br.actionPanelHeight,
		color.RGBA{40, 40, 60, 200}, "Actions")

	state := br.battleManager.GetBattleState()
	baseY := br.actionPanelY + 25

	switch state {
	case BattleStateWaitingForPlayerAction:
		if player := br.battleManager.GetCurrentPlayerEntity(); player != nil {
			name := "Player"
			if stats := player.RPGStats(); stats != nil {
				name = stats.Name
			}
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s's Turn - Choose Action:", name), br.actionPanelX+10, baseY)
			ebitenutil.DebugPrintAt(screen, "1. Physical Attack", br.actionPanelX+10, baseY+25)
			ebitenutil.DebugPrintAt(screen, "2. Magical Attack", br.actionPanelX+10, baseY+45)
			ebitenutil.DebugPrintAt(screen, "3. Defend", br.actionPanelX+10, baseY+65)
		}
	case BattleStateWaitingForTarget:
		ebitenutil.DebugPrintAt(screen, "Select Target:", br.actionPanelX+10, baseY)
		ebitenutil.DebugPrintAt(screen, "Use arrows, Enter to confirm", br.actionPanelX+10, baseY+20)
	case BattleStateVictory:
		ebitenutil.DebugPrintAt(screen, "VICTORY!", br.actionPanelX+10, baseY)
		ebitenutil.DebugPrintAt(screen, "Press any key to continue...", br.actionPanelX+10, baseY+20)
	case BattleStateDefeat:
		ebitenutil.DebugPrintAt(screen, "DEFEAT!", br.actionPanelX+10, baseY)
		ebitenutil.DebugPrintAt(screen, "Press any key to continue...", br.actionPanelX+10, baseY+20)
	}
}

// drawCombinedLogPanel draws activity queue and player info in right-top
func (br *BattleRenderer) drawCombinedLogPanel(screen *ebiten.Image) {
	br.drawPanel(screen, br.combinedLogPanelX, br.combinedLogPanelY, br.combinedLogPanelWidth, br.combinedLogPanelHeight,
		color.RGBA{30, 30, 50, 200}, "Activity Queue")

	// Draw activity queue
	queue := br.battleManager.GetActivityQueue()
	baseY := br.combinedLogPanelY + 25
	y := baseY

	for i, entry := range queue {
		if i >= 10 {
			break
		}

		name := "Unknown"
		prefix := "P:"

		if stats := entry.Entity.RPGStats(); stats != nil {
			name = stats.Name

			// Check if enemy
			if formation := br.battleManager.GetEnemyFormation(); formation != nil {
				positions := formation.GetAllPositions()
				for enemyEntity := range positions {
					if enemyEntity == entry.Entity {
						prefix = "E:"
						break
					}
				}
			}
		}

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s %s", prefix, name), br.combinedLogPanelX+10, y)
		y += 18
	}

	// Draw current player info below
	if br.battleManager.GetBattleState() == BattleStateWaitingForPlayerAction {
		if player := br.battleManager.GetCurrentPlayerEntity(); player != nil {
			if stats := player.RPGStats(); stats != nil {
				startY := br.combinedLogPanelY + 250
				ebitenutil.DebugPrintAt(screen, "Current Player:", br.combinedLogPanelX+10, startY)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Name: %s", stats.Name), br.combinedLogPanelX+10, startY+20)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d/%d", stats.CurrentHP, stats.MaxHP), br.combinedLogPanelX+10, startY+35)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MP: %d/%d", stats.CurrentMP, stats.MaxMP), br.combinedLogPanelX+10, startY+50)
			}
		}
	}
}

// drawBattleLogPanel draws combat messages in bottom-right
func (br *BattleRenderer) drawBattleLogPanel(screen *ebiten.Image) {
	br.drawPanel(screen, br.battleLogPanelX, br.battleLogPanelY, br.battleLogPanelWidth, br.battleLogPanelHeight,
		color.RGBA{50, 30, 30, 200}, "Battle Log")

	baseY := br.battleLogPanelY + 25
	y := baseY

	start := 0
	if len(br.battleMessages) > br.maxLogMessages {
		start = len(br.battleMessages) - br.maxLogMessages
	}

	for i := start; i < len(br.battleMessages); i++ {
		ebitenutil.DebugPrintAt(screen, br.battleMessages[i], br.battleLogPanelX+10, y)
		y += 18
	}
}

// drawEntity draws individual entities with name/health above sprite
func (br *BattleRenderer) drawEntity(screen *ebiten.Image, entity *ecs.Entity, pos Position, isEnemy bool) {
	stats := entity.RPGStats()
	if stats == nil {
		return
	}

	// Check if selected target
	isSelected := br.isSelectedTarget(entity)

	// Draw highlight
	if isSelected {
		ebitenutil.DrawRect(screen, pos.X-3, pos.Y-3, 54, 54, color.RGBA{255, 255, 0, 150})
	}

	// Draw sprite or rectangle
	if sprite := entity.Sprite(); sprite != nil && sprite.Sprite != nil && sprite.Sprite.Img != nil {
		op := &ebiten.DrawImageOptions{}
		if stats.CurrentHP <= 0 {
			op.ColorScale.Scale(0.5, 0.5, 0.5, 0.7)
		}
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(sprite.Sprite.Img, op)
	} else {
		entityColor := color.RGBA{100, 255, 100, 255}
		if isEnemy {
			entityColor = color.RGBA{255, 100, 100, 255}
		}
		if stats.CurrentHP <= 0 {
			entityColor = color.RGBA{100, 100, 100, 255}
		}
		ebitenutil.DrawRect(screen, pos.X, pos.Y, 48, 48, entityColor)
	}

	// Draw name and health above sprite
	name := stats.Name
	if len(name) > 8 {
		name = name[:8]
	}
	ebitenutil.DebugPrintAt(screen, name, int(pos.X), int(pos.Y-25))
	br.drawHealthBar(screen, pos.X, pos.Y-10, 48, 6, stats.CurrentHP, stats.MaxHP)
}

// isSelectedTarget checks if entity is the selected target
func (br *BattleRenderer) isSelectedTarget(entity *ecs.Entity) bool {
	if br.battleManager.GetBattleState() != BattleStateWaitingForTarget {
		return false
	}

	targets := br.battleManager.GetAvailableTargets()
	index := br.battleManager.GetTargetIndex()

	if index >= 0 && index < len(targets) {
		return targets[index] == entity
	}
	return false
}

// drawHealthBar renders health bar
func (br *BattleRenderer) drawHealthBar(screen *ebiten.Image, x, y, width, height float64, currentHP, maxHP int) {
	if maxHP <= 0 {
		return
	}

	// Background
	ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{100, 0, 0, 255})

	// Health
	percentage := float64(currentHP) / float64(maxHP)
	healthWidth := width * percentage

	var healthColor color.RGBA
	if percentage > 0.6 {
		healthColor = color.RGBA{0, 255, 0, 255}
	} else if percentage > 0.3 {
		healthColor = color.RGBA{255, 255, 0, 255}
	} else {
		healthColor = color.RGBA{255, 0, 0, 255}
	}

	if healthWidth > 0 {
		ebitenutil.DrawRect(screen, x, y, healthWidth, height, healthColor)
	}
}

// AddBattleMessage adds a message to the battle log
func (br *BattleRenderer) AddBattleMessage(message string) {
	br.battleMessages = append(br.battleMessages, message)
}

// SetShowActivityQueue toggles activity queue display
func (br *BattleRenderer) SetShowActivityQueue(show bool) {
	br.showActivityQueue = show
}

// SetShowBattleLog toggles battle log display
func (br *BattleRenderer) SetShowBattleLog(show bool) {
	br.showBattleLog = show
}

// drawFormationRows draws entities in back and front row formation layout
func (br *BattleRenderer) drawFormationRows(screen *ebiten.Image, formation *Formation, panelX, panelY int, isEnemy bool) {
	// Get back and front row entities
	backRowEntities := formation.GetBackRowEntities()
	frontRowEntities := formation.GetFrontRowEntities()

	// Draw back row
	br.drawFormationRow(screen, backRowEntities, panelX+constants.BackRowX, panelY+constants.BackRowY, isEnemy, "Back Row")

	// Draw front row
	br.drawFormationRow(screen, frontRowEntities, panelX+constants.FrontRowX, panelY+constants.FrontRowY, isEnemy, "Front Row")
}

// drawFormationRow draws a single row of entities
func (br *BattleRenderer) drawFormationRow(screen *ebiten.Image, entities []*ecs.Entity, rowX, rowY int, isEnemy bool, rowLabel string) {
	if len(entities) == 0 {
		return
	}

	// Calculate how many entities can fit horizontally with proper spacing
	maxEntitiesPerRow := constants.FormationRowWidth / constants.EntitySpriteBoxSize
	entitySpacing := constants.FormationRowWidth / maxEntitiesPerRow

	// Draw entities in the row, centered vertically in the 95px available space
	for i, entity := range entities {
		if i >= maxEntitiesPerRow {
			break // Don't overflow the row
		}

		// Calculate position for this entity within the 95x95 box
		entityX := rowX + (i * entitySpacing) + (entitySpacing-constants.EntitySpriteBoxSize)/2
		entityY := rowY + constants.EntityDisplayMarginTop

		br.drawEntityInBox(screen, entity, entityX, entityY, isEnemy)
	}
}

// drawEntityInBox draws an entity within a 95x95 box with sprite, health bar, and name
func (br *BattleRenderer) drawEntityInBox(screen *ebiten.Image, entity *ecs.Entity, x, y int, isEnemy bool) {
	stats := entity.RPGStats()
	if stats == nil {
		return
	}

	// Check if selected target
	isSelected := br.isSelectedTarget(entity)

	// Draw selection highlight around the entire box
	if isSelected {
		// Highlight the entire 95x95 box
		vector.FillRect(screen, float32(x-2), float32(y-2), float32(constants.EntitySpriteBoxSize+4), float32(constants.EntitySpriteBoxSize+4), color.RGBA{255, 255, 0, 150}, false)
	}

	// Calculate sprite position (centered horizontally, with room for health bar and name)
	spriteSize := 32 // Standard sprite size
	spriteX := x + (constants.EntitySpriteBoxSize-spriteSize)/2
	spriteY := y + 20 // Leave room for name at top

	// Draw sprite or rectangle
	if sprite := entity.Sprite(); sprite != nil && sprite.Sprite != nil && sprite.Sprite.Img != nil {
		op := &ebiten.DrawImageOptions{}
		if stats.CurrentHP <= 0 {
			op.ColorScale.Scale(0.5, 0.5, 0.5, 0.7)
		}
		op.GeoM.Translate(float64(spriteX), float64(spriteY))
		screen.DrawImage(sprite.Sprite.Img, op)
	} else {
		entityColor := color.RGBA{100, 255, 100, 255}
		if isEnemy {
			entityColor = color.RGBA{255, 100, 100, 255}
		}
		if stats.CurrentHP <= 0 {
			entityColor = color.RGBA{100, 100, 100, 255}
		}
		vector.FillRect(screen, float32(spriteX), float32(spriteY), float32(spriteSize), float32(spriteSize), entityColor, false)
	}

	// Draw name above sprite (truncated to fit)
	name := stats.Name
	if len(name) > 8 {
		name = name[:8]
	}
	nameX := x + (constants.EntitySpriteBoxSize-len(name)*6)/2 // Center the text
	ebitenutil.DebugPrintAt(screen, name, nameX, y+5)

	// Draw health bar below sprite
	healthBarY := spriteY + spriteSize + 5
	healthBarWidth := constants.EntitySpriteBoxSize - 10 // Leave 5px margin on each side
	br.drawHealthBar(screen, float64(x+5), float64(healthBarY), float64(healthBarWidth), 6, stats.CurrentHP, stats.MaxHP)
}
