package vec

import (
	"iter"

	"github.com/bennicholls/tyumi/log"
)

type Circle struct {
	Radius int
	Center Coord
}

func (c Circle) EachCoordInPerimeter() iter.Seq[Coord] {
	return func(yield func(Coord) bool) {
		for arcPos := range arcIterator(c.Radius) {
			for i := range 8 {
				if !yield(transformArcToOctant(c.Center, arcPos, i)) {
					return
				}
			}
		}
	}
}

// Computes a circle, calling fn for each point on the perimeter the circle.
// fn can be a drawing function or whatever.
func CircleFunc(center Coord, radius int, fn func(pos Coord)) {
	circle := Circle{radius, center}
	for cursor := range circle.EachCoordInPerimeter() {
		fn(cursor)
	}
}

func arcIterator(radius int) iter.Seq[Coord] {
	x, y := 0, radius
	f := 1 - radius
	ddf_x, ddf_y := 1, -2*radius

	return func(yield func(Coord) bool) {
		for x <= y {
			if !yield(Coord{x, y}) {
				return
			}

			if f >= 0 {
				y -= 1
				ddf_y += 2
				f += ddf_y
			}

			x += 1
			ddf_x += 2
			f += ddf_x
		}
	}
}

// octant is 0-7, arcPos is the position returned by the arc generator
func transformArcToOctant(center, arcPos Coord, octant int) Coord {
	switch octant {
	case 0:
		return Coord{center.X + arcPos.X, center.Y + arcPos.Y}
	case 1:
		return Coord{center.X + arcPos.Y, center.Y + arcPos.X}
	case 2:
		return Coord{center.X - arcPos.Y, center.Y + arcPos.X}
	case 3:
		return Coord{center.X - arcPos.X, center.Y + arcPos.Y}
	case 4:
		return Coord{center.X - arcPos.X, center.Y - arcPos.Y}
	case 5:
		return Coord{center.X - arcPos.Y, center.Y - arcPos.X}
	case 6:
		return Coord{center.X + arcPos.Y, center.Y - arcPos.X}
	case 7:
		return Coord{center.X + arcPos.X, center.Y - arcPos.Y}
	default:
		log.Error("Octant is bigger than 7. Don't you know what octant *means*????")
		return ZERO_COORD
	}
}
