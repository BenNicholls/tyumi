// utility functions for generating random numbers, sequences, etc.
package util

import "math/rand"

// RandomDirection generates a tuple of cartesian directions (cannot be 0,0)
func RandomDirection() (int, int) {
	for {
		dx, dy := rand.Intn(3)-1, rand.Intn(3)-1
		if dx != 0 || dy != 0 {
			return dx, dy
		}
	}
}

// GenerateCoord generates a random (x,y) pair within a box defined by (x, y, w, h)
func GenerateCoord(x, y, w, h int) (int, int) {
	return rand.Intn(w) + x, rand.Intn(h) + y
}
