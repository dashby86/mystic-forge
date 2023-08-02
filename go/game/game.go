package game

import (
	"fmt"
	"image/color"
	"log"
	"mf/battle"
	"mf/events"
	"mf/forge"
	"mf/models"

	my_image "image"
	sqlService "mf/services/sql"

	"github.com/davecgh/go-spew/spew"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	Background  *ebiten.Image
	Crafted     bool
	Ui          *ebitenui.UI
	Sql         sqlService.SqlService
	Player      models.Player
	CraftedGear models.Gear
	Forge       forge.Forge
	Events      *events.EventQueue
}

func (g *Game) Update() error {
	g.Ui.Update()
	g.CharWindow()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the background image
	screen.DrawImage(g.Background, nil)

	g.Ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Set the screen size
	return 1024, 1024
}

func (g *Game) ShowCraftMenu(gear models.Gear, crafted models.Gear) {
	// Check if the forge has been clicked
	// Create the container
	myImage, _, err := ebitenutil.NewImageFromFile("assets/forge-menu.png")
	if err != nil {
		log.Fatal(err)
	}
	face, _ := g.loadFont(12)
	buttonSlice, _ := g.loadButtonImage()
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
			g.Forge.EquipGear(g.CraftedGear)
			g.Player, _ = g.Sql.GetPlayerByID(g.Player.Id)
			window.Close()
			g.Crafted = false
		}),
	)
	windowContainer.AddChild(cb)

	//g.ui.Container.BackgroundImage.Draw(nineSlice)
	//buttonStackedLayout.AddChild(button)
	//innerContainer.AddChild(button)
	g.Ui.AddWindow(window)
	//g.ui.Container.AddChild(window)
}

func (g *Game) CharWindow() {
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
	face, _ := g.loadFont(12)
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("HP: %d/%d", g.Player.CurrentHP, g.Player.MaxHP), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Attack: %d", g.Player.Attack), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Defense: %d", g.Player.Defense), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Attack Speed: %f", g.Player.Speed), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Crit: %%%f", g.Player.Crit), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Dodge: %%%f", g.Player.Dodge), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Block: %%%f", g.Player.Block), face, color.Color(color.Black)),
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
	g.Ui.AddWindow(window)
}

func (g *Game) BattleButton() {
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
	g.Ui.AddWindow(bwindow)
}

func (g *Game) Anvil() {
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
	g.Ui.AddWindow(window)
}

func (g *Game) Battle() {

	battler := battle.Battle{
		Player: g.Player,
		Enemy:  g.generateNextEnemy(g.Player),
		Sql:    g.Sql,
	}
	battler.SimBattle(g.Events)
}

func (g *Game) forge() {
	g.Ui.Container.Children()
	if !g.Crafted && g.Player.Ore > 0 {
		// Craft equipment
		fmt.Println("Crafting equipment...")
		//gear := models.Gear{}
		g.Player, _ = g.Sql.GetPlayerByID(g.Player.Id)
		gear := g.Forge.CraftGear()
		_, _ = g.Sql.SpendOre(1)
		current, err := g.Sql.GetEquippedGearBySlot(1, gear.SlotId)
		if err != nil {
			log.Fatal(err)
		}
		//spew.Dump(current)
		g.ShowCraftMenu(current, gear)
		g.CraftedGear = gear

		spew.Dump(gear)
		g.Crafted = true
	}
}

func (g *Game) loadFont(size float64) (font.Face, error) {
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

func (g *Game) loadButtonImage() (*widget.ButtonImage, error) {
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

func (g *Game) loadButtonIcon() *ebiten.Image {
	// we'll use a circle as an icon image
	// in reality it could be an arbitrary *ebiten.Image
	icon := ebiten.NewImage(32, 32)
	vector.DrawFilledCircle(icon, 16.0, 16.0, 16.0, color.RGBA{R: 0x71, G: 0x56, B: 0xbd, A: 255}, false)
	return icon
}

func (g *Game) generateNextEnemy(player models.Player) models.Enemy {
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
