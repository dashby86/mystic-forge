package main

import (
	"database/sql"
	"fmt"
	"log"
	"mf/events"
	"mf/forge"
	"mf/game"
	sqlService "mf/services/sql"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Anvil struct {
	X, Y, Width, Height int
}

const (
	DBUser     = "app"       // Database username
	DBPassword = "admin"     // Database password
	DBHost     = "127.0.0.1" // Docker container hostname
	DBPort     = "3309"      // MySQL port
	DBName     = "app"       // Name of the database
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	sService := sqlService.SqlService{
		DB: db,
	}
	player, err := sService.GetPlayerByID(1)
	if err != nil {
		log.Fatalf("Failed to query the database: %v", err)
	}

	defer db.Close()
	// Check the connection
	err = db.Ping()

	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Initialize the game window
	ebiten.SetWindowSize(1024, 1024)
	ebiten.SetWindowTitle("Blacksmith Game")

	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		//widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		//widget.ContainerOpts.Bi.VisibleWhen)
	)

	// Set the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Load the background image
	backgroundImg, _, _ := ebitenutil.NewImageFromFile("assets/mainforge.png")

	// Implement game logic here
	spew.Dump(player)

	Forge := forge.Forge{
		Sql:    sService,
		Player: player,
	}

	// Create a new event queue
	eventQueue := events.NewEventQueue()

	// Create a channel to handle early cancellation
	stop := make(chan struct{})

	// Create a wait group to wait for goroutines to finish
	var wg sync.WaitGroup

	game := game.Game{
		Background: backgroundImg,
		Ui:         &ui,
		Sql:        sService,
		Player:     player,
		Forge:      Forge,
		Events:     eventQueue,
	}

	// Start the Dispatcher in a separate goroutine
	go eventQueue.Dispatcher(&wg, stop)

	game.CharWindow()
	game.Anvil()
	game.BattleButton()

	// Display the stats of different equipment types
	if err := ebiten.RunGame(&game); err != nil {
		close(stop) // Send stop signal to the Dispatcher
		fmt.Println(err)
	}
}
