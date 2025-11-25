package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Pond struct {
	dt      float64
	elapsed float64

	screenCenterX float64
	screenCenterY float64

	circleX     float64
	circleY     float64
	orbitRadius float64
}

func (p *Pond) Update() error {
	p.elapsed += p.dt

	p.circleX = p.screenCenterX + math.Cos(p.elapsed*2.0)*p.orbitRadius
	p.circleY = p.screenCenterY + math.Sin(p.elapsed*2.0)*p.orbitRadius
	return nil
}

func (p *Pond) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(p.circleX), float32(p.circleY), 50, color.RGBA{R: 125, G: 200, B: 85, A: 255}, false)
}

func (p *Pond) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}

func main() {
	w, h := ebiten.Monitor().Size()
	// Transpacy doesnt work when setting window to exact size
	// of monitor. My guess is that it thinks that nothing needs to render
	// behind it and that changes the properties of the window
	ebiten.SetWindowSize(w-1, h-1)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowDecorated(false)
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
