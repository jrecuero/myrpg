// Package ui provides character statistics display widget for the RPG game.
// This widget shows detailed character information organized into categories
// with tabbed navigation and progress bars.
package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// Character Stats Widget Constants
const (
	// Layout dimensions
	StatsWidgetWidth          = 600 // Widget width in pixels
	StatsWidgetHeight         = 450 // Widget height in pixels
	StatsWidgetBorderWidth    = 2   // Border thickness
	StatsWidgetShadowOffset   = 4   // Shadow offset distance
	StatsWidgetPadding        = 15  // Internal padding
	StatsWidgetLineHeight     = 18  // Height of each text line
	StatsWidgetSectionSpacing = 25  // Space between sections
	StatsWidgetColumnWidth    = 280 // Width of each stat column
	StatsWidgetColumnSpacing  = 20  // Space between columns
	StatsWidgetHeaderHeight   = 30  // Height of section headers
	StatsWidgetTitleY         = 20  // Y offset for main title
	StatsWidgetContentStartY  = 50  // Y offset where content begins
	StatsWidgetBottomReserved = 40  // Space reserved at bottom for help text

	// Tab navigation
	StatsWidgetTabHeight  = 25 // Height of tab buttons
	StatsWidgetTabSpacing = 5  // Space between tabs
	StatsWidgetTabPadding = 10 // Internal tab padding

	// Progress bar dimensions
	StatsWidgetBarWidth  = 200 // Width of HP/MP/XP bars
	StatsWidgetBarHeight = 12  // Height of progress bars
	StatsWidgetBarBorder = 1   // Border width for bars
)

// Color constants - RGBA values for easy customization
const (
	// Background and border colors
	StatsWidgetBackgroundR = 25
	StatsWidgetBackgroundG = 25
	StatsWidgetBackgroundB = 35
	StatsWidgetBackgroundA = 245

	StatsWidgetBorderR = 120
	StatsWidgetBorderG = 120
	StatsWidgetBorderB = 120
	StatsWidgetBorderA = 255

	StatsWidgetShadowR = 0
	StatsWidgetShadowG = 0
	StatsWidgetShadowB = 0
	StatsWidgetShadowA = 120

	// Text colors
	StatsWidgetTitleR = 255
	StatsWidgetTitleG = 255
	StatsWidgetTitleB = 100
	StatsWidgetTitleA = 255

	StatsWidgetTextR = 255
	StatsWidgetTextG = 255
	StatsWidgetTextB = 255
	StatsWidgetTextA = 255

	StatsWidgetValueR = 200
	StatsWidgetValueG = 255
	StatsWidgetValueB = 200
	StatsWidgetValueA = 255

	// Section header colors
	StatsWidgetCoreHeaderR = 100
	StatsWidgetCoreHeaderG = 200
	StatsWidgetCoreHeaderB = 255
	StatsWidgetCoreHeaderA = 255

	StatsWidgetCombatHeaderR = 255
	StatsWidgetCombatHeaderG = 150
	StatsWidgetCombatHeaderB = 100
	StatsWidgetCombatHeaderA = 255

	StatsWidgetTacticalHeaderR = 150
	StatsWidgetTacticalHeaderG = 255
	StatsWidgetTacticalHeaderB = 150
	StatsWidgetTacticalHeaderA = 255

	// Progress bar colors
	StatsWidgetHPBarR = 200
	StatsWidgetHPBarG = 50
	StatsWidgetHPBarB = 50
	StatsWidgetHPBarA = 255

	StatsWidgetMPBarR = 50
	StatsWidgetMPBarG = 100
	StatsWidgetMPBarB = 200
	StatsWidgetMPBarA = 255

	StatsWidgetXPBarR = 255
	StatsWidgetXPBarG = 200
	StatsWidgetXPBarB = 50
	StatsWidgetXPBarA = 255

	StatsWidgetBarBackgroundR = 40
	StatsWidgetBarBackgroundG = 40
	StatsWidgetBarBackgroundB = 40
	StatsWidgetBarBackgroundA = 255

	// Tab colors
	StatsWidgetActiveTabR = 80
	StatsWidgetActiveTabG = 80
	StatsWidgetActiveTabB = 120
	StatsWidgetActiveTabA = 220

	StatsWidgetInactiveTabR = 50
	StatsWidgetInactiveTabG = 50
	StatsWidgetInactiveTabB = 70
	StatsWidgetInactiveTabA = 180
)

