// utility functions for generating random numbers, sequences, etc.
package util

import (
	"math/rand"
	"slices"

	"golang.org/x/exp/constraints"
)

// PickOne picks a random element from the slice and returns it.
func PickOne[S ~[]E, E any](slice S) E {
	return slice[rand.Intn(len(slice))]
}

func DeleteElement[S ~[]E, E comparable](slice S, element E) S {
	return slices.DeleteFunc(slice, func(e E) bool {
		return e == element
	})
}

// SetAll sets all elements in the slice S to the provided value.
func SetAll[S ~[]E, E any](slice S, value E) {
	for i := range slice {
		slice[i] = value
	}
}

// OrAll return all elements of the slice or'd together.
func OrAll[S ~[]E, E constraints.Integer](slice S) (val E) {
	for _, elem := range slice {
		val |= elem
	}

	return
}
