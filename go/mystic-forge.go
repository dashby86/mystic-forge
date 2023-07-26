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
	my_image "image"
	"image/color"
	"log"
	"mf/forge"
	"mf/models"
	sqlService "mf/services/sql"
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
	anvil := Anvil{X: 200, Y: 200, Width: 200, Height: 200}

	// Implement game logic here
	spew.Dump(player)

	Forge := forge.Forge{
		Sql:    sService,
		Player: player,
	}

	game := game{
		Background: backgroundImg,
		Anvil:      anvil,
		ui:         &ui,
		sql:        sService,
		player:     player,
		Forge:      Forge,
	}
	// Display the stats of different equipment types
	if err := ebiten.RunGame(&game); err != nil {
		fmt.Println(err)
	}
}

type game struct {
	Background  *ebiten.Image
	Anvil       Anvil
	Crafted     bool
	ui          *ebitenui.UI
	sql         sqlService.SqlService
	player      models.Player
	craftedGear models.Gear
	Forge       forge.Forge
}

func (g *game) Update() error {
	// Check if the forge has been clicked and if a craft has not already been triggered
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !g.Crafted {
		x, y := ebiten.CursorPosition()
		forge := g.Anvil
		if x >= forge.X && x <= forge.X+forge.Width && y >= forge.Y && y <= forge.Y+forge.Height {
			g.ShowCraftMenu()
			// Craft equipment
			fmt.Println("Crafting equipment...")
			//gear := models.Gear{}
			g.player, _ = g.sql.GetPlayerByID()
			gear := g.Forge.CraftGear()
			g.craftedGear = gear
			/**
			enemy := models.Enemy{
				Name:    "Goblin",
				HP:      40416,
				Attack:  6178,
				Defense: 560,
				Speed:   281,
				Crit:    20,
				Dodge:   50,
				Block:   1,
			}

			battler := battle.Battle{
				Player: g.player,
				Enemey: enemy,
			}
			battler.SimBattle()

			*/

			spew.Dump(gear)
			g.Crafted = true
		}
	}
	g.charWindow()
	g.ui.Update()

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Draw the background image
	screen.DrawImage(g.Background, nil)

	// Draw the forge
	forge := g.Anvil
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
	face, _ := loadFont(12)
	buttonSlice, _ := loadButtonImage()
	nineSlice := image.NewNineSlice(myImage, [3]int{1020, 1020, 1020}, [3]int{555, 555, 555})
	/**
	innerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),

		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)

	*/
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice),
		//widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    40,
					Left:   40,
					Right:  40,
					Bottom: 20,
				}),
				widget.GridLayoutOpts.Spacing(30, 150),
			),
		),
	)
	c.AddChild(widget.NewText(
		widget.TextOpts.Text("This window blocks all input to widgets below it.", face, color.Color(color.Black)),
	))
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	c.AddChild(windowContainer)
	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(c),
		//Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
		//Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		//that is not a part of the window object
		//widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//Indicates that the window is draggable. It must have a TitleBar for this to work
		//widget.WindowOpts.Draggable(),
		//Set the window resizeable
		//widget.WindowOpts.Resizeable(),
		//Set the minimum size the window can be
		//widget.WindowOpts.MinSize(200, 100),
		//Set the maximum size a window can be
		//widget.WindowOpts.MaxSize(200, 100),
		//Set the callback that triggers when a move is complete
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		//Set the callback that triggers when a resize is complete
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),
	)
	x, y := window.Contents.PreferredSize()
	fmt.Println(x, y)
	//Create a rect with the preferred size of the content
	r := my_image.Rect(0, 0, 1020, 555)
	//Use the Add method to move the window to the specified point
	r = r.Add(my_image.Point{0, 400})
	//Set the windows location to the rect.
	window.SetLocation(r)
	/**
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
	*/

	o2b := widget.NewButton(
		widget.ButtonOpts.Image(buttonSlice),
		//widget.ButtonOpts.TextPadding(res.button.padding),
		//widget.ButtonOpts.Text("Open Another", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			g.Forge.EquipGear(g.craftedGear)
			g.player, _ = g.sql.GetPlayerByID()
			window.Close()
			g.Crafted = false
		}),
	)
	windowContainer.AddChild(o2b)

	cb := widget.NewButton(
		widget.ButtonOpts.Image(buttonSlice),
		//widget.ButtonOpts.TextPadding(res.button.padding),
		//widget.ButtonOpts.Text("Close", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			window.Close()
			g.Crafted = false
		}),
	)
	windowContainer.AddChild(cb)
	/**
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonSlice),
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				//Specify where within the row or column this element should be positioned.
				Position: widget.RowLayoutPositionEnd,
				//Should this widget be stretched across the row or column
				Stretch: false,
				//How wide can this element grow to (override preferred widget size)
				//MaxWidth: 100,
				//How tall can this element grow to (override preferred widget size)
				//MaxHeight: 100,
			}),
		),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			//g.ui.Container.RemoveChildren()
			window.Close()
			println("button clicked")
		}),
	)
	button2 := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonSlice),
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				//Specify where within the row or column this element should be positioned.
				Position: widget.RowLayoutPositionStart,
				//Should this widget be stretched across the row or column
				Stretch: false,
				//How wide can this element grow to (override preferred widget size)
				//MaxWidth: 100,
				//How tall can this element grow to (override preferred widget size)
				//MaxHeight: 100,
			}),
		),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			//g.ui.Container.RemoveChildren()
			window.Close()
			println("button clicked")
		}),
	)
	windowContainer.AddChild(button)
	windowContainer.AddChild(button2)

	*/

	//g.ui.Container.BackgroundImage.Draw(nineSlice)
	//buttonStackedLayout.AddChild(button)
	//innerContainer.AddChild(button)
	g.ui.AddWindow(window)
	//g.ui.Container.AddChild(window)
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
}

