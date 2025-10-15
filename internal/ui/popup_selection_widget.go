// Package ui provides user interface components including widgets for user interaction
package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PopupSelectionWidget layout constants
const (
	// Layout dimensions
	SelectionWidgetItemHeight    = 20 // Height of each option item
	SelectionWidgetReservedSpace = 40 // Space reserved for title and padding
	SelectionWidgetTitleY        = 15 // Y offset for title from top
	SelectionWidgetContentY      = 35 // Y offset for content from top
	SelectionWidgetSidePadding   = 10 // Left/right padding for text
	SelectionWidgetItemPadding   = 2  // Padding between items
	SelectionWidgetBorderWidth   = 2  // Border thickness in pixels
	SelectionWidgetShadowOffset  = 4  // Shadow offset distance

	// Scrollbar constants
	SelectionWidgetScrollbarWidth    = 10 // Width of the scrollbar track
	SelectionWidgetScrollbarRight    = 15 // Distance from right edge to scrollbar
	SelectionWidgetScrollbarMargin   = 20 // Vertical margin for scrollbar area
	SelectionWidgetScrollbarMinThumb = 10 // Minimum scrollbar thumb height
	SelectionWidgetScrollbarPadding  = 1  // Internal padding for scrollbar thumb

	// Draw method layout constants
	SelectionWidgetTitlePadding     = 10 // Left padding for title text
	SelectionWidgetTitleMargin      = 10 // Top margin for title text
	SelectionWidgetNoTitleMargin    = 10 // Top margin when no title is present
	SelectionWidgetBottomMargin     = 10 // Bottom margin for instructions area
	SelectionWidgetInstrMargin      = 15 // Bottom margin for instruction text
	SelectionWidgetHighlightMargin  = 5  // Left margin for highlight rectangle
	SelectionWidgetHighlightPadding = 2  // Vertical padding for highlight rectangle
	SelectionWidgetTextPadding      = 10 // Left padding for option text
	SelectionWidgetListReserved     = 60 // Reserved space for title and instructions

	// Default colors (can be overridden per widget)
	SelectionWidgetDefaultBgR      = 40  // Background red
	SelectionWidgetDefaultBgG      = 40  // Background green
	SelectionWidgetDefaultBgB      = 40  // Background blue
	SelectionWidgetDefaultBgA      = 240 // Background alpha
	SelectionWidgetDefaultBorderR  = 100 // Border red
	SelectionWidgetDefaultBorderG  = 100 // Border green
	SelectionWidgetDefaultBorderB  = 100 // Border blue
	SelectionWidgetDefaultTextR    = 255 // Text red
	SelectionWidgetDefaultTextG    = 255 // Text green
	SelectionWidgetDefaultTextB    = 255 // Text blue
	SelectionWidgetDefaultTitleR   = 255 // Title red
	SelectionWidgetDefaultTitleG   = 255 // Title green
	SelectionWidgetDefaultTitleB   = 0   // Title blue (yellow title)
	SelectionWidgetScrollbarColorR = 160 // Scrollbar thumb red component
	SelectionWidgetScrollbarColorG = 160 // Scrollbar thumb green component
	SelectionWidgetScrollbarColorB = 160 // Scrollbar thumb blue component
	SelectionWidgetScrollbarColorA = 255 // Scrollbar thumb alpha component

	// Scrollbar background color constants
	SelectionWidgetScrollBgColorR = 60  // Scrollbar background red component
	SelectionWidgetScrollBgColorG = 60  // Scrollbar background green component
	SelectionWidgetScrollBgColorB = 60  // Scrollbar background blue component
	SelectionWidgetScrollBgColorA = 255 // Scrollbar background alpha component
	SelectionWidgetDefaultShadowA = 100 // Shadow alpha
	SelectionWidgetHighlightR     = 70  // Highlight red (steel blue)
	SelectionWidgetHighlightG     = 130 // Highlight green
	SelectionWidgetHighlightB     = 180 // Highlight blue
	SelectionWidgetHighlightA     = 200 // Highlight alpha
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
	maxVisible := (height - SelectionWidgetReservedSpace) / SelectionWidgetItemHeight
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
		BackgroundColor: color.RGBA{SelectionWidgetDefaultBgR, SelectionWidgetDefaultBgG, SelectionWidgetDefaultBgB, SelectionWidgetDefaultBgA},
		BorderColor:     color.RGBA{SelectionWidgetDefaultBorderR, SelectionWidgetDefaultBorderG, SelectionWidgetDefaultBorderB, 255},
		TextColor:       color.RGBA{SelectionWidgetDefaultTextR, SelectionWidgetDefaultTextG, SelectionWidgetDefaultTextB, 255},
		HighlightColor:  color.RGBA{SelectionWidgetHighlightR, SelectionWidgetHighlightG, SelectionWidgetHighlightB, SelectionWidgetHighlightA},
		TitleColor:      color.RGBA{SelectionWidgetDefaultTitleR, SelectionWidgetDefaultTitleG, SelectionWidgetDefaultTitleB, 255},
		ScrollbarColor:  color.RGBA{SelectionWidgetScrollbarColorR, SelectionWidgetScrollbarColorG, SelectionWidgetScrollbarColorB, SelectionWidgetScrollbarColorA},
		ShadowColor:     color.RGBA{0, 0, 0, SelectionWidgetDefaultShadowA},
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
	vector.FillRect(screen,
		float32(p.X+SelectionWidgetShadowOffset), float32(p.Y+SelectionWidgetShadowOffset),
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
		SelectionWidgetBorderWidth, p.BorderColor, false)

	// Draw title
	if p.Title != "" {
		ebitenutil.DebugPrintAt(screen, p.Title, p.X+SelectionWidgetTitlePadding, p.Y+SelectionWidgetTitleMargin)
	}

	// Calculate list area
	listStartY := p.Y + SelectionWidgetContentY
	if p.Title == "" {
		listStartY = p.Y + SelectionWidgetNoTitleMargin
	}
	listHeight := p.Height - (listStartY - p.Y) - SelectionWidgetBottomMargin

	// Draw options
	p.drawOptions(screen, listStartY, listHeight)

	// Draw scrollbar if needed
	if len(p.Options) > p.MaxVisibleItems {
		p.drawScrollbar(screen, listStartY, listHeight)
	}

	// Draw instructions at bottom
	instructionY := p.Y + p.Height - SelectionWidgetInstrMargin
	ebitenutil.DebugPrintAt(screen, "↑↓: Navigate | Enter: Select | Esc: Cancel", p.X+SelectionWidgetTitlePadding, instructionY)
}

