package lilypad

import (
	"bytes"
	_ "embed"
	"image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed lendo.png
var lendoLogoData []byte

func NewUsingCirclePacking(cfg Config) *Lilypads {
	cfg = cfg.WithDefaults()

	img, err := png.Decode(bytes.NewReader(lendoLogoData))
	if err != nil {
		panic(err)
	}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	pixels := make([]byte, 4*width*height)

	target := ebiten.NewImageFromImage(img)
	target.ReadPixels(pixels)

	positions := make([]mathutil.Point, 0, cfg.Amount)
	sizes := make([]int, 0, cfg.Amount)
	rotations := make([]float64, 0, cfg.Amount)

	for len(positions) < int(cfg.Amount) {
		x, y := cfg.Rand.IntN(width), cfg.Rand.IntN(height)
		positions = append(positions, mathutil.NewPoint(float64(x), float64(y)))

		a := cfg.Rand.Float64()
		sizes = append(sizes, int(cfg.MinRadius+(cfg.MaxRadius-cfg.MinRadius)*a))
		rotations = append(rotations, a*math.Pi*2)
	}

	return New(width, height, positions, sizes, rotations)
}
