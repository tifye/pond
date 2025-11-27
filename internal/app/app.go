package app

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tifye/pond/pkg/agent"
	"github.com/tifye/pond/pkg/mathutil"
	"github.com/tifye/pond/pkg/mathutil/fabrik"
)

type fish struct {
	bones       []mathutil.Point
	boneLengths []float64
}

type App struct {
	dt      float64
	elapsed float64

	debugColor color.Color

	mousePos                     mathutil.Point
	screenCenterX, screenCenterY float64

	koi    fish
	agents *agent.Agents
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

	screenCenterX := float64(w / 2)
	screenCenterY := float64(h / 2)

	koiBones := make([]mathutil.Point, 20)
	koiBoneLengths := make([]float64, len(koiBones)-1)
	koiBones[0].X = screenCenterX
	koiBones[0].Y = screenCenterY
	for i := 1; i < len(koiBones); i++ {
		koiBones[i].X = screenCenterX
		koiBones[i].Y = screenCenterY + 35.0*math.Pow(float64(i), 0.75)
		koiBoneLengths[i-1] = koiBones[i].Distance(koiBones[i-1])
	}

	agents := agent.NewAgents(1, 1)
	agents.AddBehaviour(agent.NewWander(agents.Num(), 1, 50, 50, math.Pi))
	agents.AddBehaviour(agent.Boundry(float64(w), float64(h), 200, 0.05))

	return &App{
		debugColor: color.RGBA{R: 125, G: 200, B: 85, A: 255},

		screenCenterX: screenCenterX,
		screenCenterY: screenCenterY,
		dt:            1.0 / 60.0, // default for ebiten. It runs on fixed update rate

		koi: fish{
			bones:       koiBones,
			boneLengths: koiBoneLengths,
		},
		agents: agents,
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

	mx, my := ebiten.CursorPosition()
	a.mousePos.X = float64(mx)
	a.mousePos.Y = float64(my)

	a.agents.Update(a.elapsed, a.dt)
	a.mousePos = a.agents.Position(0)

	fabrik.SolveFABRIK(
		a.koi.bones,
		a.koi.boneLengths,
		a.mousePos,
		math.Pi*0.5,
	)

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	for _, bone := range a.koi.bones {
		vector.FillCircle(screen, float32(bone.X), float32(bone.Y), 5, a.debugColor, false)
	}

	for i := range a.koi.boneLengths {
		// face the front
		curBone := a.koi.bones[i+1]
		nextBone := a.koi.bones[i]

		segment := curBone.Subtract(nextBone)
		left := segment.RotateCounterClockwise().Add(nextBone)
		right := segment.RotateClockwise().Add(nextBone)

		vector.FillCircle(screen, float32(left.X), float32(left.Y), 3, a.debugColor, false)
		vector.FillCircle(screen, float32(right.X), float32(right.Y), 3, a.debugColor, false)
	}
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}
