package vec

import "fmt"

// Dims represents a set of dimensions in 2D.
type Dims struct {
	W, H int
}

func (d Dims) Area() int {
	return d.W * d.H
}

func (d Dims) String() string {
	return fmt.Sprintf("(W: %d, H: %d)", d.W, d.H)
}

func (d Dims) Grow(dw, dh int) Dims {
	return Dims{d.W + dw, d.H + dh}
}

func (d Dims) Shrink(dw, dh int) Dims {
	return Dims{d.W - dw, d.H - dh}
}

// Returns a rect with dimensions d, positioned at (0, 0)
func (d Dims) Bounds() Rect {
	return Rect{ZERO_COORD, d}
}
