package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/skills"
)

// SkillNodeUI represents a visual skill node in the tree
type SkillNodeUI struct {
	Skill       *components.Skill
	X, Y        int
	Width       int
	Height      int
	IsLearned   bool
	IsAvailable bool
	IsHovered   bool
	IsSelected  bool
}

// SkillsWidget represents the skills/abilities interface
type SkillsWidget struct {
	// Widget properties
	X, Y          int
	Width, Height int
	Visible       bool
	Enabled       bool

	entity          *ecs.Entity
	skillsComponent *components.SkillsComponent
	currentJobTree  *components.SkillTree
	skillNodes      []*SkillNodeUI
	selectedSkill   *components.Skill
	showTooltip     bool
	tooltipX        int
	tooltipY        int

	// Layout constants
	nodeWidth    int
	nodeHeight   int
	nodeSpacingX int
	nodeSpacingY int
	treeOffsetX  int
	treeOffsetY  int

	// Colors
	colorLearned    color.Color
	colorAvailable  color.Color
	colorLocked     color.Color
	colorHovered    color.Color
	colorSelected   color.Color
	colorBackground color.Color
}

// NewSkillsWidget creates a new skills widget
func NewSkillsWidget(x, y, width, height int, entity *ecs.Entity) *SkillsWidget {
	widget := &SkillsWidget{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Visible: false,
		Enabled: true,

		entity:          entity,
		skillsComponent: entity.Skills(),
		skillNodes:      make([]*SkillNodeUI, 0),

		// Layout settings
		nodeWidth:    80,
		nodeHeight:   60,
		nodeSpacingX: 100,
		nodeSpacingY: 80,
		treeOffsetX:  50,
		treeOffsetY:  80,

		// Colors
		colorLearned:    color.RGBA{0, 255, 0, 255},     // Green
		colorAvailable:  color.RGBA{255, 255, 0, 255},   // Yellow
		colorLocked:     color.RGBA{128, 128, 128, 255}, // Gray
		colorHovered:    color.RGBA{255, 255, 255, 255}, // White
		colorSelected:   color.RGBA{0, 255, 255, 255},   // Cyan
		colorBackground: color.RGBA{20, 20, 20, 245},    // Much more opaque dark background
	}

	// Initialize with current job's skill tree if available
	if widget.skillsComponent != nil {
		if stats := entity.RPGStats(); stats != nil {
			registry := skills.GetGlobalSkillRegistry()
			if tree, exists := registry.GetSkillTree(stats.Job); exists {
				widget.currentJobTree = tree
				widget.buildSkillNodeUI()
			}
		}
	}

	return widget
}

// buildSkillNodeUI creates the visual skill nodes from the skill tree
func (sw *SkillsWidget) buildSkillNodeUI() {
	if sw.currentJobTree == nil {
		return
	}

	sw.skillNodes = make([]*SkillNodeUI, 0)

	for _, row := range sw.currentJobTree.Layout {
		for _, skillNode := range row {
			if skillNode == nil {
				continue
			}

			// Calculate position
			nodeX := sw.X + sw.treeOffsetX + (skillNode.X * sw.nodeSpacingX)
			nodeY := sw.Y + sw.treeOffsetY + (skillNode.Y * sw.nodeSpacingY)

			// Check if skill is learned
			isLearned := false
			if sw.skillsComponent != nil {
				_, isLearned = sw.skillsComponent.LearnedSkills[skillNode.Skill.ID]
			}

			// Check if skill is available to learn
			isAvailable := false
			if sw.skillsComponent != nil {
				isAvailable = sw.skillsComponent.CanLearnSkill(skillNode.Skill)
			}

			nodeUI := &SkillNodeUI{
				Skill:       skillNode.Skill,
				X:           nodeX,
				Y:           nodeY,
				Width:       sw.nodeWidth,
				Height:      sw.nodeHeight,
				IsLearned:   isLearned,
				IsAvailable: isAvailable,
			}

			sw.skillNodes = append(sw.skillNodes, nodeUI)
		}
	}
}

