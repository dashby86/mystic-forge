package main

import (
	"database/sql"
	"fmt"
	my_image "image"
	"image/color"
	"log"
	"mf/battle"
	"mf/events"
	"mf/forge"
	"mf/models"
	sqlService "mf/services/sql"

	"github.com/davecgh/go-spew/spew"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
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
	backgroundImg, _, _ := ebitenutil.NewImageFromFile("assets/mainforge.png")

	// Implement game logic here
	spew.Dump(player)

	Forge := forge.Forge{
		Sql:    sService,
		Player: player,
	}

	// Create an events queue
	eventQueue := &events.EventQueue{}

	game := game{
		Background: backgroundImg,
		ui:         &ui,
		sql:        sService,
		player:     player,
		Forge:      Forge,
		Events:     eventQueue,
	}
	game.charWindow()
	game.anvil()
	game.battleButton()
	go game.Events.EventDispatcher()

	// Display the stats of different equipment types
	if err := ebiten.RunGame(&game); err != nil {
		fmt.Println(err)
	}
}

type game struct {
	Background  *ebiten.Image
	Crafted     bool
	ui          *ebitenui.UI
	sql         sqlService.SqlService
	player      models.Player
	craftedGear models.Gear
	Forge       forge.Forge
	Events      *events.EventQueue
}

func (g *game) Update() error {
	g.ui.Update()
	g.charWindow()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Draw the background image
	screen.DrawImage(g.Background, nil)

	g.ui.Draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Set the screen size
	return 1024, 1024
}

func (g *game) ShowCraftMenu(gear models.Gear, crafted models.Gear) {
	// Check if the forge has been clicked
	// Create the container
	myImage, _, err := ebitenutil.NewImageFromFile("assets/forge-menu.png")
	if err != nil {
		log.Fatal(err)
	}
	face, _ := loadFont(12)
	buttonSlice, _ := loadButtonImage()
	nineSlice := image.NewNineSlice(myImage, [3]int{950, 950, 950}, [3]int{635, 635, 635})
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
					Left:   45,
					Right:  40,
					Bottom: 40,
				}),
				widget.GridLayoutOpts.Spacing(30, 160),
			),
		),
	)
	c2 := widget.NewContainer(
		//widget.ContainerOpts.BackgroundImage(nineSlice),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true, true, true, true, true, true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    150,
					Left:   250,
					Right:  40,
					Bottom: 80,
				}),
				widget.GridLayoutOpts.Spacing(30, 10),
			),
		),
	)
	//face, _ := loadFont(12)
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("HP: %d", gear.HP), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Attack: %d", gear.Attack), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Defense: %d", gear.Defense), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Attack Speed: %.2f", gear.Speed), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Crit: %.2f%%", gear.Crit), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Dodge: %.2f%%", gear.Dodge), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Block: %.2f%%", gear.Block), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))

	c3 := widget.NewContainer(
		//widget.ContainerOpts.BackgroundImage(nineSlice),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true, true, true, true, true, true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    150,
					Left:   50,
					Right:  40,
					Bottom: 80,
				}),
				widget.GridLayoutOpts.Spacing(30, 10),
			),
		),
	)
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("HP: %d", crafted.HP), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Attack: %d", crafted.Attack), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Defense: %d", crafted.Defense), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Attack Speed: %.2f", crafted.Speed), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Crit: %.2f%%", crafted.Crit), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Dodge: %.2f%%", crafted.Dodge), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Block: %.2f%%", crafted.Block), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))

	equipmentLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(150),
		)),
	)
	equipmentLayout.AddChild(c2)
	equipmentLayout.AddChild(c3)
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(400),
		)),
	)
	c.AddChild(equipmentLayout)
	c.AddChild(windowContainer)
	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(c),
		//Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
	)
	x, y := window.Contents.PreferredSize()
	fmt.Println(x, y)
	//Create a rect with the preferred size of the content
	r := my_image.Rect(0, 0, 950, 635)
	//Use the Add method to move the window to the specified point
	r = r.Add(my_image.Point{40, 350})
	//Set the windows location to the rect.
	window.SetLocation(r)

	o2b := widget.NewButton(
		widget.ButtonOpts.Image(buttonSlice),
		//widget.ButtonOpts.TextPadding(res.button.padding),
		//widget.ButtonOpts.Text("Open Another", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
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
			g.Forge.EquipGear(g.craftedGear)
			g.player, _ = g.sql.GetPlayerByID()
			window.Close()
			g.Crafted = false
		}),
	)
	windowContainer.AddChild(cb)

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
	buttonImage, _, err := ebitenutil.NewImageFromFile("assets/cta.png")
	if err != nil {
		log.Fatal(err)
	}
	//buttonSlice := image.NewNineSlice(buttonImage, [3]int{310, 310, 310}, [3]int{270, 270, 270})
	idle := image.NewNineSlice(buttonImage, [3]int{150, 150, 150}, [3]int{150, 150, 150})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   idle,
		Pressed: idle,
	}, nil
}

func loadButtonIcon() *ebiten.Image {
	// we'll use a circle as an icon image
	// in reality it could be an arbitrary *ebiten.Image
	icon := ebiten.NewImage(32, 32)
	vector.DrawFilledCircle(icon, 16.0, 16.0, 16.0, color.RGBA{R: 0x71, G: 0x56, B: 0xbd, A: 255}, false)
	return icon
}

