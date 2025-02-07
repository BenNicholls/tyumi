package vec

import (
	"iter"

	"github.com/bennicholls/tyumi/util"
)

// Line represents a line between 2 points on a square grid.
type Line struct {
	Start Coord
	End   Coord
}

// Length computes the precise cartesian length of the line. Consider using LengthSq() if you don't actually need this
// absolute value and are just doing comparisons.
func (l Line) Length() float64 {
	return Vec2i(l.Start.Subtract(l.End)).Mag()
}

// LengthSq returns the length of the line squared. Good for comparison purposes in loops where a nasty sqrt() makes
// you go :(
func (l Line) LengthSq() int {
	return (l.Start.X-l.End.X)*2 + (l.Start.Y-l.End.Y)*2
}

// EachCoord returns an iterator that produces Coords representing the line from start to end inclusive.
func (l Line) EachCoord() iter.Seq[Coord] {
	//detect if we're drawing a straight line or not. if we are, we can skip the more expensive bresenham algorithm in
	//lieu of a simple loop stepping in the right direction until done
	var step_dir Direction
	switch {
	case l.dy() == 0: //horizontal line
		step_dir = DIR_RIGHT
		if l.Start.X > l.End.X {
			step_dir = DIR_LEFT
		}
	case l.dx() == 0: //vertical line
		step_dir = DIR_DOWN
		if l.Start.Y > l.End.Y {
			step_dir = DIR_UP
		}
	case l.dx() == l.dy(): //UPLEFT to DOWNRIGHT
		step_dir = Direction{-1, -1}
		if l.Start.X < l.End.X {
			step_dir = Direction{1, 1}
		}
	case l.dx() == -l.dy(): //DOWNLEFT to UPRIGHT
		step_dir = Direction{-1, 1}
		if l.Start.X < l.End.X {
			step_dir = Direction{1, -1}
		}
	}

	if step_dir != DIR_NONE { // return the simple direction-stepping iterator
		return func(yield func(Coord) bool) {
			for c := l.Start; c != l.End.Step(step_dir); c = c.Step(step_dir) {
				if !yield(c) {
					return
				}
			}
		}
	}

	// setup and return an iterator implementing bresenham's algorithm
	dx := util.Abs(l.dx())
	dy := -util.Abs(l.dy())
	e := dx + dy

	sx := -1
	if l.Start.X < l.End.X {
		sx = 1
	}

	sy := -1
	if l.Start.Y < l.End.Y {
		sy = 1
	}

	return func(yield func(Coord) bool) {
		for {
			if !yield(l.Start) {
				return
			}
			e2 := 2 * e
			if e2 >= dy {
				if l.Start.X == l.End.X {
					return
				}
				e = e + dy
				l.Start.X = l.Start.X + sx
			}
			if e2 <= dx {
				if l.Start.Y == l.End.Y {
					return
				}
				e = e + dx
				l.Start.Y = l.Start.Y + sy
			}
		}
	}
}

func (l Line) dx() int {
	return l.End.X - l.Start.X
}

func (l Line) dy() int {
	return l.End.Y - l.Start.Y
}
