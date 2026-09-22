package util

import (
	"slices"
	"testing"
)

func TestPickOne(t *testing.T) {
	s := []int{42}
	if got := PickOne(s); got != 42 {
		t.Errorf("PickOne single element = %d, want 42", got)
	}

	s2 := []int{1, 2, 3, 4, 5}
	for range 50 {
		got := PickOne(s2)
		if !slices.Contains(s2, got) {
			t.Errorf("PickOne(%v) = %d, not in slice", s2, got)
		}
	}
}

func TestDeleteElement(t *testing.T) {
	tests := []struct {
		name    string
		slice   []int
		element int
		want    []int
	}{
		{"middle element", []int{1, 2, 3, 4, 5}, 3, []int{1, 2, 4, 5}},
		{"first element", []int{1, 2, 3}, 1, []int{2, 3}},
		{"last element", []int{1, 2, 3}, 3, []int{1, 2}},
		{"not present", []int{1, 2, 3}, 9, []int{1, 2, 3}},
		{"empty slice", []int{}, 1, []int{}},
		{"duplicate elements removed", []int{1, 2, 1, 3, 1}, 1, []int{2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeleteElement(tt.slice, tt.element)
			if !slices.Equal(got, tt.want) {
				t.Errorf("DeleteElement(%v, %d) = %v, want %v", tt.slice, tt.element, got, tt.want)
			}
		})
	}
}

func TestSetAll(t *testing.T) {
	s := make([]int, 5)
	SetAll(s, 7)

	for i, v := range s {
		if v != 7 {
			t.Errorf("s[%d] = %d, want 7", i, v)
		}
	}
}

func TestSetAllEmpty(t *testing.T) {
	s := []int{}
	SetAll(s, 7) // should not panic
	if len(s) != 0 {
		t.Errorf("expected empty slice to remain empty")
	}
}

func TestOrAll(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  int
	}{
		{"basic or", []int{0b0001, 0b0010, 0b0100}, 0b0111},
		{"empty", []int{}, 0},
		{"single", []int{5}, 5},
		{"overlapping bits", []int{0b0011, 0b0110}, 0b0111},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrAll(tt.slice); got != tt.want {
				t.Errorf("OrAll(%v) = %d, want %d", tt.slice, got, tt.want)
			}
		})
	}
}
