package vec

import (
	"iter"

	"github.com/bennicholls/tyumi/util"
)

// Bounded defines objects that can report a bounding box of some kind.
type Bounded interface {
	Bounds() Rect
}

// Rect is your standard rectangle object, with position (X,Y) in the top left corner.
type Rect struct {
	Coord //position (X, Y)
	Dims  //size (W, H)
}

// Goofy... but needed to satisfy Bounded interface.
func (r Rect) Bounds() Rect {
	return r
}

// Returns an iterator producing a sequence of all Coords within the Rect r.
func EachCoord(b Bounded) iter.Seq[Coord] {
	return func(yield func(Coord) bool) {
		r := b.Bounds()
		for i := 0; i < r.Area(); i++ {
			if !yield(Coord{r.X + (i % r.W), r.Y + (i / r.W)}) {
				return
			}
		}
	}
}

// IsInside checks if the point (x, y) is within the bounds of object b.
func IsInside(x, y int, b Bounded) bool {
	r := b.Bounds()
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

// FindIntersectionRect calculates the intersection of two rectangularly-bound objects as a rect if no intersection,
// returns Rect{0,0,0,0}
func FindIntersectionRect(r1, r2 Bounded) (r Rect) {
	b1, b2 := r1.Bounds(), r2.Bounds()

	//check for intersection
	if b1.X >= b2.X+b2.W || b2.X >= b1.X+b1.W || b1.Y >= b2.Y+b2.H || b2.Y >= b1.Y+b1.H {
		return
	}

	r.X, r.Y = util.Max(b1.X, b2.X), util.Max(b1.Y, b2.Y)
	r.W, r.H = util.Min(b1.X+b1.W, b2.X+b2.W) - r.X, util.Min(b1.Y+b1.H, b2.Y+b2.H) - r.Y

	return
}