// StatCategory represents different categories of character stats
type StatCategory int

const (
	StatCategoryOverview StatCategory = iota // General overview with bars and basic info
	StatCategoryCore                         // Core stats: Level, XP, HP, MP
	StatCategoryCombat                       // Combat stats: ATK, DEF, Magic stats
	StatCategoryTactical                     // Tactical stats: Speed, Accuracy, Movement
)

// String returns the string representation of a StatCategory
func (sc StatCategory) String() string {
	switch sc {
	case StatCategoryOverview:
		return "Overview"
	case StatCategoryCore:
		return "Core Stats"
	case StatCategoryCombat:
		return "Combat Stats"
	case StatCategoryTactical:
		return "Tactical Stats"
	default:
		return "Unknown"
	}
}

// CharacterStatsWidget displays detailed character statistics in an organized layout
type CharacterStatsWidget struct {
	X, Y          int                           // Widget position
	Width, Height int                           // Widget dimensions
	Title         string                        // Widget title
	Visible       bool                          // Visibility state
	Character     *components.RPGStatsComponent // Character data to display

	// Navigation state
	CurrentCategory StatCategory   // Currently selected category
	Categories      []StatCategory // Available categories

	// Visual properties
	BackgroundColor color.RGBA
	BorderColor     color.RGBA
	ShadowColor     color.RGBA
	TitleColor      color.RGBA
	TextColor       color.RGBA
	ValueColor      color.RGBA
}

// NewCharacterStatsWidget creates a new character stats widget
func NewCharacterStatsWidget(x, y int, character *components.RPGStatsComponent) *CharacterStatsWidget {
	return &CharacterStatsWidget{
		X:         x,
		Y:         y,
		Width:     StatsWidgetWidth,
		Height:    StatsWidgetHeight,
		Title:     "Character Statistics",
		Visible:   false,
		Character: character,

		// Initialize with overview category
		CurrentCategory: StatCategoryOverview,
		Categories: []StatCategory{
			StatCategoryOverview,
			StatCategoryCore,
			StatCategoryCombat,
			StatCategoryTactical,
		},

		// Set default colors
		BackgroundColor: color.RGBA{StatsWidgetBackgroundR, StatsWidgetBackgroundG, StatsWidgetBackgroundB, StatsWidgetBackgroundA},
		BorderColor:     color.RGBA{StatsWidgetBorderR, StatsWidgetBorderG, StatsWidgetBorderB, StatsWidgetBorderA},
		ShadowColor:     color.RGBA{StatsWidgetShadowR, StatsWidgetShadowG, StatsWidgetShadowB, StatsWidgetShadowA},
		TitleColor:      color.RGBA{StatsWidgetTitleR, StatsWidgetTitleG, StatsWidgetTitleB, StatsWidgetTitleA},
		TextColor:       color.RGBA{StatsWidgetTextR, StatsWidgetTextG, StatsWidgetTextB, StatsWidgetTextA},
		ValueColor:      color.RGBA{StatsWidgetValueR, StatsWidgetValueG, StatsWidgetValueB, StatsWidgetValueA},
	}
}

// Show displays the widget
func (csw *CharacterStatsWidget) Show() {
	csw.Visible = true
}

// Hide hides the widget
func (csw *CharacterStatsWidget) Hide() {
	csw.Visible = false
}

// IsVisible returns whether the widget is currently visible
func (csw *CharacterStatsWidget) IsVisible() bool {
	return csw.Visible
}

// SetCharacter updates the character data displayed by the widget
func (csw *CharacterStatsWidget) SetCharacter(character *components.RPGStatsComponent) {
	csw.Character = character
}

