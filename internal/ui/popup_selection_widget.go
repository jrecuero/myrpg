// Package ui provides user interface components including widgets for user interaction
package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PopupSelectionWidget represents a popup dialog with a scrollable list of options
type PopupSelectionWidget struct {
	// Configuration
	Title   string
	Options []string
	X, Y    int // Position of the popup
	Width   int
	Height  int

	// State
	IsVisible       bool
	SelectedIndex   int
	ScrollOffset    int
	MaxVisibleItems int

	// Callbacks
	OnSelection func(index int, option string) // Called when user presses Enter
	OnCancel    func()                         // Called when user presses Escape

	// Styling
	BackgroundColor color.RGBA
	BorderColor     color.RGBA
	TextColor       color.RGBA
	HighlightColor  color.RGBA
	TitleColor      color.RGBA
	ScrollbarColor  color.RGBA
	ShadowColor     color.RGBA

	// Internal state for input handling
	lastUpPressed     bool
	lastDownPressed   bool
	lastEnterPressed  bool
	lastEscapePressed bool
}

// NewPopupSelectionWidget creates a new popup selection widget
func NewPopupSelectionWidget(title string, options []string, x, y, width, height int) *PopupSelectionWidget {
	maxVisible := (height - 40) / 20 // Estimate based on font height + spacing
	if maxVisible < 1 {
		maxVisible = 1
	}

	return &PopupSelectionWidget{
		Title:           title,
		Options:         options,
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		IsVisible:       false,
		SelectedIndex:   0,
		ScrollOffset:    0,
		MaxVisibleItems: maxVisible,

		// Default styling
		BackgroundColor: color.RGBA{40, 40, 40, 240},    // Semi-transparent dark
		BorderColor:     color.RGBA{100, 100, 100, 255}, // Gray border
		TextColor:       color.RGBA{255, 255, 255, 255}, // White text
		HighlightColor:  color.RGBA{70, 130, 180, 200},  // Steel blue highlight
		TitleColor:      color.RGBA{255, 255, 0, 255},   // Yellow title
		ScrollbarColor:  color.RGBA{160, 160, 160, 255}, // Light gray scrollbar
		ShadowColor:     color.RGBA{0, 0, 0, 100},       // Semi-transparent shadow
	}
}

// Show displays the popup with the given options
func (p *PopupSelectionWidget) Show(title string, options []string) {
	p.Title = title
	p.Options = options
	p.IsVisible = true
	p.SelectedIndex = 0
	p.ScrollOffset = 0
	p.updateMaxVisibleItems()
	p.ensureSelectedVisible()
}

// Hide closes the popup
func (p *PopupSelectionWidget) Hide() {
	p.IsVisible = false
}

// Update handles input and updates the widget state
func (p *PopupSelectionWidget) Update() {
	if !p.IsVisible || len(p.Options) == 0 {
		return
	}

	// Handle Up/Down navigation
	upPressed := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	downPressed := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	enterPressed := ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace)
	escapePressed := ebiten.IsKeyPressed(ebiten.KeyEscape)

	// Navigate up (with key repeat prevention)
	if upPressed && !p.lastUpPressed {
		if p.SelectedIndex > 0 {
			p.SelectedIndex--
			p.ensureSelectedVisible()
		}
	}

	// Navigate down (with key repeat prevention)
	if downPressed && !p.lastDownPressed {
		if p.SelectedIndex < len(p.Options)-1 {
			p.SelectedIndex++
			p.ensureSelectedVisible()
		}
	}

	// Select current option
	if enterPressed && !p.lastEnterPressed {
		if p.OnSelection != nil {
			p.OnSelection(p.SelectedIndex, p.Options[p.SelectedIndex])
		}
		p.Hide()
	}

	// Cancel selection
	if escapePressed && !p.lastEscapePressed {
		if p.OnCancel != nil {
			p.OnCancel()
		}
		p.Hide()
	}

	// Update input state
	p.lastUpPressed = upPressed
	p.lastDownPressed = downPressed
	p.lastEnterPressed = enterPressed
	p.lastEscapePressed = escapePressed
}

