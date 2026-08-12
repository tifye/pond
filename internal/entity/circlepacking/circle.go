package circlepacking

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tifye/pond/colors"
	"github.com/tifye/pond/pkg/mathutil"
)

type Circle struct {
	Position mathutil.Point
	Radius   float64
}

func (c *Circle) Draw(target *ebiten.Image) {
	vector.FillCircle(
		target,
		float32(c.Position.X),
		float32(c.Position.Y),
		float32(c.Radius),
		colors.Violet600,
		false,
	)
}

func (c *Circle) Grow() {
	c.Radius += 1.5
}

func (c *Circle) left() float64 {
	return c.Position.X - c.Radius
}

func (c *Circle) right() float64 {
	return c.Position.X + c.Radius
}

func (c *Circle) top() float64 {
	return c.Position.Y - c.Radius
}

func (c *Circle) bottom() float64 {
	return c.Position.Y + c.Radius
}

func (c *Circle) isOverlapping(o *Circle) bool {
	return c.Position.Subtract(o.Position).Magnitude() < o.Radius+c.Radius
}

func (c *Circle) isInside(p mathutil.Point) bool {
	return c.Position.Subtract(p).Magnitude() < c.Radius
}
