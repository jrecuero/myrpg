// Package classic implements rendering for Dragon Quest-style battles
package classic

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jrecuero/myrpg/internal/ecs"
)

// BattleRenderer handles rendering of the classic battle system
type BattleRenderer struct {
	battleManager *BattleManager

	// Screen dimensions
	screenWidth  int
	screenHeight int

	// UI areas
	battleAreaY      int
	battleAreaHeight int
	uiAreaY          int
	uiAreaHeight     int

	// Colors
	backgroundColor color.Color
	playerAreaColor color.Color
	enemyAreaColor  color.Color
	textColor       color.Color
	highlightColor  color.Color

	// UI state
	showActivityQueue bool
	showBattleLog     bool
	battleMessages    []string
	maxLogMessages    int
}

// NewBattleRenderer creates a new renderer for classic battles
func NewBattleRenderer(battleManager *BattleManager, screenWidth, screenHeight int) *BattleRenderer {
	return &BattleRenderer{
		battleManager:     battleManager,
		screenWidth:       screenWidth,
		screenHeight:      screenHeight,
		battleAreaY:       50,
		battleAreaHeight:  400,
		uiAreaY:           450,
		uiAreaHeight:      150,
		backgroundColor:   color.RGBA{20, 20, 40, 255},    // Dark blue
		playerAreaColor:   color.RGBA{20, 40, 20, 100},    // Dark green (transparent)
		enemyAreaColor:    color.RGBA{40, 20, 20, 100},    // Dark red (transparent)
		textColor:         color.RGBA{255, 255, 255, 255}, // White
		highlightColor:    color.RGBA{255, 255, 0, 255},   // Yellow
		showActivityQueue: true,
		showBattleLog:     true,
		battleMessages:    make([]string, 0),
		maxLogMessages:    5,
	}
}

// Render draws the entire battle scene
func (br *BattleRenderer) Render(screen *ebiten.Image) {
	if !br.battleManager.IsActive() {
		return
	}

	// Draw background
	br.drawBackground(screen)

	// Draw battle area
	br.drawBattleArea(screen)

	// Draw formations
	br.drawFormations(screen)

	// Draw UI
	br.drawUI(screen)

	// Draw battle messages
	if br.showBattleLog {
		br.drawBattleLog(screen)
	}

	// Draw activity queue
	if br.showActivityQueue {
		br.drawActivityQueue(screen)
	}
}

// drawBackground renders the battle background
func (br *BattleRenderer) drawBackground(screen *ebiten.Image) {
	screen.Fill(br.backgroundColor)

	// Draw battle area separator
	ebitenutil.DrawLine(screen,
		0, float64(br.battleAreaY+br.battleAreaHeight),
		float64(br.screenWidth), float64(br.battleAreaY+br.battleAreaHeight),
		color.RGBA{100, 100, 100, 255})
}

// drawBattleArea renders the main battle visualization area
func (br *BattleRenderer) drawBattleArea(screen *ebiten.Image) {
	// Draw player area (bottom section)
	playerAreaY := br.battleAreaY + br.battleAreaHeight - 150
	ebitenutil.DrawRect(screen,
		0, float64(playerAreaY),
		float64(br.screenWidth), 150,
		br.playerAreaColor)

	// Draw enemy area (top section)
	enemyAreaHeight := float64(br.battleAreaHeight - 150)
	ebitenutil.DrawRect(screen,
		0, float64(br.battleAreaY),
		float64(br.screenWidth), enemyAreaHeight,
		br.enemyAreaColor)
}

// drawFormations renders both player and enemy formations
func (br *BattleRenderer) drawFormations(screen *ebiten.Image) {
	// Draw enemy formation
	if enemyFormation := br.battleManager.GetEnemyFormation(); enemyFormation != nil {
		br.drawFormation(screen, enemyFormation, true)
	}

	// Draw player formation
	if playerFormation := br.battleManager.GetPlayerFormation(); playerFormation != nil {
		br.drawFormation(screen, playerFormation, false)
	}
}

// drawFormation renders a specific formation
func (br *BattleRenderer) drawFormation(screen *ebiten.Image, formation *Formation, isEnemy bool) {
	positions := formation.GetAllPositions()

	for entity, pos := range positions {
		br.drawEntity(screen, entity, pos, isEnemy)
	}
}

// drawEntity renders a single entity in the formation
func (br *BattleRenderer) drawEntity(screen *ebiten.Image, entity *ecs.Entity, pos Position, isEnemy bool) {
	// Get entity stats to check if alive
	stats := entity.RPGStats()
	if stats == nil {
		return
	}

	// Determine entity color based on status
	entityColor := br.textColor
	if stats.CurrentHP <= 0 {
		entityColor = color.RGBA{100, 100, 100, 255} // Gray for defeated
	} else if isEnemy {
		entityColor = color.RGBA{255, 100, 100, 255} // Red for enemies
	} else {
		entityColor = color.RGBA{100, 255, 100, 255} // Green for players
	}

	// Try to render sprite first
	if sprite := entity.Sprite(); sprite != nil && sprite.Sprite != nil && sprite.Sprite.Img != nil {
		op := &ebiten.DrawImageOptions{}

		// Apply color tint based on status
		if stats.CurrentHP <= 0 {
			op.ColorScale.Scale(0.5, 0.5, 0.5, 0.7) // Darken defeated entities
		}

		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(sprite.Sprite.Img, op)
	} else {
		// Fallback to colored rectangle
		width, height := 48.0, 48.0
		ebitenutil.DrawRect(screen, pos.X, pos.Y, width, height, entityColor)
	}

	// Draw entity name and HP
	br.drawEntityInfo(screen, entity, pos, isEnemy)
}

