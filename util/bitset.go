package util

import "math/bits"

type Bitset struct {
	bits     []uint64
	capacity int
}

// Initializes the bitset. Capacity is the number of bits in the bitset.
func (bs *Bitset) Init(capacity int) {
	bs.bits = make([]uint64, capacity/64+1)
	bs.capacity = capacity
}

// Sets the bit at index idx to 1.
func (bs *Bitset) Set(idx int) {
	if idx >= bs.capacity || idx < 0 {
		return
	}

	bs.bits[idx/64] |= (uint64(0x1) << (idx % 64))
}

// Sets all bits in the bitset.
func (bs *Bitset) SetAll() {
	for i := range bs.bits {
		bs.bits[i] = 0xFFFFFFFFFFFFFFFF
	}
}

// Sets the bit at index idx to 0.
func (bs *Bitset) Unset(idx int) {
	if idx >= bs.capacity || idx < 0 {
		return
	}

	bs.bits[idx/64] &^= (uint64(0x1) << (idx % 64))
}

// Sets the bit at index idx to the value indicated (0 if false, 1 if true)
func (bs *Bitset) SetTo(idx int, value bool) {
	if value {
		bs.Set(idx)
	} else {
		bs.Unset(idx)
	}
}

// Gets the bit at index idx.
func (bs Bitset) Get(idx int) bool {
	if idx >= bs.capacity || idx < 0 {
		return false
	}

	bit := bs.bits[idx/64]
	mask := uint64(0x1) << (idx % 64)

	return (bit & mask) != 0
}

// Clears the bitset (sets all bits to zero).
func (bs *Bitset) Clear() {
	clear(bs.bits)
}

// IsEmpty is true if all bits in the set are 0.
func (bs Bitset) IsEmpty() bool {
	for _, bit := range bs.bits {
		if bit != 0 {
			return false
		}
	}

	return true
}

// Count reports the total number of set bits.
func (bs Bitset) Count() (total int) {
	for _, bit := range bs.bits {
		total += bits.OnesCount64(bit)
	}

	return
}
