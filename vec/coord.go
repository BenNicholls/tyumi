package vec

import "github.com/bennicholls/tyumi/util"

var (
	ZERO_COORD Coord = Coord{0, 0}
)

//Coord is an (X, Y) pair that represents a spot on some 2d grid. Effectively just an implenetation of
//vec.2D using ints.
type Coord struct {
	X, Y int
}

func (c Coord) Get() (int, int) {
	return c.X, c.Y
}

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

func (c1 Coord) Sub(c2 Coord) Coord {
	return Coord{c1.X - c2.X, c1.Y - c2.Y}
}

func (c Coord) Mag() int {
	return int(c.ToVector().Mag())
}

func (c Coord) ToVector() Vec2 {
	return Vec2{float64(c.X), float64(c.Y)}
}

//DistanceSquared calculates the distance squared (sqrt unnecessary usually)
func DistanceSquared(c1, c2 Coord) int {
	return (c1.X-c2.X)*(c1.X-c2.X) + (c1.Y-c2.Y)*(c1.Y-c2.Y)
}

//ManhattanDistance calculates the manhattan (or taxicab) distance on a square grid.
func ManhattanDistance(c1, c2 Coord) int {
	return util.Abs(c2.X-c1.X) + util.Abs(c2.Y-c1.Y)
}
