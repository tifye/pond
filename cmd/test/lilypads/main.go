package main

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/internal/app/scenes"
	"github.com/tifye/pond/internal/test"
	"github.com/tifye/pond/pkg/mathutil"
)

func main() {
	app := test.NewTestApp(test.TestAppConfig{})

	scene := &scene{
		Width:  app.Width,
		Height: app.Height,
	}

	app.Run(&scenes.BookstrapScene{
		Next: scene,
	})
}

type scene struct {
	Width  int
	Height int

	Noise *ebiten.Image
}

func (s *scene) Initialize() {
	seed1, seed2 := rand.Uint64(), rand.Uint64()
	rsrc := rand.NewPCG(seed1, seed2)

	perlin := mathutil.NewPerlin(rand.New(rsrc))

	pixelBuffer := make([]byte, s.Width*s.Height*4)
	buffIdx := 0
	noiseScale := 3.0
	for y := range s.Height {
		fy := (float64(y) / float64(s.Height)) * noiseScale
		for x := range s.Width {
			fx := (float64(x) / float64(s.Height)) * noiseScale

			v := perlin.FBM(fy, fx, 0, 4, 0.35, 3)

			v = math.Abs(v) / v         // 1 / v => cool effect
			beta := (v*0.5 + 0.5)       // 0 if v < 0, 1 if v > 0
			alpha := ((v*-1)*0.5 + 0.5) // 0 if v > 0, 1 if v < 0
			v = beta*255 + alpha*0

			pixelBuffer[buffIdx+0] = byte(v)
			pixelBuffer[buffIdx+1] = byte(v)
			pixelBuffer[buffIdx+2] = byte(v)
			pixelBuffer[buffIdx+3] = 255
			buffIdx += 4
		}
	}

	noise := ebiten.NewImage(s.Width, s.Height)
	noise.WritePixels(pixelBuffer)
	s.Noise = noise
}

func (s *scene) Update() scenes.Scene {
	return nil
}

func (s *scene) Draw(target *ebiten.Image) {
	target.DrawImage(s.Noise, nil)
}
