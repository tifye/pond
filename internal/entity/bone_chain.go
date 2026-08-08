package entity

import (
	"math"

	"github.com/tifye/pond/pkg/mathutil"
	"github.com/tifye/pond/pkg/mathutil/fabrik"
)

type BoneChainUpdater interface {
	Update(bc *BoneChain)
}

type BoneChain struct {
	Position mathutil.Point

	Joints  []mathutil.Point
	Lengths []float64

	Updater BoneChainUpdater
}

func (bc *BoneChain) Update() {
	bc.Updater.Update(bc)
}

func (bc *BoneChain) Curvature() float64 {
	dotSum := 0.0
	for i := 0; i < len(bc.Joints)-2; i++ {
		var (
			a = bc.Joints[i]
			b = bc.Joints[i+1]
			c = bc.Joints[i+2]
		)

		ba := a.Subtract(b).Normalize()
		cb := b.Subtract(c).Normalize()

		dotSum += ba.Cross(cb)
	}

	return dotSum
}

type BoneChainFABRIKUpdater struct {
	MinAngle float64
}

func (b BoneChainFABRIKUpdater) Update(bc *BoneChain) {
	fabrik.SolveFABRIK(
		bc.Joints,
		bc.Lengths,
		bc.Position,
		b.MinAngle,
	)
}

// BoneChainLineUpdater lays the chain out as a straight horizontal line
// trailing behind the chain's Position, spacing each joint by its segment
// length.
type BoneChainLineUpdater struct{}

func (BoneChainLineUpdater) Update(bc *BoneChain) {
	if len(bc.Joints) == 0 {
		return
	}

	bc.Joints[0] = bc.Position
	for i := 1; i < len(bc.Joints); i++ {
		bc.Joints[i] = mathutil.Point{
			Y: bc.Joints[i-1].Y + bc.Lengths[i-1],
			X: bc.Position.X,
		}
	}
}

// BoneChainSineUpdater lays the chain out along a sine wave trailing behind
// the chain's Position. X advances by each segment length while Y oscillates.
// The wave travels over time via an internal phase advanced by Speed each
// update; set Speed to 0 for a static wave.
type BoneChainSineUpdater struct {
	Amplitude float64 // peak height of the wave, in pixels
	Frequency float64 // wave cycles per unit of length
	Speed     float64 // phase advanced per update (0 = static)

	phase float64
}

func (u *BoneChainSineUpdater) Update(bc *BoneChain) {
	if len(bc.Joints) == 0 {
		return
	}

	u.phase += u.Speed

	bc.Joints[0] = bc.Position
	dist := 0.0
	for i := 1; i < len(bc.Joints); i++ {
		dist += bc.Lengths[i-1]
		bc.Joints[i] = mathutil.Point{
			X: bc.Position.X - dist,
			Y: bc.Position.Y + math.Sin(dist*u.Frequency+u.phase)*u.Amplitude,
		}
	}
}
