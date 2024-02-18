package systems

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	"os"
)

type menusystem struct {
	UI        *ebitenui.UI
	menuEntry *donburi.Entry
}

var Menu = &menusystem{
	UI: &ebitenui.UI{Container: setupMenu()},
}

func setupMenu() *widget.Container {

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 0, 175})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	menuContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(30),
			),
		),
	)

	newGameButton := createMenuButton("New Game", func(args *widget.ButtonClickedEventArgs) {
		println("options button clicked")
	})
	optionsButton := createMenuButton("Options", func(args *widget.ButtonClickedEventArgs) {
		println("options button clicked")
	})
	quitButton := createMenuButton("Quit", func(args *widget.ButtonClickedEventArgs) {
		os.Exit(1)
	})

	// add the newGameButton as a child of the container
	menuContainer.AddChild(newGameButton)
	menuContainer.AddChild(optionsButton)
	menuContainer.AddChild(quitButton)
	rootContainer.AddChild(menuContainer)

	return rootContainer
}

func createMenuButton(buttonText string, handler func(args *widget.ButtonClickedEventArgs)) *widget.Button {
	buttonImage, _ := loadButtonImage()
	face, _ := loadFont(80)
	// construct a newGameButton
	button := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the newGameButton both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the newGameButton's text, the font face, and the color
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the newGameButton's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the newGameButton
		widget.ButtonOpts.ClickedHandler(handler),
	)
	return button
}

func (system *menusystem) Update(e *ecs.ECS) {
	// ui.Update() must be called in ebiten Update function, to handle user input and other things
	var ok bool
	if system.menuEntry == nil {
		if system.menuEntry, ok = components.Menu.First(e.World); !ok {
			panic("no menu found")
		}
	}

	menu := components.Menu.Get(system.menuEntry)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		log.Print(menu.IsOpen)

		menu.ToggleMenu()

		//Update function doesn't seem to respect IsPaused yet
		if e.IsPaused() {
			e.Resume()
		} else {
			e.Pause()
		}

		log.Print(menu.IsOpen)
	}

	system.UI.Update()
}

func (system *menusystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	// ui.Draw() should be called in the ebiten Draw function, to draw the UI onto the screen.
	// It should also be called after all other rendering for your game so that it shows up on top of your game world.
	system.UI.Draw(screen)
}

func (system *menusystem) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
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
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 255, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 255, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 255, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
