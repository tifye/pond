package agent

import (
	"math"
	"math/rand/v2"

	"github.com/tifye/pond/pkg/mathutil"
)

type Wander struct {
	force          float64
	pivotDistance  float64
	steerRadius    float64
	maxAngleChange float64

	angle  []float64
	center []mathutil.Point
	target []mathutil.Point
}

func NewWander(numAgents uint, force, pivotDistance, steerRadius, maxAngleChange float64) *Wander {
	return &Wander{
		force:          force,
		pivotDistance:  pivotDistance,
		steerRadius:    steerRadius,
		maxAngleChange: maxAngleChange,

		angle:  make([]float64, numAgents),
		center: make([]mathutil.Point, numAgents),
		target: make([]mathutil.Point, numAgents),
	}
}

func (w *Wander) Update(agents *Agents, idx uint, lastTime, deltaTime float64) {
	position := agents.Position(idx)
	velocity := agents.Velocity(idx)
	velocityMag := velocity.Magnitude() + mathutil.Epsilon

	wanderCenter := position.Add(velocity.MultiplyScalar(1 / velocityMag).MultiplyScalar(w.pivotDistance))
	w.center[idx].X = wanderCenter.X
	w.center[idx].Y = wanderCenter.Y

	// todo: inject random instance
	w.angle[idx] += (rand.Float64() - 0.5) * w.maxAngleChange * 2
	w.target[idx] = w.center[idx]
	// todo: does the compiler optimize away two slice accesses?
	w.target[idx].X += math.Cos(w.angle[idx]) * w.steerRadius
	w.target[idx].Y += math.Sin(w.angle[idx]) * w.steerRadius

	steerVec := w.target[idx].Subtract(position)
	steerMag := steerVec.Magnitude() + mathutil.Epsilon
	steerVec = steerVec.MultiplyScalar(w.force / steerMag)

	agents.ApplyForce(idx, steerVec)
}
