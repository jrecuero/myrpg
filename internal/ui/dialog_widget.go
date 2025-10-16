// Package ui provides dialog widget for NPC conversations and branching narratives
package ui

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jrecuero/myrpg/internal/logger"
)

// Dialog Widget Constants
const (
	// Layout dimensions
	DialogWidgetWidth        = 700 // Widget width in pixels
	DialogWidgetHeight       = 400 // Widget height in pixels
	DialogWidgetBorderWidth  = 2   // Border thickness
	DialogWidgetShadowOffset = 6   // Shadow offset distance
	DialogWidgetPadding      = 20  // Internal padding

	// Portrait area
	DialogPortraitWidth  = 120 // Portrait width in pixels
	DialogPortraitHeight = 120 // Portrait height in pixels
	DialogPortraitX      = 30  // X position from widget left
	DialogPortraitY      = 30  // Y position from widget top
	DialogPortraitBorder = 2   // Portrait border thickness

	// Text area
	DialogTextX          = 170 // Text area X position (after portrait)
	DialogTextY          = 30  // Text area Y position
	DialogTextWidth      = 500 // Text area width
	DialogTextHeight     = 180 // Text area height
	DialogTextLineHeight = 18  // Height per text line
	DialogTextPadding    = 10  // Text area internal padding

	// Speaker name area
	DialogSpeakerNameX      = 170 // Speaker name X position
	DialogSpeakerNameY      = 35  // Speaker name Y position
	DialogSpeakerNameHeight = 25  // Height reserved for speaker name

	// Dialog text content
	DialogContentY      = 65  // Dialog content Y position (after speaker name)
	DialogContentHeight = 140 // Height available for dialog content
	DialogMaxLines      = 7   // Maximum visible text lines

	// Choice buttons area
	DialogChoicesY          = 230 // Choices area Y position
	DialogChoicesHeight     = 140 // Height available for choices
	DialogChoiceHeight      = 25  // Height per choice button
	DialogChoiceSpacing     = 5   // Space between choice buttons
	DialogChoicePadding     = 8   // Internal padding for choice buttons
	DialogMaxVisibleChoices = 5   // Maximum visible choices before scrolling

	// Typewriter effect
	DialogTypewriterSpeed     = 50   // Milliseconds per character
	DialogTypewriterFastSpeed = 20   // Fast mode milliseconds per character
	DialogAutoAdvanceDelay    = 3000 // Auto-advance delay in milliseconds
)

// Color constants - RGBA values for easy customization
const (
	// Background colors
	DialogWidgetBackgroundR = 25
	DialogWidgetBackgroundG = 25
	DialogWidgetBackgroundB = 35
	DialogWidgetBackgroundA = 240

	// Border colors
	DialogWidgetBorderR = 100
	DialogWidgetBorderG = 100
	DialogWidgetBorderB = 120
	DialogWidgetBorderA = 255

	// Text colors
	DialogSpeakerNameR = 255
	DialogSpeakerNameG = 220
	DialogSpeakerNameB = 100
	DialogSpeakerNameA = 255
	DialogContentTextR = 240
	DialogContentTextG = 240
	DialogContentTextB = 240
	DialogContentTextA = 255

	// Choice button colors
	DialogChoiceNormalR   = 60
	DialogChoiceNormalG   = 60
	DialogChoiceNormalB   = 80
	DialogChoiceNormalA   = 255
	DialogChoiceSelectedR = 100
	DialogChoiceSelectedG = 150
	DialogChoiceSelectedB = 200
	DialogChoiceSelectedA = 255
	DialogChoiceTextR     = 255
	DialogChoiceTextG     = 255
	DialogChoiceTextB     = 255
	DialogChoiceTextA     = 255

	// Portrait colors
	DialogPortraitBorderR      = 150
	DialogPortraitBorderG      = 150
	DialogPortraitBorderB      = 150
	DialogPortraitBorderA      = 255
	DialogPortraitPlaceholderR = 80
	DialogPortraitPlaceholderG = 80
	DialogPortraitPlaceholderB = 100
	DialogPortraitPlaceholderA = 255

	// Shadow colors
	DialogWidgetShadowR = 0
	DialogWidgetShadowG = 0
	DialogWidgetShadowB = 0
	DialogWidgetShadowA = 120
)

