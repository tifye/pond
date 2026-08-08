package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	orbitPoint  mathutil.Point
	followPoint mathutil.Point
	fish        entity.Fish

	stillFish entity.Fish

	offscreen *ebiten.Image
}

func newGame() *test {
	return &test{
		orbitPoint: mathutil.Point{X: width / 2, Y: height / 2},
		background: color.RGBA{38, 38, 38, 255},
		fish: *entity.NewFish(entity.FishConfig{
			BoneChainUpdater: entity.BoneChainFABRIKUpdater{
				MinAngle: math.Pi * 0.85,
			},
		}),
		stillFish: *entity.NewFish(entity.FishConfig{
			BoneChainUpdater: entity.BoneChainLineUpdater{},
		}),

		offscreen: ebiten.NewImage(width, height),
	}
}

func (g *test) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	speed := 0.03
	amplitude := 50.0
	phase := float64(ebiten.Tick())
	g.followPoint.Y = math.Sin(phase*speed*2)*amplitude + g.orbitPoint.Y
	g.followPoint.X = math.Cos(phase*speed)*amplitude*math.Pi + g.orbitPoint.X

	g.fish.Position = g.followPoint
	g.fish.Update()

	g.stillFish.Position.X = 50
	g.stillFish.Position.Y = 150
	g.stillFish.Update()

	return nil
}

func (g *test) Draw(screen *ebiten.Image) {
	// screen.Fill(g.background)

	vector.FillCircle(
		screen,
		float32(g.followPoint.X),
		float32(g.followPoint.Y),
		12,
		color.RGBA{244, 63, 94, 255},
		false,
	)

	g.fish.Draw(screen)

	g.stillFish.Draw(screen)
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
