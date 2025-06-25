package rl

import (
	"fmt"

	"github.com/bennicholls/tyumi/util"
	"golang.org/x/exp/constraints"
)

type Stat[T constraints.Float | constraints.Integer] struct {
	value    T
	min, max T
}

func NewStat[T constraints.Float | constraints.Integer](value, min, max T) Stat[T] {
	return Stat[T]{
		value: util.Clamp(value, min, max),
		min:   min,
		max:   max,
	}
}

// Returns a stat with min 0, max = value.
func NewBasicStat[T constraints.Float | constraints.Integer](value T) Stat[T] {
	return Stat[T]{
		value: value,
		max:   value,
	}
}

func (s Stat[T]) Get() T {
	return s.value
}

func (s *Stat[T]) Set(v T) {
	s.value = util.Clamp(v, s.min, s.max)
}

func (s *Stat[T]) Mod(dv T) {
	s.Set(s.value + dv)
}

func (s Stat[T]) Max() T {
	return s.max
}

func (s Stat[T]) Min() T {
	return s.min
}

// Sets a new minimum. If this would make min > max, does nothing.
func (s *Stat[T]) SetMin(m T) {
	if m > s.max {
		return
	}

	s.min = m
	s.value = util.Clamp(s.value, s.min, s.max)
}

// Sets a new maximum. If this would make max < min, does nothing.
func (s *Stat[T]) SetMax(m T) {
	if m < s.min {
		return
	}

	s.max = m
	s.value = util.Clamp(s.value, s.min, s.max)
}

// Modifies the minimum. Takes a delta. Follows same rules as SetMin().
func (s *Stat[T]) ModMin(d T) {
	s.SetMin(s.min + d)
}

// Modifies the maximum. Takes a delta. Follows same rules as SetMax().
func (s *Stat[T]) ModMax(d T) {
	s.SetMax(s.max + d)
}

func (s Stat[T]) IsMax() bool {
	return s.value == s.max
}

func (s Stat[T]) IsMin() bool {
	return s.value == s.min
}

// returns a % (0-100) for the stat[T]. If min == val == max, returns 0.
func (s Stat[T]) GetPct() int {
	if s.min == s.max {
		return 0
	}

	return int(100 * (float32(s.value-s.min) / float32(s.max-s.min)))
}

func (s Stat[T]) String() string {
	return fmt.Sprintf("%v/%v", s.value, s.max)
}