// Update handles input and widget state updates
// Returns true if the widget consumed the ESC key in this frame
func (csw *CharacterStatsWidget) Update() bool {
	if !csw.Visible {
		return false
	}

	// Handle escape key to close
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		csw.Hide()
		return true // ESC key was consumed to close the widget
	}

	// Handle tab navigation (Left/Right arrows or Tab/Shift+Tab)
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		csw.previousCategory()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		csw.nextCategory()
	}

	return false // No ESC key consumed
}

// nextCategory switches to the next stat category
func (csw *CharacterStatsWidget) nextCategory() {
	currentIndex := int(csw.CurrentCategory)
	nextIndex := (currentIndex + 1) % len(csw.Categories)
	csw.CurrentCategory = csw.Categories[nextIndex]
}

// previousCategory switches to the previous stat category
func (csw *CharacterStatsWidget) previousCategory() {
	currentIndex := int(csw.CurrentCategory)
	prevIndex := currentIndex - 1
	if prevIndex < 0 {
		prevIndex = len(csw.Categories) - 1
	}
	csw.CurrentCategory = csw.Categories[prevIndex]
}

// Draw renders the character stats widget
func (csw *CharacterStatsWidget) Draw(screen *ebiten.Image) {
	if !csw.Visible || csw.Character == nil {
		return
	}

	// Draw shadow
	shadowX := float32(csw.X + StatsWidgetShadowOffset)
	shadowY := float32(csw.Y + StatsWidgetShadowOffset)
	vector.DrawFilledRect(screen, shadowX, shadowY, float32(csw.Width), float32(csw.Height), csw.ShadowColor, false)

	// Draw background
	vector.DrawFilledRect(screen, float32(csw.X), float32(csw.Y), float32(csw.Width), float32(csw.Height), csw.BackgroundColor, false)

	// Draw border
	vector.StrokeRect(screen, float32(csw.X), float32(csw.Y), float32(csw.Width), float32(csw.Height), StatsWidgetBorderWidth, csw.BorderColor, false)

	// Draw title
	titleText := fmt.Sprintf("%s - %s", csw.Title, csw.Character.Name)
	ebitenutil.DebugPrintAt(screen, titleText, csw.X+StatsWidgetPadding, csw.Y+StatsWidgetTitleY)

	// Draw category tabs
	csw.drawTabs(screen)

	// Draw content area wrapper (positioned below tabs with spacing)
	contentSpacing := 10 // Add spacing between tabs and content area
	contentY := float32(csw.Y + StatsWidgetContentStartY + StatsWidgetTabHeight + contentSpacing)
	contentHeight := float32(csw.Height - StatsWidgetContentStartY - StatsWidgetTabHeight - contentSpacing - StatsWidgetBottomReserved)
	contentX := float32(csw.X + StatsWidgetPadding)
	contentWidth := float32(csw.Width - 2*StatsWidgetPadding)

	// Draw subtle content area background
	contentBgColor := color.RGBA{
		StatsWidgetBackgroundR + 5,
		StatsWidgetBackgroundG + 5,
		StatsWidgetBackgroundB + 5,
		StatsWidgetBackgroundA - 20,
	}
	vector.DrawFilledRect(screen, contentX, contentY, contentWidth, contentHeight, contentBgColor, false)

	// Draw content area border
	contentBorderColor := color.RGBA{StatsWidgetBorderR, StatsWidgetBorderG, StatsWidgetBorderB, 100}
	vector.StrokeRect(screen, contentX, contentY, contentWidth, contentHeight, 1, contentBorderColor, false)

	// Draw content based on current category
	switch csw.CurrentCategory {
	case StatCategoryOverview:
		csw.drawOverview(screen)
	case StatCategoryCore:
		csw.drawCoreStats(screen)
	case StatCategoryCombat:
		csw.drawCombatStats(screen)
	case StatCategoryTactical:
		csw.drawTacticalStats(screen)
	}

	// Draw help text at bottom
	helpText := "Press ← → to switch tabs, ESC to close"
	ebitenutil.DebugPrintAt(screen, helpText, csw.X+StatsWidgetPadding, csw.Y+csw.Height-25)
}

