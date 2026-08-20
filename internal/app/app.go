package app

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tifye/pond/internal/app/scenes"
)

type App struct {
	activeScene scenes.Scene
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

	scene := &scenes.BookstrapScene{
		Next: &scenes.LoadingScene{
			ScreenWidth:  float64(w),
			ScreenHeight: float64(h),
		},
	}

	return &App{
		activeScene: scene,
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

	for {
		scene := a.activeScene.Update()
		if scene == nil {
			break
		}

		a.activeScene = scene

		if s, ok := scene.(scenes.Initializer); ok {
			s.Initialize()
		}
	}

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	a.activeScene.Draw(screen)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}
