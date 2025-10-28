package ui

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jrecuero/myrpg/internal/ecs/components"
)

// QuestJournalWidget provides a visual interface for quest management
type QuestJournalWidget struct {
	// Widget properties
	X, Y          int
	Width, Height int
	Visible       bool
	Enabled       bool

	questJournal     *components.QuestJournalComponent
	selectedTab      int // 0 = Active, 1 = Completed
	selectedQuestIdx int
	scrollOffset     int
	showDetails      bool
	selectedQuest    *components.Quest

	// UI Layout
	tabHeight    int
	questHeight  int
	maxVisible   int
	detailWidth  int
	detailHeight int

	// Colors
	colorBackground    color.RGBA
	colorBorder        color.RGBA
	colorText          color.RGBA
	colorTabActive     color.RGBA
	colorTabInactive   color.RGBA
	colorQuestActive   color.RGBA
	colorQuestComplete color.RGBA
	colorObjective     color.RGBA
	colorSelected      color.RGBA
}

// NewQuestJournalWidget creates a new quest journal widget
func NewQuestJournalWidget(x, y, width, height int, questJournal *components.QuestJournalComponent) *QuestJournalWidget {
	widget := &QuestJournalWidget{
		X:                x,
		Y:                y,
		Width:            width,
		Height:           height,
		Visible:          true,
		Enabled:          true,
		questJournal:     questJournal,
		selectedTab:      0,
		selectedQuestIdx: 0,
		scrollOffset:     0,
		showDetails:      false,

		// Layout
		tabHeight:    30,
		questHeight:  60,
		maxVisible:   (height - 50) / 60, // Account for tabs and padding
		detailWidth:  width/2 - 15,       // Half width minus padding when showing details
		detailHeight: height - 100,

		// Colors
		colorBackground:    color.RGBA{15, 15, 25, 240},
		colorBorder:        color.RGBA{100, 100, 120, 255},
		colorText:          color.RGBA{220, 220, 220, 255},
		colorTabActive:     color.RGBA{40, 60, 100, 255},
		colorTabInactive:   color.RGBA{25, 25, 35, 255},
		colorQuestActive:   color.RGBA{50, 70, 50, 255},
		colorQuestComplete: color.RGBA{70, 50, 70, 255},
		colorObjective:     color.RGBA{180, 180, 180, 255},
		colorSelected:      color.RGBA{80, 100, 140, 255},
	}

	return widget
}

// Update handles input and state changes for the quest journal
func (qw *QuestJournalWidget) Update() InputResult {
	result := NewInputResult()

	if !qw.Visible || !qw.Enabled {
		return result
	}

	// Handle ESC key to close quest journal
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		qw.Visible = false
		result.EscConsumed = true
		return result
	}

	// Get mouse position
	mouseX, mouseY := ebiten.CursorPosition()

	// Check if mouse is over the widget area
	isMouseOverWidget := mouseX >= qw.X && mouseX <= qw.X+qw.Width &&
		mouseY >= qw.Y && mouseY <= qw.Y+qw.Height

	if isMouseOverWidget {
		result.MouseConsumed = true
	}

	// Handle keyboard navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		qw.selectedTab = (qw.selectedTab + 1) % 2
		qw.selectedQuestIdx = 0
		qw.scrollOffset = 0
		qw.showDetails = false
	}

	// Handle quest navigation
	quests := qw.getCurrentQuests()
	if len(quests) > 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			if qw.selectedQuestIdx > 0 {
				qw.selectedQuestIdx--
			}
			qw.adjustScroll()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			if qw.selectedQuestIdx < len(quests)-1 {
				qw.selectedQuestIdx++
			}
			qw.adjustScroll()
		}

		// Handle Enter key to show/hide details
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if qw.selectedQuestIdx < len(quests) {
				qw.selectedQuest = quests[qw.selectedQuestIdx]
				qw.showDetails = !qw.showDetails
			}
		}
	}

	// Handle mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && isMouseOverWidget {
		qw.handleMouseClick(mouseX, mouseY)
	}

	return result
}

// getCurrentQuests returns quests for the currently selected tab
func (qw *QuestJournalWidget) getCurrentQuests() []*components.Quest {
	if qw.questJournal == nil {
		return []*components.Quest{}
	}

	var quests []*components.Quest
	if qw.selectedTab == 0 {
		quests = qw.questJournal.GetActiveQuests()
	} else {
		quests = qw.questJournal.GetCompletedQuests()
	}

	// Sort by quest type and level
	sort.Slice(quests, func(i, j int) bool {
		if quests[i].Type != quests[j].Type {
			return quests[i].Type < quests[j].Type
		}
		return quests[i].Level < quests[j].Level
	})

	return quests
}

// adjustScroll adjusts scroll offset to keep selected quest visible
func (qw *QuestJournalWidget) adjustScroll() {
	if qw.selectedQuestIdx < qw.scrollOffset {
		qw.scrollOffset = qw.selectedQuestIdx
	}
	if qw.selectedQuestIdx >= qw.scrollOffset+qw.maxVisible {
		qw.scrollOffset = qw.selectedQuestIdx - qw.maxVisible + 1
	}
}

