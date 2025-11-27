package agent

import "github.com/tifye/pond/pkg/mathutil"

func Boundry(width, height, offset, force float64) BehaviourFunc {
	effectiveWidth := width - offset
	effectiveHeight := height - offset
	return func(agents *Agents, idx uint, lastTime, deltaTime float64) {
		position := agents.Position(idx)
		if position.X < offset {
			agents.Seek(idx, mathutil.Point{X: width, Y: position.Y}, force)
		} else if position.X > effectiveWidth {
			agents.Seek(idx, mathutil.Point{X: 0, Y: position.Y}, force)
		}

		if position.Y < offset {
			agents.Seek(idx, mathutil.Point{X: position.X, Y: height}, force)
		} else if position.Y > effectiveHeight {
			agents.Seek(idx, mathutil.Point{X: position.X, Y: 0}, force)
		}
	}
}
