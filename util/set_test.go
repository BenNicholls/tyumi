package util

import "testing"

func TestSetAddContains(t *testing.T) {
	var s Set[int]

	if s.Contains(1) {
		t.Error("empty set should not contain 1")
	}

	s.Add(1, 2, 3)
	if !s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
		t.Error("set should contain added elements")
	}

	if s.Count() != 3 {
		t.Errorf("Count() = %d, want 3", s.Count())
	}

	s.Add(1) // duplicate add is a no-op
	if s.Count() != 3 {
		t.Errorf("Count() after duplicate add = %d, want 3", s.Count())
	}

	s.Add() // empty add is a no-op, should not panic even on nil map
	if s.Count() != 3 {
		t.Errorf("Count() after empty add = %d, want 3", s.Count())
	}
}

func TestSetAddOnNil(t *testing.T) {
	var s Set[int]
	s.Add() // should not initialize the map or panic
	if s.Count() != 0 {
		t.Errorf("Count() = %d, want 0", s.Count())
	}
}

func TestSetRemove(t *testing.T) {
	var s Set[int]
	s.Add(1, 2, 3)
	s.Remove(2)

	if s.Contains(2) {
		t.Error("set should not contain removed element")
	}
	if s.Count() != 2 {
		t.Errorf("Count() = %d, want 2", s.Count())
	}

	var empty Set[int]
	empty.Remove(1) // should not panic on nil map
}

func TestSetRemoveAll(t *testing.T) {
	var s Set[int]
	s.Add(1, 2, 3)
	s.RemoveAll()

	if s.Count() != 0 {
		t.Errorf("Count() = %d, want 0", s.Count())
	}

	var empty Set[int]
	empty.RemoveAll() // should not panic
}

func TestSetContainsAllAny(t *testing.T) {
	var s Set[int]
	s.Add(1, 2, 3)

	if !s.ContainsAll(1, 2) {
		t.Error("ContainsAll(1, 2) should be true")
	}
	if s.ContainsAll(1, 9) {
		t.Error("ContainsAll(1, 9) should be false")
	}
	if s.ContainsAll(1, 2, 3, 4) {
		t.Error("ContainsAll with more elements than set has should be false")
	}

	if !s.ContainsAny(9, 2) {
		t.Error("ContainsAny(9, 2) should be true")
	}
	if s.ContainsAny(8, 9) {
		t.Error("ContainsAny(8, 9) should be false")
	}

	var empty Set[int]
	if empty.ContainsAll(1) {
		t.Error("empty set ContainsAll should be false")
	}
	if empty.ContainsAny(1) {
		t.Error("empty set ContainsAny should be false")
	}
}

func TestSetEquals(t *testing.T) {
	var s1, s2 Set[int]
	s1.Add(1, 2, 3)
	s2.Add(3, 2, 1)

	if !s1.Equals(s2) {
		t.Error("sets with same elements in different order should be equal")
	}

	s2.Add(4)
	if s1.Equals(s2) {
		t.Error("sets with different elements should not be equal")
	}
}

func TestSetIntersectionUnionDifference(t *testing.T) {
	var s1, s2 Set[int]
	s1.Add(1, 2, 3)
	s2.Add(2, 3, 4)

	intersection := s1.Intersection(s2)
	if !intersection.Equals(setOf(2, 3)) {
		t.Errorf("Intersection = %v elements, want {2, 3}", intersection.Count())
	}

	union := s1.Union(s2)
	if !union.Equals(setOf(1, 2, 3, 4)) {
		t.Error("Union does not match expected {1, 2, 3, 4}")
	}

	difference := s1.Difference(s2)
	if !difference.Equals(setOf(1)) {
		t.Error("Difference does not match expected {1}")
	}
}

func TestSetIntersectionUnionEmpty(t *testing.T) {
	var s1, empty Set[int]
	s1.Add(1, 2, 3)

	if s1.Intersection(empty).Count() != 0 {
		t.Error("Intersection with empty set should be empty")
	}

	if !s1.Union(empty).Equals(s1) {
		t.Error("Union with empty set should equal original set")
	}

	if !s1.Difference(empty).Equals(s1) {
		t.Error("Difference with empty set should equal original set")
	}
}

func TestSetAddSetRemoveSet(t *testing.T) {
	var s1, s2 Set[int]
	s1.Add(1, 2)
	s2.Add(2, 3)

	s1.AddSet(s2)
	if !s1.Equals(setOf(1, 2, 3)) {
		t.Error("AddSet did not merge sets correctly")
	}

	s1.RemoveSet(s2)
	if !s1.Equals(setOf(1)) {
		t.Error("RemoveSet did not remove elements correctly")
	}
}

