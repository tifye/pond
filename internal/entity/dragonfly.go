package entity

import (
	_ "embed"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/pkg/agent"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed dragonfly.kage
var dragonflyShaderSrc []byte
var DragonflyShader *ebiten.Shader

func init() {
	if s, err := ebiten.NewShader(dragonflyShaderSrc); err != nil {
		panic(err)
	} else {
		DragonflyShader = s
	}
}

type DragonFlies struct {
	agents *agent.Agents

	sprite *ebiten.Image
	size   mathutil.Size
}

func NewDragonFlies(n uint, width, height float64, targets []mathutil.Point) *DragonFlies {
	agents := agent.NewAgents(n, n)
	agents.MaxSpeed = 6
	agents.AddBehaviours(
		// agent.NewWander(n, 1, 50, 250, math.Pi*1.5),
		// agent.Boundry(width, height, 200, 0.05),
		agent.NewCycleStaticTargets(int(agents.Num()), targets),
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

	return &DragonFlies{
		agents: agents,
		sprite: NewDragonFlySprite(50, 50),
		size:   mathutil.NewPoint(50, 50),
	}
}

func (f *DragonFlies) Update() {
	elapsed := float64(ebiten.Tick()) / float64(ebiten.TPS())
	delta := 1 / float64(ebiten.TPS())

	f.agents.Update(elapsed, delta)
}

func (f *DragonFlies) Draw(target *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	for i := range f.agents.Num() {
		pos := f.agents.Position(i)
		vel := f.agents.Velocity(i)

		opts.GeoM.Reset()

		opts.GeoM.Translate(-f.size.X*0.5, -f.size.Y*0.5)
		opts.GeoM.Rotate(vel.Angle() + math.Pi*0.5)
		opts.GeoM.Translate(pos.X, pos.Y)

		target.DrawImage(f.sprite, opts)
	}
}

func NewDragonFlySprite(width, height int) *ebiten.Image {
	fwidth, fheight := float32(width), float32(height)

	sprite := ebiten.NewImage(width, height)
	sprite.DrawTrianglesShader(
		[]ebiten.Vertex{
			{DstX: 0, DstY: 0, SrcX: 0, SrcY: 0},
			{DstX: fwidth, DstY: 0, SrcX: 1, SrcY: 0},
			{DstX: fwidth, DstY: fheight, SrcX: 1, SrcY: 1},
			{DstX: 0, DstY: fheight, SrcX: 0, SrcY: 1},
		},
		[]uint16{
			0, 1, 3, 1, 2, 3,
		},
		DragonflyShader,
		nil,
	)

	return sprite
}
