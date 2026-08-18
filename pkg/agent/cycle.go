package agent

import (
	"math/rand/v2"

	"github.com/tifye/pond/pkg/mathutil"
)

type CycleStaticTargets struct {
	Targets        []mathutil.Point
	Force          float64
	MinDuraction   float64
	DuractionRange float64

	agentTargets []mathutil.Point

	// Elapsed time of each agent
	agentDurations []float64
	agentRelatedAt []float64
}

func NewCycleStaticTargets(numAgents int, targets []mathutil.Point) *CycleStaticTargets {
	return &CycleStaticTargets{
		Targets:        targets,
		agentTargets:   make([]mathutil.Point, numAgents),
		agentDurations: make([]float64, numAgents),
		agentRelatedAt: make([]float64, numAgents),
		MinDuraction:   3,
		DuractionRange: 5,
		Force:          0.15,
	}
}

func (c *CycleStaticTargets) Update(agents *Agents, idx uint, lastTime, deltaTime float64) {
	timeSinceTargetChange := (lastTime + deltaTime) - c.agentRelatedAt[idx]
	if timeSinceTargetChange > c.agentDurations[idx] {
		c.agentDurations[idx] = c.MinDuraction + rand.Float64()*c.DuractionRange
		c.agentRelatedAt[idx] = (lastTime + deltaTime)
		c.agentTargets[idx] = c.Targets[rand.IntN(len(c.Targets))]
	}

	target := c.agentTargets[idx]
	agents.Arrive(idx, target, 100, c.Force)
}
