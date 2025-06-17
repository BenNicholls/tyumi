package vec

import (
	"fmt"
	"iter"

	"github.com/bennicholls/tyumi/util"
)

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
	return fmt.Sprintf("{%s, %s}", r.Coord, r.Dims)
}

func (r Rect) Translated(coord Coord) Rect {
	return Rect{r.Coord.Add(coord), r.Dims}
}

// Returns a rect that is the rect r extended by some amount of cells in some direction.
func (r Rect) Extended(amount int, dir Direction) (extended Rect) {
	extended = r
	if x := dir.Coord().X; x < 0 {
		extended.X += x * amount
	}

	if y := dir.Coord().Y; y < 0 {
		extended.Y += y * amount
	}

	extended.W += util.Abs(dir.Coord().X) * amount
	extended.H += util.Abs(dir.Coord().Y) * amount

	return
}

// Returns a rect that is the rect r expanded in size by amount. Think of it like applying a layer of cells around the
// perimeter of the rect, repeated amount times.
// Example: expanding a 2x2 rect at the origin by 1 results in a 4x4 rect at position (-1, -1)
func (r Rect) Expanded(amount int) Rect {
	return Rect{
		Coord: r.Coord.Subtract(Coord{amount, amount}),
		Dims:  r.Dims.Grow(amount*2, amount*2),
	}
}

// Returns a rect that is the rect r contracted in size by amount. Think of it like removing the perimeter of the rect,
// repeated amount times.
// Example: contracting a 4x4 rect at the origin by 1 results in a 2x2 rect at position (1, 1)
func (r Rect) Contracted(amount int) Rect {
	return Rect{
		Coord: r.Coord.Add(Coord{amount, amount}),
		Dims:  r.Dims.Shrink(amount*2, amount*2),
	}
}

// Returns the coordinates of the 4 corners of the rect, starting in the top left and going clockwise.
func (r Rect) Corners() (corners [4]Coord) {
	corners[0] = r.Coord                             //TOPLEFT
	corners[1] = Coord{r.X + r.W - 1, r.Y}           //TOPRIGHT
	corners[2] = Coord{r.X + r.W - 1, r.Y + r.H - 1} //BOTTOMRIGHT
	corners[3] = Coord{r.X, r.Y + r.H - 1}           //BOTTOMLEFT

	return
}

func (r Rect) Sides() (sides [4]Line) {
	sides[0] = Line{r.Coord, Coord{r.X + r.W - 1, r.Y}}                             //TOP
	sides[1] = Line{Coord{r.X + r.W - 1, r.Y}, Coord{r.X + r.W - 1, r.Y + r.H - 1}} //RIGHT
	sides[2] = Line{Coord{r.X + r.W - 1, r.Y + r.H - 1}, Coord{r.X, r.Y + r.H - 1}} //BOTTOM
	sides[3] = Line{Coord{r.X, r.Y + r.H - 1}, r.Coord}                             //LEFT

	return
}

// IsInside calculates whether the rect is within the bounds of the provided rect.
// NOTE: for now, identical rects are NOT reported as being within eachother. May have to rethink this...
func (r Rect) IsInside(r2 Rect) bool {
	if r == r2 || r.Area() >= r2.Area() {
		return false
	}

	for _, corner := range r.Corners() {
		if !r2.Contains(corner) {
			return false
		}
	}

	return true
}

// Intersects returns true if r and r2 intersect in some way.
func (r Rect) Intersects(r2 Rect) bool {
	if r.Area() == 0 || r2.Area() == 0 { // zero-sized rects cannot intersect with anything
		return false
	}

	if r.X >= r2.X+r2.W || r2.X >= r.X+r.W || r.Y >= r2.Y+r2.H || r2.Y >= r.Y+r.H {
		return false
	}

	return true
}

// Intersection returns the intersection of r and r2.
func (r Rect) Intersection(r2 Rect) (i Rect) {
	if !r.Intersects(r2) {
		return
	}

	i.X, i.Y = max(r.X, r2.X), max(r.Y, r2.Y)
	i.W, i.H = min(r.X+r.W, r2.X+r2.W)-i.X, min(r.Y+r.H, r2.Y+r2.H)-i.Y
	return
}

// Returns the center of the rect. Since we're all integers 'round these parts this won't be exact unless both width
// and height are odd numbers, so be aware.
func (r Rect) Center() Coord {
	return Coord{r.X + r.W/2, r.Y + r.H/2}
}

// Contains calculates whether the provided coord is within the bounds of the rect.
func (r Rect) Contains(c Coord) bool {
	return c.X >= r.X && c.Y >= r.Y && c.X < r.X+r.W && c.Y < r.Y+r.H
}

func (r Rect) EachCoord() iter.Seq[Coord] {
	return func(yield func(Coord) bool) {
		for j := range r.H {
			for i := range r.W {
				if !yield(Coord{r.X + i, r.Y + j}) {
					return
				}
			}
		}
	}
}

// Returns an iterator producing a sequence of all coords in the perimeter of a bounded area.
func (r Rect) EachCoordInPerimeter() iter.Seq[Coord] {
	return func(yield func(Coord) bool) {
		if r.Area() == 0 { //0x0 box, aka a nothing
			return
		}

		if r.Area() == 1 { //1x1 box, aka 1 cell
			yield(r.Coord)
			return
		}

		if r.W == 1 || r.H == 1 { // 1D box, aka a line.
			corners := r.Corners()
			line := Line{corners[0], corners[2]}
			for coord := range line.EachCoord() {
				if !yield(coord) {
					return
				}
			}
			return
		}

		sides := r.Sides()
		for _, side := range sides {
			for coord := range side.EachCoord() {
				if coord == side.End { // don't yield the end coord, the next side will do it.
					break
				}

				if !yield(coord) {
					return
				}
			}
		}
	}
}

// CalcExtendedRect returns a rect that is the original rect r minimally extended so that it encompasses the provided
// coord.
func (r Rect) CalcExtendedRect(coord Coord) (extended Rect) {
	extended = r
	if coord.X < r.X {
		extended.X = coord.X
		extended.W += r.X - coord.X
	} else if coord.X >= r.X+r.W {
		extended.W = coord.X - r.X + 1
	}

	if coord.Y < r.Y {
		extended.Y = coord.Y
		extended.H += r.Y - coord.Y
	} else if coord.Y >= r.Y+r.H {
		extended.H = coord.Y - r.Y + 1
	}

	return
}

// CalcRectContainingCoord returns the smallest rect that contains the provided coords.
func CalcRectContainingCoords(c1, c2 Coord) Rect {
	if c1 == c2 {
		return Rect{c1, Dims{1, 1}}
	}

	minX, maxX := min(c1.X, c2.X), max(c1.X, c2.X)
	minY, maxY := min(c1.Y, c1.Y), max(c1.Y, c1.Y)

	return Rect{Coord{minX, minY}, Dims{maxX - minX + 1, maxY - minY + 1}}
}
