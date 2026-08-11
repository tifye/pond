package lilypad

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed lilypad.kage
var lilypadShaderSrc []byte
var lilpadShader *ebiten.Shader

func init() {
	if s, err := ebiten.NewShader(lilypadShaderSrc); err != nil {
		panic(err)
	} else {
		lilpadShader = s
	}
}

type Lilypads struct {
	Positions []mathutil.Point

	vertices []ebiten.Vertex
	indices  []uint16

	img       *ebiten.Image
	translate mathutil.Vector
}

func New(width, height int, v []mathutil.Point, sizes []int, angles []float64) *Lilypads {
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

	img := ebiten.NewImage(width, height)
	op := &ebiten.DrawTrianglesShaderOptions{}
	img.DrawTrianglesShader(vertices, indices, lilpadShader, op)

	return &Lilypads{
		Positions: v,

		vertices: vertices,
		indices:  indices,

		img: img,
	}
}

func (l *Lilypads) Draw(target *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		float64(target.Bounds().Dx())*0.5-float64(l.img.Bounds().Dx())*0.5,
		float64(target.Bounds().Dy())*0.5-float64(l.img.Bounds().Dy())*0.5,
	)
	target.DrawImage(l.img, opts)
}
