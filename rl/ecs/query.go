package ecs

import (
	"iter"

	"github.com/bennicholls/tyumi/log"
)

// EachComponent is an iterator that iterates over all active components of a certain type. The 2nd returned iterator
// value is the component's entity.
// WARNING: do NOT remove components of this type while iterating!
func EachComponent[T componentType]() iter.Seq2[*T, Entity] {
	cache := getComponentCache[T]()
	return func(yield func(*T, Entity) bool) {
		for i := range cache.components {
			if !yield(&cache.components[i], cache.components[i].GetEntity()) {
				return
			}
		}
	}
}

// EachEntityWith is an iterator that returns all of the entities with a certain component.
func EachEntityWith[T componentType]() iter.Seq[Entity] {
	cache := getComponentCache[T]()
	return func(yield func(Entity) bool) {
		for i := range cache.components {
			if !yield(cache.components[i].GetEntity()) {
				return
			}
		}
	}
}

func EachEntityWithTag(tag Tag) iter.Seq[Entity] {
	if !tag.isValid() {
		log.Debug("Invalid tag!")
		return func(yield func(Entity) bool) { return }
	}

	return tagCaches[tag].EachElement()
}
