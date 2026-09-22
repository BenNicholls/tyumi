package util

import (
	"math"
	"testing"
)

func TestPow(t *testing.T) {
	tests := []struct {
		value, exponent, want int
	}{
		{2, 0, 1},
		{2, 3, 8},
		{5, 1, 5},
		{-2, 3, -8},
		{0, 5, 0},
	}

	for _, tt := range tests {
		if got := Pow(tt.value, tt.exponent); got != tt.want {
			t.Errorf("Pow(%d, %d) = %d, want %d", tt.value, tt.exponent, got, tt.want)
		}
	}
}

func TestAbs(t *testing.T) {
	if got := Abs(-5); got != 5 {
		t.Errorf("Abs(-5) = %d, want 5", got)
	}
	if got := Abs(5); got != 5 {
		t.Errorf("Abs(5) = %d, want 5", got)
	}
	if got := Abs(0); got != 0 {
		t.Errorf("Abs(0) = %d, want 0", got)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		value, min, max, want int
	}{
		{5, 0, 10, 5},
		{-5, 0, 10, 0},
		{15, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
		{5, 5, 5, 5},   // min == max
		{5, 10, 0, 5},  // min > max, gets swapped
		{-5, 10, 0, 0}, // min > max, gets swapped, clamps to swapped min
	}

	for _, tt := range tests {
		if got := Clamp(tt.value, tt.min, tt.max); got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.value, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestCycleClamp(t *testing.T) {
	tests := []struct {
		value, min, max, want int
	}{
		{5, 0, 10, 5},
		{11, 0, 10, 0},
		{-1, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
		{5, 5, 5, 5}, // min == max
		{12, 0, 10, 1},
		{-2, 0, 10, 9},
	}

	for _, tt := range tests {
		if got := CycleClamp(tt.value, tt.min, tt.max); got != tt.want {
			t.Errorf("CycleClamp(%d, %d, %d) = %d, want %d", tt.value, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestCycleClampWithOverflow(t *testing.T) {
	tests := []struct {
		value, min, max, wantVal, wantOverflow int
	}{
		{5, 0, 10, 5, 0},
		{12, 0, 10, 1, 1},
		{23, 0, 10, 1, 2},
		{-1, 0, 10, 10, -1},
		{-12, 0, 10, 10, -2},
	}

	for _, tt := range tests {
		gotVal, gotOverflow := CycleClampWithOverflow(tt.value, tt.min, tt.max)
		if gotVal != tt.wantVal || gotOverflow != tt.wantOverflow {
			t.Errorf("CycleClampWithOverflow(%d, %d, %d) = (%d, %d), want (%d, %d)", tt.value, tt.min, tt.max, gotVal, gotOverflow, tt.wantVal, tt.wantOverflow)
		}
	}
}

func TestRoundFloatToInt(t *testing.T) {
	tests := []struct {
		in   float64
		want int
	}{
		{1.4, 1},
		{1.5, 2},
		{1.6, 2},
		{-1.4, -1},
		{-1.5, -2},
		{-1.6, -2},
		{0, 0},
	}

	for _, tt := range tests {
		if got := RoundFloatToInt(tt.in); got != tt.want {
			t.Errorf("RoundFloatToInt(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		start, end, val, steps, want int
	}{
		{0, 10, 0, 10, 0},
		{0, 10, 10, 10, 10},
		{0, 10, 5, 10, 5},
		{10, 0, 5, 10, 5},
	}

	for _, tt := range tests {
		if got := Lerp(tt.start, tt.end, tt.val, tt.steps); got != tt.want {
			t.Errorf("Lerp(%d, %d, %d, %d) = %d, want %d", tt.start, tt.end, tt.val, tt.steps, got, tt.want)
		}
	}

	if got := Lerp(0.0, 10.0, 5, 10); math.Abs(got-5.0) > 0.0001 {
		t.Errorf("Lerp float = %v, want 5.0", got)
	}
}