func loadButtonIcon() *ebiten.Image {
	// we'll use a circle as an icon image
	// in reality it could be an arbitrary *ebiten.Image
	icon := ebiten.NewImage(32, 32)
	ebitenutil.DrawCircle(icon, 16, 16, 16, color.RGBA{R: 0x71, G: 0x56, B: 0xbd, A: 255})
	return icon
}

func (g *game) charWindow() {
	myImage, _, err := ebitenutil.NewImageFromFile("assets/char_menu.png")
	if err != nil {
		log.Fatal(err)
	}
	nineSlice := image.NewNineSlice(myImage, [3]int{372, 372, 372}, [3]int{323, 323, 323})
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true, true, true, true, true, true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    80,
					Left:   80,
					Right:  40,
					Bottom: 80,
				}),
				widget.GridLayoutOpts.Spacing(30, 20),
			),
		),
	)
	face, _ := loadFont(12)
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("HP: %d", g.player.HP), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Attack: %d", g.player.Attack), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Defense: %d", g.player.Defense), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Speed: %d", g.player.Speed), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Crit: %%%d", g.player.Crit), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Dodge: %%%d", g.player.Dodge), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Block: %%%d", g.player.Block), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))

	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(c),
		//Set the window above everything else and block input elsewhere
		//widget.WindowOpts.Modal(),
		//Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		//that is not a part of the window object
		//widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//Indicates that the window is draggable. It must have a TitleBar for this to work
		//widget.WindowOpts.Draggable(),
		//Set the window resizeable
		//widget.WindowOpts.Resizeable(),
		//Set the minimum size the window can be
		//widget.WindowOpts.MinSize(200, 100),
		//Set the maximum size a window can be
		//widget.WindowOpts.MaxSize(200, 100),
		//Set the callback that triggers when a move is complete
	)
	//Create a rect with the preferred size of the content
	r := my_image.Rect(0, 0, 370, 320)
	//Use the Add method to move the window to the specified point
	r = r.Add(my_image.Point{655, 15})
	//Set the windows location to the rect.
	window.SetLocation(r)
	g.ui.AddWindow(window)
}
