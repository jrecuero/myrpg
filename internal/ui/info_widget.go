// Package ui provides UI components for displaying event-related information.
package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// InfoWidget displays information to the player (chest contents, signs, etc.)
type InfoWidget struct {
	X, Y          int
	Width, Height int
	Visible       bool
	Title         string
	Message       string
	ImagePath     string

	// Visual properties
	backgroundColor color.Color
	borderColor     color.Color
	textColor       color.Color

	// Input handling
	closeRequested bool
}

// NewInfoWidget creates a new information display widget
func NewInfoWidget(x, y, width, height int) *InfoWidget {
	return &InfoWidget{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		Visible:         false,
		backgroundColor: color.RGBA{40, 40, 60, 240},    // Semi-transparent dark blue
		borderColor:     color.RGBA{100, 150, 200, 255}, // Light blue border
		textColor:       color.RGBA{255, 255, 255, 255}, // White text
		closeRequested:  false,
	}
}

// Show displays the info widget with the specified content
func (iw *InfoWidget) Show(title, message, imagePath string) {
	iw.Title = title
	iw.Message = message
	iw.ImagePath = imagePath
	iw.Visible = true
	iw.closeRequested = false
}

// Hide hides the info widget
func (iw *InfoWidget) Hide() {
	iw.Visible = false
	iw.closeRequested = false
}

// Update handles input for the info widget
func (iw *InfoWidget) Update() InputResult {
	if !iw.Visible {
		return NewInputResult()
	}

	result := NewInputResult()

	// Check for close input (Escape, Enter, or Space)
	if ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		ebiten.IsKeyPressed(ebiten.KeyEnter) ||
		ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !iw.closeRequested {
			iw.closeRequested = true
			iw.Hide()
			result.EscConsumed = true
			return result
		}
	} else {
		iw.closeRequested = false
	}

	// Check for mouse click to close
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()

		// If clicked inside widget area, close it
		if mouseX >= iw.X && mouseX <= iw.X+iw.Width &&
			mouseY >= iw.Y && mouseY <= iw.Y+iw.Height {
			iw.Hide()
			result.MouseConsumed = true
			return result
		}
	}

	// Consume all input while visible to prevent interaction with game
	result.EscConsumed = true
	result.MouseConsumed = true
	return result
}

// Draw renders the info widget
func (iw *InfoWidget) Draw(screen *ebiten.Image) {
	if !iw.Visible {
		return
	}

	// Draw background
	ebitenutil.DrawRect(screen,
		float64(iw.X), float64(iw.Y),
		float64(iw.Width), float64(iw.Height),
		iw.backgroundColor)

	// Draw border
	borderWidth := 2
	// Top border
	ebitenutil.DrawRect(screen,
		float64(iw.X), float64(iw.Y),
		float64(iw.Width), float64(borderWidth),
		iw.borderColor)
	// Bottom border
	ebitenutil.DrawRect(screen,
		float64(iw.X), float64(iw.Y+iw.Height-borderWidth),
		float64(iw.Width), float64(borderWidth),
		iw.borderColor)
	// Left border
	ebitenutil.DrawRect(screen,
		float64(iw.X), float64(iw.Y),
		float64(borderWidth), float64(iw.Height),
		iw.borderColor)
	// Right border
	ebitenutil.DrawRect(screen,
		float64(iw.X+iw.Width-borderWidth), float64(iw.Y),
		float64(borderWidth), float64(iw.Height),
		iw.borderColor)

	// Draw title
	if iw.Title != "" {
		titleY := iw.Y + 15
		ebitenutil.DebugPrintAt(screen, iw.Title, iw.X+10, titleY)

		// Draw underline for title
		titleLineY := titleY + 15
		ebitenutil.DrawRect(screen,
			float64(iw.X+10), float64(titleLineY),
			float64(len(iw.Title)*6), 1,
			iw.textColor)
	}

	// Draw message content
	if iw.Message != "" {
		messageY := iw.Y + 45
		messageLines := iw.wrapText(iw.Message, (iw.Width-20)/8) // Approximate character width

		for i, line := range messageLines {
			if messageY+(i*15) < iw.Y+iw.Height-30 { // Leave space at bottom
				ebitenutil.DebugPrintAt(screen, line, iw.X+10, messageY+(i*15))
			}
		}
	}

	// Draw instructions at the bottom
	instructionY := iw.Y + iw.Height - 20
	ebitenutil.DebugPrintAt(screen, "Press ESC, ENTER, or SPACE to close", iw.X+10, instructionY)
}

// wrapText wraps text to fit within a specified character width
func (iw *InfoWidget) wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := make([]string, 0)
	currentWord := ""

	// Split text into words, handling newlines
	for _, char := range text {
		if char == ' ' || char == '\n' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
			if char == '\n' {
				words = append(words, "\n")
			}
		} else {
			currentWord += string(char)
		}
	}
	if currentWord != "" {
		words = append(words, currentWord)
	}

	lines := make([]string, 0)
	currentLine := ""

	for _, word := range words {
		if word == "\n" {
			lines = append(lines, currentLine)
			currentLine = ""
			continue
		}

		testLine := currentLine
		if currentLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= maxWidth {
			currentLine = testLine
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

// IsVisible returns whether the widget is currently visible
func (iw *InfoWidget) IsVisible() bool {
	return iw.Visible
}

// SetPosition sets the widget position
func (iw *InfoWidget) SetPosition(x, y int) {
	iw.X = x
	iw.Y = y
}

// SetSize sets the widget size
func (iw *InfoWidget) SetSize(width, height int) {
	iw.Width = width
	iw.Height = height
}

// CenterOn centers the widget on the specified position
func (iw *InfoWidget) CenterOn(screenWidth, screenHeight int) {
	iw.X = (screenWidth - iw.Width) / 2
	iw.Y = (screenHeight - iw.Height) / 2
}
