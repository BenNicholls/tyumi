package vec

import (
	"iter"
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

func (r Rect) String() string {
	return "{" + r.Coord.String() + " " + r.Dims.String() + "}"
}

func (r Rect) Translated(c Coord) Rect {
	return Rect{r.Coord.Add(c), r.Dims}
}

// Returns the coordinates of the 4 corners of the rect, starting in the top left and going clockwise.
func (r Rect) Corners() (corners [4]Coord) {
	corners[0] = r.Coord                             //TOPLEFT
	corners[1] = Coord{r.X + r.W - 1, r.Y}           //TOPRIGHT
	corners[2] = Coord{r.X + r.W - 1, r.Y + r.H - 1} //BOTTOMRIGHT
	corners[3] = Coord{r.X, r.Y + r.H - 1}           //BOTTOMLEFT

	return
}

// Returns an iterator producing a sequence of all Coords within the Rect r, starting in the top-left corner and
// proceeding to the right, going line by line (like how you'd read)
func EachCoordInArea(b Bounded) iter.Seq[Coord] {
	return func(yield func(Coord) bool) {
		r := b.Bounds()
		for i := range r.Area() {
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

// FindIntersectionRect calculates the intersection of two rectangularly-bound objects. if no intersection
// returns Rect{0,0,0,0}
func FindIntersectionRect(r1, r2 Bounded) (r Rect) {
	if !Intersects(r1, r2) {
		return
	}

	b1, b2 := r1.Bounds(), r2.Bounds()
	r.X, r.Y = max(b1.X, b2.X), max(b1.Y, b2.Y)
	r.W, r.H = min(b1.X+b1.W, b2.X+b2.W)-r.X, min(b1.Y+b1.H, b2.Y+b2.H)-r.Y

	return
}

// Intersects returns true if the two provided Bounded areas intersect
func Intersects(r1, r2 Bounded) bool {
	b1, b2 := r1.Bounds(), r2.Bounds()
	if b1.X >= b2.X+b2.W || b2.X >= b1.X+b1.W || b1.Y >= b2.Y+b2.H || b2.Y >= b1.Y+b1.H {
		return false
	}

	return true
}
