// Package ui provides user interface components including widgets for user interaction
package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PopupInfoWidget layout constants
const (
	// Layout dimensions
	InfoWidgetLineHeight    = 16 // Height of each text line in pixels
	InfoWidgetTitleY        = 20 // Y offset for title from top
	InfoWidgetTitleSpacing  = 30 // Space between title and content
	InfoWidgetBottomPadding = 30 // Space reserved for help text at bottom
	InfoWidgetSidePadding   = 10 // Left/right padding for text
	InfoWidgetExtraPadding  = 10 // Additional padding for layout
	InfoWidgetBorderWidth   = 2  // Border thickness
	InfoWidgetShadowOffset  = 4  // Shadow offset distance

	// Calculated total reserved space
	InfoWidgetReservedSpace = InfoWidgetTitleY + InfoWidgetTitleSpacing + InfoWidgetBottomPadding + InfoWidgetExtraPadding // 90px total

	// Scrollbar dimensions
	InfoWidgetScrollbarWidth = 8  // Width of scrollbar
	InfoWidgetScrollbarRight = 15 // Distance from right edge
	InfoWidgetMinThumbSize   = 20 // Minimum scrollbar thumb height

	// Default colors (can be overridden per widget)
	InfoWidgetDefaultBgR     = 30  // Background red
	InfoWidgetDefaultBgG     = 30  // Background green
	InfoWidgetDefaultBgB     = 30  // Background blue
	InfoWidgetDefaultBgA     = 240 // Background alpha
	InfoWidgetDefaultBorderR = 100 // Border red
	InfoWidgetDefaultBorderG = 100 // Border green
	InfoWidgetDefaultBorderB = 100 // Border blue
	InfoWidgetDefaultTextR   = 255 // Text red
	InfoWidgetDefaultTextG   = 255 // Text green
	InfoWidgetDefaultTextB   = 255 // Text blue
	InfoWidgetDefaultTitleR  = 255 // Title red
	InfoWidgetDefaultTitleG  = 255 // Title green
	InfoWidgetDefaultTitleB  = 0   // Title blue (yellow title)
	InfoWidgetDefaultScrollR = 160 // Scrollbar red
	InfoWidgetDefaultScrollG = 160 // Scrollbar green
	InfoWidgetDefaultScrollB = 160 // Scrollbar blue
	InfoWidgetDefaultShadowA = 100 // Shadow alpha
	InfoWidgetScrollBgR      = 60  // Scrollbar background red
	InfoWidgetScrollBgG      = 60  // Scrollbar background green
	InfoWidgetScrollBgB      = 60  // Scrollbar background blue
)

// PopupInfoWidget represents a popup dialog for displaying plain information/text
type PopupInfoWidget struct {
	// Configuration
	Title   string
	Content string // Main text content to display
	X, Y    int    // Position of the popup
	Width   int
	Height  int

	// State
	IsVisible    bool
	ScrollOffset int // For scrolling through long content
	LineHeight   int // Height of each text line
	MaxLines     int // Maximum lines that fit in the widget

	// Callbacks
	OnClose func() // Called when user presses Escape

	// Styling
	BackgroundColor color.RGBA
	BorderColor     color.RGBA
	TextColor       color.RGBA
	TitleColor      color.RGBA
	ScrollbarColor  color.RGBA
	ShadowColor     color.RGBA

	// Internal state for input handling
	lastEscapePressed bool
	lastUpPressed     bool
	lastDownPressed   bool

	// Processed content
	lines []string // Content split into lines for display
}

// NewPopupInfoWidget creates a new popup information widget
func NewPopupInfoWidget(title string, content string, x, y, width, height int) *PopupInfoWidget {
	maxLines := (height - InfoWidgetReservedSpace) / InfoWidgetLineHeight
	if maxLines < 1 {
		maxLines = 1
	}

	widget := &PopupInfoWidget{
		Title:        title,
		Content:      content,
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
		IsVisible:    false,
		ScrollOffset: 0,
		LineHeight:   InfoWidgetLineHeight,
		MaxLines:     maxLines,

		// Default styling
		BackgroundColor: color.RGBA{InfoWidgetDefaultBgR, InfoWidgetDefaultBgG, InfoWidgetDefaultBgB, InfoWidgetDefaultBgA},
		BorderColor:     color.RGBA{InfoWidgetDefaultBorderR, InfoWidgetDefaultBorderG, InfoWidgetDefaultBorderB, 255},
		TextColor:       color.RGBA{InfoWidgetDefaultTextR, InfoWidgetDefaultTextG, InfoWidgetDefaultTextB, 255},
		TitleColor:      color.RGBA{InfoWidgetDefaultTitleR, InfoWidgetDefaultTitleG, InfoWidgetDefaultTitleB, 255},
		ScrollbarColor:  color.RGBA{InfoWidgetDefaultScrollR, InfoWidgetDefaultScrollG, InfoWidgetDefaultScrollB, 255},
		ShadowColor:     color.RGBA{0, 0, 0, InfoWidgetDefaultShadowA},
	}

	widget.processContent()
	return widget
}

// Show displays the popup with the given title and content
func (p *PopupInfoWidget) Show(title string, content string) {
	p.Title = title
	p.Content = content
	p.IsVisible = true
	p.ScrollOffset = 0
	p.processContent()
}

