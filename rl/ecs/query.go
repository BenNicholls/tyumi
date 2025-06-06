package ecs

import "iter"

// EachComponent is an iterator that iterates over all active components of a certain type.
// WARNING: do NOT remove components of this type while iterating!
func EachComponent[T componentType]() iter.Seq[*T] {
	cache := getComponentCache[T]()
	return func(yield func(*T) bool) {
		for i := range cache.components {
			if !yield(&cache.components[i]) {
				return
			}
		}
	}
}
