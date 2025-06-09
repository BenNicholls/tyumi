package vec

import (
	"iter"
	"math/rand"
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

func (r Rect) Translated(coord Coord) Rect {
	return Rect{r.Coord.Add(coord), r.Dims}
}

// Returns the coordinates of the 4 corners of the rect, starting in the top left and going clockwise.
func (r Rect) Corners() (corners [4]Coord) {
	corners[0] = r.Coord                             //TOPLEFT
	corners[1] = Coord{r.X + r.W - 1, r.Y}           //TOPRIGHT
	corners[2] = Coord{r.X + r.W - 1, r.Y + r.H - 1} //BOTTOMRIGHT
	corners[3] = Coord{r.X, r.Y + r.H - 1}           //BOTTOMLEFT

	return
}

// Returns the center of the rect. Since we're all integers 'round these parts this won't be exact unless both width
// and height are odd numbers, so be aware.
func (r Rect) Center() Coord {
	return Coord{r.X + r.W/2, r.Y + r.H/2}
}

// Contains calculates whether the provided coord is within the bounds of the rect.
func (r Rect) Contains(c Coord) bool {
	return !(c.X < r.X || c.Y < r.Y || c.X >= r.X+r.W || c.Y >= r.Y+r.H)
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

// Returns an iterator producing a sequence of all Coords that are contained within the intersection of all provided
// bounded areas.
func EachCoordInIntersection(areas ...Bounded) iter.Seq[Coord] {
	switch len(areas) {
	case 0:
		return EachCoordInArea(Rect{ZERO_COORD, Dims{0, 0}})
	case 1:
		return EachCoordInArea(areas[0])
	default:
		intersection := areas[0].Bounds()
		for i := 1; i < len(areas); i++ {
			if intersection.Area() == 0 {
				break
			}
			intersection = FindIntersectionRect(intersection, areas[i])
		}
		return EachCoordInArea(intersection)
	}
}

// Returns an iterator producing a sequence of all coords in the perimeter of a bounded area.
func EachCoordInPerimeter(b Bounded) iter.Seq[Coord] {
	return func(yield func(Coord) bool) {
		r := b.Bounds()
		if r.Area() == 0 { //0x0 box, aka a nothing
			return
		}

		if r.Area() == 1 { //1x1 box, aka 1 cell
			yield(r.Coord)
			return
		}

		corners := r.Corners()

		if r.W == 1 || r.H == 1 { // 1D box, aka a line.
			line := Line{corners[0], corners[2]}
			for coord := range line.EachCoord() {
				if !yield(coord) {
					return
				}
			}
			return
		}

		var sides [4]Line
		sides[0] = Line{r.Coord, corners[1].Step(DIR_LEFT)}     // top
		sides[1] = Line{corners[1], corners[2].Step(DIR_UP)}    // right
		sides[2] = Line{corners[2], corners[3].Step(DIR_RIGHT)} // bottom
		sides[3] = Line{corners[3], r.Coord.Step(DIR_DOWN)}     // right
		for _, side := range sides {
			for coord := range side.EachCoord() {
				if !yield(coord) {
					return
				}
			}
		}
	}
}

// FindIntersectionRect calculates the intersection of two rectangularly-bound objects. if no intersection
// returns Rect{0,0,0,0}
func FindIntersectionRect(b1, b2 Bounded) (r Rect) {
	if !Intersects(b1, b2) {
		return
	}

	r1, r2 := b1.Bounds(), b2.Bounds()
	r.X, r.Y = max(r1.X, r2.X), max(r1.Y, r2.Y)
	r.W, r.H = min(r1.X+r1.W, r2.X+r2.W)-r.X, min(r1.Y+r1.H, r2.Y+r2.H)-r.Y

	return
}

// Intersects returns true if the two provided Bounded areas intersect
func Intersects(b1, b2 Bounded) bool {
	r1, r2 := b1.Bounds(), b2.Bounds()
	if r1.Area() == 0 || r2.Area() == 0 { // zero-sized rects cannot intersect with anything
		return false
	}

	if r1.X >= r2.X+r2.W || r2.X >= r1.X+r1.W || r1.Y >= r2.Y+r2.H || r2.Y >= r1.Y+r1.H {
		return false
	}

	return true
}

func RandomCoordInArea(area Bounded) (c Coord) {
	r := area.Bounds()
	c.X = rand.Intn(r.W) + r.X
	c.Y = rand.Intn(r.H) + r.Y

	return
}
