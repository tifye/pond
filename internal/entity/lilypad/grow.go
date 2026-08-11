package lilypad

import (
	"math"
	"math/rand/v2"

	"github.com/tifye/pond/pkg/mathutil"
)

type Config struct {
	Rand      *rand.Rand
	Amount    uint
	MinRadius float64
	MaxRadius float64
}

var DefaultFromNoiseConfig = Config{
	Amount:    250,
	MinRadius: 25,
	MaxRadius: 50,
	Rand:      mathutil.GlobalRand,
}

func (cfg Config) WithDefaults() Config {
	if cfg.Amount == 0 {
		cfg.Amount = DefaultFromNoiseConfig.Amount
	}

	if cfg.MinRadius <= 0 {
		cfg.MinRadius = DefaultFromNoiseConfig.MinRadius
	}

	if cfg.MaxRadius <= 0 {
		cfg.MaxRadius = DefaultFromNoiseConfig.MaxRadius
	} else if cfg.MaxRadius < cfg.MinRadius {
		cfg.MaxRadius = cfg.MinRadius
	}

	if cfg.Rand == nil {
		cfg.Rand = DefaultFromNoiseConfig.Rand
	}

	return cfg
}

// data should be row major
func NewUsingNoiseThreshold(data []float64, width, height int, scale, threshold float64, cfg Config) *Lilypads {
	cfg = cfg.WithDefaults()

	positions := make([]mathutil.Point, 0, cfg.Amount)
	sizes := make([]int, 0, cfg.Amount)
	rotations := make([]float64, 0, cfg.Amount)

	for len(positions) < int(cfg.Amount) {
		x, y := cfg.Rand.IntN(width), cfg.Rand.IntN(height)
		idx := y*width + x

		noiseValue := data[idx]
		if noiseValue > threshold {
			positions = append(positions, mathutil.NewPoint(scale*float64(x), scale*float64(y)))

			a := cfg.Rand.Float64()
			sizes = append(sizes, int(cfg.MinRadius+(cfg.MaxRadius-cfg.MinRadius)*a))
			rotations = append(rotations, a*math.Pi*2)
		}
	}

	return New(int(scale*float64(width)), int(scale*float64(height)), positions, sizes, rotations)
}

func NewUsingNoiseChance(data []float64, width, height int, scale float64, f func(v float64) float64, cfg Config) *Lilypads {
	cfg = cfg.WithDefaults()

	positions := make([]mathutil.Point, 0, cfg.Amount)
	sizes := make([]int, 0, cfg.Amount)
	rotations := make([]float64, 0, cfg.Amount)

	minChance := 0.0
	maxChance := 1.0
	threshold := 1.0 - (maxChance - minChance)

	for len(positions) < int(cfg.Amount) {
		x, y := cfg.Rand.IntN(width), cfg.Rand.IntN(height)
		idx := y*width + x

		noiseValue := data[idx]

		a := cfg.Rand.Float64()
		if a >= threshold*f(noiseValue) {
			positions = append(positions, mathutil.NewPoint(scale*float64(x), scale*float64(y)))

			a := cfg.Rand.Float64()
			sizes = append(sizes, int(cfg.MinRadius+(cfg.MaxRadius-cfg.MinRadius)*a))
			rotations = append(rotations, a*math.Pi*2)
		}
	}

	return New(int(scale*float64(width)), int(scale*float64(height)), positions, sizes, rotations)
}
