package mathutil

import "math"

const (
	Pi2 = math.Pi * 2
)

type Point struct {
	X float64
	Y float64
}

func (p Point) Subtract(o Point) Point {
	return Point{
		X: p.X - o.X,
		Y: p.Y - o.Y,
	}
}

func (p Point) Add(o Point) Point {
	return Point{
		X: p.X + o.X,
		Y: p.Y + o.Y,
	}
}

func (p Point) Distance(o Point) float64 {
	temp := p.Subtract(o)
	return math.Sqrt(temp.X*temp.X + temp.Y*temp.Y)
}

func (p Point) DistanceSquared(o Point) float64 {
	temp := p.Subtract(o)
	return temp.X*temp.X + temp.Y*temp.Y
}

// todo: optimize
func (p Point) Normalize() Point {
	sqr := math.Sqrt(p.X*p.X + p.Y*p.Y)
	return Point{
		X: p.X / sqr,
		Y: p.Y / sqr,
	}
}

func (p Point) MultiplyScalar(v float64) Point {
	return Point{
		X: p.X * v,
		Y: p.Y * v,
	}
}

func (p Point) Dot(o Point) float64 {
	return p.X*o.X + p.Y*o.Y
}

func (p Point) Cross(o Point) float64 {
	return p.X*o.Y - p.Y*o.X
}

func (p Point) RotateCounterClockwise() Point {
	return Point{
		X: -p.Y,
		Y: p.X,
	}
}

func (p Point) RotateClockwise() Point {
	return Point{
		X: p.Y,
		Y: -p.X,
	}
}

func (p Point) Follow(o Point, distance float64) Point {
	return p.Subtract(o).
		Normalize().
		MultiplyScalar(distance).
		Add(o)
}

// todo: optimize
func (p Point) AngleBetween(o Point) float64 {
	temp := p.Subtract(o)
	return math.Atan2(temp.Y, temp.X) + math.Pi
}
