package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	width  = 600
	height = 400

	// renderScale renders the noise field at 1/renderScale resolution and then
	// upscales it, which keeps things smooth when octaves get expensive. Set to
	// 1 for a pixel-perfect (but heavier) field.
	renderScale = 2

	// Noise tuning — these are the knobs to play with.
	noiseScale  = 0.008 // zoom: smaller = smoother / more zoomed in
	octaves     = 1     // layers of detail stacked together
	persistence = 0.75  // amplitude falloff per octave (< 1)
	lacunarity  = 2.0   // frequency growth per octave (> 1)
	timeSpeed   = 0.01  // how fast the field animates
	threshold   = 0.5
)

// test is a sandbox for messing around with Perlin noise: it renders an
// animated fractal-noise field to a scratch buffer and blits it to the screen.
type test struct {
	perlin *Perlin

	bufW, bufH int
	pixels     []byte
	buffer     *ebiten.Image

	t      float64
	paused bool
}

func newGame() *test {
	bufW := width / renderScale
	bufH := height / renderScale

	g := &test{
		perlin: NewPerlin(1),
		bufW:   bufW,
		bufH:   bufH,
		pixels: make([]byte, bufW*bufH*4),
		buffer: ebiten.NewImage(bufW, bufH),
	}

	// Render one frame up front so a paused start still shows something.
	g.render()
	return g
}

func (g *test) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}

	if g.paused {
		return nil
	}

	g.t += timeSpeed
	g.render()
	return nil
}

// render samples the noise field into the pixel buffer and uploads it to the
// scratch image.
func (g *test) render() {
	i := 0
	for y := 0; y < g.bufH; y++ {
		fy := float64(y*renderScale) * noiseScale
		for x := 0; x < g.bufW; x++ {
			fx := float64(x*renderScale) * noiseScale

			n := g.perlin.FBM(fx, fy, g.t, octaves, persistence, lacunarity)

			// Map [-1, 1] -> [0, 1] and clamp.
			n = n*0.5 + 0.5

			var c byte
			if n > threshold {
				c = 255
			} else {
				c = 0
			}
			g.pixels[i] = c
			g.pixels[i+1] = c
			g.pixels[i+2] = c
			g.pixels[i+3] = 255
			i += 4
		}
	}

	g.buffer.WritePixels(g.pixels)
}

func (g *test) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(renderScale, renderScale)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(g.buffer, op)

	paused := ""
	if g.paused {
		paused = "  (paused)"
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"perlin noise sandbox%s\n"+
			"scale %.4f  octaves %d  persistence %.2f  lacunarity %.1f\n"+
			"speed %.3f  t %.1f\n"+
			"[space] pause   [esc] quit",
		paused, noiseScale, octaves, persistence, lacunarity, timeSpeed, g.t,
	))
}

func (g *test) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func main() {
	mWidth, _ := ebiten.Monitor().Size()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowPosition(mWidth-width, 300)
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
