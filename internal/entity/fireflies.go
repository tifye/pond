package entity

import (
	"math"
	"math/rand/v2"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/pkg/agent"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed firefly.kage
var fireflyShaderSrc []byte
var fireflyShader *ebiten.Shader

func init() {
	if s, err := ebiten.NewShader(fireflyShaderSrc); err != nil {
		panic(err)
	} else {
		fireflyShader = s
	}
}

type FireFlies struct {
	agents *agent.Agents

	sprite *ebiten.Image
}

func NewFireFlies(n uint, width, height float64) *FireFlies {
	agents := agent.NewAgents(n, n)
	agents.MaxSpeed = 3
	agents.AddBehaviours(
		agent.NewWander(n, 1, 50, 50, math.Pi*1.5),
		agent.Boundry(width, height, 200, 0.05),
	)

	for i := range agents.Num() {
		rx := rand.IntN(int(width))
		ry := rand.IntN(int(height))
		agents.Set(i,
			mathutil.NewPoint(float64(rx), float64(ry)),
			mathutil.Point{},
			mathutil.Point{},
		)
	}

	sprite := ebiten.NewImage(50, 50)
	sprite.DrawTrianglesShader(
		[]ebiten.Vertex{
			{DstX: 0, DstY: 0, SrcX: 0, SrcY: 0},
			{DstX: 50, DstY: 0, SrcX: 1, SrcY: 0},
			{DstX: 50, DstY: 50, SrcX: 1, SrcY: 1},
			{DstX: 0, DstY: 50, SrcX: 0, SrcY: 1},
		},
		[]uint16{
			0, 1, 3, 1, 2, 3,
		},
		fireflyShader,
		nil,
	)

	return &FireFlies{
		agents: agents,
		sprite: sprite,
	}
}

func (f *FireFlies) Update() {
	elapsed := float64(ebiten.Tick()) / float64(ebiten.TPS())
	delta := 1 / float64(ebiten.TPS())

	f.agents.Update(elapsed, delta)
}

func (f *FireFlies) Draw(target *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	for i := range f.agents.Num() {
		pos := f.agents.Position(i)

		opts.GeoM.Reset()
		opts.GeoM.Translate(pos.X-25, pos.Y-25)
		target.DrawImage(f.sprite, opts)
	}
}
