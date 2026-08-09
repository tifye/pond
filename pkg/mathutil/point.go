package mathutil

import (
	"fmt"
	"math"
)

const (
	Pi2             = math.Pi * 2
	Epsilon float64 = 0.000001
)

type Point struct {
	X float64
	Y float64
}

type Vector = Point
type Size = Point

func NewPoint(x, y float64) Point {
	return Point{x, y}
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

func (p Point) Limit(max float64) Point {
	magSqrd := p.MagnitudeSquared()

	if magSqrd < max*max {
		return p
	}

	scale := max/math.Sqrt(magSqrd) + Epsilon
	return p.MultiplyScalar(scale)
}

func (p Point) MagnitudeSquared() float64 {
	return p.X*p.X + p.Y*p.Y
}

func (p Point) Magnitude() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Point) Distance(o Point) float64 {
	temp := p.Subtract(o)
	return math.Sqrt(temp.X*temp.X + temp.Y*temp.Y)
}

func (p Point) DistanceSquared(o Point) float64 {
	temp := p.Subtract(o)
	return temp.X*temp.X + temp.Y*temp.Y
}

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

func (p Point) Rotate(deltaAngle float64) Point {
	newAngle := math.Atan2(p.Y, p.X) + deltaAngle
	p.X = math.Cos(newAngle)
	p.Y = math.Sin(newAngle)
	return p
}

func (p Point) RotateAround(o Point, deltaAngle float64) Point {
	var t Point

	p = p.Subtract(o)
	t.X = p.X*math.Cos(deltaAngle) - p.Y*math.Sin(deltaAngle)
	t.Y = p.X*math.Sin(deltaAngle) + p.Y*math.Cos(deltaAngle)

	return t.Add(o)
}

func (p Point) Follow(o Point, distance float64) Point {
	return p.Subtract(o).
		Normalize().
		MultiplyScalar(distance).
		Add(o)
}

func (p Point) AngleBetween(o Point) float64 {
	temp := p.Subtract(o)
	return math.Atan2(temp.Y, temp.X) + math.Pi
}

func (p Point) Angle() float64 {
	return math.Atan2(p.Y, p.Y)
}

func (p Point) String() string {
	return fmt.Sprintf("%b, %b", p.X, p.Y)
}