// DialogState represents the current state of the dialog widget
type DialogState int

const (
	DialogStateHidden DialogState = iota
	DialogStateTyping
	DialogStateWaiting
	DialogStateChoices
	DialogStateAutoAdvance
)

// String returns the string representation of DialogState
func (s DialogState) String() string {
	switch s {
	case DialogStateHidden:
		return "Hidden"
	case DialogStateTyping:
		return "Typing"
	case DialogStateWaiting:
		return "Waiting"
	case DialogStateChoices:
		return "Choices"
	case DialogStateAutoAdvance:
		return "AutoAdvance"
	default:
		return "Unknown"
	}
}

// Condition represents a condition for dialog branching
type Condition struct {
	Type     string      `json:"type"`
	Name     string      `json:"name,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	ItemID   string      `json:"item_id,omitempty"`
	Quantity int         `json:"quantity,omitempty"`
	QuestID  string      `json:"quest_id,omitempty"`
	Status   string      `json:"status,omitempty"`
	JobType  string      `json:"job_type,omitempty"`
}

// Action represents an action to perform during dialog
type Action struct {
	Type     string      `json:"type"`
	Name     string      `json:"name,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	Amount   int         `json:"amount,omitempty"`
	ItemID   string      `json:"item_id,omitempty"`
	Quantity int         `json:"quantity,omitempty"`
	QuestID  string      `json:"quest_id,omitempty"`
	SoundID  string      `json:"sound_id,omitempty"`
	ShopID   string      `json:"shop_id,omitempty"`
	Flag     string      `json:"flag,omitempty"`
}

// Choice represents a player choice in dialog
type Choice struct {
	Text       string      `json:"text"`
	Target     string      `json:"target"`
	Conditions []Condition `json:"conditions,omitempty"`
	Actions    []Action    `json:"actions,omitempty"`
}

// DialogNode represents a single dialog node
type DialogNode struct {
	Speaker    string      `json:"speaker"`
	Text       string      `json:"text"`
	Conditions []Condition `json:"conditions,omitempty"`
	Actions    []Action    `json:"actions,omitempty"`
	Choices    []Choice    `json:"choices,omitempty"`
	End        bool        `json:"end,omitempty"`
	Next       string      `json:"next,omitempty"`
}

// Character represents a character definition
type Character struct {
	Name        string `json:"name"`
	Portrait    string `json:"portrait"`
	VoiceStyle  string `json:"voice_style"`
	Description string `json:"description"`
}

