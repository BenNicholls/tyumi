package vec

import "testing"

func TestDimsArea(t *testing.T) {
	tests := []struct {
		d    Dims
		want int
	}{
		{Dims{3, 4}, 12},
		{Dims{0, 5}, 0},
		{Dims{5, 5}, 25},
	}

	for _, tt := range tests {
		if got := tt.d.Area(); got != tt.want {
			t.Errorf("%v.Area() = %d, want %d", tt.d, got, tt.want)
		}
	}
}

func TestDimsGrow(t *testing.T) {
	tests := []struct {
		d      Dims
		dw, dh int
		want   Dims
	}{
		{Dims{3, 4}, 2, 2, Dims{5, 6}},
		{Dims{3, 4}, -10, -10, Dims{0, 0}}, // clamped at 0
	}

	for _, tt := range tests {
		if got := tt.d.Grow(tt.dw, tt.dh); got != tt.want {
			t.Errorf("%v.Grow(%d, %d) = %v, want %v", tt.d, tt.dw, tt.dh, got, tt.want)
		}
	}
}

func TestDimsShrink(t *testing.T) {
	tests := []struct {
		d      Dims
		dw, dh int
		want   Dims
	}{
		{Dims{5, 6}, 2, 2, Dims{3, 4}},
		{Dims{3, 4}, 10, 10, Dims{0, 0}}, // clamped at 0
	}

	for _, tt := range tests {
		if got := tt.d.Shrink(tt.dw, tt.dh); got != tt.want {
			t.Errorf("%v.Shrink(%d, %d) = %v, want %v", tt.d, tt.dw, tt.dh, got, tt.want)
		}
	}
}

func TestDimsBounds(t *testing.T) {
	d := Dims{4, 5}
	got := d.Bounds()
	want := Rect{Coord{0, 0}, Dims{4, 5}}
	if got != want {
		t.Errorf("Bounds() = %v, want %v", got, want)
	}
}

func TestDimsString(t *testing.T) {
	got := Dims{3, 4}.String()
	want := "(W: 3, H: 4)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
