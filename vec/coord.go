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

func (c Coord) Step(d Direction) Coord {
	return c.Add(Coord(d))
}

//ManhattanDistance calculates the manhattan (or taxicab) distance on a square grid.
func ManhattanDistance(c1, c2 Coord) int {
	return util.Abs(c2.X-c1.X) + util.Abs(c2.Y-c1.Y)
}


type Direction Vec2i

var (
	DIR_NONE  Direction = Direction{0, 0}
	DIR_UP    Direction = Direction{0, -1}
	DIR_DOWN  Direction = Direction{0, 1}
	DIR_LEFT  Direction = Direction{-1, 0}
	DIR_RIGHT Direction = Direction{1, 0}
)

func (d Direction) Inverted() Direction {
	return Direction{-d.X, -d.Y}
}

func (d Direction) RotateCW() Direction {
	switch d {
	case DIR_UP:
		return DIR_RIGHT
	case DIR_RIGHT:
		return DIR_DOWN
	case DIR_DOWN:
		return DIR_LEFT
	case DIR_LEFT:
		return DIR_UP
	}
	return DIR_NONE
}

func (d Direction) RotateCCW() Direction {
	switch d {
	case DIR_UP:
		return DIR_LEFT
	case DIR_RIGHT:
		return DIR_UP
	case DIR_DOWN:
		return DIR_RIGHT
	case DIR_LEFT:
		return DIR_DOWN
	}
	return DIR_NONE
}

//TODO this should be somewhere else...
type Dims struct {
	W, H int
}

func (d Dims) Area() int {
	return d.W * d.H
}