// handleMouseClick processes mouse clicks on the widget
func (qw *QuestJournalWidget) handleMouseClick(mouseX, mouseY int) {
	// Check tab clicks
	tabY := qw.Y + 5
	if mouseY >= tabY && mouseY <= tabY+qw.tabHeight {
		tabWidth := qw.Width / 2
		if mouseX >= qw.X+5 && mouseX <= qw.X+5+tabWidth {
			qw.selectedTab = 0 // Active tab
			qw.selectedQuestIdx = 0
			qw.scrollOffset = 0
			qw.showDetails = false
		} else if mouseX >= qw.X+5+tabWidth && mouseX <= qw.X+qw.Width-5 {
			qw.selectedTab = 1 // Completed tab
			qw.selectedQuestIdx = 0
			qw.scrollOffset = 0
			qw.showDetails = false
		}
		return
	}

	// Check quest clicks
	questsY := qw.Y + qw.tabHeight + 15
	quests := qw.getCurrentQuests()

	for i := 0; i < qw.maxVisible && i+qw.scrollOffset < len(quests); i++ {
		questY := questsY + i*qw.questHeight
		if mouseY >= questY && mouseY <= questY+qw.questHeight {
			qw.selectedQuestIdx = i + qw.scrollOffset
			qw.selectedQuest = quests[qw.selectedQuestIdx]
			qw.showDetails = !qw.showDetails
			break
		}
	}
}

// Draw renders the quest journal widget
func (qw *QuestJournalWidget) Draw(screen *ebiten.Image) {
	if !qw.Visible {
		return
	}

	// Draw main background
	ebitenutil.DrawRect(screen, float64(qw.X), float64(qw.Y), float64(qw.Width), float64(qw.Height), qw.colorBackground)

	// Draw border
	borderThickness := 2
	ebitenutil.DrawRect(screen, float64(qw.X), float64(qw.Y), float64(qw.Width), float64(borderThickness), qw.colorBorder)
	ebitenutil.DrawRect(screen, float64(qw.X), float64(qw.Y), float64(borderThickness), float64(qw.Height), qw.colorBorder)
	ebitenutil.DrawRect(screen, float64(qw.X+qw.Width-borderThickness), float64(qw.Y), float64(borderThickness), float64(qw.Height), qw.colorBorder)
	ebitenutil.DrawRect(screen, float64(qw.X), float64(qw.Y+qw.Height-borderThickness), float64(qw.Width), float64(borderThickness), qw.colorBorder)

	// Draw title
	ebitenutil.DebugPrintAt(screen, "Quest Journal", qw.X+10, qw.Y+5)

	// Draw tabs
	qw.drawTabs(screen)

	// Draw quest list
	qw.drawQuestList(screen)

	// Draw quest details if selected
	if qw.showDetails && qw.selectedQuest != nil {
		qw.drawQuestDetails(screen)
	}
}

// drawTabs renders the tab buttons
func (qw *QuestJournalWidget) drawTabs(screen *ebiten.Image) {
	tabY := qw.Y + 25
	tabWidth := qw.Width / 2

	// Active tab
	activeColor := qw.colorTabInactive
	if qw.selectedTab == 0 {
		activeColor = qw.colorTabActive
	}
	ebitenutil.DrawRect(screen, float64(qw.X+5), float64(tabY), float64(tabWidth-5), float64(qw.tabHeight), activeColor)
	ebitenutil.DebugPrintAt(screen, "Active", qw.X+10, tabY+8)

	// Completed tab
	completedColor := qw.colorTabInactive
	if qw.selectedTab == 1 {
		completedColor = qw.colorTabActive
	}
	ebitenutil.DrawRect(screen, float64(qw.X+tabWidth), float64(tabY), float64(tabWidth-5), float64(qw.tabHeight), completedColor)
	ebitenutil.DebugPrintAt(screen, "Completed", qw.X+tabWidth+5, tabY+8)
}

