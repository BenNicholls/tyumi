// utility functions for math. the Go standard library has some of these as well, but in a lot of situations they are too
// generic and weighed down by unnecessary checks and the like. the functions here are expected to be run in tight loops,
// so we strip out unneeded cruft. other things are just useful functions that I made up as I needed them.
package util

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Pow is an integer power function. Doesn't ~~do~~ negative exponents. Totally does 0 though.
func Pow(val, exp int) int {
	v := 1
	for i := 0; i < exp; i++ {
		v = v * val
	}
	return v
}

// Abs returns the absolute value of val
func Abs[T constraints.Signed](val T) T {
	if val < 0 {
		return val * (-1)
	}
	return val
}

// Max returns the max of two numbers. Duh.
func Max[T constraints.Ordered](i, j T) T {
	if i < j {
		return j
	}
	return i
}

// Min is the opposite of max.
func Min[T constraints.Ordered](i, j T) T {
	if i > j {
		return j
	}
	return i
}

// Clamp checks if min <= val <= max. If val < min, returns min. If val > max, returns max. Otherwise returns val.
func Clamp[T constraints.Ordered](val, min, max T) T {
	if min > max {
		min, max = max, min
	}

	if val <= min {
		return min
	} else if val >= max {
		return max
	}
	return val
}

// CycleClamp is like clamp but instead of clamping at the endpoints, it overflows/underflows back to the other side of
// the range. This range of the function is INCLUSIVE of min and max, so min <= val <= max.
func CycleClamp(val, min, max int) int {
	clamped_val, _ := CycleClampWithOverflow(val, min, max)
	return clamped_val
}

// CycleClampWithOverflow is like CycleClamp but also returns the number of overflow cycles. negative for underflow, 
// 0 for none, positive for overflow. 
//NOTE: this function kind of doesn't work as expected since it is inclusive at the end points.
//THINK: should this just be removed? it's pretty niche.
func CycleClampWithOverflow(val, min, max int) (int, int) {
	if min > max {
		min, max = max, min
	}

	if val < min {
		r := max - min + 1
		underflows := (min-val-1)/r + 1
		return val + r*underflows, -underflows
	} else if val > max {
		r := max - min + 1
		overflows := (val-max-1)/r + 1
		return val - r*overflows, overflows
	}
	return val, 0
}

// RoundFloatToInt rounds a float to an int in the way you'd expect. It's the way I expect anyways.
func RoundFloatToInt(f float64) int {
	return int(f + math.Copysign(0.5, f))
}

// Lerp linearly interpolates a range (start-end) over (steps) intervals, and returns the (val)th step.
func Lerp[T constraints.Integer|constraints.Float](start, end T, val, steps int) T {
	if val >= steps {
		return end
	} else if val <= 0 {
		return start
	}

	stepVal := float64(end-start) / float64(steps)
	return T(float64(start) + stepVal*float64(val))
}
