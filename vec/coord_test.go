package vec

import "testing"

func TestCoordMove(t *testing.T) {
	c := Coord{1, 1}
	c.Move(2, -1)
	if c != (Coord{3, 0}) {
		t.Errorf("Move() = %v, want {3 0}", c)
	}
}

func TestCoordMoveTo(t *testing.T) {
	c := Coord{1, 1}
	c.MoveTo(5, 6)
	if c != (Coord{5, 6}) {
		t.Errorf("MoveTo() = %v, want {5 6}", c)
	}
}

func TestCoordAddSubtract(t *testing.T) {
	c1 := Coord{1, 2}
	c2 := Coord{3, 4}

	if got := c1.Add(c2); got != (Coord{4, 6}) {
		t.Errorf("Add = %v, want {4 6}", got)
	}
	if got := c1.Subtract(c2); got != (Coord{-2, -2}) {
		t.Errorf("Subtract = %v, want {-2 -2}", got)
	}
}

func TestCoordStep(t *testing.T) {
	c := Coord{5, 5}
	if got := c.Step(DIR_UP); got != (Coord{5, 4}) {
		t.Errorf("Step(UP) = %v, want {5 4}", got)
	}
	if got := c.Step(DIR_NONE); got != c {
		t.Errorf("Step(NONE) = %v, want %v (unchanged)", got, c)
	}
}

func TestCoordStepN(t *testing.T) {
	c := Coord{5, 5}
	if got := c.StepN(DIR_RIGHT, 3); got != (Coord{8, 5}) {
		t.Errorf("StepN(RIGHT, 3) = %v, want {8 5}", got)
	}
}

func TestCoordScale(t *testing.T) {
	got := Coord{2, 3}.Scale(4)
	want := Coord{8, 12}
	if got != want {
		t.Errorf("Scale(4) = %v, want %v", got, want)
	}
}

func TestCoordToIndex(t *testing.T) {
	tests := []struct {
		c      Coord
		stride int
		want   int
	}{
		{Coord{0, 0}, 10, 0},
		{Coord{3, 0}, 10, 3},
		{Coord{3, 2}, 10, 23},
	}

	for _, tt := range tests {
		if got := tt.c.ToIndex(tt.stride); got != tt.want {
			t.Errorf("%v.ToIndex(%d) = %d, want %d", tt.c, tt.stride, got, tt.want)
		}
	}
}

func TestIndexToCoord(t *testing.T) {
	tests := []struct {
		index, stride int
		want          Coord
	}{
		{0, 10, Coord{0, 0}},
		{3, 10, Coord{3, 0}},
		{23, 10, Coord{3, 2}},
	}

	for _, tt := range tests {
		if got := IndexToCoord(tt.index, tt.stride); got != tt.want {
			t.Errorf("IndexToCoord(%d, %d) = %v, want %v", tt.index, tt.stride, got, tt.want)
		}
	}
}

func TestCoordIsInside(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{10, 10}}

	if !(Coord{5, 5}).IsInside(r) {
		t.Error("{5,5} should be inside {0,0,10,10}")
	}
	if (Coord{10, 10}).IsInside(r) {
		t.Error("{10,10} should not be inside {0,0,10,10} (exclusive of far edge)")
	}
	if (Coord{-1, 5}).IsInside(r) {
		t.Error("{-1,5} should not be inside {0,0,10,10}")
	}
}

func TestCoordIsInPerimeter(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{5, 5}}

	if !(Coord{0, 2}).IsInPerimeter(r) {
		t.Error("{0,2} should be in perimeter (left edge)")
	}
	if !(Coord{4, 4}).IsInPerimeter(r) {
		t.Error("{4,4} should be in perimeter (bottom right corner)")
	}
	if (Coord{2, 2}).IsInPerimeter(r) {
		t.Error("{2,2} should not be in perimeter (center)")
	}
}

func TestCoordString(t *testing.T) {
	got := Coord{3, 4}.String()
	want := "(X: 3, Y: 4)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestCoordDistanceTo(t *testing.T) {
	c1 := Coord{0, 0}
	c2 := Coord{3, 4}

	if got := c1.DistanceTo(c2); !almostEqual(got, 5) {
		t.Errorf("DistanceTo = %v, want 5", got)
	}
}

func TestCoordDistanceSqTo(t *testing.T) {
	c1 := Coord{0, 0}
	c2 := Coord{3, 4}

	if got := c1.DistanceSqTo(c2); got != 25 {
		t.Errorf("DistanceSqTo = %d, want 25", got)
	}
}

func TestCoordManhattanDistanceTo(t *testing.T) {
	c1 := Coord{0, 0}
	c2 := Coord{3, -4}

	if got := c1.ManhattanDistanceTo(c2); got != 7 {
		t.Errorf("ManhattanDistanceTo = %d, want 7", got)
	}
}

func TestCoordLerp(t *testing.T) {
	c1 := Coord{0, 0}
	c2 := Coord{10, 10}

	if got := c1.Lerp(c2, 0, 10); got != c1 {
		t.Errorf("Lerp at val=0 = %v, want %v", got, c1)
	}
	if got := c1.Lerp(c2, 10, 10); got != c2 {
		t.Errorf("Lerp at val=steps = %v, want %v", got, c2)
	}
	if got := c1.Lerp(c2, 5, 10); got != (Coord{5, 5}) {
		t.Errorf("Lerp at val=5/10 = %v, want {5 5}", got)
	}
}

