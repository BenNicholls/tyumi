package util

import (
	"iter"

	"golang.org/x/exp/constraints"
)

// RingBuffer is an append-only buffer of objects. Use Init() to specify a size for the buffer. When trying to append
// past the limit, it will instead overwrite the oldest element in the buffer.
type RingBuffer[T any] struct {
	buffer    []T
	nextIndex int
}

func (rb *RingBuffer[T]) Init(capacity int) {
	rb.buffer = make([]T, 0, capacity)
	rb.nextIndex = 0
}

func (rb *RingBuffer[T]) Append(item T) {
	if len(rb.buffer) == cap(rb.buffer) {
		rb.buffer[rb.nextIndex] = item
	} else {
		rb.buffer = append(rb.buffer, item)
	}
	rb.nextIndex = CycleClamp(rb.nextIndex+1, 0, cap(rb.buffer)-1)
}

func (rb *RingBuffer[T]) Clear() {
	clear(rb.buffer)
	rb.nextIndex = 0
}

func (rb RingBuffer[T]) Len() int {
	return len(rb.buffer)
}

func (rb RingBuffer[T]) EachElement() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range rb.buffer {
			if !yield(rb.buffer[CycleClamp(i+rb.nextIndex, 0, len(rb.buffer))]) {
				return
			}
		}
	}
}

// SummedRingBuffer is a ring buffer that keeps track of the sum of its elements.
type SummedRingBuffer[T constraints.Integer | constraints.Float] struct {
	RingBuffer[T]

	sum T
}

func (srb *SummedRingBuffer[T]) Init(capacity int) {
	srb.RingBuffer.Init(capacity)
	srb.sum = T(0)
}

func (srb *SummedRingBuffer[T]) Append(item T) {

	if len(srb.buffer) == cap(srb.buffer) {
		srb.sum -= srb.buffer[srb.nextIndex]
	}

	srb.sum += item

	srb.RingBuffer.Append(item)
}

func (srb *SummedRingBuffer[T]) Clear() {
	srb.RingBuffer.Clear()
	srb.sum = 0
}

func (srb SummedRingBuffer[T]) Sum() T {
	return srb.sum
}

func (srb SummedRingBuffer[T]) Avg() T {
	return srb.sum / T(srb.Len())
}

func (srb SummedRingBuffer[T]) AvgF() float64 {
	return float64(srb.sum) / float64(srb.Len())
}
