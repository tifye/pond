package circlepacking

import (
	"math"
	"math/rand/v2"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/pkg/mathutil"
)

const (
	// spawnAttemptsPerTick is how many new circles we try to seed each Update.
	// Seeding a big batch per tick front-loads the placement so most circles
	// start growing together and reach similar sizes, instead of a few early
	// circles growing large and later ones squeezing into the gaps as tiny
	// circles. Turn this up if you still see too many small circles.
	spawnAttemptsPerTick = 64

	// spawnClearanceFactor sets how much empty room (as a fraction of maxRadius)
	// a new seed needs from its neighbours. Higher = fewer, more uniform circles
	// but a looser packing; lower = denser packing with more small circles.
	spawnClearanceFactor = 0.15
)

type Packing struct {
	all     []*Circle
	growing []*Circle

	Image  *ebiten.Image
	Width  int
	Height int

	grid *spatialGrid
	mask []byte

	// maxRadius caps how large a circle may grow and minRadius is the amount of
	// clearance a new seed needs. Together they keep the resulting sizes fairly
	// uniform: no runaway early circles and no tiny squeezed ones.
	maxRadius float64
	minRadius float64
}

func New(image *ebiten.Image, maxRadius float64) *Packing {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	if maxRadius < 1 {
		maxRadius = 1
	}
	minRadius := maxRadius * spawnClearanceFactor

	return &Packing{
		all:       make([]*Circle, 0),
		growing:   make([]*Circle, 0),
		Image:     image,
		Width:     width,
		Height:    height,
		grid:      newSpatialGrid(width, height, maxRadius*2),
		maxRadius: maxRadius,
		minRadius: minRadius,
	}
}

func (c *Packing) Update() {
	c.ensureMask()

	for range spawnAttemptsPerTick {
		c.trySpawn()
	}

	for _, circle := range c.growing {
		circle.Grow()
	}

	c.growing = slices.DeleteFunc(c.growing, c.shouldStopGrowing)
}

// trySpawn attempts to seed a single circle at a random free location.
func (c *Packing) trySpawn() {
	x := rand.IntN(c.Width)
	y := rand.IntN(c.Height)

	if !c.canPlaceAt(x, y) {
		return
	}

	pos := mathutil.Point{X: float64(x), Y: float64(y)}

	// Reject candidates that don't have room to reach at least minRadius. This
	// avoids spawning tiny circles wedged between existing ones and keeps the
	// overall size distribution tighter.
	crowded := c.grid.neighbours(pos, func(o *Circle) bool {
		return pos.Distance(o.Position) < o.Radius+c.minRadius
	})
	if crowded {
		return
	}

	circle := &Circle{Position: pos}
	c.all = append(c.all, circle)
	c.growing = append(c.growing, circle)
	c.grid.insert(circle)
}

func (c *Packing) Draw(target *ebiten.Image) {
	for _, circle := range c.all {
		circle.Draw(target)
	}
}

func (c *Packing) shouldStopGrowing(x *Circle) bool {
	if x.Radius >= c.maxRadius {
		x.Radius = c.maxRadius
		return true
	}

	switch {
	case x.left() < 0,
		x.right() > float64(c.Width),
		x.top() < 0,
		x.bottom() > float64(c.Height):
		return true
	}

	return c.isOverlapping(x)
}

func (c *Packing) isOverlapping(x *Circle) bool {
	return c.grid.neighbours(x.Position, func(o *Circle) bool {
		return o != x && o.isOverlapping(x)
	})
}

// ensureMask lazily reads the image into a CPU-side placement mask on first use.
// The image is treated as static, so this happens once and avoids the expensive
// per-frame GPU->CPU sync of Image.At. It's done here rather than in New because
// ReadPixels is only valid once the game has started; New may run during setup
// (before the loop), whereas Update never does.
func (c *Packing) ensureMask() {
	if c.mask != nil {
		return
	}

	c.mask = make([]byte, 4*c.Width*c.Height)
	c.Image.ReadPixels(c.mask)
}

func (c *Packing) canPlaceAt(x, y int) bool {
	return c.mask[(y*c.Width+x)*4] == 0
}

// Compile returns the current packing as parallel slices of circle positions
// and their integer radii, ready to be fed into the lilypad generator. Circles
// that round down to a radius below 1 are skipped.
func (c *Packing) Compile() ([]mathutil.Point, []int) {
	positions := make([]mathutil.Point, 0, len(c.all))
	sizes := make([]int, 0, len(c.all))

	for _, circle := range c.all {
		r := int(math.Round(circle.Radius))
		if r < 1 {
			continue
		}

		positions = append(positions, circle.Position)
		sizes = append(sizes, r)
	}

	return positions, sizes
}
