package vec

import (
	"slices"
	"testing"
)

type testBounded struct {
	bounds Rect
}

func (tb testBounded) Bounds() Rect {
	return tb.bounds
}

func TestEachCoordInArea(t *testing.T) {
	b := testBounded{Rect{Coord{0, 0}, Dims{2, 2}}}

	var got []Coord
	for c := range EachCoordInArea(b) {
		got = append(got, c)
	}

	want := []Coord{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoordInArea() = %v, want %v", got, want)
	}
}

func TestEachCoordInIntersectionNoArgs(t *testing.T) {
	count := 0
	for range EachCoordInIntersection() {
		count++
	}
	if count != 0 {
		t.Errorf("EachCoordInIntersection() with no args produced %d coords, want 0", count)
	}
}

func TestEachCoordInIntersectionSingleArg(t *testing.T) {
	b := testBounded{Rect{Coord{0, 0}, Dims{2, 2}}}

	var got []Coord
	for c := range EachCoordInIntersection(b) {
		got = append(got, c)
	}

	want := []Coord{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoordInIntersection(single) = %v, want %v", got, want)
	}
}

func TestEachCoordInIntersectionMultiple(t *testing.T) {
	b1 := testBounded{Rect{Coord{0, 0}, Dims{5, 5}}}
	b2 := testBounded{Rect{Coord{2, 2}, Dims{5, 5}}}

	var got []Coord
	for c := range EachCoordInIntersection(b1, b2) {
		got = append(got, c)
	}

	want := []Coord{{2, 2}, {3, 2}, {4, 2}, {2, 3}, {3, 3}, {4, 3}, {2, 4}, {3, 4}, {4, 4}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoordInIntersection(multiple) = %v, want %v", got, want)
	}
}

func TestEachCoordInIntersectionNoOverlap(t *testing.T) {
	b1 := testBounded{Rect{Coord{0, 0}, Dims{2, 2}}}
	b2 := testBounded{Rect{Coord{100, 100}, Dims{2, 2}}}

	count := 0
	for range EachCoordInIntersection(b1, b2) {
		count++
	}
	if count != 0 {
		t.Errorf("EachCoordInIntersection with no overlap produced %d coords, want 0", count)
	}
}

func TestEachCoordInPerimeter(t *testing.T) {
	b := testBounded{Rect{Coord{0, 0}, Dims{1, 1}}}

	var got []Coord
	for c := range EachCoordInPerimeter(b) {
		got = append(got, c)
	}

	want := []Coord{{0, 0}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoordInPerimeter() = %v, want %v", got, want)
	}
}

func TestBoundedIntersects(t *testing.T) {
	b1 := testBounded{Rect{Coord{0, 0}, Dims{5, 5}}}
	b2 := testBounded{Rect{Coord{3, 3}, Dims{5, 5}}}
	b3 := testBounded{Rect{Coord{100, 100}, Dims{5, 5}}}

	if !Intersects(b1, b2) {
		t.Error("expected b1 and b2 to intersect")
	}
	if Intersects(b1, b3) {
		t.Error("expected b1 and b3 to not intersect")
	}
}

func TestFindIntersectionRect(t *testing.T) {
	b1 := testBounded{Rect{Coord{0, 0}, Dims{5, 5}}}
	b2 := testBounded{Rect{Coord{3, 3}, Dims{5, 5}}}

	got := FindIntersectionRect(b1, b2)
	want := Rect{Coord{3, 3}, Dims{2, 2}}
	if got != want {
		t.Errorf("FindIntersectionRect() = %v, want %v", got, want)
	}

	b3 := testBounded{Rect{Coord{100, 100}, Dims{2, 2}}}
	got2 := FindIntersectionRect(b1, b3)
	want2 := Rect{}
	if got2 != want2 {
		t.Errorf("FindIntersectionRect() of non-overlapping = %v, want %v", got2, want2)
	}
}

func TestRandomCoordInArea(t *testing.T) {
	b := testBounded{Rect{Coord{5, 5}, Dims{10, 10}}}

	for range 50 {
		c := RandomCoordInArea(b)
		if !c.IsInside(b) {
			t.Errorf("RandomCoordInArea() = %v, not inside %v", c, b.bounds)
		}
	}
}
