package classic

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

// NewBattleRenderer creates a new renderer with panel layout
func NewBattleRenderer(battleManager *BattleManager, screenWidth, screenHeight int) *BattleRenderer {
	// Calculate panel dimensions with adjusted widths
	// Make left side wider (65%) for enemy/player panels, right side narrower (35%) for logs
	leftWidth := int(float64(screenWidth) * 0.65)
	rightWidth := screenWidth - leftWidth
	topHeight := screenHeight / 3
	middleHeight := screenHeight / 3
	bottomHeight := screenHeight - topHeight - middleHeight

	// Action panel takes only 40% of left width, allowing battle log to expand
	actionPanelWidth := int(float64(leftWidth) * 0.4)
	battleLogStartX := 10 + actionPanelWidth + 10        // Start after action panel with gap
	battleLogWidth := screenWidth - battleLogStartX - 10 // Extend to screen edge

	return &BattleRenderer{
		battleManager: battleManager,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,

		// Left column panels - now wider (65% of screen)
		enemyPanelX: 10, enemyPanelY: 10,
		enemyPanelWidth: leftWidth - 20, enemyPanelHeight: topHeight - 10,

		playerPanelX: 10, playerPanelY: topHeight + 5,
		playerPanelWidth: leftWidth - 20, playerPanelHeight: middleHeight - 10,

		// Action panel - narrower (40% of left width)
		actionPanelX: 10, actionPanelY: topHeight + middleHeight + 5,
		actionPanelWidth: actionPanelWidth - 20, actionPanelHeight: bottomHeight - 15,

		// Right column panels - narrower (35% of screen)
		combinedLogPanelX: leftWidth + 5, combinedLogPanelY: 10,
		combinedLogPanelWidth: rightWidth - 15, combinedLogPanelHeight: topHeight + middleHeight - 10,

		// Battle log panel - expanded to use remaining space
		battleLogPanelX: battleLogStartX, battleLogPanelY: topHeight + middleHeight + 5,
		battleLogPanelWidth: battleLogWidth, battleLogPanelHeight: bottomHeight - 15,

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

// drawEnemyPanel draws enemies in top-left
func (br *BattleRenderer) drawEnemyPanel(screen *ebiten.Image) {
	br.drawPanel(screen, br.enemyPanelX, br.enemyPanelY, br.enemyPanelWidth, br.enemyPanelHeight,
		color.RGBA{40, 20, 20, 100}, "Enemies")

	// Draw enemy formations
	if formation := br.battleManager.GetEnemyFormation(); formation != nil {
		positions := formation.GetAllPositions()
		for entity, pos := range positions {
			br.drawEntity(screen, entity, pos, true)
		}
	}
}

// drawPlayerPanel draws players in middle-left
func (br *BattleRenderer) drawPlayerPanel(screen *ebiten.Image) {
	br.drawPanel(screen, br.playerPanelX, br.playerPanelY, br.playerPanelWidth, br.playerPanelHeight,
		color.RGBA{20, 40, 20, 100}, "Players")

	// Draw player formations
	if formation := br.battleManager.GetPlayerFormation(); formation != nil {
		positions := formation.GetAllPositions()
		for entity, pos := range positions {
			br.drawEntity(screen, entity, pos, false)
		}
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
