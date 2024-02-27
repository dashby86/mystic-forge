package game

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	my_image "image"
	"image/color"
	"log"
	"mf/battle"
	"mf/enemy"
	"mf/forge"
	"mf/models"
	sqlService "mf/services/sql"
	"runtime"
)

type Game struct {
	Background  *ebiten.Image
	Crafted     bool
	Ui          *ebitenui.UI
	expBar      *widget.ProgressBar
	Sql         sqlService.SqlService
	Player      models.Player
	CraftedGear models.Gear
	Forge       forge.Forge
}

func (g *Game) Update() error {
	g.Ui.Update()
	g.CharWindow()
	runtime.GC()
	return nil
}

func (g *Game) UpdateProgressBar() int {
	if g.expBar != nil {
		currentExp := g.Player.Experience
		reqExp := g.Player.CalculateRequiredExp(g.Player.Level + 1)

		// Calculate the percentage filled
		percentageFilled := int(float64(currentExp) / float64(reqExp) * 100)

		// Ensure the percentage is within the range of 0-100
		if percentageFilled > 100 {
			percentageFilled = 100
		} else if percentageFilled < 0 {
			percentageFilled = 0
		}
		return percentageFilled

		//g.expBar.SetCurrent(percentageFilled)
	}
	return 0
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the background image
	screen.DrawImage(g.Background, nil)

	g.Ui.Draw(screen)
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
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
					Top:    150,
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
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true, true, true, true, true, true, true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    45,
					Left:   25,
					Right:  40,
					Bottom: 80,
				}),
				widget.GridLayoutOpts.Spacing(30, 10),
			),
		),
	)

	itemNineSlice := g.fetchGearIcon(g.CraftedGear)
	item := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(itemNineSlice),
		//widget.ContainerOpts.Layout()
	)

	currentItemGrid := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    100,
					Left:   150,
					Right:  0,
					Bottom: 40,
				}),
				widget.GridLayoutOpts.Spacing(30, 160),
			),
		),
	)

	currentItemGrid.AddChild(item)

	itemNineSlice = g.fetchGearIcon(g.CraftedGear)
	craftedItem := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(itemNineSlice),
		//widget.ContainerOpts.Layout()
	)
	craftedItemGrid := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    100,
					Left:   0,
					Right:  0,
					Bottom: 40,
				}),
				widget.GridLayoutOpts.Spacing(30, 160),
			),
		),
	)

	craftedItemGrid.AddChild(craftedItem)

	//face, _ := loadFont(12)
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Level %d: %s", gear.Level, gear.GetRarity(gear.Rarity)), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
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
		widget.TextOpts.Text(fmt.Sprintf("Speed: %d", gear.Speed), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Crit: %%%.2f", gear.Crit), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Dodge: %%%.2f", gear.Dodge), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c2.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Block: %%%.2f", gear.Block), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))

	c3 := widget.NewContainer(
		//widget.ContainerOpts.BackgroundImage(nineSlice),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true, true, true, true, true, true, true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    45,
					Left:   45,
					Right:  15,
					Bottom: 40,
				}),
				widget.GridLayoutOpts.Spacing(30, 10),
			),
		),
	)
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Level %d: %s", crafted.Level, crafted.GetRarity(crafted.Rarity)), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
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
		widget.TextOpts.Text(fmt.Sprintf("Speed: %d", crafted.Speed), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Crit: %%%.2f", crafted.Crit), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Dodge: %%%.2f", crafted.Dodge), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c3.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Block: %%%.2f", crafted.Block), face, color.Color(color.White)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))

	equipmentLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(50),
		)),
	)
	equipmentLayout.AddChild(currentItemGrid)
	equipmentLayout.AddChild(c2)
	equipmentLayout.AddChild(c3)
	equipmentLayout.AddChild(craftedItemGrid)
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
			g.Player, _ = g.Sql.GetPlayerByID()
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
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true, true, true, true, true, true, true, true, true}),
				//widget.GridLayoutOpts.Padding(15),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    80,
					Left:   95,
					Right:  40,
					Bottom: 80,
				}),
				widget.GridLayoutOpts.Spacing(30, 10),
			),
		),
	)
	face, _ := loadFont(12)
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("%s", g.Player.Name), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))

	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Level: %d Exp: %d", g.Player.Level, g.Player.Experience), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("HP: %d", g.Player.HP), face, color.Color(color.Black)),
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
		widget.TextOpts.Text(fmt.Sprintf("Speed: %d", g.Player.Speed), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Crit: %%%.2f", g.Player.Crit), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Dodge: %%%.2f", g.Player.Dodge), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))
	c.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Block: %%%.2f", g.Player.Block), face, color.Color(color.Black)),
		//widget.TextOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	))

	// Experience Bar
	//currentExp := g.Player.Experience
	//expRequiredForNextLevel := g.Player.CalculateRequiredExp(g.Player.Level + 1)

	expBarContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout()))

	// Create the progress bar
	g.expBar = widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this exampe since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(200, 20),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the progress images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 100, g.UpdateProgressBar()),
		// Set how much of the track is displayed when the bar is overlayed.
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
	)
	//g.updateProgressBar() // Initial update

	expBarContainer.AddChild(g.expBar)
	// ... add other elements to expBarContainer ...

	c.AddChild(expBarContainer)
	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(c),
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
	idle := image.NewNineSlice(buttonImage, [3]int{400, 400, 400}, [3]int{400, 400, 400})
	button := &widget.ButtonImage{
		Idle:    idle,
		Hover:   idle,
		Pressed: idle,
	}

	cb := widget.NewButton(
		widget.ButtonOpts.Image(button),
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
	r := my_image.Rect(0, 0, 400, 400)
	//Use the Add method to move the window to the specified point
	r = r.Add(my_image.Point{330, 600})
	//Set the windows location to the rect.
	window.SetLocation(r)
	g.Ui.AddWindow(window)
}

