package agent

import (
	"github.com/tifye/pond/internal/app/assert"
	"github.com/tifye/pond/pkg/mathutil"
)

type Agents struct {
	len, cap           uint
	maxSpeed, maxForce float64
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

	var capacity uint = 0
	if len(cap) > 0 {
		capacity = cap[0]
	}

	return &Agents{
		len:      num,
		cap:      uint(capacity),
		maxSpeed: 10,
		maxForce: 3,
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
		agents.velocity[i] = agents.velocity[i].Limit(agents.maxSpeed)
		agents.position[i] = agents.position[i].Add(agents.velocity[i])
		agents.acceleration[i].X = 0
		agents.acceleration[i].Y = 0
	}
}

func (agents *Agents) ApplyForce(idx uint, force mathutil.Vector) {
	force = force.Limit(agents.maxForce)
	agents.acceleration[idx] = agents.acceleration[idx].Add(force)
}

func (agents *Agents) Position(idx uint) mathutil.Point {
	return agents.position[idx]
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

func (agents *Agents) Seek(idx uint, target mathutil.Point, strength float64) {
	desired := target.Subtract(agents.Position(idx)).
		Normalize().
		MultiplyScalar(agents.maxSpeed)
	steer := desired.Subtract(agents.Velocity(idx))
	agents.ApplyForce(idx, steer.MultiplyScalar(strength))
}
