package vec

import (
	"slices"
	"testing"
)

func TestRectBounds(t *testing.T) {
	r := Rect{Coord{1, 2}, Dims{3, 4}}
	if got := r.Bounds(); got != r {
		t.Errorf("Bounds() = %v, want %v", got, r)
	}
}

func TestRectString(t *testing.T) {
	r := Rect{Coord{1, 2}, Dims{3, 4}}
	want := "{(X: 1, Y: 2), (W: 3, H: 4)}"
	if got := r.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestRectTranslated(t *testing.T) {
	r := Rect{Coord{1, 1}, Dims{3, 3}}
	got := r.Translated(Coord{2, 3})
	want := Rect{Coord{3, 4}, Dims{3, 3}}
	if got != want {
		t.Errorf("Translated() = %v, want %v", got, want)
	}
}

func TestRectExtended(t *testing.T) {
	r := Rect{Coord{2, 2}, Dims{3, 3}}

	tests := []struct {
		name   string
		dir    Direction
		amount int
		want   Rect
	}{
		{"right", DIR_RIGHT, 2, Rect{Coord{2, 2}, Dims{5, 3}}},
		{"left", DIR_LEFT, 2, Rect{Coord{0, 2}, Dims{5, 3}}},
		{"down", DIR_DOWN, 1, Rect{Coord{2, 2}, Dims{3, 4}}},
		{"up", DIR_UP, 1, Rect{Coord{2, 1}, Dims{3, 4}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Extended(tt.amount, tt.dir); got != tt.want {
				t.Errorf("Extended(%d, %v) = %v, want %v", tt.amount, tt.dir, got, tt.want)
			}
		})
	}
}

func TestRectExpanded(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{2, 2}}
	got := r.Expanded(1)
	want := Rect{Coord{-1, -1}, Dims{4, 4}}
	if got != want {
		t.Errorf("Expanded(1) = %v, want %v", got, want)
	}
}

func TestRectContracted(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{4, 4}}
	got := r.Contracted(1)
	want := Rect{Coord{1, 1}, Dims{2, 2}}
	if got != want {
		t.Errorf("Contracted(1) = %v, want %v", got, want)
	}
}

func TestRectCorners(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{3, 3}}
	got := r.Corners()
	want := [4]Coord{{0, 0}, {2, 0}, {2, 2}, {0, 2}}
	if got != want {
		t.Errorf("Corners() = %v, want %v", got, want)
	}
}

func TestRectSides(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{3, 3}}
	got := r.Sides()
	want := [4]Line{
		{Coord{0, 0}, Coord{2, 0}},
		{Coord{2, 0}, Coord{2, 2}},
		{Coord{2, 2}, Coord{0, 2}},
		{Coord{0, 2}, Coord{0, 0}},
	}
	if got != want {
		t.Errorf("Sides() = %v, want %v", got, want)
	}
}

func TestRectIsInside(t *testing.T) {
	outer := Rect{Coord{0, 0}, Dims{10, 10}}
	inner := Rect{Coord{2, 2}, Dims{3, 3}}

	if !inner.IsInside(outer) {
		t.Error("inner rect should be inside outer rect")
	}
	if outer.IsInside(inner) {
		t.Error("outer rect should not be inside inner rect")
	}
	if outer.IsInside(outer) {
		t.Error("identical rects should not be reported as inside each other")
	}
}

func TestRectIntersects(t *testing.T) {
	tests := []struct {
		name   string
		r1, r2 Rect
		want   bool
	}{
		{"overlapping", Rect{Coord{0, 0}, Dims{5, 5}}, Rect{Coord{3, 3}, Dims{5, 5}}, true},
		{"non-overlapping", Rect{Coord{0, 0}, Dims{2, 2}}, Rect{Coord{10, 10}, Dims{2, 2}}, false},
		{"touching edges", Rect{Coord{0, 0}, Dims{2, 2}}, Rect{Coord{2, 0}, Dims{2, 2}}, false},
		{"zero area", Rect{Coord{0, 0}, Dims{0, 0}}, Rect{Coord{0, 0}, Dims{5, 5}}, false},
		{"identical", Rect{Coord{0, 0}, Dims{5, 5}}, Rect{Coord{0, 0}, Dims{5, 5}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r1.Intersects(tt.r2); got != tt.want {
				t.Errorf("%v.Intersects(%v) = %v, want %v", tt.r1, tt.r2, got, tt.want)
			}
		})
	}
}

func TestRectIntersection(t *testing.T) {
	r1 := Rect{Coord{0, 0}, Dims{5, 5}}
	r2 := Rect{Coord{3, 3}, Dims{5, 5}}

	got := r1.Intersection(r2)
	want := Rect{Coord{3, 3}, Dims{2, 2}}
	if got != want {
		t.Errorf("Intersection() = %v, want %v", got, want)
	}

	noOverlap := Rect{Coord{100, 100}, Dims{2, 2}}
	got2 := r1.Intersection(noOverlap)
	want2 := Rect{}
	if got2 != want2 {
		t.Errorf("Intersection() of non-overlapping rects = %v, want zero value %v", got2, want2)
	}
}

