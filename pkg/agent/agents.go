package agent

import (
	"slices"

	"github.com/tifye/pond/internal/app/assert"
	"github.com/tifye/pond/pkg/mathutil"
)

type Agents struct {
	len, cap           uint
	MaxSpeed, MaxForce float64
	position           []mathutil.Point
	velocity           []mathutil.Point
	acceleration       []mathutil.Point

	behaviours []Behaviour
}

type Behaviour interface {
	Update(agents *Agents, idx uint, lastTime, deltaTime float64)
}

type BehaviourFunc func(agents *Agents, idx uint, lastTime, deltaTime float64)

func (f BehaviourFunc) Update(agents *Agents, idx uint, lastTime, deltaTime float64) {
	f(agents, idx, lastTime, deltaTime)
}

func NewAgents(num uint, cap ...uint) *Agents {
	assert.Assert(len(cap) < 2, "expected only one")

	var capacity uint = num
	if len(cap) > 0 {
		capacity = cap[0]
	}

	return &Agents{
		len:      num,
		cap:      uint(capacity),
		MaxSpeed: 7,
		MaxForce: 3,
		// todo: test whether to go for three separate or just one
		position:     make([]mathutil.Point, num, capacity),
		velocity:     make([]mathutil.Point, num, capacity),
		acceleration: make([]mathutil.Point, num, capacity),

		behaviours: make([]Behaviour, 0),
	}
}

func (agents *Agents) AddBehaviour(b Behaviour) {
	agents.behaviours = append(agents.behaviours, b)
}

func (agents *Agents) Update(lastTime, deltaTime float64) {
	for _, b := range agents.behaviours {
		for i := range agents.len {
			// todo: does convert betwen in and uint affect performance?
			b.Update(agents, i, lastTime, deltaTime)
		}
	}

	for i := range agents.len {
		agents.velocity[i] = agents.velocity[i].Add(agents.acceleration[i])
		agents.velocity[i] = agents.velocity[i].Limit(agents.MaxSpeed)
		agents.position[i] = agents.position[i].Add(agents.velocity[i])
		agents.acceleration[i].X = 0
		agents.acceleration[i].Y = 0
	}
}

func (agents *Agents) ApplyForce(idx uint, force mathutil.Vector) {
	force = force.Limit(agents.MaxForce)
	agents.acceleration[idx] = agents.acceleration[idx].Add(force)
}

func (agents *Agents) Position(idx uint) mathutil.Point {
	return agents.position[idx]
}

func (agents *Agents) Set(idx uint, position, velocity, acceleration mathutil.Point) {
	agents.position[idx] = position
	agents.velocity[idx] = velocity
	agents.acceleration[idx] = acceleration
}

func (agents *Agents) Velocity(idx uint) mathutil.Point {
	return agents.velocity[idx]
}

func (agents *Agents) Acceleration(idx uint) mathutil.Point {
	return agents.acceleration[idx]
}

func (agents *Agents) Num() uint {
	return agents.len
}

type GrowAbleBehaviour interface {
	Grow(n int)
}

// Grow increases the agents capacity to support at least
// n agents.
func (agents *Agents) Grow(n int) {
	assert.Assert(n >= int(agents.len), "cannot grow to a smaller size")

	agents.cap = uint(n)
	agents.position = slices.Grow(agents.position, n)
	agents.velocity = slices.Grow(agents.velocity, n)
	agents.acceleration = slices.Grow(agents.acceleration, n)

	for _, b := range agents.behaviours {
		growable, ok := b.(GrowAbleBehaviour)
		if ok {
			growable.Grow(n)
		}
	}
}

func (agents *Agents) Seek(idx uint, target mathutil.Point, strength float64) {
	desired := target.Subtract(agents.Position(idx)).
		Normalize().
		MultiplyScalar(agents.MaxSpeed)
	steer := desired.Subtract(agents.Velocity(idx))
	agents.ApplyForce(idx, steer.MultiplyScalar(strength))
}
