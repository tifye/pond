package app

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type App struct {
	dt      float64
	elapsed float64

	mouseX, mouseY               float64
	screenCenterX, screenCenterY float64
	circleX, circleY             float64
	orbitRadius                  float64
}

func NewApp() *App {
	w, h := ebiten.Monitor().Size()
	// Transpacy doesnt work when setting window to exact size
	// of monitor. My guess is that it thinks that nothing needs to render
	// behind it and that changes the properties of the window
	ebiten.SetWindowSize(w-1, h-1)
	// ebiten.SetFullscreen(true)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowTitle("Pond")

	return &App{
		screenCenterX: float64(w / 2),
		screenCenterY: float64(h / 2),
		dt:            1.0 / 60.0, // default for ebiten. It runs on fixed update rate
		orbitRadius:   50,
	}
}

func (a *App) Run() error {
	runOpts := ebiten.RunGameOptions{
		ScreenTransparent: true,
	}
	return ebiten.RunGameWithOptions(a, &runOpts)
}

func (a *App) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	a.elapsed += a.dt

	a.circleX = a.screenCenterX + math.Cos(a.elapsed*2.0)*a.orbitRadius
	a.circleY = a.screenCenterY + math.Sin(a.elapsed*2.0)*a.orbitRadius

	mx, my := ebiten.CursorPosition()
	a.mouseX = float64(mx)
	a.mouseY = float64(my)

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(a.circleX), float32(a.circleY), 50, color.RGBA{R: 125, G: 200, B: 85, A: 255}, false)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}