// Hide closes the popup
func (p *PopupInfoWidget) Hide() {
	p.IsVisible = false
	if p.OnClose != nil {
		p.OnClose()
	}
}

// processContent splits the content into lines that fit within the widget width
func (p *PopupInfoWidget) processContent() {
	if p.Content == "" {
		p.lines = []string{}
		return
	}

	// For now, simple line splitting by newlines
	// TODO: Add word wrapping based on widget width
	p.lines = strings.Split(p.Content, "\n")
}

// Update handles input and updates the widget state
func (p *PopupInfoWidget) Update() {
	if !p.IsVisible {
		return
	}

	// Handle Escape to close
	escapePressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if escapePressed && !p.lastEscapePressed {
		p.Hide()
	}
	p.lastEscapePressed = escapePressed

	// Handle scrolling if content is longer than widget
	if len(p.lines) > p.MaxLines {
		// Handle Up/Down scrolling
		upPressed := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
		downPressed := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)

		if upPressed && !p.lastUpPressed {
			if p.ScrollOffset > 0 {
				p.ScrollOffset--
			}
		}

		if downPressed && !p.lastDownPressed {
			maxScroll := len(p.lines) - p.MaxLines
			if p.ScrollOffset < maxScroll {
				p.ScrollOffset++
			}
		}

		p.lastUpPressed = upPressed
		p.lastDownPressed = downPressed
	}
}

// Draw renders the popup widget to the screen
func (p *PopupInfoWidget) Draw(screen *ebiten.Image) {
	if !p.IsVisible {
		return
	}

	// Draw shadow
	vector.FillRect(screen,
		float32(p.X+InfoWidgetShadowOffset), float32(p.Y+InfoWidgetShadowOffset),
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
		InfoWidgetBorderWidth, p.BorderColor, false)

	// Draw title
	titleY := p.Y + InfoWidgetTitleY
	ebitenutil.DebugPrintAt(screen, p.Title, p.X+InfoWidgetSidePadding, titleY)

	// Draw content lines
	contentStartY := titleY + InfoWidgetTitleSpacing
	visibleLines := p.MaxLines
	if len(p.lines) < visibleLines {
		visibleLines = len(p.lines)
	}

	for i := 0; i < visibleLines; i++ {
		lineIndex := i + p.ScrollOffset
		if lineIndex >= len(p.lines) {
			break
		}

		line := p.lines[lineIndex]
		lineY := contentStartY + (i * p.LineHeight)
		ebitenutil.DebugPrintAt(screen, line, p.X+InfoWidgetSidePadding, lineY)
	}

	// Draw scrollbar if needed
	if len(p.lines) > p.MaxLines {
		p.drawScrollbar(screen, contentStartY)
	}

	// Draw help text at bottom
	helpText := "Press ESC to close"
	if len(p.lines) > p.MaxLines {
		helpText += " | ↑↓ to scroll"
	}
	helpY := p.Y + p.Height - InfoWidgetTitleY // Use same spacing as title
	ebitenutil.DebugPrintAt(screen, helpText, p.X+InfoWidgetSidePadding, helpY)
}

// drawScrollbar renders a scrollbar indicating current scroll position
func (p *PopupInfoWidget) drawScrollbar(screen *ebiten.Image, contentStartY int) {
	scrollbarX := float32(p.X + p.Width - InfoWidgetScrollbarRight)
	scrollbarY := float32(contentStartY)
	scrollbarWidth := float32(InfoWidgetScrollbarWidth)
	scrollbarHeight := float32(p.MaxLines * p.LineHeight)

	// Draw scrollbar background
	vector.FillRect(screen, scrollbarX, scrollbarY, scrollbarWidth, scrollbarHeight,
		color.RGBA{InfoWidgetScrollBgR, InfoWidgetScrollBgG, InfoWidgetScrollBgB, 255}, false)

	// Calculate scrollbar thumb position and size
	totalLines := len(p.lines)
	thumbHeight := (float32(p.MaxLines) / float32(totalLines)) * scrollbarHeight
	if thumbHeight < InfoWidgetMinThumbSize {
		thumbHeight = InfoWidgetMinThumbSize
	}

	scrollProgress := float32(p.ScrollOffset) / float32(totalLines-p.MaxLines)
	thumbY := scrollbarY + scrollProgress*(scrollbarHeight-thumbHeight)

	// Draw scrollbar thumb
	vector.FillRect(screen, scrollbarX, thumbY, scrollbarWidth, thumbHeight,
		p.ScrollbarColor, false)
}

// SetPosition updates the widget's position
func (p *PopupInfoWidget) SetPosition(x, y int) {
	p.X = x
	p.Y = y
}

// SetSize updates the widget's size and recalculates layout
func (p *PopupInfoWidget) SetSize(width, height int) {
	p.Width = width
	p.Height = height
	p.MaxLines = (height - InfoWidgetReservedSpace) / p.LineHeight
	if p.MaxLines < 1 {
		p.MaxLines = 1
	}
}

// SetContent updates the displayed content
func (p *PopupInfoWidget) SetContent(title string, content string) {
	p.Title = title
	p.Content = content
	p.ScrollOffset = 0
	p.processContent()
}