// Update handles input and updates the widget state
// Returns InputResult indicating what input was consumed
func (sw *SkillsWidget) Update() InputResult {
	result := NewInputResult()

	if !sw.Visible || !sw.Enabled {
		return result
	}

	// Handle ESC key to close skills widget
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sw.Visible = false
		result.EscConsumed = true
		return result
	}

	// Get mouse position
	mouseX, mouseY := ebiten.CursorPosition()

	// Reset hover states
	for _, node := range sw.skillNodes {
		node.IsHovered = false
	}

	// Check if mouse is over the widget area first
	isMouseOverWidget := mouseX >= sw.X && mouseX <= sw.X+sw.Width &&
		mouseY >= sw.Y && mouseY <= sw.Y+sw.Height

	if isMouseOverWidget {
		result.MouseConsumed = true

		// Check for node hover and clicks
		for _, node := range sw.skillNodes {
			if mouseX >= node.X && mouseX <= node.X+node.Width &&
				mouseY >= node.Y && mouseY <= node.Y+node.Height {

				node.IsHovered = true

				// Show tooltip
				sw.showTooltip = true
				sw.tooltipX = mouseX + 10
				sw.tooltipY = mouseY + 10
				sw.selectedSkill = node.Skill

				// Handle click to learn skill
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					if node.IsAvailable && !node.IsLearned {
						sw.learnSkill(node.Skill)
					}
				}
				break
			}
		}
	}

	// Hide tooltip if not hovering over any node
	if !sw.isHoveringAnyNode() {
		sw.showTooltip = false
		sw.selectedSkill = nil
	}

	return result
}

// isHoveringAnyNode checks if mouse is hovering over any skill node
func (sw *SkillsWidget) isHoveringAnyNode() bool {
	for _, node := range sw.skillNodes {
		if node.IsHovered {
			return true
		}
	}
	return false
}

// learnSkill attempts to learn a skill
func (sw *SkillsWidget) learnSkill(skill *components.Skill) {
	if sw.skillsComponent == nil {
		return
	}

	if sw.skillsComponent.LearnSkill(skill) {
		// Rebuild UI to reflect changes
		sw.buildSkillNodeUI()

		// Apply skill effects (this would be handled by a skill system)
		sw.applySkillEffects(skill)
	}
}

// applySkillEffects applies the effects of a learned skill to the character
func (sw *SkillsWidget) applySkillEffects(skill *components.Skill) {
	if sw.entity == nil {
		return
	}

	stats := sw.entity.RPGStats()
	if stats == nil {
		return
	}

	// Apply each effect
	for _, effect := range skill.Effects {
		switch effect.Type {
		case "stat_bonus":
			switch effect.Target {
			case "MaxHP":
				stats.MaxHP += effect.Value
				stats.CurrentHP += effect.Value // Also heal current HP
			case "MaxMP":
				stats.MaxMP += effect.Value
				stats.CurrentMP += effect.Value // Also restore current MP
			case "Attack":
				stats.Attack += effect.Value
			case "Defense":
				stats.Defense += effect.Value
			case "MagicAttack":
				// Handle magic attack if implemented
			}
		case "ability_unlock":
			// Active abilities would be handled by combat system
		case "passive_effect":
			// Passive effects would be handled by appropriate systems
		}
	}
}

