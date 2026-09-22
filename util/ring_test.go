package util

import "testing"

// BUG: RingBuffer.EachElement panics with an index-out-of-range whenever the buffer has been appended to but is not
// currently exactly full (i.e. len(buffer) is not a multiple of its capacity). This is because it wraps the
// iteration index with CycleClamp(i+nextIndex, 0, len(buffer)), which is INCLUSIVE of len(buffer) - a valid index is
// only 0..len(buffer)-1. This is a real, easily reproducible bug affecting the primary use case of a partially
// filled ring buffer, documented here rather than fixed since fixing it would change public API behaviour.
// TestRingBufferEachElementPanicsWhenNotFull below demonstrates it; the other iteration tests avoid triggering it by
// only calling EachElement when the buffer is exactly at capacity.

func TestRingBufferAppendLen(t *testing.T) {
	var rb RingBuffer[int]
	rb.Init(5)

	rb.Append(1)
	rb.Append(2)
	rb.Append(3)

	if rb.Len() != 3 {
		t.Errorf("Len() = %d, want 3", rb.Len())
	}
}

func TestRingBufferEachElementPanicsWhenNotFull(t *testing.T) {
	var rb RingBuffer[int]
	rb.Init(5)
	rb.Append(1)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected EachElement to panic on a partially-filled buffer (see BUG note above); if this no longer panics, the underlying bug may have been fixed and this test should be updated")
		}
	}()

	for range rb.EachElement() {
	}
}

func TestRingBufferEachElementWhenExactlyFull(t *testing.T) {
	var rb RingBuffer[int]
	rb.Init(3)

	rb.Append(1)
	rb.Append(2)
	rb.Append(3)

	var got []int
	for v := range rb.EachElement() {
		got = append(got, v)
	}

	want := []int{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("EachElement produced %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("EachElement()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestRingBufferOverwritesOldest(t *testing.T) {
	var rb RingBuffer[int]
	rb.Init(3)

	rb.Append(1)
	rb.Append(2)
	rb.Append(3)
	rb.Append(4) // overwrites 1
	rb.Append(5) // overwrites 2
	rb.Append(6) // overwrites 3, buffer now exactly full again: [4, 5, 6]

	if rb.Len() != 3 {
		t.Errorf("Len() = %d, want 3", rb.Len())
	}

	var got []int
	for v := range rb.EachElement() {
		got = append(got, v)
	}
	want := []int{4, 5, 6}
	if len(got) != len(want) {
		t.Fatalf("EachElement produced %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("EachElement()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

// BUG: Clear() calls the builtin clear() on the underlying slice, which zeroes the elements but does not change the
// slice's length, so Len() is unaffected by Clear(). Documenting current behaviour rather than fixing it, since
// fixing it would change public API behaviour.
func TestRingBufferClear(t *testing.T) {
	var rb RingBuffer[int]
	rb.Init(5)
	rb.Append(1)
	rb.Append(2)

	rb.Clear()

	if rb.Len() != 2 {
		t.Errorf("Len() after Clear() = %d, want 2 (Clear does not reset length, see BUG note above)", rb.Len())
	}
}

func TestSummedRingBufferSum(t *testing.T) {
	var srb SummedRingBuffer[int]
	srb.Init(3)

	srb.Append(10)
	srb.Append(20)
	srb.Append(30)

	if got := srb.Sum(); got != 60 {
		t.Errorf("Sum() = %d, want 60", got)
	}

	if got := srb.Avg(); got != 20 {
		t.Errorf("Avg() = %d, want 20", got)
	}
}

func TestSummedRingBufferOverwriteAdjustsSum(t *testing.T) {
	var srb SummedRingBuffer[int]
	srb.Init(3)

	srb.Append(10)
	srb.Append(20)
	srb.Append(30)
	srb.Append(40) // overwrites 10

	if got := srb.Sum(); got != 90 {
		t.Errorf("Sum() = %d, want 90 (20+30+40)", got)
	}
}

func TestSummedRingBufferAvgF(t *testing.T) {
	var srb SummedRingBuffer[int]
	srb.Init(4)

	srb.Append(1)
	srb.Append(2)
	srb.Append(3)

	if got := srb.AvgF(); got != 2.0 {
		t.Errorf("AvgF() = %v, want 2.0", got)
	}
}

// See the BUG note on TestRingBufferClear: Len() is not reset by Clear(), only Sum() is (SummedRingBuffer.Clear
// explicitly resets its own sum field).
func TestSummedRingBufferClear(t *testing.T) {
	var srb SummedRingBuffer[int]
	srb.Init(3)
	srb.Append(10)
	srb.Append(20)

	srb.Clear()

	if got := srb.Sum(); got != 0 {
		t.Errorf("Sum() after Clear() = %d, want 0", got)
	}
	if got := srb.Len(); got != 2 {
		t.Errorf("Len() after Clear() = %d, want 2 (Clear does not reset length)", got)
	}
}
