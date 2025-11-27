package fabrik

import (
	"math"

	"github.com/tifye/pond/pkg/mathutil"
)

func SolveFABRIK(
	points []mathutil.Point,
	segmentLengths []float64,
	target mathutil.Point,
	minAngle float64,
) {
	n := len(points)
	if n < 3 {
		return
	}
	if n != len(segmentLengths)+1 {
		return
	}

	points[0] = target
	points[1] = points[1].Follow(points[0], segmentLengths[1])
	for i := 2; i < n; i++ {
		points[i], points[i-1] = jointWithMinAngle(
			points[i], points[i-1], points[i-2],
			segmentLengths[i-1], segmentLengths[i-2],
			minAngle,
		)
	}

	// points[n-1] = target
	// points[n-2] = points[n-2].Follow(points[n-1], segmentLengths[n-2])
	// for i := n - 3; i >= 0; i-- {
	// 	points[i], points[i+1] = jointWithMinAngle(
	// 		points[i], points[i+1], points[i+2],
	// 		segmentLengths[i], segmentLengths[i+1],
	// 		minAngle,
	// 	)
	// }
}

func jointWithMinAngle(
	a, b, c mathutil.Point,
	abLen, bcLen float64,
	minAngle float64,
) (newA, newB mathutil.Point) {
	b = b.Follow(c, bcLen)
	bcAngle := b.AngleBetween(c)

	ab := b.Subtract(a)
	bc := c.Subtract(b)

	dot := ab.Dot(bc)
	cross := ab.Cross(bc)
	// todo: optimize
	angle := math.Atan2(cross, dot) + math.Pi

	maxAngle := mathutil.Pi2 - minAngle

	// todo: make branchless
	if angle < minAngle {
		newAngle := bcAngle - minAngle
		return mathutil.Point{
			X: math.Cos(newAngle)*abLen + b.X,
			Y: math.Sin(newAngle)*abLen + b.Y,
		}, b
	} else if angle > maxAngle {
		newAngle := bcAngle - maxAngle
		return mathutil.Point{
			X: math.Cos(newAngle)*abLen + b.X,
			Y: math.Sin(newAngle)*abLen + b.Y,
		}, b
	}

	return a.Follow(b, abLen), b
}
