package vec

import "github.com/bennicholls/tyumi/util"

var (
	ZERO_COORD Coord = Coord{0, 0}
)

//Coord is an (X, Y) pair that represents a spot on some 2d grid.
type Coord Vec2i

func (c *Coord) Move(dx, dy int) {
	c.X += dx
	c.Y += dy
}

func (c *Coord) MoveTo(x, y int) {
	c.X, c.Y = x, y
}

//TODO: work around this somehow. Generic add function for vec2i types? Interface stuff?
func (c1 Coord) Add(c2 Coord) Coord {
	return Coord(Vec2i(c1).Add(Vec2i(c2)))
}

func (c1 Coord) Subtract(c2 Coord) Coord {
	return Coord{c1.X - c2.X, c1.Y - c2.Y}
}

//ManhattanDistance calculates the manhattan (or taxicab) distance on a square grid.
func ManhattanDistance(c1, c2 Coord) int {
	return util.Abs(c2.X-c1.X) + util.Abs(c2.Y-c1.Y)
}

//TODO this should be somewhere else...
type Dims struct {
	W, H int
}

func (d Dims) Area() int {
	return d.W * d.H
}
