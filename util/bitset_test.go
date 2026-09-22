package util

import "testing"

func TestBitsetSetGet(t *testing.T) {
	var bs Bitset
	bs.Init(100)

	if bs.Get(5) {
		t.Error("expected bit 5 to be unset initially")
	}

	bs.Set(5)
	if !bs.Get(5) {
		t.Error("expected bit 5 to be set")
	}

	bs.Unset(5)
	if bs.Get(5) {
		t.Error("expected bit 5 to be unset after Unset")
	}
}

func TestBitsetSetAcrossWords(t *testing.T) {
	var bs Bitset
	bs.Init(200)

	indices := []int{0, 63, 64, 65, 127, 128, 199}
	for _, i := range indices {
		bs.Set(i)
	}

	for _, i := range indices {
		if !bs.Get(i) {
			t.Errorf("expected bit %d to be set", i)
		}
	}

	if bs.Get(1) {
		t.Error("expected bit 1 to remain unset")
	}
}

func TestBitsetOutOfRange(t *testing.T) {
	var bs Bitset
	bs.Init(10)

	// should not panic
	bs.Set(-1)
	bs.Set(10)
	bs.Set(1000)
	bs.Unset(-1)
	bs.Unset(1000)

	if bs.Get(-1) || bs.Get(10) || bs.Get(1000) {
		t.Error("out of range Get should always return false")
	}
}

func TestBitsetSetTo(t *testing.T) {
	var bs Bitset
	bs.Init(10)

	bs.SetTo(3, true)
	if !bs.Get(3) {
		t.Error("expected bit 3 to be set via SetTo(true)")
	}

	bs.SetTo(3, false)
	if bs.Get(3) {
		t.Error("expected bit 3 to be unset via SetTo(false)")
	}
}

func TestBitsetSetAll(t *testing.T) {
	var bs Bitset
	bs.Init(10)

	bs.SetAll()
	for i := range 10 {
		if !bs.Get(i) {
			t.Errorf("expected bit %d to be set after SetAll", i)
		}
	}
}

func TestBitsetClear(t *testing.T) {
	var bs Bitset
	bs.Init(10)

	bs.SetAll()
	bs.Clear()

	for i := range 10 {
		if bs.Get(i) {
			t.Errorf("expected bit %d to be unset after Clear", i)
		}
	}

	if !bs.IsEmpty() {
		t.Error("expected bitset to be empty after Clear")
	}
}

func TestBitsetIsEmpty(t *testing.T) {
	var bs Bitset
	bs.Init(10)

	if !bs.IsEmpty() {
		t.Error("expected freshly initialized bitset to be empty")
	}

	bs.Set(5)
	if bs.IsEmpty() {
		t.Error("expected bitset to not be empty after Set")
	}
}

func TestBitsetCount(t *testing.T) {
	var bs Bitset
	bs.Init(100)

	if got := bs.Count(); got != 0 {
		t.Errorf("Count() = %d, want 0", got)
	}

	bs.Set(0)
	bs.Set(50)
	bs.Set(99)

	if got := bs.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}

	bs.Unset(50)
	if got := bs.Count(); got != 2 {
		t.Errorf("Count() = %d, want 2", got)
	}
}