func TestCoordDirectionTo(t *testing.T) {
	c := Coord{5, 5}

	tests := []struct {
		to   Coord
		want Direction
	}{
		{Coord{5, 5}, DIR_NONE},
		{Coord{5, 0}, DIR_UP},
		{Coord{5, 10}, DIR_DOWN},
		{Coord{0, 5}, DIR_LEFT},
		{Coord{10, 5}, DIR_RIGHT},
		// non-45-degree diagonals, unambiguously closer to a diagonal direction
		{Coord{11, 0}, DIR_UPRIGHT},
		{Coord{0, -1}, DIR_UPLEFT},
		{Coord{11, 10}, DIR_DOWNRIGHT},
		{Coord{0, 11}, DIR_DOWNLEFT},
		// NOTE: at an exact 45-degree angle (|dx| == |dy|), DirectionTo is biased towards the horizontal
		// direction rather than the diagonal one, per the strict "<" comparisons in the implementation.
		{Coord{10, 0}, DIR_RIGHT},
		{Coord{0, 0}, DIR_LEFT},
		{Coord{10, 10}, DIR_RIGHT},
		{Coord{0, 10}, DIR_LEFT},
	}

	for _, tt := range tests {
		if got := c.DirectionTo(tt.to); got != tt.want {
			t.Errorf("DirectionTo(%v) = %v, want %v", tt.to, got, tt.want)
		}
	}
}

func TestDirectionCoord(t *testing.T) {
	tests := []struct {
		d    Direction
		want Coord
	}{
		{DIR_UP, Coord{0, -1}},
		{DIR_DOWN, Coord{0, 1}},
		{DIR_LEFT, Coord{-1, 0}},
		{DIR_RIGHT, Coord{1, 0}},
		{DIR_NONE, Coord{0, 0}},
	}

	for _, tt := range tests {
		if got := tt.d.Coord(); got != tt.want {
			t.Errorf("%v.Coord() = %v, want %v", tt.d, got, tt.want)
		}
	}
}

func TestDirectionInverted(t *testing.T) {
	tests := []struct {
		d    Direction
		want Direction
	}{
		{DIR_UP, DIR_DOWN},
		{DIR_DOWN, DIR_UP},
		{DIR_LEFT, DIR_RIGHT},
		{DIR_RIGHT, DIR_LEFT},
		{DIR_UPLEFT, DIR_DOWNRIGHT},
	}

	for _, tt := range tests {
		if got := tt.d.Inverted(); got != tt.want {
			t.Errorf("%v.Inverted() = %v, want %v", tt.d, got, tt.want)
		}
	}
}

func TestDirectionRotateCW(t *testing.T) {
	if got := DIR_UP.RotateCW(); got != DIR_UPRIGHT {
		t.Errorf("RotateCW() = %v, want DIR_UPRIGHT", got)
	}
	if got := DIR_UPLEFT.RotateCW(); got != DIR_UP {
		t.Errorf("RotateCW() wraps = %v, want DIR_UP", got)
	}
	if got := DIR_NONE.RotateCW(); got != DIR_NONE {
		t.Errorf("RotateCW() on DIR_NONE = %v, want DIR_NONE", got)
	}
}

func TestDirectionRotateCCW(t *testing.T) {
	if got := DIR_UP.RotateCCW(); got != DIR_UPLEFT {
		t.Errorf("RotateCCW() = %v, want DIR_UPLEFT", got)
	}
	if got := DIR_NONE.RotateCCW(); got != DIR_NONE {
		t.Errorf("RotateCCW() on DIR_NONE = %v, want DIR_NONE", got)
	}
}

func TestDirectionRotateCW90(t *testing.T) {
	if got := DIR_UP.RotateCW90(); got != DIR_RIGHT {
		t.Errorf("RotateCW90() = %v, want DIR_RIGHT", got)
	}
	if got := DIR_NONE.RotateCW90(); got != DIR_NONE {
		t.Errorf("RotateCW90() on DIR_NONE = %v, want DIR_NONE", got)
	}
}

func TestDirectionRotateCCW90(t *testing.T) {
	if got := DIR_UP.RotateCCW90(); got != DIR_LEFT {
		t.Errorf("RotateCCW90() = %v, want DIR_LEFT", got)
	}
	if got := DIR_NONE.RotateCCW90(); got != DIR_NONE {
		t.Errorf("RotateCCW90() on DIR_NONE = %v, want DIR_NONE", got)
	}
}

func TestRandomDirection(t *testing.T) {
	for range 20 {
		d := RandomDirection()
		found := false
		for _, valid := range Directions {
			if d == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("RandomDirection() = %v, not in Directions", d)
		}
	}
}

func TestRandomCardinalDirection(t *testing.T) {
	for range 20 {
		d := RandomCardinalDirection()
		found := false
		for _, valid := range CardinalDirections {
			if d == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("RandomCardinalDirection() = %v, not in CardinalDirections", d)
		}
	}
}
