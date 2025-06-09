package vec

import (
	"iter"
	"math/rand"
)

// Bounded defines objects that can report a bounding box of some kind.
type Bounded interface {
	Bounds() Rect
}

// Returns an iterator producing a sequence of all Coords within the bounded object b, starting in the top-left corner
// and proceeding to the right, going line by line (like how you'd read)
func EachCoordInArea(b Bounded) iter.Seq[Coord] {
	return b.Bounds().EachCoord()
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
			intersection = intersection.Intersection(areas[i].Bounds())
		}
		return intersection.EachCoord()
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

// Intersects returns true if the two provided Bounded areas intersect
func Intersects(b1, b2 Bounded) bool {
	return b1.Bounds().Intersects(b2.Bounds())
}

// FindIntersectionRect calculates the intersection of two rectangularly-bound objects. if no intersection
// returns Rect{0,0,0,0}
func FindIntersectionRect(b1, b2 Bounded) (r Rect) {
	return b1.Bounds().Intersection(b2.Bounds())
}

func RandomCoordInArea(area Bounded) (c Coord) {
	r := area.Bounds()
	c.X = rand.Intn(r.W) + r.X
	c.Y = rand.Intn(r.H) + r.Y

	return
}
