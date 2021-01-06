//utility functions for math. the Go standard library has some of these as well, but in a lot of situations they are too
//generic and weighed down by unnecessary checks and the like. the functions here are expected to be run in tight loops,
//so we strip out unneeded cruft. other things are just useful functions that I made up as I needed them.
package util

import "math"

//Pow is an integer power function. Doesn't ~~do~~ negative exponents. Totally does 0 though.
func Pow(val, exp int) int {
	v := 1
	for i := 0; i < exp; i++ {
		v = v * val
	}
	return v
}

//Abs returns the absolute value of val
func Abs(val int) int {
	if val < 0 {
		return val * (-1)
	}
	return val
}

//Max returns the max of two integers. Duh.
func Max(i, j int) int {
	if i < j {
		return j
	}
	return i
}

func MaxF(i, j float64) float64 {
	if i < j {
		return j
	}
	return i
}

//Min is the opposite of max.
func Min(i, j int) int {
	if i > j {
		return j
	}
	return i
}

func MinF(i, j float64) float64 {
	if i > j {
		return j
	}
	return i
}

//Clamp checks if min <= val <= max. If val < min, returns min. If val > max, returns max. Otherwise returns val.
func Clamp(val, min, max int) int {
	if val <= min {
		return min
	} else if val >= max {
		return max
	}
	return val
}

//ModularClamp is like clamp but instead of clamping at the endpoints, it overflows/underflows back to the other side of
//the range. The second return param is the number of overflow cycles. negative for underflow, 0 for none, positive for 
//overflow. This kind of function probably has an actual name but hell if I know what it is.
func ModularClamp(val, min, max int) (int, int) {
	if min > max {
		//if someone foolishly puts their min higher than max, swap
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

//RoundFloatToInt rounds a float to an int in the way you'd expect. It's the way I expect anyways.
func RoundFloatToInt(f float64) int {
	return int(f + math.Copysign(0.5, f))
}

//Lerp linearly interpolates a range (min-max) over (steps) intervals, and returns the (val)th step. Currently does this
//via a conversion to float64, so there might be some rounding errors in here I don't know about.
func Lerp(min, max, val, steps int) int {
	if val >= steps {
		return max
	} else if val <= 0 {
		return min
	}

	stepVal := float64(max-min) / float64(steps)
	return int(float64(min) + stepVal*float64(val))
}