// drawTabs renders the category tab buttons
func (csw *CharacterStatsWidget) drawTabs(screen *ebiten.Image) {
	tabY := float32(csw.Y + StatsWidgetContentStartY)
	tabX := float32(csw.X + StatsWidgetPadding)

	for _, category := range csw.Categories {
		tabWidth := float32(120) // Fixed width for each tab

		// Determine tab color based on selection
		var tabColor color.RGBA
		if category == csw.CurrentCategory {
			tabColor = color.RGBA{StatsWidgetActiveTabR, StatsWidgetActiveTabG, StatsWidgetActiveTabB, StatsWidgetActiveTabA}
		} else {
			tabColor = color.RGBA{StatsWidgetInactiveTabR, StatsWidgetInactiveTabG, StatsWidgetInactiveTabB, StatsWidgetInactiveTabA}
		}

		// Draw tab background
		vector.DrawFilledRect(screen, tabX, tabY, tabWidth, float32(StatsWidgetTabHeight), tabColor, false)

		// Draw tab border
		borderColor := color.RGBA{StatsWidgetBorderR, StatsWidgetBorderG, StatsWidgetBorderB, 150}
		vector.StrokeRect(screen, tabX, tabY, tabWidth, float32(StatsWidgetTabHeight), 1, borderColor, false)

		// Draw tab text
		tabText := category.String()
		textX := int(tabX) + StatsWidgetTabPadding
		textY := int(tabY) + 8 // Center text vertically (better positioning for 25px height)
		ebitenutil.DebugPrintAt(screen, tabText, textX, textY)

		tabX += tabWidth + float32(StatsWidgetTabSpacing)
	}
}