// Draw renders the popup widget
func (p *PopupSelectionWidget) Draw(screen *ebiten.Image) {
	if !p.IsVisible {
		return
	}

	// Draw shadow
	shadowOffset := 3
	vector.FillRect(screen,
		float32(p.X+shadowOffset), float32(p.Y+shadowOffset),
		float32(p.Width), float32(p.Height),
		p.ShadowColor, false)

	// Draw background
	vector.FillRect(screen,
		float32(p.X), float32(p.Y),
		float32(p.Width), float32(p.Height),
		p.BackgroundColor, false)

	// Draw border
	vector.StrokeRect(screen,
		float32(p.X), float32(p.Y),
		float32(p.Width), float32(p.Height),
		2, p.BorderColor, false)

	// Draw title
	if p.Title != "" {
		ebitenutil.DebugPrintAt(screen, p.Title, p.X+10, p.Y+10)
	}

	// Calculate list area
	listStartY := p.Y + 30
	if p.Title == "" {
		listStartY = p.Y + 10
	}
	listHeight := p.Height - (listStartY - p.Y) - 10

	// Draw options
	p.drawOptions(screen, listStartY, listHeight)

	// Draw scrollbar if needed
	if len(p.Options) > p.MaxVisibleItems {
		p.drawScrollbar(screen, listStartY, listHeight)
	}

	// Draw instructions at bottom
	instructionY := p.Y + p.Height - 15
	ebitenutil.DebugPrintAt(screen, "↑↓: Navigate | Enter: Select | Esc: Cancel", p.X+10, instructionY)
}

// drawOptions renders the visible options list
func (p *PopupSelectionWidget) drawOptions(screen *ebiten.Image, startY, height int) {
	_ = height
	itemHeight := 20
	y := startY

	for i := 0; i < p.MaxVisibleItems && i+p.ScrollOffset < len(p.Options); i++ {
		optionIndex := i + p.ScrollOffset
		option := p.Options[optionIndex]

		// Highlight selected item
		if optionIndex == p.SelectedIndex {
			vector.FillRect(screen,
				float32(p.X+5), float32(y-2),
				float32(p.Width-10), float32(itemHeight),
				p.HighlightColor, false)
		}

		// Draw option text
		ebitenutil.DebugPrintAt(screen, option, p.X+10, y)

		y += itemHeight
	}
}

// drawScrollbar renders the scrollbar if the list is scrollable
func (p *PopupSelectionWidget) drawScrollbar(screen *ebiten.Image, startY, height int) {
	scrollbarX := p.X + p.Width - 15
	scrollbarWidth := 10
	scrollbarHeight := height - 20

	// Draw scrollbar background
	vector.FillRect(screen,
		float32(scrollbarX), float32(startY),
		float32(scrollbarWidth), float32(scrollbarHeight),
		color.RGBA{60, 60, 60, 255}, false)

	// Calculate thumb position and size
	totalItems := len(p.Options)
	thumbHeight := (p.MaxVisibleItems * scrollbarHeight) / totalItems
	if thumbHeight < 10 {
		thumbHeight = 10
	}

	thumbY := startY + (p.ScrollOffset*scrollbarHeight)/totalItems

	// Draw scrollbar thumb
	vector.FillRect(screen,
		float32(scrollbarX+1), float32(thumbY),
		float32(scrollbarWidth-2), float32(thumbHeight),
		p.ScrollbarColor, false)
}

// ensureSelectedVisible adjusts scroll offset to keep selected item visible
func (p *PopupSelectionWidget) ensureSelectedVisible() {
	if p.SelectedIndex < p.ScrollOffset {
		p.ScrollOffset = p.SelectedIndex
	} else if p.SelectedIndex >= p.ScrollOffset+p.MaxVisibleItems {
		p.ScrollOffset = p.SelectedIndex - p.MaxVisibleItems + 1
	}
}

// updateMaxVisibleItems recalculates how many items can be displayed
func (p *PopupSelectionWidget) updateMaxVisibleItems() {
	itemHeight := 20
	listHeight := p.Height - 60 // Account for title and instructions
	p.MaxVisibleItems = listHeight / itemHeight
	if p.MaxVisibleItems < 1 {
		p.MaxVisibleItems = 1
	}
}

// GetSelectedOption returns the currently selected option
func (p *PopupSelectionWidget) GetSelectedOption() (int, string) {
	if p.SelectedIndex >= 0 && p.SelectedIndex < len(p.Options) {
		return p.SelectedIndex, p.Options[p.SelectedIndex]
	}
	return -1, ""
}

// SetPosition updates the popup position
func (p *PopupSelectionWidget) SetPosition(x, y int) {
	p.X = x
	p.Y = y
}

// SetSize updates the popup dimensions
func (p *PopupSelectionWidget) SetSize(width, height int) {
	p.Width = width
	p.Height = height
	p.updateMaxVisibleItems()
}
