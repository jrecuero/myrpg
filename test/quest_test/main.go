package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jrecuero/myrpg/internal/ecs"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/logger"
	"github.com/jrecuero/myrpg/internal/quests"
	"github.com/jrecuero/myrpg/internal/ui"
)

// QuestTestGame demonstrates the quest journal system
type QuestTestGame struct {
	player      *ecs.Entity
	questWidget *ui.QuestJournalWidget
	initialized bool
}

// Update handles the game logic and input
func (g *QuestTestGame) Update() error {
	if !g.initialized {
		g.initialize()
		g.initialized = true
	}

	// Update the quest widget
	if g.questWidget != nil {
		g.questWidget.Update()
	}

	return nil
}

// initialize sets up the test environment
func (g *QuestTestGame) initialize() {
	fmt.Println("Quest Journal System Test")
	fmt.Println("========================")

	// Initialize quest registry
	quests.InitializeQuestRegistry()
	registry := quests.GetGlobalQuestRegistry()

	// Create a test player
	g.player = ecs.NewEntity("Test Player")
	g.player.AddComponent(ecs.ComponentRPGStats, &components.RPGStatsComponent{
		Name:  "Hero",
		Job:   components.JobWarrior,
		Level: 1,
	})

	// Create quest journal component
	questJournal := components.NewQuestJournalComponent()

	// Add some test quests
	tutorialQuest, _ := registry.CreateQuestInstance("tutorial_first_steps")
	mainQuest, _ := registry.CreateQuestInstance("main_threat_emerges")
	sideQuest, _ := registry.CreateQuestInstance("side_gather_resources")
	equipQuest, _ := registry.CreateQuestInstance("side_equipment_mastery")

	// Start some quests
	tutorialQuest.StartQuest()
	questJournal.AddQuest(tutorialQuest)

	mainQuest.StartQuest()
	questJournal.AddQuest(mainQuest)

	sideQuest.StartQuest()
	questJournal.AddQuest(sideQuest)

	equipQuest.StartQuest()
	questJournal.AddQuest(equipQuest)

	// Complete tutorial quest for demo
	tutorialQuest.UpdateObjective("tutorial_move", 1)
	tutorialQuest.UpdateObjective("tutorial_combat", 2)
	if tutorialQuest.CheckCompletion() {
		questJournal.CompleteQuest("tutorial_first_steps")
		fmt.Println("Tutorial quest completed!")
	}

	// Make progress on main quest
	mainQuest.UpdateObjective("eliminate_threats", 3)
	fmt.Println("Progress on main quest: 3/5 enemies defeated")

	// Make progress on side quest
	sideQuest.UpdateObjective("collect_iron", 2)
	sideQuest.UpdateObjective("collect_wood", 3)
	fmt.Println("Progress on side quest: 2/3 iron, 3/5 wood")

	g.player.AddComponent(ecs.ComponentQuestJournal, questJournal)

	// Create quest widget
	g.questWidget = ui.NewQuestJournalWidget(50, 50, 800, 600, questJournal)
	g.questWidget.Visible = true

	fmt.Println("\nQuest Journal Widget Created!")
	fmt.Printf("Active Quests: %d\n", len(questJournal.GetActiveQuests()))
	fmt.Printf("Completed Quests: %d\n", len(questJournal.GetCompletedQuests()))

	fmt.Println("\nActive Quests:")
	for _, quest := range questJournal.GetActiveQuests() {
		fmt.Printf("- %s (%s, Level %d)\n", quest.Title, getQuestTypeText(quest.Type), quest.Level)
		for _, obj := range quest.Objectives {
			fmt.Printf("  %s\n", obj.GetProgressText())
		}
	}

	fmt.Println("\nCompleted Quests:")
	for _, quest := range questJournal.GetCompletedQuests() {
		fmt.Printf("- %s (%s, Level %d) ✓\n", quest.Title, getQuestTypeText(quest.Type), quest.Level)
	}

	fmt.Println("\nControls:")
	fmt.Println("- TAB: Switch between Active/Completed tabs")
	fmt.Println("- Arrow Keys: Navigate quests")
	fmt.Println("- Enter: Show/hide quest details")
	fmt.Println("- ESC: Close quest journal")
}

// getQuestTypeText returns display text for quest types
func getQuestTypeText(questType components.QuestType) string {
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

// Draw renders the game
func (g *QuestTestGame) Draw(screen *ebiten.Image) {
	if g.questWidget != nil {
		g.questWidget.Draw(screen)
	}
}

// Layout defines the screen size
func (g *QuestTestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1024, 768
}

func main() {
	// Initialize logger
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Set window properties
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Quest Journal System Test")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Create and run the game
	game := &QuestTestGame{}

	fmt.Println("Starting Quest Journal Test...")
	fmt.Println("The quest journal widget will appear in the game window.")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
