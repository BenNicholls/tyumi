package col

import "testing"

func TestBlendMultiply(t *testing.T) {
	c1 := MakeOpaque(255, 255, 255)
	c2 := MakeOpaque(100, 150, 200)

	// multiplying by full white should return the other colour unchanged
	got := Blend(c1, c2, BLEND_MULTIPLY)
	if got.R() != 100 || got.G() != 150 || got.B() != 200 {
		t.Errorf("Blend(white, c2, MULTIPLY) = %v, want RGB (100, 150, 200)", got)
	}
}

func TestBlendMultiplyBlack(t *testing.T) {
	c1 := MakeOpaque(0, 0, 0)
	c2 := MakeOpaque(100, 150, 200)

	got := Blend(c1, c2, BLEND_MULTIPLY)
	if got.R() != 0 || got.G() != 0 || got.B() != 0 {
		t.Errorf("Blend(black, c2, MULTIPLY) = %v, want RGB (0, 0, 0)", got)
	}
}

func TestBlendScreen(t *testing.T) {
	c1 := MakeOpaque(0, 0, 0)
	c2 := MakeOpaque(100, 150, 200)

	// screen blending with black should return the other colour unchanged
	got := Blend(c1, c2, BLEND_SCREEN)
	if got.R() != 100 || got.G() != 150 || got.B() != 200 {
		t.Errorf("Blend(black, c2, SCREEN) = %v, want RGB (100, 150, 200)", got)
	}
}

func TestBlendScreenWhite(t *testing.T) {
	c1 := MakeOpaque(255, 255, 255)
	c2 := MakeOpaque(100, 150, 200)

	got := Blend(c1, c2, BLEND_SCREEN)
	if got.R() != 255 || got.G() != 255 || got.B() != 255 {
		t.Errorf("Blend(white, c2, SCREEN) = %v, want RGB (255, 255, 255)", got)
	}
}