func TestRectCenter(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{4, 6}}
	got := r.Center()
	want := Coord{2, 3}
	if got != want {
		t.Errorf("Center() = %v, want %v", got, want)
	}
}

func TestRectContains(t *testing.T) {
	r := Rect{Coord{1, 1}, Dims{3, 3}}

	tests := []struct {
		c    Coord
		want bool
	}{
		{Coord{1, 1}, true},
		{Coord{3, 3}, true},
		{Coord{4, 4}, false}, // one past the far edge
		{Coord{0, 1}, false},
	}

	for _, tt := range tests {
		if got := r.Contains(tt.c); got != tt.want {
			t.Errorf("Contains(%v) = %v, want %v", tt.c, got, tt.want)
		}
	}
}

func TestRectEachCoord(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{2, 2}}

	var got []Coord
	for c := range r.EachCoord() {
		got = append(got, c)
	}

	want := []Coord{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoord() = %v, want %v", got, want)
	}
}

func TestRectEachCoordEarlyStop(t *testing.T) {
	r := Rect{Coord{0, 0}, Dims{5, 5}}

	count := 0
	for range r.EachCoord() {
		count++
		if count == 3 {
			break
		}
	}

	if count != 3 {
		t.Errorf("expected iteration to stop after 3 coords, got %d", count)
	}
}

func TestRectEachCoordInPerimeter(t *testing.T) {
	tests := []struct {
		name string
		r    Rect
		want []Coord
	}{
		{"zero area", Rect{Coord{0, 0}, Dims{0, 0}}, nil},
		{"single cell", Rect{Coord{2, 2}, Dims{1, 1}}, []Coord{{2, 2}}},
		{"1D horizontal line", Rect{Coord{0, 0}, Dims{3, 1}}, []Coord{{0, 0}, {1, 0}, {2, 0}}},
		{"1D vertical line", Rect{Coord{0, 0}, Dims{1, 3}}, []Coord{{0, 0}, {0, 1}, {0, 2}}},
		{
			"3x3 square perimeter",
			Rect{Coord{0, 0}, Dims{3, 3}},
			[]Coord{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}, {1, 2}, {0, 2}, {0, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []Coord
			for c := range tt.r.EachCoordInPerimeter() {
				got = append(got, c)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("EachCoordInPerimeter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRectCalcExtendedRect(t *testing.T) {
	tests := []struct {
		name  string
		r     Rect
		coord Coord
		want  Rect
	}{
		{"zero area rect", Rect{}, Coord{5, 5}, Rect{Coord{5, 5}, Dims{1, 1}}},
		{"extend left", Rect{Coord{3, 3}, Dims{2, 2}}, Coord{0, 3}, Rect{Coord{0, 3}, Dims{5, 2}}},
		{"extend right", Rect{Coord{3, 3}, Dims{2, 2}}, Coord{10, 3}, Rect{Coord{3, 3}, Dims{8, 2}}},
		{"extend up", Rect{Coord{3, 3}, Dims{2, 2}}, Coord{3, 0}, Rect{Coord{3, 0}, Dims{2, 5}}},
		{"extend down", Rect{Coord{3, 3}, Dims{2, 2}}, Coord{3, 10}, Rect{Coord{3, 3}, Dims{2, 8}}},
		{"coord already inside", Rect{Coord{0, 0}, Dims{5, 5}}, Coord{2, 2}, Rect{Coord{0, 0}, Dims{5, 5}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.CalcExtendedRect(tt.coord); got != tt.want {
				t.Errorf("CalcExtendedRect(%v) = %v, want %v", tt.coord, got, tt.want)
			}
		})
	}
}

func TestCalcRectContainingCoords(t *testing.T) {
	tests := []struct {
		name   string
		c1, c2 Coord
		want   Rect
	}{
		{"same coord", Coord{3, 3}, Coord{3, 3}, Rect{Coord{3, 3}, Dims{1, 1}}},
		{"horizontal only, same X", Coord{1, 5}, Coord{9, 5}, Rect{Coord{1, 5}, Dims{9, 1}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcRectContainingCoords(tt.c1, tt.c2); got != tt.want {
				t.Errorf("CalcRectContainingCoords(%v, %v) = %v, want %v", tt.c1, tt.c2, got, tt.want)
			}
		})
	}
}

// BUG: CalcRectContainingCoords computes minY/maxY from c1.Y compared against itself (min(c1.Y, c1.Y),
// max(c1.Y, c1.Y)) instead of comparing c1.Y against c2.Y, so the Y dimension of the result is always based
// entirely on c1's Y coordinate and never actually accounts for c2's Y. This test documents that current
// (likely unintended) behaviour rather than the probably-intended one, since fixing it would change public API
// behaviour.
func TestCalcRectContainingCoordsYAxisBug(t *testing.T) {
	got := CalcRectContainingCoords(Coord{1, 1}, Coord{9, 9})
	want := Rect{Coord{1, 1}, Dims{9, 1}} // a correct implementation would produce Dims{9, 9}
	if got != want {
		t.Errorf("CalcRectContainingCoords({1,1}, {9,9}) = %v, want %v (see BUG note above)", got, want)
	}
}