// drawEntityInfo renders entity name and health information
func (br *BattleRenderer) drawEntityInfo(screen *ebiten.Image, entity *ecs.Entity, pos Position, isEnemy bool) {
	stats := entity.RPGStats()
	if stats == nil {
		return
	}

	// Entity name
	nameY := pos.Y - 15
	if isEnemy {
		nameY = pos.Y + 60 // Below enemy sprites
	}

	// Simple text rendering (you might want to use a proper font)
	name := stats.Name
	if len(name) > 8 {
		name = name[:8] // Truncate long names
	}

	ebitenutil.DebugPrintAt(screen, name, int(pos.X), int(nameY))

	// HP bar
	hpBarY := nameY + 15
	br.drawHealthBar(screen, pos.X, hpBarY, 48, 6, stats.CurrentHP, stats.MaxHP)
}

// drawHealthBar renders a health bar
func (br *BattleRenderer) drawHealthBar(screen *ebiten.Image, x, y, width, height float64, currentHP, maxHP int) {
	// Background (empty health)
	ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{100, 0, 0, 255})

	// Foreground (current health)
	if maxHP > 0 {
		healthPercent := float64(currentHP) / float64(maxHP)
		healthWidth := width * healthPercent

		healthColor := color.RGBA{0, 255, 0, 255} // Green
		if healthPercent < 0.3 {
			healthColor = color.RGBA{255, 0, 0, 255} // Red
		} else if healthPercent < 0.6 {
			healthColor = color.RGBA{255, 255, 0, 255} // Yellow
		}

		ebitenutil.DrawRect(screen, x, y, healthWidth, height, healthColor)
	}

	// Border
	ebitenutil.DrawRect(screen, x-1, y-1, width+2, 1, br.textColor)      // Top
	ebitenutil.DrawRect(screen, x-1, y+height, width+2, 1, br.textColor) // Bottom
	ebitenutil.DrawRect(screen, x-1, y, 1, height, br.textColor)         // Left
	ebitenutil.DrawRect(screen, x+width, y, 1, height, br.textColor)     // Right
}

// drawUI renders the bottom UI area
func (br *BattleRenderer) drawUI(screen *ebiten.Image) {
	// UI background
	ebitenutil.DrawRect(screen,
		0, float64(br.uiAreaY),
		float64(br.screenWidth), float64(br.uiAreaHeight),
		color.RGBA{40, 40, 60, 255})

	// Battle state info
	state := br.battleManager.GetState()
	stateText := br.getBattleStateText(state)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Battle State: %s", stateText), 10, br.uiAreaY+10)

	// Instructions
	ebitenutil.DebugPrintAt(screen, "Classic Battle System - Dragon Quest Style", 10, br.uiAreaY+30)
	ebitenutil.DebugPrintAt(screen, "Speed-based turns, formations, activity queue", 10, br.uiAreaY+45)
}

// drawActivityQueue renders the activity queue
func (br *BattleRenderer) drawActivityQueue(screen *ebiten.Image) {
	queue := br.battleManager.GetActivityQueue()
	if len(queue) == 0 {
		return
	}

	// Draw activity queue title
	queueX := br.screenWidth - 200
	queueY := 60
	ebitenutil.DebugPrintAt(screen, "Activity Queue:", queueX, queueY)

	// Draw top 5 entities in queue
	maxDisplay := 5
	if len(queue) < maxDisplay {
		maxDisplay = len(queue)
	}

	for i := 0; i < maxDisplay; i++ {
		entry := queue[i]
		y := queueY + 20 + (i * 15)

		name := "Unknown"
		if stats := entry.Entity.RPGStats(); stats != nil {
			name = stats.Name
		}

		timeUntilAction := entry.NextActionAt.Sub(br.battleManager.battleTime)

		text := fmt.Sprintf("%d. %s (%.1fs)", i+1, name, timeUntilAction.Seconds())
		ebitenutil.DebugPrintAt(screen, text, queueX, y)
	}
}

// drawBattleLog renders recent battle messages
func (br *BattleRenderer) drawBattleLog(screen *ebiten.Image) {
	if len(br.battleMessages) == 0 {
		return
	}

	logX := 10
	logY := br.uiAreaY + 70

	ebitenutil.DebugPrintAt(screen, "Battle Log:", logX, logY)

	for i, message := range br.battleMessages {
		y := logY + 15 + (i * 12)
		ebitenutil.DebugPrintAt(screen, message, logX, y)
	}
}

// getBattleStateText converts battle state to readable text
func (br *BattleRenderer) getBattleStateText(state BattleState) string {
	switch state {
	case BattleStateIdle:
		return "Idle"
	case BattleStatePlayerTurn:
		return "Player Turn"
	case BattleStateEnemyTurn:
		return "Enemy Turn"
	case BattleStateAnimation:
		return "Animation"
	case BattleStateVictory:
		return "Victory!"
	case BattleStateDefeat:
		return "Defeat"
	case BattleStateEscaped:
		return "Escaped"
	default:
		return "Unknown"
	}
}

// AddBattleMessage adds a message to the battle log
func (br *BattleRenderer) AddBattleMessage(message string) {
	br.battleMessages = append(br.battleMessages, message)

	// Keep only the most recent messages
	if len(br.battleMessages) > br.maxLogMessages {
		br.battleMessages = br.battleMessages[1:]
	}
}

// SetShowActivityQueue toggles activity queue display
func (br *BattleRenderer) SetShowActivityQueue(show bool) {
	br.showActivityQueue = show
}

// SetShowBattleLog toggles battle log display
func (br *BattleRenderer) SetShowBattleLog(show bool) {
	br.showBattleLog = show
}