// Draw renders the skills widget
func (sw *SkillsWidget) Draw(screen *ebiten.Image) {
	if !sw.Visible {
		return
	}

	// Draw background
	ebitenutil.DrawRect(screen, float64(sw.X), float64(sw.Y), float64(sw.Width), float64(sw.Height), sw.colorBackground)

	// Draw title
	title := "Skills & Abilities"
	if sw.currentJobTree != nil {
		title = sw.currentJobTree.Name
	}
	ebitenutil.DebugPrintAt(screen, title, sw.X+10, sw.Y+10)

	// Draw skill points info
	if sw.skillsComponent != nil {
		skillPointsText := fmt.Sprintf("Available Skill Points: %d", sw.skillsComponent.AvailablePoints)
		ebitenutil.DebugPrintAt(screen, skillPointsText, sw.X+10, sw.Y+30)
	}

	// Draw skill nodes
	for _, node := range sw.skillNodes {
		sw.drawSkillNode(screen, node)
	}

	// Draw connections between nodes
	sw.drawSkillConnections(screen)

	// Draw tooltip
	if sw.showTooltip && sw.selectedSkill != nil {
		sw.drawTooltip(screen, sw.selectedSkill, sw.tooltipX, sw.tooltipY)
	}
}

// drawSkillNode renders a single skill node
func (sw *SkillsWidget) drawSkillNode(screen *ebiten.Image, node *SkillNodeUI) {
	// Determine node color
	nodeColor := sw.colorLocked
	if node.IsLearned {
		nodeColor = sw.colorLearned
	} else if node.IsAvailable {
		nodeColor = sw.colorAvailable
	}

	// Add hover or selection effects
	if node.IsHovered {
		// Draw border for hover effect
		ebitenutil.DrawRect(screen, float64(node.X-2), float64(node.Y-2),
			float64(node.Width+4), float64(node.Height+4), sw.colorHovered)
	}

	// Draw node background
	ebitenutil.DrawRect(screen, float64(node.X), float64(node.Y),
		float64(node.Width), float64(node.Height), nodeColor)

	// Draw node border - thicker for better visibility
	borderColor := color.RGBA{50, 50, 50, 255}
	borderThickness := 3
	ebitenutil.DrawRect(screen, float64(node.X), float64(node.Y), float64(node.Width), float64(borderThickness), borderColor)
	ebitenutil.DrawRect(screen, float64(node.X), float64(node.Y), float64(borderThickness), float64(node.Height), borderColor)
	ebitenutil.DrawRect(screen, float64(node.X+node.Width-borderThickness), float64(node.Y), float64(borderThickness), float64(node.Height), borderColor)
	ebitenutil.DrawRect(screen, float64(node.X), float64(node.Y+node.Height-borderThickness), float64(node.Width), float64(borderThickness), borderColor)

	// Draw text background for better readability
	textBgColor := color.RGBA{0, 0, 0, 180}
	ebitenutil.DrawRect(screen, float64(node.X+2), float64(node.Y+2), float64(node.Width-4), 15, textBgColor)

	// Draw skill name (truncated to fit)
	skillName := node.Skill.Name
	if len(skillName) > 10 {
		skillName = skillName[:8] + ".."
	}
	ebitenutil.DebugPrintAt(screen, skillName, node.X+5, node.Y+5)

	// Draw tier number with background
	ebitenutil.DrawRect(screen, float64(node.X+2), float64(node.Y+node.Height-18), 20, 15, textBgColor)
	tierText := fmt.Sprintf("T%d", node.Skill.Tier)
	ebitenutil.DebugPrintAt(screen, tierText, node.X+5, node.Y+node.Height-15)

	// Draw cost with background
	ebitenutil.DrawRect(screen, float64(node.X+node.Width-27), float64(node.Y+node.Height-18), 25, 15, textBgColor)
	costText := fmt.Sprintf("%dSP", node.Skill.SkillPoints)
	ebitenutil.DebugPrintAt(screen, costText, node.X+node.Width-25, node.Y+node.Height-15)
}