// DialogScript represents a complete dialog script
type DialogScript struct {
	DialogID    string                 `json:"dialog_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Variables   map[string]interface{} `json:"variables"`
	Nodes       map[string]DialogNode  `json:"nodes"`
}

// CharacterDefinitions represents the character definitions file
type CharacterDefinitions struct {
	Characters map[string]Character `json:"characters"`
}

// DialogWidget handles NPC conversations and branching narratives
type DialogWidget struct {
	// Widget properties
	X, Y          int
	Width, Height int
	Visible       bool

	// Dialog system
	Script      *DialogScript
	Characters  map[string]Character
	Variables   map[string]interface{}
	CurrentNode string
	State       DialogState

	// Visual state
	DisplayedText      string
	FullText           string
	TypewriterIndex    int
	LastTypewriterTime time.Time
	TypewriterSpeed    time.Duration
	FastMode           bool

	// Choice system
	AvailableChoices   []Choice
	SelectedChoice     int
	ChoiceScrollOffset int

	// Callbacks
	OnDialogEnd     func()
	OnVariableSet   func(name string, value interface{})
	OnActionExecute func(action Action)

	// Input state - removed unused fields
}

// NewDialogWidget creates a new dialog widget
func NewDialogWidget(screenWidth, screenHeight int) *DialogWidget {
	x := (screenWidth - DialogWidgetWidth) / 2
	y := (screenHeight - DialogWidgetHeight) / 2

	return &DialogWidget{
		X:               x,
		Y:               y,
		Width:           DialogWidgetWidth,
		Height:          DialogWidgetHeight,
		Visible:         false,
		State:           DialogStateHidden,
		Characters:      make(map[string]Character),
		Variables:       make(map[string]interface{}),
		TypewriterSpeed: time.Millisecond * DialogTypewriterSpeed,
		SelectedChoice:  0,
	}
}

// LoadCharacters loads character definitions from JSON file
func (dw *DialogWidget) LoadCharacters(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read characters file: %w", err)
	}

	var charDefs CharacterDefinitions
	if err := json.Unmarshal(data, &charDefs); err != nil {
		return fmt.Errorf("failed to parse characters JSON: %w", err)
	}

	dw.Characters = charDefs.Characters
	logger.Info("Loaded %d character definitions", len(dw.Characters))
	return nil
}

// LoadDialog loads a dialog script from JSON file
func (dw *DialogWidget) LoadDialog(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read dialog file: %w", err)
	}

	var script DialogScript
	if err := json.Unmarshal(data, &script); err != nil {
		return fmt.Errorf("failed to parse dialog JSON: %w", err)
	}

	dw.Script = &script

	// Initialize variables from script
	if script.Variables != nil {
		for k, v := range script.Variables {
			dw.Variables[k] = v
		}
	}

	logger.Info("Loaded dialog script: %s (%s)", script.DialogID, script.Title)
	return nil
}

// StartDialog begins a dialog from the specified node
func (dw *DialogWidget) StartDialog(nodeID string) {
	if dw.Script == nil {
		logger.Error("Cannot start dialog: no script loaded")
		return
	}

	node, exists := dw.Script.Nodes[nodeID]
	if !exists {
		logger.Error("Dialog node not found: %s", nodeID)
		return
	}

	dw.Visible = true
	dw.CurrentNode = nodeID
	dw.processNode(node)
	logger.Info("Started dialog at node: %s", nodeID)
}

// processNode handles the logic for a dialog node
func (dw *DialogWidget) processNode(node DialogNode) {
	// Execute actions
	for _, action := range node.Actions {
		dw.executeAction(action)
	}

	// Set up text display
	dw.FullText = dw.processText(node.Text)
	dw.DisplayedText = ""
	dw.TypewriterIndex = 0
	dw.LastTypewriterTime = time.Now()

	// Handle node type
	if len(node.Choices) > 0 {
		// Show choices after text is complete
		dw.State = DialogStateTyping
		dw.AvailableChoices = dw.filterChoices(node.Choices)
		dw.SelectedChoice = 0
		dw.ChoiceScrollOffset = 0
	} else if node.End {
		// Add automatic "End dialog" choice for end nodes
		dw.State = DialogStateTyping
		endChoice := Choice{
			Text:   "End dialog",
			Target: "__END_DIALOG__", // Special target to close dialog
		}
		dw.AvailableChoices = []Choice{endChoice}
		dw.SelectedChoice = 0
		dw.ChoiceScrollOffset = 0
	} else if node.Next != "" {
		// Auto-advance to next node
		dw.State = DialogStateTyping
	} else {
		// Wait for user input
		dw.State = DialogStateTyping
	}
}

// processText handles variable substitution in dialog text
func (dw *DialogWidget) processText(text string) string {
	// Replace variables in format ${variable_name}
	for name, value := range dw.Variables {
		placeholder := fmt.Sprintf("${%s}", name)
		replacement := fmt.Sprintf("%v", value)
		text = strings.ReplaceAll(text, placeholder, replacement)
	}

	// TODO: Add conditional text blocks {{if condition}}text{{endif}}

	return text
}

// filterChoices returns only choices whose conditions are met
func (dw *DialogWidget) filterChoices(choices []Choice) []Choice {
	var available []Choice
	for _, choice := range choices {
		if dw.checkConditions(choice.Conditions) {
			available = append(available, choice)
		}
	}
	return available
}

// checkConditions evaluates whether all conditions in a list are met
func (dw *DialogWidget) checkConditions(conditions []Condition) bool {
	for _, condition := range conditions {
		if !dw.checkCondition(condition) {
			return false
		}
	}
	return true
}

// checkCondition evaluates a single condition
func (dw *DialogWidget) checkCondition(condition Condition) bool {
	switch condition.Type {
	case "variable":
		return dw.checkVariableCondition(condition)
	case "level":
		return dw.checkLevelCondition(condition)
	case "has_item":
		return dw.checkItemCondition(condition)
	case "quest_status":
		return dw.checkQuestCondition(condition)
	case "job":
		return dw.checkJobCondition(condition)
	default:
		logger.Error("Unknown condition type: %s", condition.Type)
		return false
	}
}

// checkVariableCondition evaluates variable-based conditions
func (dw *DialogWidget) checkVariableCondition(condition Condition) bool {
	value, exists := dw.Variables[condition.Name]
	if !exists {
		return false
	}

	switch condition.Operator {
	case "equals":
		return value == condition.Value
	case "not_equals":
		return value != condition.Value
	case "greater_than":
		if v, ok := value.(float64); ok {
			if cv, ok := condition.Value.(float64); ok {
				return v > cv
			}
		}
	case "less_than":
		if v, ok := value.(float64); ok {
			if cv, ok := condition.Value.(float64); ok {
				return v < cv
			}
		}
	}

	return false
}

// Placeholder condition checks - these would integrate with actual game systems
func (dw *DialogWidget) checkLevelCondition(_ Condition) bool {
	// TODO: Integrate with player level system
	return true
}

func (dw *DialogWidget) checkItemCondition(_ Condition) bool {
	// TODO: Integrate with inventory system
	return true
}

func (dw *DialogWidget) checkQuestCondition(_ Condition) bool {
	// TODO: Integrate with quest system
	return true
}

func (dw *DialogWidget) checkJobCondition(_ Condition) bool {
	// TODO: Integrate with player job/class system
	return true
}

// executeAction performs a dialog action
func (dw *DialogWidget) executeAction(action Action) {
	switch action.Type {
	case "set_variable":
		dw.Variables[action.Name] = action.Value
		if dw.OnVariableSet != nil {
			dw.OnVariableSet(action.Name, action.Value)
		}
		logger.Debug("Set dialog variable %s = %v", action.Name, action.Value)
	case "add_experience":
		// TODO: Integrate with experience system
		logger.Info("Would add %d experience", action.Amount)
	case "give_item":
		// TODO: Integrate with inventory system
		logger.Info("Would give item %s x%d", action.ItemID, action.Quantity)
	case "start_quest":
		// TODO: Integrate with quest system
		logger.Info("Would start quest: %s", action.QuestID)
	case "play_sound":
		// TODO: Integrate with audio system
		logger.Info("Would play sound: %s", action.SoundID)
	case "open_shop":
		// TODO: Integrate with shop system
		logger.Info("Would open shop: %s", action.ShopID)
	case "set_flag":
		// TODO: Integrate with game flag system
		logger.Info("Would set flag %s = %v", action.Flag, action.Value)
	default:
		logger.Error("Unknown action type: %s", action.Type)
	}

	if dw.OnActionExecute != nil {
		dw.OnActionExecute(action)
	}
}

// Update handles input and updates dialog state
func (dw *DialogWidget) Update() bool {
	if !dw.Visible {
		return false
	}

	escConsumed := false

	// Handle ESC key to close dialog
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		dw.Hide()
		escConsumed = true
	}

	// Update typewriter effect
	if dw.State == DialogStateTyping {
		dw.updateTypewriter()
	}

	// Handle input based on state
	switch dw.State {
	case DialogStateTyping:
		dw.handleTypingInput()
	case DialogStateWaiting:
		dw.handleWaitingInput()
	case DialogStateChoices:
		dw.handleChoicesInput()
	}

	return escConsumed
}

// updateTypewriter handles the typewriter text effect
func (dw *DialogWidget) updateTypewriter() {
	if dw.TypewriterIndex >= len(dw.FullText) {
		// Typing complete
		if len(dw.AvailableChoices) > 0 {
			dw.State = DialogStateChoices
		} else {
			dw.State = DialogStateWaiting
		}
		return
	}

	now := time.Now()
	if now.Sub(dw.LastTypewriterTime) >= dw.TypewriterSpeed {
		dw.TypewriterIndex++
		dw.DisplayedText = dw.FullText[:dw.TypewriterIndex]
		dw.LastTypewriterTime = now
	}
}

// handleTypingInput handles input during typewriter effect
func (dw *DialogWidget) handleTypingInput() {
	// Space or Enter to skip typewriter
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		dw.TypewriterIndex = len(dw.FullText)
		dw.DisplayedText = dw.FullText

		if len(dw.AvailableChoices) > 0 {
			dw.State = DialogStateChoices
		} else {
			dw.State = DialogStateWaiting
		}
	}

	// Tab to toggle fast mode
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		dw.FastMode = !dw.FastMode
		if dw.FastMode {
			dw.TypewriterSpeed = time.Millisecond * DialogTypewriterFastSpeed
		} else {
			dw.TypewriterSpeed = time.Millisecond * DialogTypewriterSpeed
		}
	}
}

// handleWaitingInput handles input when waiting for user to continue
func (dw *DialogWidget) handleWaitingInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		dw.advanceDialog()
	}
}

// handleChoicesInput handles input when showing choices
func (dw *DialogWidget) handleChoicesInput() {
	// Navigate choices
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if dw.SelectedChoice > 0 {
			dw.SelectedChoice--
			dw.updateChoiceScroll()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if dw.SelectedChoice < len(dw.AvailableChoices)-1 {
			dw.SelectedChoice++
			dw.updateChoiceScroll()
		}
	}

	// Select choice
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if dw.SelectedChoice < len(dw.AvailableChoices) {
			dw.selectChoice(dw.AvailableChoices[dw.SelectedChoice])
		}
	}
}

// updateChoiceScroll updates the choice scroll offset
func (dw *DialogWidget) updateChoiceScroll() {
	if dw.SelectedChoice < dw.ChoiceScrollOffset {
		dw.ChoiceScrollOffset = dw.SelectedChoice
	} else if dw.SelectedChoice >= dw.ChoiceScrollOffset+DialogMaxVisibleChoices {
		dw.ChoiceScrollOffset = dw.SelectedChoice - DialogMaxVisibleChoices + 1
	}
}

// selectChoice handles choice selection
func (dw *DialogWidget) selectChoice(choice Choice) {
	// Execute choice actions
	for _, action := range choice.Actions {
		dw.executeAction(action)
	}

	// Move to target node
	if choice.Target != "" {
		// Check for special end dialog target
		if choice.Target == "__END_DIALOG__" {
			logger.Info("Dialog ended by user choice")
			dw.Hide()
			return
		}

		if node, exists := dw.Script.Nodes[choice.Target]; exists {
			dw.CurrentNode = choice.Target
			dw.processNode(node)
		} else {
			logger.Error("Choice target node not found: %s", choice.Target)
			dw.Hide()
		}
	} else {
		dw.Hide()
	}
}

// advanceDialog advances to the next node or ends dialog
func (dw *DialogWidget) advanceDialog() {
	if dw.Script == nil {
		dw.Hide()
		return
	}

	currentNode, exists := dw.Script.Nodes[dw.CurrentNode]
	if !exists {
		dw.Hide()
		return
	}

	if currentNode.End {
		dw.Hide()
	} else if currentNode.Next != "" {
		if nextNode, exists := dw.Script.Nodes[currentNode.Next]; exists {
			dw.CurrentNode = currentNode.Next
			dw.processNode(nextNode)
		} else {
			logger.Error("Next node not found: %s", currentNode.Next)
			dw.Hide()
		}
	} else {
		dw.Hide()
	}
}

// Show makes the dialog widget visible
func (dw *DialogWidget) Show() {
	dw.Visible = true
	dw.State = DialogStateWaiting
}

// Hide closes the dialog widget
func (dw *DialogWidget) Hide() {
	dw.Visible = false
	dw.State = DialogStateHidden
	dw.DisplayedText = ""
	dw.FullText = ""
	dw.TypewriterIndex = 0
	dw.AvailableChoices = nil
	dw.SelectedChoice = 0
	dw.ChoiceScrollOffset = 0

	if dw.OnDialogEnd != nil {
		dw.OnDialogEnd()
	}
}

// IsVisible returns whether the widget is visible
func (dw *DialogWidget) IsVisible() bool {
	return dw.Visible
}

// Draw renders the dialog widget
func (dw *DialogWidget) Draw(screen *ebiten.Image) {
	if !dw.Visible || dw.Script == nil {
		return
	}

	// Draw shadow
	shadowColor := color.RGBA{DialogWidgetShadowR, DialogWidgetShadowG, DialogWidgetShadowB, DialogWidgetShadowA}
	vector.FillRect(screen,
		float32(dw.X+DialogWidgetShadowOffset), float32(dw.Y+DialogWidgetShadowOffset),
		float32(dw.Width), float32(dw.Height),
		shadowColor, false)

	// Draw main background
	bgColor := color.RGBA{DialogWidgetBackgroundR, DialogWidgetBackgroundG, DialogWidgetBackgroundB, DialogWidgetBackgroundA}
	vector.FillRect(screen,
		float32(dw.X), float32(dw.Y),
		float32(dw.Width), float32(dw.Height),
		bgColor, false)

	// Draw border
	borderColor := color.RGBA{DialogWidgetBorderR, DialogWidgetBorderG, DialogWidgetBorderB, DialogWidgetBorderA}
	for i := 0; i < DialogWidgetBorderWidth; i++ {
		vector.StrokeRect(screen,
			float32(dw.X+i), float32(dw.Y+i),
			float32(dw.Width-i*2), float32(dw.Height-i*2),
			1, borderColor, false)
	}

	// Draw portrait
	dw.drawPortrait(screen)

	// Draw speaker name and dialog text
	dw.drawDialogText(screen)

	// Draw choices if available
	if dw.State == DialogStateChoices && len(dw.AvailableChoices) > 0 {
		dw.drawChoices(screen)
	}

	// Draw help text
	dw.drawHelpText(screen)
}

// drawPortrait renders the character portrait
func (dw *DialogWidget) drawPortrait(screen *ebiten.Image) {
	portraitX := dw.X + DialogPortraitX
	portraitY := dw.Y + DialogPortraitY

	// Draw portrait background (placeholder)
	placeholderColor := color.RGBA{DialogPortraitPlaceholderR, DialogPortraitPlaceholderG, DialogPortraitPlaceholderB, DialogPortraitPlaceholderA}
	vector.FillRect(screen,
		float32(portraitX), float32(portraitY),
		float32(DialogPortraitWidth), float32(DialogPortraitHeight),
		placeholderColor, false)

	// Draw portrait border
	borderColor := color.RGBA{DialogPortraitBorderR, DialogPortraitBorderG, DialogPortraitBorderB, DialogPortraitBorderA}
	for i := 0; i < DialogPortraitBorder; i++ {
		vector.StrokeRect(screen,
			float32(portraitX-i), float32(portraitY-i),
			float32(DialogPortraitWidth+i*2), float32(DialogPortraitHeight+i*2),
			1, borderColor, false)
	}

	// TODO: Load and draw actual portrait image
	// For now, show character initial or placeholder text
	if dw.Script != nil {
		if currentNode, exists := dw.Script.Nodes[dw.CurrentNode]; exists {
			if character, exists := dw.Characters[currentNode.Speaker]; exists {
				// Draw character initial as placeholder
				initial := strings.ToUpper(string(character.Name[0]))
				ebitenutil.DebugPrintAt(screen, initial, portraitX+DialogPortraitWidth/2-5, portraitY+DialogPortraitHeight/2-5)
			}
		}
	}
}

// drawDialogText renders the speaker name and dialog content
func (dw *DialogWidget) drawDialogText(screen *ebiten.Image) {
	if dw.Script == nil {
		return
	}

	currentNode, exists := dw.Script.Nodes[dw.CurrentNode]
	if !exists {
		return
	}

	textX := dw.X + DialogSpeakerNameX
	textY := dw.Y + DialogSpeakerNameY

	// Draw speaker name
	character, hasCharacter := dw.Characters[currentNode.Speaker]
	speakerName := currentNode.Speaker
	if hasCharacter {
		speakerName = character.Name
	}

	ebitenutil.DebugPrintAt(screen, speakerName, textX, textY)

	// Draw dialog content
	contentY := dw.Y + DialogContentY
	dw.drawWrappedText(screen, dw.DisplayedText, textX, contentY, DialogTextWidth, DialogMaxLines)
}

// drawWrappedText renders text with word wrapping
func (dw *DialogWidget) drawWrappedText(screen *ebiten.Image, text string, x, y, maxWidth, maxLines int) {
	if text == "" {
		return
	}

	words := strings.Fields(text)
	lines := []string{""}
	currentLine := 0

	// Simple word wrapping (approximate character width of 6 pixels)
	charWidth := 6
	maxCharsPerLine := maxWidth / charWidth

	for _, word := range words {
		if len(lines[currentLine])+len(word)+1 <= maxCharsPerLine {
			if lines[currentLine] != "" {
				lines[currentLine] += " "
			}
			lines[currentLine] += word
		} else {
			if currentLine < maxLines-1 {
				currentLine++
				lines = append(lines, word)
			} else {
				// Truncate if too many lines
				break
			}
		}
	}

	// Draw lines
	for i, line := range lines {
		if i < maxLines {
			ebitenutil.DebugPrintAt(screen, line, x, y+i*DialogTextLineHeight)
		}
	}
}

// drawChoices renders the choice buttons
func (dw *DialogWidget) drawChoices(screen *ebiten.Image) {
	choicesY := dw.Y + DialogChoicesY

	for i := 0; i < DialogMaxVisibleChoices && i < len(dw.AvailableChoices); i++ {
		choiceIndex := dw.ChoiceScrollOffset + i
		if choiceIndex >= len(dw.AvailableChoices) {
			break
		}

		choice := dw.AvailableChoices[choiceIndex]
		buttonY := choicesY + i*(DialogChoiceHeight+DialogChoiceSpacing)
		buttonX := dw.X + DialogWidgetPadding

		// Determine button color
		var buttonColor color.RGBA
		if choiceIndex == dw.SelectedChoice {
			buttonColor = color.RGBA{DialogChoiceSelectedR, DialogChoiceSelectedG, DialogChoiceSelectedB, DialogChoiceSelectedA}
		} else {
			buttonColor = color.RGBA{DialogChoiceNormalR, DialogChoiceNormalG, DialogChoiceNormalB, DialogChoiceNormalA}
		}

		// Draw button background
		buttonWidth := dw.Width - DialogWidgetPadding*2
		vector.FillRect(screen,
			float32(buttonX), float32(buttonY),
			float32(buttonWidth), float32(DialogChoiceHeight),
			buttonColor, false)

		// Draw button text
		ebitenutil.DebugPrintAt(screen, choice.Text, buttonX+DialogChoicePadding, buttonY+6)
	}

	// Draw scroll indicator if needed
	if len(dw.AvailableChoices) > DialogMaxVisibleChoices {
		scrollY := choicesY + DialogMaxVisibleChoices*(DialogChoiceHeight+DialogChoiceSpacing) + 10
		scrollText := fmt.Sprintf("(%d/%d choices)", dw.ChoiceScrollOffset+DialogMaxVisibleChoices, len(dw.AvailableChoices))
		ebitenutil.DebugPrintAt(screen, scrollText, dw.X+DialogWidgetPadding, scrollY)
	}
}

// drawHelpText renders contextual help text
func (dw *DialogWidget) drawHelpText(screen *ebiten.Image) {
	helpY := dw.Y + dw.Height - 25
	var helpText string

	switch dw.State {
	case DialogStateTyping:
		helpText = "SPACE: Skip typing  TAB: Fast mode  ESC: Close"
	case DialogStateWaiting:
		helpText = "SPACE/ENTER: Continue  ESC: Close"
	case DialogStateChoices:
		helpText = "UP/DOWN: Select choice  ENTER: Choose  ESC: Close"
	default:
		helpText = "ESC: Close"
	}

	ebitenutil.DebugPrintAt(screen, helpText, dw.X+DialogWidgetPadding, helpY)
}

// GetVariable returns a dialog variable value
func (dw *DialogWidget) GetVariable(name string) (interface{}, bool) {
	value, exists := dw.Variables[name]
	return value, exists
}

// SetVariable sets a dialog variable value
func (dw *DialogWidget) SetVariable(name string, value interface{}) {
	dw.Variables[name] = value
	if dw.OnVariableSet != nil {
		dw.OnVariableSet(name, value)
	}
}
