// utility functions for generating random numbers, sequences, etc.
package util

import "math/rand"

// Picks a random element from the slice and returns it.
func PickOne[S ~[]E, E any](slice S) E {
	return slice[rand.Intn(len(slice))]
}
