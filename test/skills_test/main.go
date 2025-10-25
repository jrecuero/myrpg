package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/skills"
	"github.com/jrecuero/myrpg/internal/ui"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Game represents the skills test application
type Game struct {
	skillsWidget *ui.SkillsWidget
	testEntity   *ecs.Entity
	showHelp     bool
}

// NewGame creates a new skills test game
func NewGame() *Game {
	// Initialize skill registry
	skills.InitializeSkillRegistry()

	// Create test entity with warrior job
	entity := ecs.NewEntity("TestWarrior")
	
	// Add RPG stats component
	rpgStats := components.NewRPGStatsComponent("Test Warrior", components.JobWarrior, 1)
	entity.AddComponent(ecs.ComponentRPGStats, rpgStats)

	// Add skills component
	skillsComp := components.NewSkillsComponent(components.JobWarrior)
	skillsComp.AddSkillPoints(10) // Give some skill points to spend
	entity.AddComponent(ecs.ComponentSkills, skillsComp)

	// Initialize skill tree for the entity
	registry := skills.GetGlobalSkillRegistry()
	if tree, exists := registry.GetSkillTree(components.JobWarrior); exists {
		skillsComp.SkillTrees[components.JobWarrior] = tree
	}

	// Create skills widget
	skillsWidget := ui.NewSkillsWidget(50, 50, 700, 500, entity)
	skillsWidget.Toggle() // Show it by default

	return &Game{
		skillsWidget: skillsWidget,
		testEntity:   entity,
		showHelp:     true,
	}
}

// Update handles game logic updates
func (g *Game) Update() error {
	// Handle help toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.showHelp = !g.showHelp
	}

	// Handle skills widget toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.skillsWidget.Toggle()
	}

	// Handle adding skill points for testing
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.skillsWidget.AddSkillPoints(5)
	}

	// Handle job class switching for testing
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.switchJobClass(components.JobWarrior)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.switchJobClass(components.JobMage)
	}

	// Handle quit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Update skills widget
	g.skillsWidget.Update()

	return nil
}

// switchJobClass changes the test entity's job class and updates the skill tree
func (g *Game) switchJobClass(newJob components.JobType) {
	if stats := g.testEntity.RPGStats(); stats != nil {
		stats.Job = newJob
		
		// Create new skills component for the new job
		skillsComp := components.NewSkillsComponent(newJob)
		skillsComp.AddSkillPoints(10) // Give skill points
		g.testEntity.AddComponent(ecs.ComponentSkills, skillsComp)

		// Update skill tree
		registry := skills.GetGlobalSkillRegistry()
		if tree, exists := registry.GetSkillTree(newJob); exists {
			skillsComp.SkillTrees[newJob] = tree
		}

		// Create new skills widget
		g.skillsWidget = ui.NewSkillsWidget(50, 50, 700, 500, g.testEntity)
		g.skillsWidget.Toggle() // Show it
	}
}

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{26, 26, 26, 255})

	// Draw title
	ebitenutil.DebugPrintAt(screen, "Skills System Test", 10, 10)

	// Draw current character info
	if stats := g.testEntity.RPGStats(); stats != nil {
		info := fmt.Sprintf("Character: %s | Job: %s | Level: %d", 
			stats.Name, stats.Job.String(), stats.Level)
		ebitenutil.DebugPrintAt(screen, info, 10, 30)
	}

	// Draw skills widget
	g.skillsWidget.Draw(screen)

	// Draw help text
	if g.showHelp {
		g.drawHelp(screen)
	}
}

// drawHelp renders help text
func (g *Game) drawHelp(screen *ebiten.Image) {
	helpY := ScreenHeight - 120
	helpLines := []string{
		"Controls:",
		"H - Toggle help",
		"S - Toggle skills window",
		"P - Add 5 skill points",
		"1 - Switch to Warrior",
		"2 - Switch to Mage",
		"Click skills to learn them",
		"ESC - Quit",
	}

	for i, line := range helpLines {
		ebitenutil.DebugPrintAt(screen, line, 10, helpY+(i*15))
	}
}

// Layout sets the game screen layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Skills System Test")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}