// drawOptions renders the visible options list
func (p *PopupSelectionWidget) drawOptions(screen *ebiten.Image, startY, height int) {
	_ = height
	y := startY

	for i := 0; i < p.MaxVisibleItems && i+p.ScrollOffset < len(p.Options); i++ {
		optionIndex := i + p.ScrollOffset
		option := p.Options[optionIndex]

		// Highlight selected item
		if optionIndex == p.SelectedIndex {
			vector.FillRect(screen,
				float32(p.X+SelectionWidgetHighlightMargin), float32(y-SelectionWidgetHighlightPadding),
				float32(p.Width-SelectionWidgetHighlightMargin*2), float32(SelectionWidgetItemHeight),
				p.HighlightColor, false)
		}

		// Draw option text
		ebitenutil.DebugPrintAt(screen, option, p.X+SelectionWidgetTextPadding, y)

		y += SelectionWidgetItemHeight
	}
}

// drawScrollbar renders the scrollbar if the list is scrollable
func (p *PopupSelectionWidget) drawScrollbar(screen *ebiten.Image, startY, height int) {
	scrollbarX := p.X + p.Width - SelectionWidgetScrollbarRight
	scrollbarHeight := height - SelectionWidgetScrollbarMargin

	// Draw scrollbar background
	vector.FillRect(screen,
		float32(scrollbarX), float32(startY),
		float32(SelectionWidgetScrollbarWidth), float32(scrollbarHeight),
		color.RGBA{SelectionWidgetScrollBgColorR, SelectionWidgetScrollBgColorG, SelectionWidgetScrollBgColorB, SelectionWidgetScrollBgColorA}, false)

	// Calculate thumb position and size
	totalItems := len(p.Options)
	thumbHeight := (p.MaxVisibleItems * scrollbarHeight) / totalItems
	if thumbHeight < SelectionWidgetScrollbarMinThumb {
		thumbHeight = SelectionWidgetScrollbarMinThumb
	}

	thumbY := startY + (p.ScrollOffset*scrollbarHeight)/totalItems

	// Draw scrollbar thumb
	vector.FillRect(screen,
		float32(scrollbarX+SelectionWidgetScrollbarPadding), float32(thumbY),
		float32(SelectionWidgetScrollbarWidth-SelectionWidgetScrollbarPadding*2), float32(thumbHeight),
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
	listHeight := p.Height - SelectionWidgetListReserved // Account for title and instructions
	p.MaxVisibleItems = listHeight / SelectionWidgetItemHeight
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
