package vec

import (
	"slices"
	"testing"
)

func TestLineLength(t *testing.T) {
	l := Line{Coord{0, 0}, Coord{3, 4}}
	if got := l.Length(); !almostEqual(got, 5) {
		t.Errorf("Length() = %v, want 5", got)
	}
}

func TestLineLengthSq(t *testing.T) {
	l := Line{Coord{0, 0}, Coord{3, 4}}
	if got := l.LengthSq(); got != 25 {
		t.Errorf("LengthSq() = %d, want 25", got)
	}
}

func TestLineEachCoordHorizontal(t *testing.T) {
	l := Line{Coord{0, 0}, Coord{3, 0}}
	var got []Coord
	for c := range l.EachCoord() {
		got = append(got, c)
	}
	want := []Coord{{0, 0}, {1, 0}, {2, 0}, {3, 0}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoord() = %v, want %v", got, want)
	}
}

func TestLineEachCoordVertical(t *testing.T) {
	l := Line{Coord{0, 3}, Coord{0, 0}}
	var got []Coord
	for c := range l.EachCoord() {
		got = append(got, c)
	}
	want := []Coord{{0, 3}, {0, 2}, {0, 1}, {0, 0}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoord() = %v, want %v", got, want)
	}
}

func TestLineEachCoordDiagonal(t *testing.T) {
	l := Line{Coord{0, 0}, Coord{3, 3}}
	var got []Coord
	for c := range l.EachCoord() {
		got = append(got, c)
	}
	want := []Coord{{0, 0}, {1, 1}, {2, 2}, {3, 3}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoord() = %v, want %v", got, want)
	}
}

func TestLineEachCoordAntiDiagonal(t *testing.T) {
	l := Line{Coord{0, 3}, Coord{3, 0}}
	var got []Coord
	for c := range l.EachCoord() {
		got = append(got, c)
	}
	want := []Coord{{0, 3}, {1, 2}, {2, 1}, {3, 0}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoord() = %v, want %v", got, want)
	}
}

func TestLineEachCoordSinglePoint(t *testing.T) {
	l := Line{Coord{5, 5}, Coord{5, 5}}
	var got []Coord
	for c := range l.EachCoord() {
		got = append(got, c)
	}
	want := []Coord{{5, 5}}
	if !slices.Equal(got, want) {
		t.Errorf("EachCoord() = %v, want %v", got, want)
	}
}

func TestLineEachCoordBresenham(t *testing.T) {
	// non-45-degree, non-axis-aligned line, exercises the bresenham branch
	l := Line{Coord{0, 0}, Coord{5, 2}}
	var got []Coord
	for c := range l.EachCoord() {
		got = append(got, c)
	}

	if len(got) == 0 {
		t.Fatal("EachCoord() produced no points")
	}
	if got[0] != l.Start {
		t.Errorf("first point = %v, want start %v", got[0], l.Start)
	}
	if got[len(got)-1] != l.End {
		t.Errorf("last point = %v, want end %v", got[len(got)-1], l.End)
	}

	// bresenham line should be continuous: each step moves by at most 1 in each axis
	for i := 1; i < len(got); i++ {
		dx := Abs(got[i].X - got[i-1].X)
		dy := Abs(got[i].Y - got[i-1].Y)
		if dx > 1 || dy > 1 {
			t.Errorf("non-contiguous step from %v to %v", got[i-1], got[i])
		}
	}
}

func TestLineEachCoordEarlyStop(t *testing.T) {
	l := Line{Coord{0, 0}, Coord{10, 0}}
	count := 0
	for range l.EachCoord() {
		count++
		if count == 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("expected iteration to stop after 2 coords, got %d", count)
	}
}

func TestLineContracted(t *testing.T) {
	l := Line{Coord{0, 0}, Coord{5, 0}}
	got := l.Contracted(1)
	want := Line{Coord{1, 0}, Coord{4, 0}}
	if got != want {
		t.Errorf("Contracted(1) = %v, want %v", got, want)
	}
}

func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
