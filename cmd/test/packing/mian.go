package main

import (
	"bytes"
	_ "embed"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tifye/pond/internal/entity/circlepacking"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed image.png
var imageBytes []byte

const (
	width  = 1300
	height = 900
)

type test struct {
	background color.Color
	packing    *circlepacking.Packing
}

func newGame() *test {
	img, _ := png.Decode(bytes.NewReader(imageBytes))
	target := ebiten.NewImageFromImage(img)

	return &test{
		background: color.RGBA{38, 38, 38, 255},
		packing:    circlepacking.New(target, 20),
	}
}

func (g *test) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	g.packing.Update()

	return nil
}

func (g *test) Draw(screen *ebiten.Image) {
	g.packing.Draw(screen)
}

func debugCircle(t *ebiten.Image, p mathutil.Point, c color.Color) {
	vector.FillCircle(
		t,
		float32(p.X),
		float32(p.Y),
		4,
		c,
		false,
	)
}

func (g *test) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowTitle("Pond")

	opts := &ebiten.RunGameOptions{
		ScreenTransparent: true,
	}

	if err := ebiten.RunGameWithOptions(newGame(), opts); err != nil {
		log.Fatal(err)
	}
}