func (g *game) charWindow() {
	myImage, _, err := ebitenutil.NewImageFromFile("assets/char_menu.png")
	if err != nil {
		log.Fatal(err)
	}
	nineSlice := image.NewNineSlice(myImage, [3]int{240, 240, 240}, [3]int{360, 360, 360})
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSlice),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true, true, true, true, true, true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    80,
					Left:   95,
					Right:  40,
					Bottom: 80,
				}),
				widget.GridLayoutOpts.Spacing(30, 15),
			),
		),
	)
	face, _ := loadFont(12)
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("HP: %d/%d", g.player.CurrentHP, g.player.MaxHP), face, color.Color(color.Black)),
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
		widget.TextOpts.Text(fmt.Sprintf("Attack Speed: %f", g.player.Speed), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Crit: %%%f", g.player.Crit), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Dodge: %%%f", g.player.Dodge), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Block: %%%f", g.player.Block), face, color.Color(color.Black)),
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
	r := my_image.Rect(0, 0, 240, 360)
	//Use the Add method to move the window to the specified point
	r = r.Add(my_image.Point{655, 15})
	//Set the windows location to the rect.
	window.SetLocation(r)
	g.ui.AddWindow(window)
}

func (g *game) battleButton() {
	//nineSlice := image.NewNineSlice(myImage, [3]int{150, 150, 150}, [3]int{150, 150, 150})
	innerContainer := widget.NewContainer(
		//widget.ContainerOpts.BackgroundImage(nineSlice),
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
	buttonImage, _, err := ebitenutil.NewImageFromFile("assets/battle-icon.png")
	if err != nil {
		log.Fatal(err)
	}
	//buttonSlice := image.NewNineSlice(buttonImage, [3]int{310, 310, 310}, [3]int{270, 270, 270})
	idle := image.NewNineSlice(buttonImage, [3]int{150, 150, 150}, [3]int{150, 150, 150})

	hover := image.NewNineSlice(buttonImage, [3]int{150, 150, 150}, [3]int{150, 150, 150})

	pressed := image.NewNineSlice(buttonImage, [3]int{150, 150, 150}, [3]int{150, 150, 150})

	button := &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
	cb := widget.NewButton(
		widget.ButtonOpts.Image(button),
		//widget.ButtonOpts.TextPadding(res.button.padding),
		//widget.ButtonOpts.Text("Close", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			g.Battle()
		}),
	)
	innerContainer.AddChild(cb)
	bwindow := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(innerContainer),
	)
	//Create a rect with the preferred size of the content
	r := my_image.Rect(0, 0, 150, 150)
	//Use the Add method to move the window to the specified point
	r = r.Add(my_image.Point{800, 800})
	//Set the windows location to the rect.
	bwindow.SetLocation(r)
	g.ui.AddWindow(bwindow)
}

func (g *game) anvil() {
	innerContainer := widget.NewContainer(
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

	buttonImage, _, err := ebitenutil.NewImageFromFile("assets/anvil.png")
	if err != nil {
		log.Fatal(err)
	}
	//buttonSlice := image.NewNineSlice(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}), [3]int{600, 600, 600}, [3]int{600, 600, 600})
	idle := image.NewNineSlice(buttonImage, [3]int{600, 600, 600}, [3]int{600, 600, 600})
	hover := image.NewNineSlice(buttonImage, [3]int{600, 600, 600}, [3]int{600, 600, 600})
	pressed := image.NewNineSlice(buttonImage, [3]int{600, 600, 600}, [3]int{600, 600, 600})
	button := &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}

	cb := widget.NewButton(
		widget.ButtonOpts.Image(button),
		//widget.ButtonOpts.TextPadding(res.button.padding),
		//widget.ButtonOpts.Text("Close", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			g.forge()
		}),
	)
	innerContainer.AddChild(cb)
	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(innerContainer),
	)
	//Create a rect with the preferred size of the content
	r := my_image.Rect(0, 0, 600, 600)
	//Use the Add method to move the window to the specified point
	r = r.Add(my_image.Point{330, 500})
	//Set the windows location to the rect.
	window.SetLocation(r)
	g.ui.AddWindow(window)
}

func (g *game) Battle() {

	battler := battle.Battle{
		Player: g.player,
		Enemy:  generateNextEnemy(g.player),
		Sql:    g.sql,
	}
	battler.SimBattle(g.Events)
}

func generateNextEnemy(player models.Player) models.Enemy {
	//TODO: Implement next enemy logic.
	return models.Enemy{
		Name:      "Goblin",
		MaxHP:     40416,
		CurrentHP: 40416,
		Attack:    6178,
		Defense:   560,
		Speed:     281,
		Crit:      20,
		Dodge:     50,
		Block:     1,
	}
}

func (g *game) forge() {
	g.ui.Container.Children()
	if !g.Crafted && g.player.Ore > 0 {
		// Craft equipment
		fmt.Println("Crafting equipment...")
		//gear := models.Gear{}
		g.player, _ = g.sql.GetPlayerByID()
		gear := g.Forge.CraftGear()
		_, _ = g.sql.SpendOre(1)
		current, err := g.sql.GetEquippedGearBySlot(1, gear.SlotId)
		if err != nil {
			log.Fatal(err)
		}
		//spew.Dump(current)
		g.ShowCraftMenu(current, gear)
		g.craftedGear = gear

		spew.Dump(gear)
		g.Crafted = true
	}
}

// func gearTextTemplate() {

// }
