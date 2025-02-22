package vec

import (
	"fmt"

	"github.com/bennicholls/tyumi/util"
)

var (
	ZERO_COORD Coord = Coord{0, 0}
)

// Coord is an (X, Y) pair that represents a spot on some 2d grid.
type Coord Vec2i

func (c *Coord) Move(dx, dy int) {
	c.X += dx
	c.Y += dy
}

func (c *Coord) MoveTo(x, y int) {
	c.X, c.Y = x, y
}

func (c1 Coord) Add(c2 Coord) Coord {
	return Coord{c1.X + c2.X, c1.Y + c2.Y}
}

func (c1 Coord) Subtract(c2 Coord) Coord {
	return Coord{c1.X - c2.X, c1.Y - c2.Y}
}

// returns the coordinate stepped once in the direction d
func (c Coord) Step(d Direction) Coord {
	return c.Add(Coord(d))
}

// Returns the coordinate stepped N times in the direction d
func (c Coord) StepN(d Direction, n int) Coord {
	return c.Add(Coord(d).Scale(n))
}

func (c Coord) Scale(scale int) Coord {
	return Coord{c.X * scale, c.Y * scale}
}

// ToIndex converts a Coord to a 1D index in a 2D array with the given stride.
func (c Coord) ToIndex(stride int) int {
	return c.Y*stride + c.X
}

// IsInside checks if the coord is within the bounds of object b.
func (c Coord) IsInside(b Bounded) bool {
	r := b.Bounds()
	return !(c.X < r.X || c.X >= r.X+r.W || c.Y < r.Y || c.Y >= r.Y+r.H)
}

func (c Coord) String() string {
	return fmt.Sprintf("(X: %d, Y: %d)", c.X, c.Y)
}

// ManhattanDistance calculates the manhattan (or taxicab) distance on a square grid.
func ManhattanDistance(c1, c2 Coord) int {
	return util.Abs(c2.X-c1.X) + util.Abs(c2.Y-c1.Y)
}

// IndexToCoord returns a coord representing an index from a 1D array representing a 2D grid with the given stride
func IndexToCoord(index, stride int) Coord {
	return Coord{index % stride, index / stride}
}

type Direction Vec2i

var (
	DIR_NONE  Direction = Direction{0, 0}
	DIR_UP    Direction = Direction{0, -1}
	DIR_DOWN  Direction = Direction{0, 1}
	DIR_LEFT  Direction = Direction{-1, 0}
	DIR_RIGHT Direction = Direction{1, 0}
)

var CardinalDirections [4]Direction = [4]Direction{DIR_UP, DIR_RIGHT, DIR_DOWN, DIR_LEFT}

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
