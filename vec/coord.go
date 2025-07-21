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

// Returns the coordinate stepped once in the direction d
func (c Coord) Step(d Direction) Coord {
	return c.Add(d.Coord())
}

// Returns the coordinate stepped N times in the direction d
func (c Coord) StepN(d Direction, n int) Coord {
	return c.Add(d.Coord().Scale(n))
}

// Returns the coordinate with both X and Y multiplied by scale
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
	return r.Contains(c)
}

// IsInPerimeter check if the coord lies in the perimeter of the bounded object b.
func (c Coord) IsInPerimeter(b Bounded) bool {
	r := b.Bounds()
	return c.X == r.X || c.X == r.X+r.W-1 || c.Y == r.Y || c.Y == r.Y+r.H-1
}

func (c Coord) String() string {
	return fmt.Sprintf("(X: %d, Y: %d)", c.X, c.Y)
}

// DistanceTo returns the euclidean distance bewteen c1 and c2. For comparisons, consider using DistanceSqTo instead,
// it will be much faster.
func (c1 Coord) DistanceTo(c2 Coord) float64 {
	return Vec2i(c1.Subtract(c2)).Mag()
}

// DistanceSqTo returns the euclidean distance between c1 and c2, squared. This is useful for comparing distances in
// cases where the actual distance is not important, because this is much faster than calculating the real distance.
func (c1 Coord) DistanceSqTo(c2 Coord) int {
	dx, dy := c1.X-c2.X, c1.Y-c2.Y
	return dx*dx + dy*dy
}

// ManhattanDistance calculates the manhattan (or taxicab) distance on a square grid.
func (c1 Coord) ManhattanDistanceTo(c2 Coord) int {
	return util.Abs(c2.X-c1.X) + util.Abs(c2.Y-c1.Y)
}

func (c Coord) Lerp(to Coord, val, steps int) Coord {
	return Coord{
		X: util.Lerp(c.X, to.X, val, steps),
		Y: util.Lerp(c.Y, to.Y, val, steps),
	}
}

// IndexToCoord returns a coord representing an index from a 1D array representing a 2D grid with the given stride
func IndexToCoord(index, stride int) Coord {
	return Coord{index % stride, index / stride}
}

type Direction int

const (
	DIR_UP Direction = iota
	DIR_UPRIGHT
	DIR_RIGHT
	DIR_DOWNRIGHT
	DIR_DOWN
	DIR_DOWNLEFT
	DIR_LEFT
	DIR_UPLEFT
	DIR_NONE
)

var directions []Coord = []Coord{
	{0, -1},  //Up
	{1, -1},  //Up Right
	{1, 0},   //Right
	{1, 1},   //Down Right
	{0, 1},   //Down
	{-1, 1},  //Down Left
	{-1, 0},  //Left
	{-1, -1}, //Up Left
	{0, 0},   //No Move
}

var Directions []Direction = []Direction{DIR_UP, DIR_UPRIGHT, DIR_RIGHT, DIR_DOWNRIGHT, DIR_DOWN, DIR_DOWNLEFT, DIR_LEFT, DIR_UPLEFT}
var CardinalDirections []Direction = []Direction{DIR_UP, DIR_RIGHT, DIR_DOWN, DIR_LEFT}

func (d Direction) Coord() Coord {
	return directions[int(d)]
}

func (d Direction) Inverted() Direction {
	return Direction(util.CycleClamp(int(d)+4, 0, 7))
}

func (d Direction) RotateCW() Direction {
	if d == DIR_NONE {
		return d
	}

	return Direction(util.CycleClamp(int(d)+1, 0, 7))
}

func (d Direction) RotateCCW() Direction {
	if d == DIR_NONE {
		return d
	}

	return Direction(util.CycleClamp(int(d)-1, 0, 7))
}

func (d Direction) RotateCW90() Direction {
	if d == DIR_NONE {
		return d
	}

	return Direction(util.CycleClamp(int(d)+2, 0, 7))
}

func (d Direction) RotateCCW90() Direction {
	if d == DIR_NONE {
		return d
	}

	return Direction(util.CycleClamp(int(d)-2, 0, 7))
}

func RandomDirection() Direction {
	return util.PickOne(Directions)
}

func RandomCardinalDirection() Direction {
	return util.PickOne(CardinalDirections)
}