// drawSkillConnections draws lines between prerequisite skills
func (sw *SkillsWidget) drawSkillConnections(screen *ebiten.Image) {
	if sw.currentJobTree == nil {
		return
	}

	lineColor := color.RGBA{100, 100, 100, 255}

	for _, node := range sw.skillNodes {
		// Draw lines to children
		for _, childID := range sw.currentJobTree.Nodes[node.Skill.ID].Children {
			if _, exists := sw.currentJobTree.Nodes[childID]; exists {
				// Find the corresponding UI node
				for _, childUINode := range sw.skillNodes {
					if childUINode.Skill.ID == childID {
						// Draw line from current node to child
						sw.drawLine(screen,
							node.X+node.Width/2, node.Y+node.Height,
							childUINode.X+childUINode.Width/2, childUINode.Y,
							lineColor)
						break
					}
				}
			}
		}
	}
}

// drawLine draws a simple line between two points
func (sw *SkillsWidget) drawLine(screen *ebiten.Image, x1, y1, x2, y2 int, lineColor color.Color) {
	// Simple vertical/horizontal line drawing
	if x1 == x2 {
		// Vertical line
		for y := y1; y <= y2; y++ {
			ebitenutil.DrawRect(screen, float64(x1), float64(y), 2, 1, lineColor)
		}
	} else if y1 == y2 {
		// Horizontal line
		minX, maxX := x1, x2
		if x1 > x2 {
			minX, maxX = x2, x1
		}
		for x := minX; x <= maxX; x++ {
			ebitenutil.DrawRect(screen, float64(x), float64(y1), 1, 2, lineColor)
		}
	}
}

// drawTooltip renders skill information tooltip
func (sw *SkillsWidget) drawTooltip(screen *ebiten.Image, skill *components.Skill, x, y int) {
	tooltipWidth := 250
	tooltipHeight := 120

	// Draw tooltip background - much more opaque for better text readability
	tooltipBg := color.RGBA{15, 15, 15, 250}
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(tooltipWidth), float64(tooltipHeight), tooltipBg)

	// Draw border
	borderColor := color.RGBA{220, 220, 220, 255}
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(tooltipWidth), 2, borderColor)
	ebitenutil.DrawRect(screen, float64(x), float64(y), 2, float64(tooltipHeight), borderColor)
	ebitenutil.DrawRect(screen, float64(x+tooltipWidth-2), float64(y), 2, float64(tooltipHeight), borderColor)
	ebitenutil.DrawRect(screen, float64(x), float64(y+tooltipHeight-2), float64(tooltipWidth), 2, borderColor)

	// Draw skill information
	lineHeight := 15
	currentY := y + 10

	ebitenutil.DebugPrintAt(screen, skill.Name, x+5, currentY)
	currentY += lineHeight

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Type: %s", skill.Type.String()), x+5, currentY)
	currentY += lineHeight

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cost: %d SP", skill.SkillPoints), x+5, currentY)
	currentY += lineHeight

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Tier: %d", skill.Tier), x+5, currentY)
	currentY += lineHeight

	// Draw description (wrap text if too long)
	desc := skill.Description
	if len(desc) > 35 {
		ebitenutil.DebugPrintAt(screen, desc[:35], x+5, currentY)
		currentY += lineHeight
		if len(desc) > 35 {
			remaining := desc[35:]
			if len(remaining) > 35 {
				remaining = remaining[:32] + "..."
			}
			ebitenutil.DebugPrintAt(screen, remaining, x+5, currentY)
		}
	} else {
		ebitenutil.DebugPrintAt(screen, desc, x+5, currentY)
	}
}

// Toggle shows/hides the skills widget
func (sw *SkillsWidget) Toggle() {
	sw.Visible = !sw.Visible
	if sw.Visible {
		// Refresh the skill tree display when opened
		sw.buildSkillNodeUI()
	}
}

// AddSkillPoints adds skill points to the character
func (sw *SkillsWidget) AddSkillPoints(points int) {
	if sw.skillsComponent != nil {
		sw.skillsComponent.AddSkillPoints(points)
		sw.buildSkillNodeUI() // Refresh availability
	}
}
