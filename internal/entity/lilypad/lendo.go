package lilypad

import (
	"bytes"
	_ "embed"
	"image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/internal/entity/circlepacking"
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

	target := ebiten.NewImageFromImage(img)

	packing := circlepacking.New(target, 30)
	for range 10_000 {
		packing.Update()
	}

	positions, sizes := packing.Compile()
	rotations := make([]float64, 0, len(positions))

	for range len(positions) {
		a := cfg.Rand.Float64()
		rotations = append(rotations, a*math.Pi*2)
	}

	return New(width, height, positions, sizes, rotations)
}