// drawQuestList renders the list of quests
func (qw *QuestJournalWidget) drawQuestList(screen *ebiten.Image) {
	quests := qw.getCurrentQuests()
	if len(quests) == 0 {
		ebitenutil.DebugPrintAt(screen, "No quests available", qw.X+10, qw.Y+70)
		return
	}

	questsY := qw.Y + qw.tabHeight + 35

	// Calculate quest list width based on whether details are showing
	questListWidth := qw.Width - 10
	if qw.showDetails {
		questListWidth = qw.Width/2 - 10
	}

	// Draw visible quests
	for i := 0; i < qw.maxVisible && i+qw.scrollOffset < len(quests); i++ {
		questIdx := i + qw.scrollOffset
		quest := quests[questIdx]
		questY := questsY + i*qw.questHeight

		// Quest background
		questColor := qw.colorQuestActive
		if quest.IsCompleted() {
			questColor = qw.colorQuestComplete
		}
		if questIdx == qw.selectedQuestIdx {
			questColor = qw.colorSelected
		}

		ebitenutil.DrawRect(screen, float64(qw.X+5), float64(questY), float64(questListWidth), float64(qw.questHeight-5), questColor)

		// Quest title
		title := quest.Title
		maxTitleLen := 25
		if qw.showDetails {
			maxTitleLen = 15 // Shorter title when details panel is visible
		}
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-3] + "..."
		}
		ebitenutil.DebugPrintAt(screen, title, qw.X+10, questY+5)

		// Quest type and level
		typeText := qw.getQuestTypeText(quest.Type)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s (Lv.%d)", typeText, quest.Level), qw.X+10, questY+20)

		// Progress
		if quest.IsActive() {
			completed := quest.GetCompletedObjectives()
			required := quest.GetRequiredObjectives()
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Progress: %d/%d", completed, required), qw.X+10, questY+35)
		} else {
			ebitenutil.DebugPrintAt(screen, "COMPLETED", qw.X+10, questY+35)
		}
	}

	// Draw scroll indicator
	if len(quests) > qw.maxVisible {
		scrollY := qw.Y + qw.tabHeight + 35
		scrollHeight := qw.maxVisible * qw.questHeight
		scrollBarHeight := int(float64(scrollHeight) * float64(qw.maxVisible) / float64(len(quests)))
		scrollBarY := scrollY + int(float64(scrollHeight-scrollBarHeight)*float64(qw.scrollOffset)/float64(len(quests)-qw.maxVisible))

		ebitenutil.DrawRect(screen, float64(qw.X+qw.Width-10), float64(scrollY), 5, float64(scrollHeight), color.RGBA{50, 50, 50, 255})
		ebitenutil.DrawRect(screen, float64(qw.X+qw.Width-10), float64(scrollBarY), 5, float64(scrollBarHeight), qw.colorBorder)
	}
}

// drawQuestDetails renders detailed information about the selected quest
func (qw *QuestJournalWidget) drawQuestDetails(screen *ebiten.Image) {
	if qw.selectedQuest == nil {
		return
	}

	// Details background - position in the right half of the widget
	detailX := qw.X + qw.Width/2 + 5
	detailY := qw.Y + qw.tabHeight + 35
	detailWidth := qw.Width/2 - 10
	detailHeight := qw.Height - qw.tabHeight - 35
	ebitenutil.DrawRect(screen, float64(detailX), float64(detailY), float64(detailWidth), float64(detailHeight), qw.colorBackground)

	// Border
	ebitenutil.DrawRect(screen, float64(detailX), float64(detailY), float64(detailWidth), 2, qw.colorBorder)
	ebitenutil.DrawRect(screen, float64(detailX), float64(detailY), 2, float64(detailHeight), qw.colorBorder)
	ebitenutil.DrawRect(screen, float64(detailX+detailWidth-2), float64(detailY), 2, float64(detailHeight), qw.colorBorder)
	ebitenutil.DrawRect(screen, float64(detailX), float64(detailY+detailHeight-2), float64(detailWidth), 2, qw.colorBorder)

	// Quest title
	ebitenutil.DebugPrintAt(screen, qw.selectedQuest.Title, detailX+5, detailY+5)

	// Quest description
	desc := qw.selectedQuest.Description
	lines := qw.wrapText(desc, 25) // Narrower text for split panel layout
	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, detailX+5, detailY+25+i*15)
	}

	// Objectives
	objY := detailY + 25 + len(lines)*15 + 10
	ebitenutil.DebugPrintAt(screen, "Objectives:", detailX+5, objY)
	objY += 15

	for _, obj := range qw.selectedQuest.Objectives {
		progressText := obj.GetProgressText()
		ebitenutil.DebugPrintAt(screen, progressText, detailX+10, objY)
		objY += 15
	}

	// Rewards
	if qw.selectedQuest.Reward != nil {
		objY += 10
		ebitenutil.DebugPrintAt(screen, "Rewards:", detailX+5, objY)
		objY += 15

		reward := qw.selectedQuest.Reward
		if reward.Experience > 0 {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Experience: %d", reward.Experience), detailX+10, objY)
			objY += 15
		}
		if reward.Gold > 0 {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Gold: %d", reward.Gold), detailX+10, objY)
			objY += 15
		}
		if reward.SkillPoints > 0 {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Skill Points: %d", reward.SkillPoints), detailX+10, objY)
			objY += 15
		}
	}
}

// getQuestTypeText returns a display string for quest types
func (qw *QuestJournalWidget) getQuestTypeText(questType components.QuestType) string {
	switch questType {
	case components.QuestTypeMain:
		return "Main"
	case components.QuestTypeSide:
		return "Side"
	case components.QuestTypeDaily:
		return "Daily"
	case components.QuestTypeTutorial:
		return "Tutorial"
	default:
		return "Quest"
	}
}

// wrapText wraps text to fit within a specified character width
func (qw *QuestJournalWidget) wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := make([]string, 0)

	// Simple word splitting
	currentWord := ""
	for _, char := range text {
		if char == ' ' || char == '\n' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}
	if currentWord != "" {
		words = append(words, currentWord)
	}

	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
