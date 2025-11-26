package main

import (
	"image/color"
	"log"
	"math"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Mode = string

const (
	ScreensaverModeFlag Mode = "/c"
	ConfigModeFlag      Mode = "/p"
	PreviewModeFlag     Mode = "/s"
)

type Pond struct {
	dt      float64
	elapsed float64

	mouseX, mouseY float64

	screenCenterX float64
	screenCenterY float64

	circleX     float64
	circleY     float64
	orbitRadius float64
}

func (p *Pond) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	p.elapsed += p.dt

	p.circleX = p.screenCenterX + math.Cos(p.elapsed*2.0)*p.orbitRadius
	p.circleY = p.screenCenterY + math.Sin(p.elapsed*2.0)*p.orbitRadius

	mx, my := ebiten.CursorPosition()
	p.mouseX = float64(mx)
	p.mouseY = float64(my)
	return nil
}

func (p *Pond) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(p.circleX), float32(p.circleY), 50, color.RGBA{R: 125, G: 200, B: 85, A: 255}, false)
}

func (p *Pond) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}

func main() {
	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		switch {
		case strings.HasPrefix(arg, ConfigModeFlag):
			return
		case strings.HasPrefix(arg, PreviewModeFlag):
		case strings.HasPrefix(arg, ScreensaverModeFlag):
		default:
		}
	}

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

	p := &Pond{
		screenCenterX: float64(w / 2),
		screenCenterY: float64(h / 2),
		dt:            1.0 / 60.0,
		orbitRadius:   50,
	}
	runOpts := ebiten.RunGameOptions{
		ScreenTransparent: true,
	}
	if err := ebiten.RunGameWithOptions(p, &runOpts); err != nil {
		log.Fatal(err)
	}
}
