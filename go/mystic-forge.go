package main

import (
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	myForge "mf/forge"
	sqlService "mf/services/sql"
)

type Forge struct {
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
	player, err := sService.GetPlayerByID()
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
	backgroundImg, _, _ := ebitenutil.NewImageFromFile("assets/forge-main.png")

	// Define the forge
	forge := Forge{X: 200, Y: 200, Width: 200, Height: 200}

	// Implement game logic here
	fmt.Println("Welcome to the blacksmith game! Player: ", player.Name)

	game := game{
		Background: backgroundImg,
		Forge:      forge,
		ui:         &ui,
	}
	// Display the stats of different equipment types
	if err := ebiten.RunGame(&game); err != nil {
		fmt.Println(err)
	}
}

type game struct {
	Background *ebiten.Image
	Forge      Forge
	Crafted    bool
	ui         *ebitenui.UI
}

func (g *game) Update() error {
	// Check if the forge has been clicked and if a craft has not already been triggered
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !g.Crafted {
		x, y := ebiten.CursorPosition()
		forge := g.Forge
		if x >= forge.X && x <= forge.X+forge.Width && y >= forge.Y && y <= forge.Y+forge.Height {
			g.ShowCraftMenu()
			// Craft equipment
			fmt.Println("Crafting equipment...")
			//gear := models.Gear{}
			gear := myForge.CraftGear()
			spew.Dump(gear)
			//gear := forge.CraftGear
			fmt.Println("hp:", gear.Level)
			g.Crafted = true
		}
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Reset the "Crafted" flag when the mouse button is released
		g.Crafted = false
	}
	//g.ui.Update()

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Draw the background image
	screen.DrawImage(g.Background, nil)

	// Draw the forge
	forge := g.Forge
	ebitenutil.DrawRect(screen, float64(forge.X), float64(forge.Y), float64(forge.Width), float64(forge.Height), color.NRGBA{0, 0, 0, 255})
	g.ui.Draw(screen)
	//g.ui.
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Set the screen size
	return 1024, 1024
}

func (g *game) ShowCraftMenu() {
	// Check if the forge has been clicked
	// Create the container
	myImage, _, err := ebitenutil.NewImageFromFile("assets/menu-frame.png")
	if err != nil {
		log.Fatal(err)
	}
	buttonSlice, _ := loadButtonImage()
	nineSlice := image.NewNineSlice(myImage, [3]int{310, 310, 310}, [3]int{270, 270, 270})
	innerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
				StretchVertical:    true,
			}),
		),
	)

	buttonStackedLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		// instruct the container's anchor layout to center the button both horizontally and vertically;
		// since our button is a 2-widget object, we add the anchor info to the wrapping container
		// instead of the button
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	)
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonSlice),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			//g.ui.Container.RemoveChildren()
			println("button clicked")
		}),
	)
	//g.ui.Container.BackgroundImage.Draw(nineSlice)
	buttonStackedLayout.AddChild(button)
	//buttonStackedLayout.AddChild(widget.NewGraphic(widget.GraphicOpts.Image(buttonIcon)))
	innerContainer.AddChild(buttonStackedLayout)
	g.ui.Container.AddChild(button)
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func loadButtonImage() (*widget.ButtonImage, error) {
	/**
	buttonImage, _, err := ebitenutil.NewImageFromFile("assets/sell_menu.png")
	if err != nil {
		log.Fatal(err)
	}
	//buttonSlice := image.NewNineSlice(buttonImage, [3]int{310, 310, 310}, [3]int{270, 270, 270})
	idle := image.NewNineSlice(buttonImage, [3]int{100, 100, 100}, [3]int{63, 63, 63})

	hover := image.NewNineSlice(buttonImage, [3]int{100, 100, 100}, [3]int{63, 63, 63})

	pressed := image.NewNineSlice(buttonImage, [3]int{100, 100, 100}, [3]int{63, 63, 63})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil

	*/

	idle := image.NewNineSliceColor(color.RGBA{R: 170, G: 170, B: 180, A: 255})
	hover := image.NewNineSliceColor(color.RGBA{R: 130, G: 130, B: 150, A: 255})
	pressed := image.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadButtonIcon() *ebiten.Image {
	// we'll use a circle as an icon image
	// in reality it could be an arbitrary *ebiten.Image
	icon := ebiten.NewImage(32, 32)
	ebitenutil.DrawCircle(icon, 16, 16, 16, color.RGBA{R: 0x71, G: 0x56, B: 0xbd, A: 255})
	return icon
}