// drawOverview renders the overview category with progress bars and key stats
func (csw *CharacterStatsWidget) drawOverview(screen *ebiten.Image) {
	startY := csw.Y + StatsWidgetContentStartY + StatsWidgetTabHeight + StatsWidgetSectionSpacing
	char := csw.Character

	// Character basic info
	leftX := csw.X + StatsWidgetPadding
	rightX := csw.X + StatsWidgetPadding + StatsWidgetColumnWidth + StatsWidgetColumnSpacing

	// Left column - Basic info
	currentY := startY
	basicInfo := fmt.Sprintf("Name: %s", char.Name)
	ebitenutil.DebugPrintAt(screen, basicInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	levelInfo := fmt.Sprintf("Level: %d", char.Level)
	ebitenutil.DebugPrintAt(screen, levelInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	jobInfo := fmt.Sprintf("Job: %s", char.Job.String())
	ebitenutil.DebugPrintAt(screen, jobInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight * 2

	// HP Bar
	hpText := fmt.Sprintf("HP: %d/%d", char.CurrentHP, char.MaxHP)
	ebitenutil.DebugPrintAt(screen, hpText, leftX, currentY)
	currentY += StatsWidgetLineHeight + 5
	csw.drawProgressBar(screen, float32(leftX), float32(currentY), StatsWidgetBarWidth, StatsWidgetBarHeight,
		char.CurrentHP, char.MaxHP, color.RGBA{StatsWidgetHPBarR, StatsWidgetHPBarG, StatsWidgetHPBarB, StatsWidgetHPBarA})
	currentY += StatsWidgetBarHeight + StatsWidgetLineHeight

	// MP Bar
	mpText := fmt.Sprintf("MP: %d/%d", char.CurrentMP, char.MaxMP)
	ebitenutil.DebugPrintAt(screen, mpText, leftX, currentY)
	currentY += StatsWidgetLineHeight + 5
	csw.drawProgressBar(screen, float32(leftX), float32(currentY), StatsWidgetBarWidth, StatsWidgetBarHeight,
		char.CurrentMP, char.MaxMP, color.RGBA{StatsWidgetMPBarR, StatsWidgetMPBarG, StatsWidgetMPBarB, StatsWidgetMPBarA})
	currentY += StatsWidgetBarHeight + StatsWidgetLineHeight

	// XP Bar
	xpText := fmt.Sprintf("XP: %d/%d", char.Experience, char.ExpToNext)
	ebitenutil.DebugPrintAt(screen, xpText, leftX, currentY)
	currentY += StatsWidgetLineHeight + 5
	csw.drawProgressBar(screen, float32(leftX), float32(currentY), StatsWidgetBarWidth, StatsWidgetBarHeight,
		char.Experience, char.ExpToNext, color.RGBA{StatsWidgetXPBarR, StatsWidgetXPBarG, StatsWidgetXPBarB, StatsWidgetXPBarA})

	// Right column - Key stats summary
	currentY = startY
	ebitenutil.DebugPrintAt(screen, "Combat Stats:", rightX, currentY)
	currentY += StatsWidgetLineHeight + 10

	attackInfo := fmt.Sprintf("Attack: %d", char.Attack)
	ebitenutil.DebugPrintAt(screen, attackInfo, rightX, currentY)
	currentY += StatsWidgetLineHeight

	defenseInfo := fmt.Sprintf("Defense: %d", char.Defense)
	ebitenutil.DebugPrintAt(screen, defenseInfo, rightX, currentY)
	currentY += StatsWidgetLineHeight

	speedInfo := fmt.Sprintf("Speed: %d", char.Speed)
	ebitenutil.DebugPrintAt(screen, speedInfo, rightX, currentY)
	currentY += StatsWidgetLineHeight * 2

	ebitenutil.DebugPrintAt(screen, "Tactical Info:", rightX, currentY)
	currentY += StatsWidgetLineHeight + 10

	moveInfo := fmt.Sprintf("Move Range: %d", char.MoveRange)
	ebitenutil.DebugPrintAt(screen, moveInfo, rightX, currentY)
	currentY += StatsWidgetLineHeight

	movesLeft := fmt.Sprintf("Moves Left: %d", char.MovesRemaining)
	ebitenutil.DebugPrintAt(screen, movesLeft, rightX, currentY)
}

// drawProgressBar draws a progress bar with current/max values
func (csw *CharacterStatsWidget) drawProgressBar(screen *ebiten.Image, x, y, width, height float32, current, max int, barColor color.RGBA) {
	// Background
	bgColor := color.RGBA{StatsWidgetBarBackgroundR, StatsWidgetBarBackgroundG, StatsWidgetBarBackgroundB, StatsWidgetBarBackgroundA}
	vector.DrawFilledRect(screen, x, y, width, height, bgColor, false)

	// Progress fill
	if max > 0 {
		fillWidth := width * float32(current) / float32(max)
		if fillWidth > 0 {
			vector.DrawFilledRect(screen, x, y, fillWidth, height, barColor, false)
		}
	}

	// Border
	borderColor := color.RGBA{StatsWidgetBorderR, StatsWidgetBorderG, StatsWidgetBorderB, StatsWidgetBorderA}
	vector.StrokeRect(screen, x, y, width, height, float32(StatsWidgetBarBorder), borderColor, false)
}

// drawCoreStats renders the core stats category
func (csw *CharacterStatsWidget) drawCoreStats(screen *ebiten.Image) {
	startY := csw.Y + StatsWidgetContentStartY + StatsWidgetTabHeight + StatsWidgetSectionSpacing
	char := csw.Character
	leftX := csw.X + StatsWidgetPadding

	// Core stats header
	headerColor := color.RGBA{StatsWidgetCoreHeaderR, StatsWidgetCoreHeaderG, StatsWidgetCoreHeaderB, StatsWidgetCoreHeaderA}
	csw.drawSectionHeader(screen, "Core Statistics", leftX, startY, headerColor)

	currentY := startY + StatsWidgetHeaderHeight + 10

	// Character progression
	levelInfo := fmt.Sprintf("Level: %d", char.Level)
	ebitenutil.DebugPrintAt(screen, levelInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	expInfo := fmt.Sprintf("Experience: %d / %d", char.Experience, char.ExpToNext)
	ebitenutil.DebugPrintAt(screen, expInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	expToNext := char.ExpToNext - char.Experience
	nextLevelInfo := fmt.Sprintf("To Next Level: %d XP", expToNext)
	ebitenutil.DebugPrintAt(screen, nextLevelInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight * 2

	// Health and Mana detailed
	ebitenutil.DebugPrintAt(screen, "Vitals:", leftX, currentY)
	currentY += StatsWidgetLineHeight + 5

	hpInfo := fmt.Sprintf("Health Points: %d / %d", char.CurrentHP, char.MaxHP)
	ebitenutil.DebugPrintAt(screen, hpInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	mpInfo := fmt.Sprintf("Mana Points: %d / %d", char.CurrentMP, char.MaxMP)
	ebitenutil.DebugPrintAt(screen, mpInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	// Health/Mana percentages
	hpPercent := float32(char.CurrentHP) / float32(char.MaxHP) * 100
	mpPercent := float32(char.CurrentMP) / float32(char.MaxMP) * 100

	percentInfo := fmt.Sprintf("HP: %.1f%%  MP: %.1f%%", hpPercent, mpPercent)
	ebitenutil.DebugPrintAt(screen, percentInfo, leftX, currentY)
}

// drawCombatStats renders the combat stats category
func (csw *CharacterStatsWidget) drawCombatStats(screen *ebiten.Image) {
	startY := csw.Y + StatsWidgetContentStartY + StatsWidgetTabHeight + StatsWidgetSectionSpacing
	char := csw.Character
	leftX := csw.X + StatsWidgetPadding
	rightX := csw.X + StatsWidgetPadding + StatsWidgetColumnWidth + StatsWidgetColumnSpacing

	// Combat stats header
	headerColor := color.RGBA{StatsWidgetCombatHeaderR, StatsWidgetCombatHeaderG, StatsWidgetCombatHeaderB, StatsWidgetCombatHeaderA}
	csw.drawSectionHeader(screen, "Combat Statistics", leftX, startY, headerColor)

	currentY := startY + StatsWidgetHeaderHeight + 10

	// Left column - Physical combat
	ebitenutil.DebugPrintAt(screen, "Physical Combat:", leftX, currentY)
	leftCurrentY := currentY + StatsWidgetLineHeight + 5

	attackInfo := fmt.Sprintf("Attack Power: %d", char.Attack)
	ebitenutil.DebugPrintAt(screen, attackInfo, leftX, leftCurrentY)
	leftCurrentY += StatsWidgetLineHeight

	defenseInfo := fmt.Sprintf("Defense: %d", char.Defense)
	ebitenutil.DebugPrintAt(screen, defenseInfo, leftX, leftCurrentY)
	leftCurrentY += StatsWidgetLineHeight

	speedInfo := fmt.Sprintf("Speed: %d", char.Speed)
	ebitenutil.DebugPrintAt(screen, speedInfo, leftX, leftCurrentY)
	leftCurrentY += StatsWidgetLineHeight * 2

	ebitenutil.DebugPrintAt(screen, "Combat Chances:", leftX, leftCurrentY)
	leftCurrentY += StatsWidgetLineHeight + 5

	accuracyInfo := fmt.Sprintf("Accuracy: %d%%", char.Accuracy)
	ebitenutil.DebugPrintAt(screen, accuracyInfo, leftX, leftCurrentY)
	leftCurrentY += StatsWidgetLineHeight

	critInfo := fmt.Sprintf("Critical Rate: %d%%", char.CritRate)
	ebitenutil.DebugPrintAt(screen, critInfo, leftX, leftCurrentY)

	// Right column - Magical combat
	ebitenutil.DebugPrintAt(screen, "Magical Combat:", rightX, currentY+StatsWidgetLineHeight+5)
	rightCurrentY := currentY + StatsWidgetLineHeight*2 + 10

	magicAttackInfo := fmt.Sprintf("Magic Attack: %d", char.MagicAttack)
	ebitenutil.DebugPrintAt(screen, magicAttackInfo, rightX, rightCurrentY)
	rightCurrentY += StatsWidgetLineHeight

	magicDefenseInfo := fmt.Sprintf("Magic Defense: %d", char.MagicDefense)
	ebitenutil.DebugPrintAt(screen, magicDefenseInfo, rightX, rightCurrentY)
	rightCurrentY += StatsWidgetLineHeight * 2

	ebitenutil.DebugPrintAt(screen, "Job Class:", rightX, rightCurrentY)
	rightCurrentY += StatsWidgetLineHeight + 5

	jobInfo := fmt.Sprintf("Job: %s", char.Job.String())
	ebitenutil.DebugPrintAt(screen, jobInfo, rightX, rightCurrentY)
}

// drawTacticalStats renders the tactical stats category
func (csw *CharacterStatsWidget) drawTacticalStats(screen *ebiten.Image) {
	startY := csw.Y + StatsWidgetContentStartY + StatsWidgetTabHeight + StatsWidgetSectionSpacing
	char := csw.Character
	leftX := csw.X + StatsWidgetPadding

	// Tactical stats header
	headerColor := color.RGBA{StatsWidgetTacticalHeaderR, StatsWidgetTacticalHeaderG, StatsWidgetTacticalHeaderB, StatsWidgetTacticalHeaderA}
	csw.drawSectionHeader(screen, "Tactical Statistics", leftX, startY, headerColor)

	currentY := startY + StatsWidgetHeaderHeight + 10

	// Movement information
	ebitenutil.DebugPrintAt(screen, "Movement:", leftX, currentY)
	currentY += StatsWidgetLineHeight + 5

	moveRangeInfo := fmt.Sprintf("Movement Range: %d tiles", char.MoveRange)
	ebitenutil.DebugPrintAt(screen, moveRangeInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	movesLeftInfo := fmt.Sprintf("Moves Remaining: %d", char.MovesRemaining)
	ebitenutil.DebugPrintAt(screen, movesLeftInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	moveHistoryInfo := fmt.Sprintf("Moves This Turn: %d", len(char.MoveHistory))
	ebitenutil.DebugPrintAt(screen, moveHistoryInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight * 2

	// Combat positioning
	ebitenutil.DebugPrintAt(screen, "Combat Performance:", leftX, currentY)
	currentY += StatsWidgetLineHeight + 5

	speedInfo := fmt.Sprintf("Initiative (Speed): %d", char.Speed)
	ebitenutil.DebugPrintAt(screen, speedInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	accuracyInfo := fmt.Sprintf("Hit Chance: %d%%", char.Accuracy)
	ebitenutil.DebugPrintAt(screen, accuracyInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight

	critInfo := fmt.Sprintf("Critical Chance: %d%%", char.CritRate)
	ebitenutil.DebugPrintAt(screen, critInfo, leftX, currentY)
	currentY += StatsWidgetLineHeight * 2

	// Status
	aliveStatus := "Alive"
	if !char.IsAlive() {
		aliveStatus = "Defeated"
	}
	statusInfo := fmt.Sprintf("Status: %s", aliveStatus)
	ebitenutil.DebugPrintAt(screen, statusInfo, leftX, currentY)
}

// drawSectionHeader draws a colored section header
func (csw *CharacterStatsWidget) drawSectionHeader(screen *ebiten.Image, title string, x, y int, headerColor color.RGBA) {
	// Draw header background
	headerWidth := float32(csw.Width - StatsWidgetPadding*2)
	vector.DrawFilledRect(screen, float32(x), float32(y), headerWidth, float32(StatsWidgetHeaderHeight), headerColor, false)

	// Draw header text
	ebitenutil.DebugPrintAt(screen, title, x+10, y+8)
}
