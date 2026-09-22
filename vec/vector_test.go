package vec

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}

func TestVec2fAddSub(t *testing.T) {
	v1 := Vec2f{1, 2}
	v2 := Vec2f{3, 4}

	if got := v1.Add(v2); got != (Vec2f{4, 6}) {
		t.Errorf("Add = %v, want {4 6}", got)
	}
	if got := v1.Sub(v2); got != (Vec2f{-2, -2}) {
		t.Errorf("Sub = %v, want {-2 -2}", got)
	}
}

func TestVec2fSetMod(t *testing.T) {
	v := Vec2f{1, 1}
	v.Set(5, 6)
	if v != (Vec2f{5, 6}) {
		t.Errorf("Set = %v, want {5 6}", v)
	}

	v.Mod(1, -1)
	if v != (Vec2f{6, 5}) {
		t.Errorf("Mod = %v, want {6 5}", v)
	}
}

func TestVec2fMag(t *testing.T) {
	v := Vec2f{3, 4}
	if got := v.Mag(); !almostEqual(got, 5) {
		t.Errorf("Mag() = %v, want 5", got)
	}
}

func TestVec2fNonZero(t *testing.T) {
	if (Vec2f{0, 0}).NonZero() {
		t.Error("zero vector should not be NonZero")
	}
	if !(Vec2f{1, 0}).NonZero() {
		t.Error("{1,0} should be NonZero")
	}
}

func TestVec2fToVec2i(t *testing.T) {
	got := Vec2f{1.6, -1.6}.ToVec2i()
	want := Vec2i{2, -2}
	if got != want {
		t.Errorf("ToVec2i() = %v, want %v", got, want)
	}
}

func TestVec2fToPolar(t *testing.T) {
	p := Vec2f{1, 0}.ToPolar()
	if !almostEqual(p.R, 1) || !almostEqual(p.Phi, 0) {
		t.Errorf("ToPolar() = %v, want {1 0}", p)
	}
}

func TestVec2PolarToRect(t *testing.T) {
	p := Vec2Polar{1, 0}
	got := p.ToRect()
	if !almostEqual(got.X, 1) || !almostEqual(got.Y, 0) {
		t.Errorf("ToRect() = %v, want {1 0}", got)
	}
}

func TestVec2PolarSetGet(t *testing.T) {
	var p Vec2Polar
	p.Set(2, 1.5)
	r, phi := p.Get()
	if r != 2 || phi != 1.5 {
		t.Errorf("Get() = (%v, %v), want (2, 1.5)", r, phi)
	}
}

func TestVec2PolarPos(t *testing.T) {
	p := Vec2Polar{-1, 0}
	p.Pos()
	if p.R < 0 {
		t.Errorf("Pos() left R negative: %v", p.R)
	}

	p2 := Vec2Polar{1, -1}
	p2.Pos()
	if p2.Phi < 0 || p2.Phi > 2*math.Pi {
		t.Errorf("Pos() left Phi out of range [0, 2pi]: %v", p2.Phi)
	}
}

func TestVec2PolarAngularDistance(t *testing.T) {
	v1 := Vec2Polar{1, 0}
	v2 := Vec2Polar{1, math.Pi / 2}

	got := v1.AngularDistance(v2)
	if !almostEqual(got, math.Pi/2) {
		t.Errorf("AngularDistance = %v, want pi/2", got)
	}
}

func TestVec2PolarAngularDistanceWraps(t *testing.T) {
	v1 := Vec2Polar{1, 0.1}
	v2 := Vec2Polar{1, 2*math.Pi - 0.1}

	got := v1.AngularDistance(v2)
	// the shortest path should be negative (clockwise), not almost 2*pi
	if got > 0 {
		t.Errorf("AngularDistance = %v, want a small negative value (shortest path)", got)
	}
}

func TestVec2iAddSub(t *testing.T) {
	v1 := Vec2i{1, 2}
	v2 := Vec2i{3, 4}

	if got := v1.Add(v2); got != (Vec2i{4, 6}) {
		t.Errorf("Add = %v, want {4 6}", got)
	}
	if got := v1.Sub(v2); got != (Vec2i{-2, -2}) {
		t.Errorf("Sub = %v, want {-2 -2}", got)
	}
}

func TestVec2iScale(t *testing.T) {
	got := Vec2i{2, 3}.Scale(3)
	want := Vec2i{6, 9}
	if got != want {
		t.Errorf("Scale(3) = %v, want %v", got, want)
	}
}

func TestVec2iMag(t *testing.T) {
	got := Vec2i{3, 4}.Mag()
	if !almostEqual(got, 5) {
		t.Errorf("Mag() = %v, want 5", got)
	}
}

func TestVec2iNonZero(t *testing.T) {
	if (Vec2i{0, 0}).NonZero() {
		t.Error("zero vector should not be NonZero")
	}
	if !(Vec2i{0, 1}).NonZero() {
		t.Error("{0,1} should be NonZero")
	}
}

func TestVec2iToVec2f(t *testing.T) {
	got := Vec2i{3, -4}.ToVec2f()
	want := Vec2f{3, -4}
	if got != want {
		t.Errorf("ToVec2f() = %v, want %v", got, want)
	}
}
