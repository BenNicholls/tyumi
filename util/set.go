package util

import "iter"

// Set is a container that holds elements of a type E. Elements can be added and removed, and duplicate adds are no-ops.
type Set[E comparable] struct {
	elements map[E]bool
}

func (s Set[E]) Count() int {
	return len(s.elements)
}

func (s *Set[E]) Add(elems ...E) {
	if len(elems) == 0 {
		return
	}

	if s.elements == nil {
		s.elements = make(map[E]bool)
	}

	for _, elem := range elems {
		s.elements[elem] = true
	}
}

// Adds the elements of another set to the Set. Operationally this is the same as a Union, except it works in-place.
func (s *Set[E]) AddSet(s2 Set[E]) {
	for elem := range s2.elements {
		s.Add(elem)
	}
}

// Returns true if s and s2 have precisely the same elements.
func (s Set[E]) Equals(s2 Set[E]) bool {
	if s.Count() != s2.Count() {
		return false
	}

	for elem := range s.elements {
		if !s2.Contains(elem) {
			return false
		}
	}

	return true
}

func (s Set[E]) Contains(elem E) bool {
	_, ok := s.elements[elem]
	return ok
}

func (s Set[E]) ContainsAll(elems ...E) bool {
	if len(s.elements) == 0 || len(s.elements) < len(elems) {
		return false
	}

	for _, elem := range elems {
		if !s.Contains(elem) {
			return false
		}
	}

	return true
}

func (s Set[E]) ContainsAny(elems ...E) bool {
	if len(s.elements) == 0 {
		return false
	}

	for _, elem := range elems {
		if s.Contains(elem) {
			return true
		}
	}

	return false
}

func (s *Set[E]) Remove(elems ...E) {
	if len(s.elements) == 0 {
		return
	}

	for _, elem := range elems {
		delete(s.elements, elem)
	}
}

func (s *Set[E]) RemoveSet(s2 Set[E]) {
	if len(s.elements) == 0 || len(s2.elements) == 0 {
		return
	}

	for elem := range s2.elements {
		s.Remove(elem)
	}
}

func (s *Set[E]) RemoveAll() {
	if len(s.elements) == 0 {
		return
	}

	clear(s.elements)
}

func (s Set[E]) Intersection(s2 Set[E]) (intersection Set[E]) {
	if len(s.elements) == 0 || len(s2.elements) == 0 {
		return
	}

	for elem := range s.elements {
		if s2.Contains(elem) {
			intersection.Add(elem)
		}
	}

	return
}

func (s Set[E]) Union(s2 Set[E]) (union Set[E]) {
	if len(s.elements) == 0 && len(s2.elements) == 0 {
		return
	}

	for elem := range s.elements {
		union.Add(elem)
	}

	for elem := range s2.elements {
		union.Add(elem)
	}

	return
}

// Difference returns s - s2
func (s Set[E]) Difference(s2 Set[E]) (difference Set[E]) {
	if len(s.elements) == 0 {
		return
	}

	if len(s2.elements) == 0 {
		for elem := range s.EachElement() {
			difference.Add(elem)
		}
		return
	}

	for elem := range s.EachElement() {
		if !s2.Contains(elem) {
			difference.Add(elem)
		}
	}

	return
}

func (s Set[E]) EachElement() iter.Seq[E] {
	return func(yield func(E) bool) {
		for elem := range s.elements {
			if !yield(elem) {
				return
			}
		}
	}
}
