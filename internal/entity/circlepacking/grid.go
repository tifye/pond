package circlepacking

import "github.com/tifye/pond/pkg/mathutil"

// spatialGrid is a uniform grid used to accelerate neighbour lookups while
// packing. Circles are bucketed by their centre. As long as the cell size is at
// least the largest possible circle diameter, any two overlapping circles are
// guaranteed to sit in the same cell or in adjacent cells, so a 3x3 block query
// is enough to find every possible collision. This turns the per-circle overlap
// test from O(n) into roughly O(1).
type spatialGrid struct {
	cellSize float64
	cols     int
	rows     int
	cells    [][]*Circle
}

func newSpatialGrid(width, height int, cellSize float64) *spatialGrid {
	if cellSize < 1 {
		cellSize = 1
	}

	cols := int(float64(width)/cellSize) + 1
	rows := int(float64(height)/cellSize) + 1

	return &spatialGrid{
		cellSize: cellSize,
		cols:     cols,
		rows:     rows,
		cells:    make([][]*Circle, cols*rows),
	}
}

func (g *spatialGrid) coords(p mathutil.Point) (int, int) {
	cx := clampInt(int(p.X/g.cellSize), 0, g.cols-1)
	cy := clampInt(int(p.Y/g.cellSize), 0, g.rows-1)
	return cx, cy
}

func (g *spatialGrid) insert(c *Circle) {
	cx, cy := g.coords(c.Position)
	i := cy*g.cols + cx
	g.cells[i] = append(g.cells[i], c)
}

// neighbours calls fn for every circle in the 3x3 block of cells surrounding p,
// stopping early (and returning true) as soon as fn returns true.
func (g *spatialGrid) neighbours(p mathutil.Point, fn func(*Circle) bool) bool {
	cx, cy := g.coords(p)

	for dy := -1; dy <= 1; dy++ {
		ny := cy + dy
		if ny < 0 || ny >= g.rows {
			continue
		}

		row := ny * g.cols
		for dx := -1; dx <= 1; dx++ {
			nx := cx + dx
			if nx < 0 || nx >= g.cols {
				continue
			}

			for _, c := range g.cells[row+nx] {
				if fn(c) {
					return true
				}
			}
		}
	}

	return false
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
