package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tifye/pond/colors"
	"github.com/tifye/pond/internal/entity"
	"github.com/tifye/pond/pkg/mathutil"
)

const (
	width  = 600
	height = 400
)

// test renders two fish side by side so the BoneChain updaters can be compared:
// the left fish uses BoneChainLineUpdater (a straight line) and the right fish
// uses BoneChainSineUpdater (an animated sine wave).
type test struct {
	background color.Color

	center mathutil.Point
	fp     mathutil.Point

	offscreen *ebiten.Image

	tail *entity.CaudralFin
}

func newGame() *test {
	return &test{
		background: color.RGBA{38, 38, 38, 255},
		offscreen:  ebiten.NewImage(width, height),
		center:     mathutil.NewPoint(float64(width)/2, float64(height)/2),

		tail: entity.NewTail(),
	}
}

func (g *test) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	scale := 1.0
	amplitude := 50.0
	frequency := float64(ebiten.Tick()) / float64(ebiten.TPS())
	g.fp.X = math.Cos(frequency*scale) * amplitude * math.Pi
	g.fp.Y = math.Sin(frequency*scale*2) * amplitude

	g.fp = g.fp.Add(g.center)

	g.tail.Position = g.fp
	g.tail.Update()

	return nil
}

func (g *test) Draw(screen *ebiten.Image) {
	// screen.Fill(g.background)

	debugCircle(screen, g.fp, colors.Rose600)
	g.tail.Draw(screen)
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
