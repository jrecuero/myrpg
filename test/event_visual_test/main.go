package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jrecuero/myrpg/cmd/myrpg/game/entities"
	"github.com/jrecuero/myrpg/internal/ecs/components"
	"github.com/jrecuero/myrpg/internal/engine"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// TestGame for testing event system visualization
type TestGame struct {
	game *engine.Game
}

func main() {
	fmt.Println("Starting Event System Visual Test...")

	game := engine.NewGame()

	// Create player
	player := entities.CreatePlayerWithJob("Hero", 100, 200, components.JobWarrior, 1)
	game.AddEntity(player)

	// Create various event entities to test visualization

	// Battle event (red square)
	battleEvent := entities.CreateBattleEvent("battle_test", "Test Battle", 200, 150, []string{"enemy1"}, "test_map")
	game.AddEntity(battleEvent)

	// Chest event (yellow square)
	chestEvent := entities.CreateChestEvent("chest_test", "Test Chest", 300, 150, []string{"sword"}, 100, false)
	game.AddEntity(chestEvent)

	// Dialog event (blue square)
	npcEvent := entities.CreateDialogEvent("npc_test", "Test NPC", 400, 150, "npc1", "dialog1")
	game.AddEntity(npcEvent)

	// Info event (green square)
	signEvent := entities.CreateInfoEvent("sign_test", "Test Sign", 200, 250, "Test Sign", "This is a test sign.")
	game.AddEntity(signEvent)

	// Hidden trap (should not be visible)
	trapEvent := entities.CreateTrapEvent("trap_test", "Test Trap", 300, 250, "Spike", 10)
	game.AddEntity(trapEvent)

	// Set up event handlers
	game.SetupGameEventHandlers()

	// Initialize game
	game.InitializeGame()

	testGame := &TestGame{game: game}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Event System Visual Test")

	fmt.Println("Event entities created:")
	fmt.Println("- Red square: Battle Event (touch to trigger)")
	fmt.Println("- Yellow square: Chest Event (touch to trigger)")
	fmt.Println("- Blue square: Dialog Event (proximity trigger)")
	fmt.Println("- Green square: Info Event (touch to trigger)")
	fmt.Println("- Hidden Trap: Not visible (touch to trigger)")
	fmt.Println("")
	fmt.Println("Move the player (arrow keys) to interact with events!")

	if err := ebiten.RunGame(testGame); err != nil {
		log.Fatal(err)
	}
}

func (tg *TestGame) Update() error {
	return tg.game.Update()
}

func (tg *TestGame) Draw(screen *ebiten.Image) {
	tg.game.Draw(screen)

	// Add debug information
	ebitenutil.DebugPrint(screen, "Event System Visual Test\nArrow Keys: Move Player\nTouch colored squares to trigger events!")
}

func (tg *TestGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
