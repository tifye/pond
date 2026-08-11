package main

import (
	"math"
	"math/rand/v2"
)

// Perlin implements Ken Perlin's "improved noise" (2002). It is seeded with a
// shuffled permutation table so runs are reproducible for a given seed.
type Perlin struct {
	perm [512]int
}

// NewPerlin builds a Perlin generator with a shuffled permutation table.
func NewPerlin(seed uint64) *Perlin {
	var p Perlin

	var base [256]int
	for i := range base {
		base[i] = i
	}

	// Fisher-Yates shuffle with a seeded PRNG.
	r := rand.New(rand.NewPCG(seed, seed*0x9E3779B97F4A7C15+1))
	for i := len(base) - 1; i > 0; i-- {
		j := r.IntN(i + 1)
		base[i], base[j] = base[j], base[i]
	}

	// Duplicate so we can index up to 511 without wrapping by hand.
	for i := range 512 {
		p.perm[i] = base[i&255]
	}

	return &p
}

// Noise2D returns 2D noise in roughly [-1, 1].
func (p *Perlin) Noise2D(x, y float64) float64 {
	return p.Noise3D(x, y, 0)
}

// Noise3D returns 3D noise in roughly [-1, 1]. Feed a slowly changing z (e.g.
// elapsed time) to animate a 2D field.
func (p *Perlin) Noise3D(x, y, z float64) float64 {
	xi := int(math.Floor(x)) & 255
	yi := int(math.Floor(y)) & 255
	zi := int(math.Floor(z)) & 255

	x -= math.Floor(x)
	y -= math.Floor(y)
	z -= math.Floor(z)

	u := fade(x)
	v := fade(y)
	w := fade(z)

	a := p.perm[xi] + yi
	aa := p.perm[a] + zi
	ab := p.perm[a+1] + zi
	b := p.perm[xi+1] + yi
	ba := p.perm[b] + zi
	bb := p.perm[b+1] + zi

	return lerp(w,
		lerp(v,
			lerp(u, grad(p.perm[aa], x, y, z), grad(p.perm[ba], x-1, y, z)),
			lerp(u, grad(p.perm[ab], x, y-1, z), grad(p.perm[bb], x-1, y-1, z)),
		),
		lerp(v,
			lerp(u, grad(p.perm[aa+1], x, y, z-1), grad(p.perm[ba+1], x-1, y, z-1)),
			lerp(u, grad(p.perm[ab+1], x, y-1, z-1), grad(p.perm[bb+1], x-1, y-1, z-1)),
		),
	)
}

// FBM sums several octaves of noise (fractal Brownian motion). The result is
// normalised back to roughly [-1, 1].
//
//   - octaves:     how many layers of detail to stack
//   - persistence: how much each octave contributes (amplitude falloff, < 1)
//   - lacunarity:  how much finer each octave gets (frequency growth, > 1)
func (p *Perlin) FBM(x, y, z float64, octaves int, persistence, lacunarity float64) float64 {
	var total, maxValue float64
	amplitude := 1.0
	frequency := 1.0

	for i := 0; i < octaves; i++ {
		total += p.Noise3D(x*frequency, y*frequency, z*frequency) * amplitude
		maxValue += amplitude
		amplitude *= persistence
		frequency *= lacunarity
	}

	if maxValue == 0 {
		return 0
	}
	return total / maxValue
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

// grad is Perlin's reference gradient function: it picks one of 12 gradient
// directions from the low bits of the hash and returns its dot product with
// (x, y, z).
func grad(hash int, x, y, z float64) float64 {
	h := hash & 15

	var u, v float64
	if h < 8 {
		u = x
	} else {
		u = y
	}
	switch {
	case h < 4:
		v = y
	case h == 12 || h == 14:
		v = x
	default:
		v = z
	}

	if h&1 != 0 {
		u = -u
	}
	if h&2 != 0 {
		v = -v
	}
	return u + v
}