func TestSetEachElement(t *testing.T) {
	var s Set[int]
	s.Add(1, 2, 3)

	found := make(map[int]bool)
	for elem := range s.EachElement() {
		found[elem] = true
	}

	if len(found) != 3 || !found[1] || !found[2] || !found[3] {
		t.Errorf("EachElement did not iterate all elements, got %v", found)
	}
}

func TestSetPickOne(t *testing.T) {
	var s Set[int]
	s.Add(1, 2, 3)

	for range 20 {
		picked := s.PickOne()
		if !s.Contains(picked) {
			t.Errorf("PickOne() = %d, not in set", picked)
		}
	}
}

func TestOrderedSetOrderPreserved(t *testing.T) {
	var os OrderedSet[int]
	os.Add(3, 1, 2)

	want := []int{3, 1, 2}
	for i, wantElem := range want {
		if got := os.At(i); got != wantElem {
			t.Errorf("os.At(%d) = %d, want %d", i, got, wantElem)
		}
	}

	if os.Count() != 3 {
		t.Errorf("Count() = %d, want 3", os.Count())
	}
}

func TestOrderedSetAddDuplicateNoOp(t *testing.T) {
	var os OrderedSet[int]
	os.Add(1, 2, 3)
	os.Add(2) // duplicate, should not change order/count

	if os.Count() != 3 {
		t.Errorf("Count() = %d, want 3", os.Count())
	}
	if got := os.At(1); got != 2 {
		t.Errorf("os.At(1) = %d, want 2", got)
	}
}

func TestOrderedSetRemove(t *testing.T) {
	var os OrderedSet[int]
	os.Add(1, 2, 3)
	os.Remove(2)

	if os.Contains(2) {
		t.Error("expected 2 to be removed")
	}
	if os.Count() != 2 {
		t.Errorf("Count() = %d, want 2", os.Count())
	}
	if got := os.At(1); got != 3 {
		t.Errorf("os.At(1) after removal = %d, want 3", got)
	}
}

func TestOrderedSetRemoveAt(t *testing.T) {
	var os OrderedSet[int]
	os.Add(10, 20, 30)
	os.RemoveAt(1)

	if os.Contains(20) {
		t.Error("expected element at index 1 (20) to be removed")
	}
	if os.Count() != 2 {
		t.Errorf("Count() = %d, want 2", os.Count())
	}
}

func TestOrderedSetRemoveAtPanicsOutOfRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for out of range RemoveAt")
		}
	}()

	var os OrderedSet[int]
	os.Add(1)
	os.RemoveAt(5)
}

func TestOrderedSetAtPanicsOutOfRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for out of range At")
		}
	}()

	var os OrderedSet[int]
	os.Add(1)
	os.At(5)
}

func TestOrderedSetRemoveFunc(t *testing.T) {
	var os OrderedSet[int]
	os.Add(1, 2, 3, 4, 5)

	os.RemoveFunc(func(e int) bool { return e%2 == 0 })

	if os.Count() != 3 {
		t.Errorf("Count() = %d, want 3", os.Count())
	}
	if os.Contains(2) || os.Contains(4) {
		t.Error("expected even elements to be removed")
	}
	if !os.Contains(1) || !os.Contains(3) || !os.Contains(5) {
		t.Error("expected odd elements to remain")
	}
}

func TestOrderedSetRemoveAll(t *testing.T) {
	var os OrderedSet[int]
	os.Add(1, 2, 3)
	os.RemoveAll()

	if os.Count() != 0 {
		t.Errorf("Count() = %d, want 0", os.Count())
	}
}

func TestOrderedSetEachElement(t *testing.T) {
	var os OrderedSet[int]
	os.Add(5, 6, 7)

	for i, elem := range os.EachElement() {
		if os.At(i) != elem {
			t.Errorf("EachElement mismatch at index %d: got %d, want %d", i, elem, os.At(i))
		}
	}
}

func TestOrderedSetPickOne(t *testing.T) {
	var os OrderedSet[int]
	os.Add(1, 2, 3)

	for range 20 {
		picked := os.PickOne()
		if !os.Contains(picked) {
			t.Errorf("PickOne() = %d, not in set", picked)
		}
	}
}

func setOf(elems ...int) Set[int] {
	var s Set[int]
	s.Add(elems...)
	return s
}
