package entity

import (
	_ "embed"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed shaders/lilypad_shader.go
var lilypadShaderSrc []byte
var lilpadShader *ebiten.Shader

func init() {
	if s, err := ebiten.NewShader(lilypadShaderSrc); err != nil {
		panic(err)
	} else {
		lilpadShader = s
	}
}

type LilyPads struct {
	Positions []mathutil.Point

	vertices []ebiten.Vertex
	indices  []uint16
}

func NewLilyPads(v []mathutil.Point, sizes []int, angles []float64) *LilyPads {
	var vertexOffset uint16 = 0

	vertices := make([]ebiten.Vertex, 0)
	indices := make([]uint16, 0)

	for i, lp := range v {
		a := angles[i]

		left := lp.X
		right := lp.X + float64(sizes[i])

		top := lp.Y
		bottom := lp.Y + float64(sizes[i])

		topLeft := mathutil.NewPoint(left, top).RotateAround(lp, a)
		topRight := mathutil.NewPoint(right, top).RotateAround(lp, a)
		bottomLeft := mathutil.NewPoint(left, bottom).RotateAround(lp, a)
		bottomRight := mathutil.NewPoint(right, bottom).RotateAround(lp, a)

		vertices = append(vertices,
			// Top-Left (UV: 0,0)
			ebiten.Vertex{DstX: float32(topLeft.X), DstY: float32(topLeft.Y), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			// Top-Right (UV: 1,0)
			ebiten.Vertex{DstX: float32(topRight.X), DstY: float32(topRight.Y), SrcX: 1, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			// Bottom-Left (UV: 0,1)
			ebiten.Vertex{DstX: float32(bottomLeft.X), DstY: float32(bottomLeft.Y), SrcX: 0, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			// Bottom-Right (UV: 1,1)
			ebiten.Vertex{DstX: float32(bottomRight.X), DstY: float32(bottomRight.Y), SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		)

		indices = append(indices,
			vertexOffset+0, vertexOffset+1, vertexOffset+2,
			vertexOffset+1, vertexOffset+2, vertexOffset+3,
		)

		vertexOffset += 4
	}

	return &LilyPads{
		Positions: v,

		vertices: vertices,
		indices:  indices,
	}
}

func (l *LilyPads) Draw(target *ebiten.Image) {
	op := &ebiten.DrawTrianglesShaderOptions{}
	target.DrawTrianglesShader(l.vertices, l.indices, lilpadShader, op)
}

// type GenerateLilypadsOptions struct {
// 	MinRadius float64
// 	MaxRadius float64

// 	MinAmount int
// 	MaxAmount int
// }

// func GenerateLilypads(r *rand.Rand, width, height int, opts GenerateLilypadsOptions) *LilyPads {
// 	rn := r.Float64()

// 	n := opts.MinAmount + int(math.Floor(rn*float64(opts.MaxAmount-opts.MinAmount)))

// 	positions := make([]mathutil.Point, 0, n)
// 	sizes := make([]int, 0, n)
// 	angles := make([]float64, 0, n)

// }

type LilypadPatches struct {
	Positions []mathutil.Point
	Sizes     []mathutil.Size
}

type LilypadPatchesOptions struct {
	MinAmount int
	MaxAmount int

	MinSize int
	MaxSize int
}

func GenerateLilypadPatches(r *rand.Rand, containerWidth, containerHeight int, opts LilypadPatchesOptions) *LilypadPatches {
	rn := r.Float64()
	n := opts.MinAmount + int(math.Floor(rn*float64(opts.MaxAmount-opts.MinAmount)))

	positions := make([]mathutil.Point, 0, n)
	sizes := make([]mathutil.Size, 0, n)

	sizeRange := opts.MaxSize - opts.MinSize
	for range n {
		width := opts.MinSize + r.IntN(sizeRange)
		height := opts.MinSize + r.IntN(sizeRange)

		x := r.IntN(containerWidth)
		y := r.IntN(containerHeight)

		positions = append(positions, mathutil.NewPoint(float64(x), float64(y)))
		sizes = append(sizes, mathutil.NewPoint(float64(width), float64(height)))
	}

	return &LilypadPatches{
		Positions: positions,
		Sizes:     sizes,
	}
}
