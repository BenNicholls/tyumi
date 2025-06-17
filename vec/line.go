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

// Length computes the precise cartesian length of the line. Consider using LengthSq() if you don't actually need the
// absolute value and are just doing comparisons.
func (l Line) Length() float64 {
	return l.Start.DistanceTo(l.End)
}

// LengthSq returns the length of the line squared. Good for comparison purposes in loops where a nasty sqrt() makes
// you go :(
func (l Line) LengthSq() int {
	return l.Start.DistanceSqTo(l.End)
}

// Contracted returns the line l with the first N and last N points removed (where N = amount).
func (l Line) Contracted(amount int) (contracted Line) {
	i := 0
	for cursor := range l.EachCoord() {
		if i == amount {
			contracted.Start = cursor
			break
		}
		i += 1
	}

	i = 0
	backwards := Line{l.End, l.Start}
	for cursor := range backwards.EachCoord() {
		if i == amount {
			contracted.End = cursor
			break
		}
		i += 1
	}

	return
}

// EachCoord returns an iterator that produces Coords representing the line from start to end inclusive.
func (l Line) EachCoord() iter.Seq[Coord] {
	//detect if we're drawing a straight line or not. if we are, we can skip the more expensive bresenham algorithm in
	//lieu of a simple loop stepping in the right direction until done
	var step_dir Direction = DIR_NONE
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
		step_dir = DIR_UPLEFT
		if l.Start.X < l.End.X {
			step_dir = DIR_DOWNRIGHT
		}
	case l.dx() == -l.dy(): //DOWNLEFT to UPRIGHT
		step_dir = DIR_DOWNLEFT
		if l.Start.X < l.End.X {
			step_dir = DIR_UPRIGHT
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
	dx, dy := util.Abs(l.dx()), -util.Abs(l.dy())
	e := dx + dy
	e2 := e * 2
	sx, sy := -1, -1
	if l.Start.X < l.End.X {
		sx = 1
	}
	if l.Start.Y < l.End.Y {
		sy = 1
	}

	return func(yield func(Coord) bool) {
		for {
			if !yield(l.Start) {
				return
			}
			if e2 >= dy {
				if l.Start.X == l.End.X {
					return
				}
				e = e + dy
				l.Start.X += sx
			}
			if e2 <= dx {
				if l.Start.Y == l.End.Y {
					return
				}
				e = e + dx
				l.Start.Y += sy
			}
			e2 = e * 2
		}
	}
}

func (l Line) dx() int {
	return l.End.X - l.Start.X
}

func (l Line) dy() int {
	return l.End.Y - l.Start.Y
}
