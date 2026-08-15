package lilypad

import (
	_ "embed"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed flower.kage
var flowerShaderSrc []byte
var flowerShader *ebiten.Shader

func init() {
	if s, err := ebiten.NewShader(flowerShaderSrc); err != nil {
		panic(err)
	} else {
		flowerShader = s
	}
}

type Flowers struct {
	Positions []mathutil.Point

	vertices []ebiten.Vertex
	indices  []uint16

	img *ebiten.Image
}

func NewFlowers(width, height int, spawnChance float64, v []mathutil.Point) *Flowers {
	var vertexOffset uint16 = 0

	actual := make([]mathutil.Point, 0)
	vertices := make([]ebiten.Vertex, 0)
	indices := make([]uint16, 0)

	for _, lp := range v {
		rn := rand.Float64()
		if rn > spawnChance {
			continue
		}

		a := math.Pi * 2 * rn

		left := lp.X - 10
		right := lp.X + 10

		top := lp.Y - 10
		bottom := lp.Y + 10

		topLeft := mathutil.NewPoint(left, top).RotateAround(lp, a)
		topRight := mathutil.NewPoint(right, top).RotateAround(lp, a)
		bottomLeft := mathutil.NewPoint(left, bottom).RotateAround(lp, a)
		bottomRight := mathutil.NewPoint(right, bottom).RotateAround(lp, a)

		actual = append(actual, lp)

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

	img := ebiten.NewImage(width, height)
	op := &ebiten.DrawTrianglesShaderOptions{}
	img.DrawTrianglesShader(vertices, indices, flowerShader, op)

	return &Flowers{
		Positions: actual,

		vertices: vertices,
		indices:  indices,

		img: img,
	}
}

func (l *Flowers) Draw(target *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		float64(target.Bounds().Dx())*0.5-float64(l.img.Bounds().Dx())*0.5,
		float64(target.Bounds().Dy())*0.5-float64(l.img.Bounds().Dy())*0.5,
	)
	target.DrawImage(l.img, opts)
}