func (g *Game) Battle() {
	enemy := enemy.Enemy{
		Name:    "Goblin",
		Level:   1,
		HP:      1200,
		Attack:  300,
		Defense: 120,
		Speed:   3,
		Crit:    3,
		Dodge:   2,
		Block:   1,
	}

	currentEnemy := enemy.CreateEnemy(g.Player.DungeonLevel, enemy)

	battler := battle.Battle{
		Player: g.Player,
		Enemy:  currentEnemy,
		Sql:    g.Sql,
	}
	battler.SimBattle()
	g.Player, _ = g.Sql.GetPlayerByID()
}

func (g *Game) forge() {
	g.Ui.Container.Children()
	if g.Crafted == false && g.Player.Ore > 0 {
		// Craft equipment
		fmt.Println("Crafting equipment...")
		g.Player, _ = g.Sql.GetPlayerByID()
		gear := g.Forge.CraftGear()
		_, _ = g.Sql.SpendOre(1)
		g.CraftedGear = gear
		current, err := g.Sql.GetEquipedGearBySlot(1, gear.SlotId)
		if err != nil {
			log.Fatal(err)
		}
		spew.Dump(current)
		g.ShowCraftMenu(current, gear)

		spew.Dump(gear)
		g.UpdateProgressBar()
		g.Crafted = true
	}
}

func gearTextTemplate() {

}

func (g *Game) battleScene() {

}

func (g *Game) fetchGearIcon(gear models.Gear) *image.NineSlice {

	num := 1
	path := fmt.Sprintf("assets/equipment/%s/%s/%d.png", slotNames[gear.SlotId], rarityNames[0], num)
	gearImage, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return image.NewNineSlice(gearImage, [3]int{100, 100, 100}, [3]int{100, 100, 100})
}

var rarityNames = []string{
	"junk",
	"common",
	"uncommon",
	"rare",
	"epic",
	"legendary",
	"mythic",
}

var slotNames = []string{
	"blank",
	"helmet",
	"pauldrons",
	"gloves",
	"boots",
	"greaves",
	"ring",
	"necklace",
	"weapon",
	"chest",
}
