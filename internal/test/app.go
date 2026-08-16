package test

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tifye/pond/internal/app/scenes"
)

type TestApp struct {
	Width  int
	Height int
	scene  scenes.Scene
}

type TestAppConfig struct {
	Width      int
	Height     int
	WindowName string
}

func (cfg TestAppConfig) withDefaults() TestAppConfig {
	if cfg.Height == 0 {
		cfg.Height = 400
	}

	if cfg.Width == 0 {
		cfg.Width = 600
	}

	if cfg.WindowName == "" {
		cfg.WindowName = "Pond"
	}

	return cfg
}

func NewTestApp(cfg TestAppConfig) *TestApp {
	cfg = cfg.withDefaults()

	ebiten.SetWindowSize(cfg.Width, cfg.Height)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowTitle(cfg.WindowName)

	return &TestApp{
		Width:  cfg.Width,
		Height: cfg.Height,
	}
}

func (t *TestApp) Run(scene scenes.Scene) {
	t.scene = scene

	opts := &ebiten.RunGameOptions{
		ScreenTransparent: true,
	}

	if err := ebiten.RunGameWithOptions(t, opts); err != nil {
		log.Fatal(err)
	}
}

func (t *TestApp) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if nextScene := t.scene.Update(); nextScene != nil {
		t.scene = nextScene
		t.scene.Update()
	}

	return nil
}

func (t *TestApp) Draw(screen *ebiten.Image) {
	t.scene.Draw(screen)
}

func (t *TestApp) Layout(outherWidth, outerHeight int) (int, int) {
	return t.Width, t.Height
}
