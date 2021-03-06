package vec

import "github.com/bennicholls/tyumi/util"

//Bounded defines objects that can report a bounding box of some kind.
type Bounded interface {
	Bounds() Rect
}

//Rect is your standard rectangle object, with position (X,Y) in the top left corner.
type Rect struct {
	W, H int
	X, Y int
}

func (r Rect) Pos() (int, int) {
	return r.X, r.Y
}

func (r Rect) Dims() (int, int) {
	return r.W, r.H
}

//Goofy... but needed to satisfy Bounded interface.
func (r Rect) Bounds() Rect {
	return r
}

//translates r in 2d by the vector (dx, dy)
//CONSIDER: should input be a vec.coord?
func (r *Rect) Move(dx, dy int) {
	r.X += dx
	r.Y += dy
}

//Moves the rect to location (x, y)
func (r *Rect) MoveTo(x, y int) {
	r.X = x
	r.Y = y
}

//IsInside checks if the point (x, y) is within the object b.
func IsInside(x, y int, b Bounded) bool {
	r := b.Bounds()
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

//FindIntersectionRect calculates the intersection of two rectangularly-bound objects as a rect if no intersection,
//returns Rect{0,0,0,0}
func FindIntersectionRect(r1, r2 Bounded) (r Rect) {
	b1 := r1.Bounds()
	b2 := r2.Bounds()

	//check for intersection
	if b1.X >= b2.X+b2.W || b2.X >= b1.X+b1.W || b1.Y >= b2.Y+b2.H || b2.Y >= b1.Y+b1.H {
		return
	}

	r.X = util.Max(b1.X, b2.X)
	r.Y = util.Max(b1.Y, b2.Y)
	r.W = util.Min(b1.X+b1.W, b2.X+b2.W) - r.X
	r.H = util.Min(b1.Y+b1.H, b2.Y+b2.H) - r.Y

	return
}